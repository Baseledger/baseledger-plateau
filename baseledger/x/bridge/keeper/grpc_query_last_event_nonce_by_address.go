package keeper

import (
	"context"
	"errors"

	"github.com/Baseledger/baseledger/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) LastEventNonceByAddress(goCtx context.Context, req *types.QueryLastEventNonceByAddressRequest) (*types.QueryLastEventNonceByAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	_, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, req.Address)
	}
	orchValAddr := k.GetOrchestratorValidator(ctx, req.Address)
	if orchValAddr == nil {
		return nil, sdkerrors.Wrap(errors.New("Validator not found"), "address")
	}

	return &types.QueryLastEventNonceByAddressResponse{
		EventNonce: k.GetLastEventNonceByValidator(ctx, orchValAddr.GetOperator()),
	}, nil
}
