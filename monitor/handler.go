package monitor

import (
	log "github.com/celer-network/sgn/clog"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/validator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (m *EthMonitor) handleNewBlock(header *types.Header) {
	log.Infof("Push MsgSyncBlock with block number: %d to Transactor's msgQueue for broadcast.", header.Number)
	if !m.isPuller() {
		return
	}

	msg := global.NewMsgSyncBlock(header.Number.Uint64(), m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}

func (m *EthMonitor) handleInitializeCandidate(event *mainchain.GuardInitializeCandidate) {
	log.Infof("Push initializeCandidate event to pullerQueue: candidate %x, min self stake %s, sidechain addr %x",
		event.Candidate, event.MinSelfStake.String(), event.SidechainAddr)
	m.pullerQueue.PushBack(event)
}

func (m *EthMonitor) handleDelegate(delegate *mainchain.GuardDelegate) {
	log.Infoln("New delegate", delegate.Candidate)
	if m.isPullerOrOwner(delegate.Candidate.String()) {
		m.syncDelegator(delegate.Candidate, delegate.Delegator)
	}

	if m.isValidator {
		m.syncValidator(delegate.Candidate)
	} else {
		m.ethClaimValidator(delegate)
	}
}

func (m *EthMonitor) handleValidatorChange(validatorChange *mainchain.GuardValidatorChange) {
	log.Infoln("New validator change", validatorChange.EthAddr, validatorChange.ChangeType)
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
	log.Infoln("New intend withdraw", intendWithdraw.Candidate)

	if m.isPullerOrOwner(intendWithdraw.Candidate.String()) {
		m.syncValidator(intendWithdraw.Candidate)
	}
}

func (m *EthMonitor) handleIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	log.Infoln("New intend settle", intendSettle.ChannelId)
	request, err := m.getRequest(intendSettle.ChannelId[:])
	// TODO: a "not found" error should be regarded as a normal situation for not-guarded channels
	if err != nil {
		log.Errorln("Query request err", err)
		return
	}

	if intendSettle.SeqNums[request.PeerFromIndex].Uint64() >= request.SeqNum {
		log.Infoln("Ignore the intendSettle event with a larger seqNum")
		return
	}

	m.pusherQueue.PushBack(intendSettle)
}

func (m *EthMonitor) handleInitiateWithdrawReward(ethAddr string) {
	log.Infoln("New initiate withdraw", ethAddr)

	reward, err := validator.CLIQueryReward(m.transactor.CliCtx, validator.StoreKey, ethAddr)
	if err != nil {
		log.Errorln("Query reward err", err)
		return
	}

	sig, err := m.ethClient.SignMessage(reward.RewardProtoBytes)
	if err != nil {
		log.Errorln("SignMessage err", err)
		return
	}

	msg := validator.NewMsgSignReward(ethAddr, sig, m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}

func (m *EthMonitor) handlePenalty(nonce uint64) {
	log.Infoln("New Penalty", nonce)

	penalty, err := slash.CLIQueryPenalty(m.transactor.CliCtx, slash.StoreKey, nonce)
	if err != nil {
		log.Errorln("Query penalty err", err)
		return
	}

	sig, err := m.ethClient.SignMessage(penalty.PenaltyProtoBytes)
	if err != nil {
		log.Errorln("SignMessage err", err)
		return
	}

	msg := slash.NewMsgSignPenalty(nonce, sig, m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}

func (m *EthMonitor) ethClaimValidator(delegate *mainchain.GuardDelegate) {
	minStake, err := m.ethClient.Guard.GetMinStakingPool(&bind.CallOpts{})
	if err != nil {
		log.Errorln("GetMinStakingPool err", err)
		return
	}

	if delegate.StakingPool.Uint64() <= minStake.Uint64() {
		log.Error("Not enough stake to become validator")
		return
	}

	tx, err := m.ethClient.Guard.GuardTransactor.ClaimValidator(m.ethClient.Auth)
	if err != nil {
		log.Errorln("ClaimValidator tx err", err)
		return
	}
	log.Infoln("ClaimValidator tx detail", tx)
}

func (m *EthMonitor) claimValidator() {
	log.Info("ClaimValidator")
	transactors, err := transactor.ParseTransactorAddrs(m.transactors)
	if err != nil {
		log.Errorln("parse transactors err", err)
		return
	}

	msg := validator.NewMsgClaimValidator(m.ethClient.Address.String(), m.pubkey, transactors, m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)

}

func (m *EthMonitor) syncValidator(address ethcommon.Address) {
	log.Infoln("SyncValidator", address.String())
	msg := validator.NewMsgSyncValidator(address.String(), m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}

func (m *EthMonitor) syncDelegator(candidatorAddr, delegatorAddr ethcommon.Address) {
	log.Infoln("SyncDelegator", candidatorAddr.String(), delegatorAddr.String())
	msg := validator.NewMsgSyncDelegator(candidatorAddr.String(), delegatorAddr.String(), m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}
