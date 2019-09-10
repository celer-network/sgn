package monitor

import (
	"context"
	"log"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gammazero/deque"
)

type EthMonitor struct {
	ethClient         *mainchain.EthClient
	transactor        *utils.Transactor
	cdc               *codec.Codec
	intendSettleQueue deque.Deque
	eventQueue        deque.Deque
	pubkey            string
	isValidator       bool
}

func NewEthMonitor(ethClient *mainchain.EthClient, transactor *utils.Transactor, cdc *codec.Codec, pubkey string) {
	candiateInfo, err := ethClient.Guard.GetCandidateInfo(&bind.CallOpts{}, ethClient.Address)
	if err != nil {
		log.Fatalf("GetCandidateInfo err", err)
		return
	}

	m := EthMonitor{
		ethClient:   ethClient,
		transactor:  transactor,
		cdc:         cdc,
		pubkey:      pubkey,
		isValidator: candiateInfo.IsVldt,
	}

	// TODO: initiate isValidator value
	go m.monitorBlockHead()
	go m.monitorDelegate()
	go m.monitorValidatorChange()
	go m.monitorIntendWithdraw()
	go m.monitorIntendSettle()
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
