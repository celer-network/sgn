package multinode

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func setUpSubscribe() {
	log.Infoln("set up new sgn env")
	p := &tf.SGNParams{
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MinStakingPool:         big.NewInt(0),
		SidechainGoLiveTimeout: big.NewInt(0),
		CelrAddr:               tf.E2eProfile.CelrAddr,
	}
	setupNewSGNEnv(p)
	addThreeValidators(big.NewInt(1000000000000000000))
	tf.SleepWithLog(10, "sgn syncing")
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
	conn := tf.DefaultTestEthClient.Client
	auth := tf.DefaultTestEthClient.Auth
	ethAddress := tf.DefaultTestEthClient.Address
	guardContract := tf.DefaultTestEthClient.Guard
	ledgerContract := tf.DefaultTestEthClient.Ledger
	transactor := tf.NewTransactor(
		t,
		sgnCLIHomes[0],
		sgnChainID,
		sgnNodeURIs[0],
		sgnTransactors[0],
		sgnPassphrase,
		sgnGasPrice,
	)
	Client1PrivKey, _ := crypto.HexToECDSA(tf.Client1Priv)
	client1Auth := bind.NewKeyedTransactor(Client1PrivKey)
	client1Auth.GasPrice = big.NewInt(2e9) // 2Gwei
	validatorNum := 3

	// err = tf.InitializeCandidate(auth, sgnAddr)
	// tf.ChkTestErr(t, err, "failed to initialize candidate")
	// err = tf.DelegateStake(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, ethAddress, amt)
	// tf.ChkTestErr(t, err, "failed to delegate stake")

	log.Infoln("Call subscribe on guard contract...")
	amt := new(big.Int)
	amt.SetString("100000000000000000000", 10) // 100 CELR
	tx, err := tf.E2eProfile.CelrContract.Approve(auth, tf.E2eProfile.GuardAddr, amt)
	tf.ChkTestErr(t, err, "failed to approve CELR to Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, tf.BlockDelay, "Approve CELR to Guard contract")
	tx, err = guardContract.Subscribe(auth, amt)
	tf.ChkTestErr(t, err, "failed to call subscribe of Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, tf.BlockDelay, "Subscribe on Guard contract")

	log.Infoln("Send tx on sidechain to sync mainchain subscription balance...")
	msgSubscribe := subscribe.NewMsgSubscribe(ethAddress.Hex(), transactor.Key.GetAddress())
	transactor.AddTxMsg(msgSubscribe)
	tf.SleepWithLog(10, "sgn syncing Subscribe balance from mainchain")

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
	channelId, err := tf.OpenChannel(ethAddress.Bytes(), mainchain.Hex2Bytes(tf.Client1AddrStr), tf.DefaultTestEthClient.PrivateKey, Client1PrivKey, tf.E2eProfile.CelrAddr.Bytes())
	tf.ChkTestErr(t, err, "failed to open channel")
	tf.SleepWithLog(10, "wait channelId to be in secure state")
	signedSimplexStateProto, err := tf.PrepareSignedSimplexState(10, channelId[:], ethAddress.Bytes(), tf.DefaultTestEthClient.PrivateKey, Client1PrivKey)
	tf.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateBytes, err := protobuf.Marshal(signedSimplexStateProto)
	tf.ChkTestErr(t, err, "failed to get signedSimplexStateBytes")
	msgRequestGuard := subscribe.NewMsgRequestGuard(ethAddress.Hex(), signedSimplexStateBytes, transactor.Key.GetAddress())
	transactor.AddTxMsg(msgRequestGuard)
	tf.SleepWithLog(5, "sgn processes request guard")

	log.Infoln("Query sgn to check if request has correct state proof data...")
	request, err := subscribe.CLIQueryRequest(transactor.CliCtx, subscribe.RouterKey, channelId[:])
	tf.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	// TxHash now should be empty
	expectedRes = fmt.Sprintf(`SeqNum: %d, PeerAddresses: [%s %s], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: , GuardTxHash:`, 10, tf.Client0AddrStr, tf.Client1AddrStr, 0, signedSimplexStateBytes)
	assert.Equal(t, strings.ToLower(expectedRes), strings.ToLower(request.String()), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Call intendSettle on ledger contract...")
	signedSimplexStateProto, err = tf.PrepareSignedSimplexState(1, channelId[:], ethAddress.Bytes(), tf.DefaultTestEthClient.PrivateKey, Client1PrivKey)
	tf.ChkTestErr(t, err, "failed to prepare SignedSimplexState")
	signedSimplexStateArrayBytes, err := protobuf.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{signedSimplexStateProto},
	})
	tf.ChkTestErr(t, err, "failed to get signedSimplexStateArrayBytes")
	tx, err = ledgerContract.IntendSettle(auth, signedSimplexStateArrayBytes)
	tf.ChkTestErr(t, err, "failed to IntendSettle")
	tf.WaitMinedWithChk(ctx, conn, tx, tf.BlockDelay, "IntendSettle")

	log.Infoln("Query sgn to check if validator has submitted the state proof correctly...")
	tf.SleepWithLog(15, "sgn submitting state proof")
	request, err = subscribe.CLIQueryRequest(transactor.CliCtx, subscribe.RouterKey, channelId[:])
	tf.ChkTestErr(t, err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	rstr := fmt.Sprintf(`SeqNum: %d, PeerAddresses: \[%s %s\], PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: 0x[a-f0-9]{64}, GuardTxHash: 0x[a-f0-9]{64}`, 10, tf.Client0AddrStr, tf.Client1AddrStr, 0, signedSimplexStateBytes)
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
	expectedReward := params.RequestCost.QuoRaw(int64(validatorNum))
	expectedRes = fmt.Sprintf(`MiningReward: %d, ServiceReward: %s`, 0, expectedReward.String())
	assert.Equal(t, expectedRes, reward.String(), fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	log.Infoln("Send tx on sidechain to withdraw reward...")
	msgWithdrawReward := validator.NewMsgWithdrawReward(ethAddress.Hex(), transactor.Key.GetAddress())
	transactor.AddTxMsg(msgWithdrawReward)
	tf.SleepWithLog(15, "sgn withdrawing reward")

	log.Infoln("Query sgn to check if reward gets signature...")
	reward, err = validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddress.Hex())
	tf.ChkTestErr(t, err, "failed to query reward on sgn")
	assert.Equal(t, validatorNum, len(reward.Sigs), fmt.Sprintf("The length of reward signatures should be %d", validatorNum))

	log.Infoln("Call redeemReward on guard contract...")
	tx, err = guardContract.RedeemReward(auth, reward.GetRewardRequest())
	tf.ChkTestErr(t, err, "failed to redeem reward")
	tf.WaitMinedWithChk(ctx, conn, tx, tf.BlockDelay, "redeem reward on Guard contract")
	rsr, err := guardContract.RedeemedServiceReward(&bind.CallOpts{}, ethAddress)
	tf.ChkTestErr(t, err, "failed to query redeemed service reward")
	assert.Equal(t, reward.ServiceReward.BigInt(), rsr, "reward is not redeemed")
}
