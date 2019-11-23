package transactor

import (
	"fmt"
	"strings"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/seal"
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

const (
	maxQueryRetry   = 20
	queryRetryDelay = 500 * time.Millisecond
	maxSignRetry    = 10
	signRetryDelay  = 10 * time.Millisecond
)

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

// AddTxMsg add msg into a queue before actual broadcast
func (t *Transactor) AddTxMsg(msg sdk.Msg) {
	t.msgQueue.PushBack(msg)
}

// Poll tx queue and send msgs in batch
func (t *Transactor) start() {
	for {
		if t.msgQueue.Len() == 0 {
			time.Sleep(time.Second)
			continue
		}

		txlog := seal.NewTransactorLog()
		tx, err := t.broadcastTx(txlog)
		if err != nil {
			txlog.Error = append(txlog.Error, err.Error())
			seal.CommitTransactorLog(txlog)
			continue
		}

		seal.CommitTransactorLog(txlog)

		// Make sure the transaction has been mined
		success := false
		for try := 0; try < maxQueryRetry; try++ {
			time.Sleep(queryRetryDelay)
			if _, err = utils.QueryTx(t.CliCtx, tx.TxHash); err == nil {
				success = true
				break
			}
		}
		if !success {
			log.Errorf("Transaction %s not mined within %d retry", tx.TxHash, maxQueryRetry)
		} else {
			log.Debugf("Transaction %s has been mined", tx.TxHash)
		}
	}
}

func (t *Transactor) broadcastTx(txlog *seal.TransactorLog) (*sdk.TxResponse, error) {
	txlog.MsgNum = uint32(t.msgQueue.Len())
	var msgs []sdk.Msg
	for t.msgQueue.Len() != 0 {
		msg := t.msgQueue.PopFront().(sdk.Msg)
		seal.AddTransactorMsg(txlog, msg.Type())
		msgs = append(msgs, msg)
	}

	txBytes, err := t.signTx(msgs)
	tx, err := t.CliCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, fmt.Errorf("BroadcastTx err: %s", err)
	}
	txlog.TxHash = tx.TxHash

	return &tx, nil
}

func (t *Transactor) signTx(msgs []sdk.Msg) ([]byte, error) {
	txBldr, err := utils.PrepareTxBuilder(t.TxBuilder, t.CliCtx)
	if err != nil {
		return nil, fmt.Errorf("PrepareTxBuilder err: %s", err)
	}
	var txBytes []byte
	for try := 0; try < maxSignRetry; try++ {
		txBytes, err = txBldr.BuildAndSign(t.Key.GetName(), t.Passphrase, msgs)
		if err == nil {
			return txBytes, nil
		}
		if !strings.Contains(err.Error(), "resource temporarily unavailable") {
			break
		}
		if try != maxSignRetry-1 {
			time.Sleep(signRetryDelay)
		}
	}
	return nil, fmt.Errorf("BuildAndSign err: %s", err)
}
