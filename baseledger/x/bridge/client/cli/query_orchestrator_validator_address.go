package cli

import (
    "context"
	
    "github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/Baseledger/baseledger/x/bridge/types"
)

func CmdListOrchestratorValidatorAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-orchestrator-validator-address",
		Short: "list all orchestratorValidatorAddress",
		RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx := client.GetClientContextFromCmd(cmd)

            pageReq, err := client.ReadPageRequest(cmd.Flags())
            if err != nil {
                return err
            }

            queryClient := types.NewQueryClient(clientCtx)

            params := &types.QueryAllOrchestratorValidatorAddressRequest{
                Pagination: pageReq,
            }

            res, err := queryClient.OrchestratorValidatorAddressAll(context.Background(), params)
            if err != nil {
                return err
            }

            return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

    return cmd
}

func CmdShowOrchestratorValidatorAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-orchestrator-validator-address [orchestrator-address]",
		Short: "shows a orchestratorValidatorAddress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
            clientCtx := client.GetClientContextFromCmd(cmd)

            queryClient := types.NewQueryClient(clientCtx)

             argOrchestratorAddress := args[0]
            
            params := &types.QueryGetOrchestratorValidatorAddressRequest{
                OrchestratorAddress: argOrchestratorAddress,
                
            }

            res, err := queryClient.OrchestratorValidatorAddress(context.Background(), params)
            if err != nil {
                return err
            }

            return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

    return cmd
}
