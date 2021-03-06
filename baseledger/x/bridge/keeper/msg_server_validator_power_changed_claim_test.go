package keeper_test

import (
	"testing"
	"time"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	bridge "github.com/Baseledger/baseledger/x/bridge"
	"github.com/Baseledger/baseledger/x/bridge/keeper"
	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ccrypto "github.com/cosmos/cosmos-sdk/crypto/types"
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
	// all validators, nonce 1, 50k with 8 decimals
	newAmount, _ := sdk.NewIntFromString("5000000000000")
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
	// 50k with 6 decimals
	newAmountExp, _ := sdk.NewIntFromString("50000000000")
	require.Equal(t, validator.Tokens, newAmountExp)

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
	require.Equal(t, validator.Tokens, newAmountExp)

	// all validators, correct nonce 2, this time decrease staking, 30k with 8 decimals
	newAmount, _ = sdk.NewIntFromString("3000000000000")
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
	// 30k with 6 decimals
	newAmountExp, _ = sdk.NewIntFromString("30000000000")
	require.Equal(t, validator.Tokens, newAmountExp)

	// making sure 2stake are unbonding
	undelegations := testKeepers.StakingKeeper.GetUnbondingDelegations(ctx, keepertest.FaucetAccount, 10)
	require.Equal(t, 1, len(undelegations))
	require.Equal(t, 1, len(undelegations[0].Entries))
	require.Equal(t, keepertest.FaucetAccount.String(), undelegations[0].DelegatorAddress)
	require.Equal(t, validator.OperatorAddress, undelegations[0].ValidatorAddress)
}

func TestValidatorPowerChangedClaim_NotObserved(t *testing.T) {
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
	newAmount, _ := sdk.NewIntFromString("5000000000000")
	orchSet := []sdk.AccAddress{keepertest.OrchAddrs[0], keepertest.OrchAddrs[1]}

	for _, orchAddress := range orchSet {
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
		require.False(t, attestation.Observed)
	}

	// balance did not change
	validator, _ = testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	require.Equal(t, validator.Tokens, startAmount)
}

func TestValidatorPowerChangedClaim_SpreadVotes(t *testing.T) {
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
	newAmount, _ := sdk.NewIntFromString("5000000000000")
	orchSet := []sdk.AccAddress{keepertest.OrchAddrs[0], keepertest.OrchAddrs[1]}

	for _, orchAddress := range orchSet {
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
		require.False(t, attestation.Observed)
	}

	// balance did not change
	validator, _ = testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	require.Equal(t, validator.Tokens, startAmount)

	secondOrchSet := []sdk.AccAddress{keepertest.OrchAddrs[2], keepertest.OrchAddrs[3], keepertest.OrchAddrs[4]}
	for _, orchAddress := range secondOrchSet {
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
	}

	// balance changed
	validator, _ = testKeepers.StakingKeeper.GetValidator(ctx, validatorReceiver)
	newAmountExp, _ := sdk.NewIntFromString("50000000000")
	require.Equal(t, validator.Tokens, newAmountExp)
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
	newAmount, _ := sdk.NewIntFromString("5000000000000")
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

func TestValidatorPowerChangedClaim_NonExistingOrchestratorSet(t *testing.T) {
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
	startAmount, _ := sdk.NewIntFromString("10000000")
	require.Equal(t, validator.Tokens, startAmount)

	// Orchestrator private keys
	OrchPrivKeys := []ccrypto.PrivKey{
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
	}

	// AccPubKeys holds the pub keys for the account keys
	OrchPubKeys := []ccrypto.PubKey{
		OrchPrivKeys[0].PubKey(),
		OrchPrivKeys[1].PubKey(),
		OrchPrivKeys[2].PubKey(),
		OrchPrivKeys[3].PubKey(),
		OrchPrivKeys[4].PubKey(),
	}
	// AccAddrs holds the sdk.AccAddresses
	NewOrchAddrs := []sdk.AccAddress{
		sdk.AccAddress(OrchPubKeys[0].Address()),
		sdk.AccAddress(OrchPubKeys[1].Address()),
		sdk.AccAddress(OrchPubKeys[2].Address()),
		sdk.AccAddress(OrchPubKeys[3].Address()),
		sdk.AccAddress(OrchPubKeys[4].Address()),
	}

	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)

	// all validators, nonce 1
	newAmount, _ := sdk.NewIntFromString("5000000000000")
	for _, orchAddress := range NewOrchAddrs {
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
	newAmount, _ := sdk.NewIntFromString("5000000000000")
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
	newAmount, _ := sdk.NewIntFromString("5000000000000")
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
	newAmount, _ := sdk.NewIntFromString("5000000000000")
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
