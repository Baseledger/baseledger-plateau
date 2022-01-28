package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		BaseledgerTransactionList: []BaseledgerTransaction{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in baseledgerTransaction
	baseledgerTransactionIdMap := make(map[uint64]bool)
	baseledgerTransactionCount := gs.GetBaseledgerTransactionCount()
	for _, elem := range gs.BaseledgerTransactionList {
		if _, ok := baseledgerTransactionIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for baseledgerTransaction")
		}
		if elem.Id >= baseledgerTransactionCount {
			return fmt.Errorf("baseledgerTransaction id should be lower or equal than the last id")
		}
		baseledgerTransactionIdMap[elem.Id] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
