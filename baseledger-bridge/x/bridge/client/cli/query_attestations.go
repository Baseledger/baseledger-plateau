package cli

import (
	"strconv"

	"github.com/Baseledger/baseledger-bridge/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdAttestations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attestations [limit]",
		Short: "Query attestations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqLimit, _ := strconv.ParseUint(args[0], 10, 64)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAttestationsRequest{

				Limit: reqLimit,
			}

			res, err := queryClient.Attestations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
