package bridge

import (
	"github.com/celer-network/sgn/x/bridge/types"
)

const (
	ModuleName      = types.ModuleName
	RouterKey       = types.RouterKey
	StoreKey        = types.StoreKey
	QueryEthAddress = types.QueryEthAddress
)

var (
	NewMsgSetEthAddress = types.NewMsgSetEthAddress
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
	NewEthAddress       = types.NewEthAddress
)

type (
	MsgSetEthAddress      = types.MsgSetEthAddress
	QueryResEthAddress    = types.QueryResEthAddress
	EthAddress            = types.EthAddress
	QueryEthAddressParams = types.QueryEthAddressParams
)
