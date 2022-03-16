package keeper

import (
	"testing"

	"github.com/Baseledger/baseledger/testutil/sample"
	bridgekeeper "github.com/Baseledger/baseledger/x/bridge/keeper"
	bridgetypes "github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/Baseledger/baseledger/x/proof/keeper"
	"github.com/Baseledger/baseledger/x/proof/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

var SampleAccountWithFunds = sample.AccAddressNative()

type TestProofKeepers struct {
	ProofKeeper *keeper.Keeper
	BankKeeper  bankkeeper.BaseKeeper
	Context     sdk.Context
}

func BaseledgerKeeper(t testing.TB) TestProofKeepers {

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keyBridge := sdk.NewKVStoreKey(bridgetypes.StoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	stateStore.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyBridge, sdk.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	cdc := MakeTestMarshaler()
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	paramsKeeper := paramskeeper.NewKeeper(cdc, types.Amino, keyParams, tkeyParams)
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"BaseledgerParams",
	)

	// this is also used to initialize module accounts for all the map keys
	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		types.ModuleName:               {authtypes.Minter, authtypes.Burner},
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
		nil,
	)
	bankKeeper.SetParams(ctx, banktypes.Params{
		SendEnabled:        []*banktypes.SendEnabled{},
		DefaultSendEnabled: true,
	})
	// total supply to track this
	totalSupply := sdk.NewCoins(sdk.NewInt64Coin("stake", 10000000000000000), sdk.NewInt64Coin("work", 10000000000000000))
	faucetSupply := sdk.NewCoins(sdk.NewInt64Coin("stake", 500000000000000), sdk.NewInt64Coin("work", 500000000000000))
	// set up initial accounts
	for name, perms := range maccPerms {
		mod := authtypes.NewEmptyModuleAccount(name, perms...)
		if name == stakingtypes.NotBondedPoolName {
			err := bankKeeper.MintCoins(ctx, types.ModuleName, totalSupply)
			require.NoError(t, err)
			err = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, FaucetAccount, faucetSupply)
			require.NoError(t, err)
		}
		accountKeeper.SetModuleAccount(ctx, mod)
	}

	sampleAccFunds := sdk.NewCoins(sdk.NewInt64Coin("work", 100))
	bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, SampleAccountWithFunds, sampleAccFunds)

	stakingKeeper := stakingkeeper.NewKeeper(cdc, keyStaking, accountKeeper, bankKeeper, getSubspace(paramsKeeper, stakingtypes.ModuleName))
	stakingKeeper.SetParams(ctx, TestingStakeParams)

	bridgeKeeper := bridgekeeper.NewKeeper(
		cdc,
		keyBridge,
		memStoreKey,
		paramsSubspace,
		&bankKeeper,
		&stakingKeeper,
	)

	testBridgeParams := bridgetypes.Params{
		WorktokenEurPrice:       "0.01",
		BaseledgerFaucetAddress: FaucetAccount.String(),
	}
	bridgeKeeper.SetParams(ctx, testBridgeParams)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		bankKeeper,
		bridgeKeeper,
	)

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return TestProofKeepers{
		ProofKeeper: k,
		BankKeeper:  bankKeeper,
		Context:     ctx,
	}
}
