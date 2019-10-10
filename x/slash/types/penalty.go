package types

import (
	"fmt"
	"strings"
)

type Penalty struct {
	Nounce uint64 `json:"nounce"`
}

func NewPenalty(nounce uint64) Penalty {
	return Penalty{
		Nounce: nounce,
	}
}

// implement fmt.Stringer
func (s Penalty) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Nounce: %d`, s.Nounce))
}
