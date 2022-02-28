./baseledgerd init validator --chain-id=baseledger
./baseledgerd keys add --keyring-backend file validator
./baseledgerd keys add --keyring-backend file orchestrator

VAL_KEY=$(./baseledgerd keys show validator -a --keyring-backend file)
ORCH_KEY=$(./baseledgerd keys show orchestrator -a --keyring-backend file)

./baseledgerd add-genesis-account --keyring-backend file baseledger1xgs5tamqre7rkz5q7d5fegjsdwufxxvt36w0a0 10000000000stake,10000000000work
./baseledgerd  add-genesis-account --keyring-backend file $VAL_KEY 1000000stake
./baseledgerd  add-genesis-account --keyring-backend file $ORCH_KEY 1work
./baseledgerd gentx --keyring-backend file --moniker finspot_validator --ip 68.183.221.78 --chain-id=baseledger validator 1000000stake
