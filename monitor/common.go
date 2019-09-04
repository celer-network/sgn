package monitor

import (
	"log"

	"github.com/celer-network/sgn/x/validator"
	ethcommon "github.com/ethereum/go-ethereum/common"
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

func (m *EthMonitor) syncValidator(address ethcommon.Address) {
	msg := validator.NewMsgSyncValidator(address.String(), m.pubkey, m.transactor.Key.GetAddress())
	_, err := m.transactor.BroadcastTx(msg)
	if err != nil {
		log.Printf("SyncValidator err", err)
		return
	}
}
