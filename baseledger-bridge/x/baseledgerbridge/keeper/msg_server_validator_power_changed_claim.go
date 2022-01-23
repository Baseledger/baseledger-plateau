package keeper

import (
	"context"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ValidatorPowerChangedClaim(goCtx context.Context, msg *types.MsgValidatorPowerChangedClaim) (*types.MsgValidatorPowerChangedClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.checkOrchestratorValidatorInSet(ctx, msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Could not check orchestrator validator inset")
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Could not create Any value for MsgValidatorPowerChangedClaim")
	}

	err = k.claimHandlerCommon(ctx, any, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgValidatorPowerChangedClaimResponse{}, nil
}
