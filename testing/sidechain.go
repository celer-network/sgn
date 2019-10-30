package testing

import (
	"log"

	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/flags"
	"github.com/spf13/viper"
)

var (
	Transactor *common.Transactor
)

func SetupTransactor() {
	cdc := app.MakeCodec()
	t, err := common.NewTransactor(
		app.DefaultCLIHome,
		viper.GetString(flags.FlagSgnChainID),
		viper.GetString(flags.FlagSgnNodeURI),
		viper.GetString(flags.FlagSgnTransactors),
		viper.GetString(flags.FlagSgnPassphrase),
		viper.GetString(flags.FlagSgnGasPrice),
		cdc,
	)
	if err != nil {
		log.Fatal(err)
	}

	Transactor = t
}
