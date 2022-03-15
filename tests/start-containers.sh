#!/bin/bash
set -eux

# the directory of this script, useful for allowing this script
# to be run with any PWD
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
VALIDATOR_CONTAINER_BASE_NAME="baseledger-validator-container"
ETHEREUM_CONTAINER_NAME="baseledger-ethereum-node"
NODES=${1:-3}    
# Remove existing container instance
set +e
for i in $(seq 1 $NODES);
do

docker rm -f $VALIDATOR_CONTAINER_BASE_NAME$i

done
set -e

pushd $DIR/../

# setup for Mac M1 compatibility
PLATFORM_CMD=""
if [[ "$OSTYPE" == "darwin"* ]]; then
    if [[ -n $(sysctl -a | grep brand | grep "M1") ]]; then
       echo "Setting --platform=linux/amd64 for Mac M1 compatibility"
       PLATFORM_CMD="--platform=linux/amd64"; fi
fi

# Run new test container instances

GRPC_PORT=9090
RPC_PORT=26657
API_PORT=1317
P2P_PORT=26656

docker network create baseledgernet

for i in $(seq 1 $NODES);
do

docker run --name $VALIDATOR_CONTAINER_BASE_NAME$i $PLATFORM_CMD --net baseledgernet -d --expose $GRPC_PORT --expose $RPC_PORT --expose $API_PORT --expose $P2P_PORT --publish $(($API_PORT + $i - 1)):$API_PORT --publish $(($RPC_PORT + $i - 1)):$RPC_PORT  --publish $(($GRPC_PORT + $i - 1)):$GRPC_PORT baseledger-base

done

# Assumes that the baseledger-hardhat has been built by following instructions in the baseledger-contracts repo readme file
docker run --name $ETHEREUM_CONTAINER_NAME --net baseledgernet -d -p 8545:8545 baseledger-hardhat
