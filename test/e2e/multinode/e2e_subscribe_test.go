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
	tc.SleepWithLog(10, "sgn syncing")
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
	_, auth, err := tc.GetAuth(tc.ValEthKs[1])
	err = tc.DelegateStake(auth, mainchain.Hex2Addr(tc.ValEthAddrs[0]), amt3)
	tc.ChkTestErr(t, err, "failed to delegate stake")

	turnOffMonitor(2)

	amt := new(big.Int)
	amt.SetString("1"+strings.Repeat("0", 20), 10)
	// Request cost is 1000000000000000000, validator0 has a half of stake,
	// so it is going to get 500000000000000000 to distribute to its delegators.
	// validators0 commission rate is 0.01%, so the comission fee it collections is 50000000000000
	// The self delegated stake of validator0 is 2/3 of total stake of validator0,
	// so validator0 gets (500000000000000000 - 50000000000000) * 2/3 = 333300000000000000 reward.
	// The total service reward of validator0 is 333300000000000000 + 50000000000000 = 333350000000000000
	e2ecommon.SubscribteTestCommon(t, transactor, amt, "333350000000000000", 2)

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
	for retry := 0; retry < tc.RetryLimit; retry++ {
		ci, _ := tc.Client0.DPoS.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[2]))
		poolAmt = ci.StakingPool.String()
		if poolAmt == "990000000000000000" {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	assert.Equal(t, "990000000000000000", poolAmt, fmt.Sprintf("The expected StakingPool should be 990000000000000000"))

}
