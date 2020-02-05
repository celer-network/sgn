package osp

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
	osp        *mainchain.EthClient
	user       *mainchain.EthClient
	cdc        *codec.Codec
	channelID  mainchain.CidType
	gateway    string
}

const (
	userFlag       = "user"
	ospFlag        = "osp"
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

	user, err := mainchain.NewEthClient(viper.GetString(common.FlagEthWS), viper.GetString(common.FlagEthGuardAddress), viper.GetString(common.FlagEthLedgerAddress), viper.GetString(userFlag), "")
	if err != nil {
		return
	}

	osp, err := mainchain.NewEthClient(viper.GetString(common.FlagEthWS), viper.GetString(common.FlagEthGuardAddress), viper.GetString(common.FlagEthLedgerAddress), viper.GetString(ospFlag), "")
	if err != nil {
		return
	}

	tc.DefaultTestEthClient = user
	channelID, err := tc.OpenChannel(user.Address, osp.Address, user.PrivateKey, osp.PrivateKey)
	if err != nil {
		return
	}

	log.Infof("Subscribe to sgn")
	amt := new(big.Int)
	amt.SetString("1"+strings.Repeat("0", 19), 10)
	tx, err := user.Guard.Subscribe(user.Auth, amt)
	if err != nil {
		return
	}
	tc.WaitMinedWithChk(context.Background(), user.Client, tx, viper.GetUint64(blockDelayFlag), "Subscribe on Guard contract")

	if gateway == "" {
		msgSubscribe := subscribe.NewMsgSubscribe(user.Address.Hex(), ts.Key.GetAddress())
		ts.AddTxMsg(msgSubscribe)
	} else {
		reqBody, err := json.Marshal(map[string]string{
			"ethAddr": user.Address.Hex(),
		})
		if err != nil {
			return nil, err
		}
		_, err = http.Post(rs.gateway+"/subscribe/subscribe",
			"application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, err
		}
	}

	return &RestServer{
		Mux:        r,
		logger:     logger,
		transactor: ts,
		cdc:        cdc,
		osp:        osp,
		user:       user,
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
		Use:   "osp",
		Short: "Start a local REST server talking to osp",
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

	cmd.Flags().String(userFlag, "./test/keys/client0.json", "user keystore path")
	cmd.Flags().String(ospFlag, "./test/keys/client1.json", "osp keystore path")
	cmd.Flags().String(gatewayFlag, "", "gateway url")
	cmd.Flags().Uint64(blockDelayFlag, 5, "block delay")
	return sdkFlags.RegisterRestServerFlags(cmd)
}
