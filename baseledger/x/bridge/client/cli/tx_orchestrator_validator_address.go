package cli

import (
	
    "github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/Baseledger/baseledger/x/bridge/types"
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

func CmdUpdateOrchestratorValidatorAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-orchestrator-validator-address [orchestrator-address]",
		Short: "Update a orchestratorValidatorAddress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
            // Get indexes
         indexOrchestratorAddress := args[0]
        
            // Get value arguments
		
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateOrchestratorValidatorAddress(
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

func CmdDeleteOrchestratorValidatorAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-orchestrator-validator-address [orchestrator-address]",
		Short: "Delete a orchestratorValidatorAddress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
             indexOrchestratorAddress := args[0]
            
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteOrchestratorValidatorAddress(
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