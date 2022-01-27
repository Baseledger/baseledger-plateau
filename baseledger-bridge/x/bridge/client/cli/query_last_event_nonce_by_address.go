package cli

import (
	"strconv"

	"github.com/Baseledger/baseledger-bridge/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdLastEventNonceByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-event-nonce-by-address [address]",
		Short: "Query lastEventNonceByAddress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryLastEventNonceByAddressRequest{

				Address: reqAddress,
			}

			res, err := queryClient.LastEventNonceByAddress(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
