package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/x/proof/keeper"
	"github.com/Baseledger/baseledger/x/proof/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, keepertest.TestProofKeepers, context.Context) {
	testKeepers := keepertest.BaseledgerKeeper(t)
	k := testKeepers.ProofKeeper
	ctx := testKeepers.Context
	return keeper.NewMsgServerImpl(*k), testKeepers, sdk.WrapSDKContext(ctx)
}
