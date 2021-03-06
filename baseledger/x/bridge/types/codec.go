package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUbtDepositedClaim{}, "bridge/UbtDepositedClaim", nil)
	cdc.RegisterConcrete(&MsgValidatorPowerChangedClaim{}, "bridge/ValidatorPowerChangedClaim", nil)
	cdc.RegisterConcrete(&MsgCreateOrchestratorValidatorAddress{}, "bridge/CreateOrchestratorValidatorAddress", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUbtDepositedClaim{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgValidatorPowerChangedClaim{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateOrchestratorValidatorAddress{},
	)
	// this line is used by starport scaffolding # 3

	// TODO skos: is this good protoName?
	registry.RegisterInterface(
		"Baseledger.bridge.EthereumClaim",
		(*EthereumClaim)(nil),
		&MsgUbtDepositedClaim{},
		&MsgValidatorPowerChangedClaim{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
