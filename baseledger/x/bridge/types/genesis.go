package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdk.NewInt(66)

	// ParamsStoreKeyWorktokenEurPrice stores the price of 1 worktoken in EUR
	ParamsStoreKeyWorktokenEurPrice = []byte("WorktokenEurPrice")

	// ParamsStoreKeyBaseledgerFaucetAddress stores the faucet address used to send stake and work tokens
	ParamsStoreKeyBaseledgerFaucetAddress = []byte("BaseledgerFaucetAddress")

	// Ensure that params implements the proper interface
	_ paramtypes.ParamSet = &Params{}
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		OrchestratorValidatorAddressList: []OrchestratorValidatorAddress{},
		// this line is used by starport scaffolding # genesis/types/default
		Params:       DefaultParams(),
		Attestations: []Attestation{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in orchestratorValidatorAddress
	orchestratorValidatorAddressIndexMap := make(map[string]struct{})

	for _, elem := range gs.OrchestratorValidatorAddressList {
		index := string(OrchestratorValidatorAddressKey(elem.OrchestratorAddress))
		if _, ok := orchestratorValidatorAddressIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for orchestratorValidatorAddress")
		}
		orchestratorValidatorAddressIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
