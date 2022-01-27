package baseledger

import (
	"github.com/Baseledger/baseledger/x/proof/keeper"
	"github.com/Baseledger/baseledger/x/proof/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the baseledgerTransaction
	for _, elem := range genState.BaseledgerTransactionList {
		k.SetBaseledgerTransaction(ctx, elem)
	}

	// Set baseledgerTransaction count
	k.SetBaseledgerTransactionCount(ctx, genState.BaseledgerTransactionCount)
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.BaseledgerTransactionList = k.GetAllBaseledgerTransaction(ctx)
	genesis.BaseledgerTransactionCount = k.GetBaseledgerTransactionCount(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
