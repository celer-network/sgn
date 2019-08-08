package types

const (
	QuerySubscrption = "subscription"
)

type QuerySubscrptionParams struct {
	EthAddress string
}

func NewQuerySubscrptionParams(addr string) QuerySubscrptionParams {
	return QuerySubscrptionParams{
		EthAddress: addr,
	}
}
