package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagConfig  = "config"
	FlagCLIHome = "cli-home"

	FlagStartMonitor = "startMonitor"

	FlagEthGateway         = "eth.gateway"
	FlagEthCelrAddress     = "eth.contracts.celr_address"
	FlagEthDPoSAddress     = "eth.contracts.dpos_address"
	FlagEthSGNAddress      = "eth.contracts.sgn_address"
	FlagEthLedgerAddress   = "eth.contracts.ledger_address"
	FlagEthKeystore        = "eth.keystore"
	FlagEthPassphrase      = "eth.passphrase"
	FlagEthPollInterval    = "eth.poll_interval"
	FlagEthBlockDelay      = "eth.block_delay"
	FlagEthChainID         = "eth.chain_id"
	FlagEthCheckInterval   = "eth.check_interval"
	FlagEthMinGasPriceGwei = "eth.min_gas_price_gwei"
	FlagEthAddGasPriceGwei = "eth.add_gas_price_gwei"

	FlagSgnValidatorAccount = "sgn.validator_account"
	FlagSgnTransactors      = "sgn.transactors"
	FlagSgnPassphrase       = "sgn.passphrase"
	FlagSgnPubKey           = "sgn.pubkey"
	FlagSgnChainID          = "sgn.chain_id"
	FlagSgnNodeURI          = "sgn.node_uri"
	FlagSgnBaseGasPrice     = "sgn.base_gas_price"
	FlagSgnTimeoutCommit    = "sgn.timeout_commit"
	FlagSgnKeyringBackend   = "sgn.keyring_backend"
	FlagSgnGasAdjustment    = "sgn.gas_adjustment"
	FlagSgnExecuteSlash     = "sgn.execute_slash"

	FlagLogLevel = "log.level"
	FlagLogColor = "log.color"
)

const (
	DefaultGasAdjustment = 1.15
)

func PostCommands(cmds ...*cobra.Command) []*cobra.Command {
	for _, c := range cmds {
		c.Flags().Bool(flags.FlagIndentResponse, false, "Add indent to JSON response")
		c.Flags().Bool(flags.FlagTrustNode, true, "Trust connected full node (don't verify proofs for responses)")
		c.Flags().Bool(flags.FlagDryRun, false, "ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it")

		// --gas can accept integers and "simulate"
		c.Flags().Var(&flags.GasFlagVar, "gas", fmt.Sprintf(
			"gas limit to set per-transaction; set to %q to calculate required gas automatically (default %d)",
			flags.GasFlagAuto, flags.DefaultGasLimit,
		))
		viper.BindPFlag(flags.FlagTrustNode, c.Flags().Lookup(flags.FlagTrustNode))

		c.SetErr(c.ErrOrStderr())
	}
	return cmds
}

func GetCommands(cmds ...*cobra.Command) []*cobra.Command {
	for _, c := range cmds {
		c.Flags().Bool(flags.FlagIndentResponse, false, "Add indent to JSON response")
		c.Flags().Bool(flags.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
		c.Flags().Int64(flags.FlagHeight, 0, "Use a specific height to query state at (this can error if the node is pruning state)")

		viper.BindPFlag(flags.FlagTrustNode, c.Flags().Lookup(flags.FlagTrustNode))

		c.SetErr(c.ErrOrStderr())
	}
	return cmds
}
