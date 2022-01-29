package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const TypeMsgValidatorPowerChangedClaim = "validator_power_changed_claim"

var _ sdk.Msg = &MsgValidatorPowerChangedClaim{}

func NewMsgValidatorPowerChangedClaim(creator string, eventNonce uint64, blockHeight uint64, tokenContract string, amount sdk.Int, ethereumSender string, cosmosReceiver string) *MsgValidatorPowerChangedClaim {
	return &MsgValidatorPowerChangedClaim{
		Creator:        creator,
		EventNonce:     eventNonce,
		BlockHeight:    blockHeight,
		TokenContract:  tokenContract,
		Amount:         amount,
		EthereumSender: ethereumSender,
		CosmosReceiver: cosmosReceiver,
	}
}

func (msg *MsgValidatorPowerChangedClaim) Route() string {
	return RouterKey
}

func (msg *MsgValidatorPowerChangedClaim) Type() string {
	return TypeMsgValidatorPowerChangedClaim
}

func (msg *MsgValidatorPowerChangedClaim) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgValidatorPowerChangedClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetType returns the type of the claim
func (msg *MsgValidatorPowerChangedClaim) GetType() ClaimType {
	return CLAIM_VALIDATOR_POWER_CHANGED
}

// ValidateBasic performs stateless checks
func (msg *MsgValidatorPowerChangedClaim) ValidateBasic() error {
	if err := ValidateEthAddress(msg.EthereumSender); err != nil {
		return sdkerrors.Wrap(err, "eth sender")
	}
	if err := ValidateEthAddress(msg.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "orchestrator")
	}
	if _, err := sdk.AccAddressFromBech32(msg.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(err, "cosmos receiver")
	}
	if msg.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

func (msg MsgValidatorPowerChangedClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgValidatorPowerChangedClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return val
}

func (msg *MsgValidatorPowerChangedClaim) GetUbtPriceAsInt() sdk.Int {
	// TODO: Ognjen - Cleanup, her just to make the build work
	panic("Not implemented")
}

// Hash implements BridgeDeposit.Hash
// modify this with care as it is security sensitive. If an element of the claim is not in this hash a single hostile validator
// could engineer a hash collision and execute a version of the claim with any unhashed data changed to benefit them.
// note that the Orchestrator is the only field excluded from this hash, this is because that value is used higher up in the store
// structure for who has made what claim and is verified by the msg ante-handler for signatures
func (msg *MsgValidatorPowerChangedClaim) ClaimHash() ([]byte, error) {
	path := fmt.Sprintf("%d/%d/%s/%s/%s/%s", msg.EventNonce, msg.BlockHeight, msg.TokenContract, msg.Amount, msg.EthereumSender, msg.CosmosReceiver)
	return tmhash.Sum([]byte(path)), nil
}
