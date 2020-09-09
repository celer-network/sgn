package transactor

import (
	"sync/atomic"

	"github.com/cosmos/cosmos-sdk/codec"
)

type TransactorPool struct {
	transactors []*Transactor
	index       uint64
	cliHome     string
	chainID     string
	cdc         *codec.Codec
	gpe         *GasPriceEstimator
}

func NewTransactorPool(cliHome, chainID string, cdc *codec.Codec, gpe *GasPriceEstimator) *TransactorPool {
	return &TransactorPool{
		transactors: []*Transactor{},
		index:       0,
		cliHome:     cliHome,
		chainID:     chainID,
		cdc:         cdc,
		gpe:         gpe,
	}
}

// Add transactors to the pool
func (tp *TransactorPool) AddTransactor(transactor *Transactor) {
	tp.transactors = append(tp.transactors, transactor)
}

// Add transactors to the pool
func (tp *TransactorPool) AddTransactors(nodeURI, passphrase string, ts []string) error {
	var transactors []*Transactor
	for _, t := range ts {
		transactor, err := NewTransactor(tp.cliHome, tp.chainID, nodeURI, t, passphrase, tp.cdc, tp.gpe)
		if err != nil {
			return err
		}
		transactor.Run()

		transactors = append(transactors, transactor)
	}

	tp.transactors = append(tp.transactors, transactors...)
	return nil
}

// Get a transactor from the pool
func (tp *TransactorPool) GetTransactor() *Transactor {
	if len(tp.transactors) == 0 {
		return nil
	}

	transactor := tp.transactors[tp.index%uint64(len(tp.transactors))]
	atomic.AddUint64(&tp.index, 1)
	return transactor
}
