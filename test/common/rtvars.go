package testcommon

import (
	"github.com/celer-network/sgn/mainchain"
)

// runtime variables, will be initialized before each test
var (
	// E2eProfile will be updated and used for each test
	// not support parallel tests
	E2eProfile *TestProfile

	DefaultTestEthClient = &mainchain.EthClient{}
)
