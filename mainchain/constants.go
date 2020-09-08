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
	ProposalDeposit     = 0
	GovernVoteTimeout   = 1
	SlashTimeout        = 2
	MinValidatorNum     = 3
	MaxValidatorNum     = 4
	MinStakeInPool      = 5
	AdvanceNoticePeriod = 6
	MigrationTime       = 7

	TxFailure = 0
	TxSuccess = 1

	// channel status of CelerChannel
	ChannelStatus_UNINITIALIZED uint8 = 0
	ChannelStatus_OPERABLE      uint8 = 1
	ChannelStatus_SETTLING      uint8 = 2
	ChannelStatus_CLOSED        uint8 = 3
	ChannelStatus_MIGRATED      uint8 = 4
)
