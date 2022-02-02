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

	return genesis
}
