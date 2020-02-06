package multinode

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/test/common"
	e2ecommon "github.com/celer-network/sgn/test/e2e/common"
	"github.com/celer-network/sgn/x/slash"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
)

func setUpSubscribe() {
	log.Infoln("set up new sgn env")
	p := &tc.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MinStakingPool:         big.NewInt(0),
		SidechainGoLiveTimeout: big.NewInt(0),
		CelrAddr:               tc.E2eProfile.CelrAddr,
		MaxValidatorNum:        big.NewInt(11),
	}
	setupNewSGNEnv(p)
}

func TestE2ESubscribe(t *testing.T) {
	setUpSubscribe()

	t.Run("e2e-subscribe", func(t *testing.T) {
		t.Run("subscribeTest", subscribeTest)
	})
}

func subscribeTest(t *testing.T) {
	log.Infoln("===================================================================")
	log.Infoln("======================== Test subscribe ===========================")

	transactor := tc.NewTransactor(
		t,
		tc.SgnCLIHome,
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnCLIAddr,
		tc.SgnPassphrase,
		tc.SgnGasPrice,
	)

	amt := new(big.Int)
	amt.SetString("100000000000000000000", 10) // 100 CELR
	amts := []*big.Int{amt, amt, amt}
	log.Infoln("Add validators...")
	tc.AddValidators(t, transactor, tc.ValEthKs[:], tc.SgnOperators[:], amts)
	turnOffMonitor(2)

	e2ecommon.SubscribteTestCommon(t, transactor, amt, "476190476190476190", 2)

	log.Infoln("Query sgn to check penalty")
	nonce := uint64(0)
	penalty, err := slash.CLIQueryPenalty(transactor.CliCtx, slash.StoreKey, nonce)
	tc.ChkTestErr(t, err, "failed to query penalty")
	expectedRes := fmt.Sprintf(`Nonce: %d, ValidatorAddr: %s, Reason: guard_failure`, nonce, tc.ValEthAddrs[2])
	assert.Equal(t, expectedRes, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	expectedRes = fmt.Sprintf(`Account: %s, Amount: 1000000000000000`, tc.ValEthAddrs[2])
	assert.Equal(t, expectedRes, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	assert.Equal(t, 2, len(penalty.Sigs), fmt.Sprintf("The length of validators should be 2"))

	log.Infoln("Query onchain staking pool")
	var poolAmt string
	for retry := 0; retry < 30; retry++ {
		ci, _ := tc.Client0.Guard.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[2]))
		poolAmt = ci.StakingPool.String()
		if poolAmt == "99000000000000000" {
			break
		}
		time.Sleep(time.Second)
	}
	assert.Equal(t, "99000000000000000", poolAmt, fmt.Sprintf("The expected StakingPool should be 99000000000000000"))

}
