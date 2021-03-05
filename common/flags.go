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
	FlagEthKeystore          = "eth.keystore"
	FlagEthPassphrase        = "eth.passphrase"
	FlagEthPollInterval      = "eth.poll_interval"
	FlagEthSyncBlkInterval   = "eth.sync_blk_interval"
	FlagEthAcceptedBlkRange  = "eth.accepted_blk_range"
	FlagEthBlockDelay        = "eth.block_delay"
	FlagEthChainID           = "eth.chain_id"
	FlagEthCheckInterval     = "eth.check_interval"
	FlagEthMonitorStartBlock = "eth.monitor_start_block"
	FlagEthMinGasPriceGwei   = "eth.min_gas_price_gwei"
	FlagEthAddGasPriceGwei   = "eth.add_gas_price_gwei"

	FlagEthTxEnableTxSender           = "eth.tx.enable_tx_sender"
	FlagEthTxDBPollInterval           = "eth.tx.db_poll_interval"
	FlagEthTxBlockTime                = "eth.tx.block_time"
	FlagEthTxKeyDir                   = "eth.tx.key_dir"
	FlagEthTxSenderAddress            = "eth.tx.sender_address"
	FlagEthTxRPCURL                   = "eth.tx.rpc_url"
	FlagEthTxSecondaryRPCURLs         = "eth.tx.secondary_rpc_urls"
	FlagEthTxHeadTrackerHistoryDepth  = "eth.tx.head_tracker_history_depth"
	FlagEthTxHeadTrackerMaxBufferSize = "eth.tx.head_tracker_max_buffer_size"
	FlagEthTxHeadTrackerFinalityDepth = "eth.tx.head_tracker_finality_depth"
	FlagEthTxDefaultGasPrice          = "eth.tx.default_gas_price"
	FlagEthTxMaxGasPrice              = "eth.tx.max_gas_price"
	FlagEthTxGasBumpPercent           = "eth.tx.gas_bump_percent"
	FlagEthTxGasBumpWei               = "eth.tx.gas_bump_wei"
	FlagEthTxGasBumpThreshold         = "eth.tx.gas_bump_threshold"
	FlagEthTxGasBumpTxDepth           = "eth.tx.gas_bump_tx_depth"

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
