package types

import (
	"bytes"
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		WorktokenEurPrice:       "0.1",
		BaseledgerFaucetAddress: "baseledger1xgs5tamqre7rkz5q7d5fegjsdwufxxvt36w0a0",
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamsStoreKeyWorktokenEurPrice, &p.WorktokenEurPrice, validateWorktokenEurPrice),
		paramtypes.NewParamSetPair(ParamsStoreKeyBaseledgerFaucetAddress, &p.BaseledgerFaucetAddress, validateBaseledgerFaucetAddress),
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateWorktokenEurPrice(p.WorktokenEurPrice); err != nil {
		return sdkerrors.Wrap(err, "worktoken eur price")
	}

	if err := validateBaseledgerFaucetAddress(p.BaseledgerFaucetAddress); err != nil {
		return sdkerrors.Wrap(err, "baseledger faucet address")
	}

	return nil
}

func validateWorktokenEurPrice(i interface{}) error {
	value, err := sdk.NewDecFromStr(fmt.Sprint(i))
	if err != nil {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value.LTE(sdk.ZeroDec()) {
		return fmt.Errorf("invalid parameter value: %T", i)
	}

	return nil
}

func validateBaseledgerFaucetAddress(i interface{}) error {
	_, err := sdk.AccAddressFromBech32(fmt.Sprint(i))

	if err != nil {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
