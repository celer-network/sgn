package monitor

import (
	"bytes"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type MonitorContractInfo struct {
	address common.Address
	abi     string
}

func (info *MonitorContractInfo) GetAddr() common.Address {
	return info.address
}

func (info *MonitorContractInfo) GetABI() string {
	return info.abi
}

func NewMonitorContractInfo(address common.Address, abi string) *MonitorContractInfo {
	return &MonitorContractInfo{
		address: address,
		abi:     abi,
	}
}

func (m *Monitor) isPuller() bool {
	puller, err := validator.CLIQueryPuller(m.operator.CliCtx, validator.StoreKey)
	if err != nil {
		log.Errorln("Get puller err", err)
		return false
	}

	return puller.ValidatorAddr.Equals(m.operator.Key.GetAddress())
}

func (m *Monitor) isPusher() bool {
	pusher, err := validator.CLIQueryPusher(m.operator.CliCtx, validator.StoreKey)
	if err != nil {
		log.Errorln("Get pusher err", err)
		return false
	}

	return pusher.ValidatorAddr.Equals(m.operator.Key.GetAddress())
}

func (m *Monitor) isPullerOrOwner(candidate mainchain.Addr) bool {
	return m.isPuller() || candidate == m.ethClient.Address
}

// Is the current node the guard to submit state proof
func (m *Monitor) isRequestGuard(request *guard.Request, eventBlockNumber uint64) bool {
	requestGuards := request.RequestGuards
	if len(requestGuards) == 0 {
		log.Debug("no request guards")
		return false
	}

	blkNum := m.ethMonitor.GetCurrentBlockNumber().Uint64()
	blockNumberDiff := blkNum - eventBlockNumber
	guardIndex := uint64(len(requestGuards)+1) * blockNumberDiff / request.DisputeTimeout

	log.Infoln("IsRequestGuard", blkNum, eventBlockNumber, guardIndex, requestGuards)
	// All other validators need to guard
	if guardIndex >= uint64(len(requestGuards)) {
		return true
	}

	return requestGuards[guardIndex].Equals(m.operator.Key.GetAddress())
}

func (m *Monitor) getRequest(channelId []byte, simplexReceiver string) (guard.Request, error) {
	return guard.CLIQueryRequest(m.operator.CliCtx, guard.RouterKey, channelId, simplexReceiver)
}

// get guard requests for the channel, return an array with at most two elements
func (m *Monitor) getRequests(cid mainchain.CidType) (requests []*guard.Request) {
	addresses, seqNums, err := m.ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{}, cid)
	if err != nil {
		log.Errorln("Query StateSeqNumMap err", err)
		return
	}

	for _, addr := range addresses {
		simplexReceiver := addr
		request, err := m.getRequest(cid.Bytes(), mainchain.Addr2Hex(simplexReceiver))
		if err != nil {
			continue
		}
		simplexSender := mainchain.Hex2Addr(request.SimplexSender)
		seqIndex := 0
		if bytes.Compare(simplexSender.Bytes(), simplexReceiver.Bytes()) > 0 {
			seqIndex = 1
		}
		if seqNums[seqIndex].Uint64() >= request.SeqNum {
			log.Infoln("Ignore the intendSettle event with an equal or larger seqNum")
			continue
		}

		requests = append(requests, &request)
	}

	return requests
}

func (m *Monitor) getAccount(addr sdk.AccAddress) (exported.Account, error) {
	accGetter := types.NewAccountRetriever(m.operator.CliCtx)
	return accGetter.GetAccount(addr)
}

func (m *Monitor) dbGet(key []byte) ([]byte, error) {
	m.dbLock.RLock()
	defer m.dbLock.RUnlock()
	return m.db.Get(key)
}

func (m *Monitor) dbSet(key, val []byte) error {
	m.dbLock.Lock()
	defer m.dbLock.Unlock()
	return m.db.Set(key, val)
}

func (m *Monitor) dbDelete(key []byte) error {
	m.dbLock.Lock()
	defer m.dbLock.Unlock()
	return m.db.Delete(key)
}
