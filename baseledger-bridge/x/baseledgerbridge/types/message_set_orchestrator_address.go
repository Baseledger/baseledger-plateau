package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetOrchestratorAddress = "set_orchestrator_address"

var _ sdk.Msg = &MsgSetOrchestratorAddress{}

func NewMsgSetOrchestratorAddress(validator string, orchestrator string, ethAddress string) *MsgSetOrchestratorAddress {
	return &MsgSetOrchestratorAddress{
		Validator:    validator,
		Orchestrator: orchestrator,
		EthAddress:   ethAddress,
	}
}

func (msg *MsgSetOrchestratorAddress) Route() string {
	return RouterKey
}

func (msg *MsgSetOrchestratorAddress) Type() string {
	return TypeMsgSetOrchestratorAddress
}

func (msg *MsgSetOrchestratorAddress) GetSigners() []sdk.AccAddress {
	validator, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{validator}
}

func (msg *MsgSetOrchestratorAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetOrchestratorAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}
	return nil
}
