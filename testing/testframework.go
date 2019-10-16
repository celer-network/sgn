// Copyright 2018 Celer Network

package testing

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"sync"

	ccommon "github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/utils"
	"github.com/celer-network/sgn/testing/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	ethInstance = "http://127.0.0.1:8545"
)

var (
	pendingNonceLock sync.Mutex
	// e2eProfile after deploy contracts, with hardcoded ethgw etc
	// serialized json will be saved as outRootDir/profile.json
	// tests wish to use a different profile can overrides fields like svrRpc
	// and keep contract addresses etc
	E2eProfile     *ccommon.CProfile
	GuardAddr      string
	Erc20TokenAddr string
)

var (
	// to be set by test/e2e
	outRootDir  string
	envDir      = "../../testing/env"
	etherBaseKs = envDir + "/keystore/etherbase.json"
)

func SetEnvDir(dir string) {
	envDir = dir
	etherBaseKs = envDir + "/keystore/etherbase.json"
}

func SetOutRootDir(dir string) {
	outRootDir = dir
}

func StartProcess(name string, args ...string) *os.Process {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	return cmd.Process
}

func KillProcess(process *os.Process) {
	process.Kill()
	process.Release()
}

func prepareEthClient() (
	*ethclient.Client, *bind.TransactOpts, context.Context, common.Address, error) {
	conn, err := ethclient.Dial(ethInstance)
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
		receipt, err := utils.WaitMined(ctx, conn, tx, 0)
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
