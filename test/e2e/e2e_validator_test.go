package e2e

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
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

	log.Info("=====================================================================")
	log.Info("======================== Test validator ===========================")

	ethAddress := tf.EthClient.Address
	transactor := tf.Transactor
	amt := big.NewInt(100)

	initializeCandidate()
	log.Info("Query sgn about the validator candidate...")
	candidate, err := validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, ethAddress.String())
	tf.ChkErr(err, "failed to queryCandidate")
	log.Infoln("Query sgn about the validator candidate:", candidate)
	expectedRes := fmt.Sprintf(`Operator: %s, StakingPool: %d`, client0SGNAddrStr, 0) // defined in Candidate.String()
	assert.Equal(t, expectedRes, candidate.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	delegateStake()
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

func initializeCandidate() {
	ctx := context.Background()
	conn := tf.EthClient.Client
	auth := tf.EthClient.Auth
	guardContract := tf.EthClient.Guard
	sgnAddr, err := sdk.AccAddressFromBech32(client0SGNAddrStr)
	tf.ChkErr(err, "Parse SGN address error")

	log.Info("Call initializeCandidate on guard contract using the validator eth address...")
	tx, err := guardContract.InitializeCandidate(auth, big.NewInt(1), sgnAddr.Bytes())
	tf.ChkErr(err, "failed to InitializeCandidate")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "InitializeCandidate")
	sleepWithLog(30, "sgn syncing InitializeCandidate event on mainchain")
}

func delegateStake() {
	ctx := context.Background()
	conn := tf.EthClient.Client
	auth := tf.EthClient.Auth
	ethAddress := tf.EthClient.Address
	guardContract := tf.EthClient.Guard

	log.Info("Call delegate on guard contract to delegate stake to the validator eth address...")
	amt := big.NewInt(100)
	tx, err := celrContract.Approve(auth, guardAddr, amt)
	tf.ChkErr(err, "failed to approve CELR to Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")
	tx, err = guardContract.Delegate(auth, ethAddress, amt)
	tf.ChkErr(err, "failed to call delegate of Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Delegate to validator")
	sleepWithLog(30, "sgn syncing Delegate event on mainchain")
}
