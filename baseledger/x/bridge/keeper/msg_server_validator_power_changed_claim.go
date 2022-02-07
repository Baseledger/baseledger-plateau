package keeper

import (
	"context"

	"github.com/Baseledger/baseledger/x/bridge/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ValidatorPowerChangedClaim(goCtx context.Context, msg *types.MsgValidatorPowerChangedClaim) (*types.MsgValidatorPowerChangedClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	val := k.GetOrchestratorValidator(ctx, msg.Creator)
	if val == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, "Validator not found")
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
