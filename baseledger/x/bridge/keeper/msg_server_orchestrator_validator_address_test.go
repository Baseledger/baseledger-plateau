package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/x/bridge/keeper"
	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestOrchestratorValidatorAddress_CreateSuccess(t *testing.T) {
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

func TestOrchestratorValidatorAddress_OrchestratorAlreadySet(t *testing.T) {
	// flag true will also set orch validators for all 5 validators
	testKeepers := keepertest.SetFiveValidators(t, true)
	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)
	wctx := sdk.WrapSDKContext(testKeepers.Context)
	validatorAddress := testKeepers.StakingKeeper.GetValidators(testKeepers.Context, 10)[0].OperatorAddress

	msg := &types.MsgCreateOrchestratorValidatorAddress{
		ValidatorAddress:    validatorAddress,
		OrchestratorAddress: keepertest.OrchAddrs[0].String(),
	}
	_, err := srv.CreateOrchestratorValidatorAddress(wctx, msg)
	require.Error(t, err)
	require.Equal(t, "orchestrator already set: invalid request", err.Error())
}

func TestOrchestratorValidatorAddress_ValidatorAlreadySet(t *testing.T) {
	// flag true will also set orch validators for all 5 validators
	testKeepers := keepertest.SetFiveValidators(t, true)
	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)
	wctx := sdk.WrapSDKContext(testKeepers.Context)
	validatorAddress := testKeepers.StakingKeeper.GetValidators(testKeepers.Context, 10)[0].OperatorAddress

	msg := &types.MsgCreateOrchestratorValidatorAddress{
		ValidatorAddress:    validatorAddress,
		OrchestratorAddress: keepertest.AccAddrs[0].String(),
	}
	_, err := srv.CreateOrchestratorValidatorAddress(wctx, msg)
	require.Error(t, err)
	require.Equal(t, "validator already set: invalid request", err.Error())
}

func TestOrchestratorValidatorAddress_InvalidValidatorAddress(t *testing.T) {
	// flag true will also set orch validators for all 5 validators
	testKeepers := keepertest.SetFiveValidators(t, true)
	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)
	wctx := sdk.WrapSDKContext(testKeepers.Context)
	// validatorAddress := testKeepers.StakingKeeper.GetValidators(testKeepers.Context, 10)[0].OperatorAddress

	msg := &types.MsgCreateOrchestratorValidatorAddress{
		ValidatorAddress:    "12345",
		OrchestratorAddress: keepertest.AccAddrs[0].String(),
	}
	_, err := srv.CreateOrchestratorValidatorAddress(wctx, msg)
	require.Error(t, err)
	require.Equal(t, "validator address format error: not found", err.Error())
}

func TestOrchestratorValidatorAddress_ValidatorNotFound(t *testing.T) {
	// generate valid validator address for not existing validator
	AccPrivKey := secp256k1.GenPrivKey()
	AccPubKey := AccPrivKey.PubKey()
	ValAddr := sdk.ValAddress(AccPubKey.Address())

	// flag true will also set orch validators for all 5 validators
	testKeepers := keepertest.SetFiveValidators(t, true)
	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)
	wctx := sdk.WrapSDKContext(testKeepers.Context)
	// validatorAddress := testKeepers.StakingKeeper.GetValidators(testKeepers.Context, 10)[0].OperatorAddress

	msg := &types.MsgCreateOrchestratorValidatorAddress{
		ValidatorAddress:    ValAddr.String(),
		OrchestratorAddress: keepertest.AccAddrs[0].String(),
	}
	_, err := srv.CreateOrchestratorValidatorAddress(wctx, msg)
	require.Error(t, err)
	require.Equal(t, "validator not found or not in active set: not found", err.Error())
}
