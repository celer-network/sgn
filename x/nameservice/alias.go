package nameservice

import (
	"github.com/cosmos/sdk-application-tutorial/x/nameservice/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewMsgSetNumber = types.NewMsgSetNumber
	ModuleCdc       = types.ModuleCdc
	RegisterCodec   = types.RegisterCodec
	NewSimple       = types.NewNumber
)

type (
	MsgSetNumber   = types.MsgSetNumber
	QueryResNumber = types.QueryResNumber
	Number         = types.Number
)
