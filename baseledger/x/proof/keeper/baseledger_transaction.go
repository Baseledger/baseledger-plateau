package keeper

import (
	"encoding/binary"

	"github.com/Baseledger/baseledger/x/proof/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetBaseledgerTransactionCount get the total number of baseledgerTransaction
func (k Keeper) GetBaseledgerTransactionCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.BaseledgerTransactionCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetBaseledgerTransactionCount set the total number of baseledgerTransaction
func (k Keeper) SetBaseledgerTransactionCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.BaseledgerTransactionCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendBaseledgerTransaction appends a baseledgerTransaction in the store with a new id and update the count
func (k Keeper) AppendBaseledgerTransaction(
	ctx sdk.Context,
	baseledgerTransaction types.BaseledgerTransaction,
) uint64 {
	// Create the baseledgerTransaction
	count := k.GetBaseledgerTransactionCount(ctx)

	// Set the ID of the appended value
	baseledgerTransaction.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionKey))
	appendedValue := k.cdc.MustMarshal(&baseledgerTransaction)
	store.Set(GetBaseledgerTransactionIDBytes(baseledgerTransaction.Id), appendedValue)

	// Update baseledgerTransaction count
	k.SetBaseledgerTransactionCount(ctx, count+1)

	return count
}

// SetBaseledgerTransaction set a specific baseledgerTransaction in the store
func (k Keeper) SetBaseledgerTransaction(ctx sdk.Context, baseledgerTransaction types.BaseledgerTransaction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionKey))
	b := k.cdc.MustMarshal(&baseledgerTransaction)
	store.Set(GetBaseledgerTransactionIDBytes(baseledgerTransaction.Id), b)
}

// GetBaseledgerTransaction returns a baseledgerTransaction from its id
func (k Keeper) GetBaseledgerTransaction(ctx sdk.Context, id uint64) (val types.BaseledgerTransaction, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionKey))
	b := store.Get(GetBaseledgerTransactionIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllBaseledgerTransaction returns all baseledgerTransaction
func (k Keeper) GetAllBaseledgerTransaction(ctx sdk.Context) (list []types.BaseledgerTransaction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BaseledgerTransactionKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.BaseledgerTransaction
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetBaseledgerTransactionIDBytes returns the byte representation of the ID
func GetBaseledgerTransactionIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetBaseledgerTransactionIDFromBytes returns ID in uint64 format from a byte array
func GetBaseledgerTransactionIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
