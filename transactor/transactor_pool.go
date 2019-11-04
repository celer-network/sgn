package transactor

import (
	"sync/atomic"

	"github.com/cosmos/cosmos-sdk/codec"
)

type TransactorPool struct {
	transactors []*Transactor
	index       uint64
}

func NewTransactorPool(cliHome, chainID, nodeURI, passphrase, gasPrice string, ts []string, cdc *codec.Codec) (*TransactorPool, error) {
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

// Get a transactor from the pool
func (t *TransactorPool) GetTransactor() *Transactor {
	transactor := t.transactors[t.index%uint64(len(t.transactors))]
	atomic.AddUint64(&t.index, 1)
	return transactor
}
