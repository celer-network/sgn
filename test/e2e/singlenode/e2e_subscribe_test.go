package singlenode

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	tc "github.com/celer-network/sgn/test/e2e/common"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/subscribe"
	stypes "github.com/celer-network/sgn/x/subscribe/types"
	"github.com/celer-network/sgn/x/validator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setUpSubscribe() []tf.Killable {
	p := &tf.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MinStakingPool:         big.NewInt(0),
		SidechainGoLiveTimeout: big.NewInt(0),
		CelrAddr:               tf.E2eProfile.CelrAddr,
		MaxValidatorNum:        big.NewInt(11),
	}
	res := setupNewSGNEnv(p, "subscribe")
	tf.SleepWithLog(10, "sgn being ready")

	return res
}

func TestE2ESubscribe(t *testing.T) {
	toKill := setUpSubscribe()
	defer tf.TearDown(toKill)

	t.Run("e2e-subscribe", func(t *testing.T) {
		t.Run("subscribeTest", subscribeTest)
	})
}

// TODO: dedup with multinode tests
func subscribeTest(t *testing.T) {
	// TODO: each test cases need a new and isolated sgn right now, which can't be run in parallel
	// t.Parallel()

	log.Infoln("===================================================================")
	log.Infoln("======================== Test subscribe ===========================")

	ctx := context.Background()
	conn := tf.DefaultTestEthClient.Client
	auth := tf.DefaultTestEthClient.Auth
	ethAddress := tf.DefaultTestEthClient.Address
	guardContract := tf.DefaultTestEthClient.Guard
	ledgerContract := tf.DefaultTestEthClient.Ledger
	privKey := tf.DefaultTestEthClient.PrivateKey
	transactor := tf.NewTransactor(
		t,
		CLIHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetStringSlice(common.FlagSgnTransactors)[1],
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetString(common.FlagSgnGasPrice),
	)
	Client1PrivKey, _ := crypto.HexToECDSA(tf.Client1Priv)

	amt := new(big.Int)
	amt.SetString("100000000000000000000", 10) // 100 CELR
	tc.AddCandidateWithStake(t, transactor, ethAddress, auth, tf.Client0SGNAddrStr, amt, big.NewInt(1), true)

	log.Infoln("Open channel...")
	channelId, err := tf.OpenChannel(ethAddress, mainchain.Hex2Addr(tf.Client1AddrStr), privKey, Client1PrivKey)
	tf.ChkTestErr(t, err, "failed to open channel")

	log.Infoln("Call subscribe on guard contract...")
	tx, err := tf.E2eProfile.CelrContract.Approve(auth, tf.E2eProfile.GuardAddr, new(big.Int).Mul(amt, big.NewInt(2)))
	tf.ChkTestErr(t, err, "failed to approve CELR to Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")

	_, err = guardContract.ContributeToMiningPool(auth, amt)
	tf.ChkTestErr(t, err, "failed to call ContributeToMiningPool of Guard contract")

	tx, err = guardContract.Subscribe(auth, amt)
	tf.ChkTestErr(t, err, "failed to call subscribe of Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, tf.BlockDelay, "Subscribe on Guard contract")
	tf.SleepWithLog(20, "passing subscribe event block delay")

	log.Infoln("Send tx on sidechain to sync mainchain subscription balance...")
	msgSubscribe := subscribe.NewMsgSubscribe(ethAddress.Hex(), transactor.Key.GetAddress())
	transactor.AddTxMsg(msgSubscribe)

	log.Infoln("Query sgn about the subscription info...")
	var subscription stypes.Subscription
	expectedRes := fmt.Sprintf(`Deposit: %d, Spend: %d`, amt, 0) // defined in Subscription.String()
	for retry := 0; retry < 30; retry++ {
		subscription, err = subscribe.CLIQuerySubscription(transactor.CliCtx, subscribe.RouterKey, ethAddress.Hex())
		if err == nil && expectedRes == subscription.String() {
			break
		}
		time.Sleep(time.Second)
	}
	tf.ChkTestErr(t, err, "failed to query subscription on sgn")
	log.Infoln("Query sgn about the subscription info:", subscription.String())
	assert.Equal(t, expectedRes, subscription.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Query sgn to check if epoch has correct fee
	// TODO: add this test after merging the change of pay per use

	// Query sgn about validator reward
	// TODO: add this test after merging the change of pay per use

	log.Infoln("Request guard...")
	signedSimplexStateProto, err := tf.PrepareSignedSimplexState(10, channelId[:], ethAddress.Bytes(), privKey, Client1PrivKey)
	tf.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateBytes, err := protobuf.Marshal(signedSimplexStateProto)
	tf.ChkTestErr(t, err, "failed to get signedSimplexStateBytes")
	msgRequestGuard := subscribe.NewMsgRequestGuard(ethAddress.Hex(), signedSimplexStateBytes, transactor.Key.GetAddress())
	transactor.AddTxMsg(msgRequestGuard)

	log.Infoln("Query sgn to check if request has correct state proof data...")
	var request stypes.Request
	// TxHash now should be empty
	expectedRes = fmt.Sprintf(`SeqNum: %d, PeerAddresses: [%s %s], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: , GuardTxHash:`, 10, tf.Client0AddrStr, tf.Client1AddrStr, 0, signedSimplexStateBytes)
	for retry := 0; retry < 30; retry++ {
		request, err = subscribe.CLIQueryRequest(transactor.CliCtx, subscribe.RouterKey, channelId[:], ethAddress.Hex())
		if err == nil && expectedRes == request.String() {
			break
		}
		time.Sleep(time.Second)
	}
	tf.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	assert.Equal(t, strings.ToLower(expectedRes), strings.ToLower(request.String()), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Call intendSettle on ledger contract...")
	signedSimplexStateProto, err = tf.PrepareSignedSimplexState(1, channelId[:], ethAddress.Bytes(), privKey, Client1PrivKey)
	tf.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateArrayBytes, err := protobuf.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{signedSimplexStateProto},
	})
	tf.ChkTestErr(t, err, "failed to get signedSimplexStateArrayBytes")
	tx, err = ledgerContract.IntendSettle(auth, signedSimplexStateArrayBytes)
	tf.ChkTestErr(t, err, "failed to IntendSettle")
	tf.WaitMinedWithChk(ctx, conn, tx, tf.BlockDelay, "IntendSettle")

	log.Infoln("Query sgn to check if validator has submitted the state proof correctly...")
	rstr := fmt.Sprintf(`SeqNum: %d, PeerAddresses: \[%s %s\], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: 0x[a-f0-9]{64}, GuardTxHash: 0x[a-f0-9]{64}`, 10, tf.Client0AddrStr, tf.Client1AddrStr, 0, signedSimplexStateBytes)
	r, err := regexp.Compile(strings.ToLower(rstr))
	tf.ChkTestErr(t, err, "failed to compile regexp")
	for retry := 0; retry < 60; retry++ {
		request, err = subscribe.CLIQueryRequest(transactor.CliCtx, subscribe.RouterKey, channelId[:], ethAddress.Hex())
		if err == nil && r.MatchString(strings.ToLower(request.String())) {
			break
		}
		time.Sleep(2 * time.Second)
	}
	tf.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	assert.True(t, r.MatchString(strings.ToLower(request.String())), "SGN query result is wrong")

	log.Infoln("Query sgn to check if it gets the correct reward info (without sigs)...")
	params, err := subscribe.CLIQueryParams(transactor.CliCtx, subscribe.RouterKey)
	tf.ChkTestErr(t, err, "failed to query params on sgn")
	log.Infoln("Query sgn about the params info:", params.String())
	reward, err := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddress.Hex())
	tf.ChkTestErr(t, err, "failed to query reward on sgn")
	log.Infoln("Query sgn about the reward info:", reward.String())
	assert.True(t, reward.MiningReward.IsPositive(), "Minging reward should be larger than 0")
	expectedRes = fmt.Sprintf(`Receiver: %s, MiningReward: %s, ServiceReward: %s`, mainchain.Addr2Hex(ethAddress), reward.MiningReward.String(), params.RequestCost.String())
	assert.Equal(t, expectedRes, reward.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Send tx on sidechain to withdraw reward...")
	msgWithdrawReward := validator.NewMsgWithdrawReward(ethAddress.Hex(), transactor.Key.GetAddress())
	transactor.AddTxMsg(msgWithdrawReward)
	tf.SleepWithLog(15, "sgn withdrawing reward")

	log.Infoln("Query sgn to check if reward gets signature...")
	for retry := 0; retry < 30; retry++ {
		reward, err = validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddress.Hex())
		if err == nil && len(reward.Sigs) == 1 {
			break
		}
		time.Sleep(time.Second)
	}
	tf.ChkTestErr(t, err, "failed to query reward on sgn")
	assert.Equal(t, 1, len(reward.Sigs), "The length of reward signatures should be 1")

	log.Infoln("Call redeemReward on guard contract...")
	tx, err = guardContract.RedeemReward(auth, reward.GetRewardRequest())
	tf.ChkTestErr(t, err, "failed to redeem reward")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "redeem reward on Guard contract")
	rsr, err := guardContract.RedeemedServiceReward(&bind.CallOpts{}, ethAddress)
	tf.ChkTestErr(t, err, "failed to query redeemed service reward")
	assert.Equal(t, reward.ServiceReward.BigInt(), rsr, "reward is not redeemed")
}
