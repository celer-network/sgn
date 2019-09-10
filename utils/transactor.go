package utils

import (
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gammazero/deque"
)

type Transactor struct {
	TxBuilder  types.TxBuilder
	CliCtx     context.CLIContext
	Key        keys.Info
	Passphrase string
	txQueue    deque.Deque
}

func NewTransactor(cliHome, chainID, nodeURI, accName, passphrase string, cdc *codec.Codec) (*Transactor, error) {
	kb, err := client.NewKeyBaseFromDir(cliHome)
	if err != nil {
		return nil, err
	}

	key, err := kb.Get(accName)
	if err != nil {
		return nil, err
	}

	txBldr := auth.
		NewTxBuilderFromCLI().
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID).
		WithKeybase(kb)
	cliCtx := client.
		NewCLIContext().
		WithCodec(cdc).
		WithFromAddress(key.GetAddress()).
		WithFromName(key.GetName()).
		WithNodeURI(nodeURI).
		WithTrustNode(true).
		WithBroadcastMode(flags.BroadcastSync)

	transactor := &Transactor{
		TxBuilder:  txBldr,
		CliCtx:     cliCtx,
		Key:        key,
		Passphrase: passphrase,
	}

	go transactor.start()
	return transactor, nil
}

func (t *Transactor) BroadcastTx(msg sdk.Msg) {
	t.txQueue.PushBack(msg)
}

func (t *Transactor) start() {
	for {
		if t.txQueue.Len() == 0 {
			time.Sleep(time.Second)
			continue
		}

		var msgs []sdk.Msg
		for t.txQueue.Len() != 0 {
			msg := t.txQueue.PopFront().(sdk.Msg)
			msgs = append(msgs, msg)
		}

		txBldr, err := utils.PrepareTxBuilder(t.TxBuilder, t.CliCtx)
		if err != nil {
			log.Printf("Transactor PrepareTxBuilder err", err)
			continue
		}

		txBytes, err := txBldr.BuildAndSign(t.Key.GetName(), t.Passphrase, msgs)
		if err != nil {
			log.Printf("Transactor BuildAndSign err", err)
			continue
		}

		tx, err := t.CliCtx.BroadcastTx(txBytes)
		if err != nil {
			log.Printf("Transactor BroadcastTx err", err)
			continue
		}

		log.Printf("Transactor tx", tx)
	}
}
