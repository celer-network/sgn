package testcommon

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

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
	protobuf "github.com/golang/protobuf/proto"
)

func ChkTestErr(t *testing.T, err error, msg string) {
	if err != nil {
		log.Errorln(msg, err)
		t.FailNow()
	}
}

func ChkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

// if status isn't 1 (success), log.Fatal
func ChkTxStatus(s uint64, txname string) {
	if s != 1 {
		log.Fatalln(txname, "tx failed")
	}
	log.Infoln(txname, "tx success")
}

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

func WaitMinedWithChk(ctx context.Context, conn *ethclient.Client,
	tx *ethtypes.Transaction, BlockDelay uint64, txname string) {
	ctx2, cancel := context.WithTimeout(ctx, waitMinedTimeout)
	defer cancel()
	receipt, err := mainchain.WaitMined(ctx2, conn, tx, BlockDelay)
	ChkErr(err, "WaitMined error")
	ChkTxStatus(receipt.Status, txname)
}

func LogBlkNum(conn *ethclient.Client) {
	header, err := conn.HeaderByNumber(context.Background(), nil)
	ChkErr(err, "failed to get HeaderByNumber")
	log.Infoln("Latest block number on mainchain:", header.Number)
}

func GetAddressFromKeystore(ksBytes []byte) (string, error) {
	type ksStruct struct {
		Address string
	}
	var ks ksStruct
	if err := json.Unmarshal(ksBytes, &ks); err != nil {
		return "", err
	}
	return ks.Address, nil
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

func PrepareSignedSimplexState(seqNum uint64, channelId, peerFrom []byte, prvtKey0, prvtKey1 *ecdsa.PrivateKey) (*chain.SignedSimplexState, error) {
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
