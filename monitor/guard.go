package monitor

import (
	"fmt"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/sync"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/protobuf/proto"
)

func (m *Monitor) processGuardQueue() {
	var keys, vals [][]byte
	m.dbLock.RLock()
	iterator, err := m.db.Iterator(GuardKeyPrefix, storetypes.PrefixEndBytes(GuardKeyPrefix))
	if err != nil {
		log.Errorln("Create db iterator err", err)
		return
	}
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, iterator.Key())
		vals = append(vals, iterator.Value())
	}
	iterator.Close()
	m.dbLock.RUnlock()

	for i, key := range keys {
		event := NewEventFromBytes(vals[i])
		if !event.Processing {
			log.Infoln("Process guard event", event.Name)
			var submitted bool
			switch e := event.ParseEvent(m.ethClient).(type) {
			case *mainchain.CelerLedgerIntendSettle:
				e.Raw = event.Log
				submitted, err = m.guardIntendSettle(e)
			case *mainchain.CelerLedgerIntendWithdraw:
				e.Raw = event.Log
				submitted, err = m.guardIntendWithdrawChannel(e)
			}
			if err != nil {
				log.Error(err)
				continue
			}
			if submitted {
				m.dbLock.Lock()
				v, err := m.db.Get(key)
				if err != nil {
					log.Errorln("db Get err:", err)
					m.dbLock.Unlock()
					continue
				}
				e := NewEventFromBytes(v)
				if !e.Processing {
					e.Processing = true
					err = m.db.Set(key, e.MustMarshal())
					if err != nil {
						log.Errorln("db Set err", err)
					}
				}
				m.dbLock.Unlock()
			}
		}
	}
}

func (m *Monitor) guardIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) (bool, error) {
	log.Infof("Guard IntendSettle %x, tx hash %x", intendSettle.ChannelId, intendSettle.Raw.TxHash)
	requests := m.getRequests(intendSettle.ChannelId)
	if len(requests) > 0 {
		return m.guardRequest(requests, intendSettle.Raw, IntendSettle)
	} else {
		err := m.dbDelete(GetGuardKey(intendSettle.Raw))
		return false, err
	}
}

func (m *Monitor) guardIntendWithdrawChannel(intendWithdrawChannel *mainchain.CelerLedgerIntendWithdraw) (bool, error) {
	log.Infof("Guard intendWithdrawChannel %x, tx hash %x", intendWithdrawChannel.ChannelId, intendWithdrawChannel.Raw.TxHash)
	requests := m.getRequests(intendWithdrawChannel.ChannelId)
	if len(requests) > 0 {
		return m.guardRequest(requests, intendWithdrawChannel.Raw, IntendWithdrawChannel)
	} else {
		err := m.dbDelete(GetGuardKey(intendWithdrawChannel.Raw))
		return false, err
	}
}

func (m *Monitor) guardRequest(requests []*guard.Request, rawLog ethtypes.Log, eventName EventName) (bool, error) {
	if len(requests) != 1 && len(requests) != 2 {
		return false, fmt.Errorf("invalid requests length")
	}

	isGuard := false
	for _, request := range requests {
		log.Infoln("guard request", request)
		if m.isRequestGuard(request, rawLog.BlockNumber) {
			isGuard = true
			break
		}
	}
	if !isGuard {
		log.Debug("not my turn to guard the requests")
		return false, nil
	}

	var stateArray chain.SignedSimplexStateArray
	for _, request := range requests {
		var signedSimplexState chain.SignedSimplexState
		err := proto.Unmarshal(request.SignedSimplexStateBytes, &signedSimplexState)
		if err != nil {
			log.Errorln("Unmarshal SignedSimplexState error:", err)
			continue
		}
		stateArray.SignedSimplexStates = append(stateArray.SignedSimplexStates, &signedSimplexState)
	}
	if len(stateArray.SignedSimplexStates) == 0 {
		return false, fmt.Errorf("invalid simplex states")
	}

	signedSimplexStateArrayBytes, err := proto.Marshal(&stateArray)
	if err != nil {
		log.Errorln("Marshal signedSimplexStateArrayBytes error:", err)
		return false, fmt.Errorf("marshal stateArray err %w", err)
	}

	var tx *ethtypes.Transaction
	switch eventName {
	case IntendWithdrawChannel:
		tx, err = m.ethClient.Transactor.Transact(
			m.guardTxHandler("SnapshotStates", requests, rawLog),
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return m.ethClient.Ledger.SnapshotStates(opts, signedSimplexStateArrayBytes)
			})
	case IntendSettle:
		tx, err = m.ethClient.Transactor.Transact(
			m.guardTxHandler("IntendSettle", requests, rawLog),
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return m.ethClient.Ledger.IntendSettle(opts, signedSimplexStateArrayBytes)
			})
	default:
		return false, fmt.Errorf("Invalid eventName %s", eventName)
	}
	if err != nil {
		return false, fmt.Errorf("tx err %w", err)
	} else {
		log.Infof("submitted guard tx %x", tx.Hash())
	}

	return true, nil
}

func (m *Monitor) guardTxHandler(
	description string, requests []*guard.Request, rawLog ethtypes.Log) *eth.TransactionStateHandler {
	return &eth.TransactionStateHandler{
		OnMined: func(receipt *ethtypes.Receipt) {
			if receipt.Status == ethtypes.ReceiptStatusSuccessful {
				log.Infof("%s transaction %x succeeded", description, receipt.TxHash)
				for _, request := range requests {
					guardProof := guard.NewGuardProof(
						mainchain.Bytes2Cid(request.ChannelId),
						mainchain.Hex2Addr(request.SimplexReceiver),
						receipt.TxHash,
						receipt.BlockNumber.Uint64(),
						m.ethClient.Address)
					syncData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(guardProof)
					msg := sync.NewMsgSubmitChange(sync.GuardProof, syncData, m.operator.Key.GetAddress())
					log.Infof("submit change tx: guard proof request %s", request)
					m.operator.AddTxMsg(msg)
				}
				err := m.dbDelete(GetGuardKey(rawLog))
				if err != nil {
					log.Errorln("db Delete err", err)
				}
			} else {
				log.Errorf("%s transaction %x failed", description, receipt.TxHash)
			}
		},
		OnError: func(tx *ethtypes.Transaction, err error) {
			log.Errorf("%s transaction %x err: %s", description, tx.Hash(), err)
		},
	}
}
