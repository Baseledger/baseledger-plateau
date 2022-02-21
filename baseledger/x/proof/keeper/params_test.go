package keeper_test

import (
	"testing"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/x/proof/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	testKeepers := keepertest.BaseledgerKeeper(t)
	k := testKeepers.ProofKeeper
	ctx := testKeepers.Context
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
