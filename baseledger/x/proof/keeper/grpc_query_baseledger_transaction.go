package keeper

import (
	"context"

	"github.com/Baseledger/baseledger/x/proof/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) BaseledgerTransactionAll(c context.Context, req *types.QueryAllBaseledgerTransactionRequest) (*types.QueryAllBaseledgerTransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var baseledgerTransactions []types.BaseledgerTransaction
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	baseledgerTransactionStore := prefix.NewStore(store, types.KeyPrefix(types.BaseledgerTransactionKey))

	pageRes, err := query.Paginate(baseledgerTransactionStore, req.Pagination, func(key []byte, value []byte) error {
		var baseledgerTransaction types.BaseledgerTransaction
		if err := k.cdc.Unmarshal(value, &baseledgerTransaction); err != nil {
			return err
		}

		baseledgerTransactions = append(baseledgerTransactions, baseledgerTransaction)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllBaseledgerTransactionResponse{BaseledgerTransaction: baseledgerTransactions, Pagination: pageRes}, nil
}

func (k Keeper) BaseledgerTransaction(c context.Context, req *types.QueryGetBaseledgerTransactionRequest) (*types.QueryGetBaseledgerTransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	baseledgerTransaction, found := k.GetBaseledgerTransaction(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetBaseledgerTransactionResponse{BaseledgerTransaction: baseledgerTransaction}, nil
}
