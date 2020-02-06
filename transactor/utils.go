package transactor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
