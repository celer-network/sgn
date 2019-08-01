package mainchain

import (
	// "github.com/cosmos/sdk-application-tutorial/simple"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

var (
	EthClient     = getEthClient()
	SimpleAddress = ethcommon.HexToAddress("0x2595a2907f3a895a320f817b72ecb654e6d50004")
)

func getEthClient() *ethclient.Client {
	const address = "wss://ropsten.infura.io/ws"
	rpcClient, err := ethrpc.Dial(address)
	if err != nil {
		panic(err)
	}
	return ethclient.NewClient(rpcClient)
}
