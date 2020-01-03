package testing

import (
	"math/big"
	"testing"

	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type SGNParams struct {
	CelrAddr               mainchain.Addr
	BlameTimeout           *big.Int
	MinValidatorNum        *big.Int
	MinStakingPool         *big.Int
	SidechainGoLiveTimeout *big.Int
	StartGateway           bool
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

func AddValidator(celrContract *mainchain.ERC20, guardAddr mainchain.Addr, fromAuth *bind.TransactOpts, toEthAddress mainchain.Addr, sgnAddr sdk.AccAddress, amt *big.Int) error {
	err := InitializeCandidate(fromAuth, sgnAddr)
	if err != nil {
		return err
	}

	err = DelegateStake(celrContract, guardAddr, fromAuth, toEthAddress, amt)
	if err != nil {
		return err
	}

	return nil
}
