package monitor

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/monitor/watcher"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/client"
	tTypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	txsPageLimit    = 30
	txsPullInterval = 5 // interval in seconds of pulling sidechain events
	pollingInterval = 10
)

var (
	initiateWithdrawRewardEvent = fmt.Sprintf("%s.%s='%s'", validator.ModuleName, sdk.AttributeKeyAction, validator.ActionInitiateWithdraw)
	slashEvent                  = fmt.Sprintf("%s.%s='%s'", slash.EventTypeSlash, sdk.AttributeKeyAction, slash.ActionPenalty)
)

type EthMonitor struct {
	ethClient      *mainchain.EthClient
	transactor     *transactor.Transactor
	db             *dbm.GoLevelDB
	ms             *watcher.Service
	guardContract  *watcher.BoundContract
	ledgerContract *watcher.BoundContract
	pubkey         string
	transactors    []string
	isValidator    bool
}

func NewEthMonitor(ethClient *mainchain.EthClient, transactor *transactor.Transactor, pubkey string, transactors []string) {
	dataDir := filepath.Join(viper.GetString(flags.FlagHome), "data")
	db, err := dbm.NewGoLevelDB("monitor", dataDir)
	if err != nil {
		log.Fatalln("New monitor db err", err)
	}

	st, err := watcher.NewKVStoreSQL("sqlite3", filepath.Join(dataDir, "watch.db"))
	if err != nil {
		log.Fatalln("New watch db err", err)
	}

	dal := watcher.NewDAL(st)
	ws := watcher.NewWatchService(ethClient.Client, dal, pollingInterval)
	if ws == nil {
		log.Fatalln("Cannot create watch service")
	}

	ms := watcher.NewService(ws, 0 /* blockDelay */, true /* enabled */, "" /* rpcAddr */)
	ms.Init()

	candidateInfo, err := ethClient.Guard.GetCandidateInfo(&bind.CallOpts{}, ethClient.Address)
	if err != nil {
		log.Fatalln("GetCandidateInfo err", err)
	}

	guardContract, err := watcher.NewBoundContract(ethClient.Client, ethClient.GuardAddress, mainchain.GuardABI)
	if err != nil {
		log.Fatalln("guardContract err", err)
	}

	ledgerContract, err := watcher.NewBoundContract(ethClient.Client, ethClient.LedgerAddress, mainchain.CelerLedgerABI)
	if err != nil {
		log.Fatalln("ledgerContract err", err)
	}

	m := EthMonitor{
		ethClient:      ethClient,
		transactor:     transactor,
		db:             db,
		ms:             ms,
		guardContract:  guardContract,
		ledgerContract: ledgerContract,
		pubkey:         pubkey,
		transactors:    transactors,
		isValidator:    mainchain.IsBonded(candidateInfo),
	}

	go m.monitorBlockHead()
	go m.monitorInitializeCandidate()
	go m.monitorDelegate()
	go m.monitorValidatorChange()
	go m.monitorIntendWithdraw()
	go m.monitorIntendSettle()
	go m.monitorWithdrawReward()
	go m.monitorSlash()
}

func (m *EthMonitor) monitorBlockHead() {
	headerChan := make(chan *ethTypes.Header)
	sub, err := m.ethClient.Client.SubscribeNewHead(context.Background(), headerChan)
	if err != nil {
		log.Errorln("SubscribeNewHead err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Errorln("SubscribeNewHead err", err)
		case header := <-headerChan:
			go m.handleNewBlock(header)
			go m.processQueue()
		}
	}
}

func (m *EthMonitor) monitorInitializeCandidate() {
	m.ms.Monitor(string(InitializeCandidate), m.guardContract, nil, nil, false, func(cb watcher.CallbackID, eLog ethtypes.Log) {
		event := NewEvent(InitializeCandidate, eLog)
		m.db.Set(GetEventKey(eLog), event.MustMarshal())
		log.Infof("Catch event InitializeCandidate, tx hash: %v", eLog.TxHash)
	})
}

func (m *EthMonitor) monitorDelegate() {
	m.ms.Monitor(string(Delegate), m.guardContract, nil, nil, false, func(cb watcher.CallbackID, eLog ethtypes.Log) {
		event := NewEvent(Delegate, eLog)
		m.db.Set(GetEventKey(eLog), event.MustMarshal())
		log.Infof("Catch event Delegate, tx hash: %v", eLog.TxHash)
	})
}

func (m *EthMonitor) monitorValidatorChange() {
	m.ms.Monitor(string(ValidatorChange), m.guardContract, nil, nil, false, func(cb watcher.CallbackID, eLog ethtypes.Log) {
		event := NewEvent(ValidatorChange, eLog)
		m.db.Set(GetEventKey(eLog), event.MustMarshal())
		log.Infof("Catch event ValidatorChange, tx hash: %v", eLog.TxHash)
	})
}

func (m *EthMonitor) monitorIntendWithdraw() {
	m.ms.Monitor(string(IntendWithdraw), m.guardContract, nil, nil, false, func(cb watcher.CallbackID, eLog ethtypes.Log) {
		event := NewEvent(IntendWithdraw, eLog)
		m.db.Set(GetEventKey(eLog), event.MustMarshal())
		log.Infof("Catch event IntendWithdraw, tx hash: %v", eLog.TxHash)
	})
}

func (m *EthMonitor) monitorIntendSettle() {
	m.ms.Monitor(string(IntendSettle), m.ledgerContract, nil, nil, false, func(cb watcher.CallbackID, eLog ethtypes.Log) {
		event := NewEvent(IntendSettle, eLog)
		m.db.Set(GetEventKey(eLog), event.MustMarshal())
		log.Infof("Catch event IntendSettle, tx hash: %v", eLog.TxHash)
	})
}

func (m *EthMonitor) monitorWithdrawReward() {
	m.monitorTendermintEvent(initiateWithdrawRewardEvent, func(e abci.Event) {
		event := sdk.StringifyEvent(e)
		if event.Attributes[0].Value == validator.ActionInitiateWithdraw {
			m.handleInitiateWithdrawReward(event.Attributes[1].Value)
		}
	})
}

func (m *EthMonitor) monitorSlash() {
	m.monitorTendermintEvent(slashEvent, func(e abci.Event) {
		event := sdk.StringifyEvent(e)
		if event.Attributes[0].Value == slash.ActionPenalty {
			nonce, err := strconv.ParseUint(event.Attributes[1].Value, 10, 64)
			if err != nil {
				log.Errorln("Parse penalty nonce error", err)
				return
			}

			m.handlePenalty(nonce)

			penaltyEvent := NewPenaltyEvent(nonce)
			m.db.Set(GetPenaltyKey(penaltyEvent.nonce), penaltyEvent.MustMarshal())
		}
	})
}

func (m *EthMonitor) monitorTendermintEvent(eventTag string, handleEvent func(event abci.Event)) {
	client := client.NewHTTP(m.transactor.CliCtx.NodeURI, "/websocket")
	err := client.Start()
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
		case tTypes.EventDataTx:
			for _, event := range data.TxResult.Result.Events {
				handleEvent(event)
			}
		}
	}
}
