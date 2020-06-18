package mainchain

import (
	"encoding/hex"
	"io/ioutil"
	"math/big"

	"github.com/celer-network/goutils/eth"
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
	DPoS          *DPoS
	SGNAddress    Addr
	SGN           *SGN
	LedgerAddress Addr
	Ledger        *CelerLedger
}

type TransactorConfig struct {
	BlockDelay           uint64
	BlockPollingInterval uint64
	ChainId              *big.Int
}

func NewEthClient(
	ethurl, ksfile, passphrase string,
	transactorConfig *TransactorConfig,
	dposAddrStr, sgnAddrStr, ledgerAddrStr string) (*EthClient, error) {
	ethClient := &EthClient{}

	rpcClient, err := ethrpc.Dial(ethurl)
	if err != nil {
		return nil, err
	}
	ethClient.Client = ethclient.NewClient(rpcClient)

	ksBytes, err := ioutil.ReadFile(ksfile)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(ksBytes, passphrase)
	if err != nil {
		return nil, err
	}
	ethClient.Address = key.Address

	ethClient.Signer, err = eth.NewSigner(hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)))
	if err != nil {
		return nil, err
	}

	transactor, err := eth.NewTransactor(string(ksBytes), passphrase, ethClient.Client)
	if err != nil {
		return nil, err
	}
	if transactorConfig != nil {
		eth.SetBlockDelay(transactorConfig.BlockDelay)
		eth.SetBlockPollingInterval(transactorConfig.BlockPollingInterval)
		eth.SetChainId(transactorConfig.ChainId)
		// TODO: GasLimit and WaitMinedConfig
	}
	ethClient.Transactor = transactor

	err = ethClient.setContracts(dposAddrStr, sgnAddrStr, ledgerAddrStr)
	if err != nil {
		return nil, err
	}

	return ethClient, nil
}

func (ethClient *EthClient) setContracts(dposAddrStr, sgnAddrStr, ledgerAddrStr string) error {
	ethClient.DPoSAddress = Hex2Addr(dposAddrStr)
	dpos, err := NewDPoS(ethClient.DPoSAddress, ethClient.Client)
	if err != nil {
		return err
	}

	ethClient.SGNAddress = Hex2Addr(sgnAddrStr)
	sgn, err := NewSGN(ethClient.SGNAddress, ethClient.Client)
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

func (ethClient *EthClient) SignMessage(data []byte) ([]byte, error) {
	return ethClient.Signer.SignEthMessage(data)
}
