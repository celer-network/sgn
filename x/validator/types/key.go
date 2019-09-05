package types

const (
	// module name
	ModuleName = "validator"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	PullerKey = []byte{0x01} // Prefix for puller
	PusherKey = []byte{0x02} // Prefix for pusher
)
