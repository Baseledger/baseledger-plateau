# Running and upgrading the blockchain app on a validator node

On a validator node, the app should run under (cosmovisor)[https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor] to enable ease of deployment and upgrades.
Cosmovisor is a process that monitors the app for chain upgrades and can perform the stoping, upgrading and restarting of the app.

## Installation

Prerequisite on the validator is golang latest version. 

Then install cosmovisor:

`go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@latest`


### Enviroment vars 

Cosmovisor reads configuration from environment variables. The following ones need to be setup on the node:

*DAEMON_HOME* is the location where the cosmovisor/ directory is kept that contains the genesis binary, the upgrade binaries, and any additional auxiliary files associated with each binary ($HOME/.baseledger-bridged)

*DAEMON_NAME* is the name of the binary itself (baseledger-bridged).

### Folder structure

$DAEMON_HOME/cosmovisor is expected to belong completely to cosmovisor and the subprocesses that are controlled by it. The folder content is organized as follows:

     .
     ├── current -> symbolic link to genesis or upgrades/<name>
     ├── genesis
     │   └── bin
     │       └── $DAEMON_NAME
     └── upgrades
         └── <name>
             ├── bin
             │   └── $DAEMON_NAME
             └── upgrade-info.json

It is the responsibility of the admin to prepare the genesis and upgrades folder structure (current will be created by cosmovisor), place the initial binary in the genesis bin folder and then start cosmovisor.

### Starting the node

Cosmovisor will pass whatever arguments provided to the app in current folder. Starting the chain is a simple:

`cosmovisor start`

It is the job of the admin to setup any systemd process to execute this command on host restart.


# Upgrade process

The admin places the new upgrade binary in the $DAEMON_HOME/upgrades/<name>/bin fodler.
There is a process of voting for a upgrade proposal happening on chain (shown in the example bellow).

As soon as the proposal is voted and the desired height is reached, the blockchain binary from the current folder will create a upgrade-info.json file with the upgrade description in the $DAEMON_HOME/current folder.

This file is parsed by cosmovisor and is a trigger to stop the app, switch the current symbolic link to the upgrade folder and start the app again.

## Migrations of chain data during upgrades

An example of performing a migration during chain upgrade is given on branch *chain-upgrade*. Look for SetUpgradeHandler call in the app.go that defines the upgrade name and actions to be taken.

Also look in the module.go file of the baseledgerbridge module where the migration is registered with the RegisterMigration and where the module version (ConsensusVersion) is increased from 2 to 3. 


## Example local test of chain upgrade

    baseledger-bridged unsafe-reset-all
    baseledger-bridged config chain-id test
    baseledger-bridged config keyring-backend test
    baseledger-bridged config broadcast-mode block
    baseledger-bridged init test --chain-id test --overwrite

Set voting time to 20s in genesis.json

    baseledger-bridged keys add validator
    baseledger-bridged add-genesis-account validator 5000000000stake --keyring-backend test
    baseledger-bridged gentx validator 1000000stake --chain-id test
    baseledger-bridged collect-gentxs

Install cosmovisor

    export DAEMON_NAME=baseledger-bridged
    export DAEMON_HOME=$HOME/.baseledger-bridge

    mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
place current version binary in it

    mkdir -p $DAEMON_HOME/cosmovisor/upgrades/<upgrade_name>/bin
place upgraded version binary in it


    cosmovisor start

Open a new terminal and submit an upgrade propoposal

    baseledger-bridged tx gov submit-proposal software-upgrade <upgrade_name>     --title upgrade --description upgrade --upgrade-height 50 --from validator     --yes

    baseledger-bridged tx gov deposit 1 10000000stake --from validator --yes

    baseledger-bridged tx gov vote 1 yes --from validator --yes

This will create a proposal which is voted pass. when the current version of the node reaches this height, it will halt consensus on all nodes, create a upgrade-info.json that will be read by cosmosvisor and will trigger a switch of the sym link to upgrade folder. Start of the chain will trigger migrations defined in the new app and the chain will continue consensus. 