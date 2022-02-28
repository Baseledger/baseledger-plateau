./baseledgerd keys add --keyring-backend file validator
./baseledgerd keys add --keyring-backend file orchestrator

VAL_KEY=$(./baseledgerd keys show validator -a --keyring-backend file)
ORCH_KEY=$(./baseledgerd keys show orchestrator -a --keyring-backend file)

./baseledgerd  add-genesis-account --keyring-backend file $VAL_KEY 1000000stake
./baseledgerd  add-genesis-account --keyring-backend file $ORCH_KEY 1work
./baseledgerd gentx --keyring-backend file --moniker skos_validator --ip 68.183.221.78 --chain-id=baseledger validator 1000000stake