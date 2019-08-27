package cli

import (
	"strconv"

	"github.com/celer-network/sgn/x/global/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	globalTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Bridge transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	globalTxCmd.AddCommand(client.PostCommands(
		GetCmdSyncBlock(cdc),
	)...)

	return globalTxCmd
}

// GetCmdSyncBlock is the CLI command for sending a SyncBlock transaction
func GetCmdSyncBlock(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sync-block [blockNum]",
		Short: "sync mainchain block with head",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			blockNumber, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgSyncBlock(blockNumber, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
