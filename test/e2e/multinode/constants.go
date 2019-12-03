package multinode

import "time"

const (
	// outPathPrefix is the path prefix for all output from e2e (incl. chain data, binaries etc)
	// the code will append epoch second to this and create the folder
	// the folder will be deleted after test ends successfully
	outRootDirPrefix = "/tmp/celer_e2e_"

	// etherbase and osp addr/priv key in hex
	etherBaseAddrStr  = "b5bb8b7f6f1883e0c01ffb8697024532e6f3238c"
	etherBasePriv     = "69ef4da8204644e354d759ca93b94361474259f63caac6e12d7d0abcca0063f8"
	client0AddrStr    = "6a6d2a97da1c453a4e099e8054865a0a59728863"
	client0Priv       = "a7c9fa8bcd45a86fdb5f30fecf88337f20185b0c526088f2b8e0f726cad12857"
	client0SGNAddrStr = "cosmos1ddvpnk98da5hgzz8lf5y82gnsrhvu3jd3cukpp"
	client1AddrStr    = "ba756d65a1a03f07d205749f35e2406e4a8522ad"
	client1Priv       = "c2ff7d4ce25f7448de00e21bbbb7b884bb8dc0ca642031642863e78a35cb933d"

	maxBlockDiff     = 2 // defined in sidechain's genesis file
	blockDelay       = 2
	sgnBlockInterval = 1
	defaultTimeout   = 60 * time.Second
)
