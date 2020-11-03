package common

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/gov"
	govtypes "github.com/celer-network/sgn/x/gov/types"
	"github.com/celer-network/sgn/x/slash"
	sgnval "github.com/celer-network/sgn/x/validator"
	vtypes "github.com/celer-network/sgn/x/validator/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SGNParams struct {
	CelrAddr              mainchain.Addr
	GovernProposalDeposit *big.Int
	GovernVoteTimeout     *big.Int
	SlashTimeout          *big.Int
	MinValidatorNum       *big.Int
	MaxValidatorNum       *big.Int
	MinStakingPool        *big.Int
	AdvanceNoticePeriod   *big.Int
	// TODO: rename to DposGoLiveTimeout
	SidechainGoLiveTimeout *big.Int
	MinGasPrices           string
	StartGateway           bool
}

func SetupSidechain() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(common.Bech32PrefixValAddr, common.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(common.Bech32PrefixConsAddr, common.Bech32PrefixConsPub)
	config.Seal()
}

func NewTestTransactor(t *testing.T, sgnCLIHome, sgnChainID, sgnNodeURI, sgnValAcct, sgnPassphrase string) *transactor.Transactor {
	cdc := app.MakeCodec()
	tr, err := transactor.NewTransactor(
		sgnCLIHome,
		sgnChainID,
		sgnNodeURI,
		sgnValAcct,
		sgnPassphrase,
		cdc,
	)
	require.NoError(t, err, "Failed to create new transactor.")
	tr.Run()

	return tr
}

func AddValidators(t *testing.T, transactor *transactor.Transactor, ethkss, valaccts []string, amts []*big.Int) {
	for i := 0; i < len(ethkss); i++ {
		log.Infoln("Adding validator", i)
		ethAddr, auth, err := GetAuth(ethkss[i])
		require.NoError(t, err, "failed to get auth")
		AddCandidateWithStake(t, transactor, ethAddr, auth, valaccts[i], amts[i], big.NewInt(1), big.NewInt(1), big.NewInt(10000), true)
	}
}

func AddCandidateWithStake(t *testing.T, transactor *transactor.Transactor,
	ethAddr mainchain.Addr, auth *bind.TransactOpts,
	valacct string, amt *big.Int, minAmt *big.Int, commissionRate *big.Int,
	rateLockEndTime *big.Int, checkValidator bool) {

	// get sgnAddr
	sgnAddr, err := sdk.AccAddressFromBech32(valacct)
	require.NoError(t, err, "failed to parse sgn address")

	// add candidate
	err = InitializeCandidate(auth, sgnAddr, minAmt, commissionRate, rateLockEndTime)
	require.NoError(t, err, "failed to initialize candidate")

	log.Infof("Query sgn about the validator candidate %s ...", ethAddr.Hex())
	CheckCandidate(t, transactor, ethAddr, valacct, big.NewInt(0))

	// self delegate stake
	err = DelegateStake(auth, ethAddr, amt)
	require.NoError(t, err, "failed to delegate stake")

	log.Info("Query sgn about the delegator to check if it has correct stakes...")
	CheckDelegator(t, transactor, ethAddr, ethAddr, amt)

	log.Info("Query sgn about the candidate to check if it has correct stakes...")
	CheckCandidate(t, transactor, ethAddr, valacct, amt)

	if checkValidator {
		log.Infof("Query sgn about the validator %s to check if it has correct stakes...", valacct)
		CheckValidator(t, transactor, valacct, amt, sdk.Bonded)
	}
}

func CheckDelegator(t *testing.T, transactor *transactor.Transactor, validatorAddr, delegatorAddr mainchain.Addr, expAmt *big.Int) {
	var delegator vtypes.Delegator
	var err error
	expectedRes := fmt.Sprintf(`CandidateAddr: %s, DelegatorAddr: %s, DelegatedStake: %s`,
		mainchain.Addr2Hex(validatorAddr), mainchain.Addr2Hex(delegatorAddr), expAmt) // defined in Delegator.String()
	for retry := 0; retry < RetryLimit; retry++ {
		delegator, err = sgnval.CLIQueryDelegator(transactor.CliCtx, sgnval.RouterKey, validatorAddr.Hex(), delegatorAddr.Hex())
		if err == nil && expectedRes == delegator.String() {
			break
		}
		time.Sleep(RetryPeriod)
	}
	require.NoError(t, err, "failed to queryDelegator")
	log.Infoln("Query sgn about the validator's delegator:", delegator)
	assert.Equal(t, expectedRes, delegator.String(), "The expected result should be: "+expectedRes)
}

func CheckCandidate(t *testing.T, transactor *transactor.Transactor, ethAddr mainchain.Addr, valacct string, expAmt *big.Int) {
	var candidate vtypes.Candidate
	var err error
	expectedRes := fmt.Sprintf(`ValAccount: %s, EthAddress: %x, StakingPool: %s`, valacct, ethAddr, expAmt) // defined in Candidate.String()
	for retry := 0; retry < RetryLimit; retry++ {
		candidate, err = sgnval.CLIQueryCandidate(transactor.CliCtx, sgnval.RouterKey, ethAddr.Hex())
		if err != nil {
			log.Debugln("retry due to err:", err)
		}
		if err == nil && expectedRes == candidate.String() {
			break
		}
		time.Sleep(RetryPeriod)
	}
	require.NoError(t, err, "failed to queryCandidate", err)
	log.Infoln("Query sgn about the validator candidate:", candidate)
	assert.Equal(t, expectedRes, candidate.String(), "The expected result should be: "+expectedRes)
}

func CheckValidator(t *testing.T, transactor *transactor.Transactor, valacct string, expAmt *big.Int, expStatus sdk.BondStatus) {
	var validator stypes.Validator
	var err error
	for retry := 0; retry < RetryLimit; retry++ {
		validator, err = sgnval.CLIQueryValidator(transactor.CliCtx, staking.RouterKey, valacct)
		if err == nil &&
			validator.Status == expStatus {
			expToken := sdk.NewIntFromBigInt(expAmt).QuoRaw(common.TokenDec).String()
			if expToken == validator.Tokens.String() {
				break
			}
		}
		time.Sleep(RetryPeriod)
	}
	require.NoError(t, err, "failed to queryValidator")
	expToken := sdk.NewIntFromBigInt(expAmt).QuoRaw(common.TokenDec).String()
	assert.Equal(t, expToken, validator.Tokens.String(), "validator token should be "+expToken)
	assert.Equal(t, expStatus, validator.Status, "validator should be "+sdkStatusName(validator.Status))
}

func CheckValidatorStatus(t *testing.T, transactor *transactor.Transactor, valacct string, expStatus sdk.BondStatus) {
	var validator stypes.Validator
	var err error
	for retry := 0; retry < RetryLimit; retry++ {
		validator, err = sgnval.CLIQueryValidator(transactor.CliCtx, staking.RouterKey, valacct)
		if err == nil && validator.Status == expStatus {
			break
		}
		time.Sleep(RetryPeriod)
	}
	require.NoError(t, err, "failed to queryValidator")
	assert.Equal(t, expStatus, validator.Status, "validator should be "+sdkStatusName(validator.Status))
}

func CheckValidatorNum(t *testing.T, transactor *transactor.Transactor, expNum int) {
	var validators stypes.Validators
	var err error
	for retry := 0; retry < RetryLimit; retry++ {
		validators, err = sgnval.CLIQueryBondedValidators(transactor.CliCtx, staking.RouterKey)
		if err == nil && len(validators) == expNum {
			break
		}
		time.Sleep(RetryPeriod)
	}
	require.NoError(t, err, "failed to queryValidators")
	assert.Equal(t, expNum, len(validators), "The length of validators should be: "+strconv.Itoa(expNum))
}

func QueryProposal(cliCtx context.CLIContext, proposalID uint64, status govtypes.ProposalStatus) (proposal govtypes.Proposal, err error) {
	for retry := 0; retry < RetryLimit; retry++ {
		proposal, err = gov.CLIQueryProposal(cliCtx, gov.RouterKey, proposalID)
		if err == nil && status == proposal.Status {
			break
		}
		time.Sleep(RetryPeriod)
	}

	if err != nil {
		return
	}

	if status != proposal.Status {
		err = fmt.Errorf("Proposal status %s does not match expectation", status)
	}

	return
}

func QueryPenalty(cliCtx context.CLIContext, nonce uint64, sigCount int) (penalty slash.Penalty, err error) {
	for retry := 0; retry < RetryLimit; retry++ {
		penalty, err = slash.CLIQueryPenalty(cliCtx, slash.StoreKey, nonce)
		if err == nil && len(penalty.PenaltyProtoBytes) > 0 && len(penalty.Sigs) == sigCount {
			break
		}
		time.Sleep(RetryPeriod)
	}

	if err != nil {
		return
	}

	if len(penalty.PenaltyProtoBytes) == 0 {
		err = errors.New("PenaltyProtoBytes cannot be zero")
	}

	if len(penalty.Sigs) != sigCount {
		err = errors.New("Signature count does not match expectation")
	}

	return
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
