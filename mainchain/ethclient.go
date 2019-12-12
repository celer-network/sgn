package mainchain

import (
	"crypto/ecdsa"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type EthClient struct {
	// set by SetupClient
	Client *ethclient.Client
	// set by SetupAuth
	PrivateKey *ecdsa.PrivateKey
	Address    Addr
	Auth       *bind.TransactOpts
	// set by SetupContract
	GuardAddress  Addr
	Guard         *Guard
	LedgerAddress Addr
	Ledger        *CelerLedger
}

// Get a new eth client
func NewEthClient(ws, guardAddrStr, ledgerAddrStr, ks, passphrase string) (*EthClient, error) {
	ethClient := &EthClient{}
	err := ethClient.SetupClient(ws)
	if err != nil {
		return nil, err
	}

	err = ethClient.SetupAuth(ks, passphrase)
	if err != nil {
		return nil, err
	}

	err = ethClient.SetupContract(guardAddrStr, ledgerAddrStr)
	if err != nil {
		return nil, err
	}

	return ethClient, nil
}

func (ethClient *EthClient) SetupClient(ws string) error {
	rpcClient, err := ethrpc.Dial(ws)
	if err != nil {
		return err
	}

	ethClient.Client = ethclient.NewClient(rpcClient)
	return nil
}

func (ethClient *EthClient) SetupAuth(ks, passphrase string) error {
	keystoreBytes, err := ioutil.ReadFile(ks)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keystoreBytes, passphrase)
	if err != nil {
		return err
	}

	auth, err := bind.NewTransactor(strings.NewReader(string(keystoreBytes)), passphrase)
	if err != nil {
		return err
	}

	ethClient.PrivateKey = key.PrivateKey
	ethClient.Address = key.Address
	ethClient.Auth = auth
	return nil
}

func (ethClient *EthClient) SetupContract(guardAddrStr, ledgerAddrStr string) error {
	ethClient.GuardAddress = Hex2Addr(guardAddrStr)
	guard, err := NewGuard(ethClient.GuardAddress, ethClient.Client)
	if err != nil {
		return err
	}

	ethClient.LedgerAddress = Hex2Addr(ledgerAddrStr)
	ledger, err := NewCelerLedger(ethClient.LedgerAddress, ethClient.Client)
	if err != nil {
		return err
	}

	ethClient.Guard = guard
	ethClient.Ledger = ledger
	return nil
}
