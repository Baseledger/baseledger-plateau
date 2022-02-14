package keeper

import (
	"testing"

	"github.com/Baseledger/baseledger/x/bridge/keeper"
	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ccrypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

var (
	// InitTokens holds the number of tokens to initialize an account with
	InitTokens = sdk.TokensFromConsensusPower(110, sdk.DefaultPowerReduction)

	// InitCoins holds the number of coins to initialize an account with
	InitCoins = sdk.NewCoins(sdk.NewCoin(TestingStakeParams.BondDenom, InitTokens))

	// StakingAmount holds the staking power to start a validator with
	StakingAmount = sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)

	// StakingCoins holds the staking coins to start a validator with
	StakingCoins = sdk.NewCoins(sdk.NewCoin(TestingStakeParams.BondDenom, StakingAmount))

	// TestingStakeParams is a set of staking params for testing
	TestingStakeParams = stakingtypes.Params{
		UnbondingTime:     100,
		MaxValidators:     10,
		MaxEntries:        10,
		HistoricalEntries: 10000,
		BondDenom:         "stake",
	}

	// ConsPrivKeys generate ed25519 ConsPrivKeys to be used for validator operator keys
	ConsPrivKeys = []ccrypto.PrivKey{
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
	}

	// ConsPubKeys holds the consensus public keys to be used for validator operator keys
	ConsPubKeys = []ccrypto.PubKey{
		ConsPrivKeys[0].PubKey(),
		ConsPrivKeys[1].PubKey(),
		ConsPrivKeys[2].PubKey(),
		ConsPrivKeys[3].PubKey(),
		ConsPrivKeys[4].PubKey(),
	}

	// Orchestrator private keys
	OrchPrivKeys = []ccrypto.PrivKey{
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
	}

	// AccPubKeys holds the pub keys for the account keys
	OrchPubKeys = []ccrypto.PubKey{
		OrchPrivKeys[0].PubKey(),
		OrchPrivKeys[1].PubKey(),
		OrchPrivKeys[2].PubKey(),
		OrchPrivKeys[3].PubKey(),
		OrchPrivKeys[4].PubKey(),
	}
	// AccAddrs holds the sdk.AccAddresses
	OrchAddrs = []sdk.AccAddress{
		sdk.AccAddress(OrchPubKeys[0].Address()),
		sdk.AccAddress(OrchPubKeys[1].Address()),
		sdk.AccAddress(OrchPubKeys[2].Address()),
		sdk.AccAddress(OrchPubKeys[3].Address()),
		sdk.AccAddress(OrchPubKeys[4].Address()),
	}

	// AccPrivKeys generate secp256k1 pubkeys to be used for account pub keys
	AccPrivKeys = []ccrypto.PrivKey{
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
	}

	// AccPubKeys holds the pub keys for the account keys
	AccPubKeys = []ccrypto.PubKey{
		AccPrivKeys[0].PubKey(),
		AccPrivKeys[1].PubKey(),
		AccPrivKeys[2].PubKey(),
		AccPrivKeys[3].PubKey(),
		AccPrivKeys[4].PubKey(),
	}

	// AccAddrs holds the sdk.AccAddresses
	AccAddrs = []sdk.AccAddress{
		sdk.AccAddress(AccPubKeys[0].Address()),
		sdk.AccAddress(AccPubKeys[1].Address()),
		sdk.AccAddress(AccPubKeys[2].Address()),
		sdk.AccAddress(AccPubKeys[3].Address()),
		sdk.AccAddress(AccPubKeys[4].Address()),
	}

	// ValAddrs holds the sdk.ValAddresses
	ValAddrs = []sdk.ValAddress{
		sdk.ValAddress(AccPubKeys[0].Address()),
		sdk.ValAddress(AccPubKeys[1].Address()),
		sdk.ValAddress(AccPubKeys[2].Address()),
		sdk.ValAddress(AccPubKeys[3].Address()),
		sdk.ValAddress(AccPubKeys[4].Address()),
	}
)

var (
	// ModuleBasics is a mock module basic manager for testing
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distribution.AppModuleBasic{},
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		vesting.AppModuleBasic{},
	)
)

type TestKeepers struct {
	BridgeKeeper  *keeper.Keeper
	AccountKeeper authkeeper.AccountKeeper
	StakingKeeper stakingkeeper.Keeper
	BankKeeper    bankkeeper.BaseKeeper
	Context       sdk.Context
	Cdc           codec.Codec
	LegacyAmino   *codec.LegacyAmino
}

func BaseledgerbridgeKeeper(t testing.TB) TestKeepers {
	t.Helper()
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	stateStore.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)

	require.NoError(t, stateStore.LoadLatestVersion())

	// registry := codectypes.NewInterfaceRegistry()
	// cdc := codec.NewProtoCodec(registry)
	cdc := MakeTestMarshaler()

	paramsKeeper := paramskeeper.NewKeeper(cdc, types.Amino, keyParams, tkeyParams)
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"BaseledgerbridgeParams",
	)

	// this is also used to initialize module accounts for all the map keys
	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		types.ModuleName:               {authtypes.Minter, authtypes.Burner},
	}
	blockedAddr := make(map[string]bool, len(maccPerms))
	for acc := range maccPerms {
		blockedAddr[authtypes.NewModuleAddress(acc).String()] = true
	}

	accountKeeper := authkeeper.NewAccountKeeper(
		cdc,
		keyAcc, // target store
		getSubspace(paramsKeeper, authtypes.ModuleName),
		authtypes.ProtoBaseAccount, // prototype
		maccPerms,
	)

	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		keyBank,
		accountKeeper,
		getSubspace(paramsKeeper, banktypes.ModuleName),
		blockedAddr,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	bankKeeper.SetParams(ctx, banktypes.Params{
		SendEnabled:        []*banktypes.SendEnabled{},
		DefaultSendEnabled: true,
	})
	stakingKeeper := stakingkeeper.NewKeeper(cdc, keyStaking, accountKeeper, bankKeeper, getSubspace(paramsKeeper, stakingtypes.ModuleName))
	stakingKeeper.SetParams(ctx, TestingStakeParams)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		&bankKeeper,
		&stakingKeeper,
	)

	// total supply to track this
	// totalSupply := sdk.NewCoins(sdk.NewInt64Coin("stake", 100000000), sdk.NewInt64Coin("work", 100000000))

	// set up initial accounts
	for name, perms := range maccPerms {
		mod := authtypes.NewEmptyModuleAccount(name, perms...)
		// if name == stakingtypes.NotBondedPoolName {
		// 	err := bankKeeper.MintCoins(ctx, types.ModuleName, totalSupply)
		// 	require.NoError(t, err)
		// 	fmt.Printf("SEND %v -> %v\n", types.ModuleName, mod.Name)
		// 	err = bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, mod.Name, totalSupply)
		// 	require.NoError(t, err)
		// }
		accountKeeper.SetModuleAccount(ctx, mod)
	}

	ps := types.DefaultParams()
	ps.ParamSetPairs()

	testParams := types.Params{
		WorktokenEurPrice:       "0.1",
		BaseledgerFaucetAddress: OrchAddrs[0].String(),
	}

	// Initialize params
	k.SetParams(ctx, testParams)

	return TestKeepers{
		BridgeKeeper:  k,
		Context:       ctx,
		BankKeeper:    bankKeeper,
		Cdc:           cdc,
		LegacyAmino:   types.Amino,
		AccountKeeper: accountKeeper,
		StakingKeeper: stakingKeeper,
	}
}

// getSubspace returns a param subspace for a given module name.
func getSubspace(k paramskeeper.Keeper, moduleName string) paramstypes.Subspace {
	subspace, _ := k.GetSubspace(moduleName)
	return subspace
}

// MakeTestMarshaler creates a proto codec for use in testing
func MakeTestMarshaler() codec.Codec {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	types.RegisterInterfaces(interfaceRegistry)
	return codec.NewProtoCodec(interfaceRegistry)
}

func SetFiveValidators(t *testing.T, createOrchValAddress bool) TestKeepers {
	t.Helper()

	testKeepers := BaseledgerbridgeKeeper(t)

	// Initialize each of the validators
	sh := staking.NewHandler(testKeepers.StakingKeeper)
	for i := range []int{0, 1, 2, 3, 4} {
		// Initialize the account for the key
		acc := testKeepers.AccountKeeper.NewAccount(
			testKeepers.Context,
			authtypes.NewBaseAccount(AccAddrs[i], AccPubKeys[i], uint64(i), 0),
		)

		// Set the balance for the account
		require.NoError(t, testKeepers.BankKeeper.MintCoins(testKeepers.Context, types.ModuleName, InitCoins))
		testKeepers.BankKeeper.SendCoinsFromModuleToAccount(testKeepers.Context, types.ModuleName, acc.GetAddress(), InitCoins)

		// Set the account in state
		testKeepers.AccountKeeper.SetAccount(testKeepers.Context, acc)

		// Create a validator for that account using some of the tokens in the account
		// and the staking handler
		_, err := sh(
			testKeepers.Context,
			NewTestMsgCreateValidator(ValAddrs[i], ConsPubKeys[i], StakingAmount),
		)

		// Return error if one exists
		require.NoError(t, err)
	}

	// Run the staking endblocker to ensure valset is correct in state
	staking.EndBlocker(testKeepers.Context, testKeepers.StakingKeeper)

	if createOrchValAddress {
		// Register eth addresses and orchestrator address for each validator
		for i, addr := range ValAddrs {
			testKeepers.BridgeKeeper.SetOrchestratorValidatorAddress(testKeepers.Context, types.OrchestratorValidatorAddress{
				ValidatorAddress:    addr.String(),
				OrchestratorAddress: OrchAddrs[i].String(),
			})
		}
	}

	// Return the test input
	return testKeepers
}

func NewTestMsgCreateValidator(address sdk.ValAddress, pubKey ccrypto.PubKey, amt sdk.Int) *stakingtypes.MsgCreateValidator {
	commission := stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	out, err := stakingtypes.NewMsgCreateValidator(
		address, pubKey, sdk.NewCoin("stake", amt),
		stakingtypes.Description{
			Moniker:         "",
			Identity:        "",
			Website:         "",
			SecurityContact: "",
			Details:         "",
		}, commission, sdk.OneInt(),
	)
	if err != nil {
		panic(err)
	}
	return out
}
