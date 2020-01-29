package monitor

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	EventKeyPrefix = []byte{0x01} // Key prefix for event

	PullerKeyPrefix = []byte{0x02} // Key prefix for puller

	PusherKeyPrefix = []byte{0x03} // Key prefix for pusher

	PenaltyKeyPrefix = []byte{0x04} // Key prefix for penalty
)

// get event key from log
func GetEventKey(log types.Log) []byte {
	logIndexBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(logIndexBytes, uint64(log.Index))
	return append(EventKeyPrefix, append(log.TxHash.Bytes(), logIndexBytes[:]...)...)
}

// get puller key from log
func GetPullerKey(log types.Log) []byte {
	return append(PullerKeyPrefix, log.TxHash.Bytes()...)
}

// get pusher key from log
func GetPusherKey(log types.Log) []byte {
	return append(PusherKeyPrefix, log.TxHash.Bytes()...)
}

// get penalty key from nonce
func GetPenaltyKey(nonce uint64) []byte {
	return append(PenaltyKeyPrefix, sdk.Uint64ToBigEndian(nonce)...)
}
