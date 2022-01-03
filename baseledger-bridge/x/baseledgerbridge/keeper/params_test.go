package keeper_test

import (
	"testing"

	testkeeper "github.com/Baseledger/baseledger-bridge/testutil/keeper"
	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.BaseledgerbridgeKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
