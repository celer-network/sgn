package testing

import (
	"math/big"

	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
)

type SGNParams struct {
	CelrAddr               mainchain.Addr
	BlameTimeout           *big.Int
	MinValidatorNum        *big.Int
	MinStakingPool         *big.Int
	SidechainGoLiveTimeout *big.Int
	StartGateway           bool
}

func NewTransactor(sgnCLIHome, sgnChainID, sgnNodeURI, sgnTransactor, sgnPassphrase, sgnGasPrice string) *transactor.Transactor {
	cdc := app.MakeCodec()
	t, err := transactor.NewTransactor(
		sgnCLIHome,
		sgnChainID,
		sgnNodeURI,
		sgnTransactor,
		sgnPassphrase,
		sgnGasPrice,
		cdc,
	)
	ChkErr(err, "new transactor")

	return t
}
