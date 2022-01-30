use clarity::Uint256;
use deep_space::{address::Address as CosmosAddress, Coin, Contact};
use std::time::Duration;
use tokio::time::sleep as delay_for;
use web30::client::Web3;
use crate::cosmos::query::get_last_event_nonce_for_validator;
use baseledger_proto::baseledger::query_client::QueryClient as BaseledgerQueryClient;
use tonic::transport::Channel;
use std::u64::MAX as U64MAX;

pub const RETRY_TIME: Duration = Duration::from_secs(5);

/// gets the current Ethereum block number, no matter how long it takes
pub async fn get_block_number_with_retry(web3: &Web3) -> Uint256 {
    let mut res = web3.eth_block_number().await;
    while res.is_err() {
        error!("Failed to get latest block! Is your Eth node working?");
        delay_for(RETRY_TIME).await;
        res = web3.eth_block_number().await;
    }
    res.unwrap()
}

/// gets Cosmos balances, no matter how long it takes
pub async fn get_balances_with_retry(address: CosmosAddress, contact: &Contact) -> Vec<Coin> {
    let mut res = contact.get_balances(address).await;
    while res.is_err() {
        error!("Failed to get Cosmos balances! Is your Cosmos node working?");
        delay_for(RETRY_TIME).await;
        res = contact.get_balances(address).await;
    }
    res.unwrap()
}

/// gets the net version, no matter how long it takes
pub async fn get_net_version_with_retry(web3: &Web3) -> u64 {
    let mut res = web3.net_version().await;
    while res.is_err() {
        error!("Failed to get net version! Is your Eth node working?");
        delay_for(RETRY_TIME).await;
        res = web3.net_version().await;
    }
    res.unwrap()
}

/// gets the Cosmos last event nonce, no matter how long it takes.
pub async fn get_last_event_nonce_with_retry(
    client: &mut BaseledgerQueryClient<Channel>,
    our_cosmos_address: CosmosAddress,
    prefix: String,
) -> u64 {
    let mut res =
        get_last_event_nonce_for_validator(client, our_cosmos_address, prefix.clone()).await;
    while res.is_err() {
        error!(
            "Failed to get last event nonce, is the Cosmos GRPC working? {:?}",
            res
        );
        delay_for(RETRY_TIME).await;
        res = get_last_event_nonce_for_validator(client, our_cosmos_address, prefix.clone()).await;
    }
    res.unwrap()
}

pub fn downcast_uint256(input: Uint256) -> Option<u64> {
    if input >= U64MAX.into() {
        None
    } else {
        let mut val = input.to_bytes_be();
        // pad to 8 bytes
        while val.len() < 8 {
            val.insert(0, 0);
        }
        let mut lower_bytes: [u8; 8] = [0; 8];
        // get the 'lowest' 8 bytes from a 256 bit integer
        lower_bytes.copy_from_slice(&val[0..val.len()]);
        Some(u64::from_be_bytes(lower_bytes))
    }
}