#!/bin/bash
set -eux
# your baseledger binary name
BIN=baseledgerd

CHAIN_ID="baseledger"

NODES=${1:-3} 

VALIDATOR_ALLOCATION="10000000000stake,10000000000work"
ORCHESTRATOR_ALLOCATION="1work"
FAUCET_ALLOCATION="10000000000000000000000000000stake,10000000000work"

# first we start a genesis.json with validator 1
# validator 1 will also collect the gentx's once generated
BASELEDGER_HOME="--home /validator"
ARGS="$BASELEDGER_HOME --keyring-backend test"
VALIDATOR_CONTAINER_BASE_NAME="baseledger-validator-container"
STARTING_VALIDATOR_CONTAINER_NAME=$VALIDATOR_CONTAINER_BASE_NAME"1"
docker exec $STARTING_VALIDATOR_CONTAINER_NAME  $BIN init $BASELEDGER_HOME validator --chain-id=$CHAIN_ID


## Modify generated genesis.json to our liking by editing fields using jq
## we could keep a hardcoded genesis file around but that would prevent us from
## testing the generated one with the default values provided by the module.

# add in denom metadata for both native tokens
docker exec $STARTING_VALIDATOR_CONTAINER_NAME jq '.app_state.bank.denom_metadata += [{"name": "Work Token", "symbol": "WRK", "base": "work", "description": "A test token for paying work", "denom_units": [{"denom": "work", "exponent": 0}]},{"name": "Stake Token", "symbol": "STK", "base": "stake", "description": "A staking test token", "denom_units": [{"denom": "stake", "exponent": 0}]}]' /validator/config/genesis.json > /metadata-genesis.json
docker cp /metadata-genesis.json $STARTING_VALIDATOR_CONTAINER_NAME:/metadata-genesis.json

# a 60 second voting period to allow us to pass governance proposals in the tests
docker exec $STARTING_VALIDATOR_CONTAINER_NAME jq '.app_state.gov.voting_params.voting_period = "60s"' /metadata-genesis.json > /edited-genesis.json
docker cp /edited-genesis.json $STARTING_VALIDATOR_CONTAINER_NAME:/edited-genesis.json

docker exec $STARTING_VALIDATOR_CONTAINER_NAME mv /edited-genesis.json /genesis.json

FAUCET_KEY="baseledger1xgs5tamqre7rkz5q7d5fegjsdwufxxvt36w0a0"
docker exec $STARTING_VALIDATOR_CONTAINER_NAME $BIN add-genesis-account $ARGS $FAUCET_KEY $FAUCET_ALLOCATION

# Copy genesis from starting node to host machine for gentx generation
docker cp $STARTING_VALIDATOR_CONTAINER_NAME:/validator/config/genesis.json .

rm -rf ./validator-phrases
rm -rf ./orchestrator-phrases

# Sets up an arbitrary number of validators on a single machine by docker exec-ing on respective containers
for i in $(seq 1 $NODES);
do

# Generate a validator key, orchestrator key for each validator
docker exec $VALIDATOR_CONTAINER_BASE_NAME$i $BIN keys add $ARGS validator 2>> ./validator-phrases
docker exec $VALIDATOR_CONTAINER_BASE_NAME$i $BIN keys add $ARGS orchestrator 2>> ./orchestrator-phrases

VALIDATOR_KEY=$(docker exec $VALIDATOR_CONTAINER_BASE_NAME$i $BIN keys show validator -a $ARGS)
ORCHESTRATOR_KEY=$(docker exec $VALIDATOR_CONTAINER_BASE_NAME$i $BIN keys show orchestrator -a $ARGS)

# move the genesis in
docker cp ./genesis.json $VALIDATOR_CONTAINER_BASE_NAME$i:/validator/config/genesis.json

docker exec $VALIDATOR_CONTAINER_BASE_NAME$i $BIN add-genesis-account $ARGS $VALIDATOR_KEY $VALIDATOR_ALLOCATION
docker exec $VALIDATOR_CONTAINER_BASE_NAME$i $BIN add-genesis-account $ARGS $ORCHESTRATOR_KEY $ORCHESTRATOR_ALLOCATION

# move the genesis back out
docker cp $VALIDATOR_CONTAINER_BASE_NAME$i:/validator/config/genesis.json .

done


for i in $(seq 1 $NODES);
do
# move the genesis in
docker cp ./genesis.json $VALIDATOR_CONTAINER_BASE_NAME$i:/validator/config/genesis.json

VALIDATOR_CONTAINER_IP=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $VALIDATOR_CONTAINER_BASE_NAME$i)

docker exec $VALIDATOR_CONTAINER_BASE_NAME$i $BIN gentx $ARGS --moniker validator$i --chain-id=$CHAIN_ID --ip $VALIDATOR_CONTAINER_IP validator 1000000stake

# copy gentx files to starting validator
if [ $i -gt 1 ]; then
docker cp $VALIDATOR_CONTAINER_BASE_NAME$i:/validator/config/gentx .
docker cp ./gentx $STARTING_VALIDATOR_CONTAINER_NAME:/validator/config/
fi
done

# move the genesis in to starting validator
docker cp ./genesis.json $STARTING_VALIDATOR_CONTAINER_NAME:/validator/config/genesis.json

# collect gentxs in starting validator
docker exec $STARTING_VALIDATOR_CONTAINER_NAME $BIN collect-gentxs $BASELEDGER_HOME
GENTXS=$(docker exec $STARTING_VALIDATOR_CONTAINER_NAME ls /validator/config/gentx | wc -l)

# move the genesis out
docker cp $STARTING_VALIDATOR_CONTAINER_NAME:/validator/config/genesis.json .

echo "Collected $GENTXS gentx"

rm -rf ./gentx

echo "Cleaned host gentx folder"

# put the now final genesis.json into the correct folders of each validator container
for i in $(seq 1 $NODES);
do

# enable 1317 API
docker exec $VALIDATOR_CONTAINER_BASE_NAME$i sed -i 's/enable = false/enable = true/' /validator/config/app.toml

docker cp ./genesis.json  $VALIDATOR_CONTAINER_BASE_NAME$i:/validator/config/genesis.json
done

echo "Cleaned host genesis file"
