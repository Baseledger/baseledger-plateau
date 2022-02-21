package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/testutil/sample"
	"github.com/Baseledger/baseledger/x/proof/types"
)

func TestBaseledgerTransactionMsgServerCreate_Success(t *testing.T) {
	srv, testKeepers, ctx := setupMsgServer(t)
	creator := keepertest.SampleAccountWithFunds
	startCreatorBalance := testKeepers.BankKeeper.GetBalance(testKeepers.Context, creator, "work")
	require.Equal(t, "100work", startCreatorBalance.String())

	startFaucetBalance := testKeepers.BankKeeper.GetBalance(testKeepers.Context, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", startFaucetBalance.String())
	for i := 0; i < 5; i++ {
		resp, err := srv.CreateBaseledgerTransaction(ctx, &types.MsgCreateBaseledgerTransaction{Creator: creator.String(), Payload: "A"})
		require.NoError(t, err)
		require.Equal(t, i, int(resp.Id))
		substractedBalance := 100 - i - 1
		currentCreatorBalance := testKeepers.BankKeeper.GetBalance(testKeepers.Context, creator, "work")
		require.Equal(t, fmt.Sprintf("%dwork", substractedBalance), currentCreatorBalance.String())

		faucetBalance := 50000000 + i + 1
		currentFaucetBalance := testKeepers.BankKeeper.GetBalance(testKeepers.Context, keepertest.FaucetAccount, "work")
		require.Equal(t, fmt.Sprintf("%dwork", faucetBalance), currentFaucetBalance.String())
	}
}

func TestBaseledgerTransactionMsgServerCreate_DifferentPayloads(t *testing.T) {
	srv, testKeepers, ctx := setupMsgServer(t)
	creator := keepertest.SampleAccountWithFunds
	startCreatorBalance := testKeepers.BankKeeper.GetBalance(testKeepers.Context, creator, "work")
	require.Equal(t, "100work", startCreatorBalance.String())

	// lenghts 1, 129, 257, 513, 1025
	payloads := [5]string{
		"A",
		"HktKSyzH7ahKSSHgDrv0NlE7KaGL5u85mSi8sdds4kiLJkxUzQ9NMWsmtSIwkP4hk7OIO04qG8CyNBQRgoyIGe53ibSTcdbIER4hQgTPVAiZsMeIIk5GbYsR6HOjsfw6W",
		"emciTfvp6F2vFTH43IrcgKRzmk0NqFU7wpcuyQuKmJY1pnxwCne0HG6nVnvUwaZmL5BBvcG39rmz1PeGfZ2M8Nnlfxabk2DVzToTwRj5YTkogUD92Y43ARCAPb3ZNpcbSDxoYLnHN3NRW4Z8rjxRwcdwMtO5WYTujpN18BD0v0VSzRbZuNNlrjlsUSJlyQ9RYj7nh0ktpp150V4BzyCUSSW1Q2jslVrmn3gxd6vqcTeJ0GUXg6dO81MLnZbXTOOSM",
		"dhIcmq69jpPKsKC6yYoDigFPT6T85nsl6oQr8NT7I7k2504OllfTws68jFq81nuu8vD2huPxyaZpKkEXpoFYlnb65UElk9H0JkNWEYIhe8oC0abb2Hkd9ZCLybpSpf5wTSpD4yeerTjx10k2tHkJ512Mbo5zY6knZaIIFn3I4EKakd5b7Mm17mPAY1jEyFucsoov9QEcwS5LmenlWgvQVG53QFDqmEfUJmFGBBwu6p7Dp1DEtQzVB3tgGIDTXpObRjtaBKo3dmXAXRgLFDgkU83mAJJhcU1worIpF0CtboI1H6WtN6CwI909DxowEe5eeOhVJhweMJ6IaAJUl3ptZWSkHgpH5FWzii3OvpjlCxN2eqRpLEzVW82mvUlxOnfp5bTdklvc2TJD9uIXP4jqeymlRji6NYacWGHXDaieozcuaDpmqSmCrADNm1fUEa7DCIrc3D72nhWt0VwWTlVmluo1XcjwyvrthIPYXiHdH2nNGtaQ7h3kQP0hXy6aDA3VI",
		"ufZXAdedbmIdeDIKIB5buy0vYDfA1vRvgJfMr1LrW6fEsHnMJmyFd9fNwFLJFAIv4m1TyjZdP6q6qYBxeZlUVG6yCYTnvDixiwVsJWqe20aBOnr2UGteaASyvXaqD0helWIJJioL9gv53hv27veq5pCMdKtYX5yh6HdlaaBFXIO85lfd0OvxUAWqcYB1cYPXJRKAZZBxzi6uZp6D3YdATNJ9p1GVZFC4dngRnrxLy4IQ8P9sUjIXBXrK1fWOLApPGKPxDQZ8FNssLHBQFmcn5lRezocdJliAtgz3Hbw1KjxS7AIACY8SSLIMrHTUSz9wdemyQjJl5UyGGeKljbdR5ex3ifLAI1f4FtUEFgZtHsnbuEszoSqqdFVraIYb8cxuFxJinU72tYx5pIpRRG8SCHZtaIOBpiNAMssNyiCFQPsEnDDVeT1lNoDyPX8EHRkR3HDHxx3Qs1w9HTP1Ciu2ufrwYukLNeEQF33phRIj4wNvwAgZfBHQ77qAljiLmN4W6Ufd2COTWUaWduoQPXTgp1owSP56WXaMGbaWIyjYWgdYHeob7kxwg3lwiryupW0gu984FjBZQ30QaGOkuPQP1ilvnVejWUg464NTc22Mg9ofB1UquCevunY74aiKqTJhzkAc2Vv7fdVhN22z7ZLAQwrGDPHa04olC6zQNI6aAu81GrumPsis2iqXmnoZNwssQylGIfS48H1ierbiFHqi7mMuERMrf6S0dDDW8Z9O0eEO74YYmQ1q456rc8rqPH6KGhepjXkMiP41cjaMHs4M8Jo61oHw5VRL3QXG3jHvnfPDxUG55updMiqbjbeZTSA952YTvdUkRqQ7syrmR5vQrgD1GKq0vcLMf3FNlXuuXa3gE3rT9qPQfEWC6hCpNw5F0Wr3ypHJt5XYHmG717H0elVQzwFP1gCs7715tsoJgItQX92GRFzj2Ab7sEMe0ybkmERgBF5Skuo3hZi9bzHcxzAizh6lcaBY94sXvkJAuVHCq6HWVMdsNjCc8smLzdEbr",
	}

	// cost is 1, 2, 6, 16, but here we store accumulated costs because we are testing all steps
	payloadsCost := [4]uint{
		1,
		3,
		9,
		25,
	}

	startFaucetBalance := testKeepers.BankKeeper.GetBalance(testKeepers.Context, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", startFaucetBalance.String())
	for i := 0; i < 4; i++ {
		resp, err := srv.CreateBaseledgerTransaction(ctx, &types.MsgCreateBaseledgerTransaction{Creator: creator.String(), Payload: payloads[i]})
		require.NoError(t, err)
		require.Equal(t, i, int(resp.Id))
		substractedBalance := 100 - payloadsCost[i]
		currentCreatorBalance := testKeepers.BankKeeper.GetBalance(testKeepers.Context, creator, "work")
		require.Equal(t, fmt.Sprintf("%dwork", substractedBalance), currentCreatorBalance.String())

		faucetBalance := 50000000 + payloadsCost[i]
		currentFaucetBalance := testKeepers.BankKeeper.GetBalance(testKeepers.Context, keepertest.FaucetAccount, "work")
		require.Equal(t, fmt.Sprintf("%dwork", faucetBalance), currentFaucetBalance.String())
	}

	// for payloads[4] it should fail because length is not supported
	_, err := srv.CreateBaseledgerTransaction(ctx, &types.MsgCreateBaseledgerTransaction{Creator: creator.String(), Payload: payloads[4]})
	require.Error(t, err)
}

func TestBaseledgerTransactionMsgServerCreate_NoFundsForCreator(t *testing.T) {
	srv, testKeepers, ctx := setupMsgServer(t)
	creator := sample.AccAddress()

	startFaucetBalance := testKeepers.BankKeeper.GetBalance(testKeepers.Context, keepertest.FaucetAccount, "work")
	require.Equal(t, "50000000work", startFaucetBalance.String())
	for i := 0; i < 5; i++ {
		_, err := srv.CreateBaseledgerTransaction(ctx, &types.MsgCreateBaseledgerTransaction{Creator: creator})
		require.Error(t, err)
		currentFaucetBalance := testKeepers.BankKeeper.GetBalance(testKeepers.Context, keepertest.FaucetAccount, "work")
		require.Equal(t, startFaucetBalance.String(), currentFaucetBalance.String())
	}
}
