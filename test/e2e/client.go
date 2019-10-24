package e2e

import (
	"log"
	"time"

	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/flags"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	"github.com/celer-network/sgn/utils"
	"github.com/celer-network/sgn/x/subscribe"
	ethcommon "github.com/ethereum/go-ethereum/common"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

var (
	ethClient  *mainchain.EthClient
	transactor *utils.Transactor
)

func main() {
	viper.SetConfigFile("config.json")
	err := viper.MergeInConfig()
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
		app.DefaultCLIHome, // "$HOME/.sgncli"
		viper.GetString(flags.FlagSgnChainID),
		viper.GetString(flags.FlagSgnNodeURI),
		"alice",
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
	transactor.BroadcastTx(msg)
}

func sendRequestGuardTx() {
	channelId := [32]byte{}
	copy(channelId[:], []byte{1})
	log.Println(channelId)
	log.Println(ethcommon.Bytes2Hex(channelId[:]))

	// TODO: currently, only use two same address
	simplexPaymentChannelBytes, err := protobuf.Marshal(&entity.SimplexPaymentChannel{
		SeqNum:    10,
		ChannelId: channelId[:],
		PeerFrom:  ethClient.Address.Bytes(),
	})
	if err != nil {
		log.Fatal(err)
	}

	sig, err := mainchain.SignMessage(ethClient.PrivateKey, simplexPaymentChannelBytes)
	if err != nil {
		log.Fatal(err)
	}

	signedSimplexStateBytes, err := protobuf.Marshal(&chain.SignedSimplexState{
		SimplexState: simplexPaymentChannelBytes,
		Sigs:         [][]byte{sig, sig},
	})
	if err != nil {
		log.Fatal(err)
	}

	msg := subscribe.NewMsgRequestGuard(ethClient.Address.String(), signedSimplexStateBytes, transactor.Key.GetAddress())
	transactor.BroadcastTx(msg)
	time.Sleep(2 * time.Second)
}
