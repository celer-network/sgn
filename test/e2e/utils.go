// Copyright 2018 Celer Network

package e2e

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"context"
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tf "github.com/celer-network/sgn/testing"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	protobuf "github.com/golang/protobuf/proto"
)

func sleep(second time.Duration) {
	time.Sleep(second * time.Second)
}

func sleepWithLog(second time.Duration, waitFor string) {
	log.Infof("Sleep %d seconds for %s", second, waitFor)
	sleep(second)
}

func sleepBlocksWithLog(count time.Duration, waitFor string) {
	sleepWithLog(count * sgnBlockInterval, waitFor)
}

func parseGatewayQueryResponse(resp *http.Response, cdc *codec.Codec) json.RawMessage {
	body, err := ioutil.ReadAll(resp.Body)
	tf.ChkErr(err, "failed to read http response")

	var responseWithHeight rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &responseWithHeight)
	return responseWithHeight.Result
}

func initializeCandidate(auth *bind.TransactOpts, sgnAddr sdk.AccAddress) error {
	conn := tf.EthClient.Client
	guardContract := tf.EthClient.Guard
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	log.Info("Call initializeCandidate on guard contract using the validator eth address...")
	tx, err := guardContract.InitializeCandidate(auth, big.NewInt(1), sgnAddr.Bytes())
	if err != nil {
		return err
	}

	tf.WaitMinedWithChk(ctx, conn, tx, blockDelay, "InitializeCandidate")
	sleepBlocksWithLog(5, "sgn syncing InitializeCandidate event on mainchain")
	return nil
}

func delegateStake(fromAuth *bind.TransactOpts, toEthAddress mainchain.Addr, amt *big.Int) error {
	conn := tf.EthClient.Client
	guardContract := tf.EthClient.Guard
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	log.Info("Call delegate on guard contract to delegate stake to the validator eth address...")
	tx, err := celrContract.Approve(fromAuth, guardAddr, amt)
	if err != nil {
		return err
	}
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")
	tx, err = guardContract.Delegate(fromAuth, toEthAddress, amt)
	if err != nil {
		return err
	}

<<<<<<< HEAD
	tf.WaitMinedWithChk(ctx, conn, tx, 3 * blockDelay, "Delegate to validator")
=======
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Delegate to validator")
	sleepWithLog(10, "sgn syncing Delegate event on mainchain")
>>>>>>> origin
	return nil
}

func openChannel(peer0Addr, peer1Addr []byte, peer0PrivKey, peer1PrivKey *ecdsa.PrivateKey) (channelId [32]byte, err error) {
	log.Info("Call openChannel on ledger contract...")
	conn := tf.EthClient.Client
	auth := tf.EthClient.Auth
	ledgerContract := tf.EthClient.Ledger
	tokenInfo := &entity.TokenInfo{
		TokenType:    entity.TokenType_ERC20,
		TokenAddress: mockCelerAddr.Bytes(),
	}
	lowAddrDist := &entity.AccountAmtPair{
		Account: peer0Addr,
		Amt:     big.NewInt(0).Bytes(),
	}
	highAddrDist := &entity.AccountAmtPair{
		Account: peer1Addr,
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
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	paymentChannelInitializerBytes, err := protobuf.Marshal(initializer)
	if err != nil {
		return
	}

	sig0, err := mainchain.SignMessage(peer0PrivKey, paymentChannelInitializerBytes)
	if err != nil {
		return
	}

	sig1, err := mainchain.SignMessage(peer1PrivKey, paymentChannelInitializerBytes)
	if err != nil {
		return
	}

	requestBytes, err := protobuf.Marshal(&chain.OpenChannelRequest{
		ChannelInitializer: paymentChannelInitializerBytes,
		Sigs:               [][]byte{sig0, sig1},
	})
	if err != nil {
		return
	}

	channelIdChan := make(chan [32]byte)
	go monitorOpenChannel(ledgerContract, channelIdChan)
	tx, err := ledgerContract.OpenChannel(auth, requestBytes)
	if err != nil {
		return
	}

	tf.WaitMinedWithChk(ctx, conn, tx, blockDelay+2, "OpenChannel")
	channelId = <-channelIdChan
	log.Info("channel ID: ", mainchain.Bytes2Hex(channelId[:]))

	return
}
