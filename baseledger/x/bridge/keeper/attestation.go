package keeper

import (
	"errors"
	"fmt"
	"math/big"
	"sort"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Baseledger/baseledger/x/bridge/types"
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
	val := k.GetOrchestratorValidator(ctx, claim.GetClaimer().String())
	if val == nil {
		panic("Could not find val for delegate key, should be checked by now")
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
		return nil, errors.New("last event nonce error")
	}

	// Tries to get an attestation with the same eventNonce and claim as the claim that was submitted.
	hash, err := claim.ClaimHash()

	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to compute claim hash")
	}

	ubtPrice := claim.GetUbtPriceAsInt()

	if ubtPrice.IsNil() || ubtPrice.IsNegative() || ubtPrice.IsZero() {
		// TODO: Ognjen - Log or err?
		k.Logger(ctx).Error("claim ubt price is nil, negative or zero ",
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
			"claimer", claim.GetClaimer(),
			"provided ubt price", ubtPrice.String(),
		)
		// return nil, errors.New("claim ubt price is nil")

		// TODO: Ognjen - If proper value not provided we set it to zero
		// make sure not to affect the later average calculation
		ubtPrice = sdk.ZeroInt()
	}

	att := k.GetAttestation(ctx, claim.GetEventNonce(), hash)

	// If it does not exist, create a new one.
	if att == nil {
		att = &types.Attestation{
			Observed:  false,
			Votes:     []string{},
			Height:    uint64(ctx.BlockHeight()),
			Claim:     anyClaim,
			UbtPrices: []sdk.Int{},
		}
	}

	// Add the validator's vote to this attestation
	att.Votes = append(att.Votes, valAddr.String())

	att.UbtPrices = append(att.UbtPrices, ubtPrice)

	k.SetAttestation(ctx, claim.GetEventNonce(), hash, att)

	k.SetLastEventNonceByValidator(ctx, valAddr, claim.GetEventNonce())

	return att, nil
}

func (k Keeper) GetAttestationMapping(ctx sdk.Context) (attestationMapping map[uint64][]types.Attestation, orderedKeys []uint64) {
	attestationMapping = make(map[uint64][]types.Attestation)
	k.IterateAttestations(ctx, func(_ []byte, att types.Attestation) bool {
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		if val, ok := attestationMapping[claim.GetEventNonce()]; !ok {
			attestationMapping[claim.GetEventNonce()] = []types.Attestation{att}
		} else {
			attestationMapping[claim.GetEventNonce()] = append(val, att)
		}
		return false
	})
	orderedKeys = make([]uint64, 0, len(attestationMapping))
	for k := range attestationMapping {
		orderedKeys = append(orderedKeys, k)
	}
	sort.Slice(orderedKeys, func(i, j int) bool { return orderedKeys[i] < orderedKeys[j] })

	return
}

// IterateAttestations iterates through all attestations
func (k Keeper) IterateAttestations(ctx sdk.Context, cb func([]byte, types.Attestation) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.OracleAttestationKey
	iter := store.Iterator(prefixRange([]byte(prefix)))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		att := types.Attestation{
			Observed: false,
			Votes:    []string{},
			Height:   0,
			Claim: &codectypes.Any{
				TypeUrl:              "",
				Value:                []byte{},
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     []byte{},
				XXX_sizecache:        0,
			},
		}
		k.cdc.MustUnmarshal(iter.Value(), &att)
		// cb returns true to stop early
		if cb(iter.Key(), att) {
			return
		}
	}
}

// GetMostRecentAttestations returns sorted (by nonce) attestations up to a provided limit number of attestations
// Note: calls GetAttestationMapping in the hopes that there are potentially many attestations
// which are distributed between few nonces to minimize sorting time
func (k Keeper) GetMostRecentAttestations(ctx sdk.Context, limit uint64) []types.Attestation {
	attestationMapping, keys := k.GetAttestationMapping(ctx)
	attestations := make([]types.Attestation, 0, limit)

	// Iterate the nonces and collect the attestations
	count := 0
	for _, nonce := range keys {
		if count >= int(limit) {
			break
		}
		for _, att := range attestationMapping[nonce] {
			if count >= int(limit) {
				break
			}
			attestations = append(attestations, att)
			count++
		}
	}

	return attestations
}

// TryAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) TryAttestation(ctx sdk.Context, att *types.Attestation) {
	claim, err := k.UnpackAttestationClaim(att)
	if err != nil {
		panic("could not cast to claim")
	}
	hash, err := claim.ClaimHash()
	if err != nil {
		panic("unable to compute claim hash")
	}
	// If the attestation has not yet been Observed, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the attestation from accidentally being applied twice.
	if !att.Observed {
		// Sum the current powers of all validators who have voted and see if it passes the current threshold
		// TODO: The different integer types and math here needs a careful review
		totalPower := k.StakingKeeper.GetLastTotalPower(ctx)
		requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
		attestationPower := sdk.NewInt(0)
		for _, validator := range att.Votes {
			val, err := sdk.ValAddressFromBech32(validator)
			if err != nil {
				panic(err)
			}
			validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, val)
			// Add it to the attestation power's sum
			attestationPower = attestationPower.Add(sdk.NewInt(validatorPower))
			// If the power of all the validators that have voted on the attestation is higher or equal to the threshold,
			// process the attestation, set Observed to true, and break
			if attestationPower.GTE(requiredPower) {
				lastEventNonce := k.GetLastObservedEventNonce(ctx)
				// this check is performed at the next level up so this should never panic
				// outside of programmer error.
				if claim.GetEventNonce() != lastEventNonce+1 {
					panic("attempting to apply events to state out of order")
				}
				k.SetLastObservedEventNonce(ctx, claim.GetEventNonce())
				// TODO: Ognjen - Last observed eth height is used for gravity bridge functionality
				// which is not used by us atm (cleanupTimedOutBatches, cleanupTimedOutLogicCalls)
				// reintroduce if this functionly proves to be necessary or delete the line
				// k.SetLastObservedEthereumBlockHeight(ctx, claim.GetBlockHeight())

				att.Observed = true
				k.SetAttestation(ctx, claim.GetEventNonce(), hash, att)

				k.processAttestation(ctx, att, claim)

				k.Logger(ctx).Info("Attestation observed",
					"claim type", claim.GetType(),
					"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
					"nonce", fmt.Sprint(claim.GetEventNonce()),
				)
				break
			}
		}
	} else {
		// We panic here because this should never happen
		panic("attempting to process observed attestation")
	}
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, att *types.Attestation, claim types.EthereumClaim) {
	hash, err := claim.ClaimHash()
	if err != nil {
		panic("unable to compute claim hash")
	}
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.AttestationHandler.Handle(xCtx, *att, claim); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", allowing the oracle to progress properly
		k.Logger(ctx).Error("attestation failed",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
	} else {
		commit() // persist transient storage
	}
}

// DeleteAttestation deletes the given attestation
func (k Keeper) DeleteAttestation(ctx sdk.Context, att types.Attestation) {
	claim, err := k.UnpackAttestationClaim(&att)
	if err != nil {
		panic("Bad Attestation in DeleteAttestation")
	}
	hash, err := claim.ClaimHash()
	if err != nil {
		panic(sdkerrors.Wrap(err, "unable to compute claim hash"))
	}
	store := ctx.KVStore(k.storeKey)

	store.Delete([]byte(types.GetAttestationKey(claim.GetEventNonce(), hash)))
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

// SetLastObservedEventNonce sets the latest observed event nonce
func (k Keeper) SetLastObservedEventNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.LastObservedEventNonceKey), types.UInt64Bytes(nonce))
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

// SetLastAttestationAverageUbtPrice sets the latest average ubt price for a succesfully procesed attestation
func (k Keeper) SetLastAttestationAverageUbtPrice(ctx sdk.Context, avgUbtPrice *big.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.LastAttestationAvgUbtPrice), types.UInt64Bytes(avgUbtPrice.Uint64()))
}

// GetLastAttestationAverageUbtPrice returns the latest average ubt price for a succesfully procesed attestation
func (k Keeper) GetLastAttestationAverageUbtPrice(ctx sdk.Context) *big.Int {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte(types.LastAttestationAvgUbtPrice))

	if len(bytes) == 0 {
		// this can happen only if the avg price calculation fails for the very first attestation
		// fall back to 1
		return big.NewInt(1)
	}

	return sdk.NewIntFromUint64(types.UInt64FromBytes(bytes)).BigInt()
}
