const request = require('supertest');
const Web3 = require('web3')
const fs = require('fs');
const path = require("path");

const starport_url = 'localhost:1317';

const baseledger_abi = JSON.parse(fs.readFileSync(path.join(__dirname, 'baseledger_abi.json')));

describe('test http', () => {
    it ('should send bank request', async () => {
        const orchestratorValidatorResponse = await request(starport_url).get('/Baseledger/baseledger/bridge/orchestrator_validator_address')
            .send().expect(200);

        // TODO: is there a better way to get some sample accounts? these should be as good as any others
        const parsedOrchValResponse = JSON.parse(orchestratorValidatorResponse.text);
        const baseledgerAddress = parsedOrchValResponse.orchestratorValidatorAddress[0].orchestratorAddress;
        const validatorAddress = parsedOrchValResponse.orchestratorValidatorAddress[0].validatorAddress;

        const testBankResponse = await request(starport_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}`)
        .send().expect(200);

        const parsedResponse = JSON.parse(testBankResponse.text);

        console.log('parsed response ', parsedResponse);

        console.log('stake token balance ', parsedResponse.balances[0].amount);
        console.log('work token balance ', parsedResponse.balances[1].amount);
                
        // console.log('ABI ', baseledger_abi)
        // console.log('WEB3 ', await web3.eth.getAccounts());
        const web3 = new Web3('http://localhost:8545');
        let contract = new web3.eth.Contract(baseledger_abi, "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");

        const accounts = await web3.eth.getAccounts()
        contract.methods.updatePayee(accounts[0], accounts[0], 100000010, validatorAddress).send({
            from: accounts[0]
        }).then(console.log)


        const testBankResponse2 = await request(starport_url).get(`/cosmos/bank/v1beta1/balances/${baseledgerAddress}`)
        .send().expect(200);

        const parsedResponse2 = JSON.parse(testBankResponse2.text);

        console.log('parsed response ', parsedResponse2);
        console.log('stake token balance ', parsedResponse2.balances[0].amount);
        console.log('work token balance ', parsedResponse2.balances[1].amount);
        // console.log('CONTRACT ', contract.methods);
    });
});