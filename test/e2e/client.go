package main

import (
	"log"

	app "github.com/celer-network/sgn"
	"github.com/celer-network/sgn/flags"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/spf13/viper"
)

var (
	ethClient  *mainchain.EthClient
	transactor *utils.Transactor
)

func main() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	ethClient, err = mainchain.NewEthClient()
	if err != nil {
		log.Fatal(err)
	}

	setupTransactor()
	sendSubscribeTx()
}

func setupTransactor() {
	cdc := app.MakeCodec()
	accName := viper.GetString(flags.FlagSgnName)
	t, err := utils.NewTransactor(
		app.DefaultCLIHome,
		accName,
		viper.GetString(flags.FlagSgnChainID),
		viper.GetString(flags.FlagSgnNodeURI),
		cdc,
	)
	if err != nil {
		log.Fatal(err)
	}

	t.Passphrase = viper.GetString(flags.FlagSgnPassphrase)
	transactor = t
}

func sendSubscribeTx() {
	msg := subscribe.NewMsgSubscribe(ethClient.Address.String(), transactor.Key.GetAddress())
	res, err := transactor.BroadcastTx(msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
}
