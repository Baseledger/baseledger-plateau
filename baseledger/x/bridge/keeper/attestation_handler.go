package keeper

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	baseledgercommon "github.com/Baseledger/baseledger/common"
	"github.com/Baseledger/baseledger/x/bridge/types"
)

// Check that distKeeper implements the expected type
var _ types.DistributionKeeper = (*distrkeeper.Keeper)(nil)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	// NOTE: If you add anything to this struct, add a nil check to ValidateMembers below!
	keeper     *Keeper
	bankKeeper *bankkeeper.BaseKeeper
}

// Check for nil members
func (a AttestationHandler) ValidateMembers() {
	if a.keeper == nil {
		panic("Nil keeper!")
	}
	if a.bankKeeper == nil {
		panic("Nil bankKeeper!")
	}
}

// Handle is the entry point for Attestation processing.
func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation, claim types.EthereumClaim) error {
	switch claim := claim.(type) {
	// deposit in this context means a deposit into the Ethereum side of the bridge
	case *types.MsgUbtDepositedClaim:
		invalidAddress := false
		receiverAddress, err := sdk.AccAddressFromBech32(claim.BaseledgerReceiverAccountAddress)
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

		// TODO: Ognjen - Verify big int overflows in the testing phase
		avgUbtPrice := CalculateAvgUbtPriceForAttestation(att)

		if avgUbtPrice.Cmp(big.NewInt(0)) == 0 {
			avgUbtPrice = a.keeper.GetLastAttestationAverageUbtPrice(ctx)
		} else {
			a.keeper.SetLastAttestationAverageUbtPrice(ctx, avgUbtPrice)
		}

		worktokenEurPrice, _ := sdk.NewDecFromStr(a.keeper.GetWorktokenEurPrice(ctx))
		amountOfWorkTokensToSend := CalculateAmountOfWorkTokens(worktokenEurPrice.BigInt(), claim.Amount.BigInt(), avgUbtPrice)
		// TODO: Ognjen - remove logging if obsolete after implementation
		a.keeper.Logger(ctx).Info("Worktokens are ready to be sent",
			"nonce", fmt.Sprint(claim.GetEventNonce()),
			"deposited ubt amount", fmt.Sprint(claim.Amount),
			"average ubt price", fmt.Sprint(avgUbtPrice.String()),
			"amount of worktokens", fmt.Sprint(amountOfWorkTokensToSend),
		)

		coins := sdk.Coins{sdk.NewCoin(baseledgercommon.WorkTokenDenom, sdk.NewIntFromBigInt(amountOfWorkTokensToSend))}

		// TODO: Skos - what is impossible amount? i think we can keep this check just in case even though conversion should stop it
		// Make sure that users are not bridging an impossible amount
		prevSupply := a.bankKeeper.GetSupply(ctx, baseledgercommon.WorkTokenDenom)
		newSupply := new(big.Int).Add(prevSupply.Amount.BigInt(), amountOfWorkTokensToSend)
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

		faucetAddress, err := sdk.AccAddressFromBech32(a.keeper.GetBaseledgerFaucetAddress(ctx))

		if err != nil {
			panic("Faucet address invalid")
		}

		if !invalidAddress {
			// TODO: Ognjen - remove logging if obsolete after implementation
			a.keeper.Logger(ctx).Info("Worktokens are ready to be sent",
				"nonce", fmt.Sprint(claim.GetEventNonce()),
				"deposited ubt amount", fmt.Sprint(claim.Amount),
				"average ubt price", fmt.Sprint(avgUbtPrice.String()),
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

		a.keeper.Logger(ctx).Info("MsgUbtDepositedClaim success",
			"MsgUbtDepositedAmount", amountOfWorkTokensToSend.String(),
			"MsgUbtDepositedNonce", strconv.Itoa(int(claim.GetEventNonce())),
			"MsgUbtDepositedToken", tokenAddress.GetAddress(),
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

		_, err = types.NewEthAddress(claim.RevenueAddress)
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

		valAddr, err := sdk.ValAddressFromBech32(claim.BaseledgerReceiverValidatorAddress)
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
		if !found || !validator.IsBonded() || validator.IsJailed() {
			hash, _ := claim.ClaimHash()
			a.keeper.Logger(ctx).Error("Validator not in active state",
				"cause", "validator either not found, not bonded, or jailed",
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrap(errors.New("validator not found"), "can not find validator specified on claim")
		}
		faucetAddress, err := sdk.AccAddressFromBech32(a.keeper.GetBaseledgerFaucetAddress(ctx))
		if err != nil {
			panic("Faucet address invalid")
		}

		stakingIncreased := true
		// ubt is 8 decimals and staking power is 6 decimals, so we divide by 100 to remove 2 zeros
		amount := claim.Amount.Quo(sdk.NewInt(100))

		fmt.Printf("CLAIM AMOUNT %v %v\n", claim.Amount, amount)

		if amount.LT(validator.Tokens) {
			stakingIncreased = false
		}

		stakingAmountChange := amount.Sub(validator.Tokens).Abs()

		fmt.Printf("STAKING AMOUNT CHANGE  %v\n", stakingAmountChange)

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

		a.keeper.Logger(ctx).Info("MsgValidatorPowerChangedClaim success",
			"MsgValidatorPowerChangedAmount", amount.String(),
			"MsgValidatorPowerChangedNonce", strconv.Itoa(int(claim.GetEventNonce())),
			"MsgValidatorPowerChangedToken", tokenAddress.GetAddress(),
		)
	default:
		panic(fmt.Sprintf("Invalid event type for attestations %s", claim.GetType()))
	}
	return nil
}
