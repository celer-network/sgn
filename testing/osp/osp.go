package osp

import (
	"math/big"
	"net"
	"os"
	"strings"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/transactor"
	sdkFlags "github.com/cosmos/cosmos-sdk/client/flags"
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
	channelID  mainchain.CidType
}

const (
	userFlag = "user"
	ospFlag  = "osp"
)

// NewRestServer creates a new rest server instance
func NewRestServer() (*RestServer, error) {
	log.Infof("New rest server")
	r := mux.NewRouter()
	logger := tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)).With("module", "rest-server")
	viper.Set(sdkFlags.FlagTrustNode, true)
	cdc := app.MakeCodec()
	transactor, err := transactor.NewTransactor(
		viper.GetString(common.FlagSgnCLIHome), // app.DefaultCLIHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetStringSlice(common.FlagSgnTransactors)[0],
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetString(common.FlagSgnGasPrice),
		cdc,
	)
	if err != nil {
		return nil, err
	}

	user, err := mainchain.NewEthClient(viper.GetString(common.FlagEthWS), viper.GetString(common.FlagEthGuardAddress), viper.GetString(common.FlagEthLedgerAddress), viper.GetString(userFlag), "")
	if err != nil {
		return nil, err
	}

	osp, err := mainchain.NewEthClient(viper.GetString(common.FlagEthWS), viper.GetString(common.FlagEthGuardAddress), viper.GetString(common.FlagEthLedgerAddress), viper.GetString(ospFlag), "")
	if err != nil {
		return nil, err
	}

	tf.DefaultTestEthClient = user
	channelID, err := tf.OpenChannel(user.Address, osp.Address, user.PrivateKey, osp.PrivateKey)
	if err != nil {
		return nil, err
	}

	log.Infof("Subscribe to sgn")
	amt := new(big.Int)
	amt.SetString("1"+strings.Repeat("0", 19), 10)
	_, err = user.Guard.Subscribe(user.Auth, amt)
	if err != nil {
		return nil, err
	}

	return &RestServer{
		Mux:        r,
		logger:     logger,
		transactor: transactor,
		osp:        osp,
		user:       user,
		channelID:  channelID,
	}, nil
}

// Start starts the rest server
func (rs *RestServer) Start(listenAddr string, maxOpen int, readTimeout, writeTimeout uint) (err error) {
	server.TrapSignal(func() {
		err := rs.listener.Close()
		log.Errorln("error closing listener err", err)
	})

	cfg := rpcserver.DefaultConfig()
	cfg.MaxOpenConnections = maxOpen
	cfg.ReadTimeout = time.Duration(readTimeout) * time.Second
	cfg.WriteTimeout = time.Duration(writeTimeout) * time.Second

	rs.listener, err = rpcserver.Listen(listenAddr, cfg)
	if err != nil {
		return
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
	return sdkFlags.RegisterRestServerFlags(cmd)
}
