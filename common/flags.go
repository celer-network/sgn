package common

const (
	FlagConfig  = "config"
	FlagCLIHome = "cli-home"

	FlagStartMonitor = "startMonitor"

	FlagEthGateway       = "eth.gateway"
	FlagEthDPoSAddress   = "eth.contracts.dposAddress"
	FlagEthSGNAddress    = "eth.contracts.sgnAddress"
	FlagEthLedgerAddress = "eth.contracts.ledgerAddress"
	FlagEthKeystore      = "eth.keystore"
	FlagEthPassphrase    = "eth.passphrase"
	FlagEthPollInterval  = "eth.pollInterval"
	FlagEthBlockDelay    = "eth.blockDelay"
	FlagEthChainID       = "eth.chainID"
	FlagEthCheckInterval = "eth.checkInterval"

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
