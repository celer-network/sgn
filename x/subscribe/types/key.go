package types

import (
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// module name
	ModuleName = "subscribe"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	SubscriptionKeyPrefix = []byte{0x01} // Key prefix for subscription
	RequestKeyPrefix      = []byte{0x21} // Key prefix for request
	EpochKeyPrefix        = []byte{0x31} // Key prefix for epoch
	LatestEpochKey        = []byte{0x32} // Key for latest epoch
)

// get guardian key from eth address
func GetSubscriptionKey(ethAddress string) []byte {
	return append(SubscriptionKeyPrefix, []byte(ethAddress)...)
}

// get request key from channelID
func GetRequestKey(channelId []byte, receiver string) []byte {
	return append(append(RequestKeyPrefix, mainchain.Bytes2Cid(channelId).Bytes()...), []byte(mainchain.FormatAddrHex(receiver))...)
}

// get epoch key from epochId
func GetEpochKey(epochId sdk.Int) []byte {
	return append(EpochKeyPrefix, epochId.BigInt().Bytes()...)
}

// get latest epoch key
func GetLatestEpochKey() []byte {
	return LatestEpochKey
}
