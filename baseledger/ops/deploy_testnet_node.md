https://docs.cosmos.network/master/run-node/run-node.html

# Preparation

DEVS Deploy contracts to testnet and make sure correct values configured in binaries
DEVS Prepare faucet account mnemonic and hardcode the address
DEVS run starport chain build locally for baseledger node
DEVS run cargo build --release locally for orchestrator 
DEVS Prepare the package of the compiled binaries to be shared with other node operators (BASELEDGER_PACKAGE)

# Initialization

Finspot node:

1. Install latest golang - https://www.geeksforgeeks.org/how-to-install-go-programming-language-in-linux/
2. Install cosmovisor: Follow the *deployment_and_chain_upgrade.md*
3. Place binaries in the respective folders as described in *deployment_and_chain_upgrade.md*
4. Initialize: ./baseledgerd init validator --chain-id=baseledger
5. Configure genesis: Edit the genesis file for various params (voting time, inflation, ubonding time, tokens metadata etc.)
6. Generate validator account: ./baseledgerd keys add --keyring-backend file validator (make sure to write down the address and the mnemonic)
7. Generate orchestrator account: ./baseledgerd keys add --keyring-backend file orchestrator (make sure to write down the address and the mnemonic)
8. Add faucet account with allocation: ./baseledgerd add-genesis-account --keyring-backend file <faucet_address> 10000000000stake,10000000000work
9. Add validator account with allocation: ./baseledgerd  add-genesis-account --keyring-backend file <validator_address> 2000000stake
10. Add orchestrator account with allocation: ./baseledgerd  add-genesis-account --keyring-backend file <orchestrator_address> 1work
11. Add gentx transaction: ./baseledgerd gentx --keyring-backend file --moniker finspot_validator --ip <validator_ip> --chain-id=baseledger validator 2000000stake
12. ./baseledgerd keys add --recover --keyring-backend file faucet
13. Extract the genesis and add it to the BASELEDGER_PACKAGE

Other nodes:

1. Download the BASELEDGER_PACKAGE
2. Install latest golang
3. Install cosmovisor: Follow the *deployment_and_chain_upgrade.md*
4. Place binaries in the respective folders as described in *deployment_and_chain_upgrade.md*
5. Initialize: ./baseledgerd init validator --chain-id=baseledger
6. Generate validator account: ./baseledgerd keys add --keyring-backend file validator (make sure to write down the address and the mnemonic)
7. Generate orchestrator account: ./baseledgerd keys add --keyring-backend file orchestrator (make sure to write down the address and the mnemonic)
8. Add validator account with allocation: ./baseledgerd add-genesis-account --keyring-backend file <validator_address> 1000000stake
9. Add orchestrator account with allocation: ./baseledgerd add-genesis-account --keyring-backend file <orchestrator_address> 1work
10. Add gentx transaction: ./baseledgerd gentx --keyring-backend file --moniker <organization>_validator --ip <validator_ip_address> --chain-id=baseledger validator 1000000stake
11. Extract the genesis and gentx and send over to Finspot, together with a node id (./baseledgerd tendermint show-node-id) and the static ip address of the node.


Finspot node:

1. Make sure latest genesis and gentx transactions are present
2. ./baseledgerd collect-gentxs
3. Prepare cosmovisor.service and orchestrator.service scripts with all node ids and ips and relevant params (contract address) and add to /etc/systemd/system/
4. Distribute genesis and run script to each validator

# Start

Each node:

1. Place the latest genesis in the appropriate folder
2. Adjust the orchestrator.service with your local params (infura, coinmarketapi key, coinapi key)
2. setup the cosmovisor.service and orchestrator.service to run as a systemd services
3. Init orchestrator: ./baseledger_bridge init
4. Register orchestrator key: ./baseledger_bridge keys set-orchestrator-key --phrase=<orchestrator_mnemonic>
5. systemctl start cosmovisor.service
6. Register orchestrator address: ./baseledger_bridge keys register-orchestrator-address --validator-phrase=<validator_mnemonic>
7. systemctl start orchestrator.service
