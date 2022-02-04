package keeper

import (
	"context"

	"github.com/Baseledger/baseledger/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateOrchestratorValidatorAddress(goCtx context.Context, msg *types.MsgCreateOrchestratorValidatorAddress) (*types.MsgCreateOrchestratorValidatorAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetOrchestratorValidatorAddress(
		ctx,
		msg.OrchestratorAddress,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	val, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, "validator address format error")
	}

	// check if validator with this address exists
	if k.Keeper.StakingKeeper.Validator(ctx, val) == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, "validator not found")
	}

	var orchestratorValidatorAddress = types.OrchestratorValidatorAddress{
		ValidatorAddress:    msg.ValidatorAddress,
		OrchestratorAddress: msg.OrchestratorAddress,
	}

	k.SetOrchestratorValidatorAddress(
		ctx,
		orchestratorValidatorAddress,
	)
	return &types.MsgCreateOrchestratorValidatorAddressResponse{}, nil
}

func (k msgServer) UpdateOrchestratorValidatorAddress(goCtx context.Context, msg *types.MsgUpdateOrchestratorValidatorAddress) (*types.MsgUpdateOrchestratorValidatorAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetOrchestratorValidatorAddress(
		ctx,
		msg.OrchestratorAddress,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg validatorAddress is the same as the current owner
	if msg.ValidatorAddress != valFound.ValidatorAddress {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var orchestratorValidatorAddress = types.OrchestratorValidatorAddress{
		ValidatorAddress:    msg.ValidatorAddress,
		OrchestratorAddress: msg.OrchestratorAddress,
	}

	k.SetOrchestratorValidatorAddress(ctx, orchestratorValidatorAddress)

	return &types.MsgUpdateOrchestratorValidatorAddressResponse{}, nil
}

func (k msgServer) DeleteOrchestratorValidatorAddress(goCtx context.Context, msg *types.MsgDeleteOrchestratorValidatorAddress) (*types.MsgDeleteOrchestratorValidatorAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetOrchestratorValidatorAddress(
		ctx,
		msg.OrchestratorAddress,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg validatorAddress is the same as the current owner
	if msg.ValidatorAddress != valFound.ValidatorAddress {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveOrchestratorValidatorAddress(
		ctx,
		msg.OrchestratorAddress,
	)

	return &types.MsgDeleteOrchestratorValidatorAddressResponse{}, nil
}
