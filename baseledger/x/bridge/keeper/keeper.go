package keeper

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace

		StakingKeeper *stakingkeeper.Keeper
		BankKeeper    *bankkeeper.BaseKeeper

		AttestationHandler interface {
			Handle(sdk.Context, types.Attestation, types.EthereumClaim) error
		}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper *bankkeeper.BaseKeeper,
	distKeeper *distrkeeper.Keeper,

	stakingKeeper *stakingkeeper.Keeper,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	k := &Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,

		StakingKeeper:      stakingKeeper,
		AttestationHandler: nil,
	}

	attestationHandler := AttestationHandler{
		keeper:     k,
		bankKeeper: bankKeeper,
		distKeeper: distKeeper,
	}
	attestationHandler.ValidateMembers()
	k.AttestationHandler = attestationHandler
	k.BankKeeper = bankKeeper

	return k
}

/////////////////////////////
//       PARAMETERS        //
/////////////////////////////

// GetWorktokenEurPrice returns the EUR price of a single worktoken. This is the price
// used during calculation of how many worktokens to send to a UBT depositor in the brigde contract.
//
// This parameter can be changed through a governance param change proposal.
func (k Keeper) GetWorktokenEurPrice(ctx sdk.Context) (a string) {
	k.paramstore.Get(ctx, types.ParamsStoreKeyWorktokenEurPrice, &a)
	return
}

func (k Keeper) SetWorktokenEurPrice(ctx sdk.Context, v string) {
	k.paramstore.Set(ctx, types.ParamsStoreKeyWorktokenEurPrice, v)
}

// GetBaseledgerFaucetAddress returns faucet address used to send work and stake tokens
//
// This parameter can be changed through a governance param change proposal.
func (k Keeper) GetBaseledgerFaucetAddress(ctx sdk.Context) (a string) {
	// just for dev convenience when using starport locally get bob address for faucet
	// to not need to always set param
	useTestKeyRing := viper.GetBool("DEV")
	useTestKeyRing = true
	if useTestKeyRing {
		kr, _ := keyring.New("baseledger", "test", viper.GetString("KEYRING_DIR"), nil)
		keysList, _ := kr.List()
		return string(keysList[1].GetAddress().String())
	}
	k.paramstore.Get(ctx, types.ParamsStoreKeyBaseledgerFaucetAddress, &a)
	return
}

func (k Keeper) SetBaseledgerFaucetAddress(ctx sdk.Context, v string) {
	k.paramstore.Set(ctx, types.ParamsStoreKeyBaseledgerFaucetAddress, v)
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// TODO: Ognjen - Why is this code here and not in attestation.go?
func (k Keeper) UnpackAttestationClaim(att *types.Attestation) (types.EthereumClaim, error) {
	var msg types.EthereumClaim
	err := k.cdc.UnpackAny(att.Claim, &msg)
	if err != nil {
		return nil, err
	} else {
		return msg, nil
	}
}

// prefixRange turns a prefix into a (start, end) range. The start is the given prefix value and
// the end is calculated by adding 1 bit to the start value. Nil is not allowed as prefix.
// 		Example: []byte{1, 3, 4} becomes []byte{1, 3, 5}
// 				 []byte{15, 42, 255, 255} becomes []byte{15, 43, 0, 0}
//
// In case of an overflow the end is set to nil.
//		Example: []byte{255, 255, 255, 255} becomes nil
// MARK finish-batches: this is where some crazy shit happens
func prefixRange(prefix []byte) ([]byte, []byte) {
	if prefix == nil {
		panic("nil key not allowed")
	}
	// special case: no prefix is whole range
	if len(prefix) == 0 {
		return nil, nil
	}

	// copy the prefix and update last byte
	end := make([]byte, len(prefix))
	copy(end, prefix)
	l := len(end) - 1
	end[l]++

	// wait, what if that overflowed?....
	for end[l] == 0 && l > 0 {
		l--
		end[l]++
	}

	// okay, funny guy, you gave us FFF, no end to this range...
	if l == 0 && end[0] == 0 {
		end = nil
	}
	return prefix, end
}
