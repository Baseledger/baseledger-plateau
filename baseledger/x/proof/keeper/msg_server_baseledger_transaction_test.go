package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Baseledger/baseledger/x/proof/types"
)

func TestBaseledgerTransactionMsgServerCreate(t *testing.T) {
	srv, ctx := setupMsgServer(t)
	creator := "A"
	for i := 0; i < 5; i++ {
		resp, err := srv.CreateBaseledgerTransaction(ctx, &types.MsgCreateBaseledgerTransaction{Creator: creator})
		require.NoError(t, err)
		require.Equal(t, i, int(resp.Id))
	}
}
