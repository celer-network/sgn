package gateway

import (
	"net"
	"os"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/cosmos/cosmos-sdk/client"
	sdkFlags "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tlog "github.com/tendermint/tendermint/libs/log"
	rpcserver "github.com/tendermint/tendermint/rpc/jsonrpc/server"
)

// RestServer represents the Light Client Rest server
type RestServer struct {
	Mux            *mux.Router
	transactorPool *transactor.TransactorPool
	listener       net.Listener
	logger         tlog.Logger

	dposContract   *mainchain.DPoS
	sgnContract    *mainchain.SGN
	ledgerContract *mainchain.CelerLedger
}

// NewRestServer creates a new rest server instance
func NewRestServer(cdc *codec.Codec) (*RestServer, error) {
	rpcClient, err := rpc.Dial(viper.GetString(common.FlagEthGateway))
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(rpcClient)
	dposContract, err := mainchain.NewDPoS(
		mainchain.Hex2Addr(viper.GetString(common.FlagEthDPoSAddress)), ethClient)
	if err != nil {
		return nil, err
	}
	sgnContract, err := mainchain.NewSGN(
		mainchain.Hex2Addr(viper.GetString(common.FlagEthSGNAddress)), ethClient)
	if err != nil {
		return nil, err
	}
	ledgerContract, err := mainchain.NewCelerLedger(
		mainchain.Hex2Addr(viper.GetString(common.FlagEthLedgerAddress)), ethClient)
	if err != nil {
		return nil, err
	}

	gpe := transactor.NewGasPriceEstimator(viper.GetString(common.FlagSgnNodeURI))
	transactorPool := transactor.NewTransactorPool(
		viper.GetString(sdkFlags.FlagHome),
		viper.GetString(common.FlagSgnChainID),
		cdc,
		gpe,
	)

	err = transactorPool.AddTransactors(
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetStringSlice(common.FlagSgnTransactors),
	)
	if err != nil {
		return nil, err
	}

	log.SetLevelByName(viper.GetString(common.FlagLogLevel))
	if viper.GetBool(common.FlagLogColor) {
		log.EnableColor()
	}

	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))
	logger := tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)).With("module", "rest-server")

	return &RestServer{
		Mux:            r,
		transactorPool: transactorPool,
		logger:         logger,
		dposContract:   dposContract,
		sgnContract:    sgnContract,
		ledgerContract: ledgerContract,
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

	return rpcserver.Serve(rs.listener, rs.Mux, rs.logger, cfg)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigFile("config.json")
			err := viper.MergeInConfig()
			if err != nil {
				return err
			}

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
