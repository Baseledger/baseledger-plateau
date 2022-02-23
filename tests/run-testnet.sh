#!/bin/bash
set -eux

BIN=baseledgerd
VALIDATOR_CONTAINER_BASE_NAME="baseledger-validator-container"
ETHEREUM_CONTAINER_NAME="baseledger-ethereum-node"
BASELEDGER_HOME="--home /validator"
NODES=3

# We are adding first validator as a persisted peer since we could not get pex to autodiscover with this setup
# (might be an issue with "Cannot add non-routable address fd010b69bae8fb323d9527e377497b93608c11f8@172.24.0.2:26656)
# check by relaxing requirement for routability
FIRST_VALIDATOR_NODE_ID=$(docker exec $VALIDATOR_CONTAINER_BASE_NAME"1" baseledgerd $BASELEDGER_HOME tendermint show-node-id)
FIRST_VALIDATOR_CONTAINER_IP=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $VALIDATOR_CONTAINER_BASE_NAME"1")

ETHEREUM_CONTAINER_NAME_IP=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $ETHEREUM_CONTAINER_NAME)

for i in $(seq 1 $NODES);
do
    rm -rf /validator$i
    mkdir /validator$i

    RPC_ADDRESS="--rpc.laddr tcp://0.0.0.0:26657"
    GRPC_ADDRESS="--grpc.address 0.0.0.0:9090"
    LISTEN_ADDRESS="--address tcp://0.0.0.0:26655"
    P2P_ADDRESS="--p2p.laddr tcp://0.0.0.0:26656"
    PERSISTENT_PEERS="--p2p.persistent_peers $FIRST_VALIDATOR_NODE_ID@$FIRST_VALIDATOR_CONTAINER_IP:26656"
    LOG_LEVEL="--log_level info"
    INVARIANTS_CHECK="--inv-check-period 1"
    ARGS="$BASELEDGER_HOME $LISTEN_ADDRESS $RPC_ADDRESS $GRPC_ADDRESS $LOG_LEVEL $INVARIANTS_CHECK $P2P_ADDRESS $PERSISTENT_PEERS"
    
    docker exec $VALIDATOR_CONTAINER_BASE_NAME$i $BIN $ARGS start &> /validator$i/vallogs &
done

sleep 10

for i in $(seq 1 $NODES);
do    
    # phrases are located on every 6th line
    y=$(( 6*$i ))

    VALIDATOR_PHRASE=$(sed "$y q;d" ./validator-phrases)
    ORCHESTRATOR_PHRASE=$(sed "$y q;d" ./orchestrator-phrases)

    docker exec --workdir /baseledger/orchestrator $VALIDATOR_CONTAINER_BASE_NAME$i cargo run -- keys set-orchestrator-key --phrase="$ORCHESTRATOR_PHRASE"
    
    docker exec --workdir /baseledger/orchestrator $VALIDATOR_CONTAINER_BASE_NAME$i cargo run -- keys register-orchestrator-address --validator-phrase="$VALIDATOR_PHRASE"
    
    ETH_RPC="--ethereum-rpc=http://$ETHEREUM_CONTAINER_NAME_IP:8545"
    DEPOSIT_CONTRACT_ADDRESS="--baseledger-contract-address=0xe7f1725e7734ce288f8367e1bb143e90bb3f0512"
    
    # TODO: API tokens for price
    docker exec --workdir /baseledger/orchestrator -e COINMARKETCAP_API_TOKEN=asd -e COINAPI_API_TOKEN=asd $VALIDATOR_CONTAINER_BASE_NAME$i cargo run -- orchestrator $ETH_RPC $DEPOSIT_CONTRACT_ADDRESS &> /validator$i/orclogs &
done


