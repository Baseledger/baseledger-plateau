package keeper

import (
	"context"

	"github.com/Baseledger/baseledger/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const QUERY_ATTESTATIONS_LIMIT uint64 = 1000

func (k Keeper) Attestations(goCtx context.Context, req *types.QueryAttestationsRequest) (*types.QueryAttestationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	limit := req.Limit
	if limit > QUERY_ATTESTATIONS_LIMIT {
		limit = QUERY_ATTESTATIONS_LIMIT
	}
	attestations := k.GetMostRecentAttestations(ctx, limit)

	return &types.QueryAttestationsResponse{Attestations: attestations}, nil
}
