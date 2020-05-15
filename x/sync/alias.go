package sync

// nolint

import (
	"github.com/celer-network/sgn/x/sync/client/cli"
	"github.com/celer-network/sgn/x/sync/keeper"
	"github.com/celer-network/sgn/x/sync/types"
)

const (
	DefaultPeriod       = types.DefaultPeriod
	ModuleName          = types.ModuleName
	StoreKey            = types.StoreKey
	RouterKey           = types.RouterKey
	QuerierRoute        = types.QuerierRoute
	DefaultParamspace   = types.DefaultParamspace
	TypeMsgApprove      = types.TypeMsgApprove
	TypeMsgSubmitChange = types.TypeMsgSubmitChange
	StatusNil           = types.StatusNil
	StatusVotingPeriod  = types.StatusVotingPeriod
	StatusPassed        = types.StatusPassed
	StatusFailed        = types.StatusFailed
	QueryParams         = types.QueryParams
	QueryChanges        = types.QueryChanges
	QueryChange         = types.QueryChange
	ParamVoting         = types.ParamVoting
	ParamTallying       = types.ParamTallying
	SyncBlock           = types.SyncBlock
	Subscribe           = types.Subscribe
	Request             = types.Request
	UpdateSidechainAddr = types.UpdateSidechainAddr
	SyncDelegator       = types.SyncDelegator
	SyncValidator       = types.SyncValidator
	EventTypeSync       = types.EventTypeSync
	ActionSubmitChange  = types.ActionSubmitChange
)

var (
	// functions aliases
	NewKeeper                   = keeper.NewKeeper
	NewQuerier                  = keeper.NewQuerier
	RegisterCodec               = types.RegisterCodec
	RegisterChangeTypeCodec     = types.RegisterChangeTypeCodec
	ErrUnknownChange            = types.ErrUnknownChange
	ErrInactiveChange           = types.ErrInactiveChange
	ErrInvalidChangeType        = types.ErrInvalidChangeType
	ErrInvalidGenesis           = types.ErrInvalidGenesis
	NewGenesisState             = types.NewGenesisState
	DefaultGenesisState         = types.DefaultGenesisState
	ValidateGenesis             = types.ValidateGenesis
	GetChangeIDBytes            = types.GetChangeIDBytes
	GetChangeIDFromBytes        = types.GetChangeIDFromBytes
	ChangeKey                   = types.ChangeKey
	ActiveChangeByTimeKey       = types.ActiveChangeByTimeKey
	ActiveChangeQueueKey        = types.ActiveChangeQueueKey
	SplitChangeKey              = types.SplitChangeKey
	SplitActiveChangeQueueKey   = types.SplitActiveChangeQueueKey
	SplitInactiveChangeQueueKey = types.SplitInactiveChangeQueueKey
	NewMsgSubmitChange          = types.NewMsgSubmitChange
	NewMsgApprove               = types.NewMsgApprove
	ParamKeyTable               = types.ParamKeyTable
	NewTallyParams              = types.NewTallyParams
	NewVotingParams             = types.NewVotingParams
	NewParams                   = types.NewParams
	NewChange                   = types.NewChange
	ChangeStatusFromString      = types.ChangeStatusFromString
	ValidChangeStatus           = types.ValidChangeStatus
	NewQueryChangeParams        = types.NewQueryChangeParams
	NewQueryChangesParams       = types.NewQueryChangesParams
	CLIQueryChanges             = cli.QueryChanges
	CLIQueryChange              = cli.QueryChange

	// variable aliases
	ModuleCdc                 = types.ModuleCdc
	ChangesKeyPrefix          = types.ChangesKeyPrefix
	ActiveChangeQueuePrefix   = types.ActiveChangeQueuePrefix
	ChangeIDKey               = types.ChangeIDKey
	ParamStoreKeyVotingParams = types.ParamStoreKeyVotingParams
	ParamStoreKeyTallyParams  = types.ParamStoreKeyTallyParams
)

type (
	Keeper             = keeper.Keeper
	GenesisState       = types.GenesisState
	MsgSubmitChange    = types.MsgSubmitChange
	MsgApprove         = types.MsgApprove
	TallyParams        = types.TallyParams
	VotingParams       = types.VotingParams
	Params             = types.Params
	Change             = types.Change
	Changes            = types.Changes
	ChangeQueue        = types.ChangeQueue
	ChangeStatus       = types.ChangeStatus
	QueryChangeParams  = types.QueryChangeParams
	QueryChangesParams = types.QueryChangesParams
)
