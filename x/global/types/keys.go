package types

const (
	// ModuleName is the name of the module
	ModuleName = "global"

	// StoreKey is the store key string for global
	StoreKey = ModuleName

	// RouterKey is the message route for global
	RouterKey = ModuleName

	// QuerierRoute is the querier route for global
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName
)

var (
	EthBlkNumKey = []byte{0x01}
)
