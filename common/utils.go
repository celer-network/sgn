package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ParseTransactorAddrs(ts []string) (transactors []sdk.AccAddress, err error) {
	for _, t := range ts {
		transactor, err := sdk.AccAddressFromBech32(t)
		if err != nil {
			return transactors, err
		}

		transactors = append(transactors, transactor)
	}

	return transactors, err
}
