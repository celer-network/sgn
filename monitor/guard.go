package monitor

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	esTxManager "github.com/celer-network/eth-services/txmanager"
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
	"github.com/google/uuid"
)

const (
	ChanInfoState_Null            uint8 = 0
	ChanInfoState_CaughtWithdraw  uint8 = 1
	ChanInfoState_GuardedWithdraw uint8 = 2
	ChanInfoState_CaughtSettle    uint8 = 3
	ChanInfoState_GuardedSettle   uint8 = 4
)

type ChanInfo struct {
	Cid             mainchain.CidType
	SimplexReceiver mainchain.Addr
	SeqNum          uint64
	State           uint8
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

	guardParams, err := guard.CLIQueryParams(m.Transactor.CliCtx, guard.RouterKey)
	if err != nil {
		log.Errorln("query guard params err", err)
		return
	}
	if mainchain.Hex2Addr(guardParams.LedgerAddress) != m.EthClient.GetLedger().Address {
		log.Infoln("update ledger address to", guardParams.LedgerAddress)
		m.EthClient.SetLedgerContract(guardParams.LedgerAddress)
		m.monitorCelerLedgerIntendSettle()
		m.monitorCelerLedgerIntendWithdraw()
	}

	var keys, vals [][]byte
	m.lock.RLock()
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
	m.lock.RUnlock()

	for i, key := range keys {
		chanInfo := unmarshalChanInfo(vals[i])
		if chanInfo == nil {
			err := m.dbDelete(key)
			if err != nil {
				log.Errorln("db Delete err", err)
			}
			continue
		}
		if chanInfo.State == ChanInfoState_GuardedWithdraw || chanInfo.State == ChanInfoState_GuardedSettle {
			continue
		}
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
			err = m.dbDelete(key)
			if err != nil {
				log.Errorln("db Delete err", err)
			}
			continue
		}

		if request.Status == guard.ChanStatus_Withdrawing || request.Status == guard.ChanStatus_Settling {
			shouldGuard, isAssigned := m.shouldGuardChannel(request, request.TriggerTxBlkNum)
			if !shouldGuard {
				continue
			}

			guarded, delete, err := m.guardChannel(request)
			if err != nil {
				log.Errorln("guardChannel err:", err)
				if !isAssigned {
					log.Infoln("non-assigned guard only tries once, should delete from guard queue")
					delete = true
				}
			}
			if delete {
				err = m.dbDelete(key)
				if err != nil {
					log.Errorln("db Delete err", err)
				}
			}
			if guarded {
				m.lock.Lock()
				exist, err := m.db.Has(key)
				if err != nil {
					log.Errorln("db Get err:", err)
					m.lock.Unlock()
					continue
				}
				if exist {
					val, err2 := m.db.Get(key)
					if err2 != nil {
						log.Errorln("db Get err", err2)
						m.lock.Unlock()
						continue
					}
					ci := unmarshalChanInfo(val)
					if ci != nil {
						chanInfo = ci
					}
				}
				if request.Status == guard.ChanStatus_Withdrawing {
					if chanInfo.State == ChanInfoState_CaughtWithdraw {
						chanInfo.State = ChanInfoState_GuardedWithdraw
					}
				} else {
					chanInfo.State = ChanInfoState_GuardedSettle
				}
				err = m.db.Set(key, chanInfo.marshal())
				if err != nil {
					log.Errorln("db Set err", err)
				}
				m.lock.Unlock()
			}
		}
	}
}

func (m *Monitor) guardChannel(request *guard.Request) (guarded, delete bool, err error) {
	log.Infof("start to guard %s", request)

	// tx pre-check: settle finalized time
	cid := mainchain.Bytes2Cid(request.ChannelId)
	settleFinalizedTime, err := m.EthClient.GetLedger().GetSettleFinalizedTime(&bind.CallOpts{}, cid)
	if err != nil {
		return false, false, fmt.Errorf("get settleFinalizedTime err: %w", err)
	}
	if settleFinalizedTime.Cmp(big.NewInt(0)) > 0 && m.getCurrentBlockNumber().Cmp(settleFinalizedTime) > 0 {
		log.Infof("channel %x passed settleFinalizedTime %s", cid, settleFinalizedTime)
		return false, true, nil
	}

	// tx pre-check: sequence number
	simplexSender := mainchain.Hex2Addr(request.SimplexSender)
	simplexReceiver := mainchain.Hex2Addr(request.SimplexReceiver)
	seqNum, err := mainchain.GetSimplexSeqNum(m.EthClient.GetLedger(), cid, simplexSender, simplexReceiver)
	if err != nil {
		return false, false, fmt.Errorf("GetSimplexSeqNum err: %w", err)
	}
	if request.SeqNum <= seqNum {
		log.Infof("channel %x stored seq %d no larger than mainchain value %d", cid, request.SeqNum, seqNum)
		return false, true, nil
	}

	// generate guard tx input
	var stateArray chain.SignedSimplexStateArray
	var signedSimplexState chain.SignedSimplexState
	err = proto.Unmarshal(request.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		return false, true, fmt.Errorf("Unmarshal SignedSimplexState err: %w", err)
	}
	stateArray.SignedSimplexStates = append(stateArray.SignedSimplexStates, &signedSimplexState)

	signedSimplexStateArrayBytes, err := proto.Marshal(&stateArray)
	if err != nil {
		return false, true, fmt.Errorf("marshal stateArray err %w", err)
	}

	// submit guard tx
	if m.Operator.TxSender == nil {
		var tx *ethtypes.Transaction
		switch request.Status {
		case guard.ChanStatus_Withdrawing:
			tx, err = m.EthClient.Transactor.Transact(
				m.guardTxHandler("SnapshotStates", request),
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return m.EthClient.GetLedger().SnapshotStates(opts, signedSimplexStateArrayBytes)
				})
		case guard.ChanStatus_Settling:
			tx, err = m.EthClient.Transactor.Transact(
				m.guardTxHandler("IntendSettle", request),
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return m.EthClient.GetLedger().IntendSettle(opts, signedSimplexStateArrayBytes)
				})
		default:
			return false, false, fmt.Errorf("Invalid guard state %d", request.Status)
		}
		if err != nil {
			return false, false, fmt.Errorf("tx err %w", err)
		}
		log.Infof("submitted guard tx %x", tx.Hash())
	} else {
		var guardErr error
		switch request.Status {
		case guard.ChanStatus_Withdrawing:
			guardErr = m.Operator.TxSender.SnapshotStates(signedSimplexStateArrayBytes, m.guardJobHandler("SnapshotStates", request))
		case guard.ChanStatus_Settling:
			guardErr = m.Operator.TxSender.IntendSettle(signedSimplexStateArrayBytes, m.guardJobHandler("IntendSettle", request))
		default:
			return false, false, fmt.Errorf("Invalid guard state %d", request.Status)
		}
		if guardErr != nil {
			return false, false, guardErr
		}
		log.Infof("submitted guard tx via TxSender")
	}

	return true, false, nil
}

func (m *Monitor) guardTxHandler(
	description string, request *guard.Request) *eth.TransactionStateHandler {
	guardState := guard.ChanStatus_Idle
	if description == "IntendSettle" {
		guardState = guard.ChanStatus_Settled
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
					m.EthClient.Address,
					guardState)
				syncData := m.Transactor.CliCtx.Codec.MustMarshalBinaryBare(guardProof)
				msg := sync.NewMsgSubmitChange(sync.GuardProof, syncData, m.EthClient.Client, m.Transactor.Key.GetAddress())
				log.Infof("submit change tx: guard proof request %s", request)
				m.Transactor.AddTxMsg(msg)
				err := m.dbDelete(GetGuardKey(mainchain.Bytes2Cid(request.ChannelId), mainchain.Hex2Addr(request.SimplexReceiver)))
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

func (m *Monitor) guardJobHandler(description string, request *guard.Request) esTxManager.JobHandler {
	guardState := guard.ChanStatus_Idle
	if description == "IntendSettle" {
		guardState = guard.ChanStatus_Settled
	}
	return func(jobID uuid.UUID, receipt *ethtypes.Receipt) error {
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			log.Infof("%s transaction %x succeeded", description, receipt.TxHash)
			guardProof := guard.NewGuardProof(
				mainchain.Bytes2Cid(request.ChannelId),
				mainchain.Hex2Addr(request.SimplexReceiver),
				receipt.TxHash,
				receipt.BlockNumber.Uint64(),
				m.EthClient.Address,
				guardState)
			syncData := m.Transactor.CliCtx.Codec.MustMarshalBinaryBare(guardProof)
			msg := sync.NewMsgSubmitChange(sync.GuardProof, syncData, m.EthClient.Client, m.Transactor.Key.GetAddress())
			log.Infof("submit change tx: guard proof request %s", request)
			m.Transactor.AddTxMsg(msg)
			err := m.dbDelete(GetGuardKey(mainchain.Bytes2Cid(request.ChannelId), mainchain.Hex2Addr(request.SimplexReceiver)))
			if err != nil {
				log.Errorln("db Delete err", err)
				return err
			}
		} else {
			txErr := fmt.Errorf("%s transaction %x failed", description, receipt.TxHash)
			log.Error(txErr)
			return txErr
		}
		return nil
	}
}

func (m *Monitor) setGuardEvent(eLog ethtypes.Log, state uint8) {
	var cid mainchain.CidType
	var withdrawReceiver mainchain.Addr // not used for settle event
	if state == ChanInfoState_CaughtSettle {
		e, err := m.EthClient.GetLedger().ParseIntendSettle(eLog)
		if err != nil {
			log.Errorln("ParseIntendSettle err", err)
			return
		}
		cid = e.ChannelId
	} else if state == ChanInfoState_CaughtWithdraw {
		e, err := m.EthClient.GetLedger().ParseIntendWithdraw(eLog)
		if err != nil {
			log.Errorln("ParseIntendWithdraw err", err)
			return
		}
		cid = e.ChannelId
		withdrawReceiver = e.Receiver
	} else {
		log.Errorln("invalid chanInfoState", state)
		return
	}
	addresses, seqNums, err := m.EthClient.GetLedger().GetStateSeqNumMap(&bind.CallOpts{}, cid)
	if err != nil {
		log.Errorf("Query StateSeqNumMap for cid %x err: %s", cid, err)
		return
	}
	for i, simplexReceiver := range addresses {
		if state == ChanInfoState_CaughtWithdraw && withdrawReceiver == simplexReceiver {
			// for intentWithdraw, only guard against the withdrawReceiver
			continue
		}
		m.setChanInfo(cid, simplexReceiver, state, seqNums[1-i].Uint64())
	}
}

func (m *Monitor) setChanInfo(cid mainchain.CidType, simplexReceiver mainchain.Addr, state uint8, seqNum uint64) {
	key := GetGuardKey(cid, simplexReceiver)
	m.lock.Lock()
	defer m.lock.Unlock()
	var chanInfo *ChanInfo
	exist, err := m.db.Has(key)
	if err != nil {
		log.Errorln("db Hash err", err)
		return
	}
	if exist {
		val, err2 := m.db.Get(key)
		if err2 != nil {
			log.Errorln("db Get err", err2)
			return
		}
		chanInfo = unmarshalChanInfo(val)
		log.Infof("ChanInfo for cid %x receiver %x already recorded", chanInfo.Cid, chanInfo.SimplexReceiver)
		chanInfo.SeqNum = seqNum
		if state == ChanInfoState_CaughtSettle {
			// IntendSettle has higher priority than IntendWithdraw
			if chanInfo.State == ChanInfoState_CaughtWithdraw || chanInfo.State == ChanInfoState_GuardedWithdraw {
				chanInfo.State = ChanInfoState_CaughtSettle
			}
		}
	} else {
		chanInfo = &ChanInfo{
			Cid:             cid,
			SimplexReceiver: simplexReceiver,
			SeqNum:          seqNum,
			State:           state,
		}
	}
	err = m.db.Set(key, chanInfo.marshal())
	if err != nil {
		log.Errorln("db Set err", err)
	}
}

// Is the current node the guard to submit state proof
func (m *Monitor) shouldGuardChannel(request *guard.Request, eventBlockNumber uint64) (shouldGuard, isAssigned bool) {
	reqlog := fmt.Sprintf("channel %x receiver %s", request.ChannelId, request.SimplexReceiver)
	assignedGuards := request.AssignedGuards
	if len(assignedGuards) == 0 {
		log.Debugf("no assigned guards for request %s", reqlog)
		return false, false
	}

	blkNum := m.getCurrentBlockNumber().Uint64()
	blockNumberDiff := blkNum - eventBlockNumber
	guardIndex := uint64(len(assignedGuards)+1) * blockNumberDiff / request.DisputeTimeout

	// All other validators need to guard
	if guardIndex >= uint64(len(assignedGuards)) {
		log.Debugf("should guard %s after assigned slots passed. current blk %d, event blk %d",
			reqlog, blkNum, eventBlockNumber)
		return true, false
	}
	actionlog := "not my turn to guard"
	shouldGuard = assignedGuards[guardIndex].Equals(m.Transactor.Key.GetAddress())
	if shouldGuard {
		isAssigned = true
		actionlog = "should guard"
	}
	log.Debugf("%s %s. index %d acct %s, current blk %d, event blk %d",
		actionlog, reqlog, guardIndex, assignedGuards[guardIndex], blkNum, eventBlockNumber)

	return shouldGuard, isAssigned
}
