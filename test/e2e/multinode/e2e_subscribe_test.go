package multinode

import (
	"fmt"
	"math/big"
	"strings"
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
		CelrAddr:               tc.E2eProfile.CelrAddr,
		GovernProposalDeposit:  big.NewInt(1), // TODO: use a more practical value
		GovernVoteTimeout:      big.NewInt(1), // TODO: use a more practical value
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MaxValidatorNum:        big.NewInt(11),
		MinStakingPool:         big.NewInt(0),
		IncreaseRateWaitTime:   big.NewInt(1), // TODO: use a more practical value
		SidechainGoLiveTimeout: big.NewInt(0),
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
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnCLIAddr,
		tc.SgnPassphrase,
	)

	amt1 := big.NewInt(2000000000000000000)
	amt2 := big.NewInt(2000000000000000000)
	amt3 := big.NewInt(1000000000000000000)
	amts := []*big.Int{amt1, amt2, amt3}
	log.Infoln("Add validators...")
	tc.AddValidators(t, transactor, tc.ValEthKs[:], tc.SgnOperators[:], amts)
	turnOffMonitor(2)

	amt := new(big.Int)
	amt.SetString("1"+strings.Repeat("0", 20), 10)
	e2ecommon.SubscribteTestCommon(t, transactor, amt, "400000000000000000", 2)

	log.Infoln("Query sgn to check penalty")
	nonce := uint64(0)
	penalty, err := slash.CLIQueryPenalty(transactor.CliCtx, slash.StoreKey, nonce)
	tc.ChkTestErr(t, err, "failed to query penalty")
	expectedRes := fmt.Sprintf(`Nonce: %d, ValidatorAddr: %s, Reason: guard_failure`, nonce, tc.ValEthAddrs[2])
	assert.Equal(t, expectedRes, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	expectedRes = fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.ValEthAddrs[2])
	assert.Equal(t, expectedRes, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	assert.Equal(t, 2, len(penalty.Sigs), fmt.Sprintf("The length of validators should be 2"))

	log.Infoln("Query onchain staking pool")
	var poolAmt string
	for retry := 0; retry < 30; retry++ {
		ci, _ := tc.Client0.DPoS.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[2]))
		poolAmt = ci.StakingPool.String()
		if poolAmt == "990000000000000000" {
			break
		}
		time.Sleep(time.Second)
	}
	assert.Equal(t, "990000000000000000", poolAmt, fmt.Sprintf("The expected StakingPool should be 990000000000000000"))

}
