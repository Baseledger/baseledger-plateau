package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace

		StakingKeeper *stakingkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// TODO: skos add this implementation
// Checks if the provided Ethereum address is on the Governance blacklist
func (k Keeper) IsOnBlacklist(ctx sdk.Context, addr types.EthAddress) bool {
	// params := k.GetParams(ctx)
	// // Checks the address if it's inside the blacklisted address list and marks
	// // if it's inside the list.
	// for index := 0; index < len(params.EthereumBlacklist); index++ {
	// 	baddr, err := types.NewEthAddress(params.EthereumBlacklist[index])
	// 	if err != nil {
	// 		// this should not be possible we validate on genesis load
	// 		panic("unvalidated black list address!")
	// 	}
	// 	if *baddr == addr {
	// 		return true
	// 	}
	// }
	return false
}
