package types

import "strings"

var (
	// OracleAttestationKey attestation details by nonce and validator address
	// i.e. gravityvaloper1ahx7f8wyertuus9r20284ej0asrs085ceqtfnm
	// An attestation can be thought of as the 'event to be executed' while
	// the Claims are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple claims vote on and
	// eventually executes
	OracleAttestationKey = "OracleAttestationKey"
)

// GetAttestationKey returns the following key format
// prefix     nonce                             claim-details-hash
// [0x5][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetAttestationKey(eventNonce uint64, claimHash []byte) string {
	key := make([]byte, len(OracleAttestationKey)+len(UInt64Bytes(0))+len(claimHash))
	copy(key[0:], OracleAttestationKey)
	copy(key[len(OracleAttestationKey):], UInt64Bytes(eventNonce))
	copy(key[len(OracleAttestationKey)+len(UInt64Bytes(0)):], claimHash)
	return convertByteArrToString(key)
}

func convertByteArrToString(value []byte) string {
	var ret strings.Builder
	for i := 0; i < len(value); i++ {
		ret.WriteString(string(value[i]))
	}
	return ret.String()
}
