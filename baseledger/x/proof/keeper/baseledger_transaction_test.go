package keeper_test

import (
	"testing"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/testutil/nullify"
	"github.com/Baseledger/baseledger/x/proof/keeper"
	"github.com/Baseledger/baseledger/x/proof/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNBaseledgerTransaction(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.BaseledgerTransaction {
	items := make([]types.BaseledgerTransaction, n)
	for i := range items {
		items[i].Id = keeper.AppendBaseledgerTransaction(ctx, items[i])
	}
	return items
}

func TestBaseledgerTransactionGet(t *testing.T) {
	testKeepers := keepertest.BaseledgerKeeper(t)
	keeper := testKeepers.ProofKeeper
	ctx := testKeepers.Context
	items := createNBaseledgerTransaction(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetBaseledgerTransaction(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestBaseledgerTransactionGetAll(t *testing.T) {
	testKeepers := keepertest.BaseledgerKeeper(t)
	keeper := testKeepers.ProofKeeper
	ctx := testKeepers.Context
	items := createNBaseledgerTransaction(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllBaseledgerTransaction(ctx)),
	)
}

func TestBaseledgerTransactionCount(t *testing.T) {
	testKeepers := keepertest.BaseledgerKeeper(t)
	keeper := testKeepers.ProofKeeper
	ctx := testKeepers.Context
	items := createNBaseledgerTransaction(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetBaseledgerTransactionCount(ctx))
}
