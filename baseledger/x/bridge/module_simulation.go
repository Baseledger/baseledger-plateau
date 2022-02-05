package bridge

import (
	"math/rand"

	"github.com/Baseledger/baseledger/testutil/sample"
	bridgesimulation "github.com/Baseledger/baseledger/x/bridge/simulation"
	"github.com/Baseledger/baseledger/x/bridge/types"
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
	_ = bridgesimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgUbtDepositedClaim = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUbtDepositedClaim int = 100

	opWeightMsgSetOrchestratorAddress = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetOrchestratorAddress int = 100

	opWeightMsgValidatorPowerChangedClaim = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgValidatorPowerChangedClaim int = 100

	opWeightMsgCreateOrchestratorValidatorAddress = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateOrchestratorValidatorAddress int = 100

	opWeightMsgUpdateOrchestratorValidatorAddress = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateOrchestratorValidatorAddress int = 100

	opWeightMsgDeleteOrchestratorValidatorAddress = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteOrchestratorValidatorAddress int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	bridgeGenesis := types.GenesisState{
		OrchestratorValidatorAddressList: []types.OrchestratorValidatorAddress{
			{
				ValidatorAddress:    sample.AccAddress(),
				OrchestratorAddress: "0",
			},
			{
				ValidatorAddress:    sample.AccAddress(),
				OrchestratorAddress: "1",
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&bridgeGenesis)
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

	var weightMsgUbtDepositedClaim int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUbtDepositedClaim, &weightMsgUbtDepositedClaim, nil,
		func(_ *rand.Rand) {
			weightMsgUbtDepositedClaim = defaultWeightMsgUbtDepositedClaim
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUbtDepositedClaim,
		bridgesimulation.SimulateMsgUbtDepositedClaim(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSetOrchestratorAddress int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSetOrchestratorAddress, &weightMsgSetOrchestratorAddress, nil,
		func(_ *rand.Rand) {
			weightMsgSetOrchestratorAddress = defaultWeightMsgSetOrchestratorAddress
		},
	)
	var weightMsgValidatorPowerChangedClaim int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgValidatorPowerChangedClaim, &weightMsgValidatorPowerChangedClaim, nil,
		func(_ *rand.Rand) {
			weightMsgValidatorPowerChangedClaim = defaultWeightMsgValidatorPowerChangedClaim
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgValidatorPowerChangedClaim,
		bridgesimulation.SimulateMsgValidatorPowerChangedClaim(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateOrchestratorValidatorAddress int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateOrchestratorValidatorAddress, &weightMsgCreateOrchestratorValidatorAddress, nil,
		func(_ *rand.Rand) {
			weightMsgCreateOrchestratorValidatorAddress = defaultWeightMsgCreateOrchestratorValidatorAddress
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateOrchestratorValidatorAddress,
		bridgesimulation.SimulateMsgCreateOrchestratorValidatorAddress(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateOrchestratorValidatorAddress int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateOrchestratorValidatorAddress, &weightMsgUpdateOrchestratorValidatorAddress, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateOrchestratorValidatorAddress = defaultWeightMsgUpdateOrchestratorValidatorAddress
		},
	)
	var weightMsgDeleteOrchestratorValidatorAddress int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteOrchestratorValidatorAddress, &weightMsgDeleteOrchestratorValidatorAddress, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteOrchestratorValidatorAddress = defaultWeightMsgDeleteOrchestratorValidatorAddress
		},
	)

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
