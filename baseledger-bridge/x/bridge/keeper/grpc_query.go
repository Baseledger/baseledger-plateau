package keeper

import (
	"github.com/Baseledger/baseledger-bridge/x/bridge/types"
)

var _ types.QueryServer = Keeper{}
