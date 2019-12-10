package testing

import (
	"math/big"

	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/spf13/viper"
)

type SGNParams struct {
	CelrAddr               mainchain.Addr
	BlameTimeout           *big.Int
	MinValidatorNum        *big.Int
	MinStakingPool         *big.Int
	SidechainGoLiveTimeout *big.Int
	StartGateway           bool
}

func NewTransactor(sgnNodeURI string) *transactor.Transactor {
	cdc := app.MakeCodec()
	t, err := transactor.NewTransactor(
		viper.GetString(common.FlagSgnCLIHome), // app.DefaultCLIHome,
		viper.GetString(common.FlagSgnChainID),
		sgnNodeURI,
		viper.GetStringSlice(common.FlagSgnTransactors)[0],
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetString(common.FlagSgnGasPrice),
		cdc,
	)
	ChkErr(err, "new transactor")

	return t
}
