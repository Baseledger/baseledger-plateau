poc-contract is a folder with contracts, contains erc-20 and test contract for deposit that is throwing event.

Folder orchestrator is simplified copied from gravity-bridge (https://github.com/Gravity-Bridge/Gravity-Bridge), i tried to not use anything that is not needed for first version.

Orchestrator folder structure (kept same names as gravity, we will rename when we only left what’s needed and maybe change structure a bit)

baseledger_bridge - cli program to start ethereum oracle

baseledger_proto - this is auto generated from cosmos proto files, i just copied this from gravity repo, do not review this at the moment, our will hopefully be much smaller, just type definitions, for example i needed this for grpc client types

utils - various utils needed, only copied some (util methods for connection to ethereum and cosmos, types, error types etc)

ethereum_oracle - main thing for us atm, there is loop listening to ethereum events, it is started as cli command

Biggest difference to gravity bridge is that gravity is doing more things not needed by us (at least atm) like relaying stuff to ethereum, transaction batching to ethereum, listening to all those events (we only listen to SendToCosmosEvent so everything else is removed) etc.

To make this work locally apart from starport rust is needed to be installed and then call

1. check poc_contracts to run hardhat

2. run `starport chain serve --verbose` in baseledger folder (if starting from scratch run `starport chain serve --verbose --reset-once` and copy alice and bob mnemonics for further usage)

3. `cargo build --all` in root of orchestrator folder

4. navigate to baseledger_bridge folder and execute

```shell
cargo run -- init 

cargo run -- keys set-orchestrator-key --phrase="<STARPORT_BOB_PHRASE>"

cargo run -- keys register-orchestrator-address --fees="0token" --validator-phrase="<STARPORT_ALICE_PHRASE>"

export COINMARKETCAP_API_TOKEN=<token>

cargo run -- orchestrator --fees "0token" --ethereum-rpc="http://localhost:8545" --baseledger-contract-address="<BASELEDGER_TEST_CONTRACT_ADDRESS>"
```

## Changing and building proto files

- change the proto files in baseledger bridge
- navigate to <root>/baseledger
- starport chain build --proto-all-modules
- navigate to <root>/orchestrator/proto_build
- cargo run