package monitor

import (
	"log"

	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
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

func (m *EthMonitor) isPullerOrOwner(candidate string) bool {
	return m.isPuller() || candidate == m.ethClient.Address.String()
}

func (m *EthMonitor) isRequestHandler(request subscribe.Request, latestBlockNum uint64, eventBlockNumber uint64) bool {
	requestHanlders := request.RequestHandlers
	blockNumberDiff := latestBlockNum - eventBlockNumber
	handlerIndex := uint64(len(requestHanlders)+1) * blockNumberDiff / request.DisputeTimeout

	// All other validators need to guard
	if handlerIndex >= uint64(len(requestHanlders)) {
		return true
	}

	return requestHanlders[handlerIndex].Equals(m.transactor.Key.GetAddress())
}

func (m *EthMonitor) getRequest(channelId []byte) (subscribe.Request, error) {
	return subscribe.CLIQueryRequest(m.cdc, m.transactor.CliCtx, subscribe.StoreKey, channelId)
}

func (m *EthMonitor) getLatestBlock() (global.Block, error) {
	return global.CLIQueryLatestBlock(m.cdc, m.transactor.CliCtx, global.StoreKey)
}

// Get account info
func (m *EthMonitor) getAccount(addr sdk.AccAddress) (exported.Account, error) {
	accGetter := types.NewAccountRetriever(m.transactor.CliCtx)
	return accGetter.GetAccount(addr)
}
