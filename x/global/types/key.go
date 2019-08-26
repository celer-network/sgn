package types

const (
	// module name
	ModuleName = "global"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	LatestBlockKey = []byte{0x01} // Prefix for lastest block
)
