package e2e

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/celer-network/goutils/log"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/assert"
)

func setUpValidator() []tf.Killable {
	p := &SGNParams{
		blameTimeout:           big.NewInt(10),
		minValidatorNum:        big.NewInt(1),
		minStakingPool:         big.NewInt(1),
		sidechainGoLiveTimeout: big.NewInt(0),
	}
	res := setupNewSGNEnv(p, "validator")
	sleepWithLog(20, "sgn being ready")

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

	auth := tf.EthClient.Auth
	ethAddress := tf.EthClient.Address
	transactor := tf.Transactor
	amt := big.NewInt(100)
	sgnAddr, err := sdk.AccAddressFromBech32(client0SGNAddrStr)
	tf.ChkErr(err, "failed to parse sgn address")

	err = initializeCandidate(auth, sgnAddr)
	tf.ChkErr(err, "failed to initialize candidate")

	log.Info("Query sgn about the validator candidate...")
	candidate, err := validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, ethAddress.String())
	tf.ChkErr(err, "failed to queryCandidate")
	log.Infoln("Query sgn about the validator candidate:", candidate)
	expectedRes := fmt.Sprintf(`Operator: %s, StakingPool: %d`, client0SGNAddrStr, 0) // defined in Candidate.String()
	assert.Equal(t, expectedRes, candidate.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	err = delegateStake(auth, ethAddress, amt)
	tf.ChkErr(err, "failed to delegate stake")

	log.Info("Query sgn about the delegator to check if it has correct stakes...")
	delegator, err := validator.CLIQueryDelegator(transactor.CliCtx, validator.RouterKey, ethAddress.String(), ethAddress.String())
	tf.ChkErr(err, "failed to queryDelegator")
	log.Infoln("Query sgn about the validator delegator:", delegator)
	expectedRes = fmt.Sprintf(`EthAddress: %s, DelegatedStake: %d`, ethAddress.String(), amt) // defined in Delegator.String()
	assert.Equal(t, expectedRes, delegator.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Query sgn about the candidate to check if it has correct stakes
	log.Info("Query sgn about the validator candidate...")
	candidate, err = validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, ethAddress.String())
	tf.ChkErr(err, "failed to queryCandidate")
	log.Infoln("Query sgn about the validator candidate:", candidate)
	expectedRes = fmt.Sprintf(`Operator: %s, StakingPool: %d`, client0SGNAddrStr, amt) // defined in Candidate.String()
	assert.Equal(t, expectedRes, candidate.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Info("Query sgn about the validator to check if it has correct stakes...")
	validators, err := validator.CLIQueryValidators(transactor.CliCtx, staking.RouterKey)
	tf.ChkErr(err, "failed to queryValidators")
	log.Infoln("Query sgn about the validators:", validators)
	// TODO: use a better way to assert/check the validity of the lengthy query results.
	// expectedRes = fmt.Sprintf("StakingPool: %d", amt) // defined in Candidate.String()
	// assert.Equal(t, expectedRes, validators.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
}
