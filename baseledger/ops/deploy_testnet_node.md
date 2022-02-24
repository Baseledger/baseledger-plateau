https://docs.cosmos.network/master/run-node/run-node.html


# Preparation

DEVS Deploy contracts to testnet and make sure correct values configured in binaries
DEVS run starport chain build locally for baseledger node
DEVS run cargo build --release locally for orchestrator 
DEVS Prepare the package of the compiled binaries to be shared with other node operators (BASELEDGER_PACKAGE)

# Initialization

Finspot node:

1. Install latest golang
2. Install cosmovisor: Follow the *deployment_and_chain_upgrade.md*
3. Place binaries in the respective folders as described in *deployment_and_chain_upgrade.md*
4. Initialize: cosmovisor init validator --chain-id=baseledger
5. Configure genesis: Edit the genesis file for various params (voting time, tokens metadata etc.)
6. Generate validator account: cosmovisor keys add --keyring-backend test validator (make sure to write down the address and the mnemonic)
7. Generate orchestrator account: cosmovisor keys add --keyring-backend test orchestrator (make sure to write down the address and the mnemonic)
8. Add faucet account with allocation: cosmovisor add-genesis-account --keyring-backend test baseledger1xgs5tamqre7rkz5q7d5fegjsdwufxxvt36w0a0 10000000000work
9. Add validator account with allocation: cosmovisor add-genesis-account --keyring-backend test <validator_address> 10000000000stake,10000000000work
10. Add orchestrator account with allocation: cosmovisor add-genesis-account --keyring-backend test <orchestrator_address> 1work
11. Add gentx transaction: cosmovisor gentx --keyring-backend --moniker finspot_validator --ip <validator_ip_address> validator 10000000000stake
12. Extract the genesis and add it to the BASELEDGER_PACKAGE

Other nodes:

1. Download the BASELEDGER_PACKAGE
2. Install latest golang
3. Install cosmovisor: Follow the *deployment_and_chain_upgrade.md*
4. Place binaries in the respective folders as described in *deployment_and_chain_upgrade.md*
5. Generate validator account: cosmovisor keys add --keyring-backend test validator (make sure to write down the address and the mnemonic)
6. Generate orchestrator account: cosmovisor keys add --keyring-backend test orchestrator (make sure to write down the address and the mnemonic)
7. Add validator account with allocation: cosmovisor add-genesis-account --keyring-backend test <validator_address> 10000000000stake,10000000000work
8. Add orchestrator account with allocation: cosmovisor add-genesis-account --keyring-backend test <orchestrator_address> 1work
9. Add gentx transaction: cosmovisor gentx --keyring-backend --moniker <organization>_validator --ip <validator_ip_address> validator 10000000000stake
10. Extract the genesis and gentx and send over to Finspot


Finspot node:

1. Make sure latest genesis and gentx transactions are present
2. cosmovisor collect-gentxs
3. Distrubute genesis to each validator

# Start

Each node:

1. Place the latest genesis in the appropriate folder
2. Run the node and add all persistent peers as coma delimited list: cosmovisor --p2p.persistent_peers <node1_id>@<node1_ip>:26656,<node2_id>... start
3. Register orchestrator key: baseleger_bridge -- keys set-orchestrator-key --phrase=<orchestrator_mnemonic>
4. Register orchestrator address: baseleger_bridge -- keys register-orchestrator-address --validator-phrase=<validator_mnemonic>
6. Get coinmarket cap api token and set env var COINMARKETCAP_API_TOKEN
7. Get coinapi api token and set env var COINAPI_API_TOKEN
8. Run ochestrator: baseledger_bridge -- orchestrator --ethereum-rpc=<your_eth_url_such_as_infura> --baseledger-contract-address=<baseledger_contract_address_provided_by_the_team>
