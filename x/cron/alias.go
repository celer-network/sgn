package cron

import (
	"github.com/celer-network/sgn/x/cron/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
)

var (
	ModuleCdc         = types.ModuleCdc
	RegisterCodec     = types.RegisterCodec
	DailyTimestampKey = types.DailyTimestampKey
)

type ()
