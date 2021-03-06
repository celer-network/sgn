package slash

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	Params       Params    `json:"params" yaml:"params"`
	Penalties    []Penalty `json:"penalties" yaml:"penalties"`
	PenaltyNonce uint64    `json:"penalty_nonce" yaml:"penalty_nonce"`
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams())
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetPenaltyNonce(ctx, data.PenaltyNonce)

	for _, penalty := range data.Penalties {
		keeper.SetPenalty(ctx, penalty)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	penaltyNonce := keeper.GetPenaltyNonce(ctx)
	penalties := keeper.GetPenalties(ctx)

	return GenesisState{
		Params:       params,
		PenaltyNonce: penaltyNonce,
		Penalties:    penalties,
	}
}
