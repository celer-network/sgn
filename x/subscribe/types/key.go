package types

const (
	// module name
	ModuleName = "subscribe"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	SubscriptionKey = []byte{0x01} // Prefix for subscription
	RequestKey      = []byte{0x02} // Prefix for request
)

// get guardian key from eth address
func GetSubscriptionKey(ethAddress string) []byte {
	return append(SubscriptionKey, []byte(ethAddress)...)
}

// get request key from channelID
func GetRequestKey(channelId []byte) []byte {
	return append(RequestKey, channelId...)
}
