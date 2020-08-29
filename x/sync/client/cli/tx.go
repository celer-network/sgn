package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/sync/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
)

// Change flags
const (
	FlagType   = "type"
	FlagData   = "data"
	FlagBlkNum = "blknum"
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

	syncTxCmd.AddCommand(common.PostCommands(
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
$ %s tx sync submit-change --type="sync_block" --data="My awesome change"
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txr, err := transactor.NewCliTransactor(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				log.Error(err)
				return err
			}

			changeType := viper.GetString(FlagType)
			data := viper.GetString(FlagData)
			blknum := viper.GetUint64(FlagBlkNum)
			msg := types.NewMsgSubmitChange(changeType, []byte(data), blknum, txr.Key.GetAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			txr.CliSendTxMsgWaitMined(msg)

			return nil
		},
	}

	cmd.Flags().String(FlagType, "", "type of change")
	cmd.Flags().String(FlagData, "", "data of change")
	cmd.Flags().String(FlagBlkNum, "", "mainchain block number of change")

	return cmd
}

// GetCmdVote implements creating a new approve command.
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "approve [change-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Vote for an active change",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a approve for an active change. You can
find the change-id by running "%s query sync changes".

Example:
$ %s tx sync approve 1
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txr, err := transactor.NewCliTransactor(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				log.Error(err)
				return err
			}

			// validate that the change id is a uint
			changeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("change-id %s not a valid int, please input a valid change-id", args[0])
			}

			// Build approve message and run basic validation
			msg := types.NewMsgApprove(changeID, txr.Key.GetAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			txr.CliSendTxMsgWaitMined(msg)

			return nil
		},
	}
}

// DONTCOVER
