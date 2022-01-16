
use deep_space::address::Address;
use gravity_proto::baseledger::query_client::QueryClient as GravityQueryClient;
use gravity_proto::baseledger::Attestation;
use gravity_proto::baseledger::Params;
use gravity_proto::baseledger::QueryAttestationsRequest;
use gravity_proto::baseledger::QueryLastEventNonceByAddressRequest;
use gravity_proto::baseledger::QueryParamsRequest;
use gravity_utils::error::GravityError;
use tonic::transport::Channel;

/// Gets the Gravity module parameters from the Gravity module
pub async fn get_gravity_params(
    client: &mut GravityQueryClient<Channel>,
) -> Result<Params, GravityError> {
    let request = client.params(QueryParamsRequest {}).await?.into_inner();
    Ok(request.params.unwrap())
}


/// Gets the last event nonce that a given validator has attested to, this lets us
/// catch up with what the current event nonce should be if a oracle is restarted
pub async fn get_last_event_nonce_for_validator(
    client: &mut GravityQueryClient<Channel>,
    address: Address,
    prefix: String,
) -> Result<u64, GravityError> {
    // TODO skos: commented this out to make it work
    let request = client
        .last_event_nonce_by_address(QueryLastEventNonceByAddressRequest {
            address: address.to_bech32(prefix).unwrap(),
        })
        .await?;
    Ok(request.into_inner().event_nonce)
}

pub async fn get_attestations(
    client: &mut GravityQueryClient<Channel>,
    limit: Option<u64>,
) -> Result<Vec<Attestation>, GravityError> {
    let request = client
        .attestations(QueryAttestationsRequest {
            limit: limit.or(Some(1000u64)).unwrap(),
        })
        .await?;
    let attestations = request.into_inner().attestations;
    Ok(attestations)
}