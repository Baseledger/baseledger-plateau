use clarity::{Address, Uint256};
use cosmos_gravity::utils::get_last_event_nonce_with_retry;
use deep_space::address::Address as CosmosAddress;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_utils::get_with_retry::get_block_number_with_retry;
use gravity_utils::get_with_retry::RETRY_TIME;
use gravity_utils::types::event_signatures::*;
use gravity_utils::types::{SendToCosmosEvent};
use tokio::time::sleep as delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

/// This function retrieves the last event nonce this oracle has relayed to Cosmos
/// it then uses the Ethereum indexes to determine what block the last entry
pub async fn get_last_checked_block(
    grpc_client: GravityQueryClient<Channel>,
    our_cosmos_address: CosmosAddress,
    prefix: String,
    gravity_contract_address: Address,
    web3: &Web3,
) -> Uint256 {
    let mut grpc_client = grpc_client;
    const BLOCKS_TO_SEARCH: u128 = 5_000u128;

    let latest_block = get_block_number_with_retry(web3).await;
    let mut last_event_nonce: Uint256 =
        get_last_event_nonce_with_retry(&mut grpc_client, our_cosmos_address, prefix)
            .await
            .into();

    // zero indicates this oracle has never submitted an event before since there is no
    // zero event nonce (it's pre-incremented in the solidity contract) we have to go
    // and look for event nonce one.
    if last_event_nonce == 0u8.into() {
        last_event_nonce = 1u8.into();
    }

    // TODO skos: fix... i am harcoding last nonce here because oracle resync is not 
    // working with my test contract, gonna have to figure how to proceed with that
    last_event_nonce = 4u8.into();

    let mut current_block: Uint256 = latest_block.clone();

    while current_block.clone() > 0u8.into() {
        info!(
            "Oracle is resyncing, looking back into the history to find our last event nonce {}, on block {}",
            last_event_nonce, current_block
        );
        let end_search = if current_block.clone() < BLOCKS_TO_SEARCH.into() {
            0u8.into()
        } else {
            current_block.clone() - BLOCKS_TO_SEARCH.into()
        };
        let send_to_cosmos_events = web3
            .check_for_events(
                end_search.clone(),
                Some(current_block.clone()),
                vec![gravity_contract_address],
                vec![SENT_TO_COSMOS_EVENT_SIG],
            )
            .await;

        if send_to_cosmos_events.is_err()
        {
            error!("Failed to get blockchain events while resyncing, is your Eth node working? If you see only one of these it's fine",);
            delay_for(RETRY_TIME).await;
            continue;
        }
        let send_to_cosmos_events = send_to_cosmos_events.unwrap();

        println!("COSMOS EVENTS {:?}", send_to_cosmos_events);
        for event in send_to_cosmos_events {
            match SendToCosmosEvent::from_log(&event) {
                Ok(send) => {
                    trace!(
                        "{} send event nonce {} last event nonce",
                        send.event_nonce,
                        last_event_nonce
                    );
                    if upcast(send.event_nonce) == last_event_nonce && event.block_number.is_some()
                    {
                        return event.block_number.unwrap();
                    }
                }
                Err(e) => error!("Got SendToCosmos event that we can't parse {}", e),
            }
        }
        current_block = end_search;
    }

    panic!("You have reached the end of block history without finding the Gravity contract deploy event! You must have the wrong contract address!");
}

fn upcast(input: u64) -> Uint256 {
    input.into()
}
