package singlenode

import (
	"math/big"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tc "github.com/celer-network/sgn/test/e2e/common"
	tf "github.com/celer-network/sgn/testing"
	"github.com/spf13/viper"
)

func setUpValidator() []tf.Killable {
	p := &tf.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(1),
		MinStakingPool:         big.NewInt(1),
		SidechainGoLiveTimeout: big.NewInt(0),
		CelrAddr:               tf.E2eProfile.CelrAddr,
		MaxValidatorNum:        big.NewInt(11),
	}
	res := setupNewSGNEnv(p, "validator")
	tf.SleepWithLog(10, "sgn being ready")

	return res
}

func TestE2EValidator(t *testing.T) {
	toKill := setUpValidator()
	defer tf.TearDown(toKill)

	t.Run("e2e-validator", func(t *testing.T) {
		t.Run("validatorTest", validatorTest)
	})
}

func validatorTest(t *testing.T) {
	// TODO: each test cases need a new and isolated sgn right now, which can't be run in parallel
	// t.Parallel()

	log.Info("===================================================================")
	log.Info("======================== Test validator ===========================")

	auth := tf.DefaultTestEthClient.Auth
	ethAddr := tf.DefaultTestEthClient.Address
	transactor := tf.NewTransactor(
		t,
		CLIHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetString(common.FlagSgnOperator),
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetString(common.FlagSgnGasPrice),
	)
	amt := big.NewInt(1000000000000000000)

	tc.AddCandidateWithStake(t, transactor, ethAddr, auth, tf.Client0SGNAddrStr, tf.Client0SGNValAddrStr, amt, big.NewInt(1), true)
	tc.CheckValidatorNum(t, transactor, 1)
}
