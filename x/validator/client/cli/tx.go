package cli

import (
	"bufio"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/validator/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTransactors = "transactors"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	validatorTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Bridge transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	validatorTxCmd.AddCommand(flags.PostCommands(
		GetCmdSetTransactors(cdc),
		GetCmdWithdrawReward(cdc),
	)...)

	return validatorTxCmd
}

// GetCmdSetTransactors is the CLI command for sending a SetTransactors transaction
func GetCmdSetTransactors(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-transasctors [eth-addr]",
		Short: "set transasctors for the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info(viper.GetStringSlice(flagTransactors))
			transactors, err := transactor.ParseTransactorAddrs(viper.GetStringSlice(flagTransactors))
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(bufio.NewReader(cmd.InOrStdin())).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgSetTransactors(args[0], transactors, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringSlice(flagTransactors, []string{}, "transactors")

	return cmd
}

// GetCmdWithdrawReward is the CLI command for sending a WithdrawReward transaction
func GetCmdWithdrawReward(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw-reward [eth-addr]",
		Short: "withdraw reward for the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(bufio.NewReader(cmd.InOrStdin())).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgWithdrawReward(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
