#!/bin/sh

echo "Installing Golang"
sudo wget https://golang.org/dl/go1.16.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz
sudo chown -R root:root /usr/local/go
mkdir -p $HOME/go/{bin,src}
echo "export GOPATH=$HOME/go" >> ~/.profile 
echo "export PATH=$PATH:$GOPATH/bin" >> ~/.profile
echo "export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin" >> ~/.profile
. ~/.profile

echo "Installing Cosmovisor"

go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@v1.0.0
export DAEMON_HOME=$HOME/.baseledger
export DAEMON_NAME=baseledgerd

echo "Creating Cosmovisor folders"

mkdir -p $DAEMON_HOME/cosmovisor/current
mkdir -p $DAEMON_HOME/cosmovisor/genesis
mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin/$DAEMON_NAME
mkdir -p $DAEMON_HOME/cosmovisor/upgrades/<replace_with_name>
mkdir -p $DAEMON_HOME/cosmovisor/upgrades/bin/$DAEMON_NAME
mkdir -p $DAEMON_HOME/data
echo "{}" > $DAEMON_HOME/cosmovisor/upgrades/<replace_with_name>/upgrade-info.json

echo "Creating Cosmovisor user"

TODO

echo "Assign folders to Cosmovisor user"

TODO

