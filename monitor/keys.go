package monitor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	PullerKeyPrefix  = []byte{0x01} // Key prefix for puller
	GuardKeyPrefix   = []byte{0x02} // Key prefix for guard
	PenaltyKeyPrefix = []byte{0x03} // Key prefix for penalty
)

// get puller key from log
func GetPullerKey(log types.Log) []byte {
	return append(PullerKeyPrefix, log.TxHash.Bytes()...)
}

// get pusher key from log
func GetGuardKey(log types.Log) []byte {
	return append(GuardKeyPrefix, log.TxHash.Bytes()...)
}

// get penalty key from nonce
func GetPenaltyKey(nonce uint64) []byte {
	return append(PenaltyKeyPrefix, sdk.Uint64ToBigEndian(nonce)...)
}
