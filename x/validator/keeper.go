package validator

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey      sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc           *codec.Codec // The wire codec for binary encoding/decoding.
	accountKeeper auth.AccountKeeper
	stakingKeeper staking.Keeper
	paramstore    params.Subspace
}

// NewKeeper creates new instances of the validator Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec,
	accountKeeper auth.AccountKeeper, stakingKeeper staking.Keeper, paramstore params.Subspace) Keeper {
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		accountKeeper: accountKeeper,
		stakingKeeper: stakingKeeper,
		paramstore:    paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// Get validators metadata
func (k Keeper) GetValidators(ctx sdk.Context) []staking.Validator {
	return k.stakingKeeper.GetBondedValidatorsByPower(ctx)
}

// Get a validator by consencus address
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, addr sdk.ConsAddress) (staking.Validator, bool) {
	return k.stakingKeeper.GetValidatorByConsAddr(ctx, addr)
}

// Get a validator by validator account address
func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (staking.Validator, bool) {
	return k.stakingKeeper.GetValidator(ctx, addr)
}

func (k Keeper) IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, fn)
}

// Get the entire Syncer metadata
func (k Keeper) GetSyncer(ctx sdk.Context) Syncer {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(SyncerKey) {
		return Syncer{}
	}

	value := store.Get(SyncerKey)
	var syncer Syncer
	k.cdc.MustUnmarshalBinaryBare(value, &syncer)
	return syncer
}

// Sets the entire Syncer metadata
func (k Keeper) SetSyncer(ctx sdk.Context, syncer Syncer) {
	store := ctx.KVStore(k.storeKey)
	store.Set(SyncerKey, k.cdc.MustMarshalBinaryBare(syncer))
}

// Get the entire Delegator metadata for a candidateAddr and delegatorAddr
func (k Keeper) GetDelegator(ctx sdk.Context, candidateAddr, delegatorAddr string) (delegator Delegator, found bool) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(GetDelegatorKey(candidateAddr, delegatorAddr)) {
		return delegator, false
	}

	value := store.Get(GetDelegatorKey(candidateAddr, delegatorAddr))
	k.cdc.MustUnmarshalBinaryBare(value, &delegator)
	return delegator, true
}

// Get the list of all delegators
func (k Keeper) GetAllDelegators(ctx sdk.Context, candidateAddr string) (delegators []Delegator) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, GetDelegatorsKey(candidateAddr))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var delegator Delegator
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &delegator)
		delegators = append(delegators, delegator)
	}
	return delegators
}

// Sets the entire Delegator metadata for a candidateAddr and delegatorAddr
func (k Keeper) SetDelegator(ctx sdk.Context, delegator Delegator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetDelegatorKey(delegator.CandidateAddr, delegator.DelegatorAddr), k.cdc.MustMarshalBinaryBare(delegator))
}

// Get the entire Candidate metadata
func (k Keeper) GetCandidate(ctx sdk.Context, candidateAddr string) (candidate Candidate, found bool) {
	store := ctx.KVStore(k.storeKey)
	candidateKey := GetCandidateKey(candidateAddr)

	if !store.Has(candidateKey) {
		return candidate, false
	}

	value := store.Get(candidateKey)
	k.cdc.MustUnmarshalBinaryBare(value, &candidate)
	return candidate, true
}

// Get the list of all candidates
func (k Keeper) GetAllCandidates(ctx sdk.Context) (candidates []Candidate) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, CandidateKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var candidate Candidate
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &candidate)
		candidates = append(candidates, candidate)
	}
	return candidates
}

// Sets the Candidate metadata
func (k Keeper) SetCandidate(ctx sdk.Context, candidate Candidate) {
	store := ctx.KVStore(k.storeKey)
	candidateKey := GetCandidateKey(candidate.EthAddress)
	store.Set(candidateKey, k.cdc.MustMarshalBinaryBare(candidate))
}

// Get the entire Reward metadata for ethAddress
func (k Keeper) GetReward(ctx sdk.Context, ethAddress string) (Reward, bool) {
	store := ctx.KVStore(k.storeKey)
	rewardKey := GetRewardKey(ethAddress)

	if !store.Has(rewardKey) {
		return NewReward(ethAddress), false
	}

	var reward Reward
	value := store.Get(rewardKey)
	k.cdc.MustUnmarshalBinaryBare(value, &reward)
	return reward, true
}

// Sets the Reward metadata for ethAddress
func (k Keeper) SetReward(ctx sdk.Context, reward Reward) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetRewardKey(reward.Receiver), k.cdc.MustMarshalBinaryBare(reward))
}

// AddReward add reward to a specific ethAddress
func (k Keeper) AddReward(ctx sdk.Context, ethAddress string, miningReward, serviceReward sdk.Int) {
	reward, _ := k.GetReward(ctx, ethAddress)
	reward.MiningReward = reward.MiningReward.Add(miningReward)
	reward.ServiceReward = reward.ServiceReward.Add(serviceReward)
	k.SetReward(ctx, reward)
}

func (k Keeper) DistributeReward(ctx sdk.Context) {
	epoch := k.GetRewardEpoch(ctx)
	if ctx.BlockHeight()-epoch.StartHeight < int64(k.EpochLength(ctx)) {
		return
	}

	candidates := k.GetValidatorCandidates(ctx)
	totalStake := sdk.ZeroInt()

	for _, candidate := range candidates {
		if !candidate.StakingPool.IsPositive() {
			log.Errorln("invalid candidate staking pool", candidate.EthAddress)
			return
		}
		totalStake = totalStake.Add(candidate.StakingPool)
	}

	for _, candidate := range candidates {
		candidateMiningReward := epoch.MiningReward.Mul(candidate.StakingPool).Quo(totalStake)
		miningCommission := candidate.CommissionRate.MulInt(candidateMiningReward).RoundInt()

		candidateServiceReward := sdk.ZeroInt()
		serviceCommission := sdk.ZeroInt()
		if epoch.ServiceReward.IsPositive() {
			candidateServiceReward = epoch.ServiceReward.Mul(candidate.StakingPool).Quo(totalStake)
			serviceCommission = candidate.CommissionRate.MulInt(candidateServiceReward).RoundInt()
		}

		k.AddReward(ctx, candidate.EthAddress, miningCommission, serviceCommission)

		delegators := k.GetAllDelegators(ctx, candidate.EthAddress)
		delegatorsMiningReward := candidateMiningReward.Sub(miningCommission)
		delegatorsServiceReward := candidateServiceReward.Sub(serviceCommission)
		for _, delegator := range delegators {
			k.AddReward(ctx, delegator.DelegatorAddr,
				delegatorsMiningReward.Mul(delegator.DelegatedStake).Quo(candidate.StakingPool),
				delegatorsServiceReward.Mul(delegator.DelegatedStake).Quo(candidate.StakingPool))
		}
	}

	k.ResetRewardEpoch(ctx)
}

// GetValidatorCandidates get candidates info for current validators
func (k Keeper) GetValidatorCandidates(ctx sdk.Context) (candidates []Candidate) {
	validators := k.GetValidators(ctx)

	for _, validator := range validators {
		ethAddr := mainchain.FormatAddrHex(validator.Description.Identity)
		if ethAddr == "" {
			log.Errorf("Miss eth address for validator %s", validator.OperatorAddress)
			continue
		}
		candidate, found := k.GetCandidate(ctx, ethAddr)

		if found && candidate.StakingPool.IsPositive() {
			candidates = append(candidates, candidate)
		}
	}

	return
}

func (k Keeper) InitAccount(ctx sdk.Context, accAddress sdk.AccAddress) {
	err := sdk.VerifyAddressFormat(accAddress)
	if err != nil {
		log.Errorf("InitAccount %s err: %s", accAddress, err)
		return
	}
	account := k.accountKeeper.GetAccount(ctx, accAddress)
	if account == nil {
		log.Infof("Set new account %s", accAddress)
		account = k.accountKeeper.NewAccountWithAddress(ctx, accAddress)
		k.accountKeeper.SetAccount(ctx, account)
	}
}

func (k Keeper) GetRewardEpoch(ctx sdk.Context) (epoch RewardEpoch) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(GetRewardEpochKey()) {
		epoch = NewRewardEpoch(ctx.BlockHeight())
		k.SetRewardEpoch(ctx, epoch)
		return
	}

	value := store.Get(GetRewardEpochKey())
	k.cdc.MustUnmarshalBinaryBare(value, &epoch)
	return
}

func (k Keeper) SetRewardEpoch(ctx sdk.Context, epoch RewardEpoch) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetRewardEpochKey(), k.cdc.MustMarshalBinaryBare(epoch))
}

func (k Keeper) ResetRewardEpoch(ctx sdk.Context) {
	k.SetRewardEpoch(ctx, NewRewardEpoch(ctx.BlockHeight()))
}

func (k Keeper) AddEpochMiningReward(ctx sdk.Context, amount sdk.Int) {
	epoch := k.GetRewardEpoch(ctx)
	epoch.MiningReward = epoch.MiningReward.Add(amount)
	k.SetRewardEpoch(ctx, epoch)
}

func (k Keeper) AddEpochServiceReward(ctx sdk.Context, amount sdk.Int) {
	epoch := k.GetRewardEpoch(ctx)
	epoch.ServiceReward = epoch.ServiceReward.Add(amount)
	k.SetRewardEpoch(ctx, epoch)
}
