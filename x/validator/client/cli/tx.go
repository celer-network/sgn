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
		// GetCmdInitializeCandidate(cdc),
		GetCmdUpdateSidechainAddr(cdc),
		GetCmdClaimValidator(cdc),
		GetCmdSyncValidator(cdc),
		GetCmdSyncDelegator(cdc),
		GetCmdWithdrawReward(cdc),
	)...)

	return validatorTxCmd
}

// GetCmdInitializeCandidate is the CLI command for sending a InitializeCandidate transaction
// func GetCmdInitializeCandidate(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "initialize-candidate [eth-addr]",
// 		Short: "initialize candidate for the eth address",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			txBldr := auth.NewTxBuilderFromCLI(bufio.NewReader(cmd.InOrStdin())).WithTxEncoder(utils.GetTxEncoder(cdc))
// 			msg := types.NewMsgInitializeCandidate(args[0], cliCtx.GetFromAddress())
// 			err := msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}

// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
// 		},
// 	}
// }

// GetCmdUpdateSidechainAddr is the CLI command for sending a UpdateSidechainAddr transaction
func GetCmdUpdateSidechainAddr(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update-sidechain-addr [eth-addr]",
		Short: "update sidechain address for the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(bufio.NewReader(cmd.InOrStdin())).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgUpdateSidechainAddr(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
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

// GetCmdClaimValidator is the CLI command for sending a ClaimValidator transaction
func GetCmdClaimValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-validator [eth-addr] [val-pubkey]",
		Short: "claim validator for the eth address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info(viper.GetStringSlice(flagTransactors))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(bufio.NewReader(cmd.InOrStdin())).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgClaimValidator(args[0], args[1], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdSyncValidator is the CLI command for sending a SyncValidator transaction
func GetCmdSyncValidator(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sync-validator [eth-addr]",
		Short: "sync validator for the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(bufio.NewReader(cmd.InOrStdin())).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgSyncValidator(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdSyncDelegator is the CLI command for sending a SyncDelegator transaction
func GetCmdSyncDelegator(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sync-delegator [candidate-addr] [delegator-addr]",
		Short: "sync delegator for the candidate address and delegator address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(bufio.NewReader(cmd.InOrStdin())).WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgSyncDelegator(args[0], args[1], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
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
