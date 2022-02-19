package types

import (
	"testing"

	"github.com/Baseledger/baseledger/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgValidatorPowerChangedClaim_ValidateBasic(t *testing.T) {
	tests := []struct {
		name          string
		msg           MsgValidatorPowerChangedClaim
		isErrExpected bool
	}{
		{
			name: "invalid revenue address",
			msg: MsgValidatorPowerChangedClaim{
				RevenueAddress: "invalid_eth_sender",
			},
			isErrExpected: true,
		},
		{
			name: "invalid token contract",
			msg: MsgValidatorPowerChangedClaim{
				RevenueAddress: "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
				TokenContract:  "invalid_token_contract",
			},
			isErrExpected: true,
		},
		{
			name: "invalid orch address",
			msg: MsgValidatorPowerChangedClaim{
				RevenueAddress: "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
				TokenContract:  "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
				Creator:        "invalid_creator",
			},
			isErrExpected: true,
		},
		{
			name: "invalid receiver address",
			msg: MsgValidatorPowerChangedClaim{
				RevenueAddress:                     "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
				TokenContract:                      "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
				Creator:                            sample.AccAddress(),
				BaseledgerReceiverValidatorAddress: "invalid_receiver",
			},
			isErrExpected: true,
		},
		{
			name: "invalid event nonce",
			msg: MsgValidatorPowerChangedClaim{
				RevenueAddress:                     "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
				TokenContract:                      "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
				Creator:                            sample.AccAddress(),
				BaseledgerReceiverValidatorAddress: sample.ValAddress(),
				EventNonce:                         0,
			},
			isErrExpected: true,
		},
		{
			name: "valid",
			msg: MsgValidatorPowerChangedClaim{
				RevenueAddress:                     "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
				TokenContract:                      "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
				Creator:                            sample.AccAddress(),
				BaseledgerReceiverValidatorAddress: sample.ValAddress(),
				EventNonce:                         1,
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
