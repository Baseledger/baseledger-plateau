# ADD NODE TO BASELEDGER MAINET

1.  Prepare a https://pro.coinmarketcap.com/ API token. Free plan is ok.
2.  Prepare a https://www.coinapi.io/ API token. Free plan is ok.
3.  Prepare a https://infura.io/product/ethereum RPC url. Free plan is ok.
4.  Copy the mainnet/v1.0.0 folder from the repo to your node
5.  Navigate to mainnet/v1.0.0/docs on the node
6.  Run prereqs.sh
7.  Navigate to /root/.baseledger/cosmovisor/genesis/bin
8.  Run *./baseledgerd init validator --chain-id=baseledger*
9.  Copy the mainnet/v1.0.0/genesis.json to /root/.baseledger/config/genesis.json
10.  Run *./baseledgerd keys add --keyring-backend file validator* and store address and  mnemonic in a safe place
    10.1 You will be asked to create a keyring-backend password. Make sure to write a strong password and save it somewhere safe
11. Run *./baseledgerd keys add --keyring-backend file orchestrator* and store address and mnemonic in a safe place
12. Request 1 work token for orchestrator from the faucet. Command 26 can happen only after this account has 1work balance.
13. Adjust the mainnet/v1.0.0/baseledger/cosmovisor.service by adding your node ip here *--p2p.laddr tcp://<your_node_ip>:26656*
14. Adjust the mainnet/v1.0.0/baseledger/cosmovisor.service by adding your keyring password here *Environment=KEYRING_PASSWORD=<your_keyring_password>*
15. Copy the adjusted mainnet/v1.0.0/baseledger/cosmovisor.service to /etc/systemd.system
16. Adjust the mainnet/v1.0.0/orchestrator/orchestrator.service by adding your infura url here *--ethereum-rpc=<your_infura_url>*
17. Adjust the mainnet/v1.0.0/orchestrator/orchestrator.service by adding your Coin market cap api token here *Environment=COINMARKETCAP_API_TOKEN=<your_coinmarketcap_api_token>*
18. Adjust the mainnet/v1.0.0/orchestrator/orchestrator.service by adding your Coin api api token here *Environment=COINAPI_API_TOKEN=<your_coin_api_token>*
19. Copy the adjusted mainnet/v1.0.0/orchestrator/orchestrator.service to /etc/systemd.system
20. Run *systemctl daemon-reload*
21. Run *systemctl start cosmovisor*
22. Verify by running *systemctl status cosmovisor* and *journalctl -u cosmovisor*
23. Navigate to /root/.baseledger/orchestrator
24. Run *./baseledger_bridge init*
25. Run *./baseledger_bridge keys set-orchestrator-key --phrase="<orchestrator_mnemonic>"*
26. Run *./baseledger_bridge keys register-orchestrator-address --validator-phrase="<validator_mnemonic>"*
27. Run *systemctl start orchestrator*
28. Verify by running *systemctl status orchestrator* and *journalctl -u orchestrator*
