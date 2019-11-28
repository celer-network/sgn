// Copyright 2018 Celer Network

package testing

import (
	"context"
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"strings"
	"sync"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	protobuf "github.com/golang/protobuf/proto"
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
		viper.GetString(common.FlagEthWS),
		viper.GetString(common.FlagEthGuardAddress),
		viper.GetString(common.FlagEthLedgerAddress),
		viper.GetString(common.FlagEthKeystore),
		viper.GetString(common.FlagEthPassphrase),
	)
	ChkErr(err, "setup eth client")
	EthClient = ec
}

func prepareEthClient() (
	*ethclient.Client, *bind.TransactOpts, context.Context, mainchain.Addr, error) {
	conn, err := ethclient.Dial(EthInstance)
	if err != nil {
		return nil, nil, nil, mainchain.Addr{}, err
	}
	log.Infoln("etherBaseKs", etherBaseKs)
	etherBaseKsBytes, err := ioutil.ReadFile(etherBaseKs)
	if err != nil {
		return nil, nil, nil, mainchain.Addr{}, err
	}
	etherBaseAddrStr, err := GetAddressFromKeystore(etherBaseKsBytes)
	if err != nil {
		return nil, nil, nil, mainchain.Addr{}, err
	}
	etherBaseAddr := mainchain.Hex2Addr(etherBaseAddrStr)
	auth, err := bind.NewTransactor(strings.NewReader(string(etherBaseKsBytes)), "")
	if err != nil {
		return nil, nil, nil, mainchain.Addr{}, err
	}
	return conn, auth, context.Background(), etherBaseAddr, nil
}

func FundAddr(amt string, recipients []*mainchain.Addr) error {
	conn, auth, ctx, senderAddr, err := prepareEthClient()
	if err != nil {
		return err
	}
	value := big.NewInt(0)
	value.SetString(amt, 10)
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
		if *r == mainchain.ZeroAddr {
			log.Info("Advancing block")
		} else {
			log.Infof("Sending %s wei from %x to %x, nonce %d. tx: %x", amt, senderAddr, r, nonce, tx.Hash())
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
			log.Errorf("tx failed. tx hash: %x", receipt.TxHash)
		} else {
			if *r == mainchain.ZeroAddr {
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

func OpenChannel(peer0Addr, peer1Addr []byte, peer0PrivKey, peer1PrivKey *ecdsa.PrivateKey, tokenAddr []byte) (channelId [32]byte, err error) {
	log.Info("Call openChannel on ledger contract...")
	auth := EthClient.Auth
	ledgerContract := EthClient.Ledger
	tokenInfo := &entity.TokenInfo{
		TokenType:    entity.TokenType_ERC20,
		TokenAddress: tokenAddr,
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
	go monitorOpenChannel(channelIdChan)
	tx, err := ledgerContract.OpenChannel(auth, requestBytes)
	if err != nil {
		return
	}

	channelId = <-channelIdChan
	log.Info("channel ID: ", mainchain.Bytes2Hex(channelId[:]))

	return
}

func monitorOpenChannel(channelIdChan chan [32]byte) {
	openChannelChan := make(chan *mainchain.CelerLedgerOpenChannel)
	sub, err := EthClient.Ledger.WatchOpenChannel(nil, openChannelChan, nil, nil)
	if err != nil {
		log.Errorln("WatchInitializeCandidate err: ", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Errorln("WatchInitializeCandidate err: ", err)
		case openChannel := <-openChannelChan:
			log.Infoln("Monitored a OpenChannel event")
			channelId := [32]byte{}
			copy(channelId[:], openChannel.ChannelId[:])
			channelIdChan <- channelId
			return
		}
	}
}
