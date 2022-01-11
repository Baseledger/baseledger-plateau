package keeper

import (
	"context"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DelegateKeysByEthAddress(goCtx context.Context, req *types.QueryDelegateKeysByEthAddressRequest) (*types.QueryDelegateKeysByEthAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Process the query
	_ = ctx

	return &types.QueryDelegateKeysByEthAddressResponse{}, nil
}
