package monitor

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/protobuf/proto"
)

func (m *Monitor) processGuardQueue() {
	iterator, err := m.db.Iterator(GuardKeyPrefix, storetypes.PrefixEndBytes(GuardKeyPrefix))
	if err != nil {
		log.Errorln("Create db iterator err", err)
		return
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewEventFromBytes(iterator.Value())
		log.Infoln("Process pusher event", event.Name)
		err = m.db.Delete(iterator.Key())
		if err != nil {
			log.Errorln("db Delete err", err)
			continue
		}

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.CelerLedgerIntendSettle:
			m.guardIntendSettle(e)
		case *mainchain.CelerLedgerIntendWithdraw:
			m.guardIntendWithdrawChannel(e)
		}
	}
}

func (m *Monitor) guardIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	log.Infof("Guard IntendSettle %x, tx hash %x", intendSettle.ChannelId, intendSettle.Raw.TxHash)
	requests := m.getRequests(intendSettle.ChannelId)
	for _, request := range requests {
		m.guardRequest(request, intendSettle.Raw, IntendSettle)
	}
}

func (m *Monitor) guardIntendWithdrawChannel(intendWithdrawChannel *mainchain.CelerLedgerIntendWithdraw) {
	log.Infof("Guard intendWithdrawChannel %x, tx hash %x", intendWithdrawChannel.ChannelId, intendWithdrawChannel.Raw.TxHash)
	requests := m.getRequests(intendWithdrawChannel.ChannelId)
	for _, request := range requests {
		m.guardRequest(request, intendWithdrawChannel.Raw, IntendWithdrawChannel)
	}
}

func (m *Monitor) guardRequest(request subscribe.Request, rawLog ethtypes.Log, eventName EventName) {
	log.Infoln("Guard request", request)
	if request.GuardTxHash != "" {
		log.Errorln("Request has been fulfilled")
		return
	}

	if !m.isRequestGuard(request, rawLog.BlockNumber) {
		log.Infof("Not valid guard at current mainchain block")
		event := NewEvent(IntendSettle, rawLog)
		err := m.db.Set(GetGuardKey(rawLog), event.MustMarshal())
		if err != nil {
			log.Errorln("db Set err", err)
		}
		return
	}

	var signedSimplexState chain.SignedSimplexState
	err := proto.Unmarshal(request.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		log.Errorln("Unmarshal SignedSimplexState error:", err)
		return
	}

	signedSimplexStateArrayBytes, err := proto.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{&signedSimplexState},
	})
	if err != nil {
		log.Errorln("Marshal signedSimplexStateArrayBytes error:", err)
		return
	}

	// TODO: use snapshotStates instead of intendSettle here? (need to update cChannel contract first)
	var receipt *ethtypes.Receipt
	switch eventName {
	case IntendWithdrawChannel:
		receipt, err = m.ethClient.Transactor.TransactWaitMined(
			"SnapshotStates",
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return m.ethClient.Ledger.SnapshotStates(opts, signedSimplexStateArrayBytes)
			})
	case IntendSettle:
		receipt, err = m.ethClient.Transactor.TransactWaitMined(
			"IntendSettle",
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return m.ethClient.Ledger.IntendSettle(opts, signedSimplexStateArrayBytes)
			})
	default:
		log.Errorln("Invalid eventName", eventName)
		return
	}

	if err != nil {
		log.Errorln("intendSettle/snapshotStates err", err)
		return
	}

	txHash := receipt.TxHash
	log.Infof("Add MsgGuardProof %x to transactor msgQueue", txHash)
	request.GuardTxHash = txHash.Hex()
	request.GuardTxBlkNum = receipt.BlockNumber.Uint64()
	request.GuardSender = mainchain.Addr2Hex(m.ethClient.Address)
	requestData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(request)
	msg := sync.NewMsgSubmitChange(sync.GuardProof, requestData, m.operator.Key.GetAddress())
	log.Infof("submit change tx: guard proof request %s", request)
	m.operator.AddTxMsg(msg)
}

func (m *Monitor) getRequests(cid [32]byte) (requests []subscribe.Request) {
	channelId := cid[:]
	addresses, seqNums, err := m.ethClient.Ledger.GetStateSeqNumMap(
		&bind.CallOpts{BlockNumber: sdk.NewIntFromUint64(m.secureBlkNum).BigInt()},
		cid)
	if err != nil {
		log.Errorln("Query StateSeqNumMap err", err)
		return
	}

	for _, addr := range addresses {
		peerFrom := mainchain.Addr2Hex(addr)
		request, err := m.getRequest(channelId, peerFrom)
		if err != nil {
			continue
		}

		if seqNums[request.PeerFromIndex].Uint64() >= request.SeqNum {
			log.Infoln("Ignore the intendSettle event with an equal or larger seqNum")
			continue
		}

		requests = append(requests, request)
	}

	return requests
}
