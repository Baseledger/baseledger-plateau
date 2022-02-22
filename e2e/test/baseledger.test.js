const request = require('supertest');

const starport_url = 'localhost:1317';

describe('test http', () => {
    it ('should send bank request', async () => {
        const testBankResponse = await request(starport_url).get('/cosmos/bank/v1beta1/balances/baseledger17m33f6hu59dy9a9mu6utehlmfsn5xjruzyuxer')
        .send().expect(200);

        const parsedResponse = JSON.parse(testBankResponse.text);

        console.log('parsed response ', parsedResponse);

        console.log('stake token balance ', parsedResponse.balances[0].amount);
        console.log('work token balance ', parsedResponse.balances[1].amount);
    });
});