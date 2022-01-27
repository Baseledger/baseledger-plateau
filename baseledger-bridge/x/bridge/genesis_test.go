package bridge_test

import (
	"testing"

	keepertest "github.com/Baseledger/baseledger-bridge/testutil/keeper"
	"github.com/Baseledger/baseledger-bridge/testutil/nullify"
	baseledgerbridge "github.com/Baseledger/baseledger-bridge/x/bridge"
	"github.com/Baseledger/baseledger-bridge/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BaseledgerbridgeKeeper(t)
	baseledgerbridge.InitGenesis(ctx, *k, genesisState)
	got := baseledgerbridge.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
