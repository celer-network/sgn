package monitor

import (
	"context"
	"log"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthMonitor struct {
	ethClient   *mainchain.EthClient
	transactor  *utils.Transactor
	cdc         *codec.Codec
	pubkey      string
	isValidator bool
}

func NewEthMonitor(ethClient *mainchain.EthClient, transactor *utils.Transactor, cdc *codec.Codec, pubkey string) {
	m := EthMonitor{
		ethClient:  ethClient,
		transactor: transactor,
		cdc:        cdc,
		pubkey:     pubkey,
	}

	// TODO: initiate isValidator value
	go m.monitorBlockHead()
	go m.monitorStake()
	go m.monitorValidatorUpdate()
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
		}
	}
}

func (m *EthMonitor) monitorStake() {
	stakeChan := make(chan *mainchain.GuardStake)
	sub, err := m.ethClient.Guard.WatchStake(nil, stakeChan, []ethcommon.Address{m.ethClient.Address})
	if err != nil {
		log.Printf("WatchStake err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("WatchStake err", err)
		case stake := <-stakeChan:
			m.handleStake(stake)
		}
	}
}

func (m *EthMonitor) monitorValidatorUpdate() {
	validatorUpdateChan := make(chan *mainchain.GuardValidatorUpdate)
	sub, err := m.ethClient.Guard.WatchValidatorUpdate(nil, validatorUpdateChan, []ethcommon.Address{m.ethClient.Address})
	if err != nil {
		log.Printf("WatchValidatorUpdate err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("WatchStake err", err)
		case validatorUpdate := <-validatorUpdateChan:
			m.handleValidatorUpdate(validatorUpdate)
		}
	}
}

func (m *EthMonitor) monitorIntendSettle() {
	intendSettleChan := make(chan *mainchain.CelerLedgerIntendSettle)
	sub, err := m.ethClient.Ledger.WatchIntendSettle(nil, intendSettleChan, [][32]byte{})
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
			m.handleIntendSettle(intendSettle)
		}
	}
}
