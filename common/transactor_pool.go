package common

import (
	"strings"
	"sync/atomic"

	"github.com/cosmos/cosmos-sdk/codec"
)

type TransactorPool struct {
	transactors []*Transactor
	index       uint64
}

func NewTransactorPool(cliHome, chainID, nodeURI, tsStr, passphrase, gasPrice string, cdc *codec.Codec) (*TransactorPool, error) {
	ts := strings.Split(tsStr, ",")
	var transactors []*Transactor
	for _, t := range ts {
		transactor, err := NewTransactor(cliHome, chainID, nodeURI, t, passphrase, gasPrice, cdc)
		if err != nil {
			return nil, err
		}

		transactors = append(transactors, transactor)
	}

	transactorPool := &TransactorPool{
		transactors: transactors,
		index:       0,
	}

	return transactorPool, nil
}

// Batch msg into a queue before actual broadcast
func (t *TransactorPool) GetTransactor() *Transactor {
	transactor := t.transactors[t.index]
	atomic.AddUint64(&t.index, 1)
	return transactor
}
