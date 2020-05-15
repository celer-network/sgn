package types

// Governance module event types
const (
	EventTypeSync         = "sync"
	EventTypeActiveChange = "active_change"

	AttributeKeyChangeResult   = "change_result"
	AttributeKeyChangeID       = "change_id"
	AttributeValueChangePassed = "change_passed" // met vote quorum
	AttributeValueChangeFailed = "change_failed" // error on change handler

	ActionSubmitChange = "submit_change"
)
