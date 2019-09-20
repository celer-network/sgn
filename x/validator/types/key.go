package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// module name
	ModuleName = "validator"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	PullerKey = []byte{0x01} // key for puller
	PusherKey = []byte{0x02} // key for pusher

	DelegatorKeyPrefix = []byte{0x03} // Key prefix for delegator

	CandidateKeyPrefix         = []byte{0x41} // Key prefix for candidate
	CandidateSnapshotKeyPrefix = []byte{0x42} // Key prefix for candidate snapshot
)

// get delegators key from candidate address
func GetDelegatorsKey(candidateAddr string) []byte {
	return append(DelegatorKeyPrefix, []byte(candidateAddr)...)
}

// get delegator key from candidate address and delegator address
func GetDelegatorKey(candidateAddr, delegatorAddr string) []byte {
	return append(GetDelegatorsKey(candidateAddr), []byte(delegatorAddr)...)
}

// get candidate key from candidateAddr
func GetCandidateKey(candidateAddr string) []byte {
	return append(CandidateKeyPrefix, []byte(candidateAddr)...)
}

// get candidate snapshot key for candidate with candidateAddr and seq
func GetCandidateSnapshotsKey(candidateAddr string) []byte {
	return append(CandidateSnapshotKeyPrefix, []byte(candidateAddr)...)
}

// get candidate snapshot key for candidate with candidateAddr and seq
func GetCandidateSnapshotKey(candidateAddr string, seq sdk.Int) []byte {
	return append(GetCandidateSnapshotsKey(candidateAddr), seq.BigInt().Bytes()...)
}
