package multinode

import (
	"math/big"
	"testing"

	"github.com/celer-network/goutils/log"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/slash"
	"github.com/cosmos/cosmos-sdk/client/rpc"
)

func setUpSlash() {
	log.Infoln("set up new sgn env")
	p := &tf.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MinStakingPool:         big.NewInt(0),
		SidechainGoLiveTimeout: big.NewInt(0),
		CelrAddr:               tf.E2eProfile.CelrAddr,
	}
	setupNewSGNEnv(p)
	amts := []*big.Int{big.NewInt(1000000000000000000), big.NewInt(1000000000000000000), big.NewInt(100000000000000000)}
	addValidators(ethKeystores[:], ethKeystorePps[:], sgnOperators[:], amts)
	tf.SleepWithLog(10, "sgn syncing")
}

func TestE2ESlash(t *testing.T) {
	setUpSlash()

	t.Run("e2e-slash", func(t *testing.T) {
		t.Run("slashTest", slashTest)
	})
}

func slashTest(t *testing.T) {
	log.Infoln("===================================================================")
	log.Infoln("======================== Test slash ===========================")

	transactor := tf.NewTransactor(
		t,
		sgnCLIHomes[0],
		sgnChainID,
		sgnNodeURIs[0],
		sgnTransactors[0],
		sgnPassphrase,
		sgnGasPrice,
	)

	shutdownNode("sgnnode2")
	tf.SleepWithLog(10, "shutdown node2")

	for {
		height, _ := rpc.GetChainHeight(transactor.CliCtx)
		log.Infoln("latest height", height)
		penalty, err := slash.CLIQueryPenalty(transactor.CliCtx, slash.StoreKey, 1)
		log.Errorln(penalty, err)
		tf.SleepWithLog(10, "check height")
	}
}
