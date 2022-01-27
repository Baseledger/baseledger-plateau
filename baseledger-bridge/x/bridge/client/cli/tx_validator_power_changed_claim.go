package cli

import (
	"strconv"

	"github.com/Baseledger/baseledger-bridge/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdValidatorPowerChangedClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-power-changed-claim [event-nonce] [block-height] [token-contract] [amount] [ethereum-sender] [cosmos-receiver]",
		Short: "Broadcast message validatorPowerChangedClaim",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argEventNonce, _ := strconv.ParseUint(args[0], 10, 64)
			argBlockHeight, _ := strconv.ParseUint(args[1], 10, 64)
			argTokenContract := args[2]
			argAmount, _ := sdk.NewIntFromString(args[3])
			argEthereumSender := args[4]
			argCosmosReceiver := args[5]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgValidatorPowerChangedClaim(
				clientCtx.GetFromAddress().String(),
				argEventNonce,
				argBlockHeight,
				argTokenContract,
				argAmount,
				argEthereumSender,
				argCosmosReceiver,
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
