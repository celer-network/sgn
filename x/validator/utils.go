package validator

import (
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func GetCandidateInfoFromMainchain(ctx sdk.Context, keeper Keeper, ethAddress string) (mainchain.CandidateInfo, error) {
	var candidateInfo mainchain.CandidateInfo

	ethBlknum := new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx))
	ethAddr := mainchain.Hex2Addr(ethAddress)

	dposCandidateInfo, err := keeper.ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{BlockNumber: ethBlknum}, ethAddr)
	if err != nil {
		return candidateInfo, err
	}

	sidechainAddr, err := keeper.ethClient.SGN.SidechainAddrMap(&bind.CallOpts{BlockNumber: ethBlknum}, ethAddr)
	if err != nil {
		return candidateInfo, err
	}

	candidateInfo.DPoSCandidateInfo = dposCandidateInfo
	candidateInfo.SidechainAddr = sidechainAddr

	return candidateInfo, nil
}

func InitAccount(ctx sdk.Context, keeper Keeper, accAddress sdk.AccAddress) {
	account := keeper.accountKeeper.GetAccount(ctx, accAddress)
	if account == nil {
		log.Infof("Set new account %x", accAddress)
		account = keeper.accountKeeper.NewAccountWithAddress(ctx, accAddress)
		keeper.accountKeeper.SetAccount(ctx, account)
	}
}
