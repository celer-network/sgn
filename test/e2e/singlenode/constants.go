package singlenode

import "time"

const (
	// outPathPrefix is the path prefix for all output from e2e (incl. chain data, binaries etc)
	// the code will append epoch second to this and create the folder
	// the folder will be deleted after test ends successfully
	outRootDirPrefix = "/tmp/celer_e2e_"

	// etherbase and osp addr/priv key in hex
	client0AddrStr    = "6a6d2a97da1c453a4e099e8054865a0a59728863"
	client0SGNAddrStr = "cosmos1ddvpnk98da5hgzz8lf5y82gnsrhvu3jd3cukpp"
	client1AddrStr    = "ba756d65a1a03f07d205749f35e2406e4a8522ad"
	client1Priv       = "c2ff7d4ce25f7448de00e21bbbb7b884bb8dc0ca642031642863e78a35cb933d"

	blockDelay       = 3
	sgnBlockInterval = 1
	defaultTimeout   = 60 * time.Second
)
