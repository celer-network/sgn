// Copyright 2018 Celer Network

package testcommon

import (
	"context"
	"crypto/ecdsa"
	"io/ioutil"
	"math"
	"math/big"
	"strings"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	protobuf "github.com/golang/protobuf/proto"
)

var (
	etherBaseKs = EnvDir + "/keystore/etherbase.json"

	EtherBase = &mainchain.EthClient{}
	Client0   = &mainchain.EthClient{}
	Client1   = &mainchain.EthClient{}

	DefaultTestEthClient = &mainchain.EthClient{}
)

func SetEthBaseKs(prefix string) {
	etherBaseKs = prefix + "/keystore/etherbase.json"
}

// SetupEthClients sets Client part (Client) and Auth part (PrivateKey, Address, Auth)
// Contracts part (GuardAddress, Guard, LedgerAddress, Ledger) is set after deploying Guard contracts in setupNewSGNEnv()
func SetupEthClients() {
	EtherBase = setupEthClient(etherBaseKs)
	Client0 = setupEthClient(ClientEthKs[0])
	Client1 = setupEthClient(ClientEthKs[1])
	DefaultTestEthClient = setupEthClient("../../keys/ethks0.json")
}

func setupEthClient(ksfile string) *mainchain.EthClient {
	ethClient := &mainchain.EthClient{}
	err := ethClient.SetClient(LocalGeth)
	ChkErr(err, "failed to connect to the Ethereum")
	err = ethClient.SetAuth(ksfile, "")
	ChkErr(err, "failed to create auth")
	return ethClient
}

func SetContracts(guardAddr, ledgerAddr mainchain.Addr) error {
	log.Infof("set contracts guard %x ledger %x", guardAddr, ledgerAddr)
	err := EtherBase.SetContracts(guardAddr.String(), ledgerAddr.String())
	if err != nil {
		return err
	}
	err = Client0.SetContracts(guardAddr.String(), ledgerAddr.String())
	if err != nil {
		return err
	}
	err = Client1.SetContracts(guardAddr.String(), ledgerAddr.String())
	if err != nil {
		return err
	}
	err = DefaultTestEthClient.SetContracts(guardAddr.String(), ledgerAddr.String())
	if err != nil {
		return err
	}
	return nil
}

func SetupE2eProfile() {
	ledgerAddr := DeployLedgerContract()
	// Deploy sample ERC20 contract (CELR)
	erc20Addr, erc20 := DeployERC20Contract()

	E2eProfile = &TestProfile{
		// hardcoded values
		DisputeTimeout: 10,
		// deployed addresses
		LedgerAddr:   ledgerAddr,
		CelrAddr:     erc20Addr,
		CelrContract: erc20,
	}
}

func FundAddrsETH(amt string, recipients []mainchain.Addr) error {
	conn, auth, ctx, senderAddr, err := prepareEtherBaseClient()
	if err != nil {
		return err
	}
	value := big.NewInt(0)
	value.SetString(amt, 10)
	auth.Value = value
	chainID := big.NewInt(883) // Private Mainchain Testnet
	var gasLimit uint64 = 21000
	for _, addr := range recipients {
		nonce, err := conn.PendingNonceAt(ctx, senderAddr)
		if err != nil {
			return err
		}
		gasPrice, err := conn.SuggestGasPrice(ctx)
		if err != nil {
			return err
		}
		tx := types.NewTransaction(nonce, addr, auth.Value, gasLimit, gasPrice, nil)
		tx, err = auth.Signer(types.NewEIP155Signer(chainID), senderAddr, tx)
		if err != nil {
			return err
		}
		if addr == mainchain.ZeroAddr {
			log.Info("Advancing block")
		} else {
			log.Infof("Sending ETH %s to %x from %x", amt, addr, senderAddr)
		}

		err = conn.SendTransaction(ctx, tx)
		if err != nil {
			return err
		}
		ctx2, cancel := context.WithTimeout(ctx, waitMinedTimeout)
		defer cancel()
		receipt, err := mainchain.WaitMined(ctx2, conn, tx, 0)
		if err != nil {
			log.Error(err)
		}
		if receipt.Status != 1 {
			log.Errorf("tx failed. tx hash: %x", receipt.TxHash)
		} else {
			if addr == mainchain.ZeroAddr {
				head, _ := conn.HeaderByNumber(ctx, nil)
				log.Infoln("Current block number:", head.Number.String())
			} else {
				bal, _ := conn.BalanceAt(ctx, addr, nil)
				log.Infoln("Tx done.", addr.String(), "bal:", bal.String())
			}
		}
	}
	return nil
}

func FundAddrsErc20(erc20Addr mainchain.Addr, addrs []mainchain.Addr, amount string) error {
	erc20Contract, err := mainchain.NewERC20(erc20Addr, EtherBase.Client)
	if err != nil {
		return err
	}
	tokenAmt := new(big.Int)
	tokenAmt.SetString(amount, 10)
	for _, addr := range addrs {
		tx, err := erc20Contract.Transfer(EtherBase.Auth, addr, tokenAmt)
		if err != nil {
			return err
		}
		log.Infof("Sending ERC20 %s to %x from %x", amount, addr, EtherBase.Auth.From)
		_, err = mainchain.WaitMined(context.Background(), EtherBase.Client, tx, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func OpenChannel(peer0Addr, peer1Addr mainchain.Addr, peer0PrivKey, peer1PrivKey *ecdsa.PrivateKey) (channelId [32]byte, err error) {
	log.Infoln("Call openChannel on ledger contract", mainchain.Addr2Hex(peer0Addr), mainchain.Addr2Hex(peer1Addr))
	tokenInfo := &entity.TokenInfo{
		TokenType: entity.TokenType_ETH,
	}
	lowAddrDist := &entity.AccountAmtPair{
		Account: peer1Addr.Bytes(),
		Amt:     big.NewInt(0).Bytes(),
	}
	highAddrDist := &entity.AccountAmtPair{
		Account: peer0Addr.Bytes(),
		Amt:     big.NewInt(0).Bytes(),
	}
	initializer := &entity.PaymentChannelInitializer{
		InitDistribution: &entity.TokenDistribution{
			Token: tokenInfo,
			Distribution: []*entity.AccountAmtPair{
				lowAddrDist, highAddrDist,
			},
		},
		OpenDeadline:   math.MaxUint64,
		DisputeTimeout: DisputeTimeout,
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
		Sigs:               [][]byte{sig1, sig0},
	})
	if err != nil {
		return
	}

	channelIdChan := make(chan [32]byte)
	go monitorOpenChannel(channelIdChan)
	_, err = Client0.Ledger.OpenChannel(Client0.Auth, requestBytes)
	if err != nil {
		return
	}

	channelId = <-channelIdChan
	log.Info("channel ID: ", mainchain.Bytes2Hex(channelId[:]))

	return
}

func IntendWithdraw(auth *bind.TransactOpts, candidateAddr mainchain.Addr, amt *big.Int) error {
	conn := EtherBase.Client
	guardContract := EtherBase.Guard
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	log.Info("Call intendWithdraw on guard contract using the validator eth address...")
	tx, err := guardContract.IntendWithdraw(auth, candidateAddr, amt)
	if err != nil {
		return err
	}

	WaitMinedWithChk(ctx, conn, tx, BlockDelay, "IntendWithdraw")
	return nil
}

func InitializeCandidate(auth *bind.TransactOpts, sgnAddr sdk.AccAddress, minSelfStake *big.Int) error {
	conn := EtherBase.Client
	guardContract := EtherBase.Guard
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	log.Info("Call initializeCandidate on guard contract using the validator eth address...")
	tx, err := guardContract.InitializeCandidate(auth, minSelfStake, sgnAddr.Bytes())
	if err != nil {
		return err
	}

	WaitMinedWithChk(ctx, conn, tx, BlockDelay, "InitializeCandidate")
	return nil
}

func DelegateStake(celrContract *mainchain.ERC20, guardAddr mainchain.Addr, fromAuth *bind.TransactOpts, toEthAddress mainchain.Addr, amt *big.Int) error {
	conn := EtherBase.Client
	guardContract := EtherBase.Guard
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
	WaitMinedWithChk(ctx, conn, tx, BlockDelay, "Delegate to validator")
	return nil
}

func monitorOpenChannel(channelIdChan chan [32]byte) {
	openChannelChan := make(chan *mainchain.CelerLedgerOpenChannel)
	sub, err := Client0.Ledger.WatchOpenChannel(nil, openChannelChan, nil, nil)
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

// Remove this
func prepareEtherBaseClient() (
	*ethclient.Client, *bind.TransactOpts, context.Context, mainchain.Addr, error) {
	conn, err := ethclient.Dial(LocalGeth)
	if err != nil {
		return nil, nil, nil, mainchain.Addr{}, err
	}
	log.Infoln("EtherBaseKs: ", etherBaseKs)
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
