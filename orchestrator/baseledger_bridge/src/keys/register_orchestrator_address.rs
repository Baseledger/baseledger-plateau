use std::path::PathBuf;
use std::process::exit;
use std::time::Duration;

use crate::args::RegisterOrchestratorAddressOpts;
use crate::config::config_exists;
use crate::config::load_keys;
use crate::config::save_keys;
use crate::config::KeyStorage;
use utils::cosmos::send::set_orchestrator_validator_addresses;
use deep_space::{mnemonic::Mnemonic, private_key::PrivateKey as CosmosPrivateKey};
use utils::connection_prep::check_for_tokens;
use utils::connection_prep::{create_rpc_connections, wait_for_cosmos_node_ready};

pub const TIMEOUT: Duration = Duration::from_secs(60);

pub async fn register_orchestrator_address(
    args: RegisterOrchestratorAddressOpts,
    prefix: String,
    home_dir: PathBuf,
) {
    let cosmos_grpc = args.cosmos_grpc;
    let validator_key = args.validator_phrase;
    let cosmos_phrase = args.cosmos_phrase;
    let mut generated_cosmos = None;

    if !args.no_save && !config_exists(&home_dir) {
        error!("Please run `baseledger_bridge init` before running this command!");
        exit(1);
    }

    let connections = create_rpc_connections(prefix, Some(cosmos_grpc), None, TIMEOUT).await;
    let contact = connections.contact.unwrap();
    wait_for_cosmos_node_ready(&contact).await;

    let validator_addr = validator_key.to_address(&contact.get_prefix()).unwrap();
    check_for_tokens(validator_addr, &contact).await;

    // Set the cosmos key to either the cli value, the value in the config, or a generated
    // value if the config has not been setup
    let cosmos_key = if let Some(cosmos_phrase) = cosmos_phrase.clone() {
        CosmosPrivateKey::from_phrase(&cosmos_phrase, "").expect("Failed to parse cosmos key")
    } else {
        let mut key = None;
        if config_exists(&home_dir) {
            let keys = load_keys(&home_dir);
            if let Some(phrase) = keys.orchestrator_phrase {
                key = Some(CosmosPrivateKey::from_phrase(phrase.as_str(), "").unwrap());
            }
        }
        if key.is_none() {
            let new_phrase = Mnemonic::generate(24).unwrap();
            key = Some(CosmosPrivateKey::from_phrase(new_phrase.as_str(), "").unwrap());
            generated_cosmos = Some(new_phrase);
        }
        key.unwrap()
    };

    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    let res = set_orchestrator_validator_addresses(
        &contact,
        cosmos_address,
        validator_key,
    )
    .await
    .expect("Failed to update delegate address");
    let res = contact.wait_for_tx(res, TIMEOUT).await;

    if let Err(e) = res {
        error!("Failed trying to register delegate addresses error {:?}, correct the error and try again", e);
        exit(1);
    }

    if let Some(phrase) = generated_cosmos.clone() {
        info!(
            "No Cosmos key provided, your generated key is\n {} -> {}",
            phrase.as_str(),
            cosmos_key.to_address(&contact.get_prefix()).unwrap()
        );
    }

    // TODO: should we hardcode baseledger_contract_address since that won't change?
    if !args.no_save {
        info!("Keys saved! You can now run `baseledger_bridge orchestrator --ethereum_rpc <eth_rpc_value> --baseledger-contract-address <baseledger_contract_address>`");
        let phrase = match (generated_cosmos, cosmos_phrase) {
            (Some(v), None) => v.to_string(),
            (None, Some(s)) => s,
            (_, _) => {
                // in this case the user has set keys in the config
                // and then registered them so lets just load the config
                // value again
                let keys = load_keys(&home_dir);
                keys.orchestrator_phrase.unwrap()
            }
        };
        let new_keys = KeyStorage {
            orchestrator_phrase: Some(phrase),
        };
        save_keys(&home_dir, new_keys);
    }
}
