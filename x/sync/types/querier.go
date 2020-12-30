package types

// DONTCOVER

// query endpoints supported by the sync Querier
const (
	QueryParams  = "params"
	QueryChange  = "change"
	QueryChanges = "changes"

	ParamVoting   = "voting"
	ParamTallying = "tallying"
)

// QueryChangeParams Params for queries:
// - 'custom/sync/change'
type QueryChangeParams struct {
	ID uint64
}

// NewQueryChangeParams creates a new instance of QueryChangeParams
func NewQueryChangeParams(changeID uint64) QueryChangeParams {
	return QueryChangeParams{
		ID: changeID,
	}
}
