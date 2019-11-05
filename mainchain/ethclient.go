package mainchain

import (
	"crypto/ecdsa"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type EthClient struct {
	PrivateKey    *ecdsa.PrivateKey
	Address       ethcommon.Address
	Client        *ethclient.Client
	GuardAddress  ethcommon.Address
	Guard         *Guard
	LedgerAddress ethcommon.Address
	Ledger        *CelerLedger
	Auth          *bind.TransactOpts
}

// Get a new eth client
func NewEthClient(ws, guardAddrStr, ledgerAddrStr, ks, passphrase string) (*EthClient, error) {
	rpcClient, err := ethrpc.Dial(ws)
	if err != nil {
		return nil, err
	}

	client := ethclient.NewClient(rpcClient)

	guardAddress := ethcommon.HexToAddress(guardAddrStr)
	guard, err := NewGuard(guardAddress, client)
	if err != nil {
		return nil, err
	}

	ledgerAddress := ethcommon.HexToAddress(ledgerAddrStr)
	ledger, err := NewCelerLedger(ledgerAddress, client)
	if err != nil {
		return nil, err
	}

	ethClient := &EthClient{
		Client:        client,
		GuardAddress:  guardAddress,
		Guard:         guard,
		LedgerAddress: ledgerAddress,
		Ledger:        ledger,
	}
	err = ethClient.setupAuth(ks, passphrase)
	if err != nil {
		return nil, err
	}

	return ethClient, nil
}

func (ethClient *EthClient) setupAuth(ks, passphrase string) error {
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
