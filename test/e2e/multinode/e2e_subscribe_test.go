package multinode

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/subscribe"
	stypes "github.com/celer-network/sgn/x/subscribe/types"
	"github.com/celer-network/sgn/x/validator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func setUpSubscribe() {
	log.Infoln("set up new sgn env")
	p := &tc.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MinStakingPool:         big.NewInt(0),
		SidechainGoLiveTimeout: big.NewInt(0),
		CelrAddr:               tc.E2eProfile.CelrAddr,
		MaxValidatorNum:        big.NewInt(11),
	}
	setupNewSGNEnv(p)
}

func TestE2ESubscribe(t *testing.T) {
	setUpSubscribe()

	t.Run("e2e-subscribe", func(t *testing.T) {
		t.Run("subscribeTest", subscribeTest)
	})
}

func subscribeTest(t *testing.T) {
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
		tc.SgnCLIHome,
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnCLIAddr,
		tc.SgnPassphrase,
		tc.SgnGasPrice,
	)
	Client1PrivKey, err := tc.GetEthPrivateKey(tc.ValEthKs[1])
	tc.ChkTestErr(t, err, "failed to get client 1 private key")

	log.Infoln("Add validators...")
	amts := []*big.Int{big.NewInt(1000000000000000000), big.NewInt(1000000000000000000), big.NewInt(100000000000000000)}
	tc.AddValidators(t, transactor, tc.ValEthKs[:], tc.SgnOperators[:], amts)
	turnOffMonitor(2)

	log.Infoln("Open channel...")
	channelId, err := tc.OpenChannel(ethAddress, mainchain.Hex2Addr(tc.ValEthAddrs[1]), privKey, Client1PrivKey)
	tc.ChkTestErr(t, err, "failed to open channel")

	log.Infoln("Call subscribe on guard contract...")
	amt := new(big.Int)
	amt.SetString("1"+strings.Repeat("0", 20), 10) // 100 CELR
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
	signedSimplexStateProto, err := tc.PrepareSignedSimplexState(10, channelId[:], ethAddress.Bytes(), tc.DefaultTestEthClient.PrivateKey, Client1PrivKey)
	tc.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateBytes, err := protobuf.Marshal(signedSimplexStateProto)
	tc.ChkTestErr(t, err, "failed to get signedSimplexStateBytes")
	msgRequestGuard := subscribe.NewMsgRequestGuard(ethAddress.Hex(), signedSimplexStateBytes, transactor.Key.GetAddress())
	transactor.AddTxMsg(msgRequestGuard)

	log.Infoln("Query sgn to check if request has correct state proof data...")
	var request stypes.Request
	// TxHash now should be empty
	expectedRes = fmt.Sprintf(`SeqNum: %d, PeerAddresses: [%s %s], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: , GuardTxHash:`, 10, tc.ValEthAddrs[0], tc.ValEthAddrs[1], 0, signedSimplexStateBytes)
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
	signedSimplexStateProto, err = tc.PrepareSignedSimplexState(1, channelId[:], ethAddress.Bytes(), tc.DefaultTestEthClient.PrivateKey, Client1PrivKey)
	tc.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateArrayBytes, err := protobuf.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{signedSimplexStateProto},
	})
	tc.ChkTestErr(t, err, "failed to get signedSimplexStateArrayBytes")
	tx, err = ledgerContract.IntendSettle(auth, signedSimplexStateArrayBytes)
	tc.ChkTestErr(t, err, "failed to IntendSettle")
	tc.WaitMinedWithChk(ctx, conn, tx, tc.BlockDelay, "IntendSettle")

	log.Infoln("Query sgn to check if validator has submitted the state proof correctly...")
	rstr := fmt.Sprintf(`SeqNum: %d, PeerAddresses: \[%s %s\], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: 0x[a-f0-9]{64}, GuardTxHash: 0x[a-f0-9]{64}`, 10, tc.ValEthAddrs[0], tc.ValEthAddrs[1], 0, signedSimplexStateBytes)
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
	expectedRes = fmt.Sprintf(`Receiver: %s, MiningReward: %s, ServiceReward: %s`, mainchain.Addr2Hex(ethAddress), reward.MiningReward.String(), "476190476190476190")
	assert.Equal(t, expectedRes, reward.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Send tx on sidechain to withdraw reward...")
	msgWithdrawReward := validator.NewMsgWithdrawReward(ethAddress.Hex(), transactor.Key.GetAddress())
	transactor.AddTxMsg(msgWithdrawReward)

	log.Infoln("Query sgn to check if reward gets signature...")
	for retry := 0; retry < 30; retry++ {
		reward, err = validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddress.Hex())
		if err == nil && len(reward.Sigs) == 2 {
			break
		}
		time.Sleep(time.Second)
	}
	tc.ChkTestErr(t, err, "failed to query reward on sgn")
	assert.Equal(t, 2, len(reward.Sigs), "The length of reward signatures should be 2")

	log.Infoln("Call redeemReward on guard contract...")
	tx, err = guardContract.RedeemReward(auth, reward.GetRewardRequest())
	tc.ChkTestErr(t, err, "failed to redeem reward")
	tc.WaitMinedWithChk(ctx, conn, tx, tc.BlockDelay, "redeem reward on Guard contract")
	rsr, err := guardContract.RedeemedServiceReward(&bind.CallOpts{}, ethAddress)
	tc.ChkTestErr(t, err, "failed to query redeemed service reward")
	assert.Equal(t, reward.ServiceReward.BigInt(), rsr, "reward is not redeemed")

	log.Infoln("Query sgn to check penalty")
	nonce := uint64(0)
	penalty, err := slash.CLIQueryPenalty(transactor.CliCtx, slash.StoreKey, nonce)
	tc.ChkTestErr(t, err, "failed to query penalty")
	expectedRes = fmt.Sprintf(`Nonce: %d, ValidatorAddr: %s, Reason: guard_failure`, nonce, tc.ValEthAddrs[2])
	assert.Equal(t, expectedRes, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	expectedRes = fmt.Sprintf(`Account: %s, Amount: 1000000000000000`, tc.ValEthAddrs[2])
	assert.Equal(t, expectedRes, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
	assert.Equal(t, 2, len(penalty.Sigs), fmt.Sprintf("The length of validators should be 2"))

	log.Infoln("Query onchain staking pool")
	var poolAmt string
	for retry := 0; retry < 30; retry++ {
		ci, _ := tc.DefaultTestEthClient.Guard.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[2]))
		poolAmt = ci.StakingPool.String()
		if poolAmt == "99000000000000000" {
			break
		}
		time.Sleep(time.Second)
	}
	assert.Equal(t, "99000000000000000", poolAmt, fmt.Sprintf("The expected StakingPool should be 99000000000000000"))

}
