package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// module name
	ModuleName = "global"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	LatestBlockKey = []byte{0x01} // Key for lastest block

	EpochKeyPrefix = []byte{0x21} // Key prefix for epoch
	LatestEpochKey = []byte{0x22} // Key for latest epoch
)

// get epoch key from epochId
func GetEpochKey(epochId sdk.Int) []byte {
	return append(EpochKeyPrefix, epochId.BigInt().Bytes()...)
}

// get latest epoch key
func GetLatestEpochKey() []byte {
	return LatestEpochKey
}
