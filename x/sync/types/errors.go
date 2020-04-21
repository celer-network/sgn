package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sync module sentinel errors
var (
	ErrUnknownChange     = sdkerrors.Register(ModuleName, 1, "unknown change")
	ErrInactiveChange    = sdkerrors.Register(ModuleName, 2, "inactive change")
	ErrInvalidChangeType = sdkerrors.Register(ModuleName, 3, "invalid change type")
	ErrInvalidChangeData = sdkerrors.Register(ModuleName, 4, "invalid change data")
)
