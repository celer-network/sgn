package watcher

// LogEventID tracks the position of a watch event in the event log.
type LogEventID struct {
	BlockNumber uint64 // Number of the block containing the event
	Index       int64  // Index of the event within the block
}
