package multinode

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/transactor"
	sgnval "github.com/celer-network/sgn/x/validator"
	vtypes "github.com/celer-network/sgn/x/validator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/assert"
)

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

func addValidators(t *testing.T, transactor *transactor.Transactor, ethkss, ethpps, sgnops, sgnopAddrs []string, amts []*big.Int) {
	for i := 0; i < len(ethkss); i++ {
		log.Infoln("Adding validator", i)
		ethAddr, auth, err := getAuth(ethkss[i], ethpps[i])
		tf.ChkTestErr(t, err, "failed to get auth")
		addCandidateWithStake(t, transactor, ethAddr, auth, sgnops[i], sgnopAddrs[i], amts[i], big.NewInt(1), true)
	}
}

func addCandidateWithStake(t *testing.T, transactor *transactor.Transactor,
	ethAddr mainchain.Addr, auth *bind.TransactOpts,
	sgnop, sgnopAddr string,
	amt *big.Int, minAmt *big.Int,
	isValidator bool) {

	// get sgnAddr
	sgnAddr, err := sdk.AccAddressFromBech32(sgnop)
	tf.ChkTestErr(t, err, "failed to parse sgn address")

	// add candidate
	err = tf.InitializeCandidate(auth, sgnAddr, minAmt)
	tf.ChkTestErr(t, err, "failed to initialize candidate")

	log.Infof("Query sgn about the validator candidate %s ...", ethAddr.Hex())
	checkCandidate(t, transactor, ethAddr, sgnop, big.NewInt(0))

	// self delegate stake
	err = tf.DelegateStake(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, ethAddr, amt)
	tf.ChkTestErr(t, err, "failed to delegate stake")

	log.Info("Query sgn about the delegator to check if it has correct stakes...")
	checkDelegator(t, transactor, ethAddr, ethAddr, amt)

	log.Info("Query sgn about the candidate to check if it has correct stakes...")
	checkCandidate(t, transactor, ethAddr, sgnop, amt)

	if isValidator {
		log.Info("Query sgn about the validators to check if it has correct stakes...")
		checkValidator(t, transactor, sgnopAddr, amt, sdk.Bonded)
	}
}

func checkDelegator(t *testing.T, transactor *transactor.Transactor, validatorAddr, delegatorAddr mainchain.Addr, expAmt *big.Int) {
	var delegator vtypes.Delegator
	var err error
	expectedRes := fmt.Sprintf(`EthAddress: %s, DelegatedStake: %s`, mainchain.Addr2Hex(delegatorAddr), expAmt) // defined in Delegator.String()
	for retry := 0; retry < 30; retry++ {
		delegator, err = sgnval.CLIQueryDelegator(transactor.CliCtx, sgnval.RouterKey, validatorAddr.Hex(), delegatorAddr.Hex())
		if err == nil && expectedRes == delegator.String() {
			break
		}
		time.Sleep(2 * time.Second)
	}
	tf.ChkTestErr(t, err, "failed to queryDelegator")
	log.Infoln("Query sgn about the validator's delegator:", delegator)
	assert.Equal(t, expectedRes, delegator.String(), "The expected result should be: "+expectedRes)
}

func checkCandidate(t *testing.T, transactor *transactor.Transactor, ethAddr mainchain.Addr, sgnop string, expAmt *big.Int) {
	var candidate vtypes.Candidate
	var err error
	expectedRes := fmt.Sprintf(`Operator: %s, StakingPool: %s`, sgnop, expAmt) // defined in Candidate.String()
	for retry := 0; retry < 30; retry++ {
		candidate, err = sgnval.CLIQueryCandidate(transactor.CliCtx, sgnval.RouterKey, ethAddr.Hex())
		if err == nil && expectedRes == candidate.String() {
			break
		}
		time.Sleep(2 * time.Second)
	}
	tf.ChkTestErr(t, err, "failed to queryCandidate")
	log.Infoln("Query sgn about the validator candidate:", candidate)
	assert.Equal(t, expectedRes, candidate.String(), "The expected result should be: "+expectedRes)
}

func checkValidator(t *testing.T, transactor *transactor.Transactor, sgnopAddr string, expAmt *big.Int, expStatus sdk.BondStatus) {
	var validator stypes.Validator
	var err error
	for retry := 0; retry < 30; retry++ {
		validator, err = sgnval.CLIQueryValidator(transactor.CliCtx, staking.RouterKey, sgnopAddr)
		if err == nil && validator.Status == expStatus && validator.Tokens.BigInt().Cmp(expAmt) == 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}
	tf.ChkTestErr(t, err, "failed to queryValidator")
	log.Infoln("Query sgn about the validator:\n", validator)
	assert.Equal(t, expAmt.String(), validator.Tokens.String(), "validator token should be "+expAmt.String())
	assert.Equal(t, expStatus, validator.Status, "validator should be "+sdkStatusName(validator.Status))
}

func checkValidatorStatus(t *testing.T, transactor *transactor.Transactor, sgnopAddr string, expStatus sdk.BondStatus) {
	var validator stypes.Validator
	var err error
	for retry := 0; retry < 30; retry++ {
		validator, err = sgnval.CLIQueryValidator(transactor.CliCtx, staking.RouterKey, sgnopAddr)
		if err == nil && validator.Status == expStatus {
			break
		}
		time.Sleep(2 * time.Second)
	}
	tf.ChkTestErr(t, err, "failed to queryValidator")
	log.Infoln("Query sgn about the validator:\n", validator)
	assert.Equal(t, expStatus, validator.Status, "validator should be "+sdkStatusName(validator.Status))
}

func checkValidatorNum(t *testing.T, transactor *transactor.Transactor, expNum int) {
	var validators stypes.Validators
	var err error
	for retry := 0; retry < 30; retry++ {
		validators, err = sgnval.CLIQueryBondedValidators(transactor.CliCtx, staking.RouterKey)
		if err == nil && len(validators) == expNum {
			break
		}
		time.Sleep(2 * time.Second)
	}
	tf.ChkTestErr(t, err, "failed to queryValidators")
	log.Infoln("Query sgn about the validators:\n", validators)
	assert.Equal(t, expNum, len(validators), "The length of validators should be: "+strconv.Itoa(expNum))
}

func sdkStatusName(status sdk.BondStatus) string {
	switch status {
	case sdk.Unbonded:
		return "Unbonded"
	case sdk.Unbonding:
		return "Unbonding"
	case sdk.Bonded:
		return "Bonded"
	default:
		return "Invalid"
	}
}
