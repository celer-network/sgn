package transactor

import (
	"fmt"
	"strings"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/seal"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clientKeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gammazero/deque"
)

const (
	maxQueryRetry   = 20
	queryRetryDelay = 500 * time.Millisecond
	maxSignRetry    = 10
	signRetryDelay  = 100 * time.Millisecond
)

type Transactor struct {
	TxBuilder  types.TxBuilder
	CliCtx     context.CLIContext
	Key        keys.Info
	Passphrase string
	msgQueue   deque.Deque
}

func NewTransactor(cliHome, chainID, nodeURI, accAddr, passphrase, gasPrice string, cdc *codec.Codec) (*Transactor, error) {
	kb, err := clientKeys.NewKeyBaseFromDir(cliHome)
	if err != nil {
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(accAddr)
	if err != nil {
		return nil, err
	}

	// may run into "resource temporarily unavailable" error if directly run it
	// retry when get this issue to avoid failure.
	var key keys.Info
	for try := 0; try < maxSignRetry; try++ {
		key, err = kb.GetByAddress(addr)
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "resource temporarily unavailable") {
			return nil, err
		}
		if try != maxSignRetry-1 {
			log.Debugln("Failed to call kb.GetByAddress. Will retry it.")
			time.Sleep(signRetryDelay)
		}
	}

	txBldr := NewTxBuilder().
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID).
		WithKeybase(kb).
		WithGasPrices(gasPrice)

	cliCtx := context.
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

		logEntry := seal.NewTransactorLog(t.Key.GetAddress().String())
		tx, err := t.broadcastTx(logEntry)
		if err != nil {
			logEntry.Error = append(logEntry.Error, err.Error())
			seal.CommitTransactorLog(logEntry)
			continue
		}

		seal.CommitTransactorLog(logEntry)

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
			log.Errorf("Transaction %s not mined within %d retry, err %s", tx.TxHash, maxQueryRetry, err)
		} else {
			log.Debugf("Transaction %s has been mined", tx.TxHash)
		}
	}
}

func (t *Transactor) broadcastTx(logEntry *seal.TransactorLog) (*sdk.TxResponse, error) {
	logEntry.MsgNum = uint32(t.msgQueue.Len())
	var msgs []sdk.Msg
	for t.msgQueue.Len() != 0 {
		msg := t.msgQueue.PopFront().(sdk.Msg)
		logEntry.MsgType[msg.Type()] = logEntry.MsgType[msg.Type()] + 1
		msgs = append(msgs, msg)
	}

	txBytes, err := t.signTx(msgs)
	if err != nil {
		return nil, fmt.Errorf("signTx err: %s", err)
	}
	tx, err := t.CliCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, fmt.Errorf("BroadcastTx err: %s", err)
	}
	logEntry.TxHash = tx.TxHash

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
			log.Debugln("Failed to call txBldr.BuildAndSign. Will retry it.")
			time.Sleep(signRetryDelay)
		}
	}
	return nil, fmt.Errorf("BuildAndSign err: %s", err)
}
