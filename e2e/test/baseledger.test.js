const request = require('supertest');
const Web3 = require('web3')
const fs = require('fs');
const path = require("path");
const chai = require("chai");
const expect = chai.expect;

const starport_url = 'localhost:1317';

const node1_api_url = 'localhost:1317';
const node2_api_url = 'localhost:1318';
const node3_api_url = 'localhost:1317';

const baseledger_abi = JSON.parse(fs.readFileSync(path.join(__dirname, 'baseledger_abi.json')));

const sleep = (ms) => {
  return new Promise(resolve => setTimeout(resolve, ms));
};
const TEST_TIMEOUT = 30000;

describe('validator power update', () => {
  it('should add/update validator staking power', async function() {
    this.timeout(TEST_TIMEOUT + 45000);
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
            
    // add payee
    const web3 = new Web3('http://localhost:8545');
    let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");

    const accounts = await web3.eth.getAccounts()
    contract.methods.addPayee(accounts[1], accounts[1], 100000100, validatorAddress).send({
        from: accounts[0]
    }).then(console.log)

    await sleep(15000);

    validators = await request(starport_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens ', parsedResponse.validator.tokens);
    expect(parsedResponse.validator.tokens).to.be.equal("100000100");

    // update payee, increase power
    contract.methods.updatePayee(accounts[1], accounts[1], 100000200, validatorAddress).send({
        from: accounts[0]
    }).then(console.log)

    await sleep(15000);

    validators = await request(starport_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens ', parsedResponse.validator.tokens);
    expect(parsedResponse.validator.tokens).to.be.equal("100000200");

    // update payee, decrease power
    contract.methods.updatePayee(accounts[1], accounts[1], 100000150, validatorAddress).send({
        from: accounts[0]
    }).then(console.log)

    await sleep(15000);

    validators = await request(starport_url).get(`/cosmos/staking/v1beta1/validators/${validatorAddress}`)
    .send().expect(200);

    parsedResponse = JSON.parse(validators.text);

    console.log('validator tokens ', parsedResponse.validator.tokens);
    expect(parsedResponse.validator.tokens).to.be.equal("100000150");
  });

});

describe('ubt deposit', () => {
  it('should deposit ubt to baseledger account', async function () {
    this.timeout(TEST_TIMEOUT + 20000);
    // random regular baseledger address
    const baseledgerAddress = "baseledger1xu5xhzj63ddw7pce4r5d0y3w3fuzjxtylzvucm"

    let accountBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}/by_denom?denom=work`)
    .send().expect(200);

    let parsedResponse = JSON.parse(accountBalance.text);

    // save work token balance before deposit
    const workTokenBalanceBefore = parsedResponse.balance.amount

    const web3 = new Web3('http://localhost:8545');
    let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");

    const accounts = await web3.eth.getAccounts()
    // deposit 1 ubt (8 decimals)
    contract.methods.deposit(100000000, baseledgerAddress).send({
        from: accounts[0]
    }).then(console.log);

    await sleep(20000);

    accountBalance = await request(node1_api_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}/by_denom?denom=work`)
    .send().expect(200);

    parsedResponse = JSON.parse(accountBalance.text);

    // check that balance increased by 1
    const workTokenBalanceAfter = parsedResponse.balance.amount;
    expect(+workTokenBalanceAfter).to.be.equal(+workTokenBalanceBefore + 1);
  });
});