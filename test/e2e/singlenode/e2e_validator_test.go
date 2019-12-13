package singlenode

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setUpValidator() []tf.Killable {
	p := &tf.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(1),
		MinStakingPool:         big.NewInt(1),
		SidechainGoLiveTimeout: big.NewInt(0),
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
	ethAddress := tf.DefaultTestEthClient.Address
	transactor := tf.NewTransactor(
		t,
		viper.GetString(common.FlagSgnCLIHome),
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetStringSlice(common.FlagSgnTransactors)[0],
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetString(common.FlagSgnGasPrice),
	)
	amt := big.NewInt(1000000000000000000)
	sgnAddr, err := sdk.AccAddressFromBech32(tf.Client0SGNAddrStr)
	tf.ChkTestErr(t, err, "failed to parse sgn address")

	err = tf.InitializeCandidate(auth, sgnAddr)
	tf.ChkTestErr(t, err, "failed to initialize candidate")

	log.Info("Query sgn about the validator candidate...")
	candidate, err := validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, ethAddress.Hex())
	tf.ChkTestErr(t, err, "failed to queryCandidate")
	log.Infoln("Query sgn about the validator candidate:", candidate)
	expectedRes := fmt.Sprintf(`Operator: %s, StakingPool: %d`, tf.Client0SGNAddrStr, 0) // defined in Candidate.String()
	assert.Equal(t, expectedRes, candidate.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	err = tf.DelegateStake(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, ethAddress, amt)
	tf.ChkTestErr(t, err, "failed to delegate stake")

	log.Info("Query sgn about the delegator to check if it has correct stakes...")
	delegator, err := validator.CLIQueryDelegator(transactor.CliCtx, validator.RouterKey, ethAddress.Hex(), ethAddress.Hex())
	tf.ChkTestErr(t, err, "failed to queryDelegator")
	log.Infoln("Query sgn about the validator delegator:", delegator)
	expectedRes = fmt.Sprintf(`EthAddress: %s, DelegatedStake: %d`, mainchain.Addr2Hex(ethAddress), amt) // defined in Delegator.String()
	assert.Equal(t, expectedRes, delegator.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Query sgn about the candidate to check if it has correct stakes
	log.Info("Query sgn about the validator candidate...")
	candidate, err = validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, ethAddress.Hex())
	tf.ChkTestErr(t, err, "failed to queryCandidate")
	log.Infoln("Query sgn about the validator candidate:", candidate)
	expectedRes = fmt.Sprintf(`Operator: %s, StakingPool: %d`, tf.Client0SGNAddrStr, amt) // defined in Candidate.String()
	assert.Equal(t, expectedRes, candidate.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Info("Query sgn about the validator to check if it has correct stakes...")
	validators, err := validator.CLIQueryValidators(transactor.CliCtx, staking.RouterKey)
	tf.ChkTestErr(t, err, "failed to queryValidators")
	log.Infoln("Query sgn about the validators:", validators)
	assert.Equal(t, 1, len(validators), "The length of validators should be 1")
	assert.True(t, validators[0].Tokens.Equal(sdk.NewIntFromBigInt(amt)), "validator token should be 1000000000000000000")
	assert.Equal(t, sdk.Bonded, validators[0].Status, "validator should be bonded")
}
