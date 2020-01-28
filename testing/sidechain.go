package testing

import (
	"math/big"
	"testing"

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
	MaxValidatorNum        *big.Int
}

func NewTransactor(t *testing.T, sgnCLIHome, sgnChainID, sgnNodeURI, sgnTransactor, sgnPassphrase, sgnGasPrice string) *transactor.Transactor {
	cdc := app.MakeCodec()
	tr, err := transactor.NewTransactor(
		sgnCLIHome,
		sgnChainID,
		sgnNodeURI,
		sgnTransactor,
		sgnPassphrase,
		sgnGasPrice,
		cdc,
	)
	ChkTestErr(t, err, "Failed to create new transactor.")

	return tr
}
