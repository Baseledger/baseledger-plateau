package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetOrchestratorAddressKey returns the following key format
// prefix
// [0xe8][gravity1ahx7f8wyertuus9r20284ej0asrs085ceqtfnm]
func GetOrchestratorAddressKey(orc sdk.AccAddress) string {
	if err := sdk.VerifyAddressFormat(orc); err != nil {
		panic(sdkerrors.Wrap(err, "invalid orchestrator address"))
	}
	return KeyOrchestratorAddress + string(orc.Bytes())
}
