package types

import (
	"testing"

	"github.com/Baseledger/baseledger/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateOrchestratorValidatorAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name          string
		msg           MsgCreateOrchestratorValidatorAddress
		isErrExpected bool
	}{
		{
			name: "invalid val address",
			msg: MsgCreateOrchestratorValidatorAddress{
				ValidatorAddress: "invalid_address",
			},
			isErrExpected: true,
		}, {
			name: "invalid orch address",
			msg: MsgCreateOrchestratorValidatorAddress{
				ValidatorAddress:    sample.ValAddress(),
				OrchestratorAddress: "invalid_address",
			},
			isErrExpected: true,
		}, {
			name: "valid",
			msg: MsgCreateOrchestratorValidatorAddress{
				ValidatorAddress:    sample.ValAddress(),
				OrchestratorAddress: sample.AccAddress(),
			},
			isErrExpected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.isErrExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
