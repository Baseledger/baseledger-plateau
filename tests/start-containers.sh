#!/bin/bash
set -eux

# the directory of this script, useful for allowing this script
# to be run with any PWD
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
VALIDATOR_CONTAINER_BASE_NAME="baseledger-validator-container"
NODES=3

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
LISTEN_PORT=26655
P2P_PORT=26656

docker network create baseledgernet

for i in $(seq 1 $NODES);
do

docker run --name $VALIDATOR_CONTAINER_BASE_NAME$i $PLATFORM_CMD --net baseledgernet -d --expose $GRPC_PORT --expose $RPC_PORT --expose $API_PORT --expose $LISTEN_PORT --expose $P2P_PORT baseledger-base

done
