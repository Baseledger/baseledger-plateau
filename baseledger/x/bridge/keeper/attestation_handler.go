package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	baseledgercommon "github.com/Baseledger/baseledger/common"
	"github.com/Baseledger/baseledger/x/bridge/types"
	distypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// Check that distKeeper implements the expected type
var _ types.DistributionKeeper = (*distrkeeper.Keeper)(nil)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	// NOTE: If you add anything to this struct, add a nil check to ValidateMembers below!
	keeper     *Keeper
	bankKeeper *bankkeeper.BaseKeeper
	distKeeper *distrkeeper.Keeper
}

// Check for nil members
func (a AttestationHandler) ValidateMembers() {
	if a.keeper == nil {
		panic("Nil keeper!")
	}
	if a.bankKeeper == nil {
		panic("Nil bankKeeper!")
	}
	if a.distKeeper == nil {
		panic("Nil distKeeper!")
	}
}

// TODO skos: change this to send to faucet or something like that later
// SendToCommunityPool handles sending incorrect deposits to the community pool, since the deposits
// have already been made on Ethereum there's nothing we can do to reverse them, and we should at least
// make use of the tokens which would otherwise be lost
func (a AttestationHandler) SendToCommunityPool(ctx sdk.Context, coins sdk.Coins) error {
	if err := a.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distypes.ModuleName, coins); err != nil {
		return sdkerrors.Wrap(err, "transfer to community pool failed")
	}
	feePool := (*a.distKeeper).GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(coins...)...)
	(*a.distKeeper).SetFeePool(ctx, feePool)
	return nil
}

// Handle is the entry point for Attestation processing.
func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation, claim types.EthereumClaim) error {
	switch claim := claim.(type) {
	// deposit in this context means a deposit into the Ethereum side of the bridge
	case *types.MsgUbtDepositedClaim:
		invalidAddress := false
		receiverAddress, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
		if err != nil {
			invalidAddress = true
		}
		tokenAddress, err := types.NewEthAddress(claim.TokenContract)
		// these are not possible unless the validators get together and submit
		// a bogus event, this would create lost tokens stuck in the bridge
		// and not accessible to anyone
		if err != nil {
			hash, _ := claim.ClaimHash()
			a.keeper.Logger(ctx).Error("Invalid token contract",
				"cause", err.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrap(err, "invalid token contract on claim")
		}

		_, err = types.NewEthAddress(claim.EthereumSender)
		if err != nil {
			hash, _ := claim.ClaimHash()
			a.keeper.Logger(ctx).Error("Invalid ethereum sender",
				"cause", err.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrap(err, "invalid ethereum sender on claim")
		}

		// While not strictly necessary, explicitly making the receiver a native address
		// insulates us from the implicit address conversion done in x/bank's account store iterator
		nativeReceiver, err := types.GetNativePrefixedAccAddress(receiverAddress)

		if err != nil {
			invalidAddress = true
		}

		// TODO: BAS-120 - Introduce fallback to previous att price in this one is nil, negative or zero
		amountOfWorkTokensToSend := calculateAmountOfWorkTokens(claim.Amount, att.AvgUbtPrice)
		coins := sdk.Coins{sdk.NewCoin(baseledgercommon.WorkTokenDenom, amountOfWorkTokensToSend)}

		// TODO: Skos - what is impossible amount? i think we can keep this check just in case even though conversion should stop it
		// Make sure that users are not bridging an impossible amount
		prevSupply := a.bankKeeper.GetSupply(ctx, baseledgercommon.WorkTokenDenom)
		newSupply := new(big.Int).Add(prevSupply.Amount.BigInt(), amountOfWorkTokensToSend.BigInt())
		if newSupply.BitLen() > 256 { // new supply overflows uint256
			a.keeper.Logger(ctx).Error("Deposit Overflow",
				"claim type", claim.GetType(),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
				"amount", claim.Amount,
				"prev supply", prevSupply.Amount,
				"new supply", newSupply,
			)
			return sdkerrors.Wrap(types.ErrIntOverflowAttestation, "invalid supply after UbtDeposit attestation")
		}

		faucetAddress, err := sdk.AccAddressFromBech32(baseledgercommon.UbtFaucetAddress)

		if err != nil {
			panic("Faucet address invalid")
		}

		if !invalidAddress {
			// TODO: Ognjen - remove logging if obsolete after implementation
			a.keeper.Logger(ctx).Info("Worktokens are ready to be sent",
				"nonce", fmt.Sprint(claim.GetEventNonce()),
				"deposited ubt amount", fmt.Sprint(claim.Amount),
				"average ubt price", fmt.Sprint(att.AvgUbtPrice.String()),
				"amount of worktokens", fmt.Sprint(amountOfWorkTokensToSend),
			)

			if err := a.bankKeeper.SendCoins(ctx, faucetAddress, nativeReceiver, coins); err != nil {
				// in this case sending failed, log to be able to resolve
				hash, _ := claim.ClaimHash()
				a.keeper.Logger(ctx).Error("Failed sending work tokens from faucet to receiver",
					"cause", err.Error(),
					"claim type", claim.GetType(),
					"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
					"nonce", fmt.Sprint(claim.GetEventNonce()),
				)
				return sdkerrors.Wrapf(err, "send work coins: %s", coins)
			}
		} else {
			hash, _ := claim.ClaimHash()
			a.keeper.Logger(ctx).Error("Invalid receiver address",
				"cause", err.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute("MsgUbtDepositedAmount", amountOfWorkTokensToSend.String()),
				sdk.NewAttribute("MsgUbtDepositedNonce", strconv.Itoa(int(claim.GetEventNonce()))),
				sdk.NewAttribute("MsgUbtDepositedToken", tokenAddress.GetAddress()),
			),
		)
	case *types.MsgValidatorPowerChangedClaim:
		tokenAddress, err := types.NewEthAddress(claim.TokenContract)
		// these are not possible unless the validators get together and submit
		// a bogus event, this would create lost tokens stuck in the bridge
		// and not accessible to anyone
		if err != nil {
			hash, _ := claim.ClaimHash()
			a.keeper.Logger(ctx).Error("Invalid token contract",
				"cause", err.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrap(err, "invalid token contract on claim")
		}

		_, err = types.NewEthAddress(claim.EthereumSender)
		if err != nil {
			hash, _ := claim.ClaimHash()
			a.keeper.Logger(ctx).Error("Invalid ethereum sender",
				"cause", err.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrap(err, "invalid ethereum sender on claim")
		}

		valAddr, err := sdk.ValAddressFromBech32(claim.CosmosReceiver)
		if err != nil {
			hash, _ := claim.ClaimHash()
			a.keeper.Logger(ctx).Error("Invalid validator address",
				"cause", err.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrap(err, "invalid validator address on claim")
		}

		validator, found := a.keeper.StakingKeeper.GetValidator(ctx, valAddr)

		if !found {
			hash, _ := claim.ClaimHash()
			a.keeper.Logger(ctx).Error("Validator not found",
				"cause", err.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrap(err, "can not find validator specified on claim")
		}

		faucetAddress, err := sdk.AccAddressFromBech32(baseledgercommon.UbtFaucetAddress)
		if err != nil {
			panic("Faucet address invalid")
		}

		stakingIncreased := true
		if claim.Amount.LT(validator.Tokens) {
			stakingIncreased = false
		}

		stakingAmountChange := claim.Amount.Sub(validator.Tokens).Abs()

		if stakingIncreased {
			_, err = a.keeper.StakingKeeper.Delegate(ctx, faucetAddress, stakingAmountChange, 1, validator, true)

			if err != nil {
				hash, _ := claim.ClaimHash()
				a.keeper.Logger(ctx).Error("Could not delegate to validator",
					"cause", err.Error(),
					"claim type", claim.GetType(),
					"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
					"nonce", fmt.Sprint(claim.GetEventNonce()),
				)
				return sdkerrors.Wrap(err, "could not delegate to validator specified on claim")
			}
		} else {
			_, err := a.keeper.StakingKeeper.Undelegate(ctx, faucetAddress, valAddr, sdk.NewDecFromInt(stakingAmountChange))
			if err != nil {
				hash, _ := claim.ClaimHash()
				a.keeper.Logger(ctx).Error("Could not undelegate from validator",
					"cause", err.Error(),
					"claim type", claim.GetType(),
					"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
					"nonce", fmt.Sprint(claim.GetEventNonce()),
				)
				return sdkerrors.Wrap(err, "could not undelegate from validator specified on claim")
			}
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute("MsgValidatorPowerChangedAmount", claim.Amount.String()),
				sdk.NewAttribute("MsgValidatorPowerChangedNonce", strconv.Itoa(int(claim.GetEventNonce()))),
				sdk.NewAttribute("MsgValidatorPowerChangedToken", tokenAddress.GetAddress()),
			),
		)
	default:
		panic(fmt.Sprintf("Invalid event type for attestations %s", claim.GetType()))
	}
	return nil
}

func calculateAmountOfWorkTokens(depositedUbtAmount sdk.Int, averagePrice sdk.Int) sdk.Int {
	// TODO: BAS-121 - Move this hardcoded value to config or somewhere
	// TODO: Ognjen - Verify calculation
	worktokenEurPrice, _ := sdk.NewDecFromStr("0.1")
	worktokenEurPriceInt := sdk.NewIntFromBigInt(worktokenEurPrice.BigInt())

	depositedEurValueInt := depositedUbtAmount.Mul(averagePrice)

	return depositedEurValueInt.Quo(worktokenEurPriceInt)
}