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
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "orchestrator already set")
	}

	val, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, "validator address format error")
	}

	// check if validator with this address exists in active set
	validator := k.StakingKeeper.Validator(ctx, val)
	if validator == nil || !validator.IsBonded() || validator.IsJailed() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, "validator not found or not in active set")
	}

	list := k.Keeper.GetAllOrchestratorValidatorAddresses(ctx)

	for i := range list {
		if list[i].ValidatorAddress == msg.ValidatorAddress {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator already set")
		}
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
