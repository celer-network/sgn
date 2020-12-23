package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DONTCOVER

// query endpoints supported by the sync Querier
const (
	QueryBlkNum        = "blkNum"
	QueryParams        = "params"
	QueryChange        = "change"
	QueryChanges       = "changes"
	QueryActiveChanges = "activeChanges"

	ParamVoting   = "voting"
	ParamTallying = "tallying"
)

// QueryChangeParams Params for queries:
// - 'custom/sync/change'
// - 'custom/sync/deposits'
// - 'custom/sync/tally'
type QueryChangeParams struct {
	ID uint64
}

// NewQueryChangeParams creates a new instance of QueryChangeParams
func NewQueryChangeParams(changeID uint64) QueryChangeParams {
	return QueryChangeParams{
		ID: changeID,
	}
}

// QueryChangesParams Params for query 'custom/sync/changes'
type QueryChangesParams struct {
	Page         int
	Limit        int
	Voter        sdk.AccAddress
	ChangeStatus ChangeStatus
}

// NewQueryChangesParams creates a new instance of QueryChangesParams
func NewQueryChangesParams(page, limit int, status ChangeStatus) QueryChangesParams {
	return QueryChangesParams{
		Page:         page,
		Limit:        limit,
		ChangeStatus: status,
	}
}
