package common

type GuardStatus uint8

const (
	QuotaCoinName = "quota"
	TokenDec      = 1000000000000000000

	// state of simplex channel guard request
	GuardStatus_Idle        GuardStatus = 0
	GuardStatus_Withdrawing GuardStatus = 1
	GuardStatus_Settling    GuardStatus = 2
	GuardStatus_Settled     GuardStatus = 3
)

func (status GuardStatus) String() string {
	switch status {
	case GuardStatus_Idle:
		return "Idle"
	case GuardStatus_Withdrawing:
		return "Withdraw"
	case GuardStatus_Settling:
		return "Settling"
	case GuardStatus_Settled:
		return "Settled"
	default:
		return "Invalid"
	}
}
