package monitor

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"
	"strconv"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/monitor/watcher"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/client"
	tTypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	pollingInterval = 1
)

var (
	initiateWithdrawRewardEvent = fmt.Sprintf("%s.%s='%s'", validator.ModuleName, sdk.AttributeKeyAction, validator.ActionInitiateWithdraw)
	slashEvent                  = fmt.Sprintf("%s.%s='%s'", slash.EventTypeSlash, sdk.AttributeKeyAction, slash.ActionPenalty)
	submitChangeEvent           = fmt.Sprintf("%s.%s='%s'", sync.EventTypeSync, sdk.AttributeKeyAction, sync.ActionSubmitChange)
)

type EthMonitor struct {
	ethClient      *mainchain.EthClient
	operator       *transactor.Transactor
	db             *dbm.GoLevelDB
	ms             *watcher.Service
	dposContract   *watcher.BoundContract
	sgnContract    *watcher.BoundContract
	ledgerContract *watcher.BoundContract
	blkNum         *big.Int
	secureBlkNum   uint64
	isValidator    bool
}

func NewEthMonitor(ethClient *mainchain.EthClient, operator *transactor.Transactor) {
	dataDir := filepath.Join(viper.GetString(flags.FlagHome), "data")
	db, err := dbm.NewGoLevelDB("monitor", dataDir)
	if err != nil {
		log.Fatalln("New monitor db err", err)
	}

	st, err := watcher.NewKVStoreLocal(filepath.Join(dataDir, "watch"), false)
	if err != nil {
		log.Fatalln("New watch db err", err)
	}

	dal := watcher.NewDAL(st)
	ws := watcher.NewWatchService(ethClient.Client, dal, viper.GetUint64(common.FlagEthPollInterval))
	if ws == nil {
		log.Fatalln("Cannot create watch service")
	}

	ms := watcher.NewService(ws, 0 /* blockDelay */, true /* enabled */, "" /* rpcAddr */)
	ms.Init()

	dposCandidateInfo, err := ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, ethClient.Address)
	if err != nil {
		log.Fatalln("GetCandidateInfo err", err)
	}

	dposContract, err := watcher.NewBoundContract(ethClient.Client, ethClient.DPoSAddress, mainchain.DPoSABI)
	if err != nil {
		log.Fatalln("dposContract err", err)
	}

	sgnContract, err := watcher.NewBoundContract(ethClient.Client, ethClient.SGNAddress, mainchain.SGNABI)
	if err != nil {
		log.Fatalln("sgnContract err", err)
	}

	ledgerContract, err := watcher.NewBoundContract(ethClient.Client, ethClient.LedgerAddress, mainchain.CelerLedgerABI)
	if err != nil {
		log.Fatalln("ledgerContract err", err)
	}

	m := EthMonitor{
		ethClient:      ethClient,
		operator:       operator,
		db:             db,
		ms:             ms,
		blkNum:         ms.GetCurrentBlockNumber(),
		dposContract:   dposContract,
		sgnContract:    sgnContract,
		ledgerContract: ledgerContract,
		isValidator:    mainchain.IsBonded(dposCandidateInfo),
	}

	go m.monitorBlockHead()
	go m.monitorUpdateSidechainAddr()
	go m.monitorDelegate()
	go m.monitorValidatorChange()
	go m.monitorIntendWithdraw()
	go m.monitorIntendSettle()
	go m.monitorWithdrawReward()
	go m.monitorSlash()
	go m.monitorSubmitChange()
}

func (m *EthMonitor) monitorBlockHead() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C
		blkNum := m.ms.GetCurrentBlockNumber()
		if blkNum.Cmp(m.blkNum) == 0 {
			continue
		}

		m.blkNum = blkNum
		m.secureBlkNum = blkNum.Uint64() - viper.GetUint64(common.FlagEthConfirmCount)
		m.processQueue()
	}
}

func (m *EthMonitor) monitorUpdateSidechainAddr() {
	_, err := m.ms.Monitor(string(UpdateSidechainAddr), m.sgnContract, m.blkNum, nil, false, func(cb watcher.CallbackID, eLog ethtypes.Log) {
		log.Infof("Catch event UpdateSidechainAddr, tx hash: %x", eLog.TxHash)
		event := NewEvent(UpdateSidechainAddr, eLog)
		m.db.Set(GetPullerKey(eLog), event.MustMarshal())
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *EthMonitor) monitorDelegate() {
	_, err := m.ms.Monitor(string(Delegate), m.dposContract, m.blkNum, nil, false, func(cb watcher.CallbackID, eLog ethtypes.Log) {
		log.Infof("Catch event Delegate, tx hash: %x", eLog.TxHash)
		event := NewEvent(Delegate, eLog)
		m.db.Set(GetEventKey(eLog), event.MustMarshal())
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *EthMonitor) monitorValidatorChange() {
	_, err := m.ms.Monitor(string(ValidatorChange), m.dposContract, m.blkNum, nil, false, func(cb watcher.CallbackID, eLog ethtypes.Log) {
		log.Infof("Catch event ValidatorChange, tx hash: %x", eLog.TxHash)
		event := NewEvent(ValidatorChange, eLog)
		m.db.Set(GetEventKey(eLog), event.MustMarshal())
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *EthMonitor) monitorIntendWithdraw() {
	_, err := m.ms.Monitor(string(IntendWithdraw), m.dposContract, m.blkNum, nil, false, func(cb watcher.CallbackID, eLog ethtypes.Log) {
		log.Infof("Catch event IntendWithdraw, tx hash: %x", eLog.TxHash)
		event := NewEvent(IntendWithdraw, eLog)
		m.db.Set(GetEventKey(eLog), event.MustMarshal())
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *EthMonitor) monitorIntendSettle() {
	_, err := m.ms.Monitor(string(IntendSettle), m.ledgerContract, m.blkNum, nil, false, func(cb watcher.CallbackID, eLog ethtypes.Log) {
		log.Infof("Catch event IntendSettle, tx hash: %x", eLog.TxHash)
		event := NewEvent(IntendSettle, eLog)
		m.db.Set(GetPullerKey(eLog), event.MustMarshal())
		m.db.Set(GetPusherKey(eLog), event.MustMarshal())
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *EthMonitor) monitorWithdrawReward() {
	m.monitorTendermintEvent(initiateWithdrawRewardEvent, func(e abci.Event) {
		if !m.isValidator {
			return
		}

		event := sdk.StringifyEvent(e)
		if event.Attributes[0].Value == validator.ActionInitiateWithdraw {
			m.handleInitiateWithdrawReward(event.Attributes[1].Value)
		}
	})
}

func (m *EthMonitor) monitorSlash() {
	m.monitorTendermintEvent(slashEvent, func(e abci.Event) {
		if !m.isValidator {
			return
		}

		event := sdk.StringifyEvent(e)

		if event.Attributes[0].Value == slash.ActionPenalty {
			nonce, err := strconv.ParseUint(event.Attributes[1].Value, 10, 64)
			if err != nil {
				log.Errorln("Parse penalty nonce error", err)
				return
			}

			penaltyEvent := NewPenaltyEvent(nonce)
			m.handlePenalty(penaltyEvent)
			m.db.Set(GetPenaltyKey(penaltyEvent.Nonce), penaltyEvent.MustMarshal())
		}
	})
}

func (m *EthMonitor) monitorSubmitChange() {
	m.monitorTendermintEvent(submitChangeEvent, func(e abci.Event) {
		if !m.isValidator && !viper.GetBool(common.FlagSgnBootNode) {
			return
		}

		event := sdk.StringifyEvent(e)
		if event.Type == sync.EventTypeSync && event.Attributes[0].Value == sync.ActionSubmitChange {
			changeId, err := strconv.ParseUint(event.Attributes[1].Value, 10, 64)
			if err != nil {
				log.Errorln("Parse changeId error", err)
				return
			}

			change, err := sync.CLIQueryChange(m.operator.CliCtx, sync.RouterKey, changeId)
			if err != nil {
				log.Errorln("Query change error", err)
				return
			}

			log.Infoln("Verify change", change)
			if m.verifyChange(change) {
				msg := sync.NewMsgApprove(changeId, m.operator.Key.GetAddress())
				m.operator.AddTxMsg(msg)
			}
		}
	})
}

func (m *EthMonitor) monitorTendermintEvent(eventTag string, handleEvent func(event abci.Event)) {
	client, err := client.NewHTTP(m.operator.CliCtx.NodeURI, "/websocket")
	if err != nil {
		log.Errorln("Fail to start create http client", err)
		return
	}

	err = client.Start()
	if err != nil {
		log.Errorln("Fail to start ws client", err)
		return
	}
	defer client.Stop()

	txs, err := client.Subscribe(context.Background(), "monitor", eventTag)
	if err != nil {
		log.Errorln("ws client subscribe error", err)
		return
	}

	for e := range txs {
		switch data := e.Data.(type) {
		case tTypes.EventDataNewBlock:
			for _, event := range data.ResultBeginBlock.Events {
				handleEvent(event)
			}
			for _, event := range data.ResultEndBlock.Events {
				handleEvent(event)
			}
		case tTypes.EventDataTx:
			for _, event := range data.TxResult.Result.Events {
				handleEvent(event)
			}
		}
	}
}
