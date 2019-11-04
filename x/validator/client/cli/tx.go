package cli

import (
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/x/validator/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
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

	validatorTxCmd.AddCommand(client.PostCommands(
		GetCmdInitializeCandidate(cdc),
		GetCmdClaimValidator(cdc),
		GetCmdSyncValidator(cdc),
		GetCmdSyncDelegator(cdc),
		GetCmdWithdrawReward(cdc),
	)...)

	return validatorTxCmd
}

// GetCmdInitializeCandidate is the CLI command for sending a InitializeCandidate transaction
func GetCmdInitializeCandidate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "initialize-candidate [eth-addr]",
		Short: "initialize candidate for the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgInitializeCandidate(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdClaimValidator is the CLI command for sending a SyncValidator transaction
func GetCmdClaimValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-validator [eth-addr] [val-pubkey]",
		Short: "claim validator for the eth address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			transactors, err := common.ParseTransactorAddrs(viper.GetStringSlice(flagTransactors))
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgClaimValidator(args[0], args[1], transactors, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringArray(flagTransactors, []string{}, "transactors")

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
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
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
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
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
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgWithdrawReward(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
