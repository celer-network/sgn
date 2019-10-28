package e2e

import (
	"context"
	"math/big"
	"testing"

	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/celer-network/sgn/x/subscribe"
)

func setUpSubscribe() []tf.Killable {
	p := &SGNParams{
		blameTimeout:           big.NewInt(10),
		minValidatorNum:        big.NewInt(0),
		minStakingPool:         big.NewInt(0),
		sidechainGoLiveTimeout: big.NewInt(0),
	}
	res := setupNewSGNEnv(p, "subscribe")
	sleepWithLog(20, "sgn being ready")

	return res
}

func TestE2ESubscribe(t *testing.T) {
	toKill := setUpSubscribe()
	defer tf.TearDown(toKill)

	t.Run("e2e-subscribe", func(t *testing.T) {
		t.Run("subscribeTest", subscribeTest)
	})
}

func subscribeTest(t *testing.T) {
	// TODO: each test cases need a new and isolated sgn right now, which can't be run in parallel
	// t.Parallel()

	log.Info("=====================================================================")
	log.Info("======================== Test subscribe ===========================")

	ctx := context.Background()

	conn := tf.EthClient.Client
	auth := tf.EthClient.Auth
	ethAddress := tf.EthClient.Address
	guardContract := tf.EthClient.Guard
	// ledgerContract := tf.EthClient.Ledger
	celrContract, err := mainchain.NewERC20(ctype.Hex2Addr(MockCelerAddr), conn)
	tf.ChkErr(err, "NewERC20 error")

	transactor := tf.Transactor
	// sgnAddr, err := sdk.AccAddressFromBech32(client0SGNAddrStr)
	tf.ChkErr(err, "Parse SGN address error")

	// Call subscribe on guard contract
	log.Info("Call subscribe on guard contract...")
	// TODO: use a separate subscriber address. Now still use the same address of validator
	amt := new(big.Int)
	amt.SetString("100000000000000000000", 10) // 100 CELR
	tx, err := celrContract.Approve(auth, ctype.Hex2Addr(GuardAddr), amt)
	tf.ChkErr(err, "failed to approve CELR to Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")
	tx, err = guardContract.Subscribe(auth, amt)
	tf.ChkErr(err, "failed to call subscribe of Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, maxBlockDiff+1, "Subscribe on Guard contract")

	// Send tx on sidechain to sync mainchain subscription balance
	log.Info("Send tx on sidechain to sync mainchain subscription balance...")
	msg := subscribe.NewMsgSubscribe(ethAddress.String(), transactor.Key.GetAddress())
	transactor.BroadcastTx(msg)
	sleepWithLog(10, "sgn syncing Subscribe balance from mainchain")

	// Query sgn about the subscription info
	log.Info("Query sgn about the subscription info...")
	subscription, err := subscribe.CLIQuerySubscription(transactor.CliCtx.Codec, transactor.CliCtx, subscribe.RouterKey, ethAddress.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Query sgn about the subscription info:", subscription.String())
	// TODO: add expectedRes check after merging the change of pay per use
	// expectedRes = fmt.Sprintf(`EthAddress: %s, DelegatedStake: %d`, ethAddress.String(), amt) // defined in Delegator.String()
	// assert.Equal(t, delegator.String(), expectedRes, fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Query sgn to check if epoch has correct fee
	// TODO: add this test after merging the change of pay per use

	// Query sgn about validator reward
	// TODO: add this test after merging the change of pay per use
}
