package types

import (
	"fmt"
	"strings"
)

type Block struct {
	Number uint64 `json:"number"`
	Time   uint64 `json:"time"`
}

// Returns a new Number with the minprice as the price
func NewBlock(number uint64, time uint64) Block {
	return Block{
		Number: number,
		Time:   time,
	}
}

// implement fmt.Stringer
func (b Block) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Number: %d, Time: %d`, b.Number, b.Time))
}
