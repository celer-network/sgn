// Copyright 2018 Celer Network

package testing

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/flags"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/testing/log"
	"github.com/celer-network/sgn/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

const (
	EthInstance = "http://127.0.0.1:8545"
)

var (
	pendingNonceLock sync.Mutex
	etherBaseKs      string
	EthClient        *mainchain.EthClient
)

func GetLatestBlkNum(conn *ethclient.Client) (*big.Int, error) {
	header, err := conn.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

func SetEnvDir(envDir string) {
	etherBaseKs = envDir + "/keystore/etherbase.json"
}

func SetupEthClient() {
	ec, err := mainchain.NewEthClient(
		viper.GetString(flags.FlagEthWS),
		viper.GetString(flags.FlagEthGuardAddress),
		viper.GetString(flags.FlagEthLedgerAddress),
		// viper.GetString(flags.FlagEthKeystore),
		"../../test/keys/client0.json", // relative path is different in tests
		viper.GetString(flags.FlagEthPassphrase),
	)
	if err != nil {
		log.Fatal(err)
	}

	EthClient = ec
}

func prepareEthClient() (
	*ethclient.Client, *bind.TransactOpts, context.Context, common.Address, error) {
	conn, err := ethclient.Dial(EthInstance)
	if err != nil {
		return nil, nil, nil, common.Address{}, err
	}
	fmt.Println("etherBaseKs", etherBaseKs)
	etherBaseKsBytes, err := ioutil.ReadFile(etherBaseKs)
	if err != nil {
		return nil, nil, nil, common.Address{}, err
	}
	etherBaseAddrStr, err := utils.GetAddressFromKeystore(etherBaseKsBytes)
	if err != nil {
		return nil, nil, nil, common.Address{}, err
	}
	etherBaseAddr := ctype.Hex2Addr(etherBaseAddrStr)
	auth, err := bind.NewTransactor(strings.NewReader(string(etherBaseKsBytes)), "")
	if err != nil {
		return nil, nil, nil, common.Address{}, err
	}
	return conn, auth, context.Background(), etherBaseAddr, nil
}

func fundAccount(amount string, recipients []*common.Address) error {
	conn, auth, ctx, senderAddr, err := prepareEthClient()
	if err != nil {
		return err
	}
	value := big.NewInt(0)
	value.SetString(amount, 10)
	auth.Value = value
	chainID := big.NewInt(883) // Private Mainchain Testnet
	var gasLimit uint64 = 21000
	for _, r := range recipients {
		pendingNonceLock.Lock()
		nonce, err := conn.PendingNonceAt(ctx, senderAddr)
		if err != nil {
			pendingNonceLock.Unlock()
			return err
		}
		gasPrice, err := conn.SuggestGasPrice(ctx)
		if err != nil {
			pendingNonceLock.Unlock()
			return err
		}
		tx := types.NewTransaction(nonce, *r, auth.Value, gasLimit, gasPrice, nil)
		tx, err = auth.Signer(types.NewEIP155Signer(chainID), senderAddr, tx)
		if err != nil {
			pendingNonceLock.Unlock()
			return err
		}
		if *r == ctype.ZeroAddr {
			log.Info("Advancing block")
		} else {
			log.Infof(
				"Sending %s wei from %s to %s, nonce %d. tx: %x", amount, senderAddr.Hex(), r.Hex(), nonce, tx.Hash())
		}

		err = conn.SendTransaction(ctx, tx)
		if err != nil {
			pendingNonceLock.Unlock()
			return err
		}
		pendingNonceLock.Unlock()
		receipt, err := mainchain.WaitMined(ctx, conn, tx, 0)
		if err != nil {
			log.Error(err)
		}
		if receipt.Status != 1 {
			log.Errorln("tx failed. tx hash:", receipt.TxHash.Hex())
		} else {
			if *r == ctype.ZeroAddr {
				head, _ := conn.HeaderByNumber(ctx, nil)
				log.Info("Current block number:", head.Number.String())
			} else {
				bal, _ := conn.BalanceAt(ctx, *r, nil)
				log.Infoln("tx done.", r.String(), "bal:", bal.String())
			}
		}
	}
	return nil
}

func FundAddr(amt string, recipients []*common.Address) error {
	return fundAccount(amt, recipients)
}

func AdvanceBlock() error {
	return fundAccount("0", []*common.Address{&common.Address{}})
}

func AdvanceBlocks(blockCount uint64) error {
	var i uint64
	for i = 0; i < blockCount; i++ {
		AdvanceBlock()
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func AdvanceBlocksUntilDone(done chan bool) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case <-ticker.C:
			AdvanceBlock()
		}
	}
}
