const request = require('supertest');
const Web3 = require('web3')
const fs = require('fs');
const path = require("path");

const starport_url = 'localhost:1317';

const baseledger_abi = JSON.parse(fs.readFileSync(path.join(__dirname, 'baseledger_abi.json')));

const sleep = (ms) => {
  return new Promise(resolve => setTimeout(resolve, ms));
};
const TEST_TIMEOUT = 30000;

describe('validator power update', () => {
  // TODO: test to add validator only 

  it('should update validator staking power', async () => {
    const orchestratorValidatorResponse = await request(starport_url).get('/Baseledger/baseledger/bridge/orchestrator_validator_address')
        .send().expect(200);

    // TODO: is there a better way to get some sample accounts? these should be as good as any others
    const parsedOrchValResponse = JSON.parse(orchestratorValidatorResponse.text);
    const baseledgerAddress = parsedOrchValResponse.orchestratorValidatorAddress[0].orchestratorAddress;
    const validatorAddress = parsedOrchValResponse.orchestratorValidatorAddress[0].validatorAddress;

    let validators = await request(starport_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    let parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens ', parsedResponse.validator.tokens);
            
    // console.log('ABI ', baseledger_abi)
    // console.log('WEB3 ', await web3.eth.getAccounts());
    const web3 = new Web3('http://localhost:8545');
    let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");

    const accounts = await web3.eth.getAccounts()
    contract.methods.updatePayee(accounts[0], accounts[0], 100000010, validatorAddress).send({
        from: accounts[0]
    }).then(console.log)


    accountBalance = await request(starport_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(accountBalance.text);

    console.log('parsed response ', parsedResponse);
    console.log('stake token balance ', parsedResponse.balances[0].amount);
    console.log('work token balance ', parsedResponse.balances[1].amount);

    // TODO: check staking power before and after
  });

  // TODO: test to decrease power
});

describe('ubt deposit', () => {
  it('should deposit ubt to baseledger account', async function () {
    this.timeout(TEST_TIMEOUT + 20000);
    // TODO: how to get/generate account?
    const baseledgerAddress = "baseledger1xu5xhzj63ddw7pce4r5d0y3w3fuzjxtylzvucm"

    let accountBalance = await request(starport_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}`)
    .send().expect(200);

    let parsedResponse = JSON.parse(accountBalance.text);

    console.log('parsed response ', parsedResponse);

    console.log('stake token balance ', parsedResponse.balances[0].amount);
    console.log('work token balance ', parsedResponse.balances[1].amount);

    const web3 = new Web3('http://localhost:8545');
    let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");

    const accounts = await web3.eth.getAccounts()
    contract.methods.deposit(100, baseledgerAddress).send({
        from: accounts[0]
    }).then(console.log)

    await sleep(20000)

    accountBalance = await request(starport_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(accountBalance.text);

    console.log('parsed response ', parsedResponse);

    console.log('stake token balance ', parsedResponse.balances[0].amount);
    console.log('work token balance ', parsedResponse.balances[1].amount);
  });
});