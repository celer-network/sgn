package types

const (
	// module name
	ModuleName = "guardianmanager"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	GuardianKey = []byte{0x01} // Prefix for guardian
	RequestKey  = []byte{0x02} // Prefix for request
)

// get guardian key from eth address
func GetGuardianKey(ethAddress string) []byte {
	return append(GuardianKey, []byte(ethAddress)...)
}

// get guardian key from eth address
func GetRequestKey(channelId []byte) []byte {
	return append(GuardianKey, channelId...)
}
