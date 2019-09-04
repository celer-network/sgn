package monitor

import (
	"github.com/celer-network/sgn/mainchain"
)

func (m *EthMonitor) processQueue() {
	pusher := m.getPusher()
	if !pusher.ValidatorAddr.Equals(m.transactor.Key.GetAddress()) {
		return
	}

	for m.intendSettleQueue.Len() != 0 {
		m.handleIntendSettle(m.intendSettleQueue.PopFront().(*mainchain.CelerLedgerIntendSettle))
	}
}
