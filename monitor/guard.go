package monitor

import (
	"encoding/json"
	"fmt"
	"strings"

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
	ChanState_Null           uint8 = 0
	ChanState_CatchWithdraw  uint8 = 1
	ChanState_SubmitWithdraw uint8 = 2
	ChanState_CatchSettle    uint8 = 3
	ChanState_SubmitSettle   uint8 = 4
	ChanState_Done           uint8 = 5
)

type ChanInfo struct {
	Cid             mainchain.CidType
	SimplexReceiver mainchain.Addr
	SeqNum          uint64
	State           uint8
	Guarded         bool
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
			var skip bool
			request, err := m.getGuardRequest(chanInfo.Cid.Bytes(), mainchain.Addr2Hex(chanInfo.SimplexReceiver))
			if err != nil {
				if !strings.Contains(err.Error(), common.ErrRecordNotFound.Error()) {
					log.Error(err)
					continue
				}
				log.Debugf("channel %x receiver %x not guarded by sgn", chanInfo.Cid, chanInfo.SimplexReceiver)
				skip = true
			} else if request.SeqNum <= chanInfo.SeqNum {
				log.Debugf("channel %x receiver %x does not have larger seqNum in sgn", chanInfo.Cid, chanInfo.SimplexReceiver)
				skip = true
			}
			if skip {
				err = m.dbDelete(GetGuardKey(chanInfo.Cid, chanInfo.SimplexReceiver))
				if err != nil {
					log.Errorln("db Delete err", err)
				}
				continue
			}

			if request.Status == common.GuardStatus_Withdraw || request.Status == common.GuardStatus_Settling {
				guarded, err := m.guardChannel(request, chanInfo.Cid)
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
}

func (m *Monitor) guardChannel(
	request *guard.Request, cid mainchain.CidType) (bool, error) {

	if request == nil {
		return false, fmt.Errorf("nil request")
	}
	if !m.isCurrentGuard(request, request.TriggerTxBlkNum) {
		log.Debugf("not my turn to guard request %s", request)
		return false, nil
	}

	log.Infof("Guard request %s", request)

	var stateArray chain.SignedSimplexStateArray
	var signedSimplexState chain.SignedSimplexState
	err := proto.Unmarshal(request.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		return false, fmt.Errorf("Unmarshal SignedSimplexState err: %w", err)
	}
	stateArray.SignedSimplexStates = append(stateArray.SignedSimplexStates, &signedSimplexState)

	signedSimplexStateArrayBytes, err := proto.Marshal(&stateArray)
	if err != nil {
		log.Errorln("Marshal signedSimplexStateArrayBytes error:", err)
		return false, fmt.Errorf("marshal stateArray err %w", err)
	}

	var tx *ethtypes.Transaction
	switch request.Status {
	case common.GuardStatus_Withdraw:
		tx, err = m.ethClient.Transactor.Transact(
			m.guardTxHandler("SnapshotStates", request, cid),
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return m.ethClient.Ledger.SnapshotStates(opts, signedSimplexStateArrayBytes)
			})
	case common.GuardStatus_Settling:
		tx, err = m.ethClient.Transactor.Transact(
			m.guardTxHandler("IntendSettle", request, cid),
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return m.ethClient.Ledger.IntendSettle(opts, signedSimplexStateArrayBytes)
			})
	default:
		return false, fmt.Errorf("Invalid guard state %d", request.Status)
	}
	if err != nil {
		return false, fmt.Errorf("tx err %w", err)
	} else {
		log.Infof("submitted guard tx %x", tx.Hash())
	}

	return true, nil
}

func (m *Monitor) guardTxHandler(
	description string, request *guard.Request, cid mainchain.CidType) *eth.TransactionStateHandler {
	guardState := common.GuardStatus_Idle
	if description == "IntendSettle" {
		guardState = common.GuardStatus_Settled
	}
	return &eth.TransactionStateHandler{
		OnMined: func(receipt *ethtypes.Receipt) {
			if receipt.Status == ethtypes.ReceiptStatusSuccessful {
				log.Infof("%s transaction %x succeeded", description, receipt.TxHash)
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
			} else {
				log.Errorf("%s transaction %x failed", description, receipt.TxHash)
			}
		},
		OnError: func(tx *ethtypes.Transaction, err error) {
			log.Errorf("%s transaction %x err: %s", description, tx.Hash(), err)
		},
	}
}
