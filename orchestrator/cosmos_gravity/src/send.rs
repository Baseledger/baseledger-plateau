use deep_space::error::CosmosGrpcError;
use deep_space::private_key::PrivateKey;
use deep_space::Contact;
use deep_space::Msg;
use deep_space::{coin::Coin};

use gravity_proto::cosmos_sdk_proto::cosmos::base::abci::v1beta1::TxResponse;
use gravity_proto::gravity::MsgUbtDepositedClaim;
use gravity_utils::types::*;
use std::{collections::HashMap, time::Duration};

use crate::utils::downcast_uint256;


pub const MEMO: &str = "Sent using Althea Gravity Bridge Orchestrator";
pub const TIMEOUT: Duration = Duration::from_secs(60);


#[allow(clippy::too_many_arguments)]
pub async fn send_ethereum_claims(
    contact: &Contact,
    private_key: PrivateKey,
    deposits: Vec<SendToCosmosEvent>,
    fee: Coin,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = private_key.to_address(&contact.get_prefix()).unwrap();

    // This sorts oracle messages by event nonce before submitting them. It's not a pretty implementation because
    // we're missing an intermediary layer of abstraction. We could implement 'EventTrait' and then implement sort
    // for it, but then when we go to transform 'EventTrait' objects into GravityMsg enum values we'll have all sorts
    // of issues extracting the inner object from the TraitObject. Likewise we could implement sort of GravityMsg but that
    // would require a truly horrendous (nearly 100 line) match statement to deal with all combinations. That match statement
    // could be reduced by adding two traits to sort against but really this is the easiest option.
    //
    // We index the events by event nonce in an unordered hashmap and then play them back in order into a vec
    let mut unordered_msgs = HashMap::new();
    for deposit in deposits {
        let claim = MsgUbtDepositedClaim {
            creator: our_address.to_string(),
            event_nonce: deposit.event_nonce,
            block_height: downcast_uint256(deposit.block_height).unwrap(),
            token_contract: deposit.erc20.to_string(),
            amount: deposit.amount.to_string(),
            cosmos_receiver: deposit.destination,
            ethereum_sender: deposit.sender.to_string(),
        };
        let msg = Msg::new("/Baseledger.baseledgerbridge.baseledgerbridge.MsgUbtDepositedClaim", claim);
        assert!(unordered_msgs.insert(deposit.event_nonce, msg).is_none());
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
        .send_message(&msgs, None, &[fee], Some(TIMEOUT), private_key)
        .await?)
}