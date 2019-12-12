package testing

import "github.com/celer-network/sgn/mainchain"

// runtime variables, will be initialized before each test
var (
	// used by setup_mainchain and tests
	Client0Addr = mainchain.Hex2Addr(Client0AddrStr)
	Client1Addr = mainchain.Hex2Addr(Client1AddrStr)

	// e2eProfile will be updated and used for each test
	// not support parallel tests
	E2eProfile *TestProfile
)
