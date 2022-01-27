package keeper

import (
	"github.com/Baseledger/baseledger/x/bridge/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// claimHandlerCommon is an internal function that provides common code for processing claims once they are
// translated from the message to the Ethereum claim interface
func (k msgServer) claimHandlerCommon(ctx sdk.Context, msgAny *codectypes.Any, msg types.EthereumClaim) error {
	// Add the claim to the store
	_, err := k.Attest(ctx, msg, msgAny)
	if err != nil {
		return sdkerrors.Wrap(err, "create attestation")
	}
	// TODO: BAS-106 figure out event usages

	// hash, err := msg.ClaimHash()
	// if err != nil {
	// 	return sdkerrors.Wrap(err, "unable to compute claim hash")
	// }

	// Emit the handle message event
	// ctx.EventManager().EmitEvent(
	// 	sdk.NewEvent(
	// 		sdk.EventTypeMessage,
	// 		sdk.NewAttribute(sdk.AttributeKeyModule, string(msg.GetType())),
	// 		// TODO: maybe return something better here? is this the right string representation?
	// 		sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(msg.GetEventNonce(), hash))),
	// 	),
	// )

	return nil
}

func (k msgServer) checkOrchestratorValidatorInSet(ctx sdk.Context, orchestrator string) error {
	orchestratorAddress, err := sdk.AccAddressFromBech32(orchestrator)
	if err != nil {
		return sdkerrors.Wrap(err, "orchestrator acc address invalid")
	}

	orchValidator, found := k.GetOrchestratorValidator(ctx, orchestratorAddress)
	if !found {
		return sdkerrors.Wrap(sdkerrors.Error{}, "Orchestrator address not set")
	}

	validator := k.StakingKeeper.Validator(ctx, orchValidator.GetOperator())
	if validator == nil || !validator.IsBonded() {
		return sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "Orchestrator validator not in active set")
	}

	return nil
}
