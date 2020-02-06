package testcommon

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	sgnval "github.com/celer-network/sgn/x/validator"
	vtypes "github.com/celer-network/sgn/x/validator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
)

type SGNParams struct {
	CelrAddr               mainchain.Addr
	BlameTimeout           *big.Int
	MinValidatorNum        *big.Int
	MinStakingPool         *big.Int
	SidechainGoLiveTimeout *big.Int
	StartGateway           bool
	MaxValidatorNum        *big.Int
}

func NewTransactor(t *testing.T, sgnCLIHome, sgnChainID, sgnNodeURI, sgnTransactor, sgnPassphrase, sgnGasPrice string) *transactor.Transactor {
	cdc := app.MakeCodec()
	tr, err := transactor.NewTransactor(
		sgnCLIHome,
		sgnChainID,
		sgnNodeURI,
		sgnTransactor,
		sgnPassphrase,
		sgnGasPrice,
		cdc,
	)
	ChkTestErr(t, err, "Failed to create new transactor.")

	return tr
}

func AddValidators(t *testing.T, transactor *transactor.Transactor, ethkss, sgnops []string, amts []*big.Int) {
	for i := 0; i < len(ethkss); i++ {
		log.Infoln("Adding validator", i)
		ethAddr, auth, err := GetAuth(ethkss[i])
		ChkTestErr(t, err, "failed to get auth")
		AddCandidateWithStake(t, transactor, ethAddr, auth, sgnops[i], amts[i], big.NewInt(1), true)
	}
}

func AddCandidateWithStake(t *testing.T, transactor *transactor.Transactor,
	ethAddr mainchain.Addr, auth *bind.TransactOpts,
	sgnop string, amt *big.Int, minAmt *big.Int, isValidator bool) {

	// get sgnAddr
	sgnAddr, err := sdk.AccAddressFromBech32(sgnop)
	ChkTestErr(t, err, "failed to parse sgn address")

	// add candidate
	err = InitializeCandidate(auth, sgnAddr, minAmt)
	ChkTestErr(t, err, "failed to initialize candidate")

	log.Infof("Query sgn about the validator candidate %s ...", ethAddr.Hex())
	CheckCandidate(t, transactor, ethAddr, sgnop, big.NewInt(0))

	// self delegate stake
	err = DelegateStake(E2eProfile.CelrContract, E2eProfile.GuardAddr, auth, ethAddr, amt)
	ChkTestErr(t, err, "failed to delegate stake")

	log.Info("Query sgn about the delegator to check if it has correct stakes...")
	CheckDelegator(t, transactor, ethAddr, ethAddr, amt)

	log.Info("Query sgn about the candidate to check if it has correct stakes...")
	CheckCandidate(t, transactor, ethAddr, sgnop, amt)

	if isValidator {
		log.Info("Query sgn about the validators to check if it has correct stakes...")
		CheckValidator(t, transactor, sgnop, amt, sdk.Bonded)
	}
}

func CheckDelegator(t *testing.T, transactor *transactor.Transactor, validatorAddr, delegatorAddr mainchain.Addr, expAmt *big.Int) {
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
	ChkTestErr(t, err, "failed to queryDelegator")
	log.Infoln("Query sgn about the validator's delegator:", delegator)
	assert.Equal(t, expectedRes, delegator.String(), "The expected result should be: "+expectedRes)
}

func CheckCandidate(t *testing.T, transactor *transactor.Transactor, ethAddr mainchain.Addr, sgnop string, expAmt *big.Int) {
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
	ChkTestErr(t, err, "failed to queryCandidate")
	log.Infoln("Query sgn about the validator candidate:", candidate)
	assert.Equal(t, expectedRes, candidate.String(), "The expected result should be: "+expectedRes)
}

func CheckValidator(t *testing.T, transactor *transactor.Transactor, sgnop string, expAmt *big.Int, expStatus sdk.BondStatus) {
	var validator stypes.Validator
	var err error
	for retry := 0; retry < 30; retry++ {
		validator, err = sgnval.CLIQueryValidator(transactor.CliCtx, staking.RouterKey, sgnop)
		if err == nil && validator.Status == expStatus && validator.Tokens.BigInt().Cmp(expAmt) == 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}
	ChkTestErr(t, err, "failed to queryValidator")
	log.Infoln("Query sgn about the validator:\n", validator)
	assert.Equal(t, expAmt.String(), validator.Tokens.String(), "validator token should be "+expAmt.String())
	assert.Equal(t, expStatus, validator.Status, "validator should be "+sdkStatusName(validator.Status))
}

func CheckValidatorStatus(t *testing.T, transactor *transactor.Transactor, sgnop string, expStatus sdk.BondStatus) {
	var validator stypes.Validator
	var err error
	for retry := 0; retry < 30; retry++ {
		validator, err = sgnval.CLIQueryValidator(transactor.CliCtx, staking.RouterKey, sgnop)
		if err == nil && validator.Status == expStatus {
			break
		}
		time.Sleep(2 * time.Second)
	}
	ChkTestErr(t, err, "failed to queryValidator")
	log.Infoln("Query sgn about the validator:\n", validator)
	assert.Equal(t, expStatus, validator.Status, "validator should be "+sdkStatusName(validator.Status))
}

func CheckValidatorNum(t *testing.T, transactor *transactor.Transactor, expNum int) {
	var validators stypes.Validators
	var err error
	for retry := 0; retry < 30; retry++ {
		validators, err = sgnval.CLIQueryBondedValidators(transactor.CliCtx, staking.RouterKey)
		if err == nil && len(validators) == expNum {
			break
		}
		time.Sleep(2 * time.Second)
	}
	ChkTestErr(t, err, "failed to queryValidators")
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