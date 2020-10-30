package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	Params         Params          `json:"params" yaml:"params"`
	Syncer         Syncer          `json:"syncer" yaml:"syncer"`
	Candidates     []Candidate     `json:"candidates" yaml:"candidates"`
	Delegators     []Delegator     `json:"delegators" yaml:"delegators"`
	Rewards        []Reward        `json:"rewards" yaml:"rewards"`
	PendingRewards []PendingReward `json:"pending_rewards" yaml:"pending_rewards"`
	RewardEpoch    RewardEpoch     `json:"reward_epoch" yaml:"reward_epoch"`
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
	if !data.Syncer.ValidatorAddr.Empty() {
		keeper.SetSyncer(ctx, data.Syncer)
	}

	keeper.SetRewardEpoch(ctx, data.RewardEpoch)

	for _, candidate := range data.Candidates {
		keeper.SetCandidate(ctx, candidate)
	}

	for _, delegator := range data.Delegators {
		keeper.SetDelegator(ctx, delegator)
	}

	for _, reward := range data.Rewards {
		keeper.SetReward(ctx, reward)
	}

	for _, pendingReward := range data.PendingRewards {
		keeper.SetPendingReward(ctx, pendingReward)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	syncer := keeper.GetSyncer(ctx)
	candidates := keeper.GetAllCandidates(ctx)
	delegators := []Delegator{}
	rewards := keeper.GetRewards(ctx)
	pendingRewards := keeper.GetPendingRewards(ctx)
	rewardEpoch := keeper.GetRewardEpoch(ctx)

	for _, candidate := range candidates {
		delegators = append(delegators, keeper.GetAllDelegators(ctx, candidate.EthAddress)...)
	}

	return GenesisState{
		Params:         params,
		Syncer:         syncer,
		Candidates:     candidates,
		Delegators:     delegators,
		Rewards:        rewards,
		PendingRewards: pendingRewards,
		RewardEpoch:    rewardEpoch,
	}
}
