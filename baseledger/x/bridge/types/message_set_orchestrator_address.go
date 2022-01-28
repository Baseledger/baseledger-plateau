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
	acc, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

func (msg *MsgSetOrchestratorAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetOrchestratorAddress) ValidateBasic() (err error) {
	if _, err = sdk.ValAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}
	if _, err = sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if err := ValidateEthAddress(msg.EthAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}
