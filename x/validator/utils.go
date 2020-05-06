package validator

import (
	"math/big"

	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func GetDPoSCandidateInfoFromMainchain(ctx sdk.Context, keeper Keeper, ethAddress string) (mainchain.DPoSCandidateInfo, error) {
	ethBlknum := new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx))
	ethAddr := mainchain.Hex2Addr(ethAddress)

	return keeper.ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{BlockNumber: ethBlknum}, ethAddr)
}
