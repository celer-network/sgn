package common

import "time"

const (
	// outPathPrefix is the path prefix for all output from e2e (incl. chain data, binaries etc)
	// the code will append epoch second to this and create the folder
	// the folder will be deleted after test ends successfully
	OutRootDirPrefix = "/tmp/celer_e2e_"
	EnvDir           = "../../env"
	LocalGeth        = "http://127.0.0.1:8545"

	SgnChainID    = "sgntest"
	SgnPassphrase = "12341234"
	SgnGasPrice   = ""
	SgnCLIAddr    = "sgn15h2geedmud70gvpajdwpcaxfs4qcrw4z92zlqe"
	SgnNodeURI    = "tcp://localhost:26657"

	SgnBlockInterval = 1
	DefaultTimeout   = 60 * time.Second
	waitMinedTimeout = 180 * time.Second
	BlockDelay       = 0
	PollingInterval  = time.Second
	DisputeTimeout   = 100

	RetryPeriod = 200 * time.Millisecond
	RetryLimit  = 200
)

var (
	SgnCLIHomes    = [...]string{"../../../docker-volumes/node0/sgncli", "../../../docker-volumes/node1/sgncli", "../../../docker-volumes/node2/sgncli"}
	SgnOperators   = [...]string{"sgn1qehw7sn3u3nhnjeqk8kjccj263rq5fv002l5fk", "sgn1egtta7su5jxjahtw56pe07qerz4lwvrlttac6y", "sgn19q9usqmjcmx8vynynfl5tj5n2k22gc5f6wjvd7"}
	ValEthKs       = [...]string{"../../keys/ethks0.json", "../../keys/ethks1.json", "../../keys/ethks2.json"}
	ValEthAddrs    = [...]string{"6a6d2a97da1c453a4e099e8054865a0a59728863", "ba756d65a1a03f07d205749f35e2406e4a8522ad", "f25d8b54fad6e976eb9175659ae01481665a2254"}
	ClientEthKs    = [...]string{"../../keys/ethks3.json", "../../keys/ethks4.json"}
	ClientEthAddrs = [...]string{"cb409caa43e385288d6bff2c3a0635688c7b3294", "75e912af38888643829380fa9f2c4019f5710ff5"}
)
