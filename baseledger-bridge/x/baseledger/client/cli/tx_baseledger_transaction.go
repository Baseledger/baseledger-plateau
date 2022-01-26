package cli

import (
	"strconv"

	"github.com/Baseledger/baseledger-bridge/x/baseledger/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdCreateBaseledgerTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-baseledger-transaction [baseledger-transaction-id] [payload] [op-code]",
		Short: "Create a new baseledgerTransaction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBaseledgerTransactionId := args[0]
			argPayload := args[1]
			argOpCode, err := cast.ToUint32E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateBaseledgerTransaction(clientCtx.GetFromAddress().String(), argBaseledgerTransactionId, argPayload, argOpCode)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateBaseledgerTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-baseledger-transaction [id] [baseledger-transaction-id] [payload] [op-code]",
		Short: "Update a baseledgerTransaction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			argBaseledgerTransactionId := args[1]

			argPayload := args[2]

			argOpCode, err := cast.ToUint32E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateBaseledgerTransaction(clientCtx.GetFromAddress().String(), id, argBaseledgerTransactionId, argPayload, argOpCode)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteBaseledgerTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-baseledger-transaction [id]",
		Short: "Delete a baseledgerTransaction by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteBaseledgerTransaction(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
