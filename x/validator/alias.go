package validator

import (
	"github.com/celer-network/sgn/x/validator/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewMsgClaimValidator = types.NewMsgClaimValidator
	ModuleCdc            = types.ModuleCdc
	RegisterCodec        = types.RegisterCodec
)

type (
	MsgClaimValidator = types.MsgClaimValidator
)
