package types_test

import (
	"testing"

	"github.com/Baseledger/baseledger/testutil/sample"
	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	mockDefaultGenesis := types.DefaultGenesis()
	mockDefaultGenesis.Params.BaseledgerFaucetAddress = sample.AccAddress()
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: mockDefaultGenesis,
			valid:    true,
		},
		{
			desc: "duplicated orchestratorValidatorAddress",
			genState: &types.GenesisState{
				OrchestratorValidatorAddressList: []types.OrchestratorValidatorAddress{
					{
						OrchestratorAddress: "0",
					},
					{
						OrchestratorAddress: "0",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
