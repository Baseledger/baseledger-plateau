use crate::args::OrchestratorOpts;
use crate::config::config_exists;
use crate::config::load_keys;
use deep_space::PrivateKey as CosmosPrivateKey;
use utils::connection_prep::{
    check_validator_address, wait_for_cosmos_node_ready,
};
use utils::connection_prep::{create_rpc_connections};
use ethereum_oracle::main_loop::eth_oracle_main_loop;
use ethereum_oracle::main_loop::{ETH_ORACLE_LOOP_SPEED};
use std::path::Path;
use std::process::exit;
use utils::connection_prep::check_for_fee;

pub async fn orchestrator(
    args: OrchestratorOpts,
    address_prefix: String,
    home_dir: &Path,
) {
    let fee = args.fees;
    let cosmos_grpc = args.cosmos_grpc;
    let ethereum_rpc = args.ethereum_rpc;
    let cosmos_key = args.cosmos_phrase;

    let cosmos_key = if let Some(k) = cosmos_key {
        k
    } else {
        let mut k = None;
        if config_exists(home_dir) {
            let keys = load_keys(home_dir);
            if let Some(stored_key) = keys.orchestrator_phrase {
                k = Some(CosmosPrivateKey::from_phrase(&stored_key, "").unwrap())
            }
        }
        if k.is_none() {
            error!("You must specify an Orchestrator key phrase!");
            error!("To set an already registered key use 'baseledger_bridge keys set-orchestrator-key --phrase \"your phrase\"`");
            error!("To run from the command line, with no key storage use 'baseledger_bridge orchestrator --cosmos-phrase \"your phrase\"' ");
            error!("If you have not already generated a key 'baseledger_bridge keys register-orchestrator-address' will generate one for you");
            exit(1);
        }
        k.unwrap()
    };

    let timeout = ETH_ORACLE_LOOP_SPEED;

    trace!("Probing RPC connections");
    // probe all rpc connections and see if they are valid
    let connections = create_rpc_connections(
        address_prefix,
        Some(cosmos_grpc),
        Some(ethereum_rpc),
        timeout,
    )
    .await;

    let mut grpc = connections.grpc.clone().unwrap();
    let contact = connections.contact.clone().unwrap();
    // let web3 = connections.web3.clone().unwrap();

    let public_cosmos_key = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    info!("Starting Baseledger bridge");
    info!(
        "Cosmos Address {}",
        public_cosmos_key
    );

    // check if the cosmos node is syncing, if so wait for it
    // we can't move any steps above this because they may fail on an incorrect
    // historic chain state while syncing occurs
    wait_for_cosmos_node_ready(&contact).await;

    // check if the delegate addresses are correctly configured
    check_validator_address(
        &mut grpc,
        public_cosmos_key,
        &contact.get_prefix(),
    )
    .await;

    // check if we actually have the promised balance of tokens to pay fees
    check_for_fee(&fee, public_cosmos_key, &contact).await;

    let contract_address = if let Some(c) = args.baseledger_contract_address {
        c
    } else {
        error!("The Baseledger contract address is not yet set as a chain parameter! You must specify --baseledger-contract-address");
        exit(1);
    };

    eth_oracle_main_loop(
        cosmos_key,
        connections.web3.unwrap().clone(),
        connections.contact.unwrap().clone(),
        connections.grpc.unwrap().clone(),
        contract_address,
        fee.clone(),
    )
    .await;
}
