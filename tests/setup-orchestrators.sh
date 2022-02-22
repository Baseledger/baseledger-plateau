#!/bin/bash
set -eux

NODES=3

# first we start a genesis.json with validator 1
# validator 1 will also collect the gentx's once gnerated
ORCHESTRATOR_HOME="--home /orchestrator"
VALIDATOR_CONTAINER_BASE_NAME="baseledger-validator-container"
STARTING_VALIDATOR_CONTAINER_NAME=$VALIDATOR_CONTAINER_BASE_NAME"1"

# Sets up an arbitrary number of validators on a single machine by docker exec-ing on respective containers
for i in $(seq 1 $NODES);
do

# Init each orchestrator and setup configuration
# docker exec --workdir /baseledger/orchestrator $VALIDATOR_CONTAINER_BASE_NAME$i cargo run -- init

done
