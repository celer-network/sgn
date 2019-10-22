package monitor

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/allegro/bigcache"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authUtils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gammazero/deque"
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
	transactor  *utils.Transactor
	cdc         *codec.Codec
	pusherQueue deque.Deque
	pullerQueue deque.Deque
	eventQueue  deque.Deque
	txMemo      *bigcache.BigCache
	pubkey      string
	isValidator bool
}

func NewEthMonitor(ethClient *mainchain.EthClient, transactor *utils.Transactor, cdc *codec.Codec, pubkey string) {
	txMemo, err := bigcache.NewBigCache(bigcache.DefaultConfig(24 * time.Hour))
	if err != nil {
		log.Fatalf("NewBigCache err", err)
	}

	candidateInfo, err := ethClient.Guard.GetCandidateInfo(&bind.CallOpts{}, ethClient.Address)
	if err != nil {
		log.Fatalf("GetCandidateInfo err", err)
	}

	m := EthMonitor{
		ethClient:   ethClient,
		transactor:  transactor,
		cdc:         cdc,
		txMemo:      txMemo,
		pubkey:      pubkey,
		isValidator: candidateInfo.IsVldt,
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
		log.Printf("SubscribeNewHead err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("SubscribeNewHead err", err)
		case header := <-headerChan:
			m.handleNewBlock(header)
			go m.processQueue()
		}
	}
}

func (m *EthMonitor) monitorInitializeCandidate() {
	initializeCandidateChan := make(chan *mainchain.GuardInitializeCandidate)
	sub, err := m.ethClient.Guard.WatchInitializeCandidate(nil, initializeCandidateChan, nil, nil)
	if err != nil {
		log.Printf("WatchInitializeCandidate err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("WatchInitializeCandidate err", err)
		case initializeCandiate := <-initializeCandiateChan:
			m.eventQueue.PushBack(NewEvent(initializeCandiate, initializeCandiate.Raw))
		}
	}
}

func (m *EthMonitor) monitorDelegate() {
	delegateChan := make(chan *mainchain.GuardDelegate)

	sub, err := m.ethClient.Guard.WatchDelegate(nil, delegateChan, nil, []ethcommon.Address{m.ethClient.Address})
	if err != nil {
		log.Printf("WatchDelegate err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("WatchDelegate err", err)
		case delegate := <-delegateChan:
			m.eventQueue.PushBack(NewEvent(delegate, delegate.Raw))
		}
	}
}

func (m *EthMonitor) monitorValidatorChange() {
	validatorChangeChan := make(chan *mainchain.GuardValidatorChange)
	sub, err := m.ethClient.Guard.WatchValidatorChange(nil, validatorChangeChan, nil, nil)
	if err != nil {
		log.Printf("WatchValidatorChange err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("WatchValidatorChange err", err)
		case validatorChange := <-validatorChangeChan:
			m.eventQueue.PushBack(NewEvent(validatorChange, validatorChange.Raw))
		}
	}
}

func (m *EthMonitor) monitorIntendWithdraw() {
	intendWithdrawChan := make(chan *mainchain.GuardIntendWithdraw)
	sub, err := m.ethClient.Guard.WatchIntendWithdraw(nil, intendWithdrawChan, nil, nil)
	if err != nil {
		log.Printf("WatchIntendWithdraw err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("WatchIntendWithdraw err", err)
		case intendWithdraw := <-intendWithdrawChan:
			m.eventQueue.PushBack(NewEvent(intendWithdraw, intendWithdraw.Raw))
		}
	}
}

func (m *EthMonitor) monitorIntendSettle() {
	intendSettleChan := make(chan *mainchain.CelerLedgerIntendSettle)
	sub, err := m.ethClient.Ledger.WatchIntendSettle(nil, intendSettleChan, nil)
	if err != nil {
		log.Printf("WatchIntendSettle err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("WatchIntendSettle err", err)
		case intendSettle := <-intendSettleChan:
			m.eventQueue.PushBack(NewEvent(intendSettle, intendSettle.Raw))
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
				log.Printf("Parse penalty nonce error", err)
				return
			}

			m.handlePenalty(nonce)

			penaltyEvent := NewPenaltyEvent(nonce)
			m.pusherQueue.PushBack(penaltyEvent)
		}
	})
}

func (m *EthMonitor) monitorTendermintEvent(eventTag string, handleEvent func(event sdk.StringEvent)) {
	for {
		for page := 1; ; page++ {
			hasSeenEvent := false
			txs, err := authUtils.QueryTxsByEvents(m.transactor.CliCtx, []string{eventTag}, page, txsPageLimit)
			if err != nil {
				log.Printf("QueryTxsByEvents err", err)
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
