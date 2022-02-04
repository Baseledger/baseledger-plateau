package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/Baseledger/baseledger/testutil/sample"
)

func TestMsgCreateOrchestratorValidatorAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateOrchestratorValidatorAddress
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateOrchestratorValidatorAddress{
				ValidatorAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateOrchestratorValidatorAddress{
				ValidatorAddress: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateOrchestratorValidatorAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateOrchestratorValidatorAddress
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateOrchestratorValidatorAddress{
				ValidatorAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateOrchestratorValidatorAddress{
				ValidatorAddress: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgDeleteOrchestratorValidatorAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteOrchestratorValidatorAddress
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteOrchestratorValidatorAddress{
				ValidatorAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteOrchestratorValidatorAddress{
				ValidatorAddress: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
