package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sync module sentinel errors
var (
	ErrInvalidChangeType = sdkerrors.Register(ModuleName, 1, "invalid change type")
	ErrInvalidChangeData = sdkerrors.Register(ModuleName, 2, "invalid change data")
)
