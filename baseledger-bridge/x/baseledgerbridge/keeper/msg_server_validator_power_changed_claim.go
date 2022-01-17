package keeper

import (
	"context"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ValidatorPowerChangedClaim(goCtx context.Context, msg *types.MsgValidatorPowerChangedClaim) (*types.MsgValidatorPowerChangedClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgValidatorPowerChangedClaimResponse{}, nil
}
