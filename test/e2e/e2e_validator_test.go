package e2e

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/mainchain"
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

	ctx := context.Background()

	conn := tf.EthClient.Client
	auth := tf.EthClient.Auth
	ethAddress := tf.EthClient.Address
	guardContract := tf.EthClient.Guard
	// ledgerContract := tf.EthClient.Ledger
	celrContract, err := mainchain.NewERC20(ctype.Hex2Addr(MockCelerAddr), conn)
	tf.ChkErr(err, "NewERC20 error")

	transactor := tf.Transactor
	sgnAddr, err := sdk.AccAddressFromBech32(client0SGNAddrStr)
	tf.ChkErr(err, "Parse SGN address error")

	// Call initializeCandidate on guard contract using the validator eth address
	log.Info("Call initializeCandidate on guard contract using the validator eth address...")
	tx, err := guardContract.InitializeCandidate(auth, big.NewInt(1), sgnAddr.Bytes())
	tf.ChkErr(err, "failed to InitializeCandidate")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "InitializeCandidate")
	sleepWithLog(30, "sgn syncing InitializeCandidate event on mainchain")

	// Query sgn about the validator candidate
	log.Info("Query sgn about the validator candidate...")
	candidate, err := validator.CLIQueryCandidate(transactor.CliCtx.Codec, transactor.CliCtx, validator.RouterKey, ethAddress.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Query sgn about the validator candidate:", candidate)
	expectedRes := fmt.Sprintf(`Operator: %s, StakingPool: %d`, client0SGNAddrStr, 0) // defined in Candidate.String()
	assert.Equal(t, candidate.String(), expectedRes, fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Call delegate on guard contract to delegate stake to the validator eth address
	log.Info("Call delegate on guard contract to delegate stake to the validator eth address...")
	amt := new(big.Int)
	amt.SetString("100", 10)
	tx, err = celrContract.Approve(auth, ctype.Hex2Addr(GuardAddr), amt)
	tf.ChkErr(err, "failed to approve CELR to Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")
	tx, err = guardContract.Delegate(auth, ethAddress, amt)
	tf.ChkErr(err, "failed to call delegate of Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Delegate to validator")
	sleepWithLog(30, "sgn syncing Delegate event on mainchain")

	// Query sgn about the delegator to check if it has correct stakes
	log.Info("Query sgn about the delegator to check if it has correct stakes...")
	delegator, err := validator.CLIQueryDelegator(transactor.CliCtx.Codec, transactor.CliCtx, validator.RouterKey, ethAddress.String(), ethAddress.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Query sgn about the validator delegator:", delegator)
	expectedRes = fmt.Sprintf(`EthAddress: %s, DelegatedStake: %d`, ethAddress.String(), amt) // defined in Delegator.String()
	assert.Equal(t, delegator.String(), expectedRes, fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Query sgn about the candidate to check if it has correct stakes
	log.Info("Query sgn about the validator candidate...")
	candidate, err = validator.CLIQueryCandidate(transactor.CliCtx.Codec, transactor.CliCtx, validator.RouterKey, ethAddress.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Query sgn about the validator candidate:", candidate)
	expectedRes = fmt.Sprintf(`Operator: %s, StakingPool: %d`, client0SGNAddrStr, amt) // defined in Candidate.String()
	assert.Equal(t, candidate.String(), expectedRes, fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Query sgn about the validator to check if it has correct stakes
	sleepWithLog(30, "wait for validator to claimValidator and sgn sync ValidatorChange event")
	log.Info("Query sgn about the validator to check if it has correct stakes...")
	validators, err := validator.CLIQueryValidators(transactor.CliCtx.Codec, transactor.CliCtx, staking.RouterKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Query sgn about the validators:", validators)
	// TODO: use a better way to assert/check the validity of the lengthy query results.
	// expectedRes = fmt.Sprintf("StakingPool: %d", amt) // defined in Candidate.String()
	// assert.Equal(t, validators.String(), expectedRes, fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
}
