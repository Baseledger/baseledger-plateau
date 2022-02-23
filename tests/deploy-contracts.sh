#!/bin/bash
set -eux

ETHEREUM_CONTAINER_NAME="baseledger-ethereum-node"

docker exec --workdir /usr/src/app $ETHEREUM_CONTAINER_NAME npm run compile
docker exec --workdir /usr/src/app $ETHEREUM_CONTAINER_NAME npm run contracts:migrate:local

done
