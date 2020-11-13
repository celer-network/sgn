package multinode

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tec "github.com/celer-network/sgn/test/e2e/common"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/celer-network/sgn/x/slash"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGuard() {
	log.Infoln("set up new sgn env")
	p := &tc.SGNParams{
		CelrAddr:               tc.E2eProfile.CelrAddr,
		GovernProposalDeposit:  big.NewInt(1),
		GovernVoteTimeout:      big.NewInt(1),
		SlashTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MaxValidatorNum:        big.NewInt(11),
		MinStakingPool:         big.NewInt(0),
		AdvanceNoticePeriod:    big.NewInt(1),
		SidechainGoLiveTimeout: big.NewInt(0),
	}
	tc.SetupNewSGNEnv(p, false)
	tc.SleepWithLog(10, "sgn syncing")
}

func TestE2EGuard(t *testing.T) {
	setupGuard()

	t.Run("e2e-guard", func(t *testing.T) {
		t.Run("guardTest", guardTest)
	})
}

func guardTest(t *testing.T) {
	log.Infoln("===================================================================")
	log.Infoln("======================== Test guard ===========================")

	transactor := tc.NewTestTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnCLIAddr,
		tc.SgnPassphrase,
	)

	amts := []*big.Int{
		big.NewInt(9000000000000000000), // 9 CELR
		big.NewInt(9000000000000000000), // 9 CELR
		big.NewInt(7000000000000000000), // 7 CELR
		big.NewInt(6000000000000000000), // 6 CELR
	}
	log.Infoln("Add validators...")
	tc.AddValidators(t, transactor, tc.ValEthKs[:], tc.ValAccounts[:], amts)
	log.Infoln("Additionally delegate 1 CELR to validator 0...")
	_, dAuth, err := tc.GetAuth(tc.DelEthKs[0])
	require.NoError(t, err, "failed to get delegator auth")
	err = tc.DelegateStake(dAuth, mainchain.Hex2Addr(tc.ValEthAddrs[0]), big.NewInt(1000000000000000000)) // 1 CELR
	require.NoError(t, err, "failed to delegate stake")

	amt := new(big.Int)
	amt.SetString("1"+strings.Repeat("0", 20), 10)
	tec.Subscribe(t, transactor, amt)

	restartWithConfig(0, common.FlagSgnCheckIntervalGuardQueue, 10000)
	// validator 0 will fail to guard as the queue check interval was set to 10000, so validator 1
	// will send the guard tx, and validator 0 will be slashed
	tec.TestGuard(t, transactor, tc.ValEthAddrs[1])

	/* Request cost is 1000000000000000000 * 2, validator0 has a 10/32 of stake,
	so it is going to get 625000000000000000 to distribute to its delegators.
	validators0 commission rate is 0.01%, so the comission fee it collections is 62500000000000
	The self delegated stake of validator0 is 9/10 of total stake of validator0,
	so validator0 gets (625000000000000000 - 62500000000000) * 9/10 = 562443750000000000 reward.
	The total service reward of validator0 is 562443750000000000 + 62500000000000 = 562506250000000000 */
	tec.CheckReward(t, transactor, tc.ValEthAddrs[0], "562506250000000000", 4)

	log.Infoln("Query sgn to check penalty")
	nonce := uint64(0)
	penalty, err := slash.CLIQueryPenalty(transactor.CliCtx, slash.StoreKey, nonce)
	require.NoError(t, err, "failed to query penalty")
	expectedRes := fmt.Sprintf(`Nonce: %d, ValidatorAddr: %s, Reason: guard_failure`, nonce, tc.ValEthAddrs[0])
	assert.Equal(t, expectedRes, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	expectedRes = fmt.Sprintf(`Account: %s, Amount: 90000000000000000`, tc.ValEthAddrs[0])
	assert.Equal(t, expectedRes, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	expectedRes = fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.DelEthAddrs[0])
	assert.Equal(t, expectedRes, penalty.PenalizedDelegators[1].String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	assert.Equal(t, 4, len(penalty.Sigs), fmt.Sprintf("The length of validators should be 4"))

	log.Infoln("Query onchain staking pool")
	var poolAmt string
	expPoolAmt := "9900000000000000000"
	for retry := 0; retry < tc.RetryLimit; retry++ {
		ci, _ := tc.DposContract.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[0]))
		poolAmt = ci.StakingPool.String()
		if poolAmt == expPoolAmt {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	assert.Equal(t, expPoolAmt, poolAmt, fmt.Sprintf("The expected StakingPool should be %s", expPoolAmt))

	// 2nd channel guard test, first guardian should be assigned to validator 3
	tec.TestGuard(t, transactor, tc.ValEthAddrs[3])
}
