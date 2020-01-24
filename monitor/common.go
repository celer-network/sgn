package monitor

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (m *EthMonitor) isPuller() bool {
	puller, err := validator.CLIQueryPuller(m.transactor.CliCtx, validator.StoreKey)
	if err != nil {
		log.Errorln("Get puller err", err)
		return false
	}

	return puller.ValidatorAddr.Equals(m.transactor.Key.GetAddress())
}

func (m *EthMonitor) isPusher() bool {
	pusher, err := validator.CLIQueryPusher(m.transactor.CliCtx, validator.StoreKey)
	if err != nil {
		log.Errorln("Get pusher err", err)
		return false
	}

	return pusher.ValidatorAddr.Equals(m.transactor.Key.GetAddress())
}

func (m *EthMonitor) isPullerOrOwner(candidate mainchain.Addr) bool {
	return m.isPuller() || candidate == m.ethClient.Address
}

// Is the current node the guard to submit state proof
func (m *EthMonitor) isRequestGuard(request subscribe.Request, latestBlockNum uint64, eventBlockNumber uint64) bool {
	requestGuards := request.RequestGuards
	blockNumberDiff := latestBlockNum - eventBlockNumber
	guardIndex := uint64(len(requestGuards)+1) * blockNumberDiff / request.DisputeTimeout

	log.Infoln("IsRequestGuard", latestBlockNum, eventBlockNumber, guardIndex, requestGuards)
	// All other validators need to guard
	if guardIndex >= uint64(len(requestGuards)) {
		return true
	}

	return requestGuards[guardIndex].Equals(m.transactor.Key.GetAddress())
}

func (m *EthMonitor) getRequest(channelId []byte, peerFrom string) (subscribe.Request, error) {
	return subscribe.CLIQueryRequest(m.transactor.CliCtx, subscribe.RouterKey, channelId, peerFrom)
}

func (m *EthMonitor) getLatestBlock() (global.Block, error) {
	return global.CLIQueryLatestBlock(m.transactor.CliCtx, global.RouterKey)
}

func (m *EthMonitor) getSecureBlockNum() (uint64, error) {
	return global.CLIQuerySecureBlockNum(m.transactor.CliCtx, global.RouterKey)
}

func (m *EthMonitor) getAccount(addr sdk.AccAddress) (exported.Account, error) {
	accGetter := types.NewAccountRetriever(m.transactor.CliCtx)
	return accGetter.GetAccount(addr)
}

func (m *EthMonitor) getGlobalParams() (global.Params, error) {
	return global.CLIQueryParams(m.transactor.CliCtx, global.RouterKey)
}
