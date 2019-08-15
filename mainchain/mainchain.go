package mainchain

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type EthClient struct {
	Address ethcommon.Address
	Client  *ethclient.Client
	Guard   *Guard
	Auth    *bind.TransactOpts
}

// Get a new eth client
func NewEthClient(ws, guardAddress, ks, passphrase string) (*EthClient, error) {
	rpcClient, err := ethrpc.Dial(ws)
	if err != nil {
		return nil, err
	}

	client := ethclient.NewClient(rpcClient)
	guard, err := NewGuard(ethcommon.HexToAddress(guardAddress), client)
	if err != nil {
		return nil, err
	}

	ethClient := &EthClient{
		Client: client,
		Guard:  guard,
	}
	ethClient.setupAuth(ks, passphrase)

	return ethClient, nil
}

func (ethClient *EthClient) GetLatestBlkNum() (uint64, error) {
	head, err := ethClient.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	latestBlkNum := head.Number.Uint64()
	return latestBlkNum, nil
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

	ethClient.Address = key.Address
	ethClient.Auth = auth
	return nil
}
