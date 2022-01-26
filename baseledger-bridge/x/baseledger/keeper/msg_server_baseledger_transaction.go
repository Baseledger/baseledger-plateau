package keeper

import (
	"context"
	"fmt"

	"github.com/Baseledger/baseledger-bridge/x/baseledger/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateBaseledgerTransaction(goCtx context.Context, msg *types.MsgCreateBaseledgerTransaction) (*types.MsgCreateBaseledgerTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var baseledgerTransaction = types.BaseledgerTransaction{
		Creator:                 msg.Creator,
		BaseledgerTransactionId: msg.BaseledgerTransactionId,
		Payload:                 msg.Payload,
		OpCode:                  msg.OpCode,
	}

	id := k.AppendBaseledgerTransaction(
		ctx,
		baseledgerTransaction,
	)

	return &types.MsgCreateBaseledgerTransactionResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateBaseledgerTransaction(goCtx context.Context, msg *types.MsgUpdateBaseledgerTransaction) (*types.MsgUpdateBaseledgerTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var baseledgerTransaction = types.BaseledgerTransaction{
		Creator:                 msg.Creator,
		Id:                      msg.Id,
		BaseledgerTransactionId: msg.BaseledgerTransactionId,
		Payload:                 msg.Payload,
		OpCode:                  msg.OpCode,
	}

	// Checks that the element exists
	val, found := k.GetBaseledgerTransaction(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetBaseledgerTransaction(ctx, baseledgerTransaction)

	return &types.MsgUpdateBaseledgerTransactionResponse{}, nil
}

func (k msgServer) DeleteBaseledgerTransaction(goCtx context.Context, msg *types.MsgDeleteBaseledgerTransaction) (*types.MsgDeleteBaseledgerTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetBaseledgerTransaction(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveBaseledgerTransaction(ctx, msg.Id)

	return &types.MsgDeleteBaseledgerTransactionResponse{}, nil
}
