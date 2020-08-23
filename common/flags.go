package common

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

	FlagValidatorAccount  = "sgn.validator_account"
	FlagSgnTransactors    = "sgn.transactors"
	FlagSgnPassphrase     = "sgn.passphrase"
	FlagSgnPubKey         = "sgn.pubkey"
	FlagSgnChainID        = "sgn.chain_id"
	FlagSgnNodeURI        = "sgn.node_uri"
	FlagSgnBaseGasPrice   = "sgn.base_gas_price"
	FlagSgnTimeoutCommit  = "sgn.timeout_commit"
	FlagSgnKeyringBackend = "sgn.keyring_backend"
	FlagSgnExecuteSlash   = "sgn.execute_slash"

	FlagLogLevel = "log.level"
	FlagLogColor = "log.color"
)
