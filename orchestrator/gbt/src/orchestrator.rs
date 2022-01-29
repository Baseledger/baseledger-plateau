use crate::args::OrchestratorOpts;
use crate::config::config_exists;
use crate::config::load_keys;
use deep_space::PrivateKey as CosmosPrivateKey;
use utils::connection_prep::{
    check_validator_address, wait_for_cosmos_node_ready,
};
use utils::connection_prep::{create_rpc_connections};
use orchestrator::main_loop::orchestrator_main_loop;
use orchestrator::main_loop::{ETH_ORACLE_LOOP_SPEED};
use std::path::Path;
use std::process::exit;

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
            error!("To set an already registered key use 'gbt keys set-orchestrator-key --phrase \"your phrase\"`");
            error!("To run from the command line, with no key storage use 'gbt orchestrator --cosmos-phrase \"your phrase\"' ");
            error!("If you have not already generated a key 'gbt keys register-orchestrator-address' will generate one for you");
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
    info!("Starting Gravity Validator companion binary Relayer + Oracle + Eth Signer");
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

    // TODO skos: this is unsafe to do, but we will pass contract address as args for now
    let contract_address = args.baseledger_contract_address.unwrap();
    // TODO skos: check if else branch here is needed...
    // check if we actually have the promised balance of tokens to pay fees
    // check_for_fee(&fee, public_cosmos_key, &contact).await;
    // check_for_eth(public_eth_key, &web3).await;

    // get the gravity contract address, if not provided
    // let contract_address = if let Some(c) = args.baseledger_contract_address {
    //     c
    // } else {
    //     let params = get_gravity_params(&mut grpc).await.unwrap();
    //     let c = params.bridge_ethereum_address.parse();
    //     match c {
    //         Ok(v) => {
    //             if v == *ZERO_ADDRESS {
    //                 error!("The Gravity address is not yet set as a chain parameter! You must specify --gravity-contract-address");
    //                 exit(1);
    //             }
    //             c.unwrap()
    //         }
    //         Err(_) => {
    //             error!("The Gravity address is not yet set as a chain parameter! You must specify --gravity-contract-address");
    //             exit(1);
    //         }
    //     }
    // };

    orchestrator_main_loop(
        cosmos_key,
        connections.web3.unwrap(),
        connections.contact.unwrap(),
        connections.grpc.unwrap(),
        contract_address,
        fee,
    )
    .await;
}
