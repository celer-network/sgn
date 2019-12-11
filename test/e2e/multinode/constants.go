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

	// sgnCLIHome
	sgnCLIHome0 = "../../../docker-volumes/node0/sgncli"
	sgnCLIHome1 = "../../../docker-volumes/node1/sgncli"
	sgnCLIHome2 = "../../../docker-volumes/node2/sgncli"
)
