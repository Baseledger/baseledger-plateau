# Running and upgrading the blockchain app on a validator node

On a validator node, the app should run under (cosmovisor)[https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor] to enable ease of deployment and upgrades.
Cosmovisor is a process that monitors the app for chain upgrades and can perform the stoping, upgrading and restarting of the app.

## Installation

Prerequisite on the validator is golang latest version. 

Then install cosmovisor:

`go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@latest`


### Enviroment vars 

Cosmovisor reads configuration from environment variables. The following ones need to be setup on the node:

*DAEMON_HOME* is the location where the cosmovisor/ directory is kept that contains the genesis binary, the upgrade binaries, and any additional auxiliary files associated with each binary ($HOME/.baseledger)
export DAEMON_HOME=$HOME/.baseledger

*DAEMON_NAME* is the name of the binary itself (baseledgerd).
export DAEMON_NAME=baseledgerd

### Folder structure

$DAEMON_HOME/cosmovisor is expected to belong completely to cosmovisor and the subprocesses that are controlled by it. The folder content is organized as follows:

     $DAEMON_HOME/cosmovisor
                    ├── current -> symbolic link to genesis or upgrades/<name>
                    ├── genesis
                    │   └── bin
                    │       └── $DAEMON_NAME
                    └── upgrades
                        └── <name>
                            ├── bin
                            │   └── $DAEMON_NAME
                            └── upgrade-info.json

Cosmovisor requires $DAEMON_HOME/data folder to exist as well

It is the responsibility of the admin to prepare the genesis and upgrades folder structure (current will be created by cosmovisor), place the initial binary in the genesis bin folder and then start cosmovisor.

Make sure binary in genesis is marked as executable.
chmod +x baseledgerd

Make sure current contains a json file upgrade-info.json with content {}

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

Also look in the module.go file of the bridge module where the migration is registered with the RegisterMigration and where the module version (ConsensusVersion) is increased from 2 to 3. 

## Backups TODO
Prior to the upgrade, validators are encouraged to take a full data snapshot. Snapshotting depends heavily on infrastructure, but generally this can be done by backing up the .baseledger directory. If you use Cosmovisor to upgrade, by default, Cosmovisor will backup your data upon upgrade.

It is critically important for validator operators to back-up the .baseledger/data/priv_validator_state.json file after stopping the baseledgerd process. This file is updated every block as your validator participates in consensus rounds. It is a critical file needed to prevent double-signing, in case the upgrade fails and the previous chain needs to be restarted.


## Example local test of chain upgrade

    baseledgerd unsafe-reset-all
    baseledgerd config chain-id test
    baseledgerd config keyring-backend test
    baseledgerd config broadcast-mode block
    baseledgerd init test --chain-id test --overwrite

Set voting time to 20s in /root/.baseledger/config/genesis.json

    baseledgerd keys add validator
    baseledgerd add-genesis-account validator 5000000000stake --keyring-backend test
    baseledgerd gentx validator 1000000stake --chain-id test
    baseledgerd collect-gentxs

Install cosmovisor

    export DAEMON_NAME=baseledgerd
    export DAEMON_HOME=$HOME/.baseledger

    mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
place current version binary in it

    mkdir -p $DAEMON_HOME/cosmovisor/upgrades/<upgrade_name>/bin
place upgraded version binary in it


    cosmovisor start

Open a new terminal and submit an upgrade propoposal

    baseledgerd tx gov submit-proposal software-upgrade <upgrade_name>     --title upgrade --description upgrade --upgrade-height 50 --from validator     --yes

    baseledgerd tx gov deposit 1 10000000stake --from validator --yes

    baseledgerd tx gov vote 1 yes --from validator --yes

This will create a proposal which is voted pass. when the current version of the node reaches this height, it will halt consensus on all nodes, create a upgrade-info.json that will be read by cosmosvisor and will trigger a switch of the sym link to upgrade folder. Start of the chain will trigger migrations defined in the new app and the chain will continue consensus. 


## Example local test of parameter change proposal

    baseledgerd unsafe-reset-all
    baseledgerd config chain-id baseledger
    baseledgerd config keyring-backend test
    baseledgerd config broadcast-mode block
    baseledgerd init baseledger --chain-id baseledger --overwrite

Set voting time to 20s in /root/.baseledger/config/genesis.json

    baseledgerd keys add validator
    baseledgerd add-genesis-account validator 5000000000stake --keyring-backend test
    baseledgerd gentx validator 1000000stake --chain-id baseledger
    baseledgerd collect-gentxs
    baseledgerd start


Open a new terminal and submit a param change propoposal

    copy the ./params.json to the folder where you are running the command

    baseledgerd tx gov submit-proposal param-change param.json  --from validator     --yes

    baseledgerd tx gov deposit 1 10000000stake --from validator --yes

    baseledgerd tx gov vote 1 yes --from validator --yes

Misc

    baseledgerd query bridge params - see the current value of the param

    baseledgerd query gov proposals - see status of proposals