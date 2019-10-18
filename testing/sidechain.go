package testing

import (
	"log"

	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/flags"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
	"github.com/spf13/viper"
)

var (
	EthClient  *mainchain.EthClient
	Transactor *utils.Transactor
)

func SetupEthClient() {
	ec, err := mainchain.NewEthClient(
		viper.GetString(flags.FlagEthWS),
		viper.GetString(flags.FlagEthGuardAddress),
		viper.GetString(flags.FlagEthLedgerAddress),
		viper.GetString(flags.FlagEthKeystore),
		viper.GetString(flags.FlagEthPassphrase),
	)
	if err != nil {
		log.Fatal(err)
	}

	EthClient = ec
}

func SetupTransactor() {
	cdc := app.MakeCodec()
	t, err := utils.NewTransactor(
		app.DefaultCLIHome,
		viper.GetString(flags.FlagSgnChainID),
		viper.GetString(flags.FlagSgnNodeURI),
		"jack",
		viper.GetString(flags.FlagSgnPassphrase),
		cdc,
	)
	if err != nil {
		log.Fatal(err)
	}

	Transactor = t
}
