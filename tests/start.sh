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

## TODO: We might be missing ports for P2P - check if fails to discover other nodes
GRPC_PORT=9089
RPC_PORT=26656
API_PORT=1316

for i in $(seq 1 $NODES);
do

GRPC_PORT=$(($GRPC_PORT + 1))
RPC_PORT=$(($RPC_PORT + 1))
API_PORT=$(($API_PORT + 1))

docker run --name $VALIDATOR_CONTAINER_BASE_NAME$i $PLATFORM_CMD -d -p $GRPC_PORT:9090 -p $RPC_PORT:26657 -p $API_PORT:1317 baseledger-base
done
