package mainchain

const (
	// ValidatorChangeType
	AddValidator    = 0
	RemoveValidator = 1

	// CandidateStatus
	Unbonded  = 0
	Bonded    = 1
	Unbonding = 2

	// ParamNames
	ProposalDeposit      = 0
	GovernVoteTimeout    = 1
	BlameTimeout         = 2
	MinValidatorNum      = 3
	MaxValidatorNum      = 4
	MinStakeInPool       = 5
	IncreaseRateWaitTime = 6

	TxFailure = 0
	TxSuccess = 1

	// channel status of CelerChannel
	UninitializedChannel = uint8(0)
	OperableChannel      = uint8(1)
	SettlingChannel      = uint8(2)
	ClosedChannel        = uint8(3)
	MigratedChannel      = uint8(4)
)
