package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// EthereumClaim represents a claim on ethereum state
type EthereumClaim interface {
	// All Ethereum claims that we relay from the Gravity contract and into the module
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
}

//nolint: exhaustivestruct
var (
	_ EthereumClaim = &MsgUbtDepositedClaim{}
)

// UInt64Bytes uses the SDK byte marshaling to encode a uint64
func UInt64Bytes(n uint64) []byte {
	return sdk.Uint64ToBigEndian(n)
}

// IBCAddressFromBech32 decodes an IBC-compatible Address from a Bech32
// encoded string
// TODO: This is very similar to the sdk's GetFromBech32 method, but makes no
// assertions about the bech32 prefix (aka "human readable part"), when Gravity
// IBC Forwarding is completed, this function should return invalid prefix errors
func IBCAddressFromBech32(bech32str string) ([]byte, error) {
	if len(bech32str) == 0 {
		return nil, errors.New("Bech32 empty")
	}

	_, bz, err := bech32.DecodeAndConvert(bech32str)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return bz, nil
}
