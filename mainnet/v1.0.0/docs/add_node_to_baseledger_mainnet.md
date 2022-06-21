# ADD NODE TO BASELEDGER MAINET

## Root user
It is very important to note, that you have to be logged in as `root` user - you can always check it by running `whoami`

## Firewall

IP addresses of the other nodes will be communicated to you by the council. These must be allowed in the firewall rules of your host machine under port 26656. All other incoming connections should be disabled.

Example command with ufw to add a firewall rule: 
ufw allow from <incoming_node_ip> to any port 26656

## Preparing

1. Prepare a https://pro.coinmarketcap.com/ API token. Free plan is ok for now.
2. Prepare a https://www.coinapi.io/ API token. Free plan is ok for now.
3. Prepare a https://infura.io/product/ethereum RPC url. Free plan is ok for now.
4. Copy the mainnet/v1.0.0 folder from the repo to your node

    4.1 A way to do it is through `git clone`:

    4.2 In root; create a test directory `mkdir test` and cd into it.

    4.3 Run `git clone https://github.com/Baseledger/baseledger-plateau.git`

    4.4 Move the `mainnet/v1.0.0` folder to the root with `mv baseledger-plateau/mainnet/ /root/`

    4.5 Navigate to root and remove test directory with `rm -rf test/`

5. Navigate to `mainnet/v1.0.0/docs` 
    
    5.1 To install the node to the home directory of the root user, execute `chmod +x prereqs.sh` and then `bash prereqs.sh` 
    
    5.2 To install the node to a different location, open the file `prereqs_custom_location.sh` 
        5.2.1 On Linux you can use `vim`: From root, run `vim mainnet/v1.0.0/docs/prereqs_custom_location.sh`
        5.2.2 Press `i` and navigate to `<path_where_to_install_the_node>` with arrow keys and replace with your chosen location f.e (/opt/mynode)
        5.2.3 Press `esc` to stop editing
        5.2.4 Type `:wq` to save and close and hit `enter`
        5.2.5 Execute `chmod +x prereqs_custom_location.sh` and then `bash prereqs_custom_location.sh` 


## Setting up accounts

If you have used the `prereqs_custom_location.sh` script to install the node to your custom location, you will have to run every one of the following commands with `/<path_where_to_install_the_node>/.baseledger` f.e (`/opt/mynode/.baseledger/cosmovisor/genesis/bin` instead of `/root/.baseledger/cosmovisor/genesis/bin`) if not told differently after the command!


6. Navigate to `/root/.baseledger/cosmovisor/genesis/bin`
7. Run `/root/.baseledger/cosmovisor/genesis/bin/.baseledgerd init validator --chain-id=baseledger`
8. Delete `genesis.json` in `/root/.baseledger/config/genesis.json` and copy/paste `genesis.json` from `mainnet/v1.0.0/genesis.json`

    8.1 Navigate to `/root/.baseledger/config` and run `rm genesis.json`

    8.2 Navigate to root and run `cp /root/mainnet/v1.0.0/genesis.json /root/.baseledger/config/`
        8.2.1 If you have used the `prereqs_custom_location.sh` script to install the node to your custom location, run `cp /root/mainnet/v1.0.0/genesis.json <path_where_to_install_the_node>/.baseledger/config/` 

9. Run `/root/.baseledger/cosmovisor/genesis/bin/baseledgerd keys add --keyring-backend file validator` and store address and mnemonic in a safe place
    
    9.1 You will be asked to create a keyring-backend password. Make sure to write a strong password and save it somewhere safe

10. Run `/root/.baseledger/cosmovisor/genesis/bin/baseledgerd keys add --keyring-backend file orchestrator` and store address and mnemonic in a safe place

## Requesting tokens

11. Request 1 work token for orchestrator address and staking tokens for validator address by following the step 4 from `./running_a_node_start_here.md`. You will not be able to execute step 24 as well as steps after step 26 until you receive the tokens.

## Starting the node

12. Update `mainnet/v1.0.0/baseledger/cosmovisor.service` by adding the list of persistent peers here `--p2p.persistent_peers <list_of_persistent_peers>` and your node ip here `--p2p.laddr tcp://<your_static_ip>:26656`

    12.1 List of persistent peers will be provided by request from the council f.e (not the actual persistent peer: `12345678910987654321@999.99.999.999:26656`)

    12.2 On Linux you can use `vim`: From root, run `vim mainnet/v1.0.0/baseledger/cosmovisor.service`

    12.3 Press `i` and navigate to `<list_of_persistent_peers>` with arrow keys and replace with the list provided by the council

    12.4 Navigate to `<your_static_ip>` with arrow keys and replace with your ip

13. Update now the keyring password here `Environment=KEYRING_PASSWORD=<your_keyring_password>`

    13.1 When persistent peers, ip and keyring password is added, press `esc` to stop editing

    13.2 Type `:wq` to save and close and hit `enter`

14. If you have used the `prereqs_custom_location.sh` script to install the node to your custom location, replace also the following:
    
    14.1 WorkingDirectory=/<path_where_to_install_the_node>/.baseledger/cosmovisor/
    
    14.2 Environment=DAEMON_HOME=/<path_where_to_install_the_node>/.baseledger
    
    14.3 Environment=KEYRING_DIR=/<path_where_to_install_the_node>/.baseledger

15. Go to root and run `cp mainnet/v1.0.0/baseledger/cosmovisor.service /etc/systemd/system` to copy the adjusted `mainnet/v1.0.0/baseledger/cosmovisor.service` to `/etc/systemd/system`
16. If using `vim`; follow steps as in 14 and 15
17. Update `mainnet/v1.0.0/orchestrator/orchestrator.service` by adding your infura url here `--ethereum-rpc=<infura_link>`
18. Update `mainnet/v1.0.0/orchestrator/orchestrator.service` by adding your Coin market cap api token here `Environment=COINMARKETCAP_API_TOKEN=<COINMARKETCAP_API_TOKEN>`
19. Update `mainnet/v1.0.0/orchestrator/orchestrator.service` by adding your Coin api api token here `Environment=COINAPI_API_TOKEN=<COINAPI_API_TOKEN>`
20. If you have used the `prereqs_custom_location.sh` script to install the node to your custom location, replace also the following:
    
    20.1 WorkingDirectory=/<your_new_path>/.baseledger/orchestrator
    
    20.2 ExecStart=/<your_new_path>/.baseledger/orchestrator/baseledger_bridge orchestrator --ethereum-rpc=<infura_link> --baseledger-contract-address=0xdade4688c10c05716929f91d3005c23d4e233869

21. Go to root and run `cp mainnet/v1.0.0/orchestrator/orchestrator.service /etc/systemd/system` to copy the adjusted `mainnet/v1.0.0/orchestrator/orchestrator.service` to `/etc/systemd/system`
22. Run `systemctl daemon-reload`
23. Run `systemctl start cosmovisor`
24. Verify by running `systemctl status cosmovisor` and `journalctl -u cosmovisor`

    24.1 You can get out of the status and journal with `ctrl + c`

25. Run `systemctl enable cosmovisor` to make sure it is restarted on host restart

## Registering as a validator

26. After tokens received, run `/root/.baseledger/cosmovisor/genesis/bin/baseledgerd tx staking create-validator --amount=2000000stake --pubkey=$(/root/.baseledger/cosmovisor/genesis/bin/baseledgerd tendermint show-validator) --moniker=<your_moniker> --commission-rate="0" --commission-max-rate="0" --commission-max-change-rate="0" --min-self-delegation="2000000" --from=validator --yes --keyring-backend file`. your_moniker is the name of the node which will be visible through the explorer (node_xyz, org_node, mars). You will be prompted to enter the keyring pasword you defined in step 10. Enter it!
27. Verify your validator is added by performing `/root/.baseledger/cosmovisor/genesis/bin/baseledgerd query staking validators`

## Starting the oracle

28. Navigate to /root/.baseledger/orchestrator
29. Run `/root/.baseledger/orchestrator/baseledger_bridge init`
30. Run `/root/.baseledger/orchestrator/baseledger_bridge keys set-orchestrator-key --phrase="<orchestrator_mnemonic>"`. Prepare the command in a text editor and replace <orchestrator_mnemonic>  with the mnemonic you stored when executing step 11.
31. Run `/root/.baseledger/orchestrator/baseledger_bridge keys register-orchestrator-address --validator-phrase="<validator_mnemonic>"`. Prepare the command in a text editor and replace <validator_mnemonic>  with the mnemonic you stored when executing step 10.
32. Run `systemctl start orchestrator`
33. Verify service active by running `systemctl status orchestrator` and `journalctl -u orchestrator`
34. Run `systemctl enable orchestrator` to make sure it is restarted on host restart
