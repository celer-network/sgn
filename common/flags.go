package common

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagConfig  = "config"
	FlagCLIHome = "cli-home"

	FlagEthGateway           = "eth.gateway"
	FlagEthCelrAddress       = "eth.contracts.celr_address"
	FlagEthDPoSAddress       = "eth.contracts.dpos_address"
	FlagEthSGNAddress        = "eth.contracts.sgn_address"
	FlagEthLedgerAddress     = "eth.contracts.ledger_address"
	FlagEthKeystore          = "eth.keystore"
	FlagEthPassphrase        = "eth.passphrase"
	FlagEthPollInterval      = "eth.poll_interval"
	FlagEthBlockDelay        = "eth.block_delay"
	FlagEthChainID           = "eth.chain_id"
	FlagEthCheckInterval     = "eth.check_interval"
	FlagEthMonitorStartBlock = "eth.monitor_start_block"
	FlagEthMinGasPriceGwei   = "eth.min_gas_price_gwei"
	FlagEthAddGasPriceGwei   = "eth.add_gas_price_gwei"

	FlagSgnValidatorAccount = "sgn.validator_account"
	FlagSgnTransactors      = "sgn.transactors"
	FlagSgnPassphrase       = "sgn.passphrase"
	FlagSgnPubKey           = "sgn.pubkey"
	FlagSgnChainID          = "sgn.chain_id"
	FlagSgnNodeURI          = "sgn.node_uri"
	FlagSgnTimeoutCommit    = "sgn.timeout_commit"
	FlagSgnKeyringBackend   = "sgn.keyring_backend"
	FlagSgnMinGasPrices     = "sgn.min_gas_prices"

	FlagSgnCheckIntervalGuardQueue = "sgn.check_interval.guard_queue"
	FlagSgnCheckIntervalSlashQueue = "sgn.check_interval.slash_queue"

	FlagLogLevel = "log.level"
	FlagLogColor = "log.color"
)

func PostCommands(cmds ...*cobra.Command) []*cobra.Command {
	for _, c := range cmds {
		c.Flags().Bool(flags.FlagIndentResponse, false, "Add indent to JSON response")
		c.Flags().Bool(flags.FlagTrustNode, true, "Trust connected full node (don't verify proofs for responses)")

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
