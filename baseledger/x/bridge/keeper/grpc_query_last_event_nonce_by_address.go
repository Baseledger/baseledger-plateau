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
	var ret types.QueryLastEventNonceByAddressResponse
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, req.Address)
	}
	validator, found := k.GetOrchestratorValidator(ctx, addr)
	if !found {
		return nil, sdkerrors.Wrap(errors.New("Validator not found"), "address")
	}
	if err := sdk.VerifyAddressFormat(validator.GetOperator()); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid validator address")
	}
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, validator.GetOperator())
	ret.EventNonce = lastEventNonce
	return &ret, nil
}
