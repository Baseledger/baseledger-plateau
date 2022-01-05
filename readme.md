poc-contract is a folder with contracts, contains erc-20 and test contract for deposit that is throwing event.

Folder orchestrator is simplified copied from gravity-bridge (https://github.com/Gravity-Bridge/Gravity-Bridge), i tried to not use anything that is not needed for first version.

Orchestrator folder structure (kept same names as gravity, we will rename when we only left what’s needed and maybe change structure a bit)

cosmos_gravity - utils for interaction with cosmos chain

gbt - cli program to start various parts of orchestrator (i only left starting of event listening loop)

gravity_proto - this is auto generated from cosmos proto files, i just copied this from gravity repo, do not review this at the moment, our will hopefully be much smaller, just type definitions, for example i needed this for grpc client types

gravity_utils - various utils needed, only copied some (util methods for connection to ethereum and cosmos, types, error types etc)

orchestrator - main thing for us atm, there is loop listening to ethereum events, it is started as cli command

Biggest difference to gravity bridge is that gravity is doing more things not needed by us (at least atm) like relaying stuff to ethereum, transaction batching to ethereum, listening to all those events (we only listen to SendToCosmosEvent so everything else is removed) etc.

To make this work locally apart from starport rust is needed to be installed and then call

cargo build --all in root of orchestrator folder
navigate to gbt and execute

cargo run -- orchestrator --fees "1token" --cosmos-phrase=<mnemonic-from-local-starport> --ethereum-key=<ethereum-private-key-from-keepass-mnemonic> --ethereum-rpc=<infura-kee-pass-url> --gravity-contract-address="0x9e7144C01e3B1D8f3E8127a0C4769637eBac01EA"