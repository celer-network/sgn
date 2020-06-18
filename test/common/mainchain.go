// Copyright 2018 Celer Network

package testcommon

import (
	"bytes"
	"context"
	"io/ioutil"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/golang/protobuf/proto"
)

type TestEthClient struct {
	Address mainchain.Addr
	Auth    *bind.TransactOpts
	Signer  eth.Signer
}

var (
	etherBaseKs = EnvDir + "/keystore/etherbase.json"

	EthClient      *ethclient.Client
	EtherBaseAuth  *bind.TransactOpts
	DposContract   *mainchain.DPoS
	SgnContract    *mainchain.SGN
	LedgerContract *mainchain.CelerLedger

	Client0 *TestEthClient
	Client1 *TestEthClient
)

func SetEthBaseKs(prefix string) {
	etherBaseKs = prefix + "/keystore/etherbase.json"
}

// SetupEthClients sets Client part (Client) and Auth part (PrivateKey, Address, Auth)
// Contracts part (DPoSAddress, DPoS, SGNAddress, SGN, LedgerAddress, Ledger) is set after deploying DPoS and SGN contracts in setupNewSGNEnv()
func SetupEthClients() {
	rpcClient, err := rpc.Dial(LocalGeth)
	if err != nil {
		log.Fatal(err)
	}
	EthClient = ethclient.NewClient(rpcClient)

	_, EtherBaseAuth, err = GetAuth(etherBaseKs)
	Client0, err = SetupTestEthClient(ClientEthKs[0])
	if err != nil {
		log.Fatal(err)
	}
	Client1, err = SetupTestEthClient(ClientEthKs[1])
	if err != nil {
		log.Fatal(err)
	}
}

func SetupTestEthClient(ksfile string) (*TestEthClient, error) {
	addr, auth, err := GetAuth(ksfile)
	if err != nil {
		return nil, err
	}
	testClient := &TestEthClient{
		Address: addr,
		Auth:    auth,
	}
	ksBytes, err := ioutil.ReadFile(ksfile)
	testClient.Signer, err = eth.NewSignerFromKeystore(string(ksBytes), "", nil)
	return testClient, nil
}

func SetContracts(dposAddr, sgnAddr, ledgerAddr mainchain.Addr) error {
	log.Infof("set contracts dpos %x sgn %x ledger %x", dposAddr, sgnAddr, ledgerAddr)
	var err error
	DposContract, err = mainchain.NewDPoS(dposAddr, EthClient)
	if err != nil {
		return err
	}
	SgnContract, err = mainchain.NewSGN(sgnAddr, EthClient)
	if err != nil {
		return err
	}
	LedgerContract, err = mainchain.NewCelerLedger(ledgerAddr, EthClient)
	if err != nil {
		return err
	}
	return nil
}

func SetupE2eProfile() {
	ledgerAddr := DeployLedgerContract()
	// Deploy sample ERC20 contract (CELR)
	tx, erc20Addr, erc20 := DeployERC20Contract()
	WaitMinedWithChk(context.Background(), EthClient, tx, BlockDelay, PollingInterval, "DeployERC20")

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
	conn, auth, ctx, senderAddr, connErr := prepareEtherBaseClient()
	if connErr != nil {
		return connErr
	}
	value := big.NewInt(0)
	value.SetString(amt, 10)
	auth.Value = value
	chainID := big.NewInt(883) // Private Mainchain Testnet
	var gasLimit uint64 = 21000
	var lastTx *types.Transaction
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
		lastTx = tx
	}
	ctx2, cancel := context.WithTimeout(ctx, waitMinedTimeout)
	defer cancel()
	receipt, err := eth.WaitMined(ctx2, conn, lastTx, eth.WithBlockDelay(BlockDelay), eth.WithPollingInterval(PollingInterval))
	if err != nil {
		log.Error(err)
	}
	if receipt.Status != 1 {
		log.Errorf("last tx failed. tx hash: %x", receipt.TxHash)
	} else {
		for _, addr := range recipients {
			if addr == mainchain.ZeroAddr {
				head, _ := conn.HeaderByNumber(ctx, nil)
				log.Infoln("Current block number:", head.Number.String())
			} else {
				bal, _ := conn.BalanceAt(ctx, addr, nil)
				log.Infoln("Funded.", addr.String(), "bal:", bal.String())
			}
		}
	}
	return nil
}

func FundAddrsErc20(erc20Addr mainchain.Addr, addrs []mainchain.Addr, amount string) error {
	erc20Contract, err := mainchain.NewERC20(erc20Addr, EthClient)
	if err != nil {
		return err
	}
	tokenAmt := new(big.Int)
	tokenAmt.SetString(amount, 10)
	var lastTx *types.Transaction
	for _, addr := range addrs {
		tx, transferErr := erc20Contract.Transfer(EtherBaseAuth, addr, tokenAmt)
		if transferErr != nil {
			return transferErr
		}
		lastTx = tx
		log.Infof("Sending ERC20 %s to %x from %x", amount, addr, EtherBaseAuth.From)
	}
	_, err = eth.WaitMined(context.Background(), EthClient, lastTx, eth.WithBlockDelay(BlockDelay), eth.WithPollingInterval(PollingInterval))
	return err
}

func OpenChannel(peer0, peer1 *TestEthClient) (channelId [32]byte, err error) {
	log.Infoln("Call openChannel on ledger contract", mainchain.Addr2Hex(peer0.Address), mainchain.Addr2Hex(peer1.Address))

	lo, hi := peer0, peer1
	if bytes.Compare(peer0.Address.Bytes(), peer1.Address.Bytes()) > 0 {
		lo, hi = peer1, peer0
	}

	tokenInfo := &entity.TokenInfo{
		TokenType: entity.TokenType_ETH,
	}
	loAddrDist := &entity.AccountAmtPair{
		Account: lo.Address.Bytes(),
		Amt:     big.NewInt(0).Bytes(),
	}
	hiAddrDist := &entity.AccountAmtPair{
		Account: hi.Address.Bytes(),
		Amt:     big.NewInt(0).Bytes(),
	}
	initializer := &entity.PaymentChannelInitializer{
		InitDistribution: &entity.TokenDistribution{
			Token: tokenInfo,
			Distribution: []*entity.AccountAmtPair{
				loAddrDist, hiAddrDist,
			},
		},
		// The unit of OpenDeadline is block number. time.Now() here is only used for uniqueness of each run.
		OpenDeadline:   uint64(time.Now().Unix()) + math.MaxUint64/2,
		DisputeTimeout: DisputeTimeout,
	}
	paymentChannelInitializerBytes, err := proto.Marshal(initializer)
	if err != nil {
		return
	}

	siglo, err := lo.Signer.SignEthMessage(paymentChannelInitializerBytes)
	if err != nil {
		return
	}

	sighi, err := hi.Signer.SignEthMessage(paymentChannelInitializerBytes)
	if err != nil {
		return
	}

	requestBytes, err := proto.Marshal(&chain.OpenChannelRequest{
		ChannelInitializer: paymentChannelInitializerBytes,
		Sigs:               [][]byte{siglo, sighi},
	})
	if err != nil {
		return
	}

	tx, err := LedgerContract.OpenChannel(Client0.Auth, requestBytes)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	WaitMinedWithChk(ctx, EthClient, tx, BlockDelay, PollingInterval, "OpenChannel")

	receipt, err := EthClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return
	}
	channelId = receipt.Logs[0].Topics[1]
	log.Info("channel ID: ", mainchain.Bytes2Hex(channelId[:]))

	return
}

func IntendWithdraw(auth *bind.TransactOpts, candidateAddr mainchain.Addr, amt *big.Int) error {
	conn := EthClient
	dposContract := DposContract
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	log.Info("Call intendWithdraw on dpos contract using the validator eth address...")
	tx, err := dposContract.IntendWithdraw(auth, candidateAddr, amt)
	if err != nil {
		return err
	}

	WaitMinedWithChk(ctx, conn, tx, BlockDelay, PollingInterval, "IntendWithdraw")
	return nil
}

func InitializeCandidate(auth *bind.TransactOpts, sgnAddr sdk.AccAddress, minSelfStake *big.Int, commissionRate *big.Int, rateLockEndTime *big.Int) error {
	conn := EthClient
	dposContract := DposContract
	sgnContract := SgnContract

	log.Infof("Call initializeCandidate on dpos contract using the validator eth address %x, minSelfStake: %d, commissionRate: %d, rateLockEndTime: %d", auth.From.Bytes(), minSelfStake, commissionRate, rateLockEndTime)
	_, err := dposContract.InitializeCandidate(auth, minSelfStake, commissionRate, rateLockEndTime)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	log.Infof("Call updateSidechainAddr on sgn contract using the validator eth address, sgnAddr: %x", sgnAddr.Bytes())
	auth.GasLimit = 8000000
	tx, err := sgnContract.UpdateSidechainAddr(auth, sgnAddr.Bytes())
	auth.GasLimit = 0
	if err != nil {
		return err
	}
	WaitMinedWithChk(ctx, conn, tx, BlockDelay, PollingInterval, "SGN UpdateSidechainAddr")

	return nil
}

func DelegateStake(fromAuth *bind.TransactOpts, toEthAddress mainchain.Addr, amt *big.Int) error {
	conn := EthClient
	dposContract := DposContract
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	log.Info("Call delegate on dpos contract to delegate stake to the validator eth address...")
	_, err := E2eProfile.CelrContract.Approve(fromAuth, E2eProfile.DPoSAddr, amt)
	if err != nil {
		return err
	}

	fromAuth.GasLimit = 8000000
	tx, err := dposContract.Delegate(fromAuth, toEthAddress, amt)
	fromAuth.GasLimit = 0
	if err != nil {
		return err
	}
	WaitMinedWithChk(ctx, conn, tx, BlockDelay, PollingInterval, "Delegate to validator")
	return nil
}

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
