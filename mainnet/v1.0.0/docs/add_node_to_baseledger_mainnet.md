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
12. Request 1 work token for orchestrator address and staking tokens for validator address by following the step 4 from ./running_a_node_start_here.md. You will not be able to execute step 24 as well as steps after step 26 until you receive the tokens.
14. Adjust the mainnet/v1.0.0/baseledger/cosmovisor.service by adding your node ip here *--p2p.laddr tcp://<your_node_ip>:26656*
15. Adjust the mainnet/v1.0.0/baseledger/cosmovisor.service by adding your keyring password here *Environment=KEYRING_PASSWORD=<your_keyring_password>*
16. Copy the adjusted mainnet/v1.0.0/baseledger/cosmovisor.service to /etc/systemd/system
17. Adjust the mainnet/v1.0.0/orchestrator/orchestrator.service by adding your infura url here *--ethereum-rpc=<your_infura_url>*
18. Adjust the mainnet/v1.0.0/orchestrator/orchestrator.service by adding your Coin market cap api token here *Environment=COINMARKETCAP_API_TOKEN=<your_coinmarketcap_api_token>*
19. Adjust the mainnet/v1.0.0/orchestrator/orchestrator.service by adding your Coin api api token here *Environment=COINAPI_API_TOKEN=<your_coin_api_token>*
20. Copy the adjusted mainnet/v1.0.0/orchestrator/orchestrator.service to /etc/systemd/system
21. Run *systemctl daemon-reload*
22. Run *systemctl start cosmovisor*
23. Verify by running *systemctl status cosmovisor* and *journalctl -u cosmovisor*
24. After tokens received, run *./baseledgerd tx staking create-validator --amount=2000000stake --pubkey=$(./baseledgerd tendermint show-validator) --moniker=<your_moniker> --commission-rate="0" --commission-max-rate="0" --commission-max-change-rate="0" --min-self-delegation="2000000" --from=<your_validator_address> --yes --keyring-backend file*
25. Navigate to /root/.baseledger/orchestrator
26. Run *./baseledger_bridge init*
27. Run *./baseledger_bridge keys set-orchestrator-key --phrase="<orchestrator_mnemonic>"*
28. Run *./baseledger_bridge keys register-orchestrator-address --validator-phrase="<validator_mnemonic>"*
29. Run *systemctl start orchestrator*
30. Verify by running *systemctl status orchestrator* and *journalctl -u orchestrator*
31. IP addresses of the other nodes will be communicated to you from the council. These must be allowed in the firewall rules under port 26656 and all other incoming connections prevented
