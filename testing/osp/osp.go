package osp

import (
	"net"
	"os"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
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
	Mux       *mux.Router
	listener  net.Listener
	logger    tlog.Logger
	channelID [32]byte
}

const (
	client0Flag = "client0"
	client1Flag = "client1"
)

// NewRestServer creates a new rest server instance
func NewRestServer() (*RestServer, error) {
	r := mux.NewRouter()
	logger := tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)).With("module", "rest-server")

	tf.SetupDefaultTestEthClient(viper.GetString(common.FlagEthKeystore), viper.GetString(common.FlagEthPassphrase))
	client0 := mainchain.EthClient{}
	client0.SetAuth(viper.GetString(client0Flag), "")
	client1 := mainchain.EthClient{}
	client1.SetAuth(viper.GetString(client1Flag), "")
	channelID, err := tf.OpenChannel(client0.Address.Bytes(), client1.Address.Bytes(), client0.PrivateKey, client1.PrivateKey, []byte(viper.GetString(common.FlagConfig)))
	if err != nil {
		return nil, err
	}

	return &RestServer{
		Mux:       r,
		logger:    logger,
		channelID: channelID,
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

	cmd.Flags().String(client0Flag, "./test/keys/client0.json", "client 0 keystore path")
	cmd.Flags().String(client1Flag, "./test/keys/client1.json", "client 1 keystore path")
	return sdkFlags.RegisterRestServerFlags(cmd)
}
