//! This file parses the Baseledger contract ethereum events. Note that there is no Ethereum ABI unpacking implementation. Instead each event
//! is parsed directly from it's binary representation. This is technical debt within this implementation. It's quite easy to parse any
//! individual event manually but a generic decoder can be quite challenging to implement. A proper implementation would probably closely
//! mirror Serde and perhaps even become a serde crate for Ethereum ABI decoding
//! For now reference the ABI encoding document here https://docs.soliditylang.org/en/v0.8.3/abi-spec.html

// TODO this file needs static assertions that prevent it from compiling on 16 bit systems.
// we assume a system bit width of at least 32

use crate::error::OrchestratorError;
use clarity::Address as EthAddress;
use deep_space::utils::bytes_to_hex_str;
use deep_space::Address as CosmosAddress;
use num256::Uint256;
use web30::types::Log;

/// Used to limit the length of variable length user provided inputs like
/// ERC20 names and deposit destination strings
const ONE_MEGABYTE: usize = 1000usize.pow(3);

/*
    UBTSplitter events have same signature, so we are using generic event decoding for now
    event UbtDeposited(
        address indexed token,
        address indexed sender,
        string baseledgerDestinationAddress,
        uint256 tokenAmount,
        uint256 lastEventNonce
    );

    event PayeeUpdated(
        address indexed token,
        address indexed revenueAddress,
        string baseledgerValidatorAddress,
        uint256 shares,
        uint256 lastEventNonce
    );
*/

#[derive(Serialize, Deserialize, Debug, Clone, Eq, PartialEq, Hash)]
pub struct UbtSplitterEvent {
    /// UBT contract address (TODO: do we need it?)
    pub token: EthAddress,
    /// Sender of UBT deposit or payee address
    pub payee_address: EthAddress,
    /// Baseledger destination address - this is a raw value from the Ethereum contract
    /// and therefore could be provided by an attacker. If the string is valid
    /// utf-8 it will be included here, if it is invalid utf8 we will provide
    /// an empty string. Values over 1mb of text are not permitted and will also
    /// be presented as empty
    pub baseledger_destination_address: String,
    /// the validated destination is the destination string parsed and interpreted
    /// as a valid Bech32 Cosmos address, if this is not possible the value is none
    pub validated_destination: Option<CosmosAddress>,
    /// The amount of the erc20 token that is being sent
    pub amount: Uint256,
    /// The transaction's nonce, used to make sure there can be no accidental duplication
    pub event_nonce: u64,
    /// The block height this event occurred at
    pub block_height: Uint256,
}

/// struct for holding the data encoded fields
/// of a send to Cosmos event for unit testing
#[derive(Eq, PartialEq, Debug)]
struct UbtSplitterEventData {
    /// Baseledger destination address, None for an invalid deposit address
    pub baseledger_destination_address: String,
    /// The amount of the erc20 token that is being sent
    pub amount: Uint256,
    /// The transaction's nonce, used to make sure there can be no accidental duplication
    pub event_nonce: Uint256,
}

impl UbtSplitterEvent {
    pub fn from_log(input: &Log) -> Result<UbtSplitterEvent, OrchestratorError> {
        let topics = (input.topics.get(1), input.topics.get(2));
        if let (Some(erc20_data), Some(sender_data)) = topics {
            let token = EthAddress::from_slice(&erc20_data[12..32])?;
            let payee_address = EthAddress::from_slice(&sender_data[12..32])?;
            let block_height = if let Some(bn) = input.block_number.clone() {
                if bn > u64::MAX.into() {
                    return Err(OrchestratorError::InvalidEventLogError(
                        "Block height overflow! probably incorrect parsing".to_string(),
                    ));
                } else {
                    bn
                }
            } else {
                return Err(OrchestratorError::InvalidEventLogError(
                    "Log does not have block number, we only search logs already in blocks?"
                        .to_string(),
                ));
            };

            let data = UbtSplitterEvent::decode_data_bytes(&input.data)?;
            if data.event_nonce > u64::MAX.into() || block_height > u64::MAX.into() {
                Err(OrchestratorError::InvalidEventLogError(
                    "Event nonce overflow, probably incorrect parsing".to_string(),
                ))
            } else {
                let event_nonce: u64 = data.event_nonce.to_string().parse().unwrap();
                let validated_destination = match data.baseledger_destination_address.parse() {
                    Ok(v) => Some(v),
                    Err(_) => {
                        if data.baseledger_destination_address.len() < 1000 {
                            warn!("Event nonce {} sends tokens to {} which is invalid bech32, these funds will be allocated to the faucet account", event_nonce, data.baseledger_destination_address);
                        } else {
                            warn!("Event nonce {} sends tokens to a destination which is invalid bech32, these funds will be allocated to the faucet account", event_nonce);
                        }
                        None
                    }
                };
                Ok(UbtSplitterEvent {
                    token,
                    payee_address,
                    baseledger_destination_address: data.baseledger_destination_address,
                    validated_destination,
                    amount: data.amount,
                    event_nonce,
                    block_height,
                })
            }
        } else {
            Err(OrchestratorError::InvalidEventLogError(
                "Too few topics".to_string(),
            ))
        }
    }
    fn decode_data_bytes(data: &[u8]) -> Result<UbtSplitterEventData, OrchestratorError> {
        if data.len() < 4 * 32 {
            return Err(OrchestratorError::InvalidEventLogError(
                "too short for UbtSplitterEventData".to_string(),
            ));
        }

        let amount = Uint256::from_bytes_be(&data[32..64]);
        let event_nonce = Uint256::from_bytes_be(&data[64..96]);

        // discard words three and four which contain the data type and length
        let destination_str_len_start = 3 * 32;
        let destination_str_len_end = 4 * 32;
        let destination_str_len =
            Uint256::from_bytes_be(&data[destination_str_len_start..destination_str_len_end]);
        if destination_str_len > u32::MAX.into() {
            return Err(OrchestratorError::InvalidEventLogError(
                "denom length overflow, probably incorrect parsing".to_string(),
            ));
        }
        let destination_str_len: usize = destination_str_len.to_string().parse().unwrap();

        let destination_str_start = 4 * 32;
        let destination_str_end = destination_str_start + destination_str_len;

        if data.len() < destination_str_end {
            return Err(OrchestratorError::InvalidEventLogError(
                "Incorrect length for dynamic data".to_string(),
            ));
        }

        let destination = &data[destination_str_start..destination_str_end];

        let dest = String::from_utf8(destination.to_vec());
        if dest.is_err() {
            if destination.len() < 1000 {
                warn!("Event nonce {} sends tokens to {} which is invalid utf-8, these funds will be allocated to the faucet account", event_nonce, bytes_to_hex_str(destination));
            } else {
                warn!("Event nonce {} sends tokens to a destination that is invalid utf-8, these funds will be allocated to the faucet account", event_nonce);
            }
            return Ok(UbtSplitterEventData {
                baseledger_destination_address: String::new(),
                event_nonce,
                amount,
            });
        }
        // whitespace can not be a valid part of a bech32 address, so we can safely trim it
        let dest = dest.unwrap().trim().to_string();

        if dest.as_bytes().len() > ONE_MEGABYTE {
            warn!("Event nonce {} sends tokens to a destination that exceeds the length limit, these funds will be allocated to the faucet account", event_nonce);
            Ok(UbtSplitterEventData {
                baseledger_destination_address: String::new(),
                event_nonce,
                amount,
            })
        } else {
            Ok(UbtSplitterEventData {
                baseledger_destination_address: dest,
                event_nonce,
                amount,
            })
        }
    }
    pub fn from_logs(input: &[Log]) -> Result<Vec<UbtSplitterEvent>, OrchestratorError> {
        let mut res = Vec::new();
        for item in input {
            res.push(UbtSplitterEvent::from_log(item)?);
        }
        Ok(res)
    }
    /// returns all values in the array with event nonces greater
    /// than the provided value
    pub fn filter_by_event_nonce(event_nonce: u64, input: &[Self]) -> Vec<Self> {
        let mut ret = Vec::new();
        for item in input {
            if item.event_nonce > event_nonce {
                ret.push(item.clone())
            }
        }
        ret
    }
}