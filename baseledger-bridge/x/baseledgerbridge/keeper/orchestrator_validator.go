package keeper

import (
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
)

// SetOrchestratorValidator sets the Orchestrator key for a given validator
func (k Keeper) SetOrchestratorValidator(ctx sdk.Context, val sdk.ValAddress, orch sdk.AccAddress) {
	if err := sdk.VerifyAddressFormat(val); err != nil {
		panic(sdkerrors.Wrap(err, "invalid val address"))
	}
	if err := sdk.VerifyAddressFormat(orch); err != nil {
		panic(sdkerrors.Wrap(err, "invalid orch address"))
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.GetOrchestratorAddressKey(orch)), val.Bytes())
}

// GetOrchestratorValidator returns the validator key associated with an orchestrator key
func (k Keeper) GetOrchestratorValidator(ctx sdk.Context, orch sdk.AccAddress) (validator stakingtypes.Validator, found bool) {
	if err := sdk.VerifyAddressFormat(orch); err != nil {
		ctx.Logger().Error("invalid orch address")
		return validator, false
	}
	store := ctx.KVStore(k.storeKey)
	valAddr := store.Get([]byte(types.GetOrchestratorAddressKey(orch)))
	if valAddr == nil {
		return stakingtypes.Validator{
			OperatorAddress: "",
			ConsensusPubkey: &codectypes.Any{
				TypeUrl:              "",
				Value:                []byte{},
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     []byte{},
				XXX_sizecache:        0,
			},
			Jailed:          false,
			Status:          0,
			Tokens:          sdk.Int{},
			DelegatorShares: sdk.Dec{},
			Description: stakingtypes.Description{
				Moniker:         "",
				Identity:        "",
				Website:         "",
				SecurityContact: "",
				Details:         "",
			},
			UnbondingHeight: 0,
			UnbondingTime:   time.Time{},
			Commission: stakingtypes.Commission{
				CommissionRates: stakingtypes.CommissionRates{
					Rate:          sdk.Dec{},
					MaxRate:       sdk.Dec{},
					MaxChangeRate: sdk.Dec{},
				},
				UpdateTime: time.Time{},
			},
			MinSelfDelegation: sdk.Int{},
		}, false
	}
	validator, found = k.StakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return stakingtypes.Validator{
			OperatorAddress: "",
			ConsensusPubkey: &codectypes.Any{
				TypeUrl:              "",
				Value:                []byte{},
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     []byte{},
				XXX_sizecache:        0,
			},
			Jailed:          false,
			Status:          0,
			Tokens:          sdk.Int{},
			DelegatorShares: sdk.Dec{},
			Description: stakingtypes.Description{
				Moniker:         "",
				Identity:        "",
				Website:         "",
				SecurityContact: "",
				Details:         "",
			},
			UnbondingHeight: 0,
			UnbondingTime:   time.Time{},
			Commission: stakingtypes.Commission{
				CommissionRates: stakingtypes.CommissionRates{
					Rate:          sdk.Dec{},
					MaxRate:       sdk.Dec{},
					MaxChangeRate: sdk.Dec{},
				},
				UpdateTime: time.Time{},
			},
			MinSelfDelegation: sdk.Int{},
		}, false
	}

	return validator, true
}
