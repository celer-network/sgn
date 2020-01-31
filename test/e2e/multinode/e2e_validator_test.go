package multinode

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	sgnval "github.com/celer-network/sgn/x/validator"
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
		MaxValidatorNum:        big.NewInt(11),
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
	log.Infoln("---------- It should add two validators successfully ----------")
	for i := 0; i < 2; i++ {
		log.Infoln("Adding validator", i)

		// get auth
		ethAddr, auth, err := getAuth(ethKeystores[i], ethKeystorePps[i])
		tf.ChkTestErr(t, err, "failed to get auth")

		// get sgnAddr
		sgnAddr, err := sdk.AccAddressFromBech32(sgnOperators[i])
		tf.ChkTestErr(t, err, "failed to parse sgn address")

		err = tf.InitializeCandidate(auth, sgnAddr, big.NewInt(1))
		tf.ChkTestErr(t, err, "failed to initialize candidate")

		log.Info("Query sgn about the validator candidate...")
		candidate, err := sgnval.CLIQueryCandidate(transactor.CliCtx, sgnval.RouterKey, ethAddr.Hex())
		tf.ChkTestErr(t, err, "failed to queryCandidate")
		log.Infoln("Query sgn about the validator candidate:", candidate)
		expectedRes := fmt.Sprintf(`Operator: %s, StakingPool: %d`, sgnOperators[i], 0) // defined in Candidate.String()
		assert.Equal(t, expectedRes, candidate.String(), "The expected result should be: "+expectedRes)

		err = tf.DelegateStake(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, ethAddr, amts[i])
		tf.ChkTestErr(t, err, "failed to delegate stake")

		log.Info("Query sgn about the delegator to check if it has correct stakes...")
		delegator, err := sgnval.CLIQueryDelegator(transactor.CliCtx, sgnval.RouterKey, ethAddr.Hex(), ethAddr.Hex())
		tf.ChkTestErr(t, err, "failed to queryDelegator")
		log.Infoln("Query sgn about the validator delegator:", delegator)
		expectedRes = fmt.Sprintf(`EthAddress: %s, DelegatedStake: %d`, mainchain.Addr2Hex(ethAddr), amts[i]) // defined in Delegator.String()
		assert.Equal(t, expectedRes, delegator.String(), "The expected result should be: "+expectedRes)

		log.Info("Query sgn about the candidate to check if it has correct stakes...")
		candidate, err = sgnval.CLIQueryCandidate(transactor.CliCtx, sgnval.RouterKey, ethAddr.Hex())
		tf.ChkTestErr(t, err, "failed to queryCandidate")
		log.Infoln("Query sgn about the validator candidate:", candidate)
		expectedRes = fmt.Sprintf(`Operator: %s, StakingPool: %d`, sgnOperators[i], amts[i]) // defined in Candidate.String()
		assert.Equal(t, expectedRes, candidate.String(), "The expected result should be: "+expectedRes)

		log.Info("Query sgn about the validators to check if it has correct stakes...")
		validators, err := sgnval.CLIQueryBondedValidators(transactor.CliCtx, staking.RouterKey)
		tf.ChkTestErr(t, err, "failed to queryValidators")
		log.Infoln("Query sgn about the validators:\n", validators)
		assert.Equal(t, i+1, len(validators), "The length of validators should be: "+strconv.Itoa(i+1))
		validator, err := sgnval.CLIQueryValidator(transactor.CliCtx, staking.RouterKey, sgnOperatorValAddrs[i])
		tf.ChkTestErr(t, err, "failed to queryValidator")
		log.Infoln("Query sgn about the validator:\n", validator)
		assert.Equal(t, sdk.NewIntFromBigInt(amts[i]), validator.Tokens, "validator token should be "+amts[i].String())
		assert.Equal(t, sdk.Bonded, validator.Status, "validator should be bonded")
	}

	log.Infoln("---------- It should fail to add validator 2 without enough delegation ----------")
	ethAddr, auth, err := getAuth(ethKeystores[2], ethKeystorePps[2])
	tf.ChkTestErr(t, err, "failed to get auth")
	sgnAddr, err := sdk.AccAddressFromBech32(sgnOperators[2])
	tf.ChkTestErr(t, err, "failed to parse sgn address")
	err = tf.InitializeCandidate(auth, sgnAddr, big.NewInt(10))
	tf.ChkTestErr(t, err, "failed to initialize candidate")
	initialDelegation := big.NewInt(1)
	err = tf.DelegateStake(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, ethAddr, initialDelegation)
	tf.ChkTestErr(t, err, "failed to delegate stake")

	log.Info("Query sgn about validators to check if validator 2 is not added...")
	validators, err := sgnval.CLIQueryBondedValidators(transactor.CliCtx, staking.RouterKey)
	tf.ChkTestErr(t, err, "failed to queryValidators")
	log.Infoln("Query sgn about the validators:\n", validators)
	assert.Equal(t, 2, len(validators), "The length of validators should be: 2")

	log.Infoln("---------- It should correctly add validator 2 with enough delegation ----------")
	err = tf.DelegateStake(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, ethAddr, big.NewInt(0).Sub(amts[2], initialDelegation))
	tf.ChkTestErr(t, err, "failed to delegate stake")
	log.Info("Query sgn about the validators to check if it has correct stakes...")
	validators, err = sgnval.CLIQueryBondedValidators(transactor.CliCtx, staking.RouterKey)
	tf.ChkTestErr(t, err, "failed to queryValidators")
	log.Infoln("Query sgn about the validators:\n", validators)
	assert.Equal(t, 3, len(validators), "The length of validators should be: 3")
	validator, err := sgnval.CLIQueryValidator(transactor.CliCtx, staking.RouterKey, sgnOperatorValAddrs[2])
	tf.ChkTestErr(t, err, "failed to queryValidator")
	log.Infoln("Query sgn about the validator:\n", validator)
	assert.Equal(t, sdk.NewIntFromBigInt(amts[2]), validator.Tokens, "validator token should be 1000000000000000000")
	assert.Equal(t, sdk.Bonded, validator.Status, "validator should be bonded")

	log.Infoln("---------- It should successfully remove validator 2 caused by intendWithdraw ----------")
	err = tf.IntendWithdraw(auth, ethAddr, amts[2])
	tf.ChkTestErr(t, err, "failed to intendWithdraw stake")
	log.Info("Query sgn about the validators to check if it has correct stakes...")
	validators, err = sgnval.CLIQueryBondedValidators(transactor.CliCtx, staking.RouterKey)
	tf.ChkTestErr(t, err, "failed to queryValidators")
	log.Infoln("Query sgn about the validators:\n", validators)
	assert.Equal(t, 2, len(validators), "The length of validators should be: 2")
	validator, err = sgnval.CLIQueryValidator(transactor.CliCtx, staking.RouterKey, sgnOperatorValAddrs[2])
	assert.Equal(t, sdk.Unbonding, validator.Status, "validator should be unbonding")

	// TODO: normally add back validator 1
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
