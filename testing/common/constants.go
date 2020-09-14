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
	ValAccounts    = [...]string{"sgn1qehw7sn3u3nhnjeqk8kjccj263rq5fv002l5fk", "sgn1egtta7su5jxjahtw56pe07qerz4lwvrlttac6y", "sgn19q9usqmjcmx8vynynfl5tj5n2k22gc5f6wjvd7"}
	ValEthKs       = [...]string{"../../keys/vethks0.json", "../../keys/vethks1.json", "../../keys/vethks2.json"}
	ValEthAddrs    = [...]string{"00078b31fa8b29a76bce074b5ea0d515a6aeaee7", "0015f5863ddc59ab6610d7b6d73b2eacd43e6b7e", "00290a43e5b2b151d530845b2d5a818240bc7c70"}
	ClientEthKs    = [...]string{"../../keys/cethks0.json", "../../keys/cethks1.json"}
	ClientEthAddrs = [...]string{"c06fdd796e140aee53de5111607e8ded93ebdca3", "c1699e89639adda8f39faefc0fc294ee5c3b462d"}

	// used by local manual tests
	SgnNodeURIs = [...]string{"tcp://localhost:26657", "tcp://localhost:26660", "tcp://localhost:26662"}
)
