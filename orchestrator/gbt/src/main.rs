#[macro_use]
extern crate log;
#[macro_use]
extern crate serde_derive;

use crate::args::{KeysSubcommand, SubCommand};
use crate::config::init_config;
use crate::keys::show_keys;
use crate::{orchestrator::orchestrator};
use args::Opts;
use clap::Parser;
use config::{get_home_dir};
use env_logger::Env;
use keys::register_orchestrator_address::register_orchestrator_address;
use keys::set_eth_key;
use keys::set_orchestrator_key;

mod args;
mod config;
mod keys;
mod orchestrator;
mod utils;

#[actix_rt::main]
async fn main() {
    env_logger::Builder::from_env(Env::default().default_filter_or("info")).init();
    // On Linux static builds we need to probe ssl certs path to be able to
    // do TLS stuff.
    openssl_probe::init_ssl_cert_env_vars();
    // parse the arguments
    let opts: Opts = Opts::parse();

    // handle global config here
    let address_prefix = opts.address_prefix;
    let home_dir = get_home_dir(opts.home);

    // control flow for the command structure
    match opts.subcmd {
        SubCommand::Keys(key_opts) => match key_opts.subcmd {
            KeysSubcommand::RegisterOrchestratorAddress(set_orchestrator_address_opts) => {
                register_orchestrator_address(
                    set_orchestrator_address_opts,
                    address_prefix,
                    home_dir,
                )
                .await
            }
            KeysSubcommand::Show => show_keys(&home_dir, &address_prefix),
            KeysSubcommand::SetEthereumKey(set_eth_key_opts) => {
                set_eth_key(&home_dir, set_eth_key_opts)
            }
            KeysSubcommand::SetOrchestratorKey(set_orch_key_opts) => {
                set_orchestrator_key(&home_dir, set_orch_key_opts)
            }
        },
        SubCommand::Orchestrator(orchestrator_opts) => {
            let test = orchestrator(orchestrator_opts, address_prefix, &home_dir).await;
        }
        SubCommand::Init(init_opts) => init_config(init_opts, home_dir),
    }
}
