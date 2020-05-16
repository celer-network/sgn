package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// module name
	ModuleName = "global"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	EpochKeyPrefix = []byte{0x11} // Key prefix for epoch
	LatestEpochKey = []byte{0x12} // Key for latest epoch
)

// get epoch key from epochId
func GetEpochKey(epochId sdk.Int) []byte {
	return append(EpochKeyPrefix, epochId.BigInt().Bytes()...)
}

// get latest epoch key
func GetLatestEpochKey() []byte {
	return LatestEpochKey
}
