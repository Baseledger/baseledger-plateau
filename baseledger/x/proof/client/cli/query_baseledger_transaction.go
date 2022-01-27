package cli

import (
	"context"
	"strconv"

	"github.com/Baseledger/baseledger/x/proof/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListBaseledgerTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-baseledger-transaction",
		Short: "list all baseledgerTransaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllBaseledgerTransactionRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.BaseledgerTransactionAll(context.Background(), params)
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

func CmdShowBaseledgerTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-baseledger-transaction [id]",
		Short: "shows a baseledgerTransaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetBaseledgerTransactionRequest{
				Id: id,
			}

			res, err := queryClient.BaseledgerTransaction(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
