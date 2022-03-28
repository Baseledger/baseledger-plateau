# ADD REPLICATING NODE TO BASELEDGER MAINET

## Root user
It is very important to note, that you have to be logged in as `root` user - you can always check it by running `whoami`

## Firewall

IP addresses of the other nodes will be communicated to you by the council. These must be allowed in the firewall rules of your host machine under port 26656. All other incoming connections should be disabled.

Example command with ufw to add a firewall rule: 
ufw allow from <incoming_node_ip> to any port 26656

## Preparing

1. Copy the mainnet/v1.0.0 folder from the repo to your node

    4.1 A way to do it is through `git clone`:

    4.2 In root; create a test directory `mkdir test` and cd into it.

    4.3 Run `git clone https://github.com/Baseledger/baseledger-plateau.git`

    4.4 Move the `mainnet/v1.0.0` folder to the root with `mv baseledger-plateau/mainnet/ /root/`

    4.5 Navigate to root and remove test directory with `rm -rf test/`

2. Navigate to `mainnet/v1.0.0/docs`
3. Execute `chmod +x prereqs_replicator.sh` and then `bash prereqs_replicator.sh`


## Setting up accounts

4. Navigate to `/root/.baseledger/cosmovisor/genesis/bin`
5. Run `/root/.baseledger/cosmovisor/genesis/bin/.baseledgerd init <moniker> --chain-id=baseledger`
6. Delete `genesis.json` in `/root/.baseledger/config/genesis.json` and copy/paste `genesis.json` from `mainnet/v1.0.0/genesis.json`

    6.1 Navigate to `/root/.baseledger/config` and run `rm genesis.json`

    6.2 Navigate to root and run `cp /root/mainnet/v1.0.0/genesis.json /root/.baseledger/config/`

7. Run `/root/.baseledger/cosmovisor/genesis/bin/baseledgerd keys add --keyring-backend file <key-name>` and store address and mnemonic in a safe place. This account will be used to store work tokens that you can use to drop proofs.
    
    7.1 You will be asked to create a keyring-backend password. Make sure to write a strong password and save it somewhere safe

## Starting the node

8. Update `mainnet/v1.0.0/baseledger/cosmovisor.service` by adding the list of persistent peers here `--p2p.persistent_peers <list_of_persistent_peers>` and your node ip here `--p2p.laddr tcp://<your_static_ip>:26656`

    14.1 List of persistent peers can be provided by request from the council

    14.2 On Linux you can use `vim`: From root, run `vim mainnet/v1.0.0/baseledger/cosmovisor.service`

    14.3 Press `i` and navigate to `<list_of_persistent_peers>` with arrow keys and replace with the list provided by the council

    14.4 Navigate to `<your_static_ip>` with arrow keys and replace with your ip

9. Update now the keyring password here `Environment=KEYRING_PASSWORD=<your_keyring_password>`

    9.1 When persistent peers, ip and keyring password is added, press `esc` to stop editing

    9.2 Type `:wq` to save and close and hit `enter`

10. Go to root and run `cp mainnet/v1.0.0/baseledger/cosmovisor.service /etc/systemd/system` to copy the adjusted `mainnet/v1.0.0/baseledger/cosmovisor.service` to `/etc/systemd/system`
11. Run `systemctl daemon-reload`
12. Run `systemctl start cosmovisor`
13. Verify by running `systemctl status cosmovisor` and `journalctl -u cosmovisor`

    23.1 You can get out of the status and journal with `ctrl + c`

