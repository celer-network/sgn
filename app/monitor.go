package app

import (
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EthMonitor struct {
	ethClient *mainchain.EthClient
	sgnApp    *sgnApp
	ctx       sdk.Context
	started   bool
}

func NewEthMonitor(ethClient *mainchain.EthClient, sgnApp *sgnApp) *EthMonitor {
	return &EthMonitor{
		ethClient: ethClient,
		sgnApp:    sgnApp,
	}
}

func (m *EthMonitor) Start(ctx sdk.Context) {
	if m.started {
		return
	}

	m.ctx = ctx
	m.started = true
}
