package testing

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

// if status isn't 1 (sucess), log.Fatal
func ChkTxStatus(s uint64, txname string) {
	if s != 1 {
		log.Fatalln(txname, "tx failed")
	}
	log.Infoln(txname, "tx success")
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

func InitializeCandidate(auth *bind.TransactOpts, sgnAddr sdk.AccAddress) error {
	conn := EthClient.Client
	guardContract := EthClient.Guard
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	log.Info("Call initializeCandidate on guard contract using the validator eth address...")
	tx, err := guardContract.InitializeCandidate(auth, big.NewInt(1), sgnAddr.Bytes())
	if err != nil {
		return err
	}

	WaitMinedWithChk(ctx, conn, tx, BlockDelay, "InitializeCandidate")
	SleepBlocksWithLog(6, "sgn syncing InitializeCandidate event on mainchain")
	return nil
}

func DelegateStake(celrContract *mainchain.ERC20, guardAddr mainchain.Addr, fromAuth *bind.TransactOpts, toEthAddress mainchain.Addr, amt *big.Int) error {
	conn := EthClient.Client
	guardContract := EthClient.Guard
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	log.Info("Call delegate on guard contract to delegate stake to the validator eth address...")
	tx, err := celrContract.Approve(fromAuth, guardAddr, amt)
	if err != nil {
		return err
	}
	WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")
	tx, err = guardContract.Delegate(fromAuth, toEthAddress, amt)
	if err != nil {
		return err
	}

	WaitMinedWithChk(ctx, conn, tx, 3*BlockDelay, "Delegate to validator")
	SleepWithLog(10, "sgn syncing Delegate event on mainchain")
	return nil
}
