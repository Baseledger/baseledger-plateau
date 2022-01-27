package keeper

import (
	"context"
	"errors"

	"github.com/Baseledger/baseledger-bridge/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k msgServer) SetOrchestratorAddress(goCtx context.Context, msg *types.MsgSetOrchestratorAddress) (*types.MsgSetOrchestratorAddressResponse, error) {
	// ensure that this passes validation, checks the key validity
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Key not valid")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// check the following, all should be validated in validate basic
	val, e1 := sdk.ValAddressFromBech32(msg.Validator)
	orch, e2 := sdk.AccAddressFromBech32(msg.Orchestrator)
	addr, e3 := types.NewEthAddress(msg.EthAddress)
	if e1 != nil || e2 != nil || e3 != nil {
		return nil, sdkerrors.Wrap(err, "Key not valid")
	}

	// check that the validator does not have an existing key
	_, foundExistingOrchestratorKey := k.GetOrchestratorValidator(ctx, orch)
	_, foundExistingEthAddress := k.GetEthAddressByValidator(ctx, val)

	// ensure that the validator exists
	if k.Keeper.StakingKeeper.Validator(ctx, val) == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, val.String())
	} else if foundExistingOrchestratorKey || foundExistingEthAddress {
		return nil, sdkerrors.Wrap(errors.New("Existing orch found"), val.String())
	}

	// check that neither key is a duplicate
	delegateKeys := k.GetDelegateKeys(ctx)
	for i := range delegateKeys {
		if delegateKeys[i].EthAddress == addr.GetAddress() {
			return nil, sdkerrors.Wrap(err, "Duplicate Ethereum Key")
		}
		if delegateKeys[i].Orchestrator == addr.GetAddress() {
			return nil, sdkerrors.Wrap(err, "Duplicate Orchestrator Key")
		}
	}

	// set the orchestrator address
	k.SetOrchestratorValidator(ctx, val, orch)
	// set the ethereum address
	k.SetEthAddressForValidator(ctx, val, *addr)

	// ctx.EventManager().EmitEvent(
	// 	sdk.NewEvent(
	// 		sdk.EventTypeMessage,
	// 		sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
	// 		sdk.NewAttribute(types.AttributeKeySetOperatorAddr, orch.String()),
	// 	),
	// )

	return &types.MsgSetOrchestratorAddressResponse{}, nil
}
