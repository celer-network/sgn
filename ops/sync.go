package ops

import (
	"fmt"
	"strings"
	"time"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagCandidateAddr = "candidate"
	FlagDelegatorAddr = "delegator"
	FlagConsumerAddr  = "consumer"
)

// GetSyncCmd
func GetSyncCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "sync",
		Short:                      "Sync a change from mainchain to sidechain",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(common.PostCommands(
		GetSyncUpdateSidechainAddr(cdc),
		GetCmdSyncValidator(cdc),
		GetCmdSyncDelegator(cdc),
		GetCmdSyncSubscriptionBalance(cdc),
	)...)

	return cmd
}

// GetSyncUpdateSidechainAddr
func GetSyncUpdateSidechainAddr(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-update-sidechain-addr",
		Short: "Sync sidechain address from mainchain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx submit-change sync-update-sidechain-addr --candidate="0xf75f679d958b7610bad84e3baef2f9fa3e9bd961"
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			operator, err := NewOperator(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				return
			}

			candidate := viper.GetString(FlagCandidateAddr)
			operator.SyncUpdateSidechainAddr(mainchain.Hex2Addr(candidate))
			time.Sleep(5 * time.Second)
			return
		},
	}

	cmd.Flags().String(FlagCandidateAddr, "", "Candidate address")
	cmd.MarkFlagRequired(FlagCandidateAddr)

	return cmd
}

// GetCmdSyncValidator
func GetCmdSyncValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-validator",
		Short: "Sync validator info from mainchain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx submit-change sync-validator --candidate="0xf75f679d958b7610bad84e3baef2f9fa3e9bd961"
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			operator, err := NewOperator(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				return
			}

			candidate := viper.GetString(FlagCandidateAddr)
			operator.SyncValidator(mainchain.Hex2Addr(candidate))
			time.Sleep(5 * time.Second)
			return
		},
	}

	cmd.Flags().String(FlagCandidateAddr, "", "Candidate address")
	cmd.MarkFlagRequired(FlagCandidateAddr)

	return cmd
}

// GetCmdSyncDelegator
func GetCmdSyncDelegator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-delegator",
		Short: "Sync delegator info from mainchain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx submit-change sync-delegator --candidate="0xf75f679d958b7610bad84e3baef2f9fa3e9bd961" --delegator="0xf75f679d958b7610bad84e3baef2f9fa3e9bd961"
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			operator, err := NewOperator(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				return
			}

			candidate := viper.GetString(FlagCandidateAddr)
			delegator := viper.GetString(FlagDelegatorAddr)
			operator.SyncDelegator(mainchain.Hex2Addr(candidate), mainchain.Hex2Addr(delegator))
			time.Sleep(5 * time.Second)
			return
		},
	}

	cmd.Flags().String(FlagCandidateAddr, "", "Candidate address")
	cmd.Flags().String(FlagDelegatorAddr, "", "Delegator address")
	cmd.MarkFlagRequired(FlagCandidateAddr)
	cmd.MarkFlagRequired(FlagDelegatorAddr)

	return cmd
}

// GetCmdSyncSubscriptionBalance
func GetCmdSyncSubscriptionBalance(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-subscription-balance",
		Short: "Sync subscription balance info from mainchain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx submit-change sync-subscription-balance --consumer="0xf75f679d958b7610bad84e3baef2f9fa3e9bd961"
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			operator, err := NewOperator(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				return
			}

			consumer := viper.GetString(FlagConsumerAddr)
			consumerAddr := mainchain.Hex2Addr(consumer)
			deposit, err := operator.EthClient.SGN.SubscriptionDeposits(
				&bind.CallOpts{}, consumerAddr)
			if err != nil {
				return
			}

			operator.SyncSubscriptionBalance(consumerAddr, deposit)
			time.Sleep(5 * time.Second)
			return
		},
	}

	cmd.Flags().String(FlagConsumerAddr, "", "Consumer address")
	cmd.MarkFlagRequired(FlagConsumerAddr)

	return cmd
}
