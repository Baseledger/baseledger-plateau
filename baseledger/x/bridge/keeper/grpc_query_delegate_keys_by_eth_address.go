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

func (k Keeper) DelegateKeysByEthAddress(goCtx context.Context, req *types.QueryDelegateKeysByEthAddressRequest) (*types.QueryDelegateKeysByEthAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if err := types.ValidateEthAddress(req.EthAddress); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid eth address")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	keys := k.GetDelegateKeys(ctx)

	for _, key := range keys {
		if req.EthAddress == key.EthAddress {
			return &types.QueryDelegateKeysByEthAddressResponse{
				ValidatorAddress:    key.Validator,
				OrchestratorAddress: key.Orchestrator,
			}, nil
		}
	}

	return nil, sdkerrors.Wrap(errors.New("Could not find keys by eth address"), "No validator")
}
