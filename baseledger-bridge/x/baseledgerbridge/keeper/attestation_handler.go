package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
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
		receiverAddress, addressErr := types.IBCAddressFromBech32(claim.CosmosReceiver)
		if addressErr != nil {
			invalidAddress = true
		}
		tokenAddress, errTokenAddress := types.NewEthAddress(claim.TokenContract)
		ethereumSender, errEthereumSender := types.NewEthAddress(claim.EthereumSender)
		// these are not possible unless the validators get together and submit
		// a bogus event, this would create lost tokens stuck in the bridge
		// and not accessible to anyone
		if errTokenAddress != nil {
			hash, _ := claim.ClaimHash()
			a.keeper.Logger(ctx).Error("Invalid token contract",
				"cause", errTokenAddress.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrap(errTokenAddress, "invalid token contract on claim")
		}
		if errEthereumSender != nil {
			hash, _ := claim.ClaimHash()
			a.keeper.Logger(ctx).Error("Invalid ethereum sender",
				"cause", errEthereumSender.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrap(errTokenAddress, "invalid ethereum sender on claim")
		}

		// While not strictly necessary, explicitly making the receiver a native address
		// insulates us from the implicit address conversion done in x/bank's account store iterator
		// TODO BAS-119: nativeReceiver
		_, err := types.GetNativePrefixedAccAddress(receiverAddress)

		if err != nil {
			invalidAddress = true
		}

		// Checks the address if it's inside the blacklisted address list and marks
		// if it's inside the list.
		if a.keeper.IsOnBlacklist(ctx, *ethereumSender) {
			invalidAddress = true
		}

		// TODO: skos revisit this one, it seems like we don't need to add ERC20 lookup for our specific case - BAS-119
		// Check if coin is Cosmos-originated asset and get denom
		// isCosmosOriginated, denom := a.keeper.ERC20ToDenomLookup(ctx, *tokenAddress)
		// TODO: BAS-119 changed this to true to skip minting block
		isCosmosOriginated := true
		denom := "token"
		coins := sdk.Coins{sdk.NewCoin(denom, claim.Amount)}

		if !isCosmosOriginated {
			// We need to mint eth-originated coins (aka vouchers)
			// Make sure that users are not bridging an impossible amount
			prevSupply := a.bankKeeper.GetSupply(ctx, denom)
			newSupply := new(big.Int).Add(prevSupply.Amount.BigInt(), claim.Amount.BigInt())
			if newSupply.BitLen() > 256 { // new supply overflows uint256
				a.keeper.Logger(ctx).Error("Deposit Overflow",
					"claim type", claim.GetType(),
					"nonce", fmt.Sprint(claim.GetEventNonce()),
				)
				return sdkerrors.Wrap(types.ErrIntOverflowAttestation, "invalid supply after SendToCosmos attestation")
			}

			if err := a.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
				// in this case we have lost tokens! They are in the bridge, but not
				// in the community pool our out in some users balance, every instance of this
				// error needs to be detected and resolved
				hash, _ := claim.ClaimHash()
				a.keeper.Logger(ctx).Error("Failed minting",
					"cause", err.Error(),
					"claim type", claim.GetType(),
					"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
					"nonce", fmt.Sprint(claim.GetEventNonce()),
				)
				return sdkerrors.Wrapf(err, "mint vouchers coins: %s", coins)
			}
		}

		if !invalidAddress { // valid address so far, try to lock up the coins in the requested cosmos address

			// TODO: BAS-120 - Introduce fallback to previous att price in this one is nil, negative or zero
			amountOfWorkTokensToSend := calculateAmountOfWorkTokens(claim.Amount, att.AvgUbtPrice)

			// TODO: Ognjen - remove logging if obsolete after implementation
			a.keeper.Logger(ctx).Info("Worktokens are ready to be sent",
				"nonce", fmt.Sprint(claim.GetEventNonce()),
				"deposited ubt amount", fmt.Sprint(claim.Amount),
				"average ubt price", fmt.Sprint(att.AvgUbtPrice.String()),
				"amount of worktokens", fmt.Sprint(amountOfWorkTokensToSend),
			)

			// if err := a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, nativeReceiver, coins); err != nil {
			// 	// someone attempted to send tokens to a blacklisted user from Ethereum, log and send to Community pool
			// 	hash, _ := claim.ClaimHash()
			// 	a.keeper.Logger(ctx).Error("Blacklisted deposit",
			// 		"cause", err.Error(),
			// 		"claim type", claim.GetType(),
			// 		"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
			// 		"nonce", fmt.Sprint(claim.GetEventNonce()),
			// 	)
			// 	invalidAddress = true
			// }
			fmt.Println("BAS-119 FIXME")
		}

		// for whatever reason above, blacklisted, invalid string, etc this deposit is not valid
		// we can't send the tokens back on the Ethereum side, and if we don't put them somewhere on
		// the cosmos side they will be lost an inaccessible even though they are locked in the bridge.
		// so we deposit the tokens into the community pool for later use
		if invalidAddress {
			if err = a.SendToCommunityPool(ctx, coins); err != nil {
				hash, _ := claim.ClaimHash()
				a.keeper.Logger(ctx).Error("Failed community pool send",
					"cause", err.Error(),
					"claim type", claim.GetType(),
					"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
					"nonce", fmt.Sprint(claim.GetEventNonce()),
				)
				return sdkerrors.Wrap(err, "failed to send to Community pool")
			}
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute("MsgSendToCosmosAmount", claim.Amount.String()),
				sdk.NewAttribute("MsgSendToCosmosNonce", strconv.Itoa(int(claim.GetEventNonce()))),
				sdk.NewAttribute("MsgSendToCosmosToken", tokenAddress.GetAddress()),
			),
		)
	case *types.MsgValidatorPowerChangedClaim:
		fmt.Println("BAS-119 Implement MsgValidatorPowerChangedClaim handler")
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
