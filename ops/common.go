package ops

import (
	"math/big"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/spf13/viper"
)

const (
	candidateFlag = "candidate"
	amountFlag    = "amount"
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
			AddGasPriceGwei:      viper.GetUint64(common.FlagEthAddGasPriceGwei),
			MinGasPriceGwei:      viper.GetUint64(common.FlagEthMinGasPriceGwei),
		},
		viper.GetString(common.FlagEthDPoSAddress),
		viper.GetString(common.FlagEthSGNAddress),
		viper.GetString(common.FlagEthLedgerAddress),
	)
}

func calcRawAmount() *big.Int {
	amount := new(big.Int)
	amount.SetString(viper.GetString(amountFlag), 10)
	return new(big.Int).Mul(amount, big.NewInt(common.TokenDec))
}
