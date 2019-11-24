package monitor

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/allegro/bigcache"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authUtils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	txsPageLimit = 30
)

var (
	initiateWithdrawRewardEvent = fmt.Sprintf("%s.%s='%s'", validator.ModuleName, sdk.AttributeKeyAction, validator.ActionInitiateWithdraw)
	slashEvent                  = fmt.Sprintf("%s.%s='%s'", slash.EventTypeSlash, sdk.AttributeKeyAction, slash.ActionPenalty)
)

type EthMonitor struct {
	ethClient   *mainchain.EthClient
	transactor  *transactor.Transactor
	db          *dbm.GoLevelDB
	txMemo      *bigcache.BigCache
	pubkey      string
	transactors []string
	isValidator bool
}

func NewEthMonitor(ethClient *mainchain.EthClient, transactor *transactor.Transactor, db *dbm.GoLevelDB, pubkey string, transactors []string) {
	txMemo, err := bigcache.NewBigCache(bigcache.DefaultConfig(24 * time.Hour))
	if err != nil {
		log.Fatalln("NewBigCache err", err)
	}

	candidateInfo, err := ethClient.Guard.GetCandidateInfo(&bind.CallOpts{}, ethClient.Address)
	if err != nil {
		log.Fatalln("GetCandidateInfo err", err)
	}

	m := EthMonitor{
		ethClient:   ethClient,
		transactor:  transactor,
		db:          db,
		txMemo:      txMemo,
		pubkey:      pubkey,
		transactors: transactors,
		isValidator: mainchain.IsBonded(candidateInfo),
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
	headerChan := make(chan *types.Header)
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
			m.handleNewBlock(header)
			go m.processQueue()
		}
	}
}

func (m *EthMonitor) monitorInitializeCandidate() {
	initializeCandidateChan := make(chan *mainchain.GuardInitializeCandidate)
	sub, err := m.ethClient.Guard.WatchInitializeCandidate(nil, initializeCandidateChan, nil)
	if err != nil {
		log.Errorln("WatchInitializeCandidate err:", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Errorln("WatchInitializeCandidate err: ", err)
		case initializeCandidate := <-initializeCandidateChan:
			event := NewEvent(InitializeCandidate, initializeCandidate.Raw)
			m.db.Set(GetEventKey(initializeCandidate.Raw), event.MustMarshal())
			log.Infof("Catch event InitializeCandidate, candidate %x, min self stake %s, sidechain addr %x",
				initializeCandidate.Candidate, initializeCandidate.MinSelfStake.String(), initializeCandidate.SidechainAddr)
		}
	}
}

func (m *EthMonitor) monitorDelegate() {
	delegateChan := make(chan *mainchain.GuardDelegate)

	sub, err := m.ethClient.Guard.WatchDelegate(nil, delegateChan, nil, []mainchain.Addr{m.ethClient.Address})
	if err != nil {
		log.Errorln("WatchDelegate err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Errorln("WatchDelegate err", err)
		case delegate := <-delegateChan:
			event := NewEvent(Delegate, delegate.Raw)
			m.db.Set(GetEventKey(delegate.Raw), event.MustMarshal())
			log.Infof("Catch event GuardDelegate, delegator %x, candidate %x, new stake %s, pool %s",
				delegate.Delegator, delegate.Candidate, delegate.NewStake.String(), delegate.StakingPool.String())
		}
	}
}

func (m *EthMonitor) monitorValidatorChange() {
	validatorChangeChan := make(chan *mainchain.GuardValidatorChange)
	sub, err := m.ethClient.Guard.WatchValidatorChange(nil, validatorChangeChan, nil, nil)
	if err != nil {
		log.Errorln("WatchValidatorChange err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Errorln("WatchValidatorChange err", err)
		case validatorChange := <-validatorChangeChan:
			event := NewEvent(ValidatorChange, validatorChange.Raw)
			m.db.Set(GetEventKey(validatorChange.Raw), event.MustMarshal())
			log.Infof("Catch event GuardValidatorChange, addr %x, change type %d",
				validatorChange.EthAddr, validatorChange.ChangeType)
		}
	}
}

func (m *EthMonitor) monitorIntendWithdraw() {
	intendWithdrawChan := make(chan *mainchain.GuardIntendWithdraw)
	sub, err := m.ethClient.Guard.WatchIntendWithdraw(nil, intendWithdrawChan, nil, nil)
	if err != nil {
		log.Errorln("WatchIntendWithdraw err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Errorln("WatchIntendWithdraw err", err)
		case intendWithdraw := <-intendWithdrawChan:
			event := NewEvent(IntendWithdraw, intendWithdraw.Raw)
			m.db.Set(GetEventKey(intendWithdraw.Raw), event.MustMarshal())
			log.Infof("Catch event GuardIntendWithdraw, delegator %x, candidate %x, withdraw %s, time %s",
				intendWithdraw.Delegator, intendWithdraw.Candidate, intendWithdraw.WithdrawAmount.String(), intendWithdraw.ProposedTime.String())
		}
	}
}

func (m *EthMonitor) monitorIntendSettle() {
	intendSettleChan := make(chan *mainchain.CelerLedgerIntendSettle)
	sub, err := m.ethClient.Ledger.WatchIntendSettle(nil, intendSettleChan, nil)
	if err != nil {
		log.Errorln("WatchIntendSettle err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Errorln("WatchIntendSettle err", err)
		case intendSettle := <-intendSettleChan:
			event := NewEvent(IntendSettle, intendSettle.Raw)
			m.db.Set(GetEventKey(intendSettle.Raw), event.MustMarshal())
			log.Infof("Catch event CelerLedgerIntendSettle, channel ID %x, seqnums %s %s",
				intendSettle.ChannelId, intendSettle.SeqNums[0].String(), intendSettle.SeqNums[1].String())
		}
	}
}

func (m *EthMonitor) monitorWithdrawReward() {
	m.monitorTendermintEvent(initiateWithdrawRewardEvent, func(event sdk.StringEvent) {
		if event.Attributes[0].Value == validator.ActionInitiateWithdraw {
			m.handleInitiateWithdrawReward(event.Attributes[1].Value)
		}
	})
}

func (m *EthMonitor) monitorSlash() {
	m.monitorTendermintEvent(slashEvent, func(event sdk.StringEvent) {
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

func (m *EthMonitor) monitorTendermintEvent(eventTag string, handleEvent func(event sdk.StringEvent)) {
	for {
		for page := 1; ; page++ {
			hasSeenEvent := false
			txs, err := authUtils.QueryTxsByEvents(m.transactor.CliCtx, []string{eventTag}, page, txsPageLimit)
			if err != nil {
				log.Errorln("QueryTxsByEvents err", err)
				break
			}

			for _, tx := range txs.Txs {
				// Check if the tx has been seen before
				_, err = m.txMemo.Get(tx.TxHash)
				if err == nil {
					hasSeenEvent = true
					continue
				}

				m.txMemo.Set(tx.TxHash, []byte{1})
				for _, event := range tx.Events {
					handleEvent(event)
				}
			}

			// Check if it is necessary to query next page
			if txs.PageNumber >= txs.PageTotal || hasSeenEvent {
				break
			}
		}

		time.Sleep(30 * time.Second)
	}
}
