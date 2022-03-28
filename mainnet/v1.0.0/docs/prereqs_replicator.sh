#!/bin/sh

echo "Installing Golang"
wget https://golang.org/dl/go1.16.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz
chown -R root:root /usr/local/go
mkdir -p $HOME/go/{bin,src}
echo "export GOPATH=$HOME/go" >> ~/.profile 
echo "export PATH=$PATH:$GOPATH/bin" >> ~/.profile
echo "export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin" >> ~/.profile
. ~/.profile

echo "Installing Cosmovisor"

go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@v1.0.0
export DAEMON_HOME=$HOME/.baseledger

echo "Creating Cosmovisor folders"

mkdir -p $DAEMON_HOME/cosmovisor/genesis
mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
mkdir -p $DAEMON_HOME/cosmovisor/upgrades
mkdir -p $DAEMON_HOME/data

echo "Moving baseledgerd binary to $DAEMON_HOME/cosmovisor/genesis/bin"

tar -C $DAEMON_HOME/cosmovisor/genesis/bin -xzf ../baseledger/baseledger_linux_amd64.tar.gz

echo "Execute permission on baseledgerd binary"

chmod +x $DAEMON_HOME/cosmovisor/genesis/bin/baseledgerd
