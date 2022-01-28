package keeper

import (
	"github.com/Baseledger/baseledger/x/proof/types"
)

var _ types.QueryServer = Keeper{}
