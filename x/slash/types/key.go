package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// module name
	ModuleName = "slash"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	AttributeKeyNonce = "nonce"

	ActionPenalty = "penalty"
)

var (
	PenaltyKeyPrefix = []byte{0x11} // Key prefix for Penalty
	PenaltyNonceKey  = []byte{0x12} // Key for Penalty nonce
)

// get penalty key from nonce
func GetPenaltyKey(nonce uint64) []byte {
	return append(PenaltyKeyPrefix, sdk.Uint64ToBigEndian(nonce)...)
}
