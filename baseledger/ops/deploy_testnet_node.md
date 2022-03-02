https://docs.cosmos.network/master/run-node/run-node.html

# Preparation

DEVS Deploy contracts to testnet and make sure correct values configured in binaries
DEVS run starport chain build locally for baseledger node
DEVS run cargo build --release locally for orchestrator 
DEVS Prepare the package of the compiled binaries to be shared with other node operators (BASELEDGER_PACKAGE)

# Initialization

Finspot node:

1. Install latest golang - https://www.geeksforgeeks.org/how-to-install-go-programming-language-in-linux/
2. Install cosmovisor: Follow the *deployment_and_chain_upgrade.md*
3. Place binaries in the respective folders as described in *deployment_and_chain_upgrade.md*
4. Initialize: cosmovisor init validator --chain-id=baseledger
5. Configure genesis: Edit the genesis file for various params (voting time, inflation, ubonding time, tokens metadata etc.)
6. Generate validator account: ./baseledgerd keys add --keyring-backend file validator (make sure to write down the address and the mnemonic)
7. Generate orchestrator account: ./baseledgerd keys add --keyring-backend file orchestrator (make sure to write down the address and the mnemonic)

./baseledgerd keys add faucet --recover <faucet_mnemonic>

8. TODO - Add faucet address generation and addition to genesis as a bridge param. Add faucet account with allocation: ./baseledgerd add-genesis-account --keyring-backend file baseledger1xgs5tamqre7rkz5q7d5fegjsdwufxxvt36w0a0 10000000000stake,10000000000work
9. Add validator account with allocation: ./baseledgerd  add-genesis-account --keyring-backend file <validator_address> 1000000stake
10. Add orchestrator account with allocation: ./baseledgerd  add-genesis-account --keyring-backend file <orchestrator_address> 1work
11. Add gentx transaction: ./baseledgerd gentx --keyring-backend file --moniker finspot_validator --ip <validator_ip> --chain-id=baseledger validator 1000000stake
12. Extract the genesis and add it to the BASELEDGER_PACKAGE

Other genesis nodes:

1. Download the BASELEDGER_PACKAGE
2. Install latest golang
3. Install cosmovisor: Follow the *deployment_and_chain_upgrade.md*
4. Place binaries in the respective folders as described in *deployment_and_chain_upgrade.md*
5. Generate validator account: ./baseledgerd keys add --keyring-backend file validator (make sure to write down the address and the mnemonic)
6. Generate orchestrator account: ./baseledgerd keys add --keyring-backend file orchestrator (make sure to write down the address and the mnemonic)
7. Add validator account with allocation: ./baseledgerd add-genesis-account --keyring-backend file <validator_address> 1000000stake
8. Add orchestrator account with allocation: ./baseledgerd add-genesis-account --keyring-backend file <orchestrator_address> 1work
9. Add gentx transaction: ./baseledgerd gentx --keyring-backend file --moniker <organization>_validator --ip <validator_ip_address> --chain-id=baseledger validator 1000000stake
10. Extract the genesis and gentx and send over to Finspot, together with a node id (./baseledgerd tendermint show-node-id) and the static ip address of the node.


Finspot node:

1. Make sure latest genesis and gentx transactions are present
2. ./baseledgerd collect-gentxs
3. Prepare cosmovisor start scripts with all node ids and ips 
4. Distribute genesis and run script to each validator

# Start

Each node:

1. Place the latest genesis in the appropriate folder
2. Run the node and add all persistent peers as coma delimited list: cosmovisor --p2p.persistent_peers <node1_id>@<node1_ip>:26656,<node2_id>... start
3. Init orchestrator: ./baseledger_bridge init
4. export COINMARKETCAP_API_TOKEN=<token>
5. export COINAPI_API_TOKEN=<token>
6. Register orchestrator key: ./baseledger_bridge keys set-orchestrator-key --phrase=<orchestrator_mnemonic>
7. Register orchestrator address: ./baseledger_bridge keys register-orchestrator-address --validator-phrase=<validator_mnemonic>
8. Run ochestrator: baseledger_bridge orchestrator --ethereum-rpc=<your_eth_url_such_as_infura> --baseledger-contract-address=<baseledger_contract_address_provided_by_the_team>




Additional node joining:

1. Download the BASELEDGER_PACKAGE
2. Install latest golang
3. Install cosmovisor: Follow the *deployment_and_chain_upgrade.md*
4. Place binaries in the respective folders as described in *deployment_and_chain_upgrade.md*
5. Generate validator account: ./baseledgerd keys add --keyring-backend file validator (make sure to write down the address and the mnemonic)
6. Generate orchestrator account: ./baseledgerd keys add --keyring-backend file orchestrator (make sure to write down the address and the mnemonic)
8.  Place the latest genesis in the appropriate folder
9. Run the node and add all persistent peers as coma delimited list: cosmovisor --p2p.persistent_peers <node1_id>@<node1_ip>:26656,<node2_id>... start
10. chmod +x baseledger_bridge
11. Init orchestrator: ./baseledger_bridge init
12. export COINMARKETCAP_API_TOKEN=<token>
13. export COINAPI_API_TOKEN=<token>
14. Register orchestrator key: ./baseledger_bridge keys set-orchestrator-key --phrase=<orchestrator_mnemonic>
15. Register orchestrator address: ./baseledger_bridge keys register-orchestrator-address --validator-phrase=<validator_mnemonic>
16. Run ochestrator: baseledger_bridge orchestrator --ethereum-rpc=<your_eth_url_such_as_infura> --baseledger-contract-address=<baseledger_contract_address_provided_by_the_team>


faucet:
./baseledgerd tx bank send faucet <new_node_address> 1stake --yes
./baseledgerd tx bank send faucet <new_orchestrator_address> 1work --yes

Here <new_node_address> is the receiver address obtained from baseledgerd keys list command

additional node:
./baseledgerd tx staking create-validator  --amount=1stake  --pubkey=baseledgervalconspub1zcjduepq0y6gpu79m6ltgjlxs2x0t0ygfdkhnjjxkdl75ejcslcpat3zytlqjp6sty --moniker="node55"  --commission-rate="0.10" --commission-max-rate="0.20" --commission-max-change-rate="0.01" --min-self-delegation="1" --from=node55_validator_address1 --yes 

--from = <name of the node to become validator>
--pubkey <output of tendermint show-validator on node_to_become_validator>
--moniker= <unique name for the validator>

faucet:
./baseledgerd tx staking delegate baseledgervaloper1kkf4ujsjj8vuj9575qw5tlm53nnwxufycnj9ru 100000000stake --from=node1_validator_address_1 --yes 

Params explanation:
--baseledgervaloper-address from the new validator node, can be seen in "docker exec first_node_blockchain_app_1 baseledgerd query staking validators"
--from=<our token controlling node1>