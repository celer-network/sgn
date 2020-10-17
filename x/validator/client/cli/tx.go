package cli

import (
	"io/ioutil"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
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
	flagMoniker  = "moniker"
	flagIdentity = "identity"
	flagWebsite  = "website"
	flagContact  = "contact"
	flagDetails  = "details"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	validatorTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Validator transaction subcommands",
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
		Use:   "set-transactors",
		Short: "set transactors based on transactors in config",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			transactors, err := common.ParseTransactorAddrs(viper.GetStringSlice(common.FlagSgnTransactors))
			if err != nil {
				return err
			}

			txr, err := transactor.NewCliTransactor(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				return err
			}

			msg := types.NewMsgSetTransactors(transactors, txr.Key.GetAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			txr.CliSendTxMsgWaitMined(msg)

			return nil
		},
	}

	return cmd
}

// GetCmdEditCandidateDescription is the CLI command for sending a EditCandidateDescription transaction
func GetCmdEditCandidateDescription(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-candidate-description",
		Short: "Edit candidate description",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			moniker, _ := cmd.Flags().GetString(flagMoniker)
			identity, _ := cmd.Flags().GetString(flagIdentity)
			website, _ := cmd.Flags().GetString(flagWebsite)
			contact, _ := cmd.Flags().GetString(flagContact)
			details, _ := cmd.Flags().GetString(flagDetails)
			description := staking.NewDescription(moniker, identity, website, contact, details)

			txr, err := transactor.NewCliTransactor(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				return err
			}

			ksBytes, err := ioutil.ReadFile(viper.GetString(common.FlagEthKeystore))
			if err != nil {
				return err
			}

			address, err := mainchain.GetAddressFromKeystore(ksBytes)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditCandidateDescription(address, description, txr.Key.GetAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			txr.CliSendTxMsgWaitMined(msg)

			return nil
		},
	}

	cmd.Flags().String(flagMoniker, staking.DoNotModifyDesc, "The candidate's name")
	cmd.Flags().String(flagIdentity, staking.DoNotModifyDesc, "The identity signature (ex. UPort or Keybase)")
	cmd.Flags().String(flagWebsite, staking.DoNotModifyDesc, "The candidate's website")
	cmd.Flags().String(flagContact, staking.DoNotModifyDesc, "The candidate's security contact email")
	cmd.Flags().String(flagDetails, staking.DoNotModifyDesc, "The candidate's details")

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
