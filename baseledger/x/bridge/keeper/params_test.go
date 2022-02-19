package keeper_test

import (
	"testing"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	testKeepers := keepertest.SetFiveValidators(t, true)
	params := types.Params{
		WorktokenEurPrice:       "0.1",
		BaseledgerFaucetAddress: keepertest.FaucetAccount.String(),
	}

	testKeepers.BridgeKeeper.SetParams(testKeepers.Context, params)

	require.EqualValues(t, params, testKeepers.BridgeKeeper.GetParams(testKeepers.Context))
}
