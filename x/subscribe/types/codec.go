package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSubscribe{}, "subscribe/MsgSubscribe", nil)
	cdc.RegisterConcrete(MsgRequestGuard{}, "subscribe/MsgRequestGuard", nil)
	cdc.RegisterConcrete(MsgIntendSettle{}, "subscribe/MsgIntendSettle", nil)
	cdc.RegisterConcrete(MsgGuardProof{}, "subscribe/MsgGuardProof", nil)
}
