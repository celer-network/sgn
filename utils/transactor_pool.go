package utils

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TransactorPool struct {
	transactors []Transactor
}

func NewTransactorPool(cliHome, chainID, nodeURI, transactors, passphrase, gasPrice string, cdc *codec.Codec) (*TransactorPool, error) {
	kb, err := client.NewKeyBaseFromDir(cliHome)
	if err != nil {
		return nil, err
	}

	key, err := kb.Get(transactors)
	if err != nil {
		return nil, err
	}

	transactor := &TransactorPool{}

	go transactor.start()
	return transactor, nil
}

// Batch msg into a queue before actual broadcast
func (t *TransactorPool) BroadcastTx(msg sdk.Msg) {
}
