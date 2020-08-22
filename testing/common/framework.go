package common

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
)

// runtime variables, will be initialized before each test
var (
	// E2eProfile will be updated and used for each test
	// not support parallel tests
	E2eProfile *TestProfile
)

// Killable is object that has Kill() func
type Killable interface {
	Kill() error
}

type TestProfile struct {
	DisputeTimeout uint64
	LedgerAddr     mainchain.Addr
	DPoSAddr       mainchain.Addr
	SGNAddr        mainchain.Addr
	CelrAddr       mainchain.Addr
	CelrContract   *mainchain.ERC20
}

func TearDown(tokill []Killable) {
	log.Info("Tear down Killables ing...")
	for _, p := range tokill {
		ChkErr(p.Kill(), "kill process error")
	}
}

func ChkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
