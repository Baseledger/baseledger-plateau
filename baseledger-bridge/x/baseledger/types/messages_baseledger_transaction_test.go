package types

import (
	"testing"

	"github.com/Baseledger/baseledger-bridge/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateBaseledgerTransaction_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateBaseledgerTransaction
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateBaseledgerTransaction{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateBaseledgerTransaction{
				Creator: sample.AccAddress(),
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

func TestMsgUpdateBaseledgerTransaction_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateBaseledgerTransaction
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateBaseledgerTransaction{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateBaseledgerTransaction{
				Creator: sample.AccAddress(),
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

func TestMsgDeleteBaseledgerTransaction_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteBaseledgerTransaction
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteBaseledgerTransaction{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteBaseledgerTransaction{
				Creator: sample.AccAddress(),
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
