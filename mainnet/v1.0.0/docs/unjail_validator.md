 # UNJAIL A VALIDATOR NODE


After jailing, your node will be slashed by a certain amount of stake tokens. This amount can be sent over to your validator account by the council in case that the jailing reason is not something you had control over.

If the amount of stake tokens of your validator is bellow the minimum of 2000000stake you will first need to delegate the amount of tokens needed to reach that minimum:
 
 ./baseledgerd tx staking delegate <your_valoper_address> 20000stake --from=validator --yes --keyring-backend file

 then you execute:

 ./baseledgerd tx slashing unjail --from=validator --yes --keyring-backend file