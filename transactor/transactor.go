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
	maxTxRetry     = 15
	txRetryDelay   = 1 * time.Second
	maxSignRetry   = 10
	signRetryDelay = 100 * time.Millisecond
	maxGasRetry    = 5
)

var gasErr = fmt.Errorf("out of gas")

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

	gasAdjustment := viper.GetFloat64(common.FlagSgnGasAdjustment)
	if gasAdjustment == 0 {
		gasAdjustment = common.DefaultSgnGasAdjustment
	}
	txBldr := types.NewTxBuilder(
		utils.GetTxEncoder(cdc),
		viper.GetUint64(flags.FlagAccountNumber),
		viper.GetUint64(flags.FlagSequence),
		common.DefaultSgnGasLimit,
		gasAdjustment,
		true,
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

// Poll tx queue and send msgs in batch
func (t *Transactor) start() {
	for {
		if t.msgQueue.Len() == 0 {
			time.Sleep(time.Second)
			continue
		}

		t.drainTxMsgQueue()
	}
}

func (t *Transactor) drainTxMsgQueue() {
	logEntry := seal.NewTransactorLog(t.Key.GetAddress().String())

	var msgs []sdk.Msg
	for t.msgQueue.Len() != 0 {
		msg := t.msgQueue.PopFront().(sdk.Msg)
		logEntry.MsgType[msg.Type()] = logEntry.MsgType[msg.Type()] + 1
		msgs = append(msgs, msg)
	}
	logEntry.MsgNum = uint32(len(msgs))

	txResponse, err := t.sendTxMsgs(msgs)
	logEntry.TxHash = txResponse.TxHash
	if err != nil {
		logEntry.Error = append(logEntry.Error, err.Error())
		logEntry.Status = seal.TxMsgStatus_FAILED
		seal.CommitTransactorLog(logEntry)
		return
	}
	logEntry.Status = seal.TxMsgStatus_SENT
	seal.CommitTransactorLog(logEntry)

	// wait till transaction is mined or failed
	_, err = t.waitMined(txResponse.TxHash)
	if err != nil {
		logEntry.Error = append(logEntry.Error, err.Error())
		logEntry.Status = seal.TxMsgStatus_FAILED
	} else {
		logEntry.Status = seal.TxMsgStatus_SUCCEED
	}
	seal.CommitTransactorLog(logEntry)
}

func (t *Transactor) sendTxMsgs(msgs []sdk.Msg) (*sdk.TxResponse, error) {
	var txResponseErr error
	for try := 0; try < maxTxRetry; try++ {
		txBytes, stdSignMsg, err := t.buildAndSignTx(msgs)
		if err != nil {
			return nil, fmt.Errorf("buildAndSignTx err: %w", err)
		}
		txResponse, err := t.CliCtx.BroadcastTx(txBytes)
		if err != nil {
			return nil, fmt.Errorf("BroadcastTx err: %w", err)
		}

		if txResponse.Code == sdkerrors.SuccessABCICode {
			return &txResponse, nil
		}

		txResponseErr = fmt.Errorf("BroadcastTx failed with code: %d, rawLog: %s, stdSignMsg chainId: %s acct: %s accnum: %d seq: %d",
			txResponse.Code, txResponse.RawLog, stdSignMsg.ChainID, t.Key.GetAddress(), stdSignMsg.AccountNumber, stdSignMsg.Sequence)
		if txResponse.Code == sdkerrors.ErrUnauthorized.ABCICode() {
			log.Warnln(txResponseErr.Error(), "retrying")
			time.Sleep(txRetryDelay)
		} else {
			return &txResponse, txResponseErr
		}
	}
	return nil, txResponseErr
}

func (t *Transactor) buildAndSignTx(msgs []sdk.Msg) ([]byte, *types.StdSignMsg, error) {
	if t.gpe != nil {
		t.TxBuilder = t.TxBuilder.WithGasPrices(t.gpe.GetGasPrice())
	}

	txBldr, err := utils.PrepareTxBuilder(t.TxBuilder, t.CliCtx)
	if err != nil {
		return nil, nil, fmt.Errorf("PrepareTxBuilder err: %w", err)
	}

	if txBldr.SimulateAndExecute() || t.CliCtx.Simulate {
		txBldr, err = utils.EnrichWithGas(txBldr, t.CliCtx, msgs)
		if err != nil {
			return nil, nil, fmt.Errorf("EnrichWithGas err: %w", err)
		}
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
	return nil, nil, fmt.Errorf("TxBuilder Sign err: %w", err)
}

func (t *Transactor) waitMined(txHash string) (*sdk.TxResponse, error) {
	var err error
	mined := false
	var txResponse sdk.TxResponse
	for try := 0; try < maxTxRetry; try++ {
		time.Sleep(txRetryDelay)
		if txResponse, err = utils.QueryTx(t.CliCtx, txHash); err == nil {
			mined = true
			break
		}
	}
	if !mined {
		return &txResponse, fmt.Errorf("tx not mined, err: %w", err)
	} else if txResponse.Code != sdkerrors.SuccessABCICode {
		log.Errorf("Transaction %s failed with code %d, %s", txHash, txResponse.Code, txResponse.RawLog)
		if txResponse.Code == 11 { // out of gas
			return &txResponse, gasErr
		} else {
			return &txResponse, fmt.Errorf("tx failed with code %d, %s", txResponse.Code, txResponse.RawLog)
		}
	}
	return &txResponse, nil
}

func (t *Transactor) CliSendTxMsgWaitMined(msg sdk.Msg) {
	t.CliSendTxMsgsWaitMined([]sdk.Msg{msg})
}

func (t *Transactor) CliSendTxMsgsWaitMined(msgs []sdk.Msg) {
	txResponse, err := t.sendTxMsgs(msgs)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Transaction %s sent", txResponse.TxHash)
	res, err := t.waitMined(txResponse.TxHash)
	if err != nil {
		log.Error(err)
	}
	t.CliCtx.PrintOutput(res)
}
