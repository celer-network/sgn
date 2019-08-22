package main

import (
	"log"

	"github.com/celer-network/sgn/chain"

	app "github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/entity"
	"github.com/celer-network/sgn/flags"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
	"github.com/celer-network/sgn/x/guardianmanager"
	"github.com/celer-network/sgn/x/subscribe"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/proto"
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

	ethClient, err = mainchain.NewEthClient(
		viper.GetString(flags.FlagEthWS),
		viper.GetString(flags.FlagEthGuardAddress),
		viper.GetString(flags.FlagEthLedgerAddress),
		viper.GetString(flags.FlagEthKeystore),
		viper.GetString(flags.FlagEthPassphrase),
	)
	if err != nil {
		log.Fatal(err)
	}

	setupTransactor()
	// sendSubscribeTx()
	sendRequestGuardTx()
}

func setupTransactor() {
	cdc := app.MakeCodec()
	t, err := utils.NewTransactor(
		app.DefaultCLIHome,
		viper.GetString(flags.FlagSgnChainID),
		viper.GetString(flags.FlagSgnNodeURI),
		viper.GetString(flags.FlagSgnName),
		viper.GetString(flags.FlagSgnPassphrase),
		cdc,
	)
	if err != nil {
		log.Fatal(err)
	}

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

func sendRequestGuardTx() {
	channelId := [32]byte{}
	copy(channelId[:], []byte{1})
	log.Println(channelId)
	log.Println(ethcommon.Bytes2Hex(channelId[:]))

	simplexPaymentChannelBytes, err := proto.Marshal(&entity.SimplexPaymentChannel{
		SeqNum:    10,
		ChannelId: channelId[:],
	})
	if err != nil {
		log.Fatal(err)
	}

	signedSimplexStateBytes, err := proto.Marshal(&chain.SignedSimplexState{
		SimplexState: simplexPaymentChannelBytes,
	})
	if err != nil {
		log.Fatal(err)
	}

	msg := guardianmanager.NewMsgRequestGuard(ethClient.Address.String(), signedSimplexStateBytes, transactor.Key.GetAddress())
	res, err := transactor.BroadcastTx(msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
}
