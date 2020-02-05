package multinode

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/test/e2e/common"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/slash"
	stypes "github.com/celer-network/sgn/x/slash/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
)

func setUpSlash() {
	log.Infoln("set up new sgn env")
	p := &tf.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MinStakingPool:         big.NewInt(0),
		SidechainGoLiveTimeout: big.NewInt(0),
		CelrAddr:               tf.E2eProfile.CelrAddr,
		MaxValidatorNum:        big.NewInt(11),
	}
	setupNewSGNEnv(p)
	tf.SleepWithLog(10, "sgn syncing")
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

	transactor := tf.NewTransactor(
		t,
		tf.SgnCLIHome,
		tf.SgnChainID,
		tf.SgnNodeURI,
		tf.SgnTransactor,
		tf.SgnPassphrase,
		tf.SgnGasPrice,
	)

	amts := []*big.Int{big.NewInt(1000000000000000000), big.NewInt(1000000000000000000), big.NewInt(100000000000000000)}
	tc.AddValidators(t, transactor, tf.EthKeystores[:], tf.SgnOperators[:], amts)

	shutdownNode(2)

	log.Infoln("Query sgn about penalty info...")
	var penalty stypes.Penalty
	var err error
	nonce := uint64(0)
	expRes1 := fmt.Sprintf(`Nonce: %d, ValidatorAddr: %s, Reason: missing_signature`, nonce, tf.EthAddresses[2])
	expRes2 := fmt.Sprintf(`Account: %s, Amount: 1000000000000000`, tf.EthAddresses[2])
	for retry := 0; retry < 30; retry++ {
		penalty, err = slash.CLIQueryPenalty(transactor.CliCtx, slash.StoreKey, nonce)
		if err == nil && penalty.String() == expRes1 && penalty.PenalizedDelegators[0].String() == expRes2 &&
			len(penalty.PenaltyProtoBytes) > 0 && len(penalty.Sigs) == 2 {
			break
		}
		time.Sleep(2 * time.Second)
	}
	tf.ChkTestErr(t, err, "failed to query penalty")
	log.Infoln("Query sgn about penalty info:", penalty.String())
	assert.Equal(t, expRes1, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expRes1))
	assert.Equal(t, expRes2, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes2))
	assert.Greater(t, len(penalty.PenaltyProtoBytes), 0, fmt.Sprintf("The length of penaltyProtoBytes should be larger than 0"))
	assert.Equal(t, 2, len(penalty.Sigs), fmt.Sprintf("The length of validators should be 2"))

	log.Infoln("Query onchain staking pool")
	var poolAmt string
	for retry := 0; retry < 30; retry++ {
		ci, _ := tf.DefaultTestEthClient.Guard.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(tf.EthAddresses[2]))
		poolAmt = ci.StakingPool.String()
		if poolAmt == "99000000000000000" {
			break
		}
		time.Sleep(time.Second)
	}
	assert.Equal(t, "99000000000000000", poolAmt, fmt.Sprintf("The expected StakingPool should be 99000000000000000"))
}
