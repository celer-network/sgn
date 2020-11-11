package multinode

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSlash() {
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

func TestE2ESlash(t *testing.T) {
	setupSlash()

	t.Run("e2e-slash", func(t *testing.T) {
		t.Run("slashTest", slashTest)
	})
}

// Test penalty slash when a validator is offline
func slashTest(t *testing.T) {
	log.Infoln("===================================================================")
	log.Infoln("======================== Test slash ===========================")

	transactor := tc.NewTestTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnCLIAddr,
		tc.SgnPassphrase,
	)

	amts := []*big.Int{big.NewInt(3000000000000000000), big.NewInt(3000000000000000000), big.NewInt(1000000000000000000)}
	tc.AddValidators(t, transactor, tc.ValEthKs[:3], tc.ValAccounts[:3], amts)

	// Add second delegator to third validator
	_, auth, err := tc.GetAuth(tc.ValEthKs[0])
	require.NoError(t, err, "failed to get auth")
	err = tc.DelegateStake(auth, mainchain.Hex2Addr(tc.ValEthAddrs[2]), big.NewInt(1000000000000000000))
	require.NoError(t, err, "failed to delegate stake")

	shutdownNode(2)

	log.Infoln("Query sgn about penalty info...")
	nonce := uint64(0)
	penalty, err := tc.QueryPenalty(transactor.CliCtx, nonce, 2)
	require.NoError(t, err, "failed to query penalty")
	log.Infoln("Query sgn about penalty info:", penalty.String())
	expRes1 := fmt.Sprintf(`Nonce: %d, ValidatorAddr: %s, Reason: missing_signature`, nonce, tc.ValEthAddrs[2])
	expRes2 := fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.ValEthAddrs[0])
	assert.Equal(t, expRes1, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expRes1))
	assert.Equal(t, expRes2, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes2))

	nonce = uint64(1)
	penalty, err = tc.QueryPenalty(transactor.CliCtx, nonce, 2)
	require.NoError(t, err, "failed to query penalty")
	log.Infoln("Query sgn about penalty info:", penalty.String())
	expRes1 = fmt.Sprintf(`Nonce: %d, ValidatorAddr: %s, Reason: missing_signature`, nonce, tc.ValEthAddrs[2])
	expRes2 = fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.ValEthAddrs[2])
	assert.Equal(t, expRes1, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expRes1))
	assert.Equal(t, expRes2, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes2))

	log.Infoln("Query onchain staking pool")
	var poolAmt string
	for retry := 0; retry < tc.RetryLimit; retry++ {
		ci, _ := tc.DposContract.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[2]))
		poolAmt = ci.StakingPool.String()
		if poolAmt == "1980000000000000000" {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	assert.Equal(t, "1980000000000000000", poolAmt, fmt.Sprintf("The expected StakingPool should be 1970000000000000000"))
}
