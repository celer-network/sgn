package mainchain

const (
	AddValidator    = 0
	RemoveValidator = 1

	Unbonded  = 0
	Bonded    = 1
	Unbonding = 2

	TxFailure = 0
	TxSuccess = 1

	// channel status of CelerChannel
	UninitializedChannel = uint8(0)
	OperableChannel      = uint8(1)
	SettlingChannel      = uint8(2)
	ClosedChannel        = uint8(3)
	MigratedChannel      = uint8(4)
)
