package v03

import (
	v02slash "github.com/celer-network/sgn/x/slash/legacy/v0_2"
)

func Migrate(oldGenState v02slash.GenesisState) *GenesisState {
	penalties := []Penalty{}
	for _, oldPenalty := range oldGenState.Penalties {
		penalty := NewPenalty(oldPenalty.Nonce, oldPenalty.Reason, oldPenalty.ValidatorAddr, oldPenalty.PenalizedDelegators, oldPenalty.Beneficiaries, DefaultSyncerReward)
		penalty.PenaltyProtoBytes = oldPenalty.PenaltyProtoBytes
		penalty.Sigs = oldPenalty.Sigs

		penalties = append(penalties, penalty)
	}

	return &GenesisState{
		Params: Params{
			SignedBlocksWindow:        oldGenState.Params.SignedBlocksWindow,
			PenaltyDelegatorSize:      DefaultPenaltyDelegatorSize,
			MinSignedPerWindow:        oldGenState.Params.MinSignedPerWindow,
			SlashFractionDoubleSign:   oldGenState.Params.SlashFractionDoubleSign,
			SlashFractionDowntime:     oldGenState.Params.SlashFractionDowntime,
			SlashFractionGuardFailure: oldGenState.Params.SlashFractionGuardFailure,
			FallbackGuardReward:       oldGenState.Params.FallbackGuardReward,
			SyncerReward:              DefaultSyncerReward,
		},
		Penalties:    penalties,
		PenaltyNonce: oldGenState.PenaltyNonce,
	}
}
