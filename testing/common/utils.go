package common

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/protobuf/proto"
)

func GetAuth(ksfile string) (addr mainchain.Addr, auth *bind.TransactOpts, err error) {
	keystoreBytes, err := ioutil.ReadFile(ksfile)
	if err != nil {
		return
	}
	key, err := keystore.DecryptKey(keystoreBytes, "")
	if err != nil {
		return
	}
	addr = key.Address
	auth, err = bind.NewTransactor(strings.NewReader(string(keystoreBytes)), "")
	if err != nil {
		return
	}
	return
}

func GetEthPrivateKey(ksfile string) (*ecdsa.PrivateKey, error) {
	keystoreBytes, err := ioutil.ReadFile(ksfile)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keystoreBytes, "")
	if err != nil {
		return nil, err
	}
	return key.PrivateKey, nil
}

func WaitMinedWithChk(
	ctx context.Context,
	conn *ethclient.Client,
	tx *ethtypes.Transaction,
	blockDelay uint64,
	pollingInterval time.Duration,
	txname string,
) {
	ctx2, cancel := context.WithTimeout(ctx, waitMinedTimeout)
	defer cancel()
	receipt, err := eth.WaitMined(ctx2, conn, tx, eth.WithBlockDelay(blockDelay), eth.WithPollingInterval(pollingInterval))
	ChkErr(err, "WaitMined error")
	if receipt.Status != ethtypes.ReceiptStatusSuccessful {
		log.Fatalln(txname, "tx failed")
	}
	log.Infoln(txname, "tx success")
}

func LogBlkNum(conn *ethclient.Client) {
	header, err := conn.HeaderByNumber(context.Background(), nil)
	ChkErr(err, "failed to get HeaderByNumber")
	log.Infoln("Latest block number on mainchain:", header.Number)
}

func sleep(second time.Duration) {
	time.Sleep(second * time.Second)
}

func SleepWithLog(second time.Duration, waitFor string) {
	log.Infof("Sleep %d seconds for %s", second, waitFor)
	sleep(second)
}

func SleepBlocksWithLog(count time.Duration, waitFor string) {
	SleepWithLog(count*SgnBlockInterval, waitFor)
}

func ParseGatewayQueryResponse(resp *http.Response, cdc *codec.Codec) (json.RawMessage, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseWithHeight rest.ResponseWithHeight
	err = cdc.UnmarshalJSON(body, &responseWithHeight)
	if err != nil {
		return nil, err
	}

	return responseWithHeight.Result, nil
}

func PrepareSignedSimplexState(seqNum uint64, channelId, peerFrom []byte, peer0, peer1 *TestEthClient) (*chain.SignedSimplexState, error) {
	simplexPaymentChannelBytes, err := proto.Marshal(&entity.SimplexPaymentChannel{
		SeqNum:    seqNum,
		ChannelId: channelId,
		PeerFrom:  peerFrom,
	})
	if err != nil {
		return nil, err
	}

	lo, hi := peer0, peer1
	if bytes.Compare(peer0.Address.Bytes(), peer1.Address.Bytes()) > 0 {
		lo, hi = peer1, peer0
	}

	siglo, err := lo.Signer.SignEthMessage(simplexPaymentChannelBytes)
	if err != nil {
		return nil, err
	}

	sighi, err := hi.Signer.SignEthMessage(simplexPaymentChannelBytes)
	if err != nil {
		return nil, err
	}

	signedSimplexStateProto := &chain.SignedSimplexState{
		SimplexState: simplexPaymentChannelBytes,
		Sigs:         [][]byte{siglo, sighi},
	}

	return signedSimplexStateProto, nil
}
