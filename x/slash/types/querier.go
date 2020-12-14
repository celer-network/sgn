package types

const (
	QueryPenalty        = "penalty"
	QueryPenalties      = "penalties"
	QueryPenaltyRequest = "penalty-request"
	QueryParameters     = "parameters"
)

type QueryPenaltyParams struct {
	Nonce uint64
}

func NewQueryPenaltyParams(nonce uint64) QueryPenaltyParams {
	return QueryPenaltyParams{
		Nonce: nonce,
	}
}
