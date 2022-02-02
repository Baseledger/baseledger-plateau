package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdk.NewInt(66)

	// ParamsStoreKeyWorktokenEurPrice storesd the price of 1 worktoken in EUR
	ParamsStoreKeyWorktokenEurPrice = []byte("WorktokenEurPrice")

	// Ensure that params implements the proper interface
	_ paramtypes.ParamSet = &Params{
		WorktokenEurPrice: "",
	}
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:                DefaultParams(),
		Attestations:          []Attestation{},
		OrchestratorAddresses: []MsgSetOrchestratorAddress{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
