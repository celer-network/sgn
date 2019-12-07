package testing

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
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
	tx *ethtypes.Transaction, blockDelay uint64, txname string) {
	receipt, err := mainchain.WaitMined(ctx, conn, tx, blockDelay)
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
