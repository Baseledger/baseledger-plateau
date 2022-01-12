package cli

import (
	"strconv"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdDelegateKeysByOrchestratorAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-keys-by-orchestrator-address [orchestrator-address]",
		Short: "Query delegateKeysByOrchestratorAddress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqOrchestratorAddress := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryDelegateKeysByOrchestratorAddressRequest{

				OrchestratorAddress: reqOrchestratorAddress,
			}

			res, err := queryClient.DelegateKeysByOrchestratorAddress(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
