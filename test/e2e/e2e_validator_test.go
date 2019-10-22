package e2e

import (
	"context"
	"math/big"
	"testing"

	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setUpValidator() []tf.Killable {
	p := &SGNParams{
		blameTimeout:           big.NewInt(10),
		minValidatorNum:        big.NewInt(1),
		minStakingPool:         big.NewInt(1),
		sidechainGoLiveTimeout: big.NewInt(0),
	}
	res := setupNewSGNEnv(p, "validator")
	log.Info("Sleep 40 seconds for sgn syncing")
	sleep(40)

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
	t.Parallel()

	log.Info("=====================================================================")
	log.Info("======================== Test validator ===========================")

	ctx := context.Background()

	conn := tf.EthClient.Client
	auth := tf.EthClient.Auth
	ethAddress := tf.EthClient.Address
	guardContract := tf.EthClient.Guard
	// ledgerContract := tf.EthClient.Ledger

	transactor := tf.Transactor
	sgnAddr, err := sdk.AccAddressFromBech32(client0SGNAddrStr)
	tf.ChkErr(err, "Parse SGN address error")

	// Call initializeCandidate on guard contract
	tx, err := guardContract.InitializeCandidate(auth, big.NewInt(1), sgnAddr.Bytes())
	tf.ChkErr(err, "failed to InitializeCandidate")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "InitializeCandidate")
	log.Info("Sleep 40 seconds for sgn syncing")
	sleep(40)

	// query sgn about the validator candidate
	candidate, err := validator.CLIQueryCandidate(transactor.CliCtx.Codec, transactor.CliCtx, validator.RouterKey, ethAddress.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("query sgn about the validator candidate", candidate)
}
