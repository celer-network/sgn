package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// module codec
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for
// sync.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSubmitChange{}, "cosmos-sdk/MsgSubmitChange", nil)
	cdc.RegisterConcrete(MsgApprove{}, "cosmos-sdk/MsgApprove", nil)
}

// RegisterChangeTypeCodec registers an external change content type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitChange
// to be correctly Amino encoded and decoded.
func RegisterChangeTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

// TODO determine a good place to seal this codec
func init() {
	RegisterCodec(ModuleCdc)
}
