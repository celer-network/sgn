package multinode

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/assert"
)

func setUpValidator() {
	log.Infoln("set up new sgn env")
	p := &tf.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(1),
		MinStakingPool:         big.NewInt(1),
		SidechainGoLiveTimeout: big.NewInt(0),
		CelrAddr:               tf.E2eProfile.CelrAddr,
	}
	setupNewSGNEnv(p)
	tf.SleepWithLog(10, "sgn being ready")
}

func TestE2EValidator(t *testing.T) {
	setUpValidator()

	t.Run("e2e-validator", func(t *testing.T) {
		t.Run("validatorTest", validatorTest)
	})
}

func validatorTest(t *testing.T) {
	log.Info("===================================================================")
	log.Info("======================== Test validator ===========================")

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
	for i := 0; i < 2; i++ {
		log.Infoln("Adding validator", i)

		// get auth
		keyAddr, auth, err := getAuth(ethKeystores[i], ethKeystorePps[i])
		tf.ChkTestErr(t, err, "failed to get auth")

		// get sgnAddr
		sgnAddr, err := sdk.AccAddressFromBech32(sgnOperators[i])
		tf.ChkTestErr(t, err, "failed to parse sgn address")

		err = tf.InitializeCandidate(auth, sgnAddr, big.NewInt(1))
		tf.ChkTestErr(t, err, "failed to initialize candidate")

		log.Info("Query sgn about the validator candidate...")
		candidate, err := validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, keyAddr.Hex())
		tf.ChkTestErr(t, err, "failed to queryCandidate")
		log.Infoln("Query sgn about the validator candidate:", candidate)
		expectedRes := fmt.Sprintf(`Operator: %s, StakingPool: %d`, sgnOperators[i], 0) // defined in Candidate.String()
		assert.Equal(t, expectedRes, candidate.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

		err = tf.DelegateStake(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, keyAddr, amts[i])
		tf.ChkTestErr(t, err, "failed to delegate stake")

		log.Info("Query sgn about the delegator to check if it has correct stakes...")
		delegator, err := validator.CLIQueryDelegator(transactor.CliCtx, validator.RouterKey, keyAddr.Hex(), keyAddr.Hex())
		tf.ChkTestErr(t, err, "failed to queryDelegator")
		log.Infoln("Query sgn about the validator delegator:", delegator)
		expectedRes = fmt.Sprintf(`EthAddress: %s, DelegatedStake: %d`, mainchain.Addr2Hex(keyAddr), amts[i]) // defined in Delegator.String()
		assert.Equal(t, expectedRes, delegator.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

		log.Info("Query sgn about the candidate to check if it has correct stakes...")
		candidate, err = validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, keyAddr.Hex())
		tf.ChkTestErr(t, err, "failed to queryCandidate")
		log.Infoln("Query sgn about the validator candidate:", candidate)
		expectedRes = fmt.Sprintf(`Operator: %s, StakingPool: %d`, sgnOperators[i], amts[i]) // defined in Candidate.String()
		assert.Equal(t, expectedRes, candidate.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

		log.Info("Query sgn about the validator to check if it has correct stakes...")
		validators, err := validator.CLIQueryValidators(transactor.CliCtx, staking.RouterKey)
		tf.ChkTestErr(t, err, "failed to queryValidators")
		log.Infoln("Query sgn about the validators:\n", validators)
		assert.Equal(t, i+1, len(validators), fmt.Sprintf("The length of validators should be \"%d\"", i+1))
		assert.Equal(t, sdk.NewIntFromBigInt(amts[i]), validators[i].Tokens, "validator token should be 1000000000000000000")
		assert.Equal(t, sdk.Bonded, validators[i].Status, "validator should be bonded")
	}

	// fail to add a validator 2 because it doesn't have enough delegation
	keyAddr, auth, err := getAuth(ethKeystores[2], ethKeystorePps[2])
	tf.ChkTestErr(t, err, "failed to get auth")
	sgnAddr, err := sdk.AccAddressFromBech32(sgnOperators[2])
	tf.ChkTestErr(t, err, "failed to parse sgn address")
	err = tf.InitializeCandidate(auth, sgnAddr, big.NewInt(10))
	tf.ChkTestErr(t, err, "failed to initialize candidate")
	initialDelegation := big.NewInt(1)
	err = tf.DelegateStake(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, keyAddr, initialDelegation)
	tf.ChkTestErr(t, err, "failed to delegate stake")

	log.Info("Query sgn about validators to check if validator 2 is not added...")
	validators, err := validator.CLIQueryValidators(transactor.CliCtx, staking.RouterKey)
	tf.ChkTestErr(t, err, "failed to queryValidators")
	log.Infoln("Query sgn about the validators:\n", validators)
	assert.Equal(t, 2, len(validators), fmt.Sprintf("The length of validators should be \"%d\"", 2))

	// correctly add validator 2 with enough delegation
	err = tf.DelegateStake(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, keyAddr, amts[2])
	tf.ChkTestErr(t, err, "failed to delegate stake")
	log.Info("Query sgn about the validator to check if it has correct stakes...")
	validators, err = validator.CLIQueryValidators(transactor.CliCtx, staking.RouterKey)
	tf.ChkTestErr(t, err, "failed to queryValidators")
	log.Infoln("Query sgn about the validators:\n", validators)
	assert.Equal(t, 3, len(validators), fmt.Sprintf("The length of validators should be \"%d\"", 3))

	sgnAddrTest := "cosmosvaloper122w97t8vsa3538fr3ylvz3hvuqxrgpnarnd9t6" // for test
	// validator, err := validator.CLIQueryValidator(transactor.CliCtx, staking.RouterKey, sgnOperators[2])
	validator, err := validator.CLIQueryValidator(transactor.CliCtx, staking.RouterKey, sgnAddrTest)
	tf.ChkTestErr(t, err, "failed to queryValidator")
	log.Infoln("Query sgn about the validator:\n", validator)
	assert.Equal(t, sdk.NewIntFromBigInt(big.NewInt(0).Add(initialDelegation, amts[2])), validator.Tokens, "validator token should be 1000000000000000001")
	assert.Equal(t, sdk.Bonded, validator.Status, "validator should be bonded")

	// normally remove validator 1 by intendWithdraw

	// normally add back validator 1

	// withdraw delegation to make it under limit. validator 2 should be removed

	// if a validator with more than 1/3 staking quit, the chain should halt. - validator 0
}

func getAuth(ks, pp string) (addr mainchain.Addr, auth *bind.TransactOpts, err error) {
	keystoreBytes, err := ioutil.ReadFile(ks)
	if err != nil {
		return
	}
	key, err := keystore.DecryptKey(keystoreBytes, pp)
	if err != nil {
		return
	}
	addr = key.Address
	auth, err = bind.NewTransactor(strings.NewReader(string(keystoreBytes)), pp)
	if err != nil {
		return
	}

	return
}
