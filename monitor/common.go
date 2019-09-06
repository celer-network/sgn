package monitor

import (
	"log"

	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
)

func (m *EthMonitor) isPuller() bool {
	puller, err := validator.CLIQueryPuller(m.cdc, m.transactor.CliCtx, validator.StoreKey)
	if err != nil {
		log.Printf("Get puller err", err)
		return false
	}

	return puller.ValidatorAddr.Equals(m.transactor.Key.GetAddress())
}

func (m *EthMonitor) isPusher() bool {
	pusher, err := validator.CLIQueryPuller(m.cdc, m.transactor.CliCtx, validator.StoreKey)
	if err != nil {
		log.Printf("Get pusher err", err)
		return false
	}

	return pusher.ValidatorAddr.Equals(m.transactor.Key.GetAddress())
}

func (m *EthMonitor) getRequest(channelId []byte) (subscribe.Request, error) {
	return subscribe.CLIQueryRequest(m.cdc, m.transactor.CliCtx, subscribe.StoreKey, channelId)
}

func (m *EthMonitor) getLatestBlock() (global.Block, error) {
	return global.CLIQueryLatestBlock(m.cdc, m.transactor.CliCtx, global.StoreKey)
}
