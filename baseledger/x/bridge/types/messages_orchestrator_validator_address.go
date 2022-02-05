package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateOrchestratorValidatorAddress = "create_orchestrator_validator_address"
	TypeMsgUpdateOrchestratorValidatorAddress = "update_orchestrator_validator_address"
	TypeMsgDeleteOrchestratorValidatorAddress = "delete_orchestrator_validator_address"
)

var _ sdk.Msg = &MsgCreateOrchestratorValidatorAddress{}

func NewMsgCreateOrchestratorValidatorAddress(
	validatorAddress string,
	orchestratorAddress string,

) *MsgCreateOrchestratorValidatorAddress {
	return &MsgCreateOrchestratorValidatorAddress{
		ValidatorAddress:    validatorAddress,
		OrchestratorAddress: orchestratorAddress,
	}
}

func (msg *MsgCreateOrchestratorValidatorAddress) Route() string {
	return RouterKey
}

func (msg *MsgCreateOrchestratorValidatorAddress) Type() string {
	return TypeMsgCreateOrchestratorValidatorAddress
}

func (msg *MsgCreateOrchestratorValidatorAddress) GetSigners() []sdk.AccAddress {
	validatorAddress, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(validatorAddress)}
}

func (msg *MsgCreateOrchestratorValidatorAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateOrchestratorValidatorAddress) ValidateBasic() (err error) {
	if _, err = sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.ValidatorAddress)
	}
	if _, err = sdk.AccAddressFromBech32(msg.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.OrchestratorAddress)
	}
	return nil
}
