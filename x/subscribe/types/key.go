package types

const (
	// module name
	ModuleName = "subscribe"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	SubscriptionKeyPrefix = []byte{0x01} // Key prefix for subscription
	RequestKeyPrefix      = []byte{0x02} // Key prefix for request
)

// get guardian key from eth address
func GetSubscriptionKey(ethAddress string) []byte {
	return append(SubscriptionKeyPrefix, []byte(ethAddress)...)
}

// get request key from channelID
func GetRequestKey(channelId []byte) []byte {
	return append(RequestKeyPrefix, channelId...)
}
