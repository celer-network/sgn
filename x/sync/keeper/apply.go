package keeper

import (
	"errors"
	"fmt"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync/types"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func (keeper Keeper) ApplyChange(ctx sdk.Context, change types.Change) error {
	switch change.Type {
	case types.ConfirmParamProposal:
		return keeper.ConfirmParamProposal(ctx, change)
	case types.UpdateSidechainAddr:
		return keeper.UpdateSidechainAddr(ctx, change)
	case types.SyncDelegator:
		return keeper.SyncDelegator(ctx, change)
	case types.SyncValidator:
		return keeper.SyncValidator(ctx, change)
	case types.Subscribe:
		return keeper.Subscribe(ctx, change)
	case types.InitGuardRequest:
		return keeper.InitGuardRequest(ctx, change)
	case types.TriggerGuard:
		return keeper.TriggerGuard(ctx, change)
	case types.GuardProof:
		return keeper.GuardProof(ctx, change)
	default:
		return errors.New("Invalid change type")
	}
}

func (keeper Keeper) ConfirmParamProposal(ctx sdk.Context, change types.Change) error {
	var paramChange common.ParamChange
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &paramChange)

	log.Infoln("Apply confirm param proposal", paramChange)
	switch paramChange.Record.Uint64() {
	case mainchain.MaxValidatorNum:
		ss, ok := keeper.paramsKeeper.GetSubspace(staking.DefaultParamspace)
		if !ok {
			return fmt.Errorf("Fail to get staking subspace")
		}

		err := ss.Update(ctx, staking.KeyMaxValidators, keeper.cdc.MustMarshalBinaryBare(uint16(paramChange.NewValue.Uint64())))
		return err
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

	candidate, found := keeper.validatorKeeper.GetCandidate(ctx, v.Description.Identity)
	if !found {
		return fmt.Errorf("Fail to get candidate for: %s", v.Description.Identity)
	}
	valAddress := sdk.ValAddress(candidate.Operator)

	log.Infof("apply sync validator %s ethaddr %x status %s token %s commission %s",
		candidate.Operator.String(),
		mainchain.Hex2Addr(v.Description.Identity),
		v.Status, v.Tokens, v.Commission)

	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		if !sdk.ValAddress(change.Initiator).Equals(valAddress) {
			return fmt.Errorf("Invalid change initiator %x for validator %x", change.Initiator, valAddress)
		}

		validator = staking.NewValidator(valAddress, v.ConsPubKey, v.Description)
		keeper.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	}

	keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Commission = v.Commission
	validator.Tokens = v.Tokens
	validator.Status = v.Status
	validator.DelegatorShares = v.Tokens.ToDec()
	keeper.stakingKeeper.SetValidator(ctx, validator)

	if validator.Status == sdk.Bonded {
		keeper.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
	}

	if validator.Status == sdk.Unbonded {
		validator.Tokens = sdk.ZeroInt()
		keeper.stakingKeeper.RemoveValidator(ctx, valAddress)
	}

	candidate.CommissionRate = v.Commission.CommissionRates.Rate
	keeper.validatorKeeper.SetCandidate(ctx, candidate)

	return nil
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

func (keeper Keeper) InitGuardRequest(ctx sdk.Context, change types.Change) error {
	var r subscribe.Request
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &r)

	log.Infoln("Apply init request", r)
	err := keeper.subscribeKeeper.ChargeRequestFee(ctx, r.GetPeerToAddress())
	if err != nil {
		return fmt.Errorf("Fail to charge request fee: %s", err)
	}

	_, found := keeper.subscribeKeeper.GetRequest(ctx, r.ChannelId, r.GetPeerFromAddress())
	if found {
		return fmt.Errorf("guard request already initiated")
	}

	keeper.subscribeKeeper.SetRequest(ctx, r)

	return nil
}

func (keeper Keeper) TriggerGuard(ctx sdk.Context, change types.Change) error {
	var r subscribe.Request
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &r)

	log.Infoln("Apply intend settle", r)
	request, found := keeper.subscribeKeeper.GetRequest(ctx, r.ChannelId, r.GetPeerFromAddress())
	if !found {
		return fmt.Errorf("Fail to get request with channelId %x %s", r.ChannelId, r.GetPeerFromAddress())
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
	request, found := keeper.subscribeKeeper.GetRequest(ctx, r.ChannelId, r.GetPeerFromAddress())
	if !found {
		return fmt.Errorf("Fail to get request with channelId %x %s", r.ChannelId, r.GetPeerFromAddress())
	}

	request.GuardTxHash = r.GuardTxHash
	request.GuardTxBlkNum = r.GuardTxBlkNum
	request.GuardSender = r.GuardSender
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
