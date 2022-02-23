#!/bin/bash
set -eux

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DOCKERFOLDER=$DIR/dockerfile
REPOFOLDER=$DIR/..

# change our directory so that the git archive command works as expected
pushd $REPOFOLDER

# Build base container
git archive --format=tar.gz -o $DOCKERFOLDER/baseledger.tar.gz --prefix=baseledger/ HEAD
pushd $DOCKERFOLDER

# setup for Mac M1 Compatibility 
PLATFORM_CMD=""
if [[ "$OSTYPE" == "darwin"* ]]; then
    if [[ -n $(sysctl -a | grep brand | grep "M1") ]]; then
       echo "Setting --platform=linux/amd64 for Mac M1 compatibility"
       PLATFORM_CMD="--platform=linux/amd64"; fi
fi
docker build -t baseledger-base $PLATFORM_CMD .
