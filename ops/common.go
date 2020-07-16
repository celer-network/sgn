package ops

import (
	"math/big"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/spf13/viper"
)

func initEthClient() (*mainchain.EthClient, error) {
	return mainchain.NewEthClient(
		viper.GetString(common.FlagEthGateway),
		viper.GetString(common.FlagEthKeystore),
		viper.GetString(common.FlagEthPassphrase),
		&mainchain.TransactorConfig{
			BlockDelay:           viper.GetUint64(common.FlagEthBlockDelay),
			BlockPollingInterval: viper.GetUint64(common.FlagEthPollInterval),
			ChainId:              big.NewInt(viper.GetInt64(common.FlagEthChainID)),
		},
		viper.GetString(common.FlagEthDPoSAddress),
		viper.GetString(common.FlagEthSGNAddress),
		viper.GetString(common.FlagEthLedgerAddress),
	)
}
