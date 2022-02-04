package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
)

// SetOrchestratorValidatorAddress set a specific orchestratorValidatorAddress in the store from its index
func (k Keeper) SetOrchestratorValidatorAddress(ctx sdk.Context, orchestratorValidatorAddress types.OrchestratorValidatorAddress) {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OrchestratorValidatorAddressKeyPrefix))
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

// RemoveOrchestratorValidatorAddress removes a orchestratorValidatorAddress from the store
func (k Keeper) RemoveOrchestratorValidatorAddress(
    ctx sdk.Context,
    orchestratorAddress string,
    
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OrchestratorValidatorAddressKeyPrefix))
	store.Delete(types.OrchestratorValidatorAddressKey(
	    orchestratorAddress,
    ))
}

// GetAllOrchestratorValidatorAddress returns all orchestratorValidatorAddress
func (k Keeper) GetAllOrchestratorValidatorAddress(ctx sdk.Context) (list []types.OrchestratorValidatorAddress) {
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
