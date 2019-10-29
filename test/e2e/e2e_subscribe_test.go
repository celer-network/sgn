package e2e

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	protobuf "github.com/golang/protobuf/proto"
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
	ledgerContract := tf.EthClient.Ledger
	celrContract, err := mainchain.NewERC20(ctype.Hex2Addr(MockCelerAddr), conn)
	tf.ChkErr(err, "NewERC20 error")

	client1PrivKey, _ := crypto.HexToECDSA(client1Priv)
	client1Auth := bind.NewKeyedTransactor(client1PrivKey)
	client1Auth.GasPrice = big.NewInt(2e9) // 2Gwei

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
	tf.WaitMinedWithChk(ctx, conn, tx, maxBlockDiff+2, "Subscribe on Guard contract")

	// Send tx on sidechain to sync mainchain subscription balance
	log.Info("Send tx on sidechain to sync mainchain subscription balance...")
	msgSubscribe := subscribe.NewMsgSubscribe(ethAddress.String(), transactor.Key.GetAddress())
	transactor.BroadcastTx(msgSubscribe)
	sleepWithLog(10, "sgn syncing Subscribe balance from mainchain")

	// Query sgn about the subscription info
	log.Info("Query sgn about the subscription info...")
	subscription, err := subscribe.CLIQuerySubscription(transactor.CliCtx.Codec, transactor.CliCtx, subscribe.RouterKey, ethAddress.String())
	tf.ChkErr(err, "failed to query subscription on sgn")
	log.Infoln("Query sgn about the subscription info:", subscription.String())
	// TODO: add expectedRes check after merging the change of pay per use
	// expectedRes = fmt.Sprintf(`EthAddress: %s, DelegatedStake: %d`, ethAddress.String(), amt) // defined in Delegator.String()
	// assert.Equal(t, delegator.String(), expectedRes, fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Query sgn to check if epoch has correct fee
	// TODO: add this test after merging the change of pay per use

	// Query sgn about validator reward
	// TODO: add this test after merging the change of pay per use

	// Call openChannelMockSet on ledger contract
	log.Info("Call openChannel on ledger contract...")
	// copy(channelId[:], []byte{1})
	paymentChannelInitializerBytes, err := protobuf.Marshal(prepareChannelInitializer())
	tf.ChkErr(err, "failed to get paymentChannelInitializerBytes")
	sig0, err := mainchain.SignMessage(tf.EthClient.PrivateKey, paymentChannelInitializerBytes)
	tf.ChkErr(err, "failed to get sig0")
	sig1, err := mainchain.SignMessage(client1PrivKey, paymentChannelInitializerBytes)
	tf.ChkErr(err, "failed to get sig1")
	requestBytes, err := protobuf.Marshal(&chain.OpenChannelRequest{
		ChannelInitializer: paymentChannelInitializerBytes,
		Sigs:               [][]byte{sig0, sig1},
	})
	tf.ChkErr(err, "failed to get requestBytes")
	channelIdChan := make(chan [32]byte)
	go monitorOpenChannel(ledgerContract, channelIdChan)
	tx, err = ledgerContract.OpenChannel(auth, requestBytes)
	tf.ChkErr(err, "failed to OpenChannel")
	tf.WaitMinedWithChk(ctx, conn, tx, maxBlockDiff+2, "OpenChannel")
	// channelId := [32]byte{}
	channelId := <-channelIdChan
	log.Info("channel ID:", "channelId", channelId)
	// log.Info("Call openChannelMockSet on ledger contract...")
	// tx, err = ledgerContract.OpenChannelMockSet(auth,
	// 	channelId,
	// 	big.NewInt(100),
	// 	ctype.Hex2Addr(MockCelerAddr),
	// 	big.NewInt(2),
	// 	[2]ctype.Addr{ethAddress, ctype.Hex2Addr(client1AddrStr)},
	// 	[2]*big.Int{big.NewInt(100), big.NewInt(100)},
	// )
	// tf.ChkErr(err, "failed to OpenChannel")
	// tf.WaitMinedWithChk(ctx, conn, tx, maxBlockDiff+2, "OpenChannel")

	// Submit state proof to sgn
	signedSimplexStateProto := prepareSignedSimplexState(10, channelId[:], ethAddress.Bytes(), tf.EthClient.PrivateKey, client1PrivKey)
	signedSimplexStateBytes, err := protobuf.Marshal(signedSimplexStateProto)
	tf.ChkErr(err, "failed to get signedSimplexStateBytes")
	msgRequestGuard := subscribe.NewMsgRequestGuard(ethAddress.String(), signedSimplexStateBytes, transactor.Key.GetAddress())
	transactor.BroadcastTx(msgRequestGuard)
	sleepWithLog(10, "sgn syncing Subscribe balance from mainchain")

	// Query sgn to check if request has correct stateproof data
	log.Info("Query sgn to check if request has correct stateproof data...")
	request, err := subscribe.CLIQueryRequest(transactor.CliCtx.Codec, transactor.CliCtx, subscribe.RouterKey, channelId[:])
	tf.ChkErr(err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	// expectedRes = fmt.Sprintf(`EthAddress: %s, DelegatedStake: %d`, ethAddress.String(), amt) // defined in Delegator.String()
	// assert.Equal(t, delegator.String(), expectedRes, fmt.Sprintf("The expected result should be \"%s\"", expectedRes))

	// Call intendSettle on ledger contract
	log.Info("Call intendSettle on ledger contract...")
	signedSimplexStateProto = prepareSignedSimplexState(1, channelId[:], ethAddress.Bytes(), tf.EthClient.PrivateKey, client1PrivKey)
	signedSimplexStateArrayBytes, err := protobuf.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{signedSimplexStateProto},
	})
	tf.ChkErr(err, "failed to get signedSimplexStateArrayBytes")
	tx, err = ledgerContract.IntendSettle(auth, signedSimplexStateArrayBytes)
	tf.ChkErr(err, "failed to IntendSettle")
	tf.WaitMinedWithChk(ctx, conn, tx, maxBlockDiff+2, "IntendSettle")

	// Query sgn to check if validator has submitted the stateproof correctly
	log.Info("Query sgn to check if validator has submitted the stateproof correctly...")
	sleepWithLog(20, "sgn submitting stateproof")
	request, err = subscribe.CLIQueryRequest(transactor.CliCtx.Codec, transactor.CliCtx, subscribe.RouterKey, channelId[:])
	tf.ChkErr(err, "failed to query request on sgn")
	log.Infoln("Query sgn about the request info:", request.String())
	// expectedRes = fmt.Sprintf(`EthAddress: %s, DelegatedStake: %d`, ethAddress.String(), amt) // defined in Delegator.String()
	// assert.Equal(t, delegator.String(), expectedRes, fmt.Sprintf("The expected result should be \"%s\"", expectedRes))
}

func prepareChannelInitializer() *entity.PaymentChannelInitializer {
	tokenInfo := &entity.TokenInfo{
		TokenType:    entity.TokenType_ERC20,
		TokenAddress: ctype.Hex2Bytes(MockCelerAddr),
	}

	lowAddrDist := &entity.AccountAmtPair{
		Account: tf.EthClient.Address[:],
		Amt:     big.NewInt(0).Bytes(),
	}

	highAddrDist := &entity.AccountAmtPair{
		Account: ctype.Hex2Bytes(client1AddrStr),
		Amt:     big.NewInt(0).Bytes(),
	}

	initializer := &entity.PaymentChannelInitializer{
		InitDistribution: &entity.TokenDistribution{
			Token: tokenInfo,
			Distribution: []*entity.AccountAmtPair{
				lowAddrDist, highAddrDist,
			},
		},
		OpenDeadline:   1000000,
		DisputeTimeout: 100,
	}
	return initializer
}

func prepareSignedSimplexState(seqNum uint64, channelId, peerFrom []byte, prvtKey0, prvtKey1 *ecdsa.PrivateKey) *chain.SignedSimplexState {
	simplexPaymentChannelBytes, err := protobuf.Marshal(&entity.SimplexPaymentChannel{
		SeqNum:    seqNum,
		ChannelId: channelId,
		PeerFrom:  peerFrom,
	})
	tf.ChkErr(err, "failed to get simplexPaymentChannelBytes")
	sig0, err := mainchain.SignMessage(prvtKey0, simplexPaymentChannelBytes)
	tf.ChkErr(err, "failed to get sig0")
	sig1, err := mainchain.SignMessage(prvtKey1, simplexPaymentChannelBytes)
	tf.ChkErr(err, "failed to get sig1")
	signedSimplexStateProto := &chain.SignedSimplexState{
		SimplexState: simplexPaymentChannelBytes,
		Sigs:         [][]byte{sig0, sig1},
	}

	return signedSimplexStateProto
}

func monitorOpenChannel(ledgerContract *mainchain.CelerLedger, channelIdChan chan [32]byte) {
	openChannelChan := make(chan *mainchain.CelerLedgerOpenChannel)
	sub, err := ledgerContract.WatchOpenChannel(nil, openChannelChan, nil, nil)
	if err != nil {
		log.Info("WatchInitializeCandidate err: ", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Info("WatchInitializeCandidate err: ", err)
		case openChannel := <-openChannelChan:
			log.Info("Monitored a OpenChannel event")
			channelId := [32]byte{}
			copy(channelId[:], openChannel.ChannelId[:])
			channelIdChan <- channelId
			return
		}
	}
}
