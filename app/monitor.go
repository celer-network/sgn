package app

import (
	"fmt"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
	"github.com/celer-network/sgn/x/subscribe/client/cli"
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
	sub, err := m.ethClient.Ledger.WatchIntendSettle(nil, intendSettleChan, [][32]byte{})
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
			m.handleIntendSettle(intendSettle)
		}
	}
}

func (m *EthMonitor) handleIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	// TODO: figure out query subscription by channel ID
	subscription, err := cli.QuerySubscrption(m.cdc, m.transactor.CliCtx, "subscribe", "1f7402f55e142820ea3812106d0657103fc1709e")

	tx, err := m.ethClient.Ledger.IntendSettle(m.ethClient.Auth, subscription.SignedSimplexStateBytes)
	if err != nil {
		fmt.Printf("tx err", err)
	}
	fmt.Printf("tx detail", tx)
}
