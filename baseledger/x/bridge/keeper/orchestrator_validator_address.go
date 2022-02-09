package keeper

import (
	"github.com/Baseledger/baseledger/logger"
	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// SetOrchestratorValidatorAddress set a specific orchestratorValidatorAddress in the store from its index
func (k Keeper) SetOrchestratorValidatorAddress(ctx sdk.Context, orchestratorValidatorAddress types.OrchestratorValidatorAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OrchestratorValidatorAddressKeyPrefix))
	b := k.cdc.MustMarshal(&orchestratorValidatorAddress)
	store.Set(types.OrchestratorValidatorAddressKey(
		orchestratorValidatorAddress.OrchestratorAddress,
	), b)
}

// GetOrchestratorValidatorAddress returns a orchestratorValidatorAddress from its index
func (k Keeper) GetOrchestratorValidatorAddress(
	ctx sdk.Context,
	orchestratorAddress string,

) (val types.OrchestratorValidatorAddress, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OrchestratorValidatorAddressKeyPrefix))

	b := store.Get(types.OrchestratorValidatorAddressKey(
		orchestratorAddress,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// check if validator key is set, and if validator exists and it is active
func (k Keeper) GetOrchestratorValidator(ctx sdk.Context, orchestratorAddress string) stakingtypes.ValidatorI {
	_, err := sdk.AccAddressFromBech32(orchestratorAddress)
	if err != nil {
		logger.Errorf("Orchestrator address format error %v\n", err.Error())
		return nil
	}

	orchValAddress, found := k.GetOrchestratorValidatorAddress(ctx, orchestratorAddress)
	if !found {
		logger.Errorf("Orchestrator validator address not found for key %v\n", orchestratorAddress)
		return nil
	}

	valAddr, err := sdk.ValAddressFromBech32(orchValAddress.ValidatorAddress)
	if err != nil {
		logger.Errorf("Validator address format error %v\n", err.Error())
		return nil
	}

	validator := k.StakingKeeper.Validator(ctx, valAddr)

	if validator == nil || !validator.IsBonded() {
		logger.Errorf("Validator not in active set")
		return nil
	}

	return validator
}

// GetAllOrchestratorValidatorAddresses returns all orchestratorValidatorAddresses
func (k Keeper) GetAllOrchestratorValidatorAddresses(ctx sdk.Context) (list []types.OrchestratorValidatorAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OrchestratorValidatorAddressKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OrchestratorValidatorAddress
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
