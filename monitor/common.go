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

func (m *Monitor) isSyncer() bool {
	syncer, err := validator.CLIQuerySyncer(m.Transactor.CliCtx, validator.StoreKey)
	if err != nil {
		log.Errorln("Get syncer err", err)
		return false
	}

	return syncer.ValidatorAddr.Equals(m.Transactor.Key.GetAddress())
}

func (m *Monitor) getGuardRequest(channelId []byte, simplexReceiver string) (*guard.Request, error) {
	request, err := guard.CLIQueryRequest(m.Transactor.CliCtx, guard.RouterKey, channelId, simplexReceiver)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// get guard requests for the channel, return an array with at most two elements
// TODO: return proper err
func (m *Monitor) getGuardRequests(cid mainchain.CidType) (requests []*guard.Request, seqs []uint64) {
	addresses, seqNums, err := m.EthClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{}, cid)
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
			log.Debugln("Ignore request with an equal or larger mainchain seqnum")
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
	accGetter := types.NewAccountRetriever(m.Transactor.CliCtx)
	return accGetter.GetAccount(addr)
}

func (m *Monitor) dbGet(key []byte) ([]byte, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.db.Get(key)
}

func (m *Monitor) dbSet(key, val []byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.db.Set(key, val)
}

func (m *Monitor) dbDelete(key []byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.db.Delete(key)
}

func (m *Monitor) isBonded() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.bonded
}

func (m *Monitor) setBonded() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.bonded = true
}

func (m *Monitor) setUnbonded() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.bonded = false
}

func (m *Monitor) isBootstrapped() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.bootstrapped
}

func (m *Monitor) setBootstrapped() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.bootstrapped = true
}
