package validator

import (
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

// Gets the entire Puller metadata
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

// Gets the entire Pusher metadata
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

// Gets the entire Delegator metadata for a candidateAddress and delegatorAddress
func (k Keeper) GetDelegator(ctx sdk.Context, candidateAddress, delegatorAddress string) Delegator {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(GetDelegatorKey(candidateAddress, delegatorAddress)) {
		return Delegator{}
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

// Gets the entire Candidate metadata for a candidateAddress and seq
func (k Keeper) GetCandidate(ctx sdk.Context, candidateAddress string, seq sdk.Int) (candidate Candidate, found bool) {
	store := ctx.KVStore(k.storeKey)
	candidateKey := GetCandidateKey(candidateAddress, seq)

	if !store.Has(candidateKey) {
		return candidate, false
	}

	value := store.Get(candidateKey)
	k.cdc.MustUnmarshalBinaryBare(value, &candidate)
	return candidate, true
}

// Sets the entire Candidate metadata
func (k Keeper) SetCandidate(ctx sdk.Context, candidate Candidate) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetCandidateKey(candidate.EthAddress, candidate.Seq), k.cdc.MustMarshalBinaryBare(candidate))
}

// Gets the entire latest Candidate metadata
func (k Keeper) GetLatestCandidate(ctx sdk.Context, candidateAddress string) (candidate Candidate) {
	store := ctx.KVStore(k.storeKey)
	latestCandidateKey := GetLatestCandidateKey(candidateAddress)

	if !store.Has(latestCandidateKey) {
		return
	}

	value := store.Get(latestCandidateKey)
	k.cdc.MustUnmarshalBinaryBare(value, &candidate)
	return
}

// Sets the entire LatestCandidate metadata
func (k Keeper) SetLatestCandidate(ctx sdk.Context, candidate Candidate) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetLatestCandidateKey(candidate.EthAddress), k.cdc.MustMarshalBinaryBare(candidate))
	k.SetCandidate(ctx, candidate)
}

// Take a snapshot of candidate
func (k Keeper) SnapshotCandidate(ctx sdk.Context, candidateAddr string) {
	latestCandidate := k.GetLatestCandidate(ctx, candidateAddr)
	nextSeq := latestCandidate.Seq.AddRaw(1)
	newCandidate := NewCandidate(candidateAddr, nextSeq)
	newCandidate.Delegators = k.GetAllDelegators(ctx, candidateAddr)
	k.SetLatestCandidate(ctx, newCandidate)
}
