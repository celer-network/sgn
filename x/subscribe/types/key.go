package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// module name
	ModuleName = "subscribe"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	SubscriptionKey = []byte{0x01} // Prefix for subscription
	RequestKey      = []byte{0x02} // Prefix for request
	EpochKey        = []byte{0x03} // Prefix for epoch
)

// get guardian key from eth address
func GetSubscriptionKey(ethAddress string) []byte {
	return append(SubscriptionKey, []byte(ethAddress)...)
}

// get request key from channelID
func GetRequestKey(channelId []byte) []byte {
	return append(RequestKey, channelId...)
}

// get epoch key from epochId
func GetEpochKey(epochId sdk.Int) []byte {
	return append(RequestKey, epochId.BigInt().Bytes()...)
}

// get latest epoch key
func GetLatestEpochKey() []byte {
	return RequestKey
}
