package common

const (
	QuotaCoinName = "quota"
	TokenDec      = 1000000000000000000

	// state of simplex channel guard request
	GuardState_Idle     uint8 = 0
	GuardState_Withdraw uint8 = 1
	GuardState_Settling uint8 = 2
	GuardState_Settled  uint8 = 3
)
