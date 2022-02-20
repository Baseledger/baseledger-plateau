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
GRPC_PORT=9090
RPC_PORT=26657
API_PORT=1317
LISTEN_PORT=26655
P2P_PORT=26656
CONTAINER_IP=7.7.7.

docker network create baseledgernet

for i in $(seq 1 $NODES);
do

# add this ip for loopback dialing

ip addr add 7.7.7.$i/32 dev eth0 || true # allowed to fail
docker run --name $VALIDATOR_CONTAINER_BASE_NAME$i $PLATFORM_CMD --cap-add=NET_ADMIN --net baseledgernet -d -p $CONTAINER_IP$i:$GRPC_PORT:9090 -p $CONTAINER_IP$i:$RPC_PORT:26657 -p $CONTAINER_IP$i:$API_PORT:1317 -p $CONTAINER_IP$i:$LISTEN_PORT:26655 -p $CONTAINER_IP$i:$P2P_PORT:26656 baseledger-base

# GRPC_PORT=$(($GRPC_PORT + 1))
# RPC_PORT=$(($RPC_PORT + 1))
# API_PORT=$(($API_PORT + 1))
# LISTEN_PORT=$(($LISTEN_PORT + 10))
# P2P_PORT=$(($P2P_PORT + 10))

done
