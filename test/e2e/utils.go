// Copyright 2018 Celer Network

package e2e

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"context"
	"math/big"

	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/crypto"
	protobuf "github.com/golang/protobuf/proto"
)

func sleep(second time.Duration) {
	time.Sleep(second * time.Second)
}

func sleepWithLog(second time.Duration, waitFor string) {
	log.Infof("Sleep %d seconds for %s", second, waitFor)
	sleep(second)
}

func parseGatewayQueryResponse(resp *http.Response, cdc *codec.Codec) json.RawMessage {
	body, err := ioutil.ReadAll(resp.Body)
	tf.ChkErr(err, "failed to read http response")

	var responseWithHeight rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &responseWithHeight)
	return responseWithHeight.Result
}

func buildContextWithTimeout() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 30 * time.Second)
	return ctx
}

func initializeCandidate() error {
	ctx := buildContextWithTimeout()
	conn := tf.EthClient.Client
	auth := tf.EthClient.Auth
	guardContract := tf.EthClient.Guard
	sgnAddr, err := sdk.AccAddressFromBech32(client0SGNAddrStr)
	if err != nil {
		return err
	}

	log.Info("Call initializeCandidate on guard contract using the validator eth address...")
	tx, err := guardContract.InitializeCandidate(auth, big.NewInt(1), sgnAddr.Bytes())
	if err != nil {
		return err
	}

	tf.WaitMinedWithChk(ctx, conn, tx, 0, "InitializeCandidate")
	sleepWithLog(30, "sgn syncing InitializeCandidate event on mainchain")
	return nil
}

func delegateStake() error {
	ctx := buildContextWithTimeout()
	conn := tf.EthClient.Client
	auth := tf.EthClient.Auth
	ethAddress := tf.EthClient.Address
	guardContract := tf.EthClient.Guard

	log.Info("Call delegate on guard contract to delegate stake to the validator eth address...")
	amt := big.NewInt(100)
	tx, err := celrContract.Approve(auth, guardAddr, amt)
	if err != nil {
		return err
	}
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")
	tx, err = guardContract.Delegate(auth, ethAddress, amt)
	if err != nil {
		return err
	}

	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Delegate to validator")
	sleepWithLog(30, "sgn syncing Delegate event on mainchain")
	return nil
}

func openChannel() (channelId [32]byte, err error) {
	log.Info("Call openChannel on ledger contract...")
	client1PrivKey, _ := crypto.HexToECDSA(client1Priv)
	ctx := buildContextWithTimeout()
	conn := tf.EthClient.Client
	auth := tf.EthClient.Auth
	ledgerContract := tf.EthClient.Ledger
	tokenInfo := &entity.TokenInfo{
		TokenType:    entity.TokenType_ERC20,
		TokenAddress: mockCelerAddr.Bytes(),
	}
	lowAddrDist := &entity.AccountAmtPair{
		Account: tf.EthClient.Address.Bytes(),
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
	paymentChannelInitializerBytes, err := protobuf.Marshal(initializer)
	if err != nil {
		return
	}

	sig0, err := mainchain.SignMessage(tf.EthClient.PrivateKey, paymentChannelInitializerBytes)
	if err != nil {
		return
	}

	sig1, err := mainchain.SignMessage(client1PrivKey, paymentChannelInitializerBytes)
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

	tf.WaitMinedWithChk(ctx, conn, tx, maxBlockDiff+2, "OpenChannel")
	channelId = <-channelIdChan
	log.Info("channel ID: ", ctype.Bytes2Hex(channelId[:]))

	return
}
