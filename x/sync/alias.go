package sync

// nolint

import (
	"github.com/celer-network/sgn/x/sync/client/cli"
	"github.com/celer-network/sgn/x/sync/keeper"
	"github.com/celer-network/sgn/x/sync/types"
)

const (
	DefaultPeriod        = types.DefaultPeriod
	ModuleName           = types.ModuleName
	StoreKey             = types.StoreKey
	RouterKey            = types.RouterKey
	QuerierRoute         = types.QuerierRoute
	DefaultParamspace    = types.DefaultParamspace
	TypeMsgApprove       = types.TypeMsgApprove
	TypeMsgSubmitChange  = types.TypeMsgSubmitChange
	StatusNil            = types.StatusNil
	StatusActive         = types.StatusActive
	StatusPassed         = types.StatusPassed
	StatusFailed         = types.StatusFailed
	QueryParams          = types.QueryParams
	QueryChanges         = types.QueryChanges
	QueryChange          = types.QueryChange
	ParamVoting          = types.ParamVoting
	ParamTallying        = types.ParamTallying
	SyncEthBlkNum        = types.SyncEthBlkNum
	ConfirmParamProposal = types.ConfirmParamProposal
	UpdateSidechainAddr  = types.UpdateSidechainAddr
	SyncDelegator        = types.SyncDelegator
	SyncValidator        = types.SyncValidator
	Subscribe            = types.Subscribe
	InitGuardRequest     = types.InitGuardRequest
	GuardTrigger         = types.GuardTrigger
	GuardProof           = types.GuardProof
	EventTypeSync        = types.EventTypeSync
	ActionSubmitChange   = types.ActionSubmitChange
)

var (
	// functions aliases
	NewKeeper               = keeper.NewKeeper
	NewQuerier              = keeper.NewQuerier
	RegisterCodec           = types.RegisterCodec
	RegisterChangeTypeCodec = types.RegisterChangeTypeCodec
	ErrUnknownChange        = types.ErrUnknownChange
	ErrDoubleVote           = types.ErrDoubleVote
	ErrInvalidChangeType    = types.ErrInvalidChangeType
	ErrInvalidGenesis       = types.ErrInvalidGenesis
	NewGenesisState         = types.NewGenesisState
	DefaultGenesisState     = types.DefaultGenesisState
	ValidateGenesis         = types.ValidateGenesis
	GetChangeIDBytes        = types.GetChangeIDBytes
	GetChangeIDFromBytes    = types.GetChangeIDFromBytes
	ChangeKey               = types.ChangeKey
	NewMsgApprove           = types.NewMsgApprove
	ParamKeyTable           = types.ParamKeyTable
	NewTallyParams          = types.NewTallyParams
	NewVotingParams         = types.NewVotingParams
	NewParams               = types.NewParams
	NewChange               = types.NewChange
	ChangeStatusFromString  = types.ChangeStatusFromString
	ValidChangeStatus       = types.ValidChangeStatus
	NewQueryChangeParams    = types.NewQueryChangeParams
	CLIQueryChange          = cli.QueryChange
	CLIQueryChanges         = cli.QueryChanges

	// variable aliases
	ModuleCdc                 = types.ModuleCdc
	ChangesKeyPrefix          = types.ChangesKeyPrefix
	ChangeIDKey               = types.ChangeIDKey
	ParamStoreKeyVotingParams = types.ParamStoreKeyVotingParams
	ParamStoreKeyTallyParams  = types.ParamStoreKeyTallyParams
)

type (
	Keeper            = keeper.Keeper
	GenesisState      = types.GenesisState
	MsgSubmitChange   = types.MsgSubmitChange
	MsgApprove        = types.MsgApprove
	TallyParams       = types.TallyParams
	VotingParams      = types.VotingParams
	Params            = types.Params
	Change            = types.Change
	Changes           = types.Changes
	ChangeQueue       = types.ChangeQueue
	ChangeStatus      = types.ChangeStatus
	QueryChangeParams = types.QueryChangeParams
)
