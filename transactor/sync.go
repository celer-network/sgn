package transactor

import (
	"fmt"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagCandidate     = "candidator"
	FlagDelegatorAddr = "delegator"
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

	cmd.AddCommand(flags.PostCommands(
		GetSyncUpdateSidechainAddr(cdc),
		GetCmdSyncValidator(cdc),
		GetCmdSyncDelegator(cdc),
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
$ %s tx submit-change sync-update-sidechain-addr --candidator="0xf75f679d958b7610bad84e3baef2f9fa3e9bd961"
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			operator, err := NewOperator(cdc)
			if err != nil {
				return
			}

			candidate := viper.GetString(FlagCandidate)
			operator.SyncUpdateSidechainAddr(mainchain.Hex2Addr(candidate))
			return
		},
	}

	cmd.Flags().String(FlagCandidate, "", "Candidate address")
	cmd.MarkFlagRequired(FlagCandidate)

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
$ %s tx submit-change sync-validator --candidator="0xf75f679d958b7610bad84e3baef2f9fa3e9bd961"
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			operator, err := NewOperator(cdc)
			if err != nil {
				return
			}

			candidate := viper.GetString(FlagCandidate)
			operator.SyncValidator(mainchain.Hex2Addr(candidate))
			return
		},
	}

	cmd.Flags().String(FlagCandidate, "", "Candidate address")
	cmd.MarkFlagRequired(FlagCandidate)

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
$ %s tx submit-change sync-delegator --candidator="0xf75f679d958b7610bad84e3baef2f9fa3e9bd961" --delegator="0xf75f679d958b7610bad84e3baef2f9fa3e9bd961"
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			operator, err := NewOperator(cdc)
			if err != nil {
				return
			}

			candidate := viper.GetString(FlagCandidate)
			delegator := viper.GetString(FlagDelegatorAddr)
			operator.SyncDelegator(mainchain.Hex2Addr(candidate), mainchain.Hex2Addr(delegator))
			return
		},
	}

	cmd.Flags().String(FlagCandidate, "", "Candidate address")
	cmd.Flags().String(FlagDelegatorAddr, "", "Delegator address")
	cmd.MarkFlagRequired(FlagCandidate)
	cmd.MarkFlagRequired(FlagDelegatorAddr)

	return cmd
}
