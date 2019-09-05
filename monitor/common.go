package monitor

import (
	"log"

	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
)

func (m *EthMonitor) getPuller() validator.Puller {
	puller, err := validator.CLIQueryPuller(m.cdc, m.transactor.CliCtx, validator.StoreKey)
	if err != nil {
		log.Printf("Get puller err", err)
		return validator.Puller{}
	}

	return puller
}

func (m *EthMonitor) getPusher() validator.Pusher {
	pusher, err := validator.CLIQueryPusher(m.cdc, m.transactor.CliCtx, validator.StoreKey)
	if err != nil {
		log.Printf("Get pusher err", err)
		return validator.Pusher{}
	}

	return pusher
}

func (m *EthMonitor) getRequest(channelId []byte) (subscribe.Request, error) {
	return subscribe.CLIQueryRequest(m.cdc, m.transactor.CliCtx, subscribe.StoreKey, channelId)
}
