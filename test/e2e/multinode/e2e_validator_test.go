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
	amt := big.NewInt(1000000000000000000)

	for i := 0; i < 3; i++ {
		log.Infoln("Adding validator ", i)

		// get auth
		keystoreBytes, err := ioutil.ReadFile(ethKeystores[i])
		tf.ChkTestErr(t, err, "failed to read ethKeystore")
		key, err := keystore.DecryptKey(keystoreBytes, ethKeystorePps[i])
		tf.ChkTestErr(t, err, "failed to DecryptKey")
		auth, err := bind.NewTransactor(strings.NewReader(string(keystoreBytes)), ethKeystorePps[i])
		tf.ChkTestErr(t, err, "failed to generate auth")

		// get sgnAddr
		sgnAddr, err := sdk.AccAddressFromBech32(sgnOperators[i])
		tf.ChkTestErr(t, err, "failed to parse sgn address")

		err = tf.InitializeCandidate(auth, sgnAddr)
		tf.ChkTestErr(t, err, "failed to initialize candidate")

		log.Info("Query sgn about the validator candidate...")
		candidate, err := validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, key.Address.Hex())
		tf.ChkTestErr(t, err, "failed to queryCandidate")
		log.Infoln("Query sgn about the validator candidate:", candidate)
		expectedRes := fmt.Sprintf(`Operator: %s, StakingPool: %d`, sgnOperators[i], 0) // defined in Candidate.String()
		assert.Equal(t, expectedRes, candidate.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

		err = tf.DelegateStake(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, key.Address, amt)
		tf.ChkTestErr(t, err, "failed to delegate stake")

		log.Info("Query sgn about the delegator to check if it has correct stakes...")
		delegator, err := validator.CLIQueryDelegator(transactor.CliCtx, validator.RouterKey, key.Address.Hex(), key.Address.Hex())
		tf.ChkTestErr(t, err, "failed to queryDelegator")
		log.Infoln("Query sgn about the validator delegator:", delegator)
		expectedRes = fmt.Sprintf(`EthAddress: %s, DelegatedStake: %d`, mainchain.Addr2Hex(key.Address), amt) // defined in Delegator.String()
		assert.Equal(t, expectedRes, delegator.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

		log.Info("Query sgn about the candidate to check if it has correct stakes...")
		candidate, err = validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, key.Address.Hex())
		tf.ChkTestErr(t, err, "failed to queryCandidate")
		log.Infoln("Query sgn about the validator candidate:", candidate)
		expectedRes = fmt.Sprintf(`Operator: %s, StakingPool: %d`, sgnOperators[i], amt) // defined in Candidate.String()
		assert.Equal(t, expectedRes, candidate.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

		log.Info("Query sgn about the validator to check if it has correct stakes...")
		validators, err := validator.CLIQueryValidators(transactor.CliCtx, staking.RouterKey)
		tf.ChkTestErr(t, err, "failed to queryValidators")
		log.Infoln("Query sgn about the validators:\n", validators)
		assert.Equal(t, i+1, len(validators), fmt.Sprintf("The length of validators should be \"%d\"", i+1))
		assert.True(t, validators[i].Tokens.Equal(sdk.NewIntFromBigInt(amt)), "validator token should be 1000000000000000000")
		assert.Equal(t, sdk.Bonded, validators[i].Status, "validator should be bonded")
	}
}
