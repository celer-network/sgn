package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetEthBlkNum(ctx, data.EthBlkNum)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	ethBlkNum := k.GetEthBlkNum(ctx)

	return GenesisState{
		EthBlkNum: ethBlkNum,
	}
}
