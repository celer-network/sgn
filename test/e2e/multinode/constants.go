package multinode

const (
	// common sgn parameters
	sgnChainID    = "sgnchain"
	sgnPassphrase = "12341234"
	sgnGasPrice   = ""

	// customized sgn parameters
	sgnNode0URI = "tcp://localhost:26657"
	sgnNode1URI = "tcp://localhost:26660"
	sgnNode2URI = "tcp://localhost:26662"

	// sgnTransactor
	sgnTransactor0 = "cosmos104gnaehlj7a32r9l5za7nhf463yvgnyfnwh9sc"
	sgnTransactor1 = "cosmos1k88dweqanvd3j2csaak034amvjuyzlkefgmv8t"
	sgnTransactor2 = "cosmos190xknjnht45ypumxcq7da4plsnuuvf0ajfselv"

	// sgnOperator
	sgnOperator0 = "cosmos1ddvpnk98da5hgzz8lf5y82gnsrhvu3jd3cukpp"
	sgnOperator1 = "cosmos1lh8cr9p2a9dxtunte0sn3qmjkyksdh5yc5yxph"
	sgnOperator2 = "cosmos122w97t8vsa3538fr3ylvz3hvuqxrgpnax8es8f"

	// sgnCLIHome
	sgnCLIHome0 = "../../../docker-volumes/node0/sgncli"
	sgnCLIHome1 = "../../../docker-volumes/node1/sgncli"
	sgnCLIHome2 = "../../../docker-volumes/node2/sgncli"

	// Ethereum keystore paths
	ethKeystore0 = "../../keys/client0.json"
	ethKeystore1 = "../../keys/client1.json"
	ethKeystore2 = "../../keys/client2.json"

	// Ethereum keystore passphrases
	ethKeystorePp0 = ""
	ethKeystorePp1 = ""
	ethKeystorePp2 = ""

	// Ethereum addresses
	ethAddress0 = "6a6d2a97da1c453a4e099e8054865a0a59728863"
	ethAddress1 = "ba756d65a1a03f07d205749f35e2406e4a8522ad"
	ethAddress2 = "f25d8b54fad6e976eb9175659ae01481665a2254"
)
