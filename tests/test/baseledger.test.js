const request = require('supertest');
const Web3 = require('web3')
const fs = require('fs');
const path = require("path");
const chai = require("chai");
const expect = chai.expect;

const node1_api_url = 'localhost:1317';
const node2_api_url = 'localhost:1318';
const node3_api_url = 'localhost:1319';

const baseledger_abi = JSON.parse(fs.readFileSync(path.join(__dirname, 'baseledger_abi.json')));

const sleep = (ms) => {
  return new Promise(resolve => setTimeout(resolve, ms));
};
const TEST_TIMEOUT = 30000;

describe('validator power update', () => {
  it('should add/update validator staking power', async function() {
    this.timeout(TEST_TIMEOUT + 60000);
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

    // check that balance increased by 1
    const workTokenBalanceAfter = parsedResponse.balance.amount;
    console.log('work tokens after deposit ', workTokenBalanceAfter);

    expect(+workTokenBalanceAfter).to.be.equal(+workTokenBalanceBefore + 1);
  });
});