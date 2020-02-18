package cli

import (
	"bufio"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/subscribe/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	subscribeTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "subscribe transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	subscribeTxCmd.AddCommand(flags.PostCommands(
		GetCmdSubscribe(cdc),
		GetCmdRequestGuard(cdc),
	)...)

	return subscribeTxCmd
}

// GetCmdSubscribe is the CLI command for sending a Subscribe transaction
func GetCmdSubscribe(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "subscribe [eth-addr]",
		Short: "set subscription info associated with the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(bufio.NewReader(cmd.InOrStdin())).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgSubscribe(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdRequestGuard is the CLI command for sending a request guard transaction
func GetCmdRequestGuard(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "request-guard [eth-addr] [signed-simplex-state]",
		Short: "request guard on a channel",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(bufio.NewReader(cmd.InOrStdin())).WithTxEncoder(utils.GetTxEncoder(cdc))
			signedSimplexStateBytes := mainchain.Hex2Bytes(args[1])
			msg := types.NewMsgRequestGuard(args[0], signedSimplexStateBytes, cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
