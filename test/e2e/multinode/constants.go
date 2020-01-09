package multinode

const (
	// common sgn parameters
	sgnChainID    = "sgnchain"
	sgnPassphrase = "12341234"
	sgnGasPrice   = ""
)

// An array isn't immutable by nature in Golang, so we define the following as var
var (
	sgnNodeURIs    = [...]string{"tcp://localhost:26657", "tcp://localhost:26660", "tcp://localhost:26662"}
	sgnTransactors = [...]string{"cosmos104gnaehlj7a32r9l5za7nhf463yvgnyfnwh9sc", "cosmos1k88dweqanvd3j2csaak034amvjuyzlkefgmv8t", "cosmos190xknjnht45ypumxcq7da4plsnuuvf0ajfselv"}
	sgnOperators   = [...]string{"cosmos1ddvpnk98da5hgzz8lf5y82gnsrhvu3jd3cukpp", "cosmos1lh8cr9p2a9dxtunte0sn3qmjkyksdh5yc5yxph", "cosmos122w97t8vsa3538fr3ylvz3hvuqxrgpnax8es8f"}
	// operator addresses in cosmos valoper format
	sgnOperatorValAddrs = [...]string{"cosmosvaloper1ddvpnk98da5hgzz8lf5y82gnsrhvu3jd5vgrdj", "cosmosvaloper1lh8cr9p2a9dxtunte0sn3qmjkyksdh5yaqsndy", "cosmosvaloper122w97t8vsa3538fr3ylvz3hvuqxrgpnarnd9t6"}
	sgnCLIHomes         = [...]string{"../../../docker-volumes/node0/sgncli", "../../../docker-volumes/node1/sgncli", "../../../docker-volumes/node2/sgncli"}
	ethKeystores        = [...]string{"../../keys/client0.json", "../../keys/client1.json", "../../keys/client2.json"}
	ethKeystorePps      = [...]string{"", "", ""}
	ethAddresses        = [...]string{"6a6d2a97da1c453a4e099e8054865a0a59728863", "ba756d65a1a03f07d205749f35e2406e4a8522ad", "f25d8b54fad6e976eb9175659ae01481665a2254"}
)
