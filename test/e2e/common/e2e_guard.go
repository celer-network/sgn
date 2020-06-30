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
	"github.com/celer-network/sgn/x/guard"
	stypes "github.com/celer-network/sgn/x/guard/types"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func SubscribteTestCommon(t *testing.T, transactor *transactor.Transactor, amt *big.Int, srvReward string, rewardSigLen int) {
	ctx := context.Background()

	log.Infoln("Open channel...")
	channelId, err := tc.OpenChannel(tc.Client0, tc.Client1)
	tc.ChkTestErr(t, err, "failed to open channel")

	log.Infoln("Call subscribe on sgn contract...")

	tx, err := tc.E2eProfile.CelrContract.Approve(tc.Client0.Auth, tc.E2eProfile.DPoSAddr, amt)
	tc.ChkTestErr(t, err, "failed to approve CELR to DPoS contract")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Approve CELR to DPoS contract")

	tx, err = tc.E2eProfile.CelrContract.Approve(tc.Client0.Auth, tc.E2eProfile.SGNAddr, amt)
	tc.ChkTestErr(t, err, "failed to approve CELR to SGN contract")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Approve CELR to SGN contract")

	_, err = tc.DposContract.ContributeToMiningPool(tc.Client0.Auth, amt)
	tc.ChkTestErr(t, err, "failed to call ContributeToMiningPool of DPoS contract")

	tx, err = tc.SgnContract.Subscribe(tc.Client0.Auth, amt)
	tc.ChkTestErr(t, err, "failed to call subscribe of SGN contract")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Subscribe on SGN contract")

	log.Infoln("Send tx on sidechain to sync mainchain subscription balance...")
	subscription := stypes.NewSubscription(tc.Client0.Address.Hex())
	subscription.Deposit = sdk.NewIntFromBigInt(amt)
	subscriptionData := transactor.CliCtx.Codec.MustMarshalBinaryBare(subscription)
	msgSubmitChange := sync.NewMsgSubmitChange(sync.Subscribe, subscriptionData, transactor.Key.GetAddress())
	transactor.AddTxMsg(msgSubmitChange)

	log.Infoln("Query sgn about the subscription info...")
	expectedRes := fmt.Sprintf(`EthAddress: %s, Deposit: %d, Spend: %d`, mainchain.Addr2Hex(tc.Client0.Address), amt, 0) // defined in Subscription.String()
	for retry := 0; retry < tc.RetryLimit; retry++ {
		subscription, err = guard.CLIQuerySubscription(transactor.CliCtx, guard.RouterKey, tc.Client0.Address.Hex())
		if err == nil && expectedRes == subscription.String() {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	tc.ChkTestErr(t, err, "failed to query subscription on sgn")
	log.Infoln("Query sgn about the subscription info:", subscription.String())
	assert.Equal(t, expectedRes, subscription.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Query sgn to check if epoch has correct fee
	// TODO: add this test after merging the change of pay per use

	// Query sgn about validator reward
	// TODO: add this test after merging the change of pay per use

	log.Infoln("Request to init guard ...")
	seqNum := uint64(10)
	signedSimplexStateProto, err := tc.PrepareSignedSimplexState(seqNum, channelId[:], tc.Client1.Address.Bytes(), tc.Client0, tc.Client1)
	tc.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateBytes, err := proto.Marshal(signedSimplexStateProto)
	tc.ChkTestErr(t, err, "failed to get signedSimplexStateBytes")
	requestSig, err := tc.Client0.Signer.SignEthMessage(signedSimplexStateBytes)
	tc.ChkTestErr(t, err, "failed to sign signedSimplexStateBytes")
	initRequest := guard.NewInitRequest(signedSimplexStateBytes, requestSig, tc.DisputeTimeout)
	requestData := transactor.CliCtx.Codec.MustMarshalBinaryBare(initRequest)
	msgSubmitChange = sync.NewMsgSubmitChange(sync.InitGuardRequest, requestData, transactor.Key.GetAddress())
	transactor.AddTxMsg(msgSubmitChange)

	log.Infoln("Query sgn to check if request has correct state proof data...")
	// TxHash now should be empty
	var request guard.Request
	expectedRes = fmt.Sprintf(`SeqNum: %d, SimplexSender: %s, SimplexReceiver: %s, DisputeTimeout: %d, TriggerTxHash: , TriggerTxBlkNum: 0, GuardTxHash: , GuardTxBlkNum: 0, GuardSender:`, seqNum, tc.ClientEthAddrs[1], tc.ClientEthAddrs[0], tc.DisputeTimeout)
	for retry := 0; retry < tc.RetryLimit; retry++ {
		request, err = guard.CLIQueryRequest(transactor.CliCtx, guard.RouterKey, channelId[:], tc.Client0.Address.Hex())
		if err == nil && expectedRes == request.String() {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	tc.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	assert.Equal(t, strings.ToLower(expectedRes), strings.ToLower(request.String()), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Request guard (2nd request)...")
	seqNum = uint64(12)
	signedSimplexStateProto, err = tc.PrepareSignedSimplexState(seqNum, channelId[:], tc.Client1.Address.Bytes(), tc.Client0, tc.Client1)
	tc.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateBytes, err = proto.Marshal(signedSimplexStateProto)
	tc.ChkTestErr(t, err, "failed to get signedSimplexStateBytes")
	requestSig, err = tc.Client0.Signer.SignEthMessage(signedSimplexStateBytes)
	tc.ChkTestErr(t, err, "failed to sign signedSimplexStateBytes")
	msgRequestGuard := guard.NewMsgRequestGuard(signedSimplexStateBytes, requestSig, transactor.Key.GetAddress())
	transactor.AddTxMsg(msgRequestGuard)

	log.Infoln("Query sgn to check if request has correct state proof data...")
	expectedRes = fmt.Sprintf(`SeqNum: %d, SimplexSender: %s, SimplexReceiver: %s, DisputeTimeout: %d, TriggerTxHash: , TriggerTxBlkNum: 0, GuardTxHash: , GuardTxBlkNum: 0, GuardSender:`, seqNum, tc.ClientEthAddrs[1], tc.ClientEthAddrs[0], tc.DisputeTimeout)
	for retry := 0; retry < tc.RetryLimit; retry++ {
		request, err = guard.CLIQueryRequest(transactor.CliCtx, guard.RouterKey, channelId[:], tc.Client0.Address.Hex())
		if err == nil && expectedRes == request.String() {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	tc.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	assert.Equal(t, strings.ToLower(expectedRes), strings.ToLower(request.String()), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Call intendSettle on ledger contract...")
	signedSimplexStateProto, err = tc.PrepareSignedSimplexState(1, channelId[:], tc.Client1.Address.Bytes(), tc.Client0, tc.Client1)
	tc.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateArrayBytes, err := proto.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{signedSimplexStateProto},
	})
	tc.ChkTestErr(t, err, "failed to get signedSimplexStateArrayBytes")
	tx, err = tc.LedgerContract.IntendSettle(tc.Client0.Auth, signedSimplexStateArrayBytes)
	tc.ChkTestErr(t, err, "failed to IntendSettle")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "IntendSettle")

	log.Infoln("Query sgn to check if validator has submitted the state proof correctly...")
	rstr := fmt.Sprintf(`SeqNum: %d, SimplexSender: %s, SimplexReceiver: %s, DisputeTimeout: %d, TriggerTxHash: 0x[a-f0-9]{64}, TriggerTxBlkNum: [0-9]{2,3}, GuardTxHash: 0x[a-f0-9]{64}, GuardTxBlkNum: [0-9]{2,3}, GuardSender: [a-f0-9]{40}`, seqNum, tc.ClientEthAddrs[1], tc.ClientEthAddrs[0], tc.DisputeTimeout)
	r, err := regexp.Compile(strings.ToLower(rstr))
	tc.ChkTestErr(t, err, "failed to compile regexp")
	for retry := 0; retry < tc.RetryLimit; retry++ {
		request, err = guard.CLIQueryRequest(transactor.CliCtx, guard.RouterKey, channelId[:], tc.Client0.Address.Hex())
		if err == nil && r.MatchString(strings.ToLower(request.String())) {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	tc.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	assert.True(t, r.MatchString(strings.ToLower(request.String())), "SGN query result is wrong")

	log.Infoln("Query sgn to check if it gets the correct reward info (without sigs)...")
	params, err := guard.CLIQueryParams(transactor.CliCtx, guard.RouterKey)
	tc.ChkTestErr(t, err, "failed to query params on sgn")
	log.Infoln("Query sgn about the params info:", params.String())
	reward, err := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, tc.ValEthAddrs[0])
	tc.ChkTestErr(t, err, "failed to query reward on sgn")
	log.Infoln("Query sgn about the reward info:", reward.String())
	assert.True(t, reward.MiningReward.IsPositive(), "Minging reward should be larger than 0")
	if srvReward == "" {
		srvReward = params.RequestCost.Add(params.RequestCost).String()
	}
	expectedRes = fmt.Sprintf(`Receiver: %s, MiningReward: %s, ServiceReward: %s`, tc.ValEthAddrs[0], reward.MiningReward.String(), srvReward)
	assert.Equal(t, expectedRes, reward.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Send tx on sidechain to withdraw reward...")
	msgWithdrawReward := validator.NewMsgWithdrawReward(tc.ValEthAddrs[0], transactor.Key.GetAddress())
	transactor.AddTxMsg(msgWithdrawReward)

	log.Infoln("Query sgn to check if reward gets signature...")
	for retry := 0; retry < tc.RetryLimit; retry++ {
		reward, err = validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, tc.ValEthAddrs[0])
		if err == nil && len(reward.Sigs) == rewardSigLen {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	tc.ChkTestErr(t, err, "failed to query reward on sgn")
	assert.Equal(t, rewardSigLen, len(reward.Sigs), "The length of reward signatures mismatch")

	log.Infoln("Call redeemReward on sgn contract...")
	tx, err = tc.SgnContract.RedeemReward(tc.Client0.Auth, reward.GetRewardRequest())
	tc.ChkTestErr(t, err, "failed to redeem reward")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "redeem reward on SGN contract")
	rsr, err := tc.SgnContract.RedeemedServiceReward(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[0]))
	tc.ChkTestErr(t, err, "failed to query redeemed service reward")
	assert.Equal(t, reward.ServiceReward.BigInt(), rsr, "reward is not redeemed")
}
