package monitor

import (
	"log"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/validator"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (m *EthMonitor) handleNewBlock(header *types.Header) {
	log.Printf("New block", header.Number)
	puller := m.getPuller()
	if !puller.ValidatorAddr.Equals(m.transactor.Key.GetAddress()) {
		return
	}

	msg := global.NewMsgSyncBlock(header.Number.Uint64(), m.transactor.Key.GetAddress())
	_, err := m.transactor.BroadcastTx(msg)
	if err != nil {
		log.Printf("SyncBlock err", err)
		return
	}
}

func (m *EthMonitor) handleStake(stake *mainchain.GuardStake) {
	log.Printf("New stake", stake.NewStake)
	if m.isValidator {
		m.syncValidator(stake.Candidate)
	} else {
		// TODO: Check with mainchain to make sure that the candidate can become validator
		tx, err := m.ethClient.Guard.GuardTransactor.ClaimValidator(m.ethClient.Auth, m.transactor.Key.GetAddress().Bytes())
		if err != nil {
			log.Printf("ClaimValidator tx err", err)
			return
		}
		log.Printf("ClaimValidator tx detail", tx)
	}
}

func (m *EthMonitor) handleValidatorUpdate(vu *mainchain.GuardValidatorUpdate) {
	log.Printf("New validator update", vu.SidechainAddr)
	m.isValidator = vu.Added
	if m.isValidator {
		m.syncValidator(vu.EthAddr)
	}
}

func (m *EthMonitor) handleIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	log.Printf("New intend settle", intendSettle.ChannelId)
	request, err := m.getRequest(intendSettle.ChannelId[:])
	if err != nil {
		log.Printf("Query request err", err)
		return
	}

	if intendSettle.SeqNums[request.PeerFromIndex].Uint64() >= request.SeqNum {
		log.Printf("Ignore the intendSettle event due to larger seqNum")
		return
	}

	m.intendSettleQueue.PushBack(intendSettle)
}

func (m *EthMonitor) syncValidator(address ethcommon.Address) {
	msg := validator.NewMsgSyncValidator(address.String(), m.pubkey, m.transactor.Key.GetAddress())
	_, err := m.transactor.BroadcastTx(msg)
	if err != nil {
		log.Printf("SyncValidator err", err)
		return
	}
}
