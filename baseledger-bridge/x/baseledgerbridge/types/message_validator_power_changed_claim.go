package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (msg *MsgValidatorPowerChangedClaim) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
