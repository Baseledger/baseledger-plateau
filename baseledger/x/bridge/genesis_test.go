package bridge_test

import (
	"testing"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/testutil/nullify"
	"github.com/Baseledger/baseledger/testutil/sample"
	bridge "github.com/Baseledger/baseledger/x/bridge"
	"github.com/Baseledger/baseledger/x/bridge/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.Params{
			WorktokenEurPrice:       "0.1",
			BaseledgerFaucetAddress: sample.AccAddress(),
		},

		OrchestratorValidatorAddressList: []types.OrchestratorValidatorAddress{
			{
				OrchestratorAddress: sample.AccAddress(),
				ValidatorAddress:    sample.ValAddress(),
			},
			{
				OrchestratorAddress: sample.AccAddress(),
				ValidatorAddress:    sample.ValAddress(),
			},
		},

		LastObservedNonce: 5,

		Attestations: []types.Attestation{
			{
				Observed: true,
				Votes:    []string{sample.ValAddress()},
				Height:   uint64(2),
				Claim: &codectypes.Any{
					TypeUrl:              "/Baseledger.baseledger.bridge.MsgValidatorPowerChangedClaim",
					Value:                []byte{},
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     []byte{},
					XXX_sizecache:        0,
				},
				UbtPrices: []sdk.Int{sdk.NewIntFromUint64(1), sdk.NewIntFromUint64(2)},
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	testKeepers := keepertest.BaseledgerbridgeKeeper(t)
	bridge.InitGenesis(testKeepers.Context, *testKeepers.BridgeKeeper, genesisState)
	got := bridge.ExportGenesis(testKeepers.Context, *testKeepers.BridgeKeeper)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.OrchestratorValidatorAddressList, got.OrchestratorValidatorAddressList)
	require.ElementsMatch(t, genesisState.Attestations, got.Attestations)

	require.Equal(t, uint64(5), got.LastObservedNonce)
	// this line is used by starport scaffolding # genesis/test/assert
}
