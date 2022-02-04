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

var _ sdk.Msg = &MsgUpdateOrchestratorValidatorAddress{}

func NewMsgUpdateOrchestratorValidatorAddress(
	validatorAddress string,
	orchestratorAddress string,

) *MsgUpdateOrchestratorValidatorAddress {
	return &MsgUpdateOrchestratorValidatorAddress{
		ValidatorAddress:    validatorAddress,
		OrchestratorAddress: orchestratorAddress,
	}
}

func (msg *MsgUpdateOrchestratorValidatorAddress) Route() string {
	return RouterKey
}

func (msg *MsgUpdateOrchestratorValidatorAddress) Type() string {
	return TypeMsgUpdateOrchestratorValidatorAddress
}

func (msg *MsgUpdateOrchestratorValidatorAddress) GetSigners() []sdk.AccAddress {
	validatorAddress, err := sdk.AccAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{validatorAddress}
}

func (msg *MsgUpdateOrchestratorValidatorAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateOrchestratorValidatorAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validatorAddress address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteOrchestratorValidatorAddress{}

func NewMsgDeleteOrchestratorValidatorAddress(
	validatorAddress string,
	orchestratorAddress string,

) *MsgDeleteOrchestratorValidatorAddress {
	return &MsgDeleteOrchestratorValidatorAddress{
		ValidatorAddress:    validatorAddress,
		OrchestratorAddress: orchestratorAddress,
	}
}
func (msg *MsgDeleteOrchestratorValidatorAddress) Route() string {
	return RouterKey
}

func (msg *MsgDeleteOrchestratorValidatorAddress) Type() string {
	return TypeMsgDeleteOrchestratorValidatorAddress
}

func (msg *MsgDeleteOrchestratorValidatorAddress) GetSigners() []sdk.AccAddress {
	validatorAddress, err := sdk.AccAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{validatorAddress}
}

func (msg *MsgDeleteOrchestratorValidatorAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteOrchestratorValidatorAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validatorAddress address (%s)", err)
	}
	return nil
}
