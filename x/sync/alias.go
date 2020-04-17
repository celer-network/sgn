package sync

// nolint

import (
	"github.com/celer-network/sgn/x/sync/client/cli"
	"github.com/celer-network/sgn/x/sync/keeper"
	"github.com/celer-network/sgn/x/sync/types"
)

const (
	MaxDescriptionLength  = types.MaxDescriptionLength
	MaxTitleLength        = types.MaxTitleLength
	DefaultPeriod         = types.DefaultPeriod
	ModuleName            = types.ModuleName
	StoreKey              = types.StoreKey
	RouterKey             = types.RouterKey
	QuerierRoute          = types.QuerierRoute
	DefaultParamspace     = types.DefaultParamspace
	TypeMsgDeposit        = types.TypeMsgDeposit
	TypeMsgVote           = types.TypeMsgVote
	TypeMsgSubmitChange = types.TypeMsgSubmitChange
	StatusNil             = types.StatusNil
	StatusDepositPeriod   = types.StatusDepositPeriod
	StatusVotingPeriod    = types.StatusVotingPeriod
	StatusPassed          = types.StatusPassed
	StatusRejected        = types.StatusRejected
	StatusFailed          = types.StatusFailed
	ChangeTypeText      = types.ChangeTypeText
	QueryParams           = types.QueryParams
	QueryChanges        = types.QueryChanges
	QueryChange         = types.QueryChange
	QueryDeposits         = types.QueryDeposits
	QueryDeposit          = types.QueryDeposit
	QueryVotes            = types.QueryVotes
	QueryVote             = types.QueryVote
	QueryTally            = types.QueryTally
	ParamDeposit          = types.ParamDeposit
	ParamVoting           = types.ParamVoting
	ParamTallying         = types.ParamTallying
	OptionEmpty           = types.OptionEmpty
	OptionYes             = types.OptionYes
	OptionAbstain         = types.OptionAbstain
	OptionNo              = types.OptionNo
	OptionNoWithVeto      = types.OptionNoWithVeto
)

var (
	// functions aliases
	NewKeeper                     = keeper.NewKeeper
	NewQuerier                    = keeper.NewQuerier
	RegisterCodec                 = types.RegisterCodec
	RegisterChangeTypeCodec     = types.RegisterChangeTypeCodec
	ValidateAbstract              = types.ValidateAbstract
	NewDeposit                    = types.NewDeposit
	ErrUnknownChange            = types.ErrUnknownChange
	ErrInactiveChange           = types.ErrInactiveChange
	ErrAlreadyActiveChange      = types.ErrAlreadyActiveChange
	ErrInvalidChangeContent     = types.ErrInvalidChangeContent
	ErrInvalidChangeType        = types.ErrInvalidChangeType
	ErrInvalidVote                = types.ErrInvalidVote
	ErrInvalidGenesis             = types.ErrInvalidGenesis
	ErrNoChangeHandlerExists    = types.ErrNoChangeHandlerExists
	NewGenesisState               = types.NewGenesisState
	DefaultGenesisState           = types.DefaultGenesisState
	ValidateGenesis               = types.ValidateGenesis
	GetChangeIDBytes            = types.GetChangeIDBytes
	GetChangeIDFromBytes        = types.GetChangeIDFromBytes
	ChangeKey                   = types.ChangeKey
	ActiveChangeByTimeKey       = types.ActiveChangeByTimeKey
	ActiveChangeQueueKey        = types.ActiveChangeQueueKey
	InactiveChangeByTimeKey     = types.InactiveChangeByTimeKey
	InactiveChangeQueueKey      = types.InactiveChangeQueueKey
	DepositsKey                   = types.DepositsKey
	DepositKey                    = types.DepositKey
	VotesKey                      = types.VotesKey
	VoteKey                       = types.VoteKey
	SplitChangeKey              = types.SplitChangeKey
	SplitActiveChangeQueueKey   = types.SplitActiveChangeQueueKey
	SplitInactiveChangeQueueKey = types.SplitInactiveChangeQueueKey
	SplitKeyDeposit               = types.SplitKeyDeposit
	SplitKeyVote                  = types.SplitKeyVote
	NewMsgSubmitChange          = types.NewMsgSubmitChange
	NewMsgDeposit                 = types.NewMsgDeposit
	NewMsgVote                    = types.NewMsgVote
	ParamKeyTable                 = types.ParamKeyTable
	NewDepositParams              = types.NewDepositParams
	NewTallyParams                = types.NewTallyParams
	NewVotingParams               = types.NewVotingParams
	NewParams                     = types.NewParams
	NewChange                   = types.NewChange
	NewRouter                     = types.NewRouter
	ChangeStatusFromString      = types.ChangeStatusFromString
	ValidChangeStatus           = types.ValidChangeStatus
	NewTextChange               = types.NewTextChange
	RegisterChangeType          = types.RegisterChangeType
	ContentFromChangeType       = types.ContentFromChangeType
	IsValidChangeType           = types.IsValidChangeType
	ChangeHandler               = types.ChangeHandler
	NewQueryChangeParams        = types.NewQueryChangeParams
	NewQueryDepositParams         = types.NewQueryDepositParams
	NewQueryVoteParams            = types.NewQueryVoteParams
	NewQueryChangesParams       = types.NewQueryChangesParams
	NewValidatorGovInfo           = types.NewValidatorGovInfo
	NewTallyResult                = types.NewTallyResult
	NewTallyResultFromMap         = types.NewTallyResultFromMap
	EmptyTallyResult              = types.EmptyTallyResult
	NewVote                       = types.NewVote
	VoteOptionFromString          = types.VoteOptionFromString
	ValidVoteOption               = types.ValidVoteOption
	CLIQueryChanges             = cli.QueryChanges
	CLIQueryChange              = cli.QueryChange

	// variable aliases
	ModuleCdc                   = types.ModuleCdc
	ChangesKeyPrefix          = types.ChangesKeyPrefix
	ActiveChangeQueuePrefix   = types.ActiveChangeQueuePrefix
	InactiveChangeQueuePrefix = types.InactiveChangeQueuePrefix
	ChangeIDKey               = types.ChangeIDKey
	DepositsKeyPrefix           = types.DepositsKeyPrefix
	VotesKeyPrefix              = types.VotesKeyPrefix
	ParamStoreKeyDepositParams  = types.ParamStoreKeyDepositParams
	ParamStoreKeyVotingParams   = types.ParamStoreKeyVotingParams
	ParamStoreKeyTallyParams    = types.ParamStoreKeyTallyParams
)

type (
	Keeper               = keeper.Keeper
	Content              = types.Content
	Handler              = types.Handler
	Deposit              = types.Deposit
	Deposits             = types.Deposits
	GenesisState         = types.GenesisState
	MsgSubmitChange    = types.MsgSubmitChange
	MsgDeposit           = types.MsgDeposit
	MsgVote              = types.MsgVote
	DepositParams        = types.DepositParams
	TallyParams          = types.TallyParams
	VotingParams         = types.VotingParams
	Params               = types.Params
	Change             = types.Change
	Changes            = types.Changes
	ChangeQueue        = types.ChangeQueue
	ChangeStatus       = types.ChangeStatus
	TextChange         = types.TextChange
	QueryChangeParams  = types.QueryChangeParams
	QueryDepositParams   = types.QueryDepositParams
	QueryVoteParams      = types.QueryVoteParams
	QueryChangesParams = types.QueryChangesParams
	ValidatorGovInfo     = types.ValidatorGovInfo
	TallyResult          = types.TallyResult
	Vote                 = types.Vote
	Votes                = types.Votes
	VoteOption           = types.VoteOption
)
