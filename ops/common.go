package ops

import (
	"math/big"

	"github.com/celer-network/sgn/common"
)

const (
	candidateFlag = "candidate"
	amountFlag    = "amount"
)

func calcRawAmount(amount string) *big.Int {
	rawAmount := new(big.Int)
	rawAmount.SetString(amount, 10)
	return new(big.Int).Mul(rawAmount, big.NewInt(common.TokenDec))
}
