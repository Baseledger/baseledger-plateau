package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUbtDepositedClaim = "ubt_deposited_claim"

var _ sdk.Msg = &MsgUbtDepositedClaim{}

func NewMsgUbtDepositedClaim(creator string, eventNonce uint64, blockHeight uint64, tokenContract string, amount string, ethereumSender string, cosmosReceiver string) *MsgUbtDepositedClaim {
	return &MsgUbtDepositedClaim{
		Creator:        creator,
		EventNonce:     eventNonce,
		BlockHeight:    blockHeight,
		TokenContract:  tokenContract,
		Amount:         amount,
		EthereumSender: ethereumSender,
		CosmosReceiver: cosmosReceiver,
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

func (msg *MsgUbtDepositedClaim) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
