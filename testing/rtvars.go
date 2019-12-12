package testing

import (
	"sync"

	"github.com/celer-network/sgn/mainchain"
)

// runtime variables, will be initialized before each test
var (
	// used by setup_mainchain and tests
	Client0Addr = mainchain.Hex2Addr(Client0AddrStr)
	Client1Addr = mainchain.Hex2Addr(Client1AddrStr)

	// E2eProfile will be updated and used for each test
	// not support parallel tests
	E2eProfile *TestProfile

	EthClient        = &mainchain.EthClient{}
	pendingNonceLock sync.Mutex
)
