package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// module name
	ModuleName = "validator"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	PullerKey    = []byte{0x01} // Prefix for puller
	PusherKey    = []byte{0x02} // Prefix for pusher
	DelegatorKey = []byte{0x03} // Prefix for delegator
	CandidateKey = []byte{0x04} // Prefix for candidate
)

// get delegators key from candidate address
func GetDelegatorsKey(candidateAddr string) []byte {
	return append(DelegatorKey, []byte(candidateAddr)...)
}

// get delegator key from candidate address and delegator address
func GetDelegatorKey(candidateAddr, delegatorAddr string) []byte {
	return append(GetDelegatorsKey(candidateAddr), []byte(delegatorAddr)...)
}

// get latest candidate key from candidate address
func GetLatestCandidateKey(candidateAddr string) []byte {
	return append(DelegatorKey, []byte(candidateAddr)...)
}

// get candidate key from candidate address and seq
func GetCandidateKey(candidateAddr string, seq sdk.Int) []byte {
	return append(GetLatestCandidateKey(candidateAddr), seq.BigInt().Bytes()...)
}
