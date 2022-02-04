package bridge

import (
	"fmt"

	"github.com/Baseledger/baseledger/x/bridge/keeper"
	"github.com/Baseledger/baseledger/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the orchestratorValidatorAddress
	for _, elem := range genState.OrchestratorValidatorAddressList {
		k.SetOrchestratorValidatorAddress(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
	k.SetLastObservedEventNonce(ctx, genState.LastObservedNonce)

	for _, att := range genState.Attestations {
		att := att
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		// TODO: block height?
		hash, err := claim.ClaimHash()
		if err != nil {
			panic(fmt.Errorf("error when computing ClaimHash for %v", hash))
		}
		k.SetAttestation(ctx, claim.GetEventNonce(), hash, &att)
	}

	// reset attestation state of specific validators
	// this must be done after the above to be correct
	for _, att := range genState.Attestations {
		att := att
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}
		// reconstruct the latest event nonce for every validator
		// if somehow this genesis state is saved when all attestations
		// have been cleaned up GetLastEventNonceByValidator handles that case
		//
		// if we where to save and load the last event nonce for every validator
		// then we would need to carry that state forever across all chain restarts
		// but since we've already had to handle the edge case of new validators joining
		// while all attestations have already been cleaned up we can do this instead and
		// not carry around every validators event nonce counter forever.
		for _, vote := range att.Votes {
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
			last := k.GetLastEventNonceByValidator(ctx, val)
			if claim.GetEventNonce() > last {
				k.SetLastEventNonceByValidator(ctx, val, claim.GetEventNonce())
			}
		}
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.LastObservedNonce = k.GetLastObservedEventNonce(ctx)

	attestationMap, attestationKeys := k.GetAttestationMapping(ctx)

	// export attestations from state
	for _, key := range attestationKeys {
		// TODO: set height = 0?
		genesis.Attestations = append(genesis.Attestations, attestationMap[key]...)
	}

	genesis.OrchestratorValidatorAddressList = k.GetAllOrchestratorValidatorAddress(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
