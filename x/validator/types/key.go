package types

import (
	"github.com/celer-network/sgn/mainchain"
)

const (
	// module name
	ModuleName = "validator"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	AttributeKeyEthAddress = "eth_address"

	ActionInitiateWithdraw = "initate_withdraw"
)

var (
	SyncerKey              = []byte{0x01} // key for syncer
	DelegatorsKeyPrefix    = []byte{0x03} // Key prefix for delegator
	CandidateKeyPrefix     = []byte{0x04} // Key prefix for candidate
	RewardKeyPrefix        = []byte{0x05} // Key prefix for reward
	RewardEpochKey         = []byte{0x06} // Key for reward epoch
	PendingRewardKeyPrefix = []byte{0x07} // Key for pending reward
)

// get delegators key from candidate address
func GetDelegatorsKey(candidateAddr string) []byte {
	return append(DelegatorsKeyPrefix, mainchain.Hex2Addr(candidateAddr).Bytes()...)
}

// get delegator key from candidate address and delegator address
func GetDelegatorKey(candidateAddr, delegatorAddr string) []byte {
	return append(GetDelegatorsKey(candidateAddr), mainchain.Hex2Addr(delegatorAddr).Bytes()...)
}

// get candidate key from candidateAddr
func GetCandidateKey(candidateAddr string) []byte {
	return append(CandidateKeyPrefix, mainchain.Hex2Addr(candidateAddr).Bytes()...)
}

// get reward key from ethAddr
func GetRewardKey(ethAddr string) []byte {
	return append(RewardKeyPrefix, mainchain.Hex2Addr(ethAddr).Bytes()...)
}

func GetPendingRewardKey(candidateAddr string) []byte {
	return append(PendingRewardKeyPrefix, mainchain.Hex2Addr(candidateAddr).Bytes()...)
}
