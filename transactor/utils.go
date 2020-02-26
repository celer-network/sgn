package transactor

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/viper"
)

func ParseTransactorAddrs(ts []string) ([]sdk.AccAddress, error) {
	var transactors []sdk.AccAddress
	for _, t := range ts {
		transactor, err := sdk.AccAddressFromBech32(t)
		if err != nil {
			return transactors, err
		}

		transactors = append(transactors, transactor)
	}

	return transactors, nil
}

func NewTxBuilder() types.TxBuilder {
	fees, err := sdk.ParseCoins(viper.GetString(flags.FlagFees))
	if err != nil {
		panic(err)
	}

	gasPrices, err := sdk.ParseDecCoins(viper.GetString(flags.FlagGasPrices))
	if err != nil {
		panic(err)
	}

	txbldr := types.NewTxBuilder(nil, viper.GetUint64(flags.FlagAccountNumber), viper.GetUint64(flags.FlagSequence), flags.GasFlagVar.Gas, viper.GetFloat64(flags.FlagGasAdjustment), flags.GasFlagVar.Simulate, viper.GetString(flags.FlagChainID), viper.GetString(flags.FlagMemo), fees, gasPrices)

	return txbldr
}
