package keeper

import (
	"errors"
	"fmt"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync/types"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func (keeper Keeper) ApplyChange(ctx sdk.Context, change types.Change) error {
	switch change.Type {
	case types.Subscribe:
		return keeper.Subscribe(ctx, change)
	case types.Request:
		return keeper.Request(ctx, change)
	case types.IntendSettle:
		return keeper.IntendSettle(ctx, change)
	case types.GuardProof:
		return keeper.GuardProof(ctx, change)
	case types.UpdateSidechainAddr:
		return keeper.UpdateSidechainAddr(ctx, change)
	case types.SyncDelegator:
		return keeper.SyncDelegator(ctx, change)
	case types.SyncValidator:
		return keeper.SyncValidator(ctx, change)
	default:
		return errors.New("Invalid change type")
	}
}

func (keeper Keeper) Subscribe(ctx sdk.Context, change types.Change) error {
	var s subscribe.Subscription
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &s)

	log.Infoln("Apply new subscription", s)
	subscription, found := keeper.subscribeKeeper.GetSubscription(ctx, s.EthAddress)
	if !found {
		subscription = subscribe.NewSubscription(s.EthAddress)
	}
	subscription.Deposit = s.Deposit
	keeper.subscribeKeeper.SetSubscription(ctx, subscription)

	return nil
}

func (keeper Keeper) Request(ctx sdk.Context, change types.Change) error {
	var r subscribe.Request
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &r)

	log.Infoln("Apply new request", r)
	err := keeper.subscribeKeeper.ChargeRequestFee(ctx, r.GetOwnerAddress())
	if err != nil {
		return fmt.Errorf("Failed to charge request fee: %s", err)
	}

	request, found := keeper.subscribeKeeper.GetRequest(ctx, r.ChannelId, r.GetOwnerAddress())
	if found {
		request.SeqNum = r.SeqNum
		request.SignedSimplexStateBytes = r.SignedSimplexStateBytes
	} else {
		request = r
	}

	keeper.subscribeKeeper.SetRequest(ctx, request)

	return nil
}

func (keeper Keeper) IntendSettle(ctx sdk.Context, change types.Change) error {
	var r subscribe.Request
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &r)

	log.Infoln("Apply intend settle", r)
	request, found := keeper.subscribeKeeper.GetRequest(ctx, r.ChannelId, r.GetOwnerAddress())
	if !found {
		return fmt.Errorf("failed to get request with channelId %x and owner %s", r.ChannelId, r.GetOwnerAddress())
	}

	request.TriggerTxHash = r.TriggerTxHash
	request.TriggerTxBlkNum = r.TriggerTxBlkNum
	request.DisputeTimeout = r.DisputeTimeout
	request.RequestGuards = subscribe.GetRequestGuards(ctx, keeper.subscribeKeeper)
	keeper.subscribeKeeper.SetRequest(ctx, request)

	return nil
}

func (keeper Keeper) GuardProof(ctx sdk.Context, change types.Change) error {
	var r subscribe.Request
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &r)

	log.Infoln("Apply guard proof", r)
	request, found := keeper.subscribeKeeper.GetRequest(ctx, r.ChannelId, r.GetOwnerAddress())
	if !found {
		return fmt.Errorf("failed to get request with channelId %x and owner %s", r.ChannelId, r.GetOwnerAddress())
	}

	request.GuardTxHash = r.GuardTxHash
	request.GuardTxBlkNum = r.GuardTxBlkNum
	keeper.subscribeKeeper.SetRequest(ctx, request)

	requestGuards := request.RequestGuards
	blockNumberDiff := request.GuardTxBlkNum - request.TriggerTxBlkNum
	guardIndex := (len(requestGuards) + 1) * int(blockNumberDiff) / int(request.DisputeTimeout)

	var rewardValidator sdk.AccAddress
	if guardIndex < len(requestGuards) {
		rewardValidator = request.RequestGuards[guardIndex]
	} else {
		rewardCandidate, found := keeper.validatorKeeper.GetCandidate(ctx, request.GuardSender)
		if found {
			rewardValidator = rewardCandidate.Operator
		}

		guardIndex = len(requestGuards)
	}

	// punish corresponding guards and reward corresponding validator
	for i := 0; i < guardIndex; i++ {
		keeper.slashKeeper.HandleGuardFailure(ctx, rewardValidator, request.RequestGuards[i])
	}

	return nil
}

func (keeper Keeper) UpdateSidechainAddr(ctx sdk.Context, change types.Change) error {
	var c validator.Candidate
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &c)

	log.Infoln("Apply new candidate", c)
	candidate, found := keeper.validatorKeeper.GetCandidate(ctx, c.EthAddress)
	if found {
		candidate.Operator = c.Operator
	} else {
		candidate = c
	}
	keeper.validatorKeeper.SetCandidate(ctx, candidate)
	keeper.validatorKeeper.InitAccount(ctx, c.Operator)

	return nil
}

func (keeper Keeper) SyncDelegator(ctx sdk.Context, change types.Change) error {
	var d validator.Delegator
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &d)

	log.Infoln("Apply sync delegator", d)
	delegator, found := keeper.validatorKeeper.GetDelegator(ctx, d.CandidateAddr, d.DelegatorAddr)
	if found {
		delegator.DelegatedStake = d.DelegatedStake
	} else {
		delegator = d
	}
	keeper.validatorKeeper.SetDelegator(ctx, delegator)
	keeper.validatorKeeper.SnapshotCandidate(ctx, d.CandidateAddr)

	return nil
}

func (keeper Keeper) SyncValidator(ctx sdk.Context, change types.Change) error {
	var v staking.Validator
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &v)

	log.Infoln("Apply sync validator", v)
	candidate, found := keeper.validatorKeeper.GetCandidate(ctx, v.Description.Identity)
	if !found {
		return fmt.Errorf("failed to get candidate for: %s", v.Description.Identity)
	}

	valAddress := sdk.ValAddress(candidate.Operator)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		if !sdk.ValAddress(change.Initiator).Equals(valAddress) {
			return fmt.Errorf("Invalid change iniator %x for validator %x", change.Initiator, valAddress)
		}

		validator = staking.NewValidator(valAddress, v.ConsPubKey, v.Description)
		keeper.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	}

	keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Tokens = v.Tokens
	validator.Status = v.Status
	validator.DelegatorShares = v.Tokens.ToDec()
	keeper.stakingKeeper.SetValidator(ctx, validator)

	if validator.Status == sdk.Bonded {
		keeper.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
	}

	return nil
}
