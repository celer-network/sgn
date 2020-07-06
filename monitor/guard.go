package monitor

import (
	"encoding/json"
	"fmt"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/sync"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/protobuf/proto"
)

const (
	ChanState_Null              uint8 = 0
	ChanState_SettleWaiting     uint8 = 1
	ChanState_SettleSubmitted   uint8 = 2
	ChanState_WithdrawWaiting   uint8 = 3
	ChanState_WithdrawSubmitted uint8 = 4
	ChanState_Done              uint8 = 5
)

type ChanInfo struct {
	Cid        mainchain.CidType
	PeerStates map[mainchain.Addr]uint8
	Guarded    bool
}

func (ci *ChanInfo) marshal() []byte {
	val, err := json.Marshal(ci)
	if err != nil {
		log.Errorln("Marshal chanInfo err", err)
		return nil
	}
	return val
}

func unmarshalChanInfo(input []byte) *ChanInfo {
	var chanInfo ChanInfo
	err := json.Unmarshal(input, &chanInfo)
	if err != nil {
		log.Errorln("Unmarshal chanInfo err", err)
		return nil
	}
	return &chanInfo
}

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
		chanInfo := unmarshalChanInfo(vals[i])
		if !chanInfo.Guarded {
			requests, _ := m.getGuardRequests(chanInfo.Cid)
			if len(requests) == 0 {
				log.Infof("Ignore guard cid %x", chanInfo.Cid)
				err = m.dbDelete(GetGuardKey(chanInfo.Cid))
				if err != nil {
					log.Errorln("db Delete err", err)
				}
				continue
			}
			var triggered []*guard.Request
			for _, request := range requests {
				if request.GuardState == common.GuardState_Withdraw || request.GuardState == common.GuardState_Settling {
					triggered = append(triggered, request)
				}
			}
			if len(triggered) == 0 {
				log.Debugf("guard for channel %x not triggered yet", chanInfo.Cid)
				continue
			}
			log.Infof("Process guard cid %x", chanInfo.Cid)
			guarded, err := m.guardChannel(triggered, chanInfo.Cid)
			if err != nil {
				log.Error(err)
				continue
			}

			if guarded {
				m.dbLock.Lock()
				_, err := m.db.Get(key)
				if err != nil {
					log.Errorln("db Get err:", err)
					m.dbLock.Unlock()
					continue
				}
				chanInfo.Guarded = true
				err = m.db.Set(key, chanInfo.marshal())
				if err != nil {
					log.Errorln("db Set err", err)
				}
				m.dbLock.Unlock()
			}
		}
	}
}

func (m *Monitor) guardChannel(
	requests []*guard.Request, cid mainchain.CidType) (bool, error) {

	if len(requests) != 1 && len(requests) != 2 {
		return false, fmt.Errorf("invalid requests length")
	}

	isGuard := false
	for _, request := range requests {
		log.Infoln("guard request", request)
		if m.isCurrentGuard(request, request.TriggerTxBlkNum) {
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
	switch requests[0].GuardState {
	case common.GuardState_Withdraw:
		tx, err = m.ethClient.Transactor.Transact(
			m.guardTxHandler("SnapshotStates", requests, cid),
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return m.ethClient.Ledger.SnapshotStates(opts, signedSimplexStateArrayBytes)
			})
	case common.GuardState_Settling:
		tx, err = m.ethClient.Transactor.Transact(
			m.guardTxHandler("IntendSettle", requests, cid),
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return m.ethClient.Ledger.IntendSettle(opts, signedSimplexStateArrayBytes)
			})
	default:
		return false, fmt.Errorf("Invalid guard state %d", requests[0].GuardState)
	}
	if err != nil {
		return false, fmt.Errorf("tx err %w", err)
	} else {
		log.Infof("submitted guard tx %x", tx.Hash())
	}

	return true, nil
}

func (m *Monitor) guardTxHandler(
	description string, requests []*guard.Request, cid mainchain.CidType) *eth.TransactionStateHandler {
	guardState := common.GuardState_Idle
	if description == "IntendSettle" {
		guardState = common.GuardState_Settled
	}
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
						m.ethClient.Address,
						guardState)
					syncData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(guardProof)
					msg := sync.NewMsgSubmitChange(sync.GuardProof, syncData, m.operator.Key.GetAddress())
					log.Infof("submit change tx: guard proof request %s", request)
					m.operator.AddTxMsg(msg)
				}
				err := m.dbDelete(GetGuardKey(cid))
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
