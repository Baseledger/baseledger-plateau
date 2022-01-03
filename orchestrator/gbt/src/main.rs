#[macro_use]
extern crate log;
#[macro_use]
extern crate serde_derive;

use crate::args::{SubCommand};
use crate::{orchestrator::orchestrator};
use args::Opts;
use clap::Parser;
use env_logger::Env;


mod args;
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

    // control flow for the command structure
    match opts.subcmd {
        SubCommand::Orchestrator(orchestrator_opts) => {
            orchestrator(orchestrator_opts, address_prefix).await
        }
    }
}
