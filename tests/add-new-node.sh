#!/bin/bash
set -eux
# setup for Mac M1 compatibility
PLATFORM_CMD=""
if [[ "$OSTYPE" == "darwin"* ]]; then
    if [[ -n $(sysctl -a | grep brand | grep "M1") ]]; then
       echo "Setting --platform=linux/amd64 for Mac M1 compatibility"
       PLATFORM_CMD="--platform=linux/amd64"; fi
fi

CHAIN_ID="baseledger"
BIN=baseledgerd
VALIDATOR_CONTAINER_BASE_NAME="baseledger-validator-container"
ETHEREUM_CONTAINER_NAME="baseledger-ethereum-node"
BASELEDGER_HOME="--home /validator"
NODE_ID=${1-4}
# We are adding first validator as a persisted peer since we could not get pex to autodiscover with this setup
# (might be an issue with "Cannot add non-routable address fd010b69bae8fb323d9527e377497b93608c11f8@172.24.0.2:26656)
# check by relaxing requirement for routability
FIRST_VALIDATOR_NODE_ID=$(docker exec $VALIDATOR_CONTAINER_BASE_NAME"1" baseledgerd $BASELEDGER_HOME tendermint show-node-id)
FIRST_VALIDATOR_CONTAINER_IP=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $VALIDATOR_CONTAINER_BASE_NAME"1")

ETHEREUM_CONTAINER_NAME_IP=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $ETHEREUM_CONTAINER_NAME)

FAUCET_KEY="baseledger1p8x9ud2m75dmufevmrym3uak0hgcrw58h6n872"

rm -rf /validator$NODE_ID
mkdir /validator$NODE_ID

GRPC_PORT=9090
RPC_PORT=26657
API_PORT=1317
LISTEN_PORT=26655
P2P_PORT=26656

docker run --name $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID $PLATFORM_CMD --net baseledgernet -d --expose $GRPC_PORT --expose $RPC_PORT --expose $API_PORT --expose $LISTEN_PORT --expose $P2P_PORT --publish $(($API_PORT + $NODE_ID - 1)):$API_PORT --publish $(($RPC_PORT + $NODE_ID - 1)):$RPC_PORT  --publish $(($GRPC_PORT + $NODE_ID - 1)):$GRPC_PORT baseledger-base

docker exec $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID $BIN init $BASELEDGER_HOME validator --chain-id=$CHAIN_ID

RPC_ADDRESS="--rpc.laddr tcp://0.0.0.0:26657"
GRPC_ADDRESS="--grpc.address 0.0.0.0:9090"
LISTEN_ADDRESS="--address tcp://0.0.0.0:26655"
P2P_ADDRESS="--p2p.laddr tcp://0.0.0.0:26656"
PERSISTENT_PEERS="--p2p.persistent_peers $FIRST_VALIDATOR_NODE_ID@$FIRST_VALIDATOR_CONTAINER_IP:26656"
LOG_LEVEL="--log_level info"
INVARIANTS_CHECK="--inv-check-period 1"
ARGS="$BASELEDGER_HOME $LISTEN_ADDRESS $RPC_ADDRESS $GRPC_ADDRESS $LOG_LEVEL $INVARIANTS_CHECK $P2P_ADDRESS $PERSISTENT_PEERS"

# copy genesis
docker cp $VALIDATOR_CONTAINER_BASE_NAME"1":/validator/config/genesis.json .
docker cp ./genesis.json $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID:/validator/config/genesis.json

# enable api
docker exec $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID sed -i 's/enable = false/enable = true/' /validator/config/app.toml

# Generate a validator key, orchestrator key for new validator
docker exec $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID $BIN keys add validator --keyring-backend test 2>> ./validator-phrases-new
docker exec $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID $BIN keys add orchestrator --keyring-backend test 2>> ./orchestrator-phrases-new

VALIDATOR_KEY=$(docker exec $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID $BIN keys show validator -a --keyring-backend test)
ORCHESTRATOR_KEY=$(docker exec $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID $BIN keys show orchestrator -a --keyring-backend test)

# send 2 mil stake to new validator
docker exec $VALIDATOR_CONTAINER_BASE_NAME"1" baseledgerd tx bank send $FAUCET_KEY $VALIDATOR_KEY 2000000stake --yes --keyring-backend test
# sleep because account sequence is not updated properly
sleep 5
# send 1 work to orch
docker exec $VALIDATOR_CONTAINER_BASE_NAME"1" baseledgerd tx bank send $FAUCET_KEY $ORCHESTRATOR_KEY 1work --yes --keyring-backend test

docker exec $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID $BIN $ARGS start &> /validator$NODE_ID/vallogs &

sleep 10

# create-validator with min self delegation 2 mil stake
PUB_KEY=$(docker exec $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID baseledgerd $BASELEDGER_HOME tendermint show-validator)
docker exec $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID baseledgerd tx staking create-validator --amount=2000000stake --pubkey=$PUB_KEY --moniker=new_node --commission-rate="0.10" --commission-max-rate="0.20" --commission-max-change-rate="0.01" --min-self-delegation="2000000" --from=$VALIDATOR_KEY --yes --keyring-backend test

sleep 5
# phrases are located each 6th line
VALIDATOR_PHRASE=$(sed "6 q;d" ./validator-phrases-new)
ORCHESTRATOR_PHRASE=$(sed "6 q;d" ./orchestrator-phrases-new)

docker exec --workdir /baseledger/orchestrator $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID cargo run -- keys set-orchestrator-key --phrase="$ORCHESTRATOR_PHRASE"

docker exec --workdir /baseledger/orchestrator $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID cargo run -- keys register-orchestrator-address --validator-phrase="$VALIDATOR_PHRASE"

ETH_RPC="--ethereum-rpc=http://$ETHEREUM_CONTAINER_NAME_IP:8545"
DEPOSIT_CONTRACT_ADDRESS="--baseledger-contract-address=0xe7f1725e7734ce288f8367e1bb143e90bb3f0512"

# TODO: API tokens for price
# FOR trace logging add -e RUST_LOG="trace" to line bellow
docker exec --workdir /baseledger/orchestrator -e COINMARKETCAP_API_TOKEN=asd -e COINAPI_API_TOKEN=asd $VALIDATOR_CONTAINER_BASE_NAME$NODE_ID cargo run -- orchestrator $ETH_RPC $DEPOSIT_CONTRACT_ADDRESS &> /validator$NODE_ID/orclogs &

rm -rf ./validator-phrases-new
rm -rf ./orchestrator-phrases-new
rm -rf ./genesis.json