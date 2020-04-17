package types

import (
	"encoding/binary"
	"fmt"
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
	ChangesKeyPrefix          = []byte{0x00}
	ActiveChangeQueuePrefix   = []byte{0x01}
	InactiveChangeQueuePrefix = []byte{0x02}
	ChangeIDKey               = []byte{0x03}
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

// ActiveChangeByTimeKey gets the active change queue key by endTime
func ActiveChangeByTimeKey(endTime time.Time) []byte {
	return append(ActiveChangeQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

// ActiveChangeQueueKey returns the key for a changeID in the activeChangeQueue
func ActiveChangeQueueKey(changeID uint64, endTime time.Time) []byte {
	return append(ActiveChangeByTimeKey(endTime), GetChangeIDBytes(changeID)...)
}

// InactiveChangeByTimeKey gets the inactive change queue key by endTime
func InactiveChangeByTimeKey(endTime time.Time) []byte {
	return append(InactiveChangeQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

// InactiveChangeQueueKey returns the key for a changeID in the inactiveChangeQueue
func InactiveChangeQueueKey(changeID uint64, endTime time.Time) []byte {
	return append(InactiveChangeByTimeKey(endTime), GetChangeIDBytes(changeID)...)
}

// Split keys function; used for iterators

// SplitChangeKey split the change key and returns the change id
func SplitChangeKey(key []byte) (changeID uint64) {
	if len(key[1:]) != 8 {
		panic(fmt.Sprintf("unexpected key length (%d ≠ 8)", len(key[1:])))
	}

	return GetChangeIDFromBytes(key[1:])
}

// SplitActiveChangeQueueKey split the active change key and returns the change id and endTime
func SplitActiveChangeQueueKey(key []byte) (changeID uint64, endTime time.Time) {
	return splitKeyWithTime(key)
}

// SplitInactiveChangeQueueKey split the inactive change key and returns the change id and endTime
func SplitInactiveChangeQueueKey(key []byte) (changeID uint64, endTime time.Time) {
	return splitKeyWithTime(key)
}

// private functions

func splitKeyWithTime(key []byte) (changeID uint64, endTime time.Time) {
	if len(key[1:]) != 8+lenTime {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), lenTime+8))
	}

	endTime, err := sdk.ParseTimeBytes(key[1 : 1+lenTime])
	if err != nil {
		panic(err)
	}

	changeID = GetChangeIDFromBytes(key[1+lenTime:])
	return
}

func splitKeyWithAddress(key []byte) (changeID uint64, addr sdk.AccAddress) {
	if len(key[1:]) != 8+sdk.AddrLen {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key), 8+sdk.AddrLen))
	}

	changeID = GetChangeIDFromBytes(key[1:9])
	addr = sdk.AccAddress(key[9:])
	return
}
