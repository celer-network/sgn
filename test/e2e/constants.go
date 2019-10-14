package e2e

const (
	// // outPathPrefix is the path prefix for all output from e2e (incl. chain data, binaries etc)
	// // the code will append epoch second to this and create the folder
	// // the folder will be deleted after test ends successfully
	outRootDirPrefix = "/tmp/celer_e2e_"

	// etherbase and osp addr/priv key in hex
	etherBaseAddrStr = "b5bb8b7f6f1883e0c01ffb8697024532e6f3238c"
	etherBasePriv    = "69ef4da8204644e354d759ca93b94361474259f63caac6e12d7d0abcca0063f8"
	serverAddrStr    = "6a6d2a97da1c453a4e099e8054865a0a59728863"
	serverPriv       = "a7c9fa8bcd45a86fdb5f30fecf88337f20185b0c526088f2b8e0f726cad12857"
	server2AddrStr   = "ba756d65a1a03f07d205749f35e2406e4a8522ad"
	server2Priv      = "c2ff7d4ce25f7448de00e21bbbb7b884bb8dc0ca642031642863e78a35cb933d"

	ethGateway = "http://127.0.0.1:8545"

	// try to do some allocation for port: 10xyz are osp,
	// x is osp 0-based index
	// yz are osp ports like grpc, adminweb, selfrpc etc
	sEthAddr = serverAddrStr
	sPort    = "10000"
	sSelfRPC = "localhost:10010"

	s2Addr    = "127.0.0.1:10100"
	s2Port    = "10100"
	s2SelfRPC = "localhost:10110"

	localhost = "127.0.0.1:"

	stSqlAddr = "postgresql://celer_e2e_test@localhost:26257/celer_e2e_test?sslmode=disable"

	startDbScript    = "./start_db.sh"
	stopDbScript     = "./stop_db.sh"
	cockroachDbImage = "cockroachdb/cockroach:v19.1.2"

	sStoreDirPrefix         = "/tmp/sStore"
	sShadowDir              = "/tmp/sShadowStore"
	c1StoreDir              = "/tmp/c1Store"
	c1StoreSettleDisputeDir = "/tmp/c1StoreSettleDispute"
	c2StoreDir              = "/tmp/c2Store"
	oracleStoreDir          = "/tmp/oracleStore"

	refillRecord = "/tmp/refillRecord.json"

	logName = "serverlog"

	sendAmt = "1"

	depositAmt         = "5000000000000000000"  // 5 ETH
	accountBalance     = "50000000000000000000" // 50 ETH
	initialFreeBalance = "5000000000000000000"

	noProxyWithBlockDelayProfile = "../../testing/profile/no_proxy_with_block_delay.json"
	c2NoProxyMultiServerProfile  = "../../testing/profile/c2_no_proxy_multi_server.json"

	rtConfig       = "../../testing/profile/rt_config.json"
	ospOspRtConfig = "../../testing/profile/osp_osp_rt_config.json"
	tokensConfig   = "../../testing/profile/tokens.json"

	sKeystore = "../../testing/env/keystore/server.json"
)
