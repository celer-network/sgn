package common

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

const maxTry = 5

type Transactor struct {
	TxBuilder  types.TxBuilder
	CliCtx     context.CLIContext
	Key        keys.Info
	Passphrase string
	msgQueue   deque.Deque
}

func NewTransactor(cliHome, chainID, nodeURI, accAddr, passphrase, gasPrice string, cdc *codec.Codec) (*Transactor, error) {
	kb, err := client.NewKeyBaseFromDir(cliHome)
	if err != nil {
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(accAddr)
	if err != nil {
		return nil, err
	}

	key, err := kb.GetByAddress(addr)
	if err != nil {
		return nil, err
	}

	txBldr := auth.
		NewTxBuilderFromCLI().
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID).
		WithKeybase(kb).
		WithGasPrices(gasPrice)

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

// Batch msg into a queue before actual broadcast
func (t *Transactor) BroadcastTx(msg sdk.Msg) {
	t.msgQueue.PushBack(msg)
}

// Poll tx queue and send msgs in batch
func (t *Transactor) start() {
	for {
		if t.msgQueue.Len() == 0 {
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Packing %d messages in batch for broadcast", t.msgQueue.Len())
		var msgs []sdk.Msg
		for t.msgQueue.Len() != 0 {
			msg := t.msgQueue.PopFront().(sdk.Msg)
			log.Printf("Packed msg info. Route: %s; Type: %s", msg.Route(), msg.Type())
			msgs = append(msgs, msg)
		}

		txBldr, err := utils.PrepareTxBuilder(t.TxBuilder, t.CliCtx)
		if err != nil {
			log.Printf("Transactor PrepareTxBuilder err: %v", err)
			continue
		}

		txBytes, err := txBldr.BuildAndSign(t.Key.GetName(), t.Passphrase, msgs)
		if err != nil {
			log.Printf("Transactor BuildAndSign err: %v", err)
			continue
		}

		tx, err := t.CliCtx.BroadcastTx(txBytes)
		if err != nil {
			log.Printf("Transactor BroadcastTx err: %v", err)
			continue
		}

		// Make sure the transaction has been mines
		log.Printf("Transactor broadcasted tx: %+v", tx)
		for try := 0; try < maxTry; try++ {
			if _, err = utils.QueryTx(t.CliCtx, tx.TxHash); err == nil {
				break
			}
			time.Sleep(time.Second)
		}

	}
}
