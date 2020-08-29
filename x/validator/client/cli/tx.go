package cli

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/validator/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTransactors     = "transactors"
	flagMoniker         = "moniker"
	flagIdentity        = "identity"
	flagWebsite         = "website"
	flagSecurityContact = "security-contact"
	flagDetails         = "details"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	validatorTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Validator transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	validatorTxCmd.AddCommand(common.PostCommands(
		GetCmdSetTransactors(cdc),
		GetCmdEditCandidateDescription(cdc),
		GetCmdWithdrawReward(cdc),
	)...)

	return validatorTxCmd
}

// GetCmdSetTransactors is the CLI command for sending a SetTransactors transaction
func GetCmdSetTransactors(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-transactors [eth-addr]",
		Short: "set transactors for the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info(viper.GetStringSlice(flagTransactors))
			transactors, err := common.ParseTransactorAddrs(viper.GetStringSlice(flagTransactors))
			if err != nil {
				return err
			}

			txr, err := transactor.NewCliTransactor(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				log.Error(err)
				return err
			}

			msg := types.NewMsgSetTransactors(args[0], transactors, txr.Key.GetAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			txr.CliSendTxMsgWaitMined(msg)

			return nil
		},
	}

	cmd.Flags().StringSlice(flagTransactors, []string{}, "transactors")

	return cmd
}

// GetCmdEditCandidateDescription is the CLI command for sending a EditCandidateDescription transaction
func GetCmdEditCandidateDescription(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-candidate-description [eth-addr]",
		Short: "Edit candidate description for the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			moniker, _ := cmd.Flags().GetString(flagMoniker)
			identity, _ := cmd.Flags().GetString(flagIdentity)
			website, _ := cmd.Flags().GetString(flagWebsite)
			security, _ := cmd.Flags().GetString(flagSecurityContact)
			details, _ := cmd.Flags().GetString(flagDetails)
			description := staking.NewDescription(moniker, identity, website, security, details)

			txr, err := transactor.NewCliTransactor(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				log.Error(err)
				return err
			}

			msg := types.NewMsgEditCandidateDescription(args[0], description, txr.Key.GetAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			txr.CliSendTxMsgWaitMined(msg)

			return nil
		},
	}

	cmd.Flags().String(flagMoniker, staking.DoNotModifyDesc, "The candidate's name")
	cmd.Flags().String(flagIdentity, staking.DoNotModifyDesc, "The (optional) identity signature (ex. UPort or Keybase)")
	cmd.Flags().String(flagWebsite, staking.DoNotModifyDesc, "The candidate's (optional) website")
	cmd.Flags().String(flagSecurityContact, staking.DoNotModifyDesc, "The candidate's (optional) security contact email")
	cmd.Flags().String(flagDetails, staking.DoNotModifyDesc, "The candidate's (optional) details")

	return cmd
}

// GetCmdWithdrawReward is the CLI command for sending a WithdrawReward transaction
func GetCmdWithdrawReward(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw-reward [eth-addr]",
		Short: "withdraw reward for the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txr, err := transactor.NewCliTransactor(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				log.Error(err)
				return err
			}

			msg := types.NewMsgWithdrawReward(args[0], txr.Key.GetAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			txr.CliSendTxMsgWaitMined(msg)

			return nil
		},
	}
}
