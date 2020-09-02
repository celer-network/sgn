package transactor

import (
	"fmt"
	"strings"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/seal"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gammazero/deque"
	"github.com/spf13/viper"
)

const (
	maxQueryRetry   = 15
	queryRetryDelay = 1 * time.Second
	maxSignRetry    = 10
	signRetryDelay  = 100 * time.Millisecond
)

type Transactor struct {
	TxBuilder  types.TxBuilder
	CliCtx     context.CLIContext
	Key        keys.Info
	passphrase string
	msgQueue   deque.Deque
	gpe        *GasPriceEstimator
}

func NewTransactor(cliHome, chainID, nodeURI, accAddr, passphrase string, cdc *codec.Codec, gpe *GasPriceEstimator) (*Transactor, error) {
	kb, err := keys.NewKeyringWithPassphrase(appName,
		viper.GetString(common.FlagSgnKeyringBackend), cliHome, passphrase)
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
			log.Errorln("kb.GetByAddress error:", err)
			return nil, err
		}
		if try != maxSignRetry-1 {
			log.Debugln("retry kb.GetByAddress due to error:", err)
			time.Sleep(signRetryDelay)
		}
	}

	fees, err := sdk.ParseCoins(viper.GetString(flags.FlagFees))
	if err != nil {
		panic(err)
	}

	gasPrices, err := sdk.ParseDecCoins(viper.GetString(flags.FlagGasPrices))
	if err != nil {
		panic(err)
	}

	txBldr := types.NewTxBuilder(
		utils.GetTxEncoder(cdc),
		viper.GetUint64(flags.FlagAccountNumber),
		viper.GetUint64(flags.FlagSequence),
		flags.GasFlagVar.Gas,
		flags.DefaultGasAdjustment,
		flags.GasFlagVar.Simulate,
		chainID,
		viper.GetString(flags.FlagMemo),
		fees,
		gasPrices)
	txBldr = txBldr.WithKeybase(kb)

	cliCtx := context.NewCLIContext().
		WithCodec(cdc).
		WithFromAddress(key.GetAddress()).
		WithFromName(key.GetName()).
		WithNodeURI(nodeURI).
		WithTrustNode(true).
		WithChainID(chainID).
		WithBroadcastMode(flags.BroadcastSync)

	transactor := &Transactor{
		TxBuilder:  txBldr,
		CliCtx:     cliCtx,
		Key:        key,
		passphrase: passphrase,
		gpe:        gpe,
	}

	return transactor, nil
}

func NewCliTransactor(cdc *codec.Codec, cliHome string) (*Transactor, error) {
	return NewTransactor(
		cliHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetString(common.FlagSgnValidatorAccount),
		viper.GetString(common.FlagSgnPassphrase),
		cdc,
		nil,
	)
}

func (t *Transactor) Run() {
	go t.start()
}

// AddTxMsg add msg into a queue before actual broadcast
func (t *Transactor) AddTxMsg(msg sdk.Msg) {
	t.msgQueue.PushBack(msg)
}

func (t *Transactor) SendTxMsg(msg sdk.Msg) (*sdk.TxResponse, error) {
	return t.SendTxMsgs([]sdk.Msg{msg})
}

func (t *Transactor) SendTxMsgs(msgs []sdk.Msg) (*sdk.TxResponse, error) {
	txBytes, stdSignMsg, err := t.signTx(msgs)
	if err != nil {
		return nil, fmt.Errorf("signTx err: %s", err)
	}
	txResponse, err := t.CliCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, fmt.Errorf("BroadcastTx err: %s", err)
	}

	if txResponse.Code != sdkerrors.SuccessABCICode {
		return &txResponse, fmt.Errorf("BroadcastTx failed with code: %d, rawLog: %s, stdSignMsg: %s acct %s accnum %d seq %d",
			txResponse.Code, txResponse.RawLog, t.Key.GetAddress(), stdSignMsg.ChainID, stdSignMsg.AccountNumber, stdSignMsg.Sequence)
	}

	return &txResponse, nil
}

// Poll tx queue and send msgs in batch
func (t *Transactor) start() {
	for {
		if t.msgQueue.Len() == 0 {
			time.Sleep(time.Second)
			continue
		}

		logEntry := seal.NewTransactorLog(t.Key.GetAddress().String())
		txResponse, err := t.bcastTxMsgQueue(logEntry)
		if err != nil {
			logEntry.Error = append(logEntry.Error, err.Error())
			seal.CommitTransactorLog(logEntry)
			continue
		}
		seal.CommitTransactorLog(logEntry)

		// Make sure the transaction has been mined
		t.WaitMined(txResponse.TxHash)
	}
}

func (t *Transactor) WaitMined(txHash string) (*sdk.TxResponse, error) {
	var err error
	mined := false
	var txResponse sdk.TxResponse
	for try := 0; try < maxQueryRetry; try++ {
		time.Sleep(queryRetryDelay)
		if txResponse, err = utils.QueryTx(t.CliCtx, txHash); err == nil {
			mined = true
			break
		}
	}
	if !mined {
		log.Errorf("Transaction %s not mined within %d retry, err %s", txHash, maxQueryRetry, err)
	} else if txResponse.Code != sdkerrors.SuccessABCICode {
		log.Errorf("Transaction %s failed with code %d, %s", txHash, txResponse.Code, txResponse.RawLog)
	} else {
		log.Debugf("Transaction %s succeeded", txHash)
	}
	return &txResponse, err
}

func (t *Transactor) bcastTxMsgQueue(logEntry *seal.TransactorLog) (*sdk.TxResponse, error) {
	logEntry.MsgNum = uint32(t.msgQueue.Len())
	var msgs []sdk.Msg
	for t.msgQueue.Len() != 0 {
		msg := t.msgQueue.PopFront().(sdk.Msg)
		logEntry.MsgType[msg.Type()] = logEntry.MsgType[msg.Type()] + 1
		msgs = append(msgs, msg)
	}
	txResponse, err := t.SendTxMsgs(msgs)
	if err != nil {
		return nil, err
	}
	logEntry.TxHash = txResponse.TxHash

	return txResponse, err
}

func (t *Transactor) signTx(msgs []sdk.Msg) ([]byte, *types.StdSignMsg, error) {
	if t.gpe != nil {
		t.TxBuilder = t.TxBuilder.WithGasPrices(t.gpe.GetGasPrice())
	}

	txBldr, err := utils.PrepareTxBuilder(t.TxBuilder, t.CliCtx)
	if err != nil {
		return nil, nil, fmt.Errorf("PrepareTxBuilder err: %s", err)
	}
	var txBytes []byte
	var stdSignMsg types.StdSignMsg
	for try := 0; try < maxSignRetry; try++ {
		stdSignMsg, err = txBldr.BuildSignMsg(msgs)
		if err != nil {
			return nil, nil, err
		}
		txBytes, err = txBldr.Sign(t.Key.GetName(), t.passphrase, stdSignMsg)
		if err == nil {
			return txBytes, &stdSignMsg, nil
		}
		if !strings.Contains(err.Error(), "resource temporarily unavailable") {
			break
		}
		if try != maxSignRetry-1 {
			log.Debugln("Failed to call txBldr.BuildAndSign. Will retry it.")
			time.Sleep(signRetryDelay)
		}
	}
	return nil, nil, fmt.Errorf("BuildAndSign err: %s", err)
}

func (t *Transactor) CliSendTxMsgWaitMined(msg sdk.Msg) {
	t.CliSendTxMsgsWaitMined([]sdk.Msg{msg})
}

func (t *Transactor) CliSendTxMsgsWaitMined(msgs []sdk.Msg) {
	txResponse, err := t.SendTxMsgs(msgs)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Transaction %s sent", txResponse.TxHash)
	res, err := t.WaitMined(txResponse.TxHash)
	t.CliCtx.PrintOutput(res)
}
