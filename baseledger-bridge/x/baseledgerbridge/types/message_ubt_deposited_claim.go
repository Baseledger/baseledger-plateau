package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const TypeMsgUbtDepositedClaim = "ubt_deposited_claim"

var _ sdk.Msg = &MsgUbtDepositedClaim{}

func NewMsgUbtDepositedClaim(creator string, eventNonce uint64, blockHeight uint64, tokenContract string, amount sdk.Int, ethereumSender string, cosmosReceiver string, ubtPrice sdk.Dec) *MsgUbtDepositedClaim {
	return &MsgUbtDepositedClaim{
		Creator:        creator,
		EventNonce:     eventNonce,
		BlockHeight:    blockHeight,
		TokenContract:  tokenContract,
		Amount:         amount,
		EthereumSender: ethereumSender,
		CosmosReceiver: cosmosReceiver,
		Price:          ubtPrice,
	}
}

func (msg *MsgUbtDepositedClaim) Route() string {
	return RouterKey
}

func (msg *MsgUbtDepositedClaim) Type() string {
	return TypeMsgUbtDepositedClaim
}

func (msg *MsgUbtDepositedClaim) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUbtDepositedClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetType returns the type of the claim
func (msg *MsgUbtDepositedClaim) GetType() ClaimType {
	return CLAIM_UBT_DEPOSITED
}

// ValidateBasic performs stateless checks
func (msg *MsgUbtDepositedClaim) ValidateBasic() error {
	if err := ValidateEthAddress(msg.EthereumSender); err != nil {
		return sdkerrors.Wrap(err, "eth sender")
	}
	if err := ValidateEthAddress(msg.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "orchestrator")
	}
	if _, err := IBCAddressFromBech32(msg.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(err, "cosmos receiver")
	}
	if msg.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

func (msg MsgUbtDepositedClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgUbtDepositedClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return val
}

// Hash implements BridgeDeposit.Hash
// modify this with care as it is security sensitive. If an element of the claim is not in this hash a single hostile validator
// could engineer a hash collision and execute a version of the claim with any unhashed data changed to benefit them.
// note that the Orchestrator is the only field excluded from this hash, this is because that value is used higher up in the store
// structure for who has made what claim and is verified by the msg ante-handler for signatures
func (msg *MsgUbtDepositedClaim) ClaimHash() ([]byte, error) {
	path := fmt.Sprintf("%d/%d/%s/%s/%s/%s/%s", msg.EventNonce, msg.BlockHeight, msg.TokenContract, msg.Amount, msg.EthereumSender, msg.CosmosReceiver, msg.Price)
	return tmhash.Sum([]byte(path)), nil
}
