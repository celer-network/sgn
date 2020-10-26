package guard

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	Params        Params         `json:"params" yaml:"params"`
	Subscriptions []Subscription `json:"subscriptions" yaml:"subscriptions"`
	Requests      []Request      `json:"requests" yaml:"requests"`
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

	for _, subscription := range data.Subscriptions {
		keeper.SetSubscription(ctx, subscription)
	}

	for _, request := range data.Requests {
		keeper.SetRequest(ctx, request)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	subscriptions := keeper.GetSubscriptions(ctx)
	requests := keeper.GetRequests(ctx)

	return GenesisState{
		Params:        params,
		Subscriptions: subscriptions,
		Requests:      requests,
	}
}
