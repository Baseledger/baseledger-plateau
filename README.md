# Baseledger

*Baseledger is currently under development. Below is short overview of current work.*

## Modules overview and licence disclaimer

### Proof module
This module is similar to Baseledger Lakewood (https://github.com/Baseledger/baseledger-lakewood) and it is used for storing proofs. Everyone who executes proof storing transactions would need to pay Work tokens fee, that are depending on payload size (details on this will follow). This module also exposes custom cosmos-sdk client for signing and broadcasting transactions within REST endpoint, just by sending proof as a string, and it uses preconfigured keys stored in node file keyring (https://docs.cosmos.network/master/run-node/keyring.html).

### Bridge module
This module is forked from Gravity Bridge (https://github.com/Gravity-Bridge/Gravity-Bridge). Unlike Gravity Bridge, it is one way bridge (Ethereum => Cosmos), and it is listening and handling our application-specific events.
Even though purpose is different and it is not separate chain, but only a module, structure and flow of bridge is following Gravity Bridge good practices: there is orchestrator (only with ethereum oracle in our case) that validators will need to run, that is listening to events and sending claim transactions, that are then voted within attestations.
Also, compared to Gravity Bridge, we are using starport (https://github.com/tendermint/starport) to scaffold cosmos module.

Overview of changes compared to Gravity Bridge are:
- one way bridge (Ethereum => Cosmos)
- different smart contract (we do not use Gravity.sol)
- different events
- removed everything that we don't need (relayer, ethereum key etc) - only thing left is ethereum oracle for listening to baseledger specific events
- simplified module structure in orchestrator due to simplified overall code
- cosmos module was generated using starport
...

## Quick developer start (work in progress)

To make this work locally apart from starport rust is needed to be installed and then call

1. check https://github.com/Baseledger/baseledger-contracts to run hardhat

2. run `starport chain serve --verbose` in baseledger folder (if starting from scratch run `starport chain serve --verbose --reset-once` and copy alice and bob mnemonics for further usage)

3. `cargo build --all` in root of orchestrator folder

4. navigate to baseledger_bridge folder and execute

```shell
cargo run -- init 

cargo run -- keys set-orchestrator-key --phrase="<STARPORT_BOB_PHRASE>"

cargo run -- keys register-orchestrator-address --fees="0token" --validator-phrase="<STARPORT_ALICE_PHRASE>"

export COINMARKETCAP_API_TOKEN=<token>
export COINAPI_API_TOKEN=<token>

cargo run -- orchestrator --ethereum-rpc="http://localhost:8545" --baseledger-contract-address="<BASELEDGER_TEST_CONTRACT_ADDRESS>"
```

## Changing and building proto files

- change the proto files in baseledger bridge
- navigate to <root>/baseledger
- starport chain build --proto-all-modules
- navigate to <root>/orchestrator/proto_build
- cargo run

## Running a local dockerized testnet

- Navigate to tests
- Run build-container.sh - This should be ran only once to build the docker images.
- Run start-containers.sh - Starts 3 baseledger nodes and a hardhat node.
- Run deploy-contracts.sh - This deploys the dummy UBT token contract as well as the BaseledgerUBTSplitter contract
- Run setup-validators.sh - Creates and shares genesis.json among validators and creates gentx files.
- Run run-testnet.sh - Starts the nodes, registers and starts the orchestrators. 