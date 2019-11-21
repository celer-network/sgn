package gateway

import (
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/transactor"
	"github.com/cosmos/cosmos-sdk/client"
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
	Mux            *mux.Router
	transactorPool *transactor.TransactorPool
	listener       net.Listener
	logger         tlog.Logger
}

// NewRestServer creates a new rest server instance
func NewRestServer(cdc *codec.Codec) (*RestServer, error) {
	viper.SetConfigFile("config.json")
	err := viper.MergeInConfig()
	if err != nil {
		return nil, err
	}

	transactorPool, err := transactor.NewTransactorPool(
		app.DefaultCLIHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetString(common.FlagSgnGasPrice),
		viper.GetStringSlice(common.FlagSgnTransactors),
		cdc,
	)
	if err != nil {
		return nil, err
	}

	log.SetLevelByName(viper.GetString(common.FlagLogLevel))
	if viper.GetBool(common.FlagLogColor) {
		log.EnableColor()
	}
	if viper.GetBool(common.FlagLogLongFile) {
		log.EnableLongFile()
		_, file, _, ok := runtime.Caller(0)
		if ok {
			pref := file[:strings.LastIndex(file[:strings.LastIndex(file, "/")], "/")+1]
			log.SetFilePathSplit(pref)
		}
	}

	r := mux.NewRouter()
	logger := tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)).With("module", "rest-server")

	return &RestServer{
		Mux:            r,
		transactorPool: transactorPool,
		logger:         logger,
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

func (rs *RestServer) registerRoutes() {
	client.RegisterRoutes(rs.transactorPool.GetTransactor().CliCtx, rs.Mux)
	rs.registerQueryRoutes()
	rs.registerTxRoutes()
}

// ServeCommand will start the application REST service as a blocking process. It
// takes a codec to create a RestServer object and a function to register all
// necessary routes.
func ServeCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "Start a local REST server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rs, err := NewRestServer(cdc)
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

	return sdkFlags.RegisterRestServerFlags(cmd)
}
