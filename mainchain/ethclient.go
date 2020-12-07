package mainchain

import (
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/sgn-contract/bindings/go/sgncontracts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type EthClient struct {
	// init by NewEthClient
	Client     *ethclient.Client
	Transactor *eth.Transactor
	Signer     eth.Signer
	Address    Addr

	// init by SetContracts
	DPoSAddress   Addr
	DPoS          *sgncontracts.DPoS
	SGNAddress    Addr
	SGN           *sgncontracts.SGN
	LedgerAddress Addr
	Ledger        *CelerLedger
}

type TransactorConfig struct {
	BlockDelay           uint64
	BlockPollingInterval uint64
	ChainId              *big.Int
	AddGasPriceGwei      uint64
	MinGasPriceGwei      uint64
}

func NewEthClient(
	ethurl string,
	ksfile string,
	passphrase string,
	tconfig *TransactorConfig,
	dposAddrStr string,
	sgnAddrStr string,
	ledgerAddrStr string) (*EthClient, error) {
	ethClient := &EthClient{}

	rpcClient, err := ethrpc.Dial(ethurl)
	if err != nil {
		return nil, err
	}

	ethClient.Client = ethclient.NewClient(rpcClient)
	err = ethClient.setContracts(dposAddrStr, sgnAddrStr, ledgerAddrStr)
	if err != nil {
		return nil, err
	}

	if ksfile != "" {
		err = ethClient.setTransactor(ksfile, passphrase, tconfig)
		if err != nil {
			return nil, err
		}
	}

	return ethClient, nil
}

func (ethClient *EthClient) setTransactor(ksfile string, passphrase string, tconfig *TransactorConfig) error {
	ksBytes, err := ioutil.ReadFile(ksfile)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(ksBytes, passphrase)
	if err != nil {
		return err
	}

	ethClient.Address = key.Address
	ethClient.Signer, err = eth.NewSigner(hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)), tconfig.ChainId)
	if err != nil {
		return err
	}

	ethClient.Transactor, err = eth.NewTransactor(
		string(ksBytes),
		passphrase,
		ethClient.Client,
		tconfig.ChainId,
		eth.WithBlockDelay(tconfig.BlockDelay),
		eth.WithPollingInterval(time.Duration(tconfig.BlockPollingInterval)*time.Second),
		eth.WithAddGasGwei(tconfig.AddGasPriceGwei),
		eth.WithMinGasGwei(tconfig.MinGasPriceGwei),
	)

	return err
}

func (ethClient *EthClient) setContracts(dposAddrStr, sgnAddrStr, ledgerAddrStr string) error {
	ethClient.DPoSAddress = Hex2Addr(dposAddrStr)
	dpos, err := sgncontracts.NewDPoS(ethClient.DPoSAddress, ethClient.Client)
	if err != nil {
		return err
	}

	ethClient.SGNAddress = Hex2Addr(sgnAddrStr)
	sgn, err := sgncontracts.NewSGN(ethClient.SGNAddress, ethClient.Client)
	if err != nil {
		return err
	}

	ethClient.LedgerAddress = Hex2Addr(ledgerAddrStr)
	ledger, err := NewCelerLedger(ethClient.LedgerAddress, ethClient.Client)
	if err != nil {
		return err
	}

	ethClient.DPoS = dpos
	ethClient.SGN = sgn
	ethClient.Ledger = ledger
	return nil
}

func (ethClient *EthClient) SignEthMessage(data []byte) ([]byte, error) {
	return ethClient.Signer.SignEthMessage(data)
}
