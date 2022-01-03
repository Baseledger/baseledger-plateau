
use deep_space::address::Address;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_proto::gravity::Attestation;
use gravity_proto::gravity::Params;
use gravity_proto::gravity::QueryAttestationsRequest;
use gravity_proto::gravity::QueryLastEventNonceByAddrRequest;
use gravity_proto::gravity::QueryParamsRequest;
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
    // let request = client
    //     .last_event_nonce_by_addr(QueryLastEventNonceByAddrRequest {
    //         address: address.to_bech32(prefix).unwrap(),
    //     })
    //     .await?;
    // Ok(request.into_inner().event_nonce)

    Ok(1)
}

pub async fn get_attestations(
    client: &mut GravityQueryClient<Channel>,
    limit: Option<u64>,
) -> Result<Vec<Attestation>, GravityError> {
    let request = client
        .get_attestations(QueryAttestationsRequest {
            limit: limit.or(Some(1000u64)).unwrap(),
        })
        .await?;
    let attestations = request.into_inner().attestations;
    Ok(attestations)
}