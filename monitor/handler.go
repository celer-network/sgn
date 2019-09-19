package monitor

import (
	"log"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/validator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (m *EthMonitor) handleNewBlock(header *types.Header) {
	log.Printf("New block", header.Number)
	if !m.isPuller() {
		return
	}

	msg := global.NewMsgSyncBlock(header.Number.Uint64(), m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}

func (m *EthMonitor) handleInitializeCandidate(initializeCandidate *mainchain.GuardInitializeCandidate) {
	log.Printf("New InitializeCandidate", initializeCandidate.Candidate)
	m.pullerQueue.PushBack(initializeCandidate)
}

func (m *EthMonitor) handleDelegate(delegate *mainchain.GuardDelegate) {
	log.Printf("New delegate", delegate.Candidate)
	if m.isValidator {
		m.syncValidator(delegate.Candidate)
	} else {
		m.ethClaimValidator(delegate)
	}

	if m.isPullerOrOwner(delegate.Candidate.String()) {
		m.syncDelegator(delegate.Candidate, delegate.Delegator)
	}
}

func (m *EthMonitor) handleValidatorChange(validatorChange *mainchain.GuardValidatorChange) {
	log.Printf("New validator change", validatorChange.EthAddr, validatorChange.ChangeType)
	doSync := m.isPuller()

	if validatorChange.EthAddr.String() == m.ethClient.Address.String() {
		m.isValidator = validatorChange.ChangeType == mainchain.AddValidator
		if m.isValidator {
			m.claimValidator()
			return
		}

		doSync = true
	}

	if doSync {
		m.syncValidator(validatorChange.EthAddr)
	}
}

func (m *EthMonitor) handleIntendWithdraw(intendWithdraw *mainchain.GuardIntendWithdraw) {
	log.Printf("New intend withdraw", intendWithdraw.Candidate)

	if m.isPullerOrOwner(intendWithdraw.Candidate.String()) {
		m.syncValidator(intendWithdraw.Candidate)
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

	m.pusherQueue.PushBack(intendSettle)
}

func (m *EthMonitor) ethClaimValidator(delegate *mainchain.GuardDelegate) {
	minStake, err := m.ethClient.Guard.GetMinStake(&bind.CallOpts{})
	if err != nil {
		log.Printf("GetMinStake err", err)
		return
	}

	if delegate.TotalStake.Uint64() <= minStake.Uint64() {
		log.Printf("Not enough stake to become validator")
		return
	}

	tx, err := m.ethClient.Guard.GuardTransactor.ClaimValidator(m.ethClient.Auth)
	if err != nil {
		log.Printf("ClaimValidator tx err", err)
		return
	}
	log.Printf("ClaimValidator tx detail", tx)
}

func (m *EthMonitor) claimValidator() {
	log.Printf("ClaimValidator")
	msg := validator.NewMsgClaimValidator(m.ethClient.Address.String(), m.pubkey, m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)

}

func (m *EthMonitor) syncValidator(address ethcommon.Address) {
	log.Printf("SyncValidator", address.String())
	msg := validator.NewMsgSyncValidator(address.String(), m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}

func (m *EthMonitor) syncDelegator(candidatorAddr, delegatorAddr ethcommon.Address) {
	log.Printf("SyncDelegator", candidatorAddr.String(), delegatorAddr.String())
	msg := validator.NewMsgSyncDelegator(candidatorAddr.String(), delegatorAddr.String(), m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}
