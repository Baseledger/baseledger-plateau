pub mod register_orchestrator_address;

use crate::{
    args::{SetOrchestratorKeyOpts},
    config::{config_exists, load_keys, save_keys},
};
use deep_space::PrivateKey;
use std::{path::Path, process::exit};

pub fn show_keys(home_dir: &Path, prefix: &str) {
    if !config_exists(home_dir) {
        error!("Please run `baseledger_bridge init` before running this command!");
        exit(1);
    }
    let keys = load_keys(home_dir);
    match keys.orchestrator_phrase {
        Some(v) => {
            let key = PrivateKey::from_phrase(&v, "")
                .expect("Failed to decode key in keyfile. Did you edit it manually?");
            let address = key.to_address(prefix).unwrap();
            info!("Your Orchestrator key, {}", address);
        }
        None => info!("You do not have an Orchestrator key set"),
    }
}

pub fn set_orchestrator_key(home_dir: &Path, opts: SetOrchestratorKeyOpts) {
    if !config_exists(home_dir) {
        error!("Please run `baseledger_bridge init` before running this command!");
        exit(1);
    }
    let res = PrivateKey::from_phrase(&opts.phrase, "");
    if let Err(e) = res {
        error!("Invalid Cosmos mnemonic phrase {} {:?}", opts.phrase, e);
        exit(1);
    }
    let mut keys = load_keys(home_dir);
    keys.orchestrator_phrase = Some(opts.phrase);
    save_keys(home_dir, keys);
    info!("Successfully updated Orchestrator Key")
}
