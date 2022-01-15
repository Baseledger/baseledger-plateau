package keeper

import (
	"errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
)

// TODO-JT: carefully look at atomicity of this function
func (k Keeper) Attest(
	ctx sdk.Context,
	claim types.EthereumClaim,
	anyClaim *codectypes.Any,
) (*types.Attestation, error) {
	if err := sdk.VerifyAddressFormat(claim.GetClaimer()); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid claimer address")
	}
	val, found := k.GetOrchestratorValidator(ctx, claim.GetClaimer())
	if !found {
		panic("Could not find ValAddr for delegate key, should be checked by now")
	}

	valAddr := val.GetOperator()
	if err := sdk.VerifyAddressFormat(valAddr); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid orchestrator validator address")
	}
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processAttestation as well,
	// but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two claims with the same nonce.
	// This prevents there being two attestations with the same nonce that get 2/3s of the votes
	// in the endBlocker.
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, valAddr)
	if claim.GetEventNonce() != lastEventNonce+1 {
		return nil, errors.New("Last event nonce error")
	}

	// Tries to get an attestation with the same eventNonce and claim as the claim that was submitted.
	hash, err := claim.ClaimHash()

	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to compute claim hash")
	}

	att := k.GetAttestation(ctx, claim.GetEventNonce(), hash)

	// If it does not exist, create a new one.
	if att == nil {
		att = &types.Attestation{
			Observed: false,
			Votes:    []string{},
			Height:   uint64(ctx.BlockHeight()),
			Claim:    anyClaim,
		}
	}

	// Add the validator's vote to this attestation
	att.Votes = append(att.Votes, valAddr.String())

	k.SetAttestation(ctx, claim.GetEventNonce(), hash, att)

	k.SetLastEventNonceByValidator(ctx, valAddr, claim.GetEventNonce())

	return att, nil
}

func (k Keeper) GetAttestationMapping(ctx sdk.Context) (attestationMapping map[uint64][]types.Attestation, orderedKeys []uint64) {
	// TODO: Ognjen - implement
	return
}

func (k Keeper) TryAttestation(ctx sdk.Context, att *types.Attestation) {
	// TODO: Ognjen - implement
}

// DeleteAttestation deletes the given attestation
func (k Keeper) DeleteAttestation(ctx sdk.Context, att types.Attestation) {
	// TODO: Ognjen - implement
}

// SetAttestation sets the attestation in the store
func (k Keeper) SetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte, att *types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	aKey := []byte(types.GetAttestationKey(eventNonce, claimHash))
	store.Set(aKey, k.cdc.MustMarshal(att))
}

// GetAttestation return an attestation given a nonce
func (k Keeper) GetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte) *types.Attestation {
	store := ctx.KVStore(k.storeKey)
	aKey := []byte(types.GetAttestationKey(eventNonce, claimHash))
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil
	}
	var att types.Attestation
	k.cdc.MustUnmarshal(bz, &att)
	return &att
}

// GetLastObservedEventNonce returns the latest observed event nonce
func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte(types.LastObservedEventNonceKey))

	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

// GetLastEventNonceByValidator returns the latest event nonce for a given validator
func (k Keeper) GetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress) uint64 {
	if err := sdk.VerifyAddressFormat(validator); err != nil {
		panic(sdkerrors.Wrap(err, "invalid validator address"))
	}
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte(types.GetLastEventNonceByValidatorKey(validator)))

	if len(bytes) == 0 {
		// in the case that we have no existing value this is the first
		// time a validator is submitting a claim. Since we don't want to force
		// them to replay the entire history of all events ever we can't start
		// at zero
		lastEventNonce := k.GetLastObservedEventNonce(ctx)
		if lastEventNonce >= 1 {
			return lastEventNonce - 1
		} else {
			return 0
		}
	}
	return types.UInt64FromBytes(bytes)
}

// setLastEventNonceByValidator sets the latest event nonce for a give validator
func (k Keeper) SetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress, nonce uint64) {
	if err := sdk.VerifyAddressFormat(validator); err != nil {
		panic(sdkerrors.Wrap(err, "invalid validator address"))
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.GetLastEventNonceByValidatorKey(validator)), types.UInt64Bytes(nonce))
}
