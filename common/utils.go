package common

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ParseTransactors(tsStr string) (transactors []sdk.AccAddress, err error) {
	ts := strings.Split(tsStr, ",")
	for _, t := range ts {
		transactor, err := sdk.AccAddressFromBech32(t)
		if err != nil {
			return transactors, err
		}

		transactors = append(transactors, transactor)
	}

	return transactors, err
}
