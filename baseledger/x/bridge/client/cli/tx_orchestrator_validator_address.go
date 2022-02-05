package cli

import (
	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateOrchestratorValidatorAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-orchestrator-validator-address [orchestrator-address]",
		Short: "Create a new orchestratorValidatorAddress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexOrchestratorAddress := args[0]

			// Get value arguments

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateOrchestratorValidatorAddress(
				clientCtx.GetFromAddress().String(),
				indexOrchestratorAddress,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
