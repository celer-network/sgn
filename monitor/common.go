package monitor

import (
	"math/big"
	"strings"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type MonitorContractInfo struct {
	address mainchain.Addr
	abi     string
}

func (info *MonitorContractInfo) GetAddr() mainchain.Addr {
	return info.address
}

func (info *MonitorContractInfo) GetABI() string {
	return info.abi
}

func NewMonitorContractInfo(address mainchain.Addr, abi string) *MonitorContractInfo {
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

// Is the current node the guard to submit state proof
func (m *Monitor) isCurrentGuard(request *guard.Request, eventBlockNumber uint64) bool {
	assignedGuards := request.AssignedGuards
	if len(assignedGuards) == 0 {
		log.Debug("no assigned guards")
		return false
	}

	blkNum := m.getCurrentBlockNumber().Uint64()
	blockNumberDiff := blkNum - eventBlockNumber
	guardIndex := uint64(len(assignedGuards)+1) * blockNumberDiff / request.DisputeTimeout

	// All other validators need to guard
	if guardIndex >= uint64(len(assignedGuards)) {
		log.Debugln("Assigned guard:", blkNum, eventBlockNumber, guardIndex)
		return true
	}
	log.Debugln("Assigned guard:", blkNum, eventBlockNumber, guardIndex, assignedGuards[guardIndex].String())

	return assignedGuards[guardIndex].Equals(m.operator.Key.GetAddress())
}

func (m *Monitor) getGuardRequest(channelId []byte, simplexReceiver string) (*guard.Request, error) {
	request, err := guard.CLIQueryRequest(m.operator.CliCtx, guard.RouterKey, channelId, simplexReceiver)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// get guard requests for the channel, return an array with at most two elements
// TODO: return proper err
func (m *Monitor) getGuardRequests(cid mainchain.CidType) (requests []*guard.Request, seqs []uint64) {
	addresses, seqNums, err := m.ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{}, cid)
	if err != nil {
		log.Errorln("Query StateSeqNumMap err", err)
		return
	}

	for i, simplexReceiver := range addresses {
		request, err := m.getGuardRequest(cid.Bytes(), mainchain.Addr2Hex(simplexReceiver))
		if err != nil {
			if !strings.Contains(err.Error(), common.ErrRecordNotFound.Error()) {
				log.Error(err)
			}
			continue
		}
		if seqNums[1-i].Uint64() >= request.SeqNum {
			log.Debugln("Ignore the intendSettle event with an equal or larger seqNum")
			continue
		}

		requests = append(requests, request)
		seqs = append(seqs, seqNums[1-i].Uint64())
	}

	return
}

func (m *Monitor) getCurrentBlockNumber() *big.Int {
	return m.ethMonitor.GetCurrentBlockNumber()
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

func chanInfoKey(cid mainchain.CidType, receiver mainchain.Addr) string {
	return mainchain.Cid2Hex(cid) + ":" + mainchain.Addr2Hex(receiver)
}
