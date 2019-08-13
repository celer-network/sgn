package mainchain

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
)

type EthClient struct {
	Client *ethclient.Client
	Guard  *Guard
}

func NewEthClient() (*EthClient, error) {
	rpcClient, err := ethrpc.Dial(viper.GetString("ethWs"))
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(rpcClient)

	guard, err := NewGuard(ethcommon.HexToAddress(viper.GetString("guardAddress")), client)
	if err != nil {
		return nil, err
	}

	return &EthClient{
		Client: client,
		Guard:  guard,
	}, nil
}
