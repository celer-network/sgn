package validator

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey      sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc           *codec.Codec // The wire codec for binary encoding/decoding.
	ethClient     *mainchain.EthClient
	globalKeeper  global.Keeper
	accountKeeper auth.AccountKeeper
	stakingKeeper staking.Keeper
}

// NewKeeper creates new instances of the validator Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, ethClient *mainchain.EthClient,
	globalKeeper global.Keeper, accountKeeper auth.AccountKeeper, stakingKeeper staking.Keeper) Keeper {
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		ethClient:     ethClient,
		globalKeeper:  globalKeeper,
		accountKeeper: accountKeeper,
		stakingKeeper: stakingKeeper,
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

// Get the entire Puller metadata
func (k Keeper) GetPuller(ctx sdk.Context) Puller {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(PullerKey) {
		return Puller{}
	}

	value := store.Get(PullerKey)
	var puller Puller
	k.cdc.MustUnmarshalBinaryBare(value, &puller)
	return puller
}

// Sets the entire Puller metadata
func (k Keeper) SetPuller(ctx sdk.Context, puller Puller) {
	store := ctx.KVStore(k.storeKey)
	store.Set(PullerKey, k.cdc.MustMarshalBinaryBare(puller))
}

// Get the entire Pusher metadata
func (k Keeper) GetPusher(ctx sdk.Context) Pusher {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(PusherKey) {
		return Pusher{}
	}

	value := store.Get(PusherKey)
	var pusher Pusher
	k.cdc.MustUnmarshalBinaryBare(value, &pusher)
	return pusher
}

// Sets the entire Pusher metadata
func (k Keeper) SetPusher(ctx sdk.Context, pusher Pusher) {
	store := ctx.KVStore(k.storeKey)
	store.Set(PusherKey, k.cdc.MustMarshalBinaryBare(pusher))
}

// Get the entire Delegator metadata for a candidateAddress and delegatorAddress
func (k Keeper) GetDelegator(ctx sdk.Context, candidateAddress, delegatorAddress string) Delegator {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(GetDelegatorKey(candidateAddress, delegatorAddress)) {
		return NewDelegator(delegatorAddress)
	}

	var delegator Delegator
	value := store.Get(GetDelegatorKey(candidateAddress, delegatorAddress))
	k.cdc.MustUnmarshalBinaryBare(value, &delegator)
	return delegator
}

// Get the set of all delegators with no limits
func (k Keeper) GetAllDelegators(ctx sdk.Context, candidateAddress string) (delegators []Delegator) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, GetDelegatorsKey(candidateAddress))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var delegator Delegator
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &delegator)
		delegators = append(delegators, delegator)
	}
	return delegators
}

// Sets the entire Delegator metadata for a candidateAddress and delegatorAddress
func (k Keeper) SetDelegator(ctx sdk.Context, candidateAddress, delegatorAddress string, delegator Delegator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetDelegatorKey(candidateAddress, delegatorAddress), k.cdc.MustMarshalBinaryBare(delegator))
}

// Get the entire Candidate metadata
func (k Keeper) GetCandidate(ctx sdk.Context, candidateAddress string) (candidate Candidate, found bool) {
	store := ctx.KVStore(k.storeKey)
	candidateKey := GetCandidateKey(candidateAddress)

	if !store.Has(candidateKey) {
		return candidate, false
	}

	value := store.Get(candidateKey)
	k.cdc.MustUnmarshalBinaryBare(value, &candidate)
	return candidate, true
}

// Get the set of all candidates with no limits
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
func (k Keeper) SetCandidate(ctx sdk.Context, candidateAddr string, candidate Candidate) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetCandidateKey(candidateAddr), k.cdc.MustMarshalBinaryBare(candidate))
}

// Take a snapshot of candidate
func (k Keeper) SnapshotCandidate(ctx sdk.Context, candidateAddr string) {
	candidate, _ := k.GetCandidate(ctx, candidateAddr)
	candidate.Delegators = k.GetAllDelegators(ctx, candidateAddr)

	totalStake := sdk.ZeroInt()
	for _, delegator := range candidate.Delegators {
		totalStake = totalStake.Add(delegator.DelegatedStake)
	}
	candidate.StakingPool = totalStake

	k.SetCandidate(ctx, candidateAddr, candidate)
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

// DistributeServiceReward distributes service rewards to candidates and their delegators
func (k Keeper) DistributeServiceReward(ctx sdk.Context, totalReward sdk.Int) {
	candidates := k.GetValidatorCandidates(ctx)
	totalStake := GetCandidatesTotalStake(candidates)

	for _, candidate := range candidates {
		candidateReward := totalReward.Mul(candidate.StakingPool).Quo(totalStake)
		for _, delegator := range candidate.Delegators {
			reward, _ := k.GetReward(ctx, delegator.EthAddress)
			rewardAmt := candidateReward.Mul(delegator.DelegatedStake).Quo(candidate.StakingPool)
			reward.ServiceReward = reward.ServiceReward.Add(rewardAmt)
			k.SetReward(ctx, reward)
		}
	}
}

// DistributeMiningReward distributes mining rewards to candidates and their delegators
func (k Keeper) DistributeMiningReward(ctx sdk.Context, totalReward sdk.Int) {
	candidates := k.GetValidatorCandidates(ctx)
	totalStake := GetCandidatesTotalStake(candidates)

	for _, candidate := range candidates {
		candidateReward := totalReward.Mul(candidate.StakingPool).Quo(totalStake)
		for _, delegator := range candidate.Delegators {
			reward, _ := k.GetReward(ctx, delegator.EthAddress)
			rewardAmt := candidateReward.Mul(delegator.DelegatedStake).Quo(candidate.StakingPool)
			reward.MiningReward = reward.MiningReward.Add(rewardAmt)
			k.SetReward(ctx, reward)
		}
	}
}

// GetValidatorCandidates get candidates info for current validators
func (k Keeper) GetValidatorCandidates(ctx sdk.Context) (candidates []Candidate) {
	validators := k.GetValidators(ctx)

	for _, validator := range validators {
		ethAddr := mainchain.FormatAddrHex(validator.Description.Identity)
		if ethAddr == "" {
			log.Errorf("Miss eth address for validator %x", validator.OperatorAddress)
			continue
		}
		candidate, found := k.GetCandidate(ctx, ethAddr)

		if found {
			candidates = append(candidates, candidate)
		}
	}

	return
}
