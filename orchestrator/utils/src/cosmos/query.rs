
use deep_space::address::Address;
use baseledger_proto::baseledger::query_client::QueryClient as BaseledgerQueryClient;
use baseledger_proto::baseledger::QueryLastEventNonceByAddressRequest;
use crate::error::OrchestratorError;
use tonic::transport::Channel;

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