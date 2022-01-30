use clap::Parser;
use clarity::Address as EthAddress;
use deep_space::PrivateKey as CosmosPrivateKey;
use deep_space::{Coin};
use std::path::PathBuf;

#[derive(Parser)]
#[clap(version = env!("CARGO_PKG_VERSION"))]
pub struct Opts {
    /// Increase the logging verbosity
    #[clap(short, long)]
    pub verbose: bool,
    /// Decrease the logging verbosity
    #[clap(short, long)]
    pub quiet: bool,
    #[clap(short, long, parse(from_str))]
    pub home: Option<PathBuf>,
    /// Set the address prefix for the Cosmos chain
    /// default is 'cosmos'
    #[clap(short, long, default_value = "baseledger")]
    pub address_prefix: String,
    #[clap(subcommand)]
    pub subcmd: SubCommand,
}

#[derive(Parser)]
pub enum SubCommand {
    Orchestrator(OrchestratorOpts),
    Keys(KeyOpts),
    Init(InitOpts),
}

#[derive(Parser)]
pub struct OrchestratorOpts {
    /// Cosmos mnemonic phrase that orchestrator will use to sign txs
    #[clap(short, long, parse(try_from_str))]
    pub cosmos_phrase: Option<CosmosPrivateKey>,
    /// (Optional) The Cosmos gRPC server that will be used
    #[clap(long, default_value = "http://localhost:9090")]
    pub cosmos_grpc: String,
    /// (Optional) The Ethereum RPC server that will be used
    #[clap(long, default_value = "http://localhost:8545")]
    pub ethereum_rpc: String,
    /// The address for the Baseledger contract on Ethereum
    #[clap(short, long, parse(try_from_str))]
    pub baseledger_contract_address: Option<EthAddress>,
}

/// Manage keys
#[derive(Parser)]
pub struct KeyOpts {
    #[clap(subcommand)]
    pub subcmd: KeysSubcommand,
}

#[derive(Parser)]
pub enum KeysSubcommand {
    RegisterOrchestratorAddress(RegisterOrchestratorAddressOpts),
    SetOrchestratorKey(SetOrchestratorKeyOpts),
    Show,
}

/// Register cosmos key for orchestrator
#[derive(Parser)]
pub struct RegisterOrchestratorAddressOpts {
    /// The Cosmos private key of the validator
    #[clap(short, long, parse(try_from_str))]
    pub validator_phrase: CosmosPrivateKey,
    /// (Optional) The phrase for the Cosmos key to register, will be generated if not provided.
    #[clap(short, long, parse(try_from_str))]
    pub cosmos_phrase: Option<String>,
    /// (Optional) The Cosmos gRPC server that will be used to submit the transaction
    #[clap(long, default_value = "http://localhost:9090")]
    pub cosmos_grpc: String,
    /// The Cosmos Denom and amount to pay Cosmos chain fees
    #[clap(short, long, parse(try_from_str))]
    pub fees: Coin,
    /// Do not save keys to disk for later use with `orchestrator start`
    #[clap(long)]
    pub no_save: bool,
}

/// Add a Cosmos private key to use as the Orchestrator address
#[derive(Parser)]
pub struct SetOrchestratorKeyOpts {
    #[clap(short, long)]
    pub phrase: String,
}

/// Initialize configuration
#[derive(Parser)]
pub struct InitOpts {}
