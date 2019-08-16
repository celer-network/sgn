package app

import (
	"fmt"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/cosmos/cosmos-sdk/codec"
)

type EthMonitor struct {
	ethClient  *mainchain.EthClient
	transactor *utils.Transactor
	cdc        *codec.Codec
	started    bool
}

func NewEthMonitor(ethClient *mainchain.EthClient, transactor *utils.Transactor, cdc *codec.Codec) *EthMonitor {
	return &EthMonitor{
		ethClient:  ethClient,
		transactor: transactor,
		cdc:        cdc,
	}
}

func (m *EthMonitor) Start() {
	if m.started {
		return
	}

	go m.monitorIntendSettle()
	m.started = true
}

func (m *EthMonitor) monitorIntendSettle() {
	intendSettleChan := make(chan *mainchain.CelerLedgerIntendSettle)
	sub, err := m.ethClient.Ledger.WatchIntendSettle
	if err != nil {
		fmt.Printf("WatchIntendSettle err", err)
		return
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			fmt.Printf("WatchIntendSettle err", err)
		case intendSettle := <-intendSettleChan:

		}
	}
}

func (m *EthMonitor) querySubscription(ethAddress string) {
	data, err := m.cdc.MarshalJSON(subscribe.NewQuerySubscrptionParams(ethAddress))
	if err != nil {
		return
	}
	route := fmt.Sprintf("custom/subscribe/%s", subscribe.QuerySubscrption)
	res, _, err := m.transactor.CliCtx.QueryWithData(route, data)
	if err != nil {
		fmt.Printf("query error", err)
	}

	fmt.Printf("query result", res)
}
