package keeper

import (
	"fmt"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/sync/types"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func (keeper Keeper) ApplyChange(ctx sdk.Context, change types.Change) bool {
	var applied bool
	var err error
	switch change.Type {
	case types.ConfirmParamProposal:
		applied, err = keeper.ConfirmParamProposal(ctx, change)
	case types.UpdateSidechainAddr:
		applied, err = keeper.UpdateSidechainAddr(ctx, change)
	case types.SyncDelegator:
		applied, err = keeper.SyncDelegator(ctx, change)
	case types.SyncValidator:
		applied, err = keeper.SyncValidator(ctx, change)
	case types.Subscribe:
		applied, err = keeper.Subscribe(ctx, change)
	case types.InitGuardRequest:
		applied, err = keeper.InitGuardRequest(ctx, change)
	case types.GuardTrigger:
		applied, err = keeper.GuardTrigger(ctx, change)
	case types.GuardProof:
		applied, err = keeper.GuardProof(ctx, change)
	default:
		return false
	}
	if err != nil {
		log.Errorln("Apply change err:", err)
	}
	return applied
}

func (keeper Keeper) ConfirmParamProposal(ctx sdk.Context, change types.Change) (bool, error) {
	var paramChange common.ParamChange
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &paramChange)

	log.Infoln("Apply confirm param proposal", paramChange)
	switch paramChange.Record.Uint64() {
	case mainchain.MaxValidatorNum:
		ss, ok := keeper.paramsKeeper.GetSubspace(staking.DefaultParamspace)
		if !ok {
			return false, fmt.Errorf("Fail to get staking subspace")
		}

		err := ss.Update(ctx, staking.KeyMaxValidators, []byte(paramChange.NewValue.String()))
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (keeper Keeper) UpdateSidechainAddr(ctx sdk.Context, change types.Change) (bool, error) {
	var c validator.Candidate
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &c)

	log.Infoln("Apply new candidate", c)
	candidate, found := keeper.validatorKeeper.GetCandidate(ctx, c.EthAddress)
	if found {
		candidate.ValAccount = c.ValAccount
	} else {
		candidate = c
	}
	keeper.validatorKeeper.SetCandidate(ctx, candidate)
	keeper.validatorKeeper.InitAccount(ctx, c.ValAccount)

	return true, nil
}

func (keeper Keeper) SyncDelegator(ctx sdk.Context, change types.Change) (bool, error) {
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

	return true, nil
}

func (keeper Keeper) SyncValidator(ctx sdk.Context, change types.Change) (bool, error) {
	var v staking.Validator
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &v)

	candidate, found := keeper.validatorKeeper.GetCandidate(ctx, v.Description.Identity)
	if !found {
		return false, fmt.Errorf("Fail to get candidate for: %s", v.Description.Identity)
	}
	valAddress := sdk.ValAddress(candidate.ValAccount)

	log.Infof("Apply sync validator %s ethaddr %x status %s token %s commission rate %s",
		candidate.ValAccount, mainchain.Hex2Addr(v.Description.Identity), v.Status, v.Tokens, v.Commission.CommissionRates.Rate)

	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		if v.Status == sdk.Bonded {
			if !sdk.ValAddress(change.Initiator).Equals(valAddress) {
				return false, fmt.Errorf("Bonded validator %s not found, msg sender: %s", candidate.ValAccount, change.Initiator)
			}

			validator = staking.NewValidator(valAddress, v.ConsPubKey, v.Description)
			keeper.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
		} else if v.Status == sdk.Unbonding {
			return false, fmt.Errorf("Unbonding validator %s not found, msg sender: %s", candidate.ValAccount, change.Initiator)
		} else {
			log.Infof("Validator %s already unbonded", candidate.ValAccount)
			return false, nil
		}
	}

	keeper.stakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Commission = v.Commission
	validator.Status = v.Status
	if validator.Status == sdk.Unbonded {
		validator.Tokens = sdk.ZeroInt()
	} else {
		validator.Tokens = v.Tokens
	}
	validator.DelegatorShares = v.Tokens.ToDec()
	keeper.stakingKeeper.SetValidator(ctx, validator)

	if validator.Status == sdk.Bonded {
		keeper.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
	} else if validator.Status == sdk.Unbonded {
		log.Infof("remove validator %s %s %x", valAddress, candidate.ValAccount, mainchain.Hex2Addr(v.Description.Identity))
		keeper.stakingKeeper.RemoveValidator(ctx, valAddress)
	}

	candidate.CommissionRate = v.Commission.CommissionRates.Rate
	keeper.validatorKeeper.SetCandidate(ctx, candidate)

	return true, nil
}

func (keeper Keeper) Subscribe(ctx sdk.Context, change types.Change) (bool, error) {
	var s guard.Subscription
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &s)

	log.Infoln("Apply new subscription", s)
	subscription, found := keeper.guardKeeper.GetSubscription(ctx, s.EthAddress)
	if !found {
		subscription = guard.NewSubscription(s.EthAddress)
	}
	subscription.Deposit = s.Deposit
	keeper.guardKeeper.SetSubscription(ctx, subscription)

	return true, nil
}

func (keeper Keeper) InitGuardRequest(ctx sdk.Context, change types.Change) (bool, error) {
	var r guard.InitRequest
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &r)

	_, simplexChannel, err := common.UnmarshalSignedSimplexStateBytes(r.SignedSimplexStateBytes)
	if err != nil {
		return false, fmt.Errorf("unmarshal signedSimplexStateBytes err: %w", err)
	}
	simplexReceiver, err := eth.RecoverSigner(r.SignedSimplexStateBytes, r.SimplexReceiverSig)
	if err != nil {
		return false, fmt.Errorf("recover signer err: %w", err)
	}

	log.Infof("Apply init request %s, to %x", guard.PrintSimplexChannel(simplexChannel), simplexReceiver)

	_, found := keeper.guardKeeper.GetRequest(ctx, simplexChannel.ChannelId, mainchain.Addr2Hex(simplexReceiver))
	if found {
		return false, fmt.Errorf("guard request already initiated")
	}

	err = keeper.guardKeeper.ChargeRequestFee(ctx, mainchain.Addr2Hex(simplexReceiver))
	if err != nil {
		return false, fmt.Errorf("Fail to charge request fee: %s", err)
	}

	request := guard.NewRequest(
		simplexChannel.ChannelId,
		simplexChannel.SeqNum,
		mainchain.Bytes2Addr(simplexChannel.PeerFrom),
		simplexReceiver,
		r.SignedSimplexStateBytes,
		r.DisputeTimeout)
	keeper.guardKeeper.SetRequest(ctx, request)

	return true, nil
}

func (keeper Keeper) GuardTrigger(ctx sdk.Context, change types.Change) (bool, error) {
	var trigger guard.GuardTrigger
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &trigger)

	log.Infoln("Apply guard trigger", trigger)
	request, found := keeper.guardKeeper.GetRequest(ctx, trigger.ChannelId, trigger.SimplexReceiver)
	if !found {
		return false, fmt.Errorf("Fail to get request with channelId %x %s", trigger.ChannelId, trigger.SimplexReceiver)
	}
	if request.Status != guard.ChanStatus_Idle {
		return false, fmt.Errorf("request channel %x in non-idle status: %s", trigger.ChannelId, request.Status)
	}
	request.TriggerTxHash = trigger.TriggerTxHash
	request.TriggerTxBlkNum = trigger.TriggerTxBlkNum
	if request.SeqNum > trigger.TriggerSeqNum {
		request.AssignedGuards = guard.AssignGuards(ctx, keeper.guardKeeper)
		request.Status = trigger.Status
	}
	keeper.guardKeeper.SetRequest(ctx, request)
	return true, nil
}

func (keeper Keeper) GuardProof(ctx sdk.Context, change types.Change) (bool, error) {
	var proof guard.GuardProof
	keeper.cdc.MustUnmarshalBinaryBare(change.Data, &proof)

	log.Infoln("Apply guard proof", proof)
	request, found := keeper.guardKeeper.GetRequest(ctx, proof.ChannelId, proof.SimplexReceiver)
	if !found {
		return false, fmt.Errorf("Fail to get request with channelId %x %s", proof.ChannelId, proof.SimplexReceiver)
	}

	if request.Status != guard.ChanStatus_Withdrawing && request.Status != guard.ChanStatus_Settling {
		return false, fmt.Errorf("Request not in guard pending state: %d", request.Status)
	}

	request.GuardTxHash = proof.GuardTxHash
	request.GuardTxBlkNum = proof.GuardTxBlkNum
	request.GuardSender = proof.GuardSender
	request.Status = proof.Status

	assignedGuards := request.AssignedGuards
	blockNumberDiff := request.GuardTxBlkNum - request.TriggerTxBlkNum
	guardIndex := (len(assignedGuards) + 1) * int(blockNumberDiff) / int(request.DisputeTimeout)

	var rewardValidator sdk.AccAddress
	if guardIndex < len(assignedGuards) {
		rewardValidator = request.AssignedGuards[guardIndex]
	} else {
		rewardCandidate, found := keeper.validatorKeeper.GetCandidate(ctx, request.GuardSender)
		if found {
			rewardValidator = rewardCandidate.ValAccount
		}

		guardIndex = len(assignedGuards)
	}

	// slash corresponding guards and reward corresponding validator
	for i := 0; i < guardIndex; i++ {
		keeper.slashKeeper.HandleGuardFailure(ctx, rewardValidator, request.AssignedGuards[i])
	}

	request.AssignedGuards = nil
	keeper.guardKeeper.SetRequest(ctx, request)

	return true, nil
}
