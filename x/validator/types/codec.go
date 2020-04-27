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
	cdc.RegisterConcrete(MsgUpdateSidechainAddr{}, "validator/MsgUpdateSidechainAddr", nil)
	cdc.RegisterConcrete(MsgSetTransactors{}, "validator/MsgSetTransactors", nil)
	cdc.RegisterConcrete(MsgClaimValidator{}, "validator/MsgClaimValidator", nil)
	cdc.RegisterConcrete(MsgSyncValidator{}, "validator/MsgSyncValidator", nil)
	cdc.RegisterConcrete(MsgSyncDelegator{}, "validator/MsgSyncDelegator", nil)
	cdc.RegisterConcrete(MsgWithdrawReward{}, "validator/MsgWithdrawReward", nil)
	cdc.RegisterConcrete(MsgSignReward{}, "validator/MsgSignReward", nil)
}
