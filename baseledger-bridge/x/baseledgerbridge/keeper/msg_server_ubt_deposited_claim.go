package keeper

import (
	"context"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UbtDepositedClaim(goCtx context.Context, msg *types.MsgUbtDepositedClaim) (*types.MsgUbtDepositedClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.checkOrchestratorValidatorInSet(ctx, msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Could not check orchestrator validator inset")
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Could not create Any value for MsgUbtDepositedClaim")
	}

	_ = any

	return &types.MsgUbtDepositedClaimResponse{}, nil
}

func (k msgServer) checkOrchestratorValidatorInSet(ctx sdk.Context, orchestrator string) error {
	orchestratorAddress, err := sdk.AccAddressFromBech32(orchestrator)
	if err != nil {
		return sdkerrors.Wrap(err, "orchestrator acc address invalid")
	}

	orchValidator, found := k.GetOrchestratorValidator(ctx, orchestratorAddress)
	if !found {
		return sdkerrors.Wrap(sdkerrors.Error{}, "Orchestrator address not set")
	}

	validator := k.StakingKeeper.Validator(ctx, orchValidator.GetOperator())
	if validator == nil || !validator.IsBonded() {
		return sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "Orchestrator validator not in active set")
	}

	return nil
}
