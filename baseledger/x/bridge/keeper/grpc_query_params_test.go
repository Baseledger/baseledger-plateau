package keeper_test

import (
	"testing"

	testkeeper "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/testutil/sample"
	"github.com/Baseledger/baseledger/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	testKeepers := testkeeper.BaseledgerbridgeKeeper(t)
	wctx := sdk.WrapSDKContext(testKeepers.Context)
	params := types.Params{
		BaseledgerFaucetAddress: sample.AccAddress(),
		WorktokenEurPrice:       "0.1",
	}
	testKeepers.BridgeKeeper.SetParams(testKeepers.Context, params)

	response, err := testKeepers.BridgeKeeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
