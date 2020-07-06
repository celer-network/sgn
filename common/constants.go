package common

const (
	QuotaCoinName = "quota"
	TokenDec      = 1000000000000000000

	// state of simplex channel guard request
	GuardStatus_Idle     uint8 = 0
	GuardStatus_Withdraw uint8 = 1
	GuardStatus_Settling uint8 = 2
	GuardStatus_Settled  uint8 = 3
)
