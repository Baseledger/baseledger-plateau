# Running the test suite

## Clean test network between test runs (todo: make this automatic in future)
- Run clean.sh
## Start dockerized test network 

- Run build-container.sh - This should be ran only once to build the docker images.
- Run start-containers.sh - Starts 3 baseledger nodes and a hardhat node.
- Run deploy-contracts.sh - This deploys the dummy UBT token contract as well as the BaseledgerUBTSplitter contract
- Run setup-validators.sh - Creates and shares genesis.json among validators and creates gentx files.
- Run run-testnet.sh - Starts the nodes, registers and starts the orchestrators. 

## Install dependencies
npm i

## Run tests
npm test