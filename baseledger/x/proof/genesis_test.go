package baseledger_test

import (
	"testing"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/testutil/nullify"
	baseledger "github.com/Baseledger/baseledger/x/proof"
	"github.com/Baseledger/baseledger/x/proof/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		BaseledgerTransactionList: []types.BaseledgerTransaction{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		BaseledgerTransactionCount: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BaseledgerKeeper(t)
	baseledger.InitGenesis(ctx, *k, genesisState)
	got := baseledger.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.BaseledgerTransactionList, got.BaseledgerTransactionList)
	require.Equal(t, genesisState.BaseledgerTransactionCount, got.BaseledgerTransactionCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
