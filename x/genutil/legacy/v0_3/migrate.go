package v03

import (
	v02slash "github.com/celer-network/sgn/x/slash/legacy/v0_2"
	v03slash "github.com/celer-network/sgn/x/slash/legacy/v0_3"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
)

func Migrate(appState types.AppMap) types.AppMap {
	v02Codec := codec.New()
	codec.RegisterCrypto(v02Codec)

	v03Codec := codec.New()
	codec.RegisterCrypto(v03Codec)

	if appState[v02slash.ModuleName] != nil {
		var slashGenState v02slash.GenesisState
		v02Codec.MustUnmarshalJSON(appState[v02slash.ModuleName], &slashGenState)

		delete(appState, v02slash.ModuleName)
		appState[v03slash.ModuleName] = v03Codec.MustMarshalJSON(v03slash.Migrate(slashGenState))
	}

	return appState
}
