package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sync module sentinel errors
var (
	ErrUnknownChange     = sdkerrors.Register(ModuleName, 1, "unknown change")
	ErrDoubleVote        = sdkerrors.Register(ModuleName, 2, "double vote")
	ErrInvalidChangeType = sdkerrors.Register(ModuleName, 3, "invalid change type")
	ErrInvalidChangeData = sdkerrors.Register(ModuleName, 4, "invalid change data")
	ErrInvalidGenesis    = sdkerrors.Register(ModuleName, 5, "invalid genesis state")
)
