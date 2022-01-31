package common

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const StakeTokenDenom = "stake"
const WorkTokenDenom = "work"
const WorkTokenFee = "1work"
const UbtFaucetAddress = "baseledger1xgs5tamqre7rkz5q7d5fegjsdwufxxvt36w0a0"

/* Current formula, subject to change, decide either through code update or somehow using params
   If we decide to update this using params, also modify rest endpoint
- 128: 1 transaction cost unit
- 256: 2 transaction cost units
- 512: 6 transaction cost units
- 1024: 16 transaction cost units
*/
func CalcWorkTokenFeeBasedOnPayloadSize(payload string) (sdk.Coins, error) {
	if len(payload) > 1024 {
		return nil, errors.New("Current payload size not supported, max is 1024")
	}

	workTokenAmount := 1

	if len(payload) > 512 {
		workTokenAmount = workTokenAmount * 16
		return sdk.ParseCoinsNormalized(fmt.Sprintf("%dwork", workTokenAmount))
	}

	if len(payload) > 256 {
		workTokenAmount = workTokenAmount * 6
		return sdk.ParseCoinsNormalized(fmt.Sprintf("%dwork", workTokenAmount))
	}

	if len(payload) > 128 {
		workTokenAmount = workTokenAmount * 2
		return sdk.ParseCoinsNormalized(fmt.Sprintf("%dwork", workTokenAmount))
	}

	return sdk.ParseCoinsNormalized(fmt.Sprintf("%dwork", workTokenAmount))
}
