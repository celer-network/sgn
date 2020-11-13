package common

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/guard"
	guardtypes "github.com/celer-network/sgn/x/guard/types"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// common functions for single-node and multi-node e2e guard tests

func Subscribe(t *testing.T, transactor *transactor.Transactor, amt *big.Int) {
	ctx := context.Background()
	var tx *ethtypes.Transaction

	log.Infoln("Call subscribe on sgn contract...")
	tx, err := tc.E2eProfile.CelrContract.Approve(tc.Client0.Auth, tc.E2eProfile.DPoSAddr, amt)
	require.NoError(t, err, "failed to approve CELR to DPoS contract")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Approve CELR to DPoS contract")

	tx, err = tc.E2eProfile.CelrContract.Approve(tc.Client0.Auth, tc.E2eProfile.SGNAddr, amt)
	require.NoError(t, err, "failed to approve CELR to SGN contract")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Approve CELR to SGN contract")

	_, err = tc.DposContract.ContributeToMiningPool(tc.Client0.Auth, amt)
	require.NoError(t, err, "failed to call ContributeToMiningPool of DPoS contract")

	tx, err = tc.SgnContract.Subscribe(tc.Client0.Auth, amt)
	require.NoError(t, err, "failed to call subscribe of SGN contract")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Subscribe on SGN contract")

	subscription := guardtypes.NewSubscription(tc.Client0.Address.Hex())
	subscription.Deposit = sdk.NewIntFromBigInt(amt)
	log.Infoln("Query sgn about the subscription info...")
	expectedRes := fmt.Sprintf(`EthAddress: %s, Deposit: %d, Spend: %d`, mainchain.Addr2Hex(tc.Client0.Address), amt, 0)
	for retry := 0; retry < tc.RetryLimit; retry++ {
		subscription, err = guard.CLIQuerySubscription(transactor.CliCtx, guard.RouterKey, tc.Client0.Address.Hex())
		if err == nil && expectedRes == subscription.String() {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	require.NoError(t, err, "failed to query subscription on sgn")
	log.Infoln("Query sgn about the subscription info:", subscription.String())
	assert.Equal(t, expectedRes, subscription.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
}

func TestGuard(t *testing.T, transactor *transactor.Transactor, guardSender string) {
	ctx := context.Background()
	reqsubstr := fmt.Sprintf(`SimplexSender: %s, SimplexReceiver: %s, DisputeTimeout: %d`,
		tc.ClientEthAddrs[1], tc.ClientEthAddrs[0], tc.DisputeTimeout)

	log.Infoln("Open channel...")
	cid, err := tc.OpenChannel(tc.Client0, tc.Client1)
	require.NoError(t, err, "failed to open channel")

	log.Infoln("Request to init guard ...")
	seqNum := uint64(10)
	msgSubmitChange := getGuardMsg(t, transactor, cid, seqNum, true)
	transactor.AddTxMsg(msgSubmitChange)

	log.Infoln("Query sgn to check if request has correct state proof data...")

	expectedRes := fmt.Sprintf(`ChannelId: %x, SeqNum: %d, %s, Status: Idle`, cid, seqNum, reqsubstr)
	checkRequest(t, transactor, cid, expectedRes)

	log.Infoln("Request guard (2nd request)...")
	seqNum = uint64(12)
	msgRequestGuard := getGuardMsg(t, transactor, cid, seqNum, false)
	transactor.AddTxMsg(msgRequestGuard)

	log.Infoln("Query sgn to check if request has correct state proof data...")
	expectedRes = fmt.Sprintf(`ChannelId: %x, SeqNum: %d, %s, Status: Idle`, cid, seqNum, reqsubstr)
	checkRequest(t, transactor, cid, expectedRes)

	log.Infoln("Call intendSettle on ledger contract...")
	signedSimplexStateProto, err := tc.PrepareSignedSimplexState(1, cid[:], tc.Client1.Address.Bytes(), tc.Client0, tc.Client1)
	require.NoError(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateArrayBytes, err := proto.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{signedSimplexStateProto},
	})
	require.NoError(t, err, "failed to get signedSimplexStateArrayBytes")
	tx, err := tc.LedgerContract.IntendSettle(tc.Client0.Auth, signedSimplexStateArrayBytes)
	require.NoError(t, err, "failed to IntendSettle")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "IntendSettle")

	log.Infoln("Query sgn to check if validator has submitted the state proof correctly...")
	expectedRes = fmt.Sprintf(`ChannelId: %x, SeqNum: %d, %s, Status: Settled, TriggerTxHash: %s, TriggerTxBlkNum: [0-9]{2,3}, GuardTxHash: 0x[a-f0-9]{64}, GuardTxBlkNum: [0-9]{2,3}, GuardSender: %s`,
		cid, seqNum, reqsubstr, tx.Hash().Hex(), guardSender)
	checkRequest(t, transactor, cid, expectedRes)
}

func CheckReward(t *testing.T, transactor *transactor.Transactor, valAddr, srvReward string, rewardSigLen int) {
	ctx := context.Background()
	log.Infoln("Query sgn to check if it gets the correct reward info (without sigs)...")
	reward, err := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, valAddr)
	require.NoError(t, err, "failed to query reward on sgn")
	log.Infoln("Query sgn about the reward info:", reward.String())
	assert.True(t, reward.MiningReward.IsPositive(), "Minging reward should be larger than 0")
	expectedRes := fmt.Sprintf(`Receiver: %s, MiningReward: %s, ServiceReward: %s`,
		valAddr, reward.MiningReward.String(), srvReward)
	assert.Equal(t, expectedRes, reward.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Send tx on sidechain to withdraw reward...")
	msgWithdrawReward := validator.NewMsgWithdrawReward(valAddr, transactor.Key.GetAddress())
	transactor.AddTxMsg(msgWithdrawReward)

	log.Infoln("Query sgn to check if reward gets signature...")
	for retry := 0; retry < tc.RetryLimit; retry++ {
		reward, err = validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, valAddr)
		if err == nil && len(reward.Sigs) == rewardSigLen {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	require.NoError(t, err, "failed to query reward on sgn")
	assert.Equal(t, rewardSigLen, len(reward.Sigs), "The length of reward signatures mismatch")

	log.Infoln("Call redeemReward on sgn contract...")
	tx, err := tc.SgnContract.RedeemReward(tc.Client0.Auth, reward.GetRewardRequest())
	require.NoError(t, err, "failed to redeem reward")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "redeem reward on SGN contract")
	rsr, err := tc.SgnContract.RedeemedServiceReward(&bind.CallOpts{}, mainchain.Hex2Addr(valAddr))
	require.NoError(t, err, "failed to query redeemed service reward")
	assert.Equal(t, reward.ServiceReward.BigInt(), rsr, "reward is not redeemed")
}

func getGuardMsg(t *testing.T, transactor *transactor.Transactor, cid mainchain.CidType, seqNum uint64, init bool) sdk.Msg {
	signedSimplexStateProto, err := tc.PrepareSignedSimplexState(seqNum, cid[:], tc.Client1.Address.Bytes(), tc.Client0, tc.Client1)
	require.NoError(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateBytes, err := proto.Marshal(signedSimplexStateProto)
	require.NoError(t, err, "failed to get signedSimplexStateBytes")
	requestSig, err := tc.Client0.Signer.SignEthMessage(signedSimplexStateBytes)
	require.NoError(t, err, "failed to sign signedSimplexStateBytes")
	if init {
		initRequest := guard.NewInitRequest(signedSimplexStateBytes, requestSig, tc.DisputeTimeout)
		syncData := transactor.CliCtx.Codec.MustMarshalBinaryBare(initRequest)
		return sync.NewMsgSubmitChange(sync.InitGuardRequest, syncData, tc.EthClient, transactor.Key.GetAddress())
	}
	return guard.NewMsgRequestGuard(signedSimplexStateBytes, requestSig, transactor.Key.GetAddress())
}

func checkRequest(t *testing.T, transactor *transactor.Transactor, cid mainchain.CidType, expectedRes string) {
	var request guard.Request
	rexp, err := regexp.Compile(expectedRes)
	require.NoError(t, err, "failed to compile regexp")
	for retry := 0; retry < tc.RetryLimit; retry++ {
		request, err = guard.CLIQueryRequest(transactor.CliCtx, guard.RouterKey, cid[:], tc.Client0.Address.Hex())
		if err == nil && rexp.MatchString(request.String()) {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	require.NoError(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	assert.True(t, rexp.MatchString(request.String()), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
}
