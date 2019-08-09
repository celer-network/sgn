package types

const (
	QueryGuardian = "guardian"
)

type QueryGuardianParams struct {
	EthAddress string
}

func NewQueryGuardianParams(addr string) QueryGuardianParams {
	return QueryGuardianParams{
		EthAddress: addr,
	}
}
