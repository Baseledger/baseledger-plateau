use clarity::{Address, Uint256};
use utils::get_with_retry::{get_last_event_nonce_with_retry, RETRY_TIME, get_block_number_with_retry};
use deep_space::address::Address as CosmosAddress;
use baseledger_proto::baseledger::query_client::QueryClient as BaseledgerQueryClient;
use utils::types::event_signatures::*;
use utils::types::{UbtSplitterEvent};
use tokio::time::sleep as delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

/// This function retrieves the last event nonce this oracle has relayed to Cosmos
/// it then uses the Ethereum indexes to determine what block the last entry
pub async fn get_last_checked_block(
    grpc_client: BaseledgerQueryClient<Channel>,
    our_cosmos_address: CosmosAddress,
    prefix: String,
    baseledger_contract_address: Address,
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
        let ubt_deposited_events = web3
            .check_for_events(
                end_search.clone(),
                Some(current_block.clone()),
                vec![baseledger_contract_address],
                vec![UBT_DEPOSITED_EVENT_SIG],
            )
            .await;
        
        let power_changes_events = web3
            .check_for_events(
                end_search.clone(),
                Some(current_block.clone()),
                vec![baseledger_contract_address],
                vec![VALIDATOR_POWER_CHANGE_EVENT_SIG],
            )
            .await;

        if ubt_deposited_events.is_err() || power_changes_events.is_err()
        {
            error!("Failed to get blockchain events while resyncing, is your Eth node working? If you see only one of these it's fine",);
            delay_for(RETRY_TIME).await;
            continue;
        }
        let ubt_deposited_events = ubt_deposited_events.unwrap();
        let power_changes_events = power_changes_events.unwrap();

        for event in power_changes_events {
            match UbtSplitterEvent::from_log(&event) {
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

        for event in ubt_deposited_events {
            match UbtSplitterEvent::from_log(&event) {
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

    latest_block.clone()
}

fn upcast(input: u64) -> Uint256 {
    input.into()
}
