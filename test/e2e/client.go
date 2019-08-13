package main

import (
	"log"

	app "github.com/celer-network/sgn"
	"github.com/celer-network/sgn/flags"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

	ethClient.Guard.SubscriptionExpiration(&bind.CallOpts{}, ethClient.Address)
	// tx, err := ethClient.Guard.Subscribe(ethClient.Auth, big.NewInt(1))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	setupTransactor()
}

func setupTransactor() {
	cdc := app.MakeCodec()
	accName := viper.GetString(flags.FlagSgnName)
	transactor, err := utils.NewTransactor(
		app.DefaultCLIHome,
		accName,
		viper.GetString(flags.FlagSgnChainID),
		viper.GetString(flags.FlagSgnNodeURI),
		cdc,
	)
	if err != nil {
		log.Fatal(err)
	}

	transactor.Passphrase = viper.GetString(flags.FlagSgnPassphrase)
}
