use crate::{ubt_price_fetcher::fetch_ubt_price};
use clarity::{Address as EthAddress, Uint256};
use web30::client::Web3;
use web30::jsonrpc::error::Web3Error;
use utils::get_with_retry::get_net_version_with_retry;
use utils::{
    error::OrchestratorError,
    types::{
        UbtDeposited,
        PayeeUpdated,
    },
};
use utils::get_with_retry::get_block_number_with_retry;
use utils::types::event_signatures::*;
use deep_space::Contact;
use deep_space::{private_key::PrivateKey as CosmosPrivateKey};
use tonic::transport::Channel;
use baseledger_proto::baseledger::query_client::QueryClient as BaseledgerQueryClient;
use utils::cosmos::{query::get_last_event_nonce_for_validator, send::send_ethereum_claims};

pub struct CheckedNonces {
    pub block_number: Uint256,
    pub event_nonce: Uint256,
}

#[allow(clippy::too_many_arguments)]
pub async fn check_for_events(
    web3: &Web3,
    contact: &Contact,
    grpc_client: &mut BaseledgerQueryClient<Channel>,
    baseledger_contract_address: EthAddress,
    our_private_key: CosmosPrivateKey,
    starting_block: Uint256,
    block_delay: Uint256,
) -> Result<CheckedNonces, OrchestratorError> {
    let our_cosmos_address = our_private_key.to_address(&contact.get_prefix()).unwrap();
    let latest_block = get_block_number_with_retry(web3).await;
    let latest_block = latest_block - block_delay;

    let deposits = web3
        .check_for_events(
            starting_block.clone(),
            Some(latest_block.clone()),
            vec![baseledger_contract_address],
            vec![SENT_TO_COSMOS_EVENT_SIG],
        )
        .await;
    trace!("Deposits {:?}", deposits);

    let power_changes = web3
        .check_for_events(
            starting_block.clone(),
            Some(latest_block.clone()),
            vec![baseledger_contract_address],
            vec![VALIDATOR_POWER_CHANGE_EVENT_SIG],
        )
        .await;
    trace!("Deposits {:?}", power_changes);

    if let (Ok(deposits), Ok(power_changes)) = (deposits, power_changes)
    {
        let deposits = UbtDeposited::from_logs(&deposits)?;
        trace!("parsed deposits {:?}", deposits);

        let power_changes = PayeeUpdated::from_logs(&power_changes)?;
        trace!("parsed power_changes {:?}", power_changes);

        // note that starting block overlaps with our last checked block, because we have to deal with
        // the possibility that the relayer was killed after relaying only one of multiple events in a single
        // block, so we also need this routine so make sure we don't send in the first event in this hypothetical
        // multi event block again. In theory we only send all events for every block and that will pass of fail
        // atomicly but lets not take that risk.
        let last_event_nonce = get_last_event_nonce_for_validator(
            grpc_client,
            our_cosmos_address,
            contact.get_prefix(),
        )
        .await?;

        let deposits = UbtDeposited::filter_by_event_nonce(last_event_nonce, &deposits);

        let power_changes = PayeeUpdated::filter_by_event_nonce(last_event_nonce, &power_changes);

        if !deposits.is_empty() {
            info!(
                "Oracle observed deposit with sender {}, destination {:?}, amount {}, and event nonce {}",
                deposits[0].sender, deposits[0].validated_destination, deposits[0].amount, deposits[0].event_nonce
            )
        }

        let ubt_price = fetch_ubt_price().await.unwrap();
        
        if !power_changes.is_empty() {
            info!(
                "Oracle observed power change with sender {}, destination {:?}, amount {}, and event nonce {}",
                power_changes[0].revenue_address, power_changes[0].validated_destination, power_changes[0].shares, power_changes[0].event_nonce
            )
        }

        let mut new_event_nonce: Uint256 = last_event_nonce.into();
        if !deposits.is_empty() || !power_changes.is_empty()
        {
            let res = send_ethereum_claims(
                contact,
                our_private_key,
                deposits,
                power_changes,
                ubt_price,
            )
            .await?;
            new_event_nonce = get_last_event_nonce_for_validator(
                grpc_client,
                our_cosmos_address,
                contact.get_prefix(),
            )
            .await?
            .into();

            info!("Current event nonce is {}", new_event_nonce);

            // since we can't actually trust that the above txresponse is correct we have to check here
            // we may be able to trust the tx response post grpc
            if new_event_nonce == last_event_nonce.into() {
                return Err(OrchestratorError::InvalidBridgeStateError(
                    format!("Claims did not process, trying to update but still on {}, trying again in a moment, check txhash {} for errors", last_event_nonce, res.txhash),
                ));
            } else {
                info!("Claims processed, new nonce {}", new_event_nonce);
            }
        }
        Ok(CheckedNonces {
            block_number: latest_block,
            event_nonce: new_event_nonce,
        })
    } else {
        error!("Failed to get events");
        Err(OrchestratorError::EthereumRestError(Web3Error::BadResponse(
            "Failed to get logs!".to_string(),
        )))
    }
}

/// The number of blocks behind the 'latest block' on Ethereum our event checking should be.
/// Ethereum does not have finality and as such is subject to chain reorgs and temporary forks
/// if we check for events up to the very latest block we may process an event which did not
/// 'actually occur' in the longest POW chain.
///
/// Obviously we must chose some delay in order to prevent incorrect events from being claimed
///
/// For EVM chains with finality the correct value for this is zero. As there's no need
/// to concern ourselves with re-orgs or forking. This function checks the netID of the
/// provided Ethereum RPC and adjusts the block delay accordingly
///
/// The value used here for Ethereum is a balance between being reasonably fast and reasonably secure
/// As you can see on https://etherscan.io/blocks_forked uncles (one block deep reorgs)
/// occur once every few minutes. Two deep once or twice a day.
/// https://etherscan.io/chart/uncles
/// Let's make a conservative assumption of 1% chance of an uncle being a two block deep reorg
/// (actual is closer to 0.3%) and assume that continues as we increase the depth.
/// Given an uncle every 2.8 minutes, a 6 deep reorg would be 2.8 minutes * (100^4) or one
/// 6 deep reorg every 53,272 years.
///
/// Of course the above assume that no mining attacks occur. Once we bring that potential into
/// the equation the question becomes 'how much money'. There is no depth safe from infinite
/// spending. Taking some source values from https://blog.ethereum.org/2016/05/09/on-settlement-finality/
/// we will use 13 blocks providing a 1/1_000_000 chance of an attacker with 25% of network hash
/// power succeeding
///
pub async fn get_block_delay(web3: &Web3) -> Uint256 {
    let net_version = get_net_version_with_retry(web3).await;

    match net_version {
        // Mainline Ethereum, Ethereum classic, or the Ropsten, Kotti, Mordor testnets
        // all POW Chains
        1 | 3 | 6 | 7 => 13u8.into(),
        // Dev and Hardhat respectively
        // all single signer chains with no chance of any reorgs
        2018 | 31337 => 0u8.into(),
        // Rinkeby and Goerli use Clique (POA) Consensus, finality takes
        // up to num validators blocks. Number is higher than Ethereum based
        // on experience with operational issues
        4 | 5 => 10u8.into(),
        // assume the safe option (POW) where we don't know
        _ => 13u8.into(),
    }
}
