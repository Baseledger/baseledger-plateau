package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/x/bridge/keeper"
	"github.com/Baseledger/baseledger/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestOrchestratorValidatorAddressMsgServerCreate(t *testing.T) {
	testKeepers := keepertest.SetFiveValidators(t, false)
	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)
	wctx := sdk.WrapSDKContext(testKeepers.Context)
	validatorAddress := testKeepers.StakingKeeper.GetValidators(testKeepers.Context, 10)[0].OperatorAddress

	msg := &types.MsgCreateOrchestratorValidatorAddress{
		ValidatorAddress:    validatorAddress,
		OrchestratorAddress: keepertest.AccAddrs[0].String(),
	}

	_, err := srv.CreateOrchestratorValidatorAddress(wctx, msg)
	require.NoError(t, err)
	rst, found := testKeepers.BridgeKeeper.GetOrchestratorValidatorAddress(testKeepers.Context,
		msg.OrchestratorAddress,
	)
	require.True(t, found)
	require.Equal(t, msg.ValidatorAddress, rst.ValidatorAddress)
}
