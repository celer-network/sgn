package monitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/allegro/bigcache"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
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

	candiateInfo, err := ethClient.Guard.GetCandidateInfo(&bind.CallOpts{}, ethClient.Address)
	if err != nil {
		log.Fatalf("GetCandidateInfo err", err)
	}

	m := EthMonitor{
		ethClient:   ethClient,
		transactor:  transactor,
		cdc:         cdc,
		txMemo:      txMemo,
		pubkey:      pubkey,
		isValidator: candiateInfo.IsVldt,
	}

	go m.monitorBlockHead()
	go m.monitorInitializeCandidate()
	go m.monitorDelegate()
	go m.monitorValidatorChange()
	go m.monitorIntendWithdraw()
	go m.monitorIntendSettle()
	go m.monitorWithdrawReward()
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
	initializeCandiateChan := make(chan *mainchain.GuardInitializeCandidate)
	sub, err := m.ethClient.Guard.WatchInitializeCandidate(nil, initializeCandiateChan, nil, nil)
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
	for {
		for page := 1; ; page++ {
			hasSeenEvent := false
			txs, err := authUtils.QueryTxsByEvents(m.transactor.CliCtx, []string{initiateWithdrawRewardEvent}, page, txsPageLimit)
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
					if event.Attributes[0].Value == validator.ActionInitiateWithdraw {
						m.handleInitiateWithdrawReward(event.Attributes[1].Value)
					}
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
