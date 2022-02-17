#!/bin/bash
set -eux

# the directory of this script, useful for allowing this script
# to be run with any PWD
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
STARTING_VALIDATOR_CONTAINER="baseledger-validator-container"
NODES=2

# Remove existing container instance
set +e
docker rm -f $STARTING_VALIDATOR_CONTAINER
for i in $(seq 1 $NODES);
do
docker rm -f $STARTING_VALIDATOR_CONTAINER$i
done
set -e