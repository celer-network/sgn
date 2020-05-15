package keeper

import (
	"errors"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (keeper Keeper) ApplyChange(ctx sdk.Context, change types.Change) error {
	switch change.Type {
	case types.SyncBlock:
		return keeper.SyncBlock(ctx, change.Data)
	case types.Subscribe:
		return keeper.Subscribe(ctx, change.Data)
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
	var newSubscription subscribe.Subscription
	keeper.cdc.MustUnmarshalBinaryBare(data, &newSubscription)

	log.Infoln("Apply new subscription", newSubscription)
	subscription, found := keeper.sk.GetSubscription(ctx, newSubscription.EthAddress)
	if !found {
		subscription = subscribe.NewSubscription(newSubscription.EthAddress)
	}
	subscription.Deposit = newSubscription.Deposit
	keeper.sk.SetSubscription(ctx, subscription)

	return nil
}
