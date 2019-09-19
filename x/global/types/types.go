package types

import (
	"fmt"
	"strings"
)

type Block struct {
	Number uint64 `json:"number"`
}

func NewBlock(number uint64) Block {
	return Block{
		Number: number,
	}
}

// implement fmt.Stringer
func (b Block) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Number: %d`, b.Number))
}
