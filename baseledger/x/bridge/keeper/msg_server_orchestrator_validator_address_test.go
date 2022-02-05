package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/x/bridge/keeper"
	"github.com/Baseledger/baseledger/x/bridge/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestOrchestratorValidatorAddressMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.BaseledgerbridgeKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	validatorAddress := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateOrchestratorValidatorAddress{ValidatorAddress: validatorAddress,
			OrchestratorAddress: strconv.Itoa(i),
		}
		_, err := srv.CreateOrchestratorValidatorAddress(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetOrchestratorValidatorAddress(ctx,
			expected.OrchestratorAddress,
		)
		require.True(t, found)
		require.Equal(t, expected.ValidatorAddress, rst.ValidatorAddress)
	}
}
