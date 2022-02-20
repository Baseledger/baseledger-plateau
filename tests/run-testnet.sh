#!/bin/bash
set -eux
# your gaiad binary name
BIN=baseledgerd
VALIDATOR_CONTAINER_BASE_NAME="baseledger-validator-container"
NODES=3

for i in $(seq 1 $NODES);
do
    BASELEDGER_HOME="--home /validator"
    # this implicitly caps us at ~6000 nodes for this sim
    # note that we start on 26656 the idea here is that the first
    # node (node 1) is at the expected contact address from the gentx
    # faciliating automated peer exchange
    RPC_ADDRESS="--rpc.laddr tcp://0.0.0.0:26657"
    GRPC_ADDRESS="--grpc.address 0.0.0.0:9090"
    LISTEN_ADDRESS="--address tcp://0.0.0.0:26655"
    P2P_ADDRESS="--p2p.laddr tcp://0.0.0.0:26656"
    LOG_LEVEL="--log_level info"
    INVARIANTS_CHECK="--inv-check-period 1"
    ARGS="$BASELEDGER_HOME $LISTEN_ADDRESS $RPC_ADDRESS $GRPC_ADDRESS $LOG_LEVEL $INVARIANTS_CHECK $P2P_ADDRESS"
    rm -rf /validator$i
    mkdir /validator$i && touch /validator$i/logs
    docker exec $VALIDATOR_CONTAINER_BASE_NAME$i $BIN $ARGS start &> /validator$i/logs &
done


