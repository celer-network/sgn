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

	go m.querySubscription()
	m.started = true
}

func (m *EthMonitor) querySubscription() {
	data, err := m.cdc.MarshalJSON(subscribe.NewQuerySubscrptionParams("0x674fa8ec8572f476f07b2bc7042e80a4f4d64107"))
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
