package types

import (
	"github.com/celer-network/sgn/mainchain"
)

const (
	// module name
	ModuleName = "guard"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	SubscriptionKeyPrefix = []byte{0x01} // Key prefix for subscription
	RequestKeyPrefix      = []byte{0x21} // Key prefix for request
)

// get guardian key from eth address
func GetSubscriptionKey(ethAddress string) []byte {
	return append(SubscriptionKeyPrefix, []byte(mainchain.FormatAddrHex(ethAddress))...)
}

// get request key from channelID
func GetRequestKey(channelId []byte, simplexReceiver string) []byte {
	return append(append(RequestKeyPrefix, mainchain.Bytes2Cid(channelId).Bytes()...), []byte(mainchain.FormatAddrHex(simplexReceiver))...)
}
