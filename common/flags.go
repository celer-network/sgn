package common

const (
	FlagConfig  = "config"
	FlagCLIHome = "cli-home"

	FlagStartMonitor = "startMonitor"

	FlagEthGateway         = "eth.gateway"
	FlagEthCelrAddress     = "eth.contracts.celrAddress"
	FlagEthDPoSAddress     = "eth.contracts.dposAddress"
	FlagEthSGNAddress      = "eth.contracts.sgnAddress"
	FlagEthLedgerAddress   = "eth.contracts.ledgerAddress"
	FlagEthKeystore        = "eth.keystore"
	FlagEthPassphrase      = "eth.passphrase"
	FlagEthPollInterval    = "eth.pollInterval"
	FlagEthBlockDelay      = "eth.blockDelay"
	FlagEthChainID         = "eth.chainID"
	FlagEthCheckInterval   = "eth.checkInterval"
	FlagEthMinGasPriceGwei = "eth.minGasPriceGwei"
	FlagEthAddGasPriceGwei = "eth.addGasPriceGwei"

	FlagSgnOperator       = "sgn.operator"
	FlagSgnTransactors    = "sgn.transactors"
	FlagSgnPassphrase     = "sgn.passphrase"
	FlagSgnPubKey         = "sgn.pubkey"
	FlagSgnChainID        = "sgn.chainID"
	FlagSgnNodeURI        = "sgn.nodeURI"
	FlagSgnBaseGasPrice   = "sgn.baseGasPrice"
	FlagSgnTimeoutCommit  = "sgn.timeoutCommit"
	FlagSgnKeyringBackend = "sgn.keyringBackend"

	FlagLogLevel = "log.level"
	FlagLogColor = "log.color"
)
