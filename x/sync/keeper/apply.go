package keeper

import (
	"errors"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync/types"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (keeper Keeper) ApplyChange(ctx sdk.Context, change types.Change) error {
	switch change.Type {
	case types.SyncBlock:
		return keeper.SyncBlock(ctx, change.Data)
	case types.Subscribe:
		return keeper.Subscribe(ctx, change.Data)
	case types.UpdateSidechainAddr:
		return keeper.UpdateSidechainAddr(ctx, change.Data)
	case types.SyncDelegator:
		return keeper.SyncDelegator(ctx, change.Data)
	default:
		return errors.New("Invalid change type")
	}
}

func (keeper Keeper) SyncBlock(ctx sdk.Context, data []byte) error {
	var block global.Block
	keeper.cdc.MustUnmarshalBinaryBare(data, &block)

	keeper.gk.SyncBlock(ctx, block.Number)
	log.Infoln("Sync mainchain block to", block.Number)

	return nil
}

func (keeper Keeper) Subscribe(ctx sdk.Context, data []byte) error {
	var s subscribe.Subscription
	keeper.cdc.MustUnmarshalBinaryBare(data, &s)

	log.Infoln("Apply new subscription", s)
	subscription, found := keeper.sk.GetSubscription(ctx, s.EthAddress)
	if !found {
		subscription = subscribe.NewSubscription(s.EthAddress)
	}
	subscription.Deposit = s.Deposit
	keeper.sk.SetSubscription(ctx, subscription)

	return nil
}

func (keeper Keeper) UpdateSidechainAddr(ctx sdk.Context, data []byte) error {
	var c validator.Candidate
	keeper.cdc.MustUnmarshalBinaryBare(data, &c)

	candidate, found := keeper.vk.GetCandidate(ctx, c.EthAddress)
	if found {
		candidate.Operator = c.Operator
	} else {
		candidate = c
	}
	keeper.vk.SetCandidate(ctx, candidate)
	keeper.vk.InitAccount(ctx, c.Operator)

	return nil
}

func (keeper Keeper) SyncDelegator(ctx sdk.Context, data []byte) error {
	var d validator.Delegator
	keeper.cdc.MustUnmarshalBinaryBare(data, &d)

	delegator, found := keeper.vk.GetDelegator(ctx, d.CandidateAddr, d.DelegatorAddr)
	if found {
		delegator.DelegatedStake = d.DelegatedStake
	} else {
		delegator = d
	}
	keeper.vk.SetDelegator(ctx, delegator)
	keeper.vk.SnapshotCandidate(ctx, d.CandidateAddr)

	return nil
}
