package validator

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func GetCandidateInfo(ctx sdk.Context, keeper Keeper, ethAddress string) (struct {
	Initialized   bool
	MinSelfStake  *big.Int
	SidechainAddr []byte
	StakingPool   *big.Int
	IsVldt        bool
}, error) {
	return keeper.ethClient.Guard.GetCandidateInfo(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, ethcommon.HexToAddress(ethAddress))
}
