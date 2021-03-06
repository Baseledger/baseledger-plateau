package types_test

import (
	"testing"

	"github.com/Baseledger/baseledger/x/proof/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				BaseledgerTransactionList: []types.BaseledgerTransaction{
					{
						Id: 0,
					},
					{
						Id: 1,
					},
				},
				BaseledgerTransactionCount: 2,
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated baseledgerTransaction",
			genState: &types.GenesisState{
				BaseledgerTransactionList: []types.BaseledgerTransaction{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid baseledgerTransaction count",
			genState: &types.GenesisState{
				BaseledgerTransactionList: []types.BaseledgerTransaction{
					{
						Id: 1,
					},
				},
				BaseledgerTransactionCount: 0,
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
