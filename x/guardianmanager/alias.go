package guardianmanager

import (
	"github.com/celer-network/sgn/x/guardianmanager/types"
)

const (
	ModuleName    = types.ModuleName
	RouterKey     = types.RouterKey
	StoreKey      = types.StoreKey
	QueryGuardian = types.QueryGuardian
	QueryRequest  = types.QueryRequest
)

var (
	ModuleCdc          = types.ModuleCdc
	RegisterCodec      = types.RegisterCodec
	NewGuardian        = types.NewGuardian
	NewRequest         = types.NewRequest
	NewMsgMsgDeposit   = types.NewMsgDeposit
	NewMsgRequestGuard = types.NewMsgRequestGuard
	GetGuardianKey     = types.GetGuardianKey
	GetRequestKey      = types.GetRequestKey
)

type (
	MsgDeposit          = types.MsgDeposit
	Guardian            = types.Guardian
	Request             = types.Request
	QueryGuardianParams = types.QueryGuardianParams
	QueryRequestParams  = types.QueryRequestParams
	MsgRequestGuard     = types.MsgRequestGuard
)
