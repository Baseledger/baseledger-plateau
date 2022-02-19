package keeper_test

import (
	"testing"
	"time"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	bridge "github.com/Baseledger/baseledger/x/bridge"
	"github.com/Baseledger/baseledger/x/bridge/keeper"
	"github.com/Baseledger/baseledger/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestUbtDepositedClaim(t *testing.T) {
	var (
		baseledgerTokenContract = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
		ethereumSender          = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		myBlockTime             = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	testKeepers := keepertest.SetFiveValidators(t, true)

	cosmosReceiver := keepertest.AccAddrs[0]
	ctx := testKeepers.Context

	balance := testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "0work", balance.String())

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
	balance = testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "1work", balance.String())

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
	balance = testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "1work", balance.String())

	// all validators, correct nonce 2
	for _, orchAddress := range keepertest.OrchAddrs {
		claim := types.MsgUbtDepositedClaim{
			EventNonce:                       uint64(2),
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
		attestation := testKeepers.BridgeKeeper.GetAttestation(ctx, uint64(2), hash)
		require.NotNil(t, attestation)
	}

	// balance increased again correctly
	balance = testKeepers.BankKeeper.GetBalance(ctx, cosmosReceiver, "work")
	require.Equal(t, "2work", balance.String())
}
