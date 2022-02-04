
use deep_space::address::Address;
use baseledger_proto::baseledger::{query_client::QueryClient as BaseledgerQueryClient, QueryLastEventNonceByAddressRequest, QueryGetOrchestratorValidatorAddressRequest};
use crate::error::OrchestratorError;
use std::process::exit;
use tonic::transport::Channel;
use deep_space::Address as CosmosAddress;

/// Gets the last event nonce that a given validator has attested to, this lets us
/// catch up with what the current event nonce should be if a oracle is restarted
pub async fn get_last_event_nonce_for_validator(
    client: &mut BaseledgerQueryClient<Channel>,
    address: Address,
    prefix: String,
) -> Result<u64, OrchestratorError> {
    let request = client
        .last_event_nonce_by_address(QueryLastEventNonceByAddressRequest {
            address: address.to_bech32(prefix).unwrap(),
        })
        .await?;
    Ok(request.into_inner().event_nonce)
}

/// This function checks if orchestrator and validator addresses were set
pub async fn check_validator_address(
    client: &mut BaseledgerQueryClient<Channel>,
    delegate_orchestrator_address: CosmosAddress,
    prefix: &str,
) {
    let orchestrator_response = client
        .orchestrator_validator_address(QueryGetOrchestratorValidatorAddressRequest {
            orchestrator_address: delegate_orchestrator_address.to_bech32(prefix).unwrap(),
        })
        .await;
    trace!("{:?}", orchestrator_response);
    match orchestrator_response {
        Ok(_) => {
            trace!("Validator found by orch address");
        }
        Err(e) => {
            error!("Your Orchestrator Cosmos key is incorrect, please double check your phrase. If you can't locate the correct phrase you will need to create a new validator {:?}", e);
            exit(1);
        }
    }
}