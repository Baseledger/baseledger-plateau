package keeper_test

import (
	"testing"
	"time"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	bridge "github.com/Baseledger/baseledger/x/bridge"
	"github.com/Baseledger/baseledger/x/bridge/keeper"
	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestValidatorPowerChangedClaim_Success(t *testing.T) {
	var (
		baseledgerTokenContract = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
		ethereumSender          = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		myBlockTime             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	testKeepers := keepertest.SetFiveValidators(t, true)

	validatorReceiver := keepertest.ValAddrs[0]
	ctx := testKeepers.Context

	validator, _ := testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	startAmount, _ := sdk.NewIntFromString("10000000")
	require.Equal(t, validator.Tokens, startAmount)

	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)
	// all validators, nonce 1
	newAmount, _ := sdk.NewIntFromString("10000005")
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgValidatorPowerChangedClaim{
			Creator:                            orchAddress.String(),
			EventNonce:                         uint64(1),
			TokenContract:                      baseledgerTokenContract,
			Amount:                             newAmount,
			BaseledgerReceiverValidatorAddress: validatorReceiver.String(),
			RevenueAddress:                     ethereumSender,
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.ValidatorPowerChangedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.NoError(t, err)

		hash, err := claim.ClaimHash()
		require.NoError(t, err)
		attestation := testKeepers.BridgeKeeper.GetAttestation(ctx, uint64(1), hash)
		require.NotNil(t, attestation)

		// Test to reject duplicate deposit
		ctx = ctx.WithBlockTime(myBlockTime)
		_, err = srv.ValidatorPowerChangedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.Error(t, err)
	}

	// balance increased correctly
	validator, _ = testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	require.Equal(t, validator.Tokens, newAmount)

	// all validators, nonce 3 (skipped one)
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgValidatorPowerChangedClaim{
			Creator:                            orchAddress.String(),
			EventNonce:                         uint64(3),
			TokenContract:                      baseledgerTokenContract,
			Amount:                             newAmount,
			BaseledgerReceiverValidatorAddress: validatorReceiver.String(),
			RevenueAddress:                     ethereumSender,
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.ValidatorPowerChangedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.Error(t, err)
	}

	// balance did not change after skipped nonce
	validator, _ = testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	require.Equal(t, validator.Tokens, newAmount)

	// all validators, correct nonce 2, this time decrease staking
	newAmount, _ = sdk.NewIntFromString("10000003")
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgValidatorPowerChangedClaim{
			Creator:                            orchAddress.String(),
			EventNonce:                         uint64(2),
			TokenContract:                      baseledgerTokenContract,
			Amount:                             newAmount,
			BaseledgerReceiverValidatorAddress: validatorReceiver.String(),
			RevenueAddress:                     ethereumSender,
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.ValidatorPowerChangedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.NoError(t, err)

		hash, err := claim.ClaimHash()
		require.NoError(t, err)
		attestation := testKeepers.BridgeKeeper.GetAttestation(ctx, uint64(2), hash)
		require.NotNil(t, attestation)
	}

	// balance decreased correctly
	validator, _ = testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	require.Equal(t, validator.Tokens, newAmount)

	// making sure 2stake are unbonding
	undelegations := testKeepers.StakingKeeper.GetUnbondingDelegations(ctx, keepertest.FaucetAccount, 10)
	require.Equal(t, 1, len(undelegations))
	require.Equal(t, 1, len(undelegations[0].Entries))
	require.Equal(t, keepertest.FaucetAccount.String(), undelegations[0].DelegatorAddress)
	require.Equal(t, validator.OperatorAddress, undelegations[0].ValidatorAddress)
}

func TestValidatorPowerChangedClaim_NonRegisteredOrchestratorValidator(t *testing.T) {
	var (
		baseledgerTokenContract = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
		ethereumSender          = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		myBlockTime             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	// DON'T register validators orchestrators
	testKeepers := keepertest.SetFiveValidators(t, false)

	validatorReceiver := keepertest.ValAddrs[0]
	ctx := testKeepers.Context

	validator, _ := testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	startAmount, _ := sdk.NewIntFromString("10000000")
	require.Equal(t, validator.Tokens, startAmount)

	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)

	// all validators, nonce 1
	newAmount, _ := sdk.NewIntFromString("10000005")
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgValidatorPowerChangedClaim{
			Creator:                            orchAddress.String(),
			EventNonce:                         uint64(1),
			TokenContract:                      baseledgerTokenContract,
			Amount:                             newAmount,
			BaseledgerReceiverValidatorAddress: validatorReceiver.String(),
			RevenueAddress:                     ethereumSender,
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.ValidatorPowerChangedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.Error(t, err)

		hash, err := claim.ClaimHash()
		require.NoError(t, err)
		attestation := testKeepers.BridgeKeeper.GetAttestation(ctx, uint64(1), hash)
		require.Nil(t, attestation)
	}

	// balance not changed
	validator, _ = testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	require.Equal(t, validator.Tokens, startAmount)
}

func TestValidatorPowerChangedClaim_NonExistingValidatorReceiver(t *testing.T) {
	var (
		baseledgerTokenContract = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
		ethereumSender          = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		myBlockTime             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	// register validators orchestrators
	testKeepers := keepertest.SetFiveValidators(t, true)

	// generate valid validator address for not existing validator
	AccPrivKey := secp256k1.GenPrivKey()
	AccPubKey := AccPrivKey.PubKey()
	ValAddr := sdk.ValAddress(AccPubKey.Address())

	validatorReceiver := ValAddr
	ctx := testKeepers.Context

	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)

	// all validators, nonce 1
	newAmount, _ := sdk.NewIntFromString("10000005")
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgValidatorPowerChangedClaim{
			Creator:                            orchAddress.String(),
			EventNonce:                         uint64(1),
			TokenContract:                      baseledgerTokenContract,
			Amount:                             newAmount,
			BaseledgerReceiverValidatorAddress: validatorReceiver.String(),
			RevenueAddress:                     ethereumSender,
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		srv.ValidatorPowerChangedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
	}

	// making sure there are no undelegations
	undelegations := testKeepers.StakingKeeper.GetUnbondingDelegations(ctx, keepertest.FaucetAccount, 10)
	require.Equal(t, 0, len(undelegations))
}

func TestValidatorPowerChangedClaim_NonBondedValidatorReceiver(t *testing.T) {
	var (
		baseledgerTokenContract = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
		ethereumSender          = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		myBlockTime             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	// register validators orchestrators
	testKeepers := keepertest.SetFiveValidators(t, true)

	validatorReceiver := keepertest.ValAddrs[0]
	ctx := testKeepers.Context

	validator, _ := testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	testKeepers.StakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.UpdateStatus(stakingtypes.Unbonding)
	validator.UnbondingTime = time.Now()
	testKeepers.StakingKeeper.BlockValidatorUpdates(ctx)
	bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)

	validator, _ = testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)

	startAmount, _ := sdk.NewIntFromString("10000000")
	require.Equal(t, validator.Tokens, startAmount)

	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)

	// all validators, nonce 1
	newAmount, _ := sdk.NewIntFromString("10000005")
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgValidatorPowerChangedClaim{
			Creator:                            orchAddress.String(),
			EventNonce:                         uint64(1),
			TokenContract:                      baseledgerTokenContract,
			Amount:                             newAmount,
			BaseledgerReceiverValidatorAddress: validatorReceiver.String(),
			RevenueAddress:                     ethereumSender,
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		srv.ValidatorPowerChangedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
	}

	// balance not changed
	validator, _ = testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	require.Equal(t, validator.Tokens, startAmount)
}

func TestValidatorPowerChangedClaim_JailedValidatorReceiver(t *testing.T) {
	var (
		baseledgerTokenContract = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
		ethereumSender          = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		myBlockTime             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	// register validators orchestrators
	testKeepers := keepertest.SetFiveValidators(t, true)

	validatorReceiver := keepertest.ValAddrs[0]
	ctx := testKeepers.Context

	validator, _ := testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	validatorConsAddress, _ := validator.GetConsAddr()
	testKeepers.StakingKeeper.Jail(ctx, validatorConsAddress)
	testKeepers.StakingKeeper.BlockValidatorUpdates(ctx)
	startAmount, _ := sdk.NewIntFromString("10000000")
	require.Equal(t, validator.Tokens, startAmount)

	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)

	// all validators, nonce 1
	newAmount, _ := sdk.NewIntFromString("10000005")
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgValidatorPowerChangedClaim{
			Creator:                            orchAddress.String(),
			EventNonce:                         uint64(1),
			TokenContract:                      baseledgerTokenContract,
			Amount:                             newAmount,
			BaseledgerReceiverValidatorAddress: validatorReceiver.String(),
			RevenueAddress:                     ethereumSender,
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		srv.ValidatorPowerChangedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
	}

	// balance not changed
	validator, _ = testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	require.Equal(t, validator.Tokens, startAmount)
}
