#!/bin/bash
set -eux
# your baseledger binary name
BIN=baseledgerd

CHAIN_ID="baseledger"

NODES=$1

ALLOCATION="10000000000stake,10000000000worktoken"

# first we start a genesis.json with validator 1
# validator 1 will also collect the gentx's once gnerated
STARTING_VALIDATOR=1
STARTING_VALIDATOR_HOME="--home /validator"
STARTING_VALIDATOR_CONTAINER="baseledger-validator-container"
docker exec $STARTING_VALIDATOR_CONTAINER $BIN init $STARTING_VALIDATOR_HOME validator --chain-id=$CHAIN_ID


## Modify generated genesis.json to our liking by editing fields using jq
## we could keep a hardcoded genesis file around but that would prevent us from
## testing the generated one with the default values provided by the module.

# add in denom metadata for both native tokens
jq '.app_state.bank.denom_metadata += [{"name": "Work Token", "symbol": "WRK", "base": "worktoken", "description": "A test token for paying work", "denom_units": [{"denom": "worktoken", "exponent": 0}]},{"name": "Stake Token", "symbol": "STK", "base": "stake", "description": "A staking test token", "denom_units": [{"denom": "stake", "exponent": 0}]}]' /validator/config/genesis.json > /metadata-genesis.json

# a 60 second voting period to allow us to pass governance proposals in the tests
jq '.app_state.gov.voting_params.voting_period = "60s"' /metadata-genesis.json > /edited-genesis.json

mv /edited-genesis.json /genesis.json

# Copy genesis from starting node to host machine for gentx generation
docker cp $STARTING_VALIDATOR_CONTAINER:/validator/config/genesis.json .

# Sets up an arbitrary number of validators on a single machine by docker exec-ing on respective containers
for i in $(seq 1 $NODES);
do
BASELEDGER_HOME="--home /validator"
ARGS="$BASELEDGER_HOME --keyring-backend test"

# Generate a validator key, orchestrator key for each validator
docker exec $STARTING_VALIDATOR_CONTAINER$i $BIN keys add $ARGS validator 2>> /validator-phrases
docker exec $STARTING_VALIDATOR_CONTAINER$i $BIN keys add $ARGS orchestrator 2>> /orchestrator-phrases

VALIDATOR_KEY=$(docker exec $STARTING_VALIDATOR_CONTAINER$i $BIN keys show validator -a $ARGS)
ORCHESTRATOR_KEY=$(docker exec $STARTING_VALIDATOR_CONTAINER$i $BIN keys show orchestrator -a $ARGS)

# move the genesis in
docker cp ./genesis.json $STARTING_VALIDATOR_CONTAINER$i:/validator/config/genesis.json

docker exec $STARTING_VALIDATOR_CONTAINER$i $BIN add-genesis-account $ARGS $VALIDATOR_KEY $ALLOCATION
docker exec $STARTING_VALIDATOR_CONTAINER$i $BIN add-genesis-account $ARGS $ORCHESTRATOR_KEY $ALLOCATION

# # move the genesis back out
docker cp $STARTING_VALIDATOR_CONTAINER$i:/validator/config/genesis.json .

done


# for i in $(seq 1 $NODES);
# do
# cp /genesis.json /validator$i/config/genesis.json
# GAIA_HOME="--home /validator$i"
# ARGS="$GAIA_HOME --keyring-backend test"
# ORCHESTRATOR_KEY=$($BIN keys show orchestrator$i -a $ARGS)
# ETHEREUM_KEY=$(grep address /validator-eth-keys | sed -n "$i"p | sed 's/.*://')
# # the /8 containing 7.7.7.7 is assigned to the DOD and never routable on the public internet
# # we're using it in private to prevent gaia from blacklisting it as unroutable
# # and allow local pex
# $BIN gentx $ARGS $GAIA_HOME --moniker validator$i --chain-id=$CHAIN_ID --ip 7.7.7.$i validator$i 500000000stake $ETHEREUM_KEY $ORCHESTRATOR_KEY
# # obviously we don't need to copy validator1's gentx to itself
# if [ $i -gt 1 ]; then
# cp /validator$i/config/gentx/* /validator1/config/gentx/
# fi
# done


# $BIN collect-gentxs $STARTING_VALIDATOR_HOME
# GENTXS=$(ls /validator1/config/gentx | wc -l)
# cp /validator1/config/genesis.json /genesis.json
# echo "Collected $GENTXS gentx"

# # put the now final genesis.json into the correct folders
# for i in $(seq 1 $NODES);
# do
# cp /genesis.json /validator$i/config/genesis.json
# done
