package keeper

import (
	"time"

	"github.com/Baseledger/baseledger-bridge/x/bridge/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// SetEthAddress sets the ethereum address for a given validator
func (k Keeper) SetEthAddressForValidator(ctx sdk.Context, validator sdk.ValAddress, ethAddr types.EthAddress) {
	if err := sdk.VerifyAddressFormat(validator); err != nil {
		panic(sdkerrors.Wrap(err, "invalid validator address"))
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.GetEthAddressByValidatorKey(validator)), []byte(ethAddr.GetAddress()))
	store.Set([]byte(types.GetValidatorByEthAddressKey(ethAddr)), []byte(validator))
}

// GetEthAddressByValidator returns the eth address for a given gravity validator
func (k Keeper) GetEthAddressByValidator(ctx sdk.Context, validator sdk.ValAddress) (ethAddress *types.EthAddress, found bool) {
	if err := sdk.VerifyAddressFormat(validator); err != nil {
		panic(sdkerrors.Wrap(err, "invalid validator address"))
	}
	store := ctx.KVStore(k.storeKey)
	ethAddr := store.Get([]byte(types.GetEthAddressByValidatorKey(validator)))
	if ethAddr == nil {
		return nil, false
	}

	addr, err := types.NewEthAddress(string(ethAddr))
	if err != nil {
		return nil, false
	}
	return addr, true
}

// GetValidatorByEthAddress returns the validator for a given eth address
func (k Keeper) GetValidatorByEthAddress(ctx sdk.Context, ethAddr types.EthAddress) (validator stakingtypes.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	valAddr := store.Get([]byte(types.GetValidatorByEthAddressKey(ethAddr)))
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
