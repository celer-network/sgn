package testing

import (
	"log"

	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/flags"
	"github.com/spf13/viper"
)

var (
	Transactor *transactor.Transactor
)

func SetupTransactor() {
	cdc := app.MakeCodec()
	t, err := transactor.NewTransactor(
		app.DefaultCLIHome,
		viper.GetString(flags.FlagSgnChainID),
		viper.GetString(flags.FlagSgnNodeURI),
		viper.GetStringSlice(flags.FlagSgnTransactors)[0],
		viper.GetString(flags.FlagSgnPassphrase),
		viper.GetString(flags.FlagSgnGasPrice),
		cdc,
	)
	if err != nil {
		log.Fatal(err)
	}

	Transactor = t
}
