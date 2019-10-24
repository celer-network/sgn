package types

const (
	// module name
	ModuleName = "cron"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	DailyTimestampKey = []byte{0x01} // key for daily
)
