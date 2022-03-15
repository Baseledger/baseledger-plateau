const request = require('supertest');
const Web3 = require('web3')
const fs = require('fs');
const path = require("path");
const chai = require("chai");
const expect = chai.expect;

const shell = require('shelljs')

const host = 'localhost';
const node1_api_url = 'localhost:1317';
const node2_api_url = 'localhost:1318';
const node3_api_url = 'localhost:1319';

const baseledger_abi = JSON.parse(fs.readFileSync(path.join(__dirname, 'baseledger_abi.json')));

const sleep = (ms) => {
  return new Promise(resolve => setTimeout(resolve, ms));
};
const TEST_TIMEOUT = 30000;

// these tests are using 3 nodes and 3 orchestrators - attestations should be observed
describe('attestations observed', async function() {
  this.timeout(50000);
  before(() => {
    startTestNet();
  });

  after(() => {
    cleanTestNet();
  });

  it('should add/update validator staking power', async function() {
    this.timeout(TEST_TIMEOUT + 60000);
    const startEventNonces = await getEventNonces();
    console.log('start event nonces ', await getEventNonces());

    const orchestratorValidatorResponse = await request(node1_api_url).get('/Baseledger/baseledger/bridge/orchestrator_validator_address')
        .send().expect(200);

    // using first node validator to change staking
    const parsedOrchValResponse = JSON.parse(orchestratorValidatorResponse.text);
    const validatorAddress = parsedOrchValResponse.orchestratorValidatorAddress[0].validatorAddress;

    let validators = await request(node1_api_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    let parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens at start ', parsedResponse.validator.tokens);
            
    const web3 = new Web3('http://localhost:8545');
    let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");

    const accounts = await web3.eth.getAccounts()

    // check for payee
    const payeeExists = await contract.methods.payees(accounts[1]).call()

    // add or update payee with 50k ubt (8 decimals)
    const methodToExecute = payeeExists
      ? contract.methods.updatePayee(accounts[1], accounts[1], 5000000000000, validatorAddress)
      : contract.methods.addPayee(accounts[1], accounts[1], 5000000000000, validatorAddress)
    
    methodToExecute.send({
        from: accounts[0]
    }).then(console.log)

    // sleep to wait for attestation to be observed
    await sleep(20000);

    validators = await request(node1_api_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens after setting to 50k ', parsedResponse.validator.tokens);
    // tokens should be 50k * 10^6
    expect(parsedResponse.validator.tokens).to.be.equal("50000000000");

    // update payee, increase power to 80k
    contract.methods.updatePayee(accounts[1], accounts[1], 8000000000000, validatorAddress).send({
        from: accounts[0]
    }).then(console.log)

    // sleep to wait for attestation to be observed
    await sleep(20000);

    validators = await request(node1_api_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens after increasing to 80k ', parsedResponse.validator.tokens);
    expect(parsedResponse.validator.tokens).to.be.equal("80000000000");

    // update payee, decrease power to 70k
    contract.methods.updatePayee(accounts[1], accounts[1], 7000000000000, validatorAddress).send({
        from: accounts[0]
    }).then(console.log)

    // sleep to wait for attestation to be observed
    await sleep(20000);

    validators = await request(node1_api_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens after decreasing to 70k ', parsedResponse.validator.tokens);
    expect(parsedResponse.validator.tokens).to.be.equal("70000000000");

    // check if event nonces increased by 3
    const endEventNonces = await getEventNonces();
    console.log('end event nonces ', endEventNonces);
    startEventNonces.forEach((n, i) => {
      expect(n + 3).to.be.equal(endEventNonces[i]);
    });
  });

  it('should deposit ubt to baseledger account', async function () {
    this.timeout(TEST_TIMEOUT + 20000);

    const startEventNonces = await getEventNonces();
    console.log('start event nonces ', await getEventNonces());

    // random regular baseledger address
    const baseledgerAddress = "baseledger1xu5xhzj63ddw7pce4r5d0y3w3fuzjxtylzvucm"

    let accountBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}/by_denom?denom=work`)
    .send().expect(200);

    let parsedResponse = JSON.parse(accountBalance.text);

    // save work token balance before deposit
    const workTokenBalanceBefore = parsedResponse.balance.amount
    console.log('work tokens before deposit ', workTokenBalanceBefore);

    const web3 = new Web3('http://localhost:8545');
    let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");

    const accounts = await web3.eth.getAccounts()
    // deposit 1 ubt (8 decimals)
    contract.methods.deposit(100000000, baseledgerAddress).send({
        from: accounts[0]
    }).then(console.log);

    // sleep to wait for attestation to be observed
    await sleep(20000);

    accountBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}/by_denom?denom=work`)
    .send().expect(200);

    parsedResponse = JSON.parse(accountBalance.text);

    // check that balance is increased
    const workTokenBalanceAfter = parsedResponse.balance.amount;
    console.log('work tokens after deposit ', workTokenBalanceAfter);

    expect(+workTokenBalanceBefore).to.be.below(+workTokenBalanceAfter);

    // check if event nonces increased by 1
    const endEventNonces = await getEventNonces();
    console.log('end event nonces ', endEventNonces);
    startEventNonces.forEach((n, i) => {
      expect(n + 1).to.be.equal(endEventNonces[i]);
    });
  });

  it('should jail validator and remove tokens when reducing power back to 0', async function() {
    this.timeout(TEST_TIMEOUT + 90000);

    const orchestratorValidatorResponse = await request(node1_api_url).get('/Baseledger/baseledger/bridge/orchestrator_validator_address')
        .send().expect(200);

    // using first node validator to change staking
    const parsedOrchValResponse = JSON.parse(orchestratorValidatorResponse.text);
    const validatorAddress = parsedOrchValResponse.orchestratorValidatorAddress[0].validatorAddress;

    let validators = await request(node1_api_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    let parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens at start ', parsedResponse.validator.tokens);
            
    const web3 = new Web3('http://localhost:8545');
    let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");

    const accounts = await web3.eth.getAccounts()

    // check for payee
    const payeeExists = await contract.methods.payees(accounts[1]).call()

    // add or update payee with 50k ubt (8 decimals)
    const methodToExecute = payeeExists
      ? contract.methods.updatePayee(accounts[1], accounts[1], 5000000000000, validatorAddress)
      : contract.methods.addPayee(accounts[1], accounts[1], 5000000000000, validatorAddress)
    
    methodToExecute.send({
        from: accounts[0]
    }).then(console.log)

    // sleep to wait for attestation to be observed
    await sleep(20000);

    validators = await request(node1_api_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens after setting to 50k ', parsedResponse.validator.tokens);
    // tokens should be 50k * 10^6
    expect(parsedResponse.validator.tokens).to.be.equal("50000000000");

    // update payee, decrease power to 0
    contract.methods.updatePayee(accounts[1], accounts[1], 0, validatorAddress).send({
        from: accounts[0]
    }).then(console.log)

    // sleep to wait for attestation to be observed and to move validator to unbonding
    await sleep(50000);

    validators = await request(node1_api_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens after decreasing to 0 ', parsedResponse.validator.tokens);
    expect(parsedResponse.validator.tokens).to.be.equal("0");
    expect(parsedResponse.validator.jailed).to.be.equal(true);
    expect(parsedResponse.validator.status).to.be.equal('BOND_STATUS_UNBONDING');
  });
});

// these tests are using 3 nodes and 1 orchestrator - attestations should NOT be observed
describe('attestations NOT observed', async function() {
  this.timeout(50000);
  before(() => {
    startTestNet(3, 1);
  });

  after(() => {
    cleanTestNet();
  });

  it('should NOT deposit ubt to baseledger account', async function () {
    this.timeout(TEST_TIMEOUT + 20000);

    const startEventNonces = await getEventNonces();
    console.log('start event nonces ', await getEventNonces());

    // random regular baseledger address
    const baseledgerAddress = "baseledger1xu5xhzj63ddw7pce4r5d0y3w3fuzjxtylzvucm"

    let accountBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}/by_denom?denom=work`)
    .send().expect(200);

    let parsedResponse = JSON.parse(accountBalance.text);

    // save work token balance before deposit
    const workTokenBalanceBefore = parsedResponse.balance.amount
    console.log('work tokens before deposit ', workTokenBalanceBefore);

    const web3 = new Web3('http://localhost:8545');
    let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");

    const accounts = await web3.eth.getAccounts()
    // deposit 1 ubt (8 decimals)
    contract.methods.deposit(100000000, baseledgerAddress).send({
        from: accounts[0]
    }).then(console.log);

    // sleep to wait for attestation to be observed
    await sleep(20000);

    accountBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}/by_denom?denom=work`)
    .send().expect(200);

    parsedResponse = JSON.parse(accountBalance.text);

    // check that balance was not increased
    const workTokenBalanceAfter = parsedResponse.balance.amount;
    console.log('work tokens after deposit ', workTokenBalanceAfter);

    expect(+workTokenBalanceAfter).to.be.equal(+workTokenBalanceBefore );

    // check if event nonces increased by 1
    const endEventNonces = await getEventNonces();
    console.log('end event nonces ', endEventNonces);
    startEventNonces.forEach((n, i) => {
      expect(n + 1).to.be.equal(endEventNonces[i]);
    });
  });

  it('should NOT add/update validator staking power', async function() {
    this.timeout(TEST_TIMEOUT + 60000);
    const startEventNonces = await getEventNonces();
    console.log('start event nonces ', await getEventNonces());

    const orchestratorValidatorResponse = await request(node1_api_url).get('/Baseledger/baseledger/bridge/orchestrator_validator_address')
        .send().expect(200);

    // using first node validator to change staking
    const parsedOrchValResponse = JSON.parse(orchestratorValidatorResponse.text);
    const validatorAddress = parsedOrchValResponse.orchestratorValidatorAddress[0].validatorAddress;

    let validators = await request(node1_api_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    let parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens at start ', parsedResponse.validator.tokens);
    const startValidatorTokens = parsedResponse.validator.tokens;
    expect(startValidatorTokens).to.be.equal("2000000");
            
    const web3 = new Web3('http://localhost:8545');
    let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");

    const accounts = await web3.eth.getAccounts()
    contract.methods.addPayee(accounts[1], accounts[1], 5000000000000, validatorAddress).send({
        from: accounts[0]
    }).then(console.log)

    // sleep to wait for attestation to be observed
    await sleep(20000);

    validators = await request(node1_api_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens after setting to 50k ', parsedResponse.validator.tokens);
    expect(parsedResponse.validator.tokens).to.be.equal("2000000");

    // update payee, increase power to 80k
    contract.methods.updatePayee(accounts[1], accounts[1], 8000000000000, validatorAddress).send({
        from: accounts[0]
    }).then(console.log)

    // sleep to wait for attestation to be observed
    await sleep(20000);

    validators = await request(node1_api_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens after increasing to 80k ', parsedResponse.validator.tokens);
    expect(parsedResponse.validator.tokens).to.be.equal("2000000");

    // update payee, decrease power to 70k
    contract.methods.updatePayee(accounts[1], accounts[1], 7000000000000, validatorAddress).send({
        from: accounts[0]
    }).then(console.log)

    // sleep to wait for attestation to be observed
    await sleep(20000);

    validators = await request(node1_api_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens after decreasing to 70k ', parsedResponse.validator.tokens);
    expect(parsedResponse.validator.tokens).to.be.equal("2000000");

    // check if event nonces increased by 3
    const endEventNonces = await getEventNonces();
    console.log('end event nonces ', endEventNonces);
    startEventNonces.forEach((n, i) => {
      expect(n + 3).to.be.equal(endEventNonces[i]);
    });
  });
});

describe('add new node', async function() {
  this.timeout(70000);
  before(() => {
    startTestNet();
  });

  after(() => {
    cleanTestNet(4);
  });

  it('should correctly sync new node nonce', async function() {
    this.timeout(TEST_TIMEOUT + 200000);

    const startEventNonces = await getEventNonces();
    console.log('start event nonces ', await getEventNonces());
  
    // random regular baseledger address
    const baseledgerAddress = "baseledger1xu5xhzj63ddw7pce4r5d0y3w3fuzjxtylzvucm"
  
    let accountBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}/by_denom?denom=work`)
    .send().expect(200);
  
    let parsedResponse = JSON.parse(accountBalance.text);
  
    // save work token balance before deposit
    const workTokenBalanceBefore = parsedResponse.balance.amount
    console.log('work tokens before deposit ', workTokenBalanceBefore);
  
    const web3 = new Web3('http://localhost:8545');
    let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");
  
    const accounts = await web3.eth.getAccounts()
    // deposit 1 ubt (8 decimals)
    contract.methods.deposit(100000000, baseledgerAddress).send({
        from: accounts[0]
    }).then(() => {});
  
    // sleep to wait for attestation to be observed
    await sleep(20000);

    // deposit 1 more ubt (8 decimals)
    contract.methods.deposit(100000000, baseledgerAddress).send({
        from: accounts[0]
    }).then(() => {});

    // sleep to wait for attestation to be observed
    await sleep(20000);
  
    accountBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}/by_denom?denom=work`)
    .send().expect(200);
  
    parsedResponse = JSON.parse(accountBalance.text);
  
    // check that balance is increased
    const workTokenBalanceAfter = parsedResponse.balance.amount;
    console.log('work tokens after deposit ', workTokenBalanceAfter);
  
    expect(+workTokenBalanceBefore).to.be.below(+workTokenBalanceAfter);
  
    // check if event nonces increased by 1
    const endEventNonces = await getEventNonces();
    startEventNonces.forEach((n, i) => {
      expect(n +2).to.be.equal(endEventNonces[i]);
    });

    startNewNode();

    await getEventNonces(true);
    await sleep(10000)

    // check that added node event nonce is set to 2
    const noncesAfterAddingNewNode = await getEventNonces(true);
    expect(noncesAfterAddingNewNode.length).to.be.equal(4);
    expect(noncesAfterAddingNewNode).to.have.members([2,2,2,2]);
    accountBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}/by_denom?denom=work`)
    .send().expect(200);
  
    parsedResponse = JSON.parse(accountBalance.text);
    // check that balance is the same
    const workTokenBalanceAfterNewNode = parsedResponse.balance.amount;
    expect(+workTokenBalanceAfterNewNode).to.be.equal(+workTokenBalanceAfter);
  });
});

describe('baseledger transaction', async function() {
  this.timeout(50000);
  before(() => {
    startTestNet();
  });

  after(() => {
    cleanTestNet();
  });

  it('should deposit ubt to baseledger account and use it to send proof', async function () {
    this.timeout(TEST_TIMEOUT + 20000);
    orchAddresses = await getOrchAddresses();

    let faucetBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/baseledger1p8x9ud2m75dmufevmrym3uak0hgcrw58h6n872/by_denom?denom=work`)
     .send().expect(200);

    let parsedResponse = JSON.parse(faucetBalance.text);
    const faucetBalanceStart = parsedResponse.balance.amount;

    const web3 = new Web3('http://localhost:8545');
    let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");

    const accounts = await web3.eth.getAccounts()

    // deposit ubt to all orch addresses
    orchAddresses.forEach(o => {
      contract.methods.deposit(100000000, o).send({
          from: accounts[0]
      }).then(() => {});
    })

    // sleep to wait for attestation to be observed
    await sleep(20000);

    let orchAddressesSumBefore = 0
    orchAddresses.forEach(async o => {
      let balance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/${o}/by_denom?denom=work`)
      .send().expect(200);

      let parsedResponse = JSON.parse(balance.text);
      orchAddressesSumBefore = orchAddressesSumBefore + (+parsedResponse.balance.amount);
    });

    faucetBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/baseledger1p8x9ud2m75dmufevmrym3uak0hgcrw58h6n872/by_denom?denom=work`)
     .send().expect(200);

    parsedResponse = JSON.parse(faucetBalance.text);
    const faucetBalanceBetween = parsedResponse.balance.amount;
    // can not ask for exact because of price variation, so best to just check that balance was reduced
    expect(+faucetBalanceBetween).to.be.below(+faucetBalanceStart);

    const dto = {
      transaction_id: "cbf25e6e-cac1-4afc-8dbf-504eafb3d7d8",
      payload: "f6d0cf9d716e1"
    };

    // our API just returns tx hash, not in json format, but it is set to json
    let txHash = ""
    await request(node1_api_url)
    .post('/signAndBroadcast')
    .set('Accept', 'application/json')
    .send(dto)
    .buffer(true)
    .parse((res, cb) => {
      txHash = Buffer.from("");
      res.on("data", function(chunk) {
        txHash = Buffer.concat([txHash, chunk]);
      });
      res.on("end", function() {
        cb(null, txHash.toString());
      })
    })
    .expect(200);

    await sleep(10000);

    faucetBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/baseledger1p8x9ud2m75dmufevmrym3uak0hgcrw58h6n872/by_denom?denom=work`)
     .send().expect(200);

    parsedResponse = JSON.parse(faucetBalance.text);
    const faucetBalanceEnd = parsedResponse.balance.amount;
    // here we can check that proof posting increased faucet balance for 1
    expect(+faucetBalanceEnd).to.be.equal(+faucetBalanceBetween + 1);
  
    let orchAddressesSumAfter = 0
    orchAddresses.forEach(async o => {
      let balance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/${o}/by_denom?denom=work`)
      .send().expect(200);

      let parsedResponse = JSON.parse(balance.text);
      orchAddressesSumAfter = orchAddressesSumAfter + (+parsedResponse.balance.amount);
    });

    await sleep(10000);

    // check that balance of one of these were reduced by 1
    expect(orchAddressesSumBefore).to.be.equal(orchAddressesSumAfter + 1);
  });

  it('should not be able to send proof without deposit', async function () {
    this.timeout(TEST_TIMEOUT + 20000);
    orchAddresses = await getOrchAddresses();

    let faucetBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/baseledger1p8x9ud2m75dmufevmrym3uak0hgcrw58h6n872/by_denom?denom=work`)
    .send().expect(200);

    let parsedResponse = JSON.parse(faucetBalance.text);
    const faucetBalanceStart = parsedResponse.balance.amount;

    // payload > 512 without additional deposit
    const dto = {
      transaction_id: "cbf25e6e-cac1-4afc-8dbf-504eafb3d712",
      payload: "PYcnlSjsf94yezC2zvaiLy10K1PrUSi5sS5eLzhftG3oWO5yw8rl3YABkbIWulaKqbhdZSxBuRPwUGzQydJHdrqH1t5nyT1Zmc8wcPZ3MuX9nmNWo8XwmFtDP3KwlBDiqJFnWy7aZXIyENLoJkvEHOMgZkx9oqqOgZrjo4iFarYmaR5CqU45zKe4neRmEOb3vrn3oxlro7S01aric07htAnZ450hdlHJPcj6pl14Jc1wHFkxVIbVDFFji6UIUokuLsye8Yed1eLHBtnQvY6KMUD3AggHbuwIU7Qz15KzW91XiU842ypoexg4pXbzwY2C6N9uvFSNpt7dCik4at3fuvTt5JBRUZtY2CZKOgxcOlWAlOQEawUxs1wnpdaKTOnzd0ILLFUvmAfFiFSeSx3Tw5pIxqAbVop7onMewsLndt04MnvwPx3tjiXTtWt83imUo55tjEWjxn92gwst2KPe7uD4JVSSh4EdQQASQno3ZwYRPlOuv8r9Gztuyoj2s8Har"
    };

    const response = await request(node2_api_url)
    .post('/signAndBroadcast')
    .set('Accept', 'application/json')
    .send(dto)
    .expect(500);

    expect(response.body.error).to.be.equal("not enough tokens");

    await sleep(10000);

    faucetBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/baseledger1p8x9ud2m75dmufevmrym3uak0hgcrw58h6n872/by_denom?denom=work`)
     .send().expect(200);

    parsedResponse = JSON.parse(faucetBalance.text);
    const faucetBalanceEnd = parsedResponse.balance.amount;
    // check that faucet balance is not changed
    expect(+faucetBalanceEnd).to.be.equal(+faucetBalanceStart);
  });

  it('should not be able to send proof without proper uuid', async function () {
    this.timeout(TEST_TIMEOUT + 20000);
    orchAddresses = await getOrchAddresses();

    // payload > 512 without additional deposit
    const dto = {
      transaction_id: "cbf25e6e",
      payload: "PYcnlSjsf94yezC2zvaiLy10K1r"
    };

    const response = await request(node2_api_url)
    .post('/signAndBroadcast')
    .set('Accept', 'application/json')
    .send(dto)
    .expect(400);

    expect(response.body.error).to.be.equal("transaction id must be uuid");
  });
});

// assumes build-container.sh is already executed
startTestNet = (nodes = 3, orchs = 3) => {
  shell.exec(path.join(__dirname, `../start-containers.sh ${nodes}`));
  shell.exec(path.join(__dirname, `../deploy-contracts.sh`));
  shell.exec(path.join(__dirname, `../setup-validators.sh ${nodes}`));
  shell.exec(path.join(__dirname, `../run-testnet.sh ${nodes} ${orchs}`));
}

startNewNode = (nodeId = 4) => {
  shell.exec(path.join(__dirname, `../add-new-node.sh ${nodeId}`));
}

cleanTestNet = (nodes = 3) => {
  shell.exec(path.join(__dirname, `../clean.sh ${nodes}`));
}

getEventNonces = async() => {
  const orchestratorAddresses = await getOrchAddresses();
  const eventNonces = []
  for(const [i, v] of orchestratorAddresses.entries()) {
    let nonceResponse = await request(`${host}:${i + 1317}`).get(`/Baseledger/baseledger/bridge/last_event_nonce_by_address/${v}`).send().expect(200);

    eventNonces.push(+JSON.parse(nonceResponse.text).eventNonce);
  }

  return eventNonces;
}

getOrchAddresses = async () => {
  const orchestratorValidatorResponse = await request(node1_api_url).get('/Baseledger/baseledger/bridge/orchestrator_validator_address')
    .send().expect(200);

  const parsedOrchValResponse = JSON.parse(orchestratorValidatorResponse.text);
  return parsedOrchValResponse.orchestratorValidatorAddress.map(o => o.orchestratorAddress);
}