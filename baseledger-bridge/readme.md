local test:

baseledger-bridged unsafe-reset-all

baseledger-bridged config chain-id test
baseledger-bridged config keyring-backend test
baseledger-bridged config broadcast-mode block

baseledger-bridged init test --chain-id test --overwrite

set voting time to 20s in genesis.json

baseledger-bridged keys add validator
baseledger-bridged add-genesis-account validator 5000000000stake --keyring-backend test
baseledger-bridged gentx validator 1000000stake --chain-id test
baseledger-bridged collect-gentxs

Install cosmovisor

export DAEMON_NAME=baseledger-bridged
export DAEMON_HOME=$HOME/.baseledger-bridge
export DAEMON_RESTART_AFTER_UPGRADE=true

mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
place current version binary in it

mkdir -p $DAEMON_HOME/cosmovisor/upgrades/test_plan_for_upgrade/bin
place upgraded version binary in it

start the chain with cosmovisor start - this will create a sym link $DAEMON_HOME/cosmovisor/current that points to genesis/bin

new terminal, submit an upgrade propoposal

baseledger-bridged tx gov submit-proposal software-upgrade test_plan_for_upgrade --title upgrade --description upgrade --upgrade-height 50 --from validator --yes

baseledger-bridged tx gov deposit 1 10000000stake --from validator --yes

baseledger-bridged tx gov vote 1 yes --from validator --yes

This will create a proposal which is voted pass. when the current version of the node reaches this height, it will halt consensus on all nodes, create a upgrade-info.json that will be read by cosmosvisor and will trigger a switch of the sym link to upgrade folder. Start of the chain will trigger migrations defined in the new app and the chain will continue consensus. 