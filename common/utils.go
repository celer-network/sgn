package common

import (
	"math/big"
	"time"

	"github.com/celer-network/sgn/mainchain"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

const (
	retryTimeout = 500 * time.Millisecond
)

func RobustQuery(cliCtx context.CLIContext, route string) ([]byte, error) {
	res, _, err := cliCtx.Query(route)
	if err != nil {
		time.Sleep(retryTimeout)
		res, _, err = cliCtx.Query(route)
		return res, err
	}

	return res, err
}

func RobustQueryWithData(cliCtx context.CLIContext, route string, data []byte) ([]byte, error) {
	res, _, err := cliCtx.QueryWithData(route, data)
	if err != nil {
		time.Sleep(retryTimeout)
		res, _, err = cliCtx.QueryWithData(route, data)
		return res, err
	}

	return res, err
}

func NewCommission(ethClient *mainchain.EthClient, commissionRate *big.Int) (staking.Commission, error) {
	commissionBase, err := ethClient.DPoS.COMMISSIONRATEBASE(&bind.CallOpts{})
	if err != nil {
		return staking.Commission{}, err
	}

	prec := int64(len(commissionBase.String()) - 1)
	return staking.Commission{
		CommissionRates: staking.CommissionRates{
			Rate:          sdk.NewDecFromBigIntWithPrec(commissionRate, prec),
			MaxRate:       sdk.NewDec(1),
			MaxChangeRate: sdk.NewDec(1),
		},
	}, nil
}
