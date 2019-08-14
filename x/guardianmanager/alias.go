package guardianmanager

import (
	"github.com/celer-network/sgn/x/guardianmanager/types"
)

const (
	ModuleName    = types.ModuleName
	RouterKey     = types.RouterKey
	StoreKey      = types.StoreKey
	QueryGuardian = types.QueryGuardian
)

var (
	ModuleCdc          = types.ModuleCdc
	RegisterCodec      = types.RegisterCodec
	NewGuardian        = types.NewGuardian
	NewMsgMsgDeposit   = types.NewMsgDeposit
	NewMsgRequestGuard = types.NewMsgRequestGuard
)

type (
	MsgDeposit          = types.MsgDeposit
	Guardian            = types.Guardian
	QueryGuardianParams = types.QueryGuardianParams
	MsgRequestGuard     = types.MsgRequestGuard
)
