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

echo “Prepare folder structure”

mkdir <path_where_to_install_the_node>/.baseledger
ln -s <path_where_to_install_the_node>/.baseledger $HOME/.baseledger


echo "Installing Cosmovisor"

go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@v1.0.0
export DAEMON_HOME=$HOME/.baseledger

echo "Creating Cosmovisor folders"

mkdir -p $DAEMON_HOME/cosmovisor/genesis
mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
mkdir -p $DAEMON_HOME/cosmovisor/upgrades
ln -s $DAEMON_HOME/cosmovisor/genesis/ $DAEMON_HOME/cosmovisor/current
mkdir -p $DAEMON_HOME/data

echo "Moving baseledgerd binary to $DAEMON_HOME/cosmovisor/genesis/bin"

tar -C $DAEMON_HOME/cosmovisor/genesis/bin -xzf ../baseledger/baseledger_linux_amd64.tar.gz

echo "Execute permission on baseledgerd binary"

chmod +x $DAEMON_HOME/cosmovisor/genesis/bin/baseledgerd

echo "Creating Orchestrator folders"

mkdir -p $DAEMON_HOME/orchestrator

echo "Moving baseledger_bridge binary to $DAEMON_HOME/orchestrator"

mv ../orchestrator/baseledger_bridge $DAEMON_HOME/orchestrator

echo "Execute permission on baseledger_bridge binary"

chmod +x $DAEMON_HOME/orchestrator/baseledger_bridge

echo "Execute permission on DAEMON_HOME"

chmod -R 766 <path_where_to_install_the_node>

