package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// module name
	ModuleName = "slash"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	PenaltyKeyPrefix = []byte{0x01} // Key prefix for Penalty
)

// get penalty key from nounce
func GetPenaltyKey(nounce uint64) []byte {
	return append(PenaltyKeyPrefix, sdk.Uint64ToBigEndian(nounce)...)
}
