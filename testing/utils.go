package testing

import (
	"context"

	"github.com/celer-network/sgn/testing/log"
	"github.com/celer-network/sgn/utils"
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
		log.Fatal(txname + " tx failed")
	}
	log.Info(txname + " tx success")
}

func WaitMinedWithChk(ctx context.Context, conn *ethclient.Client,
	tx *ethtypes.Transaction, blockDelay uint64, txname string) {
	receipt, err := utils.WaitMined(ctx, conn, tx, blockDelay)
	ChkErr(err, "WaitMined error")
	ChkTxStatus(receipt.Status, txname)
}

func LogBlkNum(conn *ethclient.Client) {
	blkNum, err := GetLatestBlkNum(conn)
	ChkErr(err, "failed to get HeaderByNumber")
	log.Infoln("Latest block number on mainchain: ", blkNum)
}
