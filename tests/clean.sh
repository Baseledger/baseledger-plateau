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

docker network rm baseledgernet

set -e