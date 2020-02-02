package multinode

import (
	"math/big"
	"testing"

	"github.com/celer-network/goutils/log"
	tf "github.com/celer-network/sgn/testing"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setUpValidator(maxValidatorNum *big.Int) {
	log.Infoln("set up new sgn env")
	p := &tf.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(1),
		MinStakingPool:         big.NewInt(1),
		SidechainGoLiveTimeout: big.NewInt(0),
		CelrAddr:               tf.E2eProfile.CelrAddr,
		MaxValidatorNum:        maxValidatorNum,
	}
	setupNewSGNEnv(p)
	tf.SleepWithLog(10, "sgn being ready")
}

func TestE2EValidator(t *testing.T) {
	t.Run("e2e-validator", func(t *testing.T) {
		t.Run("validatorTest", validatorTest)
		t.Run("replaceValidatorTest", replaceValidatorTest)
	})
}

func validatorTest(t *testing.T) {
	log.Info("===================================================================")
	log.Info("======================== Test validator ===========================")
	setUpValidator(big.NewInt(11))

	transactor := tf.NewTransactor(
		t,
		sgnCLIHomes[0],
		sgnChainID,
		sgnNodeURIs[0],
		sgnTransactors[0],
		sgnPassphrase,
		sgnGasPrice,
	)

	// delegation ratio. V0 : V1 : V2 = 2 : 1 : 1
	amts := []*big.Int{big.NewInt(2000000000000000000), big.NewInt(1000000000000000000), big.NewInt(1000000000000000000)}

	// add two validators, 0 and 1
	log.Infoln("---------- It should add two validators successfully ----------")
	for i := 0; i < 2; i++ {
		log.Infoln("Adding validator", i)
		// get auth
		ethAddr, auth, err := getAuth(ethKeystores[i], ethKeystorePps[i])
		tf.ChkTestErr(t, err, "failed to get auth")
		addCandidateWithStake(t, transactor, ethAddr, auth, sgnOperators[i], sgnOperatorValAddrs[i], amts[i], big.NewInt(1), true)
		checkValidatorNum(t, transactor, i+1)
	}

	log.Infoln("---------- It should fail to add validator 2 without enough delegation ----------")
	ethAddr, auth, err := getAuth(ethKeystores[2], ethKeystorePps[2])
	tf.ChkTestErr(t, err, "failed to get auth")
	initialDelegation := big.NewInt(1)
	addCandidateWithStake(t, transactor, ethAddr, auth, sgnOperators[2], sgnOperatorValAddrs[2], initialDelegation, big.NewInt(10), false)
	log.Info("Query sgn about validators to check if validator 2 is not added...")
	checkValidatorNum(t, transactor, 2)

	log.Infoln("---------- It should correctly add validator 2 with enough delegation ----------")
	err = tf.DelegateStake(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, ethAddr, big.NewInt(0).Sub(amts[2], initialDelegation))
	tf.ChkTestErr(t, err, "failed to delegate stake")
	checkValidatorNum(t, transactor, 3)
	checkValidator(t, transactor, sgnOperatorValAddrs[2], amts[2], sdk.Bonded)

	log.Infoln("---------- It should successfully remove validator 2 caused by intendWithdraw ----------")
	err = tf.IntendWithdraw(auth, ethAddr, amts[2])
	tf.ChkTestErr(t, err, "failed to intendWithdraw stake")
	log.Info("Query sgn about the validators to check if it has correct number of validators...")
	checkValidatorNum(t, transactor, 2)
	checkValidatorStatus(t, transactor, sgnOperatorValAddrs[2], sdk.Unbonding)

	// TODO: normally add back validator 1
}

func replaceValidatorTest(t *testing.T) {
	log.Info("===================================================================")
	log.Info("========================  Test replacing validator ===========================")
	setUpValidator(big.NewInt(2))

	transactor := tf.NewTransactor(
		t,
		sgnCLIHomes[0],
		sgnChainID,
		sgnNodeURIs[0],
		sgnTransactors[0],
		sgnPassphrase,
		sgnGasPrice,
	)

	amts := []*big.Int{big.NewInt(5000000000000000000), big.NewInt(1000000000000000000), big.NewInt(2000000000000000000)}

	// add two validators, 0 and 1
	addValidators(t, transactor, ethKeystores[:2], ethKeystorePps[:2], sgnOperators[:2], sgnOperatorValAddrs[:2], amts[:2])

	log.Infoln("---------- It should correctly replace validator 1 with validator 2 ----------")
	ethAddr, auth, err := getAuth(ethKeystores[2], ethKeystorePps[2])
	tf.ChkTestErr(t, err, "failed to get auth")
	addCandidateWithStake(t, transactor, ethAddr, auth, sgnOperators[2], sgnOperatorValAddrs[2], amts[2], big.NewInt(1), true)

	log.Info("Query sgn about the validators...")
	checkValidatorNum(t, transactor, 2)
	checkValidator(t, transactor, sgnOperatorValAddrs[1], amts[1], sdk.Unbonding)
}
