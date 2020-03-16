package multinode

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
)

func setUpSlash() {
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
	tc.SleepWithLog(10, "sgn syncing")
}

func TestE2ESlash(t *testing.T) {
	setUpSlash()

	t.Run("e2e-slash", func(t *testing.T) {
		t.Run("slashTest", slashTest)
	})
}

// Test penalty slash when a validator is offline
func slashTest(t *testing.T) {
	log.Infoln("===================================================================")
	log.Infoln("======================== Test slash ===========================")

	transactor := tc.NewTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnCLIAddr,
		tc.SgnPassphrase,
	)

	amts := []*big.Int{big.NewInt(2000000000000000000), big.NewInt(2000000000000000000), big.NewInt(1000000000000000000)}
	tc.AddValidators(t, transactor, tc.ValEthKs[:], tc.SgnOperators[:], amts)
	shutdownNode(2)

	log.Infoln("Query sgn about penalty info...")
	nonce := uint64(0)
	penalty, err := tc.QueryPenalty(transactor.CliCtx, nonce, 2)
	tc.ChkTestErr(t, err, "failed to query penalty")
	log.Infoln("Query sgn about penalty info:", penalty.String())
	expRes1 := fmt.Sprintf(`Nonce: %d, ValidatorAddr: %s, Reason: missing_signature`, nonce, tc.ValEthAddrs[2])
	expRes2 := fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.ValEthAddrs[2])
	assert.Equal(t, expRes1, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expRes1))
	assert.Equal(t, expRes2, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes2))

	log.Infoln("Query onchain staking pool")
	var poolAmt string
	for retry := 0; retry < 30; retry++ {
		ci, _ := tc.Client0.Guard.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[2]))
		poolAmt = ci.StakingPool.String()
		if poolAmt == "990000000000000000" {
			break
		}
		time.Sleep(time.Second)
	}
	assert.Equal(t, "990000000000000000", poolAmt, fmt.Sprintf("The expected StakingPool should be 990000000000000000"))
}
