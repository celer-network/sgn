package types

import "fmt"

// Query Result Payload for a names query
type QueryResNumber struct {
	Value []byte `json:"value"`
}

// implement fmt.Stringer
func (n QueryResNumber) String() string {
	return fmt.Sprint(n.Value)
}
