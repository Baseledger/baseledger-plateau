package types

import "encoding/binary"

var _ binary.ByteOrder

const (
    // OrchestratorValidatorAddressKeyPrefix is the prefix to retrieve all OrchestratorValidatorAddress
	OrchestratorValidatorAddressKeyPrefix = "OrchestratorValidatorAddress/value/"
)

// OrchestratorValidatorAddressKey returns the store key to retrieve a OrchestratorValidatorAddress from the index fields
func OrchestratorValidatorAddressKey(
orchestratorAddress string,
) []byte {
	var key []byte
    
    orchestratorAddressBytes := []byte(orchestratorAddress)
    key = append(key, orchestratorAddressBytes...)
    key = append(key, []byte("/")...)
    
	return key
}