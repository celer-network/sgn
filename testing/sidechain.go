package testing

import (
	"log"

	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/flags"
	"github.com/celer-network/sgn/utils"
	"github.com/spf13/viper"
)

var (
	Transactor *utils.Transactor
)

func SetupTransactor() {
	cdc := app.MakeCodec()
	t, err := utils.NewTransactor(
		app.DefaultCLIHome,
		viper.GetString(flags.FlagSgnChainID),
		viper.GetString(flags.FlagSgnNodeURI),
		"alice", // use a different account from sgn's transactor
		viper.GetString(flags.FlagSgnPassphrase),
		cdc,
	)
	if err != nil {
		log.Fatal(err)
	}

	Transactor = t
}
