package types

import (
	"encoding/binary"
	fmt "fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// EthereumClaim represents a claim on ethereum state
type EthereumClaim interface {
	// All Ethereum claims that we relay from the Baseledger contract and into the module
	// have a nonce that is monotonically increasing and unique, since this nonce is
	// issued by the Ethereum contract it is immutable and must be agreed on by all validators
	// any disagreement on what claim goes to what nonce means someone is lying.
	GetEventNonce() uint64
	// The block height that the claimed event occurred on. This EventNonce provides sufficient
	// ordering for the execution of all claims. The block height is used only for batchTimeouts + logicTimeouts
	// when we go to create a new batch we set the timeout some number of batches out from the last
	// known height plus projected block progress since then.
	GetBlockHeight() uint64
	// the delegate address of the claimer, for MsgDepositClaim and MsgWithdrawClaim
	// this is sent in as the sdk.AccAddress of the delegated key. it is up to the user
	// to disambiguate this into a sdk.ValAddress
	GetClaimer() sdk.AccAddress
	// Which type of claim this is
	GetType() ClaimType
	ValidateBasic() error
	// The claim hash of this claim. This is used to store these claims and also used to check if two different
	// validators claims agree. Therefore it's extremely important that this include all elements of the claim
	// with the exception of the orchestrator who sent it in, which will be used as a different part of the index
	ClaimHash() ([]byte, error)
	// The ubt price oracled from orchestrator accompanying the claim
	GetUbtPriceAsInt() sdk.Int
}

//nolint: exhaustivestruct
var (
	_ EthereumClaim = &MsgUbtDepositedClaim{}
	_ EthereumClaim = &MsgValidatorPowerChangedClaim{}
)

// UInt64Bytes uses the SDK byte marshaling to encode a uint64
func UInt64Bytes(n uint64) []byte {
	return sdk.Uint64ToBigEndian(n)
}

// UInt64FromBytes create uint from binary big endian representation
func UInt64FromBytes(s []byte) uint64 {
	return binary.BigEndian.Uint64(s)
}

// GetPrefixFromBech32 returns the human readable part of a bech32 string (excluding the 1 byte)
// Returns an error on too short input or when the 1 byte cannot be found
// Note: This is an excerpt from the Decode function for bech32 strings
func GetPrefixFromBech32(bech32str string) (string, error) {
	if len(bech32str) < 8 {
		return "", fmt.Errorf("invalid bech32 string length %d",
			len(bech32str))
	}
	one := strings.LastIndexByte(bech32str, '1')
	if one < 1 || one+7 > len(bech32str) {
		return "", fmt.Errorf("invalid index of 1")
	}

	return bech32str[:one], nil
}

// GetNativePrefixedAccAddressString treats the input as an AccAddress and re-prefixes the string
// with this chain's configured Bech32AccountAddrPrefix
// Returns an error when input is not a bech32 string or the original string it is already natively prefixed
func GetNativePrefixedAccAddressString(foreignStr string) (string, error) {
	prefix, err := GetPrefixFromBech32(foreignStr)
	if err != nil {
		return "", sdkerrors.Wrap(err, "invalid bech32 string")
	}
	nativePrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	if prefix == nativePrefix {
		return foreignStr, nil
	}

	return nativePrefix + foreignStr[len(prefix):], nil
}

// GetNativePrefixedAccAddress re-prefixes the input AccAddress with this chain's configured Bech32AccountAddrPrefix
func GetNativePrefixedAccAddress(foreignAddr sdk.AccAddress) (sdk.AccAddress, error) {
	nativeStr, err := GetNativePrefixedAccAddressString(foreignAddr.String())
	if err != nil {
		return nil, err
	}
	return sdk.AccAddressFromBech32(nativeStr)
}
