package cli

import (
	"github.com/Baseledger/baseledger/x/proof/types"
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
