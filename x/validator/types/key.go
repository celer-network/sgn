package types

import (
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	PullerKey          = []byte{0x01} // key for puller
	PusherKey          = []byte{0x02} // key for pusher
	DelegatorKeyPrefix = []byte{0x03} // Key prefix for delegator

	CandidateKeyPrefix           = []byte{0x40} // Key prefix for candidate
	CandidateByOperatorKeyPrefix = []byte{0x41} // Key prefix for candidate by operator

	RewardKeyPrefix = []byte{0x05} // Key prefix for reward
)

// get delegators key from candidate address
func GetDelegatorsKey(candidateAddr string) []byte {
	return append(DelegatorKeyPrefix, []byte(mainchain.FormatAddrHex(candidateAddr))...)
}

// get delegator key from candidate address and delegator address
func GetDelegatorKey(candidateAddr, delegatorAddr string) []byte {
	return append(GetDelegatorsKey(candidateAddr), []byte(mainchain.FormatAddrHex(delegatorAddr))...)
}

// get candidate key from candidateAddr
func GetCandidateKey(candidateAddr string) []byte {
	return append(CandidateKeyPrefix, []byte(mainchain.FormatAddrHex(candidateAddr))...)
}

// get candidate by operator key from operator
func GetCandidateByOperatorKey(operator sdk.AccAddress) []byte {
	return append(CandidateKeyPrefix, operator.Bytes()...)
}

// get reward key from ethAddr
func GetRewardKey(ethAddr string) []byte {
	return append(RewardKeyPrefix, []byte(mainchain.FormatAddrHex(ethAddr))...)
}
