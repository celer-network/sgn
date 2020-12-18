package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "sync"

	// StoreKey is the store key string for sync
	StoreKey = ModuleName

	// RouterKey is the message route for sync
	RouterKey = ModuleName

	// QuerierRoute is the querier route for sync
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName
)

// Keys for sync store
// Items are stored with the following key: values
//
// - 0x00<changeID_Bytes>: Change
//
// - 0x01<endTime_Bytes><changeID_Bytes>: activeChangeID
//
// - 0x02<endTime_Bytes><changeID_Bytes>: inactiveChangeID
//
// - 0x03: nextChangeID

var (
	ChangesKeyPrefix        = []byte{0x00}
	ActiveChangeQueuePrefix = []byte{0x01}
	ChangeIDKey             = []byte{0x03}
	BlkNumKey               = []byte{0x04}
)

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// GetChangeIDBytes returns the byte representation of the changeID
func GetChangeIDBytes(changeID uint64) (changeIDBz []byte) {
	changeIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(changeIDBz, changeID)
	return
}

// GetChangeIDFromBytes returns changeID in uint64 format from a byte array
func GetChangeIDFromBytes(bz []byte) (changeID uint64) {
	return binary.BigEndian.Uint64(bz)
}

// ChangeKey gets a specific change from the store
func ChangeKey(changeID uint64) []byte {
	return append(ChangesKeyPrefix, GetChangeIDBytes(changeID)...)
}

// ActiveChangeQueueKey returns the key for a changeID in the activeChangeQueue
func ActiveChangeQueueKey(changeID uint64) []byte {
	return append(ActiveChangeQueuePrefix, GetChangeIDBytes(changeID)...)
}
