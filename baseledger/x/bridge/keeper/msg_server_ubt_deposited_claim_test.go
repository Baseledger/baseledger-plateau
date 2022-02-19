package keeper_test

import (
	"math/big"
	"testing"
	"time"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	bridge "github.com/Baseledger/baseledger/x/bridge"
	"github.com/Baseledger/baseledger/x/bridge/keeper"
	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ccrypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestUbtDepositedClaim_Success(t *testing.T) {
	var (
		baseledgerTokenContract = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
		ethereumSender          = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		myBlockTime             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	testKeepers := keepertest.SetFiveValidators(t, true)

	cosmosReceiver := keepertest.AccAddrs[0]

	ctx := testKeepers.Context

	cosmosReceiverBalance := testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "0work", cosmosReceiverBalance.String())
	faucetBalance := testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", faucetBalance.String())

	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)

	// all validators, nonce 1
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgUbtDepositedClaim{
			EventNonce:                       uint64(1),
			TokenContract:                    baseledgerTokenContract,
			Amount:                           sdk.NewIntFromUint64(1),
			EthereumSender:                   ethereumSender,
			UbtPrice:                         "1",
			Creator:                          orchAddress.String(),
			BaseledgerReceiverAccountAddress: cosmosReceiver.String(),
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.UbtDepositedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.NoError(t, err)

		hash, err := claim.ClaimHash()
		require.NoError(t, err)
		attestation := testKeepers.BridgeKeeper.GetAttestation(ctx, uint64(1), hash)
		require.NotNil(t, attestation)

		// Test to reject duplicate deposit
		ctx = ctx.WithBlockTime(myBlockTime)
		_, err = srv.UbtDepositedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.Error(t, err)
	}

	// balance increased correctly
	cosmosReceiverBalance = testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "1work", cosmosReceiverBalance.String())

	faucetBalance = testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "49999999work", faucetBalance.String())

	// all validators, nonce 3 (skipped one)
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgUbtDepositedClaim{
			EventNonce:                       uint64(3),
			TokenContract:                    baseledgerTokenContract,
			Amount:                           sdk.NewIntFromUint64(1),
			EthereumSender:                   ethereumSender,
			UbtPrice:                         "1",
			Creator:                          orchAddress.String(),
			BaseledgerReceiverAccountAddress: cosmosReceiver.String(),
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.UbtDepositedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.Error(t, err)
	}

	// balance did not change after skipped nonce
	cosmosReceiverBalance = testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "1work", cosmosReceiverBalance.String())

	faucetBalance = testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "49999999work", faucetBalance.String())

	// all validators, correct nonce 2, but without price
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgUbtDepositedClaim{
			EventNonce:                       uint64(2),
			TokenContract:                    baseledgerTokenContract,
			Amount:                           sdk.NewIntFromUint64(1),
			EthereumSender:                   ethereumSender,
			UbtPrice:                         "0",
			Creator:                          orchAddress.String(),
			BaseledgerReceiverAccountAddress: cosmosReceiver.String(),
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.UbtDepositedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.NoError(t, err)

		hash, err := claim.ClaimHash()
		require.NoError(t, err)
		attestation := testKeepers.BridgeKeeper.GetAttestation(ctx, uint64(2), hash)
		require.NotNil(t, attestation)
	}

	// balance increased again correctly
	cosmosReceiverBalance = testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "2work", cosmosReceiverBalance.String())

	faucetBalance = testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "49999998work", faucetBalance.String())

	lastAvgPrice := testKeepers.BridgeKeeper.GetLastAttestationAverageUbtPrice(ctx)
	require.Equal(t, big.NewInt(1000000000000000000), lastAvgPrice)
}

func TestUbtDepositedClaim_NonRegisteredOrchestratorValidator(t *testing.T) {
	var (
		baseledgerTokenContract = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
		ethereumSender          = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		myBlockTime             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	// DON'T register validators orchestrators
	testKeepers := keepertest.SetFiveValidators(t, false)

	cosmosReceiver := keepertest.AccAddrs[0]

	ctx := testKeepers.Context

	cosmosReceiverBalance := testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "0work", cosmosReceiverBalance.String())
	faucetBalance := testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", faucetBalance.String())

	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)

	// all validators, nonce 1
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgUbtDepositedClaim{
			EventNonce:                       uint64(1),
			TokenContract:                    baseledgerTokenContract,
			Amount:                           sdk.NewIntFromUint64(1),
			EthereumSender:                   ethereumSender,
			UbtPrice:                         "1",
			Creator:                          orchAddress.String(),
			BaseledgerReceiverAccountAddress: cosmosReceiver.String(),
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.UbtDepositedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.Error(t, err)

		hash, err := claim.ClaimHash()
		require.NoError(t, err)
		attestation := testKeepers.BridgeKeeper.GetAttestation(ctx, uint64(1), hash)
		require.Nil(t, attestation)
	}

	// balance not changed
	cosmosReceiverBalance = testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "0work", cosmosReceiverBalance.String())
	faucetBalance = testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", faucetBalance.String())
}

func TestUbtDepositedClaim_NonExistingOrchestratorSet(t *testing.T) {
	var (
		baseledgerTokenContract = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
		ethereumSender          = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		myBlockTime             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	// register validators orchestrators
	testKeepers := keepertest.SetFiveValidators(t, true)

	cosmosReceiver := keepertest.AccAddrs[0]

	ctx := testKeepers.Context

	cosmosReceiverBalance := testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "0work", cosmosReceiverBalance.String())
	faucetBalance := testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", faucetBalance.String())

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

	// new not registered orch addresses
	for _, orchAddress := range NewOrchAddrs {
		claim := types.MsgUbtDepositedClaim{
			EventNonce:                       uint64(1),
			TokenContract:                    baseledgerTokenContract,
			Amount:                           sdk.NewIntFromUint64(1),
			EthereumSender:                   ethereumSender,
			UbtPrice:                         "1",
			Creator:                          orchAddress.String(),
			BaseledgerReceiverAccountAddress: cosmosReceiver.String(),
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.UbtDepositedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.Error(t, err)

		hash, err := claim.ClaimHash()
		require.NoError(t, err)
		attestation := testKeepers.BridgeKeeper.GetAttestation(ctx, uint64(1), hash)
		require.Nil(t, attestation)
	}

	// balance not changed
	cosmosReceiverBalance = testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "0work", cosmosReceiverBalance.String())
	faucetBalance = testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", faucetBalance.String())
}

func TestUbtDepositedClaim_NotObserved(t *testing.T) {
	var (
		baseledgerTokenContract = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
		ethereumSender          = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		myBlockTime             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	testKeepers := keepertest.SetFiveValidators(t, true)

	cosmosReceiver := keepertest.AccAddrs[0]

	ctx := testKeepers.Context

	cosmosReceiverBalance := testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "0work", cosmosReceiverBalance.String())
	faucetBalance := testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", faucetBalance.String())

	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)

	orchSet := []sdk.AccAddress{keepertest.OrchAddrs[0], keepertest.OrchAddrs[1]}
	// only 2 validators sending claim, nonce 1
	for _, orchAddress := range orchSet {
		claim := types.MsgUbtDepositedClaim{
			EventNonce:                       uint64(1),
			TokenContract:                    baseledgerTokenContract,
			Amount:                           sdk.NewIntFromUint64(1),
			EthereumSender:                   ethereumSender,
			UbtPrice:                         "1",
			Creator:                          orchAddress.String(),
			BaseledgerReceiverAccountAddress: cosmosReceiver.String(),
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.UbtDepositedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.NoError(t, err)

		hash, err := claim.ClaimHash()
		require.NoError(t, err)
		attestation := testKeepers.BridgeKeeper.GetAttestation(ctx, uint64(1), hash)
		require.NotNil(t, attestation)
		require.False(t, attestation.Observed)
	}

	// balances did not change
	cosmosReceiverBalance = testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "0work", cosmosReceiverBalance.String())
	faucetBalance = testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", faucetBalance.String())
}

func TestUbtDepositedClaim_SpreadVotes(t *testing.T) {
	var (
		baseledgerTokenContract = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
		ethereumSender          = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		myBlockTime             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	testKeepers := keepertest.SetFiveValidators(t, true)

	cosmosReceiver := keepertest.AccAddrs[0]

	ctx := testKeepers.Context

	cosmosReceiverBalance := testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "0work", cosmosReceiverBalance.String())
	faucetBalance := testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", faucetBalance.String())

	srv := keeper.NewMsgServerImpl(*testKeepers.BridgeKeeper)

	orchSet := []sdk.AccAddress{keepertest.OrchAddrs[0], keepertest.OrchAddrs[1]}
	// only 2 validators sending claim, nonce 1
	for _, orchAddress := range orchSet {
		claim := types.MsgUbtDepositedClaim{
			EventNonce:                       uint64(1),
			TokenContract:                    baseledgerTokenContract,
			Amount:                           sdk.NewIntFromUint64(1),
			EthereumSender:                   ethereumSender,
			UbtPrice:                         "1",
			Creator:                          orchAddress.String(),
			BaseledgerReceiverAccountAddress: cosmosReceiver.String(),
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.UbtDepositedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.NoError(t, err)

		hash, err := claim.ClaimHash()
		require.NoError(t, err)
		attestation := testKeepers.BridgeKeeper.GetAttestation(ctx, uint64(1), hash)
		require.NotNil(t, attestation)
		require.False(t, attestation.Observed)
	}

	// balances did not change
	cosmosReceiverBalance = testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "0work", cosmosReceiverBalance.String())
	faucetBalance = testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", faucetBalance.String())

	secondOrchSet := []sdk.AccAddress{keepertest.OrchAddrs[2], keepertest.OrchAddrs[3], keepertest.OrchAddrs[4]}

	// only 2 validators sending claim, nonce 1
	for _, orchAddress := range secondOrchSet {
		claim := types.MsgUbtDepositedClaim{
			EventNonce:                       uint64(1),
			TokenContract:                    baseledgerTokenContract,
			Amount:                           sdk.NewIntFromUint64(1),
			EthereumSender:                   ethereumSender,
			UbtPrice:                         "1",
			Creator:                          orchAddress.String(),
			BaseledgerReceiverAccountAddress: cosmosReceiver.String(),
		}

		ctx = ctx.WithBlockTime(myBlockTime)
		_, err := srv.UbtDepositedClaim(sdk.WrapSDKContext(ctx), &claim)
		bridge.EndBlocker(ctx, *testKeepers.BridgeKeeper)
		require.NoError(t, err)
	}

	// balance changed
	cosmosReceiverBalance = testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "1work", cosmosReceiverBalance.String())
	faucetBalance = testKeepers.BankKeeper.GetBalance(ctx, keepertest.FaucetAccount, "work")
	require.Equal(t, "49999999work", faucetBalance.String())
}
