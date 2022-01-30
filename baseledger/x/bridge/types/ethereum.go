package types

import (
	"bytes"
	"fmt"
	"regexp"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// ETHContractAddressLen is the length of contract address strings
	ETHContractAddressLen = 42

	// ZeroAddress is an EthAddress containing the zero ethereum address
	ZeroAddressString = "0x0000000000000000000000000000000000000000"
)

// Regular EthAddress
type EthAddress struct {
	address string
}

// Returns the contained address as a string
func (ea EthAddress) GetAddress() string {
	return ea.address
}

// Sets the contained address, performing validation before updating the value
func (ea EthAddress) SetAddress(address string) error {
	if err := ValidateEthAddress(address); err != nil {
		return err
	}
	ea.address = address
	return nil
}

// Creates a new EthAddress from a string, performing validation and returning any validation errors
func NewEthAddress(address string) (*EthAddress, error) {
	if err := ValidateEthAddress(address); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid input address")
	}
	addr := EthAddress{address}
	return &addr, nil
}

// Returns a new EthAddress with 0x0000000000000000000000000000000000000000 as the wrapped address
func ZeroAddress() EthAddress {
	return EthAddress{ZeroAddressString}
}

// Validates the input string as an Ethereum Address
// Addresses must not be empty, have 42 character length, start with 0x and have 40 remaining characters in [0-9a-fA-F]
func ValidateEthAddress(address string) error {
	if address == "" {
		return fmt.Errorf("empty")
	}
	if len(address) != ETHContractAddressLen {
		return fmt.Errorf("address(%s) of the wrong length exp(%d) actual(%d)", address, ETHContractAddressLen, len(address))
	}
	if !regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString(address) {
		return fmt.Errorf("address(%s) doesn't pass regex", address)
	}

	return nil
}

// Performs validation on the wrapped string
func (ea EthAddress) ValidateBasic() error {
	return ValidateEthAddress(ea.address)
}

// EthAddrLessThan migrates the Ethereum address less than function
func EthAddrLessThan(e EthAddress, o EthAddress) bool {
	return bytes.Compare([]byte(e.GetAddress())[:], []byte(o.GetAddress())[:]) == -1
}
