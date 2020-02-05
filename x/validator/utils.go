package validator

import (
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func GetCandidateInfoFromMainchain(ctx sdk.Context, keeper Keeper, ethAddress string) (mainchain.CandidateInfo, error) {
	return keeper.ethClient.Guard.GetCandidateInfo(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, mainchain.Hex2Addr(ethAddress))
}

func InitAccount(ctx sdk.Context, keeper Keeper, accAddress sdk.AccAddress) {
	account := keeper.accountKeeper.GetAccount(ctx, accAddress)
	if account == nil {
		log.Infof("Set new account %x", accAddress)
		account = keeper.accountKeeper.NewAccountWithAddress(ctx, accAddress)
		keeper.accountKeeper.SetAccount(ctx, account)
	}
}
