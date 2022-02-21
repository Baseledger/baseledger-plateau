package keeper

import (
	"context"

	"github.com/Baseledger/baseledger/common"
	"github.com/Baseledger/baseledger/x/proof/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateBaseledgerTransaction(goCtx context.Context, msg *types.MsgCreateBaseledgerTransaction) (*types.MsgCreateBaseledgerTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	txCreatorAddress, err := sdk.AccAddressFromBech32(msg.Creator)

	faucetAccAddress, err := sdk.AccAddressFromBech32(k.bridgeKeeper.GetBaseledgerFaucetAddress(ctx))

	if err != nil {
		panic(err)
	}

	coinFee, err := common.CalcWorkTokenFeeBasedOnPayloadSize(msg.Payload)
	if err != nil {
		k.Logger(ctx).Error("Calculating fee error",
			"err", err.Error(),
			"creator address", txCreatorAddress.String(),
			"fee", coinFee.String())
		return nil, err
	}

	err = k.bankKeeper.SendCoins(ctx, txCreatorAddress, faucetAccAddress, coinFee)
	if err != nil {
		k.Logger(ctx).Error("Send coins to faucet error",
			"err", err.Error(),
			"creator address", txCreatorAddress.String(),
			"faucet address", faucetAccAddress.String(),
			"fee", coinFee.String())
		return nil, err
	}

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
