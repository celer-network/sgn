package testcommon

const (
	// common sgn parameters
	SgnChainID    = "sgnchain"
	SgnPassphrase = "12341234"
	SgnGasPrice   = ""

	SgnNodeURI    = "tcp://localhost:26657"
	SgnTransactor = "cosmos1sev3nak38elq95lumnh6t2drjtfx73274vnnjh"
	SgnCLIHome    = "../../../docker-volumes/node0/sgncli"
)

// An array isn't immutable by nature in Golang, so we define the following as var
var (
	SgnOperators   = [...]string{"cosmos1ddvpnk98da5hgzz8lf5y82gnsrhvu3jd3cukpp", "cosmos1lh8cr9p2a9dxtunte0sn3qmjkyksdh5yc5yxph", "cosmos122w97t8vsa3538fr3ylvz3hvuqxrgpnax8es8f"}
	EthKeystores   = [...]string{"../../keys/client0.json", "../../keys/client1.json", "../../keys/client2.json"}
	EthKeystorePps = [...]string{"", "", ""}
	EthAddresses   = [...]string{"6a6d2a97da1c453a4e099e8054865a0a59728863", "ba756d65a1a03f07d205749f35e2406e4a8522ad", "f25d8b54fad6e976eb9175659ae01481665a2254"}
)
