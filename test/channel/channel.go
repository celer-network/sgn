package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/subscribe"
	sdkFlags "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tlog "github.com/tendermint/tendermint/libs/log"
	rpcserver "github.com/tendermint/tendermint/rpc/lib/server"
)

// RestServer represents the Light Client Rest server
type RestServer struct {
	Mux        *mux.Router
	listener   net.Listener
	logger     tlog.Logger
	transactor *transactor.Transactor
	peer1      *mainchain.EthClient
	peer2      *mainchain.EthClient
	cdc        *codec.Codec
	channelID  mainchain.CidType
	gateway    string
}

const (
	peer1Flag      = "peer1"
	peer2Flag      = "peer2"
	gatewayFlag    = "gateway"
	blockDelayFlag = "blockDelay"
)

// NewRestServer creates a new rest server instance
func NewRestServer() (rs *RestServer, err error) {
	log.Infof("New rest server")
	r := mux.NewRouter()
	logger := tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)).With("module", "rest-server")
	viper.Set(sdkFlags.FlagTrustNode, true)
	cdc := app.MakeCodec()
	gateway := viper.GetString(gatewayFlag)
	var ts *transactor.Transactor

	if gateway == "" {
		ts, err = transactor.NewTransactor(
			viper.GetString(sdkFlags.FlagHome),
			viper.GetString(common.FlagSgnChainID),
			viper.GetString(common.FlagSgnNodeURI),
			viper.GetString(common.FlagSgnOperator),
			viper.GetString(common.FlagSgnPassphrase),
			viper.GetString(common.FlagSgnGasPrice),
			cdc,
		)
		if err != nil {
			return
		}
	}

	peer1, err := mainchain.NewEthClient(viper.GetString(common.FlagEthInstance), viper.GetString(common.FlagEthGuardAddress), viper.GetString(common.FlagEthLedgerAddress), viper.GetString(peer1Flag), "")
	if err != nil {
		return
	}

	peer2, err := mainchain.NewEthClient(viper.GetString(common.FlagEthInstance), viper.GetString(common.FlagEthGuardAddress), viper.GetString(common.FlagEthLedgerAddress), viper.GetString(peer2Flag), "")
	if err != nil {
		return
	}

	tc.Client0 = peer1

	log.Infof("Subscribe to sgn")
	tokenAddr, err := peer1.Guard.CelerToken(&bind.CallOpts{})
	if err != nil {
		return
	}
	tokenContract, err := mainchain.NewERC20(tokenAddr, peer1.Client)
	if err != nil {
		return
	}

	amt := new(big.Int)
	amt.SetString("1"+strings.Repeat("0", 18), 10)
	tx, err := tokenContract.Approve(peer1.Auth, peer1.GuardAddress, amt)
	tc.ChkErr(err, "failed to approve erc20")
	tc.WaitMinedWithChk(context.Background(), peer1.Client, tx, 0, "approve erc20")

	tx, err = peer1.Guard.Subscribe(peer1.Auth, amt)
	tc.ChkErr(err, "failed to subscribe")
	tc.WaitMinedWithChk(context.Background(), peer1.Client, tx, viper.GetUint64(blockDelayFlag)+3, "Subscribe on Guard contract")

	if gateway == "" {
		msgSubscribe := subscribe.NewMsgSubscribe(peer1.Address.Hex(), ts.Key.GetAddress())
		ts.AddTxMsg(msgSubscribe)
	} else {
		reqBody, err2 := json.Marshal(map[string]string{
			"ethAddr": peer1.Address.Hex(),
		})
		if err2 != nil {
			return nil, err2
		}
		_, err2 = http.Post(gateway+"/subscribe/subscribe",
			"application/json", bytes.NewBuffer(reqBody))
		if err2 != nil {
			return nil, err2
		}
	}

	channelID, err := tc.OpenChannel(peer1, peer2)
	if err != nil {
		return
	}

	return &RestServer{
		Mux:        r,
		logger:     logger,
		transactor: ts,
		cdc:        cdc,
		peer1:      peer1,
		peer2:      peer2,
		channelID:  channelID,
		gateway:    gateway,
	}, nil
}

// Start starts the rest server
func (rs *RestServer) Start(listenAddr string, maxOpen int, readTimeout, writeTimeout uint) error {
	server.TrapSignal(func() {
		err := rs.listener.Close()
		log.Errorln("error closing listener err", err)
	})

	cfg := rpcserver.DefaultConfig()
	cfg.MaxOpenConnections = maxOpen
	cfg.ReadTimeout = time.Duration(readTimeout) * time.Second
	cfg.WriteTimeout = time.Duration(writeTimeout) * time.Second

	var err error
	rs.listener, err = rpcserver.Listen(listenAddr, cfg)
	if err != nil {
		return err
	}
	log.Infof("Starting application REST service (chain-id: %s)...", viper.GetString(sdkFlags.FlagChainID))

	return rpcserver.StartHTTPServer(rs.listener, rs.Mux, rs.logger, cfg)
}

// ServeCommand will start the application REST service as a blocking process. It
// takes a codec to create a RestServer object and a function to register all
// necessary routes.
func ServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channel",
		Short: "Start a local REST server talking to channel and sgn",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			viper.SetConfigFile(viper.GetString(common.FlagConfig))
			err = viper.ReadInConfig()
			if err != nil {
				return err
			}

			rs, err := NewRestServer()
			if err != nil {
				return err
			}

			rs.registerRoutes()

			// Start the rest server and return error if one exists
			err = rs.Start(
				viper.GetString(sdkFlags.FlagListenAddr),
				viper.GetInt(sdkFlags.FlagMaxOpenConnections),
				uint(viper.GetInt(sdkFlags.FlagRPCReadTimeout)),
				uint(viper.GetInt(sdkFlags.FlagRPCWriteTimeout)),
			)

			return err
		},
	}

	cmd.Flags().String(peer1Flag, "./test/keys/ethks0.json", "peer1 keystore path")
	cmd.Flags().String(peer2Flag, "./test/keys/ethks1.json", "peer2 keystore path")
	cmd.Flags().String(gatewayFlag, "", "gateway url")
	cmd.Flags().Uint64(blockDelayFlag, 5, "block delay")
	return sdkFlags.RegisterRestServerFlags(cmd)
}
