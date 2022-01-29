//! Handles configuration structs + saving and loading for bridge

use crate::args::InitOpts;
use std::{
    fs::{self, create_dir},
    path::{Path, PathBuf},
    process::exit,
};

/// The name of the keys file, this file is not expected
/// to be hand edited.
pub const KEYS_NAME: &str = "keys.json";
/// The folder name for the config
pub const CONFIG_FOLDER: &str = ".baseledger_bridge";

/// The keys storage struct, including encrypted and un-encrypted local keys
/// un-encrypted keys provide for orchestrator start and relayer start functions
#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Default)]
pub struct KeyStorage {
    pub orchestrator_phrase: Option<String>,
}

/// Checks if the user has setup their config environment
pub fn config_exists(home_dir: &Path) -> bool {
    let keys_file = home_dir.join(CONFIG_FOLDER).with_file_name(KEYS_NAME);
    home_dir.exists() && keys_file.exists()
}

/// Creates the config directory and default config file if it does
/// not already exist
pub fn init_config(_init_ops: InitOpts, home_dir: PathBuf) {
    if home_dir.exists() {
        warn!(
            "Config folder {} already exists!",
            home_dir.to_str().unwrap()
        );
        warn!("You can delete this folder and run init again, you will lose any keys or other config data!");
    } else {
        create_dir(home_dir.clone()).expect("Failed to create config directory!");
        fs::write(
            home_dir.join(KEYS_NAME),
            toml::to_string(&KeyStorage::default()).unwrap(),
        )
        .expect("Unable to write config file");
    }
}

pub fn get_home_dir(home_arg: Option<PathBuf>) -> PathBuf {
    match (dirs::home_dir(), home_arg) {
        (_, Some(user_home)) => PathBuf::from(&user_home),
        (Some(default_home_dir), None) => default_home_dir.join(CONFIG_FOLDER),
        (None, None) => {
            error!("Failed to automatically determine your home directory, please provide a path to the --home argument!");
            exit(1);
        }
    }
}

/// Load the keys file, this operates at runtime
pub fn load_keys(home_dir: &Path) -> KeyStorage {
    let keys_file = home_dir.join(CONFIG_FOLDER).with_file_name(KEYS_NAME);
    if !keys_file.exists() {
        error!(
            "Keys file at {} not detected, use `baseledger_bridge init` to generate a config.",
            keys_file.to_str().unwrap()
        );
        exit(1);
    }

    let keys = fs::read_to_string(keys_file).unwrap();
    match toml::from_str(&keys) {
        Ok(v) => v,
        Err(e) => {
            error!("Invalid keys! {:?}", e);
            exit(1);
        }
    }
}

/// Saves the keys file, overwriting the existing one
pub fn save_keys(home_dir: &Path, updated_keys: KeyStorage) {
    let config_file = home_dir.join(CONFIG_FOLDER).with_file_name(KEYS_NAME);
    if !config_file.exists() {
        info!(
            "Config file at {} not detected, using defaults, use `baseledger_bridge init` to generate a config.",
            config_file.to_str().unwrap()
        );
    }

    fs::write(
        home_dir.join(KEYS_NAME),
        toml::to_string(&updated_keys).unwrap(),
    )
    .expect("Unable to write config file");
}
