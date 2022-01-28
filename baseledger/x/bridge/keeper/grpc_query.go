package keeper

import (
	"github.com/Baseledger/baseledger/x/bridge/types"
)

var _ types.QueryServer = Keeper{}
