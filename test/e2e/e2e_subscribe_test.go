package e2e

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func setUpSubscribe() []tf.Killable {
	p := &SGNParams{
		blameTimeout:           big.NewInt(10),
		minValidatorNum:        big.NewInt(0),
		minStakingPool:         big.NewInt(0),
		sidechainGoLiveTimeout: big.NewInt(0),
	}
	res := setupNewSGNEnv(p, "subscribe")
	sleepWithLog(10, "sgn being ready")

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

	log.Infoln("===================================================================")
	log.Infoln("======================== Test subscribe ===========================")

	ctx := context.Background()
	conn := tf.EthClient.Client
	auth := tf.EthClient.Auth
	ethAddress := tf.EthClient.Address
	guardContract := tf.EthClient.Guard
	ledgerContract := tf.EthClient.Ledger
	transactor := tf.Transactor
	client1PrivKey, _ := crypto.HexToECDSA(client1Priv)
	client1Auth := bind.NewKeyedTransactor(client1PrivKey)
	client1Auth.GasPrice = big.NewInt(2e9) // 2Gwei
	sgnAddr, err := sdk.AccAddressFromBech32(client0SGNAddrStr)
	tf.ChkTestErr(t, err, "failed to parse sgn address")

	err = initializeCandidate(auth, sgnAddr)
	tf.ChkTestErr(t, err, "failed to initialize candidate")
	amt := new(big.Int)
	amt.SetString("100000000000000000000", 10) // 100 CELR
	err = delegateStake(auth, ethAddress, amt)
	tf.ChkTestErr(t, err, "failed to delegate stake")

	log.Infoln("Call subscribe on guard contract...")
	tx, err := celrContract.Approve(auth, guardAddr, amt)
	tf.ChkTestErr(t, err, "failed to approve CELR to Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")
	tx, err = guardContract.Subscribe(auth, amt)
	tf.ChkTestErr(t, err, "failed to call subscribe of Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, blockDelay, "Subscribe on Guard contract")

	log.Infoln("Send tx on sidechain to sync mainchain subscription balance...")
	msgSubscribe := subscribe.NewMsgSubscribe(ethAddress.Hex(), transactor.Key.GetAddress())
	transactor.AddTxMsg(msgSubscribe)
	sleepWithLog(10, "sgn syncing Subscribe balance from mainchain")

	log.Infoln("Query sgn about the subscription info...")
	subscription, err := subscribe.CLIQuerySubscription(transactor.CliCtx, subscribe.RouterKey, ethAddress.Hex())
	tf.ChkTestErr(t, err, "failed to query subscription on sgn")
	log.Infoln("Query sgn about the subscription info:", subscription.String())
	expectedRes := fmt.Sprintf(`Deposit: %d, Spend: %d`, amt, 0) // defined in Subscription.String()
	assert.Equal(t, expectedRes, subscription.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Query sgn to check if epoch has correct fee
	// TODO: add this test after merging the change of pay per use

	// Query sgn about validator reward
	// TODO: add this test after merging the change of pay per use

	log.Infoln("Prepare for requesting guard...")
	channelId, err := openChannel(ethAddress.Bytes(), mainchain.Hex2Bytes(client1AddrStr), tf.EthClient.PrivateKey, client1PrivKey)
	tf.ChkTestErr(t, err, "failed to open channel")
	signedSimplexStateProto, err := prepareSignedSimplexState(10, channelId[:], ethAddress.Bytes(), tf.EthClient.PrivateKey, client1PrivKey)
	tf.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateBytes, err := protobuf.Marshal(signedSimplexStateProto)
	tf.ChkTestErr(t, err, "failed to get signedSimplexStateBytes")
	msgRequestGuard := subscribe.NewMsgRequestGuard(ethAddress.Hex(), signedSimplexStateBytes, transactor.Key.GetAddress())
	transactor.AddTxMsg(msgRequestGuard)
	sleepWithLog(5, "sgn processes request guard")

	log.Infoln("Query sgn to check if request has correct state proof data...")
	request, err := subscribe.CLIQueryRequest(transactor.CliCtx, subscribe.RouterKey, channelId[:])
	tf.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	// TxHash now should be empty
	expectedRes = fmt.Sprintf(`SeqNum: %d, PeerAddresses: [%s %s], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: , GuardTxHash:`, 10, client0AddrStr, client1AddrStr, 0, signedSimplexStateBytes)
	assert.Equal(t, strings.ToLower(expectedRes), strings.ToLower(request.String()), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Call intendSettle on ledger contract...")
	signedSimplexStateProto, err = prepareSignedSimplexState(1, channelId[:], ethAddress.Bytes(), tf.EthClient.PrivateKey, client1PrivKey)
	tf.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateArrayBytes, err := protobuf.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{signedSimplexStateProto},
	})
	tf.ChkTestErr(t, err, "failed to get signedSimplexStateArrayBytes")
	tx, err = ledgerContract.IntendSettle(auth, signedSimplexStateArrayBytes)
	tf.ChkTestErr(t, err, "failed to IntendSettle")
	tf.WaitMinedWithChk(ctx, conn, tx, blockDelay, "IntendSettle")

	log.Infoln("Query sgn to check if validator has submitted the state proof correctly...")
	sleepWithLog(15, "sgn submitting state proof")
	request, err = subscribe.CLIQueryRequest(transactor.CliCtx, subscribe.RouterKey, channelId[:])
	tf.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	rstr := fmt.Sprintf(`SeqNum: %d, PeerAddresses: \[%s %s\], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: 0x[a-f0-9]{64}, GuardTxHash: 0x[a-f0-9]{64}`, 10, client0AddrStr, client1AddrStr, 0, signedSimplexStateBytes)
	r, err := regexp.Compile(strings.ToLower(rstr))
	tf.ChkTestErr(t, err, "failed to compile regexp")
	assert.True(t, r.MatchString(strings.ToLower(request.String())), "SGN query result is wrong")

	log.Infoln("Query sgn to check if it gets the correct reward info (without sigs)...")
	params, err := subscribe.CLIQueryParams(transactor.CliCtx, subscribe.RouterKey)
	tf.ChkTestErr(t, err, "failed to query params on sgn")
	log.Infoln("Query sgn about the params info:", params.String())
	reward, err := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddress.Hex())
	tf.ChkTestErr(t, err, "failed to query reward on sgn")
	log.Infoln("Query sgn about the reward info:", reward.String())
	expectedRes = fmt.Sprintf(`MiningReward: %d, ServiceReward: %s`, 0, params.RequestCost.String())
	assert.Equal(t, expectedRes, reward.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Send tx on sidechain to withdraw reward...")
	msgWithdrawReward := validator.NewMsgWithdrawReward(ethAddress.Hex(), transactor.Key.GetAddress())
	transactor.AddTxMsg(msgWithdrawReward)
	sleepWithLog(15, "sgn withdrawing reward")

	log.Infoln("Query sgn to check if reward gets signature...")
	reward, err = validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddress.Hex())
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

func prepareSignedSimplexState(seqNum uint64, channelId, peerFrom []byte, prvtKey0, prvtKey1 *ecdsa.PrivateKey) (*chain.SignedSimplexState, error) {
	simplexPaymentChannelBytes, err := protobuf.Marshal(&entity.SimplexPaymentChannel{
		SeqNum:    seqNum,
		ChannelId: channelId,
		PeerFrom:  peerFrom,
	})
	if err != nil {
		return nil, err
	}

	sig0, err := mainchain.SignMessage(prvtKey0, simplexPaymentChannelBytes)
	if err != nil {
		return nil, err
	}

	sig1, err := mainchain.SignMessage(prvtKey1, simplexPaymentChannelBytes)
	if err != nil {
		return nil, err
	}

	signedSimplexStateProto := &chain.SignedSimplexState{
		SimplexState: simplexPaymentChannelBytes,
		Sigs:         [][]byte{sig0, sig1},
	}

	return signedSimplexStateProto, nil
}

func monitorOpenChannel(ledgerContract *mainchain.CelerLedger, channelIdChan chan [32]byte) {
	openChannelChan := make(chan *mainchain.CelerLedgerOpenChannel)
	sub, err := ledgerContract.WatchOpenChannel(nil, openChannelChan, nil, nil)
	if err != nil {
		log.Errorln("WatchInitializeCandidate err: ", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Errorln("WatchInitializeCandidate err: ", err)
		case openChannel := <-openChannelChan:
			log.Infoln("Monitored a OpenChannel event")
			channelId := [32]byte{}
			copy(channelId[:], openChannel.ChannelId[:])
			channelIdChan <- channelId
			return
		}
	}
}
