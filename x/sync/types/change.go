package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// DefaultStartingChangeID is 1
const DefaultStartingChangeID uint64 = 1

// Change defines a struct used by the sync module to allow for voting
// on network changes.
type Change struct {
	ChangeID      uint64       `json:"id" yaml:"id"` //  ID of the change
	Type          string       `json:"type" yaml:"type"`
	Status        ChangeStatus `json:"change_status" yaml:"change_status"`     // Status of the Change {Pending, Active, Passed, Rejected}
	SubmitTime    time.Time    `json:"submit_time" yaml:"submit_time"`         // Time of the block where TxGovSubmitChange was included
	VotingEndTime time.Time    `json:"voting_end_time" yaml:"voting_end_time"` // Time that the VotingPeriod for this change will end and votes will be tallied
}

// NewChange creates a new Change instance
func NewChange(id uint64, submitTime time.Time) Change {
	return Change{
		ChangeID:   id,
		Status:     StatusVotingPeriod,
		SubmitTime: submitTime,
	}
}

// String implements stringer interface
func (c Change) String() string {
	return fmt.Sprintf(`Change %d:
  Type:               %s
  Status:             %s
  Submit Time:        %s
  Voting End Time:    %s`,
		c.ChangeID, c.Type,
		c.Status, c.VotingEndTime,
	)
}

// Changes is an array of change
type Changes []Change

// String implements stringer interface
func (changes Changes) String() string {
	out := "ID - (Status) [Type]\n"
	for _, change := range changes {
		out += fmt.Sprintf("%d - (%s) [%s]\n",
			change.ChangeID, change.Status,
			change.Type)
	}
	return strings.TrimSpace(out)
}

type (
	// ChangeQueue defines a queue for change ids
	ChangeQueue []uint64

	// ChangeStatus is a type alias that represents a change status as a byte
	ChangeStatus byte
)

// Valid Change statuses
const (
	StatusNil          ChangeStatus = 0x00
	StatusVotingPeriod ChangeStatus = 0x01
	StatusPassed       ChangeStatus = 0x02
	StatusFailed       ChangeStatus = 0x03
)

// ChangeStatusFromString turns a string into a ChangeStatus
func ChangeStatusFromString(str string) (ChangeStatus, error) {
	switch str {
	case "VotingPeriod":
		return StatusVotingPeriod, nil
	case "Passed":
		return StatusPassed, nil
	case "Failed":
		return StatusFailed, nil
	case "":
		return StatusNil, nil

	default:
		return ChangeStatus(0xff), fmt.Errorf("'%s' is not a valid change status", str)
	}
}

// ValidChangeStatus returns true if the change status is valid and false
// otherwise.
func ValidChangeStatus(status ChangeStatus) bool {
	if status == StatusVotingPeriod ||
		status == StatusPassed ||
		status == StatusFailed {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status ChangeStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *ChangeStatus) Unmarshal(data []byte) error {
	*status = ChangeStatus(data[0])
	return nil
}

// MarshalJSON Marshals to JSON using string representation of the status
func (status ChangeStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// UnmarshalJSON Unmarshals from JSON assuming Bech32 encoding
func (status *ChangeStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ChangeStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

// String implements the Stringer interface.
func (status ChangeStatus) String() string {
	switch status {
	case StatusVotingPeriod:
		return "VotingPeriod"
	case StatusPassed:
		return "Passed"
	case StatusFailed:
		return "Failed"

	default:
		return ""
	}
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (status ChangeStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}
