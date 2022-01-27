package cli

import (
	"strconv"

	"github.com/Baseledger/baseledger/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdDelegateKeysByEthAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-keys-by-eth-address [eth-address]",
		Short: "Query delegateKeysByEthAddress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqEthAddress := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryDelegateKeysByEthAddressRequest{

				EthAddress: reqEthAddress,
			}

			res, err := queryClient.DelegateKeysByEthAddress(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
