//! This module provides useful tools for handling the Contact and Web30 connections for the relayer, orchestrator and various other utilities.
//! It's a common problem to have conflicts between ipv4 and ipv6 localhost and this module is first and foremost supposed to resolve that problem
//! by trying more than one thing to handle potentially misconfigured inputs.

use deep_space::error::CosmosGrpcError;
use deep_space::Address as CosmosAddress;
use deep_space::Contact;
use deep_space::{client::ChainStatus, Coin};
use baseledger_proto::baseledger::query_client::QueryClient as BaseledgerQueryClient;
use baseledger_proto::baseledger::QueryValidatorAddressByOrchestratorAddressRequest;
use std::process::exit;
use std::time::Duration;
use tokio::time::sleep as delay_for;
use tonic::transport::Channel;
use url::Url;
use web30::client::Web3;

use crate::get_with_retry::get_balances_with_retry;

pub struct Connections {
    pub web3: Option<Web3>,
    pub grpc: Option<BaseledgerQueryClient<Channel>>,
    pub contact: Option<Contact>,
}

/// Returns the three major RPC connections required for Gravity
/// operation in a error resilient manner. TODO find some way to generalize
/// this so that it's less ugly
pub async fn create_rpc_connections(
    address_prefix: String,
    grpc_url: Option<String>,
    eth_rpc_url: Option<String>,
    timeout: Duration,
) -> Connections {
    let mut web3 = None;
    let mut grpc = None;
    let mut contact = None;
    if let Some(grpc_url) = grpc_url {
        let url = Url::parse(&grpc_url)
            .unwrap_or_else(|_| panic!("Invalid Cosmos gRPC url {}", grpc_url));
        check_scheme(&url, &grpc_url);
        let cosmos_grpc_url = grpc_url.trim_end_matches('/').to_string();
        // try the base url first.
        let try_base = BaseledgerQueryClient::connect(cosmos_grpc_url.clone()).await;
        match try_base {
            // it worked, lets go!
            Ok(val) => {
                grpc = Some(val);
                contact = Some(Contact::new(&cosmos_grpc_url, timeout, &address_prefix).unwrap());
            }
            // did not work, now we check if it's localhost
            Err(e) => {
                warn!(
                    "Failed to access Cosmos gRPC with {:?} trying fallback options",
                    e
                );
                if grpc_url.to_lowercase().contains("localhost") {
                    let port = url.port().unwrap_or(80);
                    // this should be http or https
                    let prefix = url.scheme();
                    let ipv6_url = format!("{}://::1:{}", prefix, port);
                    let ipv4_url = format!("{}://127.0.0.1:{}", prefix, port);
                    let ipv6 = BaseledgerQueryClient::connect(ipv6_url.clone()).await;
                    let ipv4 = BaseledgerQueryClient::connect(ipv4_url.clone()).await;
                    warn!("Trying fallback urls {} {}", ipv6_url, ipv4_url);
                    match (ipv4, ipv6) {
                        (Ok(v), Err(_)) => {
                            info!("Url fallback succeeded, your cosmos gRPC url {} has been corrected to {}", grpc_url, ipv4_url);
                            contact = Some(Contact::new(&ipv4_url, timeout, &address_prefix).unwrap());
                            grpc = Some(v)
                        },
                        (Err(_), Ok(v)) => {
                            info!("Url fallback succeeded, your cosmos gRPC url {} has been corrected to {}", grpc_url, ipv6_url);
                            contact = Some(Contact::new(&ipv6_url, timeout, &address_prefix).unwrap());
                            grpc = Some(v)
                        },
                        (Ok(_), Ok(_)) => panic!("This should never happen? Why didn't things work the first time?"),
                        (Err(_), Err(_)) => panic!("Could not connect to Cosmos gRPC, are you sure it's running and on the specified port? {}", grpc_url)
                    }
                } else if url.port().is_none() || url.scheme() == "http" {
                    let body = url.host_str().unwrap_or_else(|| {
                        panic!("Cosmos gRPC url contains no host? {}", grpc_url)
                    });
                    // transparently upgrade to https if available, we can't transparently downgrade for obvious security reasons
                    let https_on_80_url = format!("https://{}:80", body);
                    let https_on_443_url = format!("https://{}:443", body);
                    let https_on_80 = BaseledgerQueryClient::connect(https_on_80_url.clone()).await;
                    let https_on_443 = BaseledgerQueryClient::connect(https_on_443_url.clone()).await;
                    warn!(
                        "Trying fallback urls {} {}",
                        https_on_443_url, https_on_80_url
                    );
                    match (https_on_80, https_on_443) {
                        (Ok(v), Err(_)) => {
                            info!("Https upgrade succeeded, your cosmos gRPC url {} has been corrected to {}", grpc_url, https_on_80_url);
                            contact = Some(Contact::new(&https_on_80_url, timeout, &address_prefix).unwrap());
                            grpc = Some(v)
                        },
                        (Err(_), Ok(v)) => {
                            info!("Https upgrade succeeded, your cosmos gRPC url {} has been corrected to {}", grpc_url, https_on_443_url);
                            contact = Some(Contact::new(&https_on_443_url, timeout, &address_prefix).unwrap());
                            grpc = Some(v)
                        },
                        (Ok(_), Ok(_)) => panic!("This should never happen? Why didn't things work the first time?"),
                        (Err(_), Err(_)) => panic!("Could not connect to Cosmos gRPC, are you sure it's running and on the specified port? {}", grpc_url)
                    }
                } else {
                    panic!("Could not connect to Cosmos gRPC! please check your grpc url {} for errors {:?}", grpc_url, e)
                }
            }
        }
    }
    if let Some(eth_rpc_url) = eth_rpc_url {
        let url = Url::parse(&eth_rpc_url)
            .unwrap_or_else(|_| panic!("Invalid Ethereum RPC url {}", eth_rpc_url));
        check_scheme(&url, &eth_rpc_url);
        let eth_url = eth_rpc_url.trim_end_matches('/');
        let base_web30 = Web3::new(eth_url, timeout);
        let try_base = base_web30.eth_block_number().await;
        match try_base {
            // it worked, lets go!
            Ok(_) => web3 = Some(base_web30),
            // did not work, now we check if it's localhost
            Err(e) => {
                warn!(
                    "Failed to access Ethereum RPC with {:?} trying fallback options",
                    e
                );
                if eth_url.to_lowercase().contains("localhost") {
                    let port = url.port().unwrap_or(80);
                    // this should be http or https
                    let prefix = url.scheme();
                    let ipv6_url = format!("{}://::1:{}", prefix, port);
                    let ipv4_url = format!("{}://127.0.0.1:{}", prefix, port);
                    let ipv6_web3 = Web3::new(&ipv6_url, timeout);
                    let ipv4_web3 = Web3::new(&ipv4_url, timeout);
                    let ipv6_test = ipv6_web3.eth_block_number().await;
                    let ipv4_test = ipv4_web3.eth_block_number().await;
                    warn!("Trying fallback urls {} {}", ipv6_url, ipv4_url);
                    match (ipv4_test, ipv6_test) {
                        (Ok(_), Err(_)) => {
                            info!("Url fallback succeeded, your Ethereum rpc url {} has been corrected to {}", eth_rpc_url, ipv4_url);
                            web3 = Some(ipv4_web3)
                        }
                        (Err(_), Ok(_)) => {
                            info!("Url fallback succeeded, your Ethereum  rpc url {} has been corrected to {}", eth_rpc_url, ipv6_url);
                            web3 = Some(ipv6_web3)
                        },
                        (Ok(_), Ok(_)) => panic!("This should never happen? Why didn't things work the first time?"),
                        (Err(_), Err(_)) => panic!("Could not connect to Ethereum rpc, are you sure it's running and on the specified port? {}", eth_rpc_url)
                    }
                } else if url.port().is_none() || url.scheme() == "http" {
                    let body = url.host_str().unwrap_or_else(|| {
                        panic!("Ethereum rpc url contains no host? {}", eth_rpc_url)
                    });
                    // transparently upgrade to https if available, we can't transparently downgrade for obvious security reasons
                    let https_on_80_url = format!("https://{}:80", body);
                    let https_on_443_url = format!("https://{}:443", body);
                    let https_on_80_web3 = Web3::new(&https_on_80_url, timeout);
                    let https_on_443_web3 = Web3::new(&https_on_443_url, timeout);
                    let https_on_80_test = https_on_80_web3.eth_block_number().await;
                    let https_on_443_test = https_on_443_web3.eth_block_number().await;
                    warn!(
                        "Trying fallback urls {} {}",
                        https_on_443_url, https_on_80_url
                    );
                    match (https_on_80_test, https_on_443_test) {
                        (Ok(_), Err(_)) => {
                            info!("Https upgrade succeeded, your Ethereum rpc url {} has been corrected to {}", eth_rpc_url, https_on_80_url);
                            web3 = Some(https_on_80_web3)
                        },
                        (Err(_), Ok(_)) => {
                            info!("Https upgrade succeeded, your Ethereum rpc url {} has been corrected to {}", eth_rpc_url, https_on_443_url);
                            web3 = Some(https_on_443_web3)
                        },
                        (Ok(_), Ok(_)) => panic!("This should never happen? Why didn't things work the first time?"),
                        (Err(_), Err(_)) => panic!("Could not connect to Ethereum rpc, are you sure it's running and on the specified port? {}", eth_rpc_url)
                    }
                } else {
                    panic!("Could not connect to Ethereum rpc! please check your grpc url {} for errors {:?}", eth_rpc_url, e)
                }
            }
        }
    }

    Connections {
        web3,
        grpc,
        contact,
    }
}

/// Verify that a url has an http or https prefix
fn check_scheme(input: &Url, original_string: &str) {
    if !(input.scheme() == "http" || input.scheme() == "https") {
        panic!(
            "Your url {} has an invalid scheme, please chose http or https",
            original_string
        )
    }
}

/// This function will wait until the Cosmos node is ready, this is intended
/// for situations such as when a node is syncing or when a node is waiting on
/// a halted chain.
pub async fn wait_for_cosmos_node_ready(contact: &Contact) {
    const WAIT_TIME: Duration = Duration::from_secs(10);
    loop {
        let res = contact.get_chain_status().await;
        match res {
            Ok(ChainStatus::Syncing) => {
                info!("Cosmos node is syncing Standing by")
            }
            Ok(ChainStatus::WaitingToStart) => {
                info!("Cosmos node is waiting for the chain to start, Standing by")
            }
            Ok(ChainStatus::Moving { .. }) => {
                break;
            }
            Err(e) => warn!(
                "Could not get syncing status, is your Cosmos node up? {:?}",
                e
            ),
        }
        delay_for(WAIT_TIME).await;
    }
}

/// This function checks if orchestrator and validator addresses were set
pub async fn check_validator_address(
    client: &mut BaseledgerQueryClient<Channel>,
    delegate_orchestrator_address: CosmosAddress,
    prefix: &str,
) {
    let orchestrator_response = client
        .validator_address_by_orchestrator_address(QueryValidatorAddressByOrchestratorAddressRequest {
            orchestrator_address: delegate_orchestrator_address.to_bech32(prefix).unwrap(),
        })
        .await;
    trace!("{:?}", orchestrator_response);
    match orchestrator_response {
        Ok(_) => {
            trace!("Validator found by orch address");
        }
        Err(e) => {
            error!("Your Gravity Orchestrator Cosmos key is incorrect, please double check your phrase. If you can't locate the correct phrase you will need to create a new validator {:?}", e);
            exit(1);
        }
    }
}

/// Checks if a given Coin, used for fees is in the provided address in a sufficient quantity
pub async fn check_for_fee(fee: &Coin, address: CosmosAddress, contact: &Contact) {
    // if we decide to pay no fees it doesn't matter, but we do need some coin balance
    if fee.amount == 0u8.into() {
        if let Err(CosmosGrpcError::NoToken) = contact.get_account_info(address).await {
            error!("Your Orchestrator address has no tokens of any kind. Even if you are paying zero fees this account needs to be 'initialized' by depositing tokens");
            error!(
                "Send the smallest possible unit of any token to {} to resolve this error",
                address
            );
            exit(1);
        }
        return;
    }
    let balances = get_balances_with_retry(address, contact).await;
    for balance in balances {
        if balance.denom.contains(&fee.denom) {
            if balance.amount < fee.amount {
                error!("You have specified a fee that is greater than your balance of that coin! {}{} > {}{} ", fee.amount, fee.denom, balance.amount, balance.denom);
                exit(1);
            } else {
                return;
            }
        }
    }
    error!("You have specified that fees should be paid in {} but account {} has no balance of that token!", fee.denom, address);
    exit(1);
}
