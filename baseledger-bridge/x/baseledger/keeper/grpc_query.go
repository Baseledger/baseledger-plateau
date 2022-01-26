package keeper

import (
	"github.com/Baseledger/baseledger-bridge/x/baseledger/types"
)

var _ types.QueryServer = Keeper{}
