package keeper

import (
	"context"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UbtDepositedClaim(goCtx context.Context, msg *types.MsgUbtDepositedClaim) (*types.MsgUbtDepositedClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgUbtDepositedClaimResponse{}, nil
}
