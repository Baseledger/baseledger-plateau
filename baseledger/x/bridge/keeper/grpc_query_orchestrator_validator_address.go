package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/Baseledger/baseledger/x/bridge/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OrchestratorValidatorAddressAll(c context.Context, req *types.QueryAllOrchestratorValidatorAddressRequest) (*types.QueryAllOrchestratorValidatorAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var orchestratorValidatorAddresss []types.OrchestratorValidatorAddress
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	orchestratorValidatorAddressStore := prefix.NewStore(store, types.KeyPrefix(types.OrchestratorValidatorAddressKeyPrefix))

	pageRes, err := query.Paginate(orchestratorValidatorAddressStore, req.Pagination, func(key []byte, value []byte) error {
		var orchestratorValidatorAddress types.OrchestratorValidatorAddress
		if err := k.cdc.Unmarshal(value, &orchestratorValidatorAddress); err != nil {
			return err
		}

		orchestratorValidatorAddresss = append(orchestratorValidatorAddresss, orchestratorValidatorAddress)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOrchestratorValidatorAddressResponse{OrchestratorValidatorAddress: orchestratorValidatorAddresss, Pagination: pageRes}, nil
}

func (k Keeper) OrchestratorValidatorAddress(c context.Context, req *types.QueryGetOrchestratorValidatorAddressRequest) (*types.QueryGetOrchestratorValidatorAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetOrchestratorValidatorAddress(
	    ctx,
	    req.OrchestratorAddress,
        )
	if !found {
	    return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetOrchestratorValidatorAddressResponse{OrchestratorValidatorAddress: val}, nil
}