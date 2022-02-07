use deep_space::error::CosmosGrpcError;
use deep_space::private_key::PrivateKey;
use deep_space::Contact;
use deep_space::Msg;
use deep_space::{coin::Coin};
use deep_space::address::Address;
use baseledger_proto::cosmos_sdk_proto::cosmos::base::abci::v1beta1::TxResponse;
use baseledger_proto::baseledger::{MsgUbtDepositedClaim, MsgValidatorPowerChangedClaim, MsgCreateOrchestratorValidatorAddress};
use crate::types::*;
use std::{collections::HashMap, time::Duration};
use num256::Uint256;
use std::str::FromStr;

use crate::get_with_retry::downcast_uint256;


pub const MEMO: &str = "Sent using Baseledger Orchestrator";
pub const TIMEOUT: Duration = Duration::from_secs(60);

/// Send a transaction updating the eth address for the sending
/// Cosmos address. The sending Cosmos address should be a validator
/// this can only be called once! Key rotation code is possible but
/// not currently implemented
pub async fn set_orchestrator_validator_addresses(
    contact: &Contact,
    delegate_cosmos_address: Address,
    private_key: PrivateKey,
) -> Result<TxResponse, CosmosGrpcError> {
    trace!("Updating Orchestrator/Validator addresses");
    let our_valoper_address = private_key
        .to_address(&contact.get_prefix())
        .unwrap()
        // This works so long as the format set by the cosmos hub is maintained
        // having a main prefix followed by a series of titles for specific keys
        // this will not work if that convention is broken. This will be resolved when
        // GRPC exposes prefix endpoints (coming to upstream cosmos sdk soon)
        .to_bech32(format!("{}valoper", contact.get_prefix()))
        .unwrap();

    let msg_set_orch_address = MsgCreateOrchestratorValidatorAddress {
        orchestrator_address: delegate_cosmos_address.to_string(),
        validator_address: our_valoper_address.to_string()
    };

    let msg = Msg::new(
        "/Baseledger.baseledger.bridge.MsgCreateOrchestratorValidatorAddress",
        msg_set_orch_address,
    );

    contact
        .send_message(
            &[msg],
            Some(MEMO.to_string()),
            &[get_zero_fee_coin()],
            Some(TIMEOUT),
            private_key,
        )
        .await
}

#[allow(clippy::too_many_arguments)]
pub async fn send_ethereum_claims(
    contact: &Contact,
    private_key: PrivateKey,
    deposits: Vec<UbtDeposited>,
    power_changes: Vec<ValidatorPowerChangeEvent>,
    ubt_price: f32,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = private_key.to_address(&contact.get_prefix()).unwrap();

    let mut unordered_msgs = HashMap::new();
    
    for deposit in deposits {
        println!("ubt token amount: {}", deposit.amount.to_string());
        println!("ubt price string: {}", ubt_price.to_string());
        let claim = MsgUbtDepositedClaim {
            creator: our_address.to_string(),
            event_nonce: deposit.event_nonce,
            block_height: downcast_uint256(deposit.block_height).unwrap(),
            token_contract: deposit.token.to_string(),
            amount: deposit.amount.to_string(),
            cosmos_receiver: deposit.baseledger_destination_address,
            ethereum_sender: deposit.sender.to_string(),
            ubt_price: ubt_price.to_string(),
        };
        let msg = Msg::new("/Baseledger.baseledger.bridge.MsgUbtDepositedClaim", claim);
        assert!(unordered_msgs.insert(deposit.event_nonce, msg).is_none());
    }

    for power_change in power_changes {
        let claim = MsgValidatorPowerChangedClaim {
            creator: our_address.to_string(),
            event_nonce: power_change.event_nonce,
            block_height: downcast_uint256(power_change.block_height).unwrap(),
            token_contract: power_change.erc20.to_string(),
            amount: power_change.amount.to_string(),
            cosmos_receiver: power_change.destination,
            ethereum_sender: power_change.sender.to_string(),
        };
        let msg = Msg::new("/Baseledger.baseledger.bridge.MsgValidatorPowerChangedClaim", claim);
        assert!(unordered_msgs.insert(power_change.event_nonce, msg).is_none());
    }

    let mut keys = Vec::new();
    for (key, _) in unordered_msgs.iter() {
        keys.push(*key);
    }
    keys.sort_unstable();

    let mut msgs = Vec::new();
    for i in keys {
        msgs.push(unordered_msgs.remove_entry(&i).unwrap().1);
    }

    Ok(contact
        .send_message(&msgs, None, &[get_zero_fee_coin()], Some(TIMEOUT), private_key)
        .await?)
}

// we won't have fees for these transactions
fn get_zero_fee_coin() -> Coin {
    return Coin {
        amount: Uint256::from_str("0").unwrap(),
        denom: "work".to_owned()
    };
}