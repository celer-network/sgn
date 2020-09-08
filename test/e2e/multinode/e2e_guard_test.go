package multinode

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	e2ecommon "github.com/celer-network/sgn/test/e2e/common"
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
		GovernProposalDeposit:  big.NewInt(1), // TODO: use a more practical value
		GovernVoteTimeout:      big.NewInt(1), // TODO: use a more practical value
		SlashTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MaxValidatorNum:        big.NewInt(11),
		MinStakingPool:         big.NewInt(0),
		AdvanceNoticePeriod:    big.NewInt(1), // TODO: use a more practical value
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

	amt1 := big.NewInt(2000000000000000000)
	amt2 := big.NewInt(2000000000000000000)
	amt3 := big.NewInt(1000000000000000000)
	amts := []*big.Int{amt1, amt2, amt3}
	log.Infoln("Add validators...")
	tc.AddValidators(t, transactor, tc.ValEthKs[:], tc.ValAccounts[:], amts)
	_, auth, err := tc.GetAuth(tc.ValEthKs[1])
	err = tc.DelegateStake(auth, mainchain.Hex2Addr(tc.ValEthAddrs[0]), amt3)
	require.NoError(t, err, "failed to delegate stake")

	turnOffMonitor(2)

	amt := new(big.Int)
	amt.SetString("1"+strings.Repeat("0", 20), 10)
	// Request cost is 1000000000000000000 * 2, validator0 has a half of stake,
	// so it is going to get 1000000000000000000 to distribute to its delegators.
	// validators0 commission rate is 0.01%, so the comission fee it collections is 100000000000000
	// The self delegated stake of validator0 is 2/3 of total stake of validator0,
	// so validator0 gets (1000000000000000000 - 100000000000000) * 2/3 = 666600000000000000 reward.
	// The total service reward of validator0 is 666600000000000000 + 100000000000000 = 666700000000000000
	e2ecommon.GuardTestCommon(t, transactor, amt, "666700000000000000", 2)

	log.Infoln("Query sgn to check penalty")
	nonce := uint64(0)
	penalty, err := slash.CLIQueryPenalty(transactor.CliCtx, slash.StoreKey, nonce)
	require.NoError(t, err, "failed to query penalty")
	expectedRes := fmt.Sprintf(`Nonce: %d, ValidatorAddr: %s, Reason: guard_failure`, nonce, tc.ValEthAddrs[2])
	assert.Equal(t, expectedRes, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	expectedRes = fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.ValEthAddrs[2])
	assert.Equal(t, expectedRes, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	assert.Equal(t, 2, len(penalty.Sigs), fmt.Sprintf("The length of validators should be 2"))

	log.Infoln("Query onchain staking pool")
	var poolAmt string
	for retry := 0; retry < tc.RetryLimit; retry++ {
		ci, _ := tc.DposContract.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[2]))
		poolAmt = ci.StakingPool.String()
		if poolAmt == "990000000000000000" {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	assert.Equal(t, "990000000000000000", poolAmt, fmt.Sprintf("The expected StakingPool should be 990000000000000000"))

}
