package validator

import (
	"math/big"

	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func GetCandidateInfoFromMainchain(ctx sdk.Context, keeper Keeper, ethAddress string) (mainchain.CandidateInfo, error) {
	return keeper.ethClient.Guard.GetCandidateInfo(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, ethcommon.HexToAddress(ethAddress))
}
