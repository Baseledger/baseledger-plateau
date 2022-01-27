package baseledger

import (
	"math/rand"

	"github.com/Baseledger/baseledger/testutil/sample"
	baseledgersimulation "github.com/Baseledger/baseledger/x/proof/simulation"
	"github.com/Baseledger/baseledger/x/proof/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = baseledgersimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCreateBaseledgerTransaction = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateBaseledgerTransaction int = 100

	opWeightMsgUpdateBaseledgerTransaction = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateBaseledgerTransaction int = 100

	opWeightMsgDeleteBaseledgerTransaction = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteBaseledgerTransaction int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	baseledgerGenesis := types.GenesisState{
		BaseledgerTransactionList: []types.BaseledgerTransaction{
			{
				Id:      0,
				Creator: sample.AccAddress(),
			},
			{
				Id:      1,
				Creator: sample.AccAddress(),
			},
		},
		BaseledgerTransactionCount: 2,
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&baseledgerGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateBaseledgerTransaction int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateBaseledgerTransaction, &weightMsgCreateBaseledgerTransaction, nil,
		func(_ *rand.Rand) {
			weightMsgCreateBaseledgerTransaction = defaultWeightMsgCreateBaseledgerTransaction
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateBaseledgerTransaction,
		baseledgersimulation.SimulateMsgCreateBaseledgerTransaction(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateBaseledgerTransaction int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateBaseledgerTransaction, &weightMsgUpdateBaseledgerTransaction, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateBaseledgerTransaction = defaultWeightMsgUpdateBaseledgerTransaction
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateBaseledgerTransaction,
		baseledgersimulation.SimulateMsgUpdateBaseledgerTransaction(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteBaseledgerTransaction int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteBaseledgerTransaction, &weightMsgDeleteBaseledgerTransaction, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteBaseledgerTransaction = defaultWeightMsgDeleteBaseledgerTransaction
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteBaseledgerTransaction,
		baseledgersimulation.SimulateMsgDeleteBaseledgerTransaction(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
