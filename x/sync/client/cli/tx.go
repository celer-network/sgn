package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/celer-network/sgn/x/sync/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

// Change flags
const (
	FlagType = "type"
	FlagData = "data"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	syncTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Sync transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	syncTxCmd.AddCommand(flags.PostCommands(
		GetCmdSubmitChange(cdc),
		GetCmdVote(cdc),
	)...)

	return syncTxCmd
}

// GetCmdSubmitChange implements submitting a change transaction command.
func GetCmdSubmitChange(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-change",
		Short: "Submit a change",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a change along with type and data.

Example:
$ %s tx sync submit-change --type="sync_block" --data="My awesome change" --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			change, err := parseSubmitProposalFlags()
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(change.Deposit)
			if !ok {
				return err
			}

			content := types.ContentFromProposalType(change.Title, change.Description, change.Type)

			msg := types.NewMsgSubmitProposal(content, amount, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagType, "", "type of change")
	cmd.Flags().String(FlagData, "", "data of change")

	return cmd
}

// GetCmdVote implements creating a new vote command.
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "vote [change-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Vote for an active change",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a vote for an active change. You can
find the change-id by running "%s query sync changes".


Example:
$ %s tx sync vote 1 --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			// validate that the change id is a uint
			changeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("change-id %s not a valid int, please input a valid change-id", args[0])
			}

			// Build vote message and run basic validation
			msg := types.NewMsgVote(cliCtx.GetFromAddress(), changeID)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// DONTCOVER
