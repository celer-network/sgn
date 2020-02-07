package e2ecommon

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
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/subscribe"
	stypes "github.com/celer-network/sgn/x/subscribe/types"
	"github.com/celer-network/sgn/x/validator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func SubscribteTestCommon(t *testing.T, transactor *transactor.Transactor, amt *big.Int, srvReward string, rewardSigLen int) {
	ctx := context.Background()

	log.Infoln("Open channel...")
	channelId, err := tc.OpenChannel(tc.Client0, tc.Client1)
	tc.ChkTestErr(t, err, "failed to open channel")

	log.Infoln("Call subscribe on guard contract...")
	tx, err := tc.E2eProfile.CelrContract.Approve(tc.Client0.Auth, tc.E2eProfile.GuardAddr, new(big.Int).Mul(amt, big.NewInt(2)))
	tc.ChkTestErr(t, err, "failed to approve CELR to Guard contract")
	tc.WaitMinedWithChk(ctx, tc.Client0.Client, tx, 0, "Approve CELR to Guard contract")

	_, err = tc.EtherBase.Guard.ContributeToMiningPool(tc.Client0.Auth, amt)
	tc.ChkTestErr(t, err, "failed to call ContributeToMiningPool of Guard contract")

	tx, err = tc.Client0.Guard.Subscribe(tc.Client0.Auth, amt)
	tc.ChkTestErr(t, err, "failed to call subscribe of Guard contract")
	tc.WaitMinedWithChk(ctx, tc.Client0.Client, tx, tc.BlockDelay, "Subscribe on Guard contract")
	tc.SleepWithLog(20, "passing subscribe event block delay")

	log.Infoln("Send tx on sidechain to sync mainchain subscription balance...")
	msgSubscribe := subscribe.NewMsgSubscribe(tc.Client0.Address.Hex(), transactor.Key.GetAddress())
	transactor.AddTxMsg(msgSubscribe)

	log.Infoln("Query sgn about the subscription info...")
	var subscription stypes.Subscription
	expectedRes := fmt.Sprintf(`Deposit: %d, Spend: %d`, amt, 0) // defined in Subscription.String()
	for retry := 0; retry < 30; retry++ {
		subscription, err = subscribe.CLIQuerySubscription(transactor.CliCtx, subscribe.RouterKey, tc.Client0.Address.Hex())
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
	signedSimplexStateProto, err := tc.PrepareSignedSimplexState(10, channelId[:], tc.Client0.Address.Bytes(), tc.Client0, tc.Client1)
	tc.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateBytes, err := proto.Marshal(signedSimplexStateProto)
	tc.ChkTestErr(t, err, "failed to get signedSimplexStateBytes")
	msgRequestGuard := subscribe.NewMsgRequestGuard(tc.Client0.Address.Hex(), signedSimplexStateBytes, transactor.Key.GetAddress())
	transactor.AddTxMsg(msgRequestGuard)

	log.Infoln("Query sgn to check if request has correct state proof data...")
	var request stypes.Request
	// TxHash now should be empty
	expectedRes = fmt.Sprintf(`SeqNum: %d, PeerAddresses: [%s %s], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: , GuardTxHash:`, 10, tc.ClientEthAddrs[1], tc.ClientEthAddrs[0], 1, signedSimplexStateBytes)
	for retry := 0; retry < 30; retry++ {
		request, err = subscribe.CLIQueryRequest(transactor.CliCtx, subscribe.RouterKey, channelId[:], tc.Client0.Address.Hex())
		if err == nil && expectedRes == request.String() {
			break
		}
		time.Sleep(time.Second)
	}
	tc.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	assert.Equal(t, strings.ToLower(expectedRes), strings.ToLower(request.String()), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Call intendSettle on ledger contract...")
	signedSimplexStateProto, err = tc.PrepareSignedSimplexState(1, channelId[:], tc.Client0.Address.Bytes(), tc.Client0, tc.Client1)
	tc.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateArrayBytes, err := proto.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{signedSimplexStateProto},
	})
	tc.ChkTestErr(t, err, "failed to get signedSimplexStateArrayBytes")
	tx, err = tc.Client0.Ledger.IntendSettle(tc.Client0.Auth, signedSimplexStateArrayBytes)
	tc.ChkTestErr(t, err, "failed to IntendSettle")
	tc.WaitMinedWithChk(ctx, tc.Client0.Client, tx, tc.BlockDelay, "IntendSettle")

	log.Infoln("Query sgn to check if validator has submitted the state proof correctly...")
	rstr := fmt.Sprintf(`SeqNum: %d, PeerAddresses: \[%s %s\], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: 0x[a-f0-9]{64}, GuardTxHash: 0x[a-f0-9]{64}`, 10, tc.ClientEthAddrs[1], tc.ClientEthAddrs[0], 1, signedSimplexStateBytes)
	r, err := regexp.Compile(strings.ToLower(rstr))
	tc.ChkTestErr(t, err, "failed to compile regexp")
	for retry := 0; retry < 60; retry++ {
		request, err = subscribe.CLIQueryRequest(transactor.CliCtx, subscribe.RouterKey, channelId[:], tc.Client0.Address.Hex())
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
	reward, err := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, tc.ValEthAddrs[0])
	tc.ChkTestErr(t, err, "failed to query reward on sgn")
	log.Infoln("Query sgn about the reward info:", reward.String())
	assert.True(t, reward.MiningReward.IsPositive(), "Minging reward should be larger than 0")
	if srvReward == "" {
		srvReward = params.RequestCost.String()
	}
	expectedRes = fmt.Sprintf(`Receiver: %s, MiningReward: %s, ServiceReward: %s`, tc.ValEthAddrs[0], reward.MiningReward.String(), srvReward)
	assert.Equal(t, expectedRes, reward.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Send tx on sidechain to withdraw reward...")
	msgWithdrawReward := validator.NewMsgWithdrawReward(tc.ValEthAddrs[0], transactor.Key.GetAddress())
	transactor.AddTxMsg(msgWithdrawReward)

	log.Infoln("Query sgn to check if reward gets signature...")
	for retry := 0; retry < 30; retry++ {
		reward, err = validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, tc.ValEthAddrs[0])
		if err == nil && len(reward.Sigs) == rewardSigLen {
			break
		}
		time.Sleep(time.Second)
	}
	tc.ChkTestErr(t, err, "failed to query reward on sgn")
	assert.Equal(t, rewardSigLen, len(reward.Sigs), "The length of reward signatures mismatch")

	log.Infoln("Call redeemReward on guard contract...")
	tx, err = tc.Client0.Guard.RedeemReward(tc.Client0.Auth, reward.GetRewardRequest())
	tc.ChkTestErr(t, err, "failed to redeem reward")
	tc.WaitMinedWithChk(ctx, tc.Client0.Client, tx, tc.BlockDelay, "redeem reward on Guard contract")
	rsr, err := tc.Client0.Guard.RedeemedServiceReward(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[0]))
	tc.ChkTestErr(t, err, "failed to query redeemed service reward")
	assert.Equal(t, reward.ServiceReward.BigInt(), rsr, "reward is not redeemed")
}
