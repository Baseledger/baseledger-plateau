# Basic Sample Baseledger POC project

## Deploy test erc20 and baseledger contracts

Deploy script will also allow 100 tokens to be deposited to baseledger bridge
```shell
npx hardhat node

npx hardhat run --network localhost scripts/deploy_baseledger_test.js

npx hardhat console --network localhost
```

## Test deposit event

```shell
npx hardhat console --network localhost

const Baseledger = await ethers.getContractFactory("BaseledgerTest")
const baseledger = await Baseledger.attach(BASELEDGER_CONTRACT_ADDRESS)

await baseledger.deposit(1, COSMOS_WALLET_ADDRESS)
```