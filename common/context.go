package common

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/viper"
)

func NewQueryCLIContext(cdc *codec.Codec) context.CLIContext {
	ctx := context.NewCLIContext().
		WithCodec(cdc).
		WithNodeURI(viper.GetString(FlagSgnNodeURI)).
		WithChainID(viper.GetString(FlagSgnChainID))
	return ctx
}
