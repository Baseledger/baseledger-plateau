package keeper_test

import (
	"testing"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetOrchestratorValidator_GetSuccessful(t *testing.T) {
	testKeepers := keepertest.SetFiveValidators(t, true)

	val := testKeepers.BridgeKeeper.GetOrchestratorValidator(testKeepers.Context, keepertest.OrchAddrs[0].String())
	require.NotNil(t, val)
}

func TestGetOrchestratorValidator_GetFailWrongOrchestratorAddressFormat(t *testing.T) {
	testKeepers := keepertest.SetFiveValidators(t, true)

	val := testKeepers.BridgeKeeper.GetOrchestratorValidator(testKeepers.Context, "12345")
	require.Nil(t, val)
}

func TestGetOrchestratorValidator_GetFailWrongOrchestratorAddress(t *testing.T) {
	testKeepers := keepertest.SetFiveValidators(t, false)

	val := testKeepers.BridgeKeeper.GetOrchestratorValidator(testKeepers.Context, keepertest.OrchAddrs[0].String())
	require.Nil(t, val)
}

func TestGetOrchestratorValidator_GetFailWrongValidatorAddressFormat(t *testing.T) {
	testKeepers := keepertest.SetFiveValidators(t, false)
	// this will never happen because of validate basic, but just to cover the case
	testKeepers.BridgeKeeper.SetOrchestratorValidatorAddress(testKeepers.Context, types.OrchestratorValidatorAddress{
		OrchestratorAddress: keepertest.OrchAddrs[0].String(),
		ValidatorAddress:    "12345",
	})

	val := testKeepers.BridgeKeeper.GetOrchestratorValidator(testKeepers.Context, keepertest.OrchAddrs[0].String())
	require.Nil(t, val)
}

func TestGetOrchestratorValidator_GetFailWrongValidatorAddress(t *testing.T) {
	testKeepers := keepertest.SetFiveValidators(t, false)
	// generate valid validator address for not existing validator
	AccPrivKey := secp256k1.GenPrivKey()
	AccPubKey := AccPrivKey.PubKey()
	ValAddr := sdk.ValAddress(AccPubKey.Address())
	// this will never happen because of validate basic, but just to cover the case
	testKeepers.BridgeKeeper.SetOrchestratorValidatorAddress(testKeepers.Context, types.OrchestratorValidatorAddress{
		OrchestratorAddress: keepertest.OrchAddrs[0].String(),
		ValidatorAddress:    ValAddr.String(),
	})

	val := testKeepers.BridgeKeeper.GetOrchestratorValidator(testKeepers.Context, keepertest.OrchAddrs[0].String())
	require.Nil(t, val)
}
