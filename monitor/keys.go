package monitor

import (
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	PullerKeyPrefix  = []byte{0x01} // Key prefix for puller
	GuardKeyPrefix   = []byte{0x02} // Key prefix for guard
	PenaltyKeyPrefix = []byte{0x03} // Key prefix for penalty
)

// get puller key from mainchain txHash
func GetPullerKey(txHash mainchain.HashType) []byte {
	return append(PullerKeyPrefix, txHash.Bytes()...)
}

// get pusher key from mainchain txHash
func GetGuardKey(cid mainchain.CidType) []byte {
	return append(GuardKeyPrefix, cid.Bytes()...)
}

// get penalty key from nonce
func GetPenaltyKey(nonce uint64) []byte {
	return append(PenaltyKeyPrefix, sdk.Uint64ToBigEndian(nonce)...)
}
