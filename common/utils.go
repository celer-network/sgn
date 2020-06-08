package common

import (
	"math/big"
	"runtime"
	"strings"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

const (
	retryTimeout = 500 * time.Millisecond
)

// EnableLogLongFile set the log file splitter from the sgn root folder
func EnableLogLongFile() {
	log.EnableLongFile()
	_, file, _, ok := runtime.Caller(0)
	if ok {
		pref := file[:strings.LastIndex(file[:strings.LastIndex(file, "/")], "/")+1]
		log.SetFilePathSplit(pref)
	}
}

func RobustQuery(cliCtx context.CLIContext, route string) (res []byte, err error) {
	res, _, err = cliCtx.Query(route)
	if err != nil {
		time.Sleep(retryTimeout)
		res, _, err = cliCtx.Query(route)
		return
	}

	return
}

func RobustQueryWithData(cliCtx context.CLIContext, route string, data []byte) (res []byte, err error) {
	res, _, err = cliCtx.QueryWithData(route, data)
	if err != nil {
		time.Sleep(retryTimeout)
		res, _, err = cliCtx.QueryWithData(route, data)
		return
	}

	return
}

func NewCommission(ethClient *mainchain.EthClient, commissionRate *big.Int) (staking.Commission, error) {
	commissionBase, err := ethClient.DPoS.COMMISSIONRATEBASE(&bind.CallOpts{})
	if err != nil {
		return staking.Commission{}, err
	}

	prec := int64(len(commissionBase.String()) - 1)

	return staking.Commission{
		CommissionRates: staking.CommissionRates{
			Rate: sdk.NewDecFromBigIntWithPrec(commissionRate, prec),
		},
	}, nil
}
