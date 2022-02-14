package keeper_test

import (
	"fmt"
	"strconv"
	"testing"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/x/bridge/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestOrchestratorValidatorAddressMsgServerCreate(t *testing.T) {
	testKeepers := keepertest.SetFiveValidators(t)
	keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)
	sdk.WrapSDKContext(testKeepers.Context)
	validatorAddress := testKeepers.StakingKeeper.GetValidators(testKeepers.Context, 10)[0].OperatorAddress

	fmt.Printf("VALIDATOR ADDRESS %v\n", validatorAddress)
	// for i := 0; i < 5; i++ {
	// 	expected := &types.MsgCreateOrchestratorValidatorAddress{ValidatorAddress: validatorAddress,
	// 		OrchestratorAddress: strconv.Itoa(i),
	// 	}
	// 	_, err := srv.CreateOrchestratorValidatorAddress(wctx, expected)
	// 	require.NoError(t, err)
	// 	rst, found := k.GetOrchestratorValidatorAddress(ctx,
	// 		expected.OrchestratorAddress,
	// 	)
	// 	require.True(t, found)
	// 	require.Equal(t, expected.ValidatorAddress, rst.ValidatorAddress)
	// }
}
