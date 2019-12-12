package testing

import "github.com/celer-network/sgn/mainchain"

var (
	// used by setup_mainchain and tests
	Client0Addr = mainchain.Hex2Addr(Client0AddrStr)
	Client1Addr = mainchain.Hex2Addr(Client1AddrStr)
)
