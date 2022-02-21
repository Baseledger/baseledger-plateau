#!/bin/bash
set -eux
# your gaiad binary name
BIN=baseledgerd
VALIDATOR_CONTAINER_BASE_NAME="baseledger-validator-container"
BASELEDGER_HOME="--home /validator"
NODES=3

# We are adding first validator as a persisted peer since we could not get pex to autodiscover with this setup
# (might be an issue with "Cannot add non-routable address fd010b69bae8fb323d9527e377497b93608c11f8@172.24.0.2:26656)
# check by relaxing requirement for routability
FIRST_VALIDATOR_NODE_ID=$(docker exec $VALIDATOR_CONTAINER_BASE_NAME"1" baseledgerd $BASELEDGER_HOME tendermint show-node-id)
FIRST_VALIDATOR_CONTAINER_IP=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $VALIDATOR_CONTAINER_BASE_NAME"1")

for i in $(seq 1 $NODES);
do
    rm -rf /validator$i
    mkdir /validator$i

    # this implicitly caps us at ~6000 nodes for this sim
    # note that we start on 26656 the idea here is that the first
    # node (node 1) is at the expected contact address from the gentx
    # faciliating automated peer exchange
    RPC_ADDRESS="--rpc.laddr tcp://0.0.0.0:26657"
    GRPC_ADDRESS="--grpc.address 0.0.0.0:9090"
    LISTEN_ADDRESS="--address tcp://0.0.0.0:26655"
    P2P_ADDRESS="--p2p.laddr tcp://0.0.0.0:26656"
    PERSISTENT_PEERS="--p2p.persistent_peers $FIRST_VALIDATOR_NODE_ID@$FIRST_VALIDATOR_CONTAINER_IP:26656"
    LOG_LEVEL="--log_level info"
    INVARIANTS_CHECK="--inv-check-period 1"
    ARGS="$BASELEDGER_HOME $LISTEN_ADDRESS $RPC_ADDRESS $GRPC_ADDRESS $LOG_LEVEL $INVARIANTS_CHECK $P2P_ADDRESS $PERSISTENT_PEERS"
    
    docker exec $VALIDATOR_CONTAINER_BASE_NAME$i $BIN $ARGS start &> /validator$i/vallogs &

    FEES="--fees '0token'"
    ETH_RPC="--ethereum-rpc='http://localhost:8545'"
    DEPOSIT_CONTRACT_ADDRESS="--baseledger-contract-address='<BASELEDGER_TEST_CONTRACT_ADDRESS>'"
    # TODO: ADD COIN PRICE APIs 

    docker exec --workdir /baseledger/orchestrator $VALIDATOR_CONTAINER_BASE_NAME$i cargo run -- orchestrator $FEES $ETH_RPC $DEPOSIT_CONTRACT_ADDRESS &> /validator$i/orclogs &
done


