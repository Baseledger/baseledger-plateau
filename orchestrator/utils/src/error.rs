//! for things that don't belong in the cosmos or ethereum libraries but also don't belong
//! in a function specific library

use clarity::Error as ClarityError;
use deep_space::error::AddressError as CosmosAddressError;
use deep_space::error::CosmosGrpcError;
use num_bigint::ParseBigIntError;
use std::fmt::{self, Debug};
use tokio::time::error::Elapsed;
use tonic::Status;
use web30::jsonrpc::error::Web3Error;

#[derive(Debug)]
#[allow(clippy::large_enum_variant)]
pub enum OrchestratorError {
    InvalidBigInt(ParseBigIntError),
    CosmosGrpcError(CosmosGrpcError),
    CosmosAddressError(CosmosAddressError),
    EthereumRestError(Web3Error),
    InvalidBridgeStateError(String),
    FailedToUpdateValset,
    EthereumContractError(String),
    InvalidOptionsError(String),
    ClarityError(ClarityError),
    TimeoutError,
    InvalidEventLogError(String),
    OrchestratorGrpcError(Status),
    InsufficientVotingPowerToPass(String),
    ParseBigIntError(ParseBigIntError),
}

impl fmt::Display for OrchestratorError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            OrchestratorError::OrchestratorGrpcError(val) => write!(f, "Orchestrator gRPC error {}", val),
            OrchestratorError::CosmosGrpcError(val) => write!(f, "Cosmos gRPC error {}", val),
            OrchestratorError::InvalidBigInt(val) => {
                write!(f, "Got invalid BigInt from cosmos! {}", val)
            }
            OrchestratorError::CosmosAddressError(val) => write!(f, "Cosmos Address error {}", val),
            OrchestratorError::EthereumRestError(val) => write!(f, "Ethereum REST error {}", val),
            OrchestratorError::InvalidOptionsError(val) => {
                write!(f, "Invalid TX options for this call {}", val)
            }
            OrchestratorError::InvalidBridgeStateError(val) => {
                write!(f, "Invalid bridge state! {}", val)
            }
            OrchestratorError::FailedToUpdateValset => write!(f, "ValidatorSetUpdate Failed!"),
            OrchestratorError::TimeoutError => write!(f, "Operation timed out!"),
            OrchestratorError::ClarityError(val) => write!(f, "Clarity Error {}", val),
            OrchestratorError::InvalidEventLogError(val) => write!(f, "InvalidEvent: {}", val),
            OrchestratorError::EthereumContractError(val) => {
                write!(f, "Contract operation failed: {}", val)
            }
            OrchestratorError::InsufficientVotingPowerToPass(val) => {
                write!(f, "{}", val)
            }
            OrchestratorError::ParseBigIntError(val) => write!(f, "Failed to parse big integer {}", val),
        }
    }
}

impl std::error::Error for OrchestratorError {}

impl From<CosmosGrpcError> for OrchestratorError {
    fn from(error: CosmosGrpcError) -> Self {
        OrchestratorError::CosmosGrpcError(error)
    }
}

impl From<Elapsed> for OrchestratorError {
    fn from(_error: Elapsed) -> Self {
        OrchestratorError::TimeoutError
    }
}

impl From<ClarityError> for OrchestratorError {
    fn from(error: ClarityError) -> Self {
        OrchestratorError::ClarityError(error)
    }
}

impl From<Web3Error> for OrchestratorError {
    fn from(error: Web3Error) -> Self {
        OrchestratorError::EthereumRestError(error)
    }
}
impl From<Status> for OrchestratorError {
    fn from(error: Status) -> Self {
        OrchestratorError::OrchestratorGrpcError(error)
    }
}
impl From<CosmosAddressError> for OrchestratorError {
    fn from(error: CosmosAddressError) -> Self {
        OrchestratorError::CosmosAddressError(error)
    }
}
impl From<ParseBigIntError> for OrchestratorError {
    fn from(error: ParseBigIntError) -> Self {
        OrchestratorError::InvalidBigInt(error)
    }
}
