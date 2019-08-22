package app

import (
	"fmt"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
	"github.com/celer-network/sgn/x/guardianmanager/client/cli"
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

	// TODO: make sure only submit intendSettle when sequence num is smaller
	// for {
	select {
	case err := <-sub.Err():
		fmt.Printf("WatchIntendSettle err", err)
	case intendSettle := <-intendSettleChan:
		m.handleIntendSettle(intendSettle)
	}
	// }
}

func (m *EthMonitor) handleIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	request, err := cli.QueryRequest(m.cdc, m.transactor.CliCtx, "guardianmanager", intendSettle.ChannelId[:])
	if err != nil {
		fmt.Printf("query request err", err)
		return
	}

	tx, err := m.ethClient.Ledger.IntendSettle(m.ethClient.Auth, request.SignedSimplexStateBytes)
	if err != nil {
		fmt.Printf("tx err", err)
		return
	}
	fmt.Printf("tx detail", tx)
}
