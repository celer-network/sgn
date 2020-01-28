package multinode

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/slash"
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

// Test penalty slash when a validator is offline
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

	shutdownNode(2)
	tf.SleepWithLog(25, "wait for slash")

	nonce := uint64(0)
	penalty, err := slash.CLIQueryPenalty(transactor.CliCtx, slash.StoreKey, nonce)
	tf.ChkTestErr(t, err, "failed to query penalty")
	expectedRes := fmt.Sprintf(`Nonce: %d, ValidatorAddr: %s, Reason: missing_signature`, nonce, ethAddresses[2])
	assert.Equal(t, expectedRes, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	expectedRes = fmt.Sprintf(`Account: %s, Amount: 1000000000000000`, ethAddresses[2])
	assert.Equal(t, expectedRes, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	assert.Greater(t, len(penalty.PenaltyProtoBytes), 0, fmt.Sprintf("The length of penaltyProtoBytes should be larger than 0"))
	assert.Equal(t, 2, len(penalty.Sigs), fmt.Sprintf("The length of validators should be 2"))

	ci, _ := tf.DefaultTestEthClient.Guard.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(ethAddresses[2]))
	assert.Equal(t, "99000000000000000", ci.StakingPool.String(), fmt.Sprintf("The expected StakingPool should be 99000000000000000"))
}
