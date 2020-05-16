package types

const (
	QueryEpoch      = "epoch"
	QueryParameters = "parameters"
)

type QueryEpochParams struct {
	EpochId int64
}

func NewQueryEpochParams(epochId int64) QueryEpochParams {
	return QueryEpochParams{
		EpochId: epochId,
	}
}
