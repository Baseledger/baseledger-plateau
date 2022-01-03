package keeper

import (
	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
)

var _ types.QueryServer = Keeper{}
