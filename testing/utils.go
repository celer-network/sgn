package testing

import (
	"context"
	"encoding/json"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ChkErr(e error, msg string) {
	if e != nil {
		log.Fatalln(msg, e)
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
	blkNum, err := GetLatestBlkNum(conn)
	ChkErr(err, "failed to get HeaderByNumber")
	log.Infoln("Latest block number on mainchain:", blkNum)
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
