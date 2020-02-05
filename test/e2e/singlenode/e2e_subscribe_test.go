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
	tc "github.com/celer-network/sgn/test/common"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/celer-network/sgn/x/subscribe"
	stypes "github.com/celer-network/sgn/x/subscribe/types"
	"github.com/celer-network/sgn/x/validator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setUpSubscribe() []tc.Killable {
	p := &tc.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MinStakingPool:         big.NewInt(0),
		SidechainGoLiveTimeout: big.NewInt(0),
		CelrAddr:               tc.E2eProfile.CelrAddr,
		MaxValidatorNum:        big.NewInt(11),
	}
	res := setupNewSGNEnv(p, "subscribe")
	tc.SleepWithLog(10, "sgn being ready")

	return res
}

func TestE2ESubscribe(t *testing.T) {
	toKill := setUpSubscribe()
	defer tc.TearDown(toKill)

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
	conn := tc.DefaultTestEthClient.Client
	auth := tc.DefaultTestEthClient.Auth
	ethAddress := tc.DefaultTestEthClient.Address
	guardContract := tc.DefaultTestEthClient.Guard
	ledgerContract := tc.DefaultTestEthClient.Ledger
	privKey := tc.DefaultTestEthClient.PrivateKey
	transactor := tc.NewTransactor(
		t,
		CLIHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetStringSlice(common.FlagSgnTransactors)[1],
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetString(common.FlagSgnGasPrice),
	)
	Client1PrivKey, err := tc.GetEthPrivateKey(tc.EthKeystores[1])
	tc.ChkTestErr(t, err, "failed to get client 1 private key")

	amt := new(big.Int)
	amt.SetString("100000000000000000000", 10) // 100 CELR
	tc.AddCandidateWithStake(t, transactor, ethAddress, auth, tc.SgnOperators[0], amt, big.NewInt(1), true)

	log.Infoln("Open channel...")
	channelId, err := tc.OpenChannel(ethAddress, mainchain.Hex2Addr(tc.EthAddresses[1]), privKey, Client1PrivKey)
	tc.ChkTestErr(t, err, "failed to open channel")

	log.Infoln("Call subscribe on guard contract...")
	tx, err := tc.E2eProfile.CelrContract.Approve(auth, tc.E2eProfile.GuardAddr, new(big.Int).Mul(amt, big.NewInt(2)))
	tc.ChkTestErr(t, err, "failed to approve CELR to Guard contract")
	tc.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")

	_, err = guardContract.ContributeToMiningPool(auth, amt)
	tc.ChkTestErr(t, err, "failed to call ContributeToMiningPool of Guard contract")

	tx, err = guardContract.Subscribe(auth, amt)
	tc.ChkTestErr(t, err, "failed to call subscribe of Guard contract")
	tc.WaitMinedWithChk(ctx, conn, tx, tc.BlockDelay, "Subscribe on Guard contract")
	tc.SleepWithLog(20, "passing subscribe event block delay")

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
	tc.ChkTestErr(t, err, "failed to query subscription on sgn")
	log.Infoln("Query sgn about the subscription info:", subscription.String())
	assert.Equal(t, expectedRes, subscription.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Query sgn to check if epoch has correct fee
	// TODO: add this test after merging the change of pay per use

	// Query sgn about validator reward
	// TODO: add this test after merging the change of pay per use

	log.Infoln("Request guard...")
	signedSimplexStateProto, err := tc.PrepareSignedSimplexState(10, channelId[:], ethAddress.Bytes(), privKey, Client1PrivKey)
	tc.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateBytes, err := protobuf.Marshal(signedSimplexStateProto)
	tc.ChkTestErr(t, err, "failed to get signedSimplexStateBytes")
	msgRequestGuard := subscribe.NewMsgRequestGuard(ethAddress.Hex(), signedSimplexStateBytes, transactor.Key.GetAddress())
	transactor.AddTxMsg(msgRequestGuard)

	log.Infoln("Query sgn to check if request has correct state proof data...")
	var request stypes.Request
	// TxHash now should be empty
	expectedRes = fmt.Sprintf(`SeqNum: %d, PeerAddresses: [%s %s], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: , GuardTxHash:`, 10, tc.EthAddresses[0], tc.EthAddresses[1], 0, signedSimplexStateBytes)
	for retry := 0; retry < 30; retry++ {
		request, err = subscribe.CLIQueryRequest(transactor.CliCtx, subscribe.RouterKey, channelId[:], ethAddress.Hex())
		if err == nil && expectedRes == request.String() {
			break
		}
		time.Sleep(time.Second)
	}
	tc.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	assert.Equal(t, strings.ToLower(expectedRes), strings.ToLower(request.String()), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Call intendSettle on ledger contract...")
	signedSimplexStateProto, err = tc.PrepareSignedSimplexState(1, channelId[:], ethAddress.Bytes(), privKey, Client1PrivKey)
	tc.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateArrayBytes, err := protobuf.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{signedSimplexStateProto},
	})
	tc.ChkTestErr(t, err, "failed to get signedSimplexStateArrayBytes")
	tx, err = ledgerContract.IntendSettle(auth, signedSimplexStateArrayBytes)
	tc.ChkTestErr(t, err, "failed to IntendSettle")
	tc.WaitMinedWithChk(ctx, conn, tx, tc.BlockDelay, "IntendSettle")

	log.Infoln("Query sgn to check if validator has submitted the state proof correctly...")
	rstr := fmt.Sprintf(`SeqNum: %d, PeerAddresses: \[%s %s\], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: 0x[a-f0-9]{64}, GuardTxHash: 0x[a-f0-9]{64}`, 10, tc.EthAddresses[0], tc.EthAddresses[1], 0, signedSimplexStateBytes)
	r, err := regexp.Compile(strings.ToLower(rstr))
	tc.ChkTestErr(t, err, "failed to compile regexp")
	for retry := 0; retry < 60; retry++ {
		request, err = subscribe.CLIQueryRequest(transactor.CliCtx, subscribe.RouterKey, channelId[:], ethAddress.Hex())
		if err == nil && r.MatchString(strings.ToLower(request.String())) {
			break
		}
		time.Sleep(2 * time.Second)
	}
	tc.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	assert.True(t, r.MatchString(strings.ToLower(request.String())), "SGN query result is wrong")

	log.Infoln("Query sgn to check if it gets the correct reward info (without sigs)...")
	params, err := subscribe.CLIQueryParams(transactor.CliCtx, subscribe.RouterKey)
	tc.ChkTestErr(t, err, "failed to query params on sgn")
	log.Infoln("Query sgn about the params info:", params.String())
	reward, err := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddress.Hex())
	tc.ChkTestErr(t, err, "failed to query reward on sgn")
	log.Infoln("Query sgn about the reward info:", reward.String())
	assert.True(t, reward.MiningReward.IsPositive(), "Minging reward should be larger than 0")
	expectedRes = fmt.Sprintf(`Receiver: %s, MiningReward: %s, ServiceReward: %s`, mainchain.Addr2Hex(ethAddress), reward.MiningReward.String(), params.RequestCost.String())
	assert.Equal(t, expectedRes, reward.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Send tx on sidechain to withdraw reward...")
	msgWithdrawReward := validator.NewMsgWithdrawReward(ethAddress.Hex(), transactor.Key.GetAddress())
	transactor.AddTxMsg(msgWithdrawReward)
	tc.SleepWithLog(15, "sgn withdrawing reward")

	log.Infoln("Query sgn to check if reward gets signature...")
	for retry := 0; retry < 30; retry++ {
		reward, err = validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddress.Hex())
		if err == nil && len(reward.Sigs) == 1 {
			break
		}
		time.Sleep(time.Second)
	}
	tc.ChkTestErr(t, err, "failed to query reward on sgn")
	assert.Equal(t, 1, len(reward.Sigs), "The length of reward signatures should be 1")

	log.Infoln("Call redeemReward on guard contract...")
	tx, err = guardContract.RedeemReward(auth, reward.GetRewardRequest())
	tc.ChkTestErr(t, err, "failed to redeem reward")
	tc.WaitMinedWithChk(ctx, conn, tx, 0, "redeem reward on Guard contract")
	rsr, err := guardContract.RedeemedServiceReward(&bind.CallOpts{}, ethAddress)
	tc.ChkTestErr(t, err, "failed to query redeemed service reward")
	assert.Equal(t, reward.ServiceReward.BigInt(), rsr, "reward is not redeemed")
}
