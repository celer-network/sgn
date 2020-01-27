package monitor

import (
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func (m *EthMonitor) handleNewBlock(header *types.Header) {
	log.Infoln("Catch new mainchain block", header.Number)
	if !m.isPuller() {
		return
	}

	params, err := m.getGlobalParams()
	if err != nil {
		log.Errorln("Query global params", err)
		return
	}

	time.Sleep(time.Duration(params.BlkTimeDiffLower+1) * time.Second)
	log.Infof("Add MsgSyncBlock %d to transactor msgQueue", header.Number)
	msg := global.NewMsgSyncBlock(header.Number.Uint64(), m.transactor.Key.GetAddress())
	m.transactor.AddTxMsg(msg)
}

func (m *EthMonitor) handleDelegate(delegate *mainchain.GuardDelegate) {
	log.Infof("Handle new delegate from delegator %x to candidate %x, stake %s pool %s",
		delegate.Delegator, delegate.Candidate, delegate.NewStake.String(), delegate.StakingPool.String())
	m.syncDelegator(delegate.Candidate, delegate.Delegator)

	if m.isValidator {
		m.syncValidator(delegate.Candidate)
	} else {
		m.claimValidatorOnMainchain(delegate)
	}
}

func (m *EthMonitor) handleValidatorChange(validatorChange *mainchain.GuardValidatorChange) {
	log.Infof("New validator change %x type %d", validatorChange.EthAddr, validatorChange.ChangeType)
	doSync := m.isPuller()

	if validatorChange.EthAddr == m.ethClient.Address {
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
	log.Infof("New intend withdraw %x", intendWithdraw.Candidate)

	if m.isPullerOrOwner(intendWithdraw.Candidate) {
		m.syncValidator(intendWithdraw.Candidate)
	}
}

func (m *EthMonitor) handleIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	channelId := intendSettle.ChannelId[:]
	log.Infof("New intend settle %x", channelId)
	doGuard := false
	addresses, seqNums, err := m.ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{}, intendSettle.ChannelId)
	if err != nil {
		log.Errorln("Query StateSeqNumMap err", err)
		return
	}

	for i := 0; i < 2; i++ {
		peerFrom := mainchain.Addr2Hex(addresses[i])
		request, err := m.getRequest(channelId, peerFrom)
		if err != nil {
			log.Errorln("Query request err", err)
			return
		}

		if seqNums[request.PeerFromIndex].Uint64() >= request.SeqNum {
			log.Infoln("Ignore the intendSettle event with a larger seqNum")
			return
		}

		doGuard = true
		msg := subscribe.NewMsgIntendSettle(channelId, peerFrom, intendSettle.Raw.TxHash.Hex(), m.transactor.Key.GetAddress())
		m.transactor.AddTxMsg(msg)
	}

	if doGuard {
		event := NewEvent(IntendSettle, intendSettle.Raw)
		m.db.Set(GetPusherKey(intendSettle.Raw), event.MustMarshal())
	}
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
	m.transactor.AddTxMsg(msg)
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
	m.transactor.AddTxMsg(msg)
}

func (m *EthMonitor) claimValidatorOnMainchain(delegate *mainchain.GuardDelegate) {
	minStake, err := m.ethClient.Guard.GetMinStakingPool(&bind.CallOpts{})
	if err != nil {
		log.Errorln("GetMinStakingPool err", err)
		return
	}

	if delegate.StakingPool.Uint64() <= minStake.Uint64() {
		log.Debug("Not enough stake to become validator")
		return
	}

	_, err = m.ethClient.Guard.ClaimValidator(m.ethClient.Auth)
	if err != nil {
		log.Errorln("ClaimValidator tx err", err)
		return
	}
	log.Infof("Claimed validator %x on mainchain", delegate.Candidate)
}

func (m *EthMonitor) claimValidator() {
	log.Infof("Claim self as a validator on sidechain, self address %x", m.ethClient.Address)
	transactors, err := transactor.ParseTransactorAddrs(m.transactors)
	if err != nil {
		log.Errorln("parse transactors err", err)
		return
	}

	msg := validator.NewMsgClaimValidator(
		mainchain.Addr2Hex(m.ethClient.Address), m.pubkey, transactors, m.transactor.Key.GetAddress())
	m.transactor.AddTxMsg(msg)

}

func (m *EthMonitor) syncValidator(address mainchain.Addr) {
	log.Infof("SyncValidator %x", address)
	msg := validator.NewMsgSyncValidator(mainchain.Addr2Hex(address), m.transactor.Key.GetAddress())
	m.transactor.AddTxMsg(msg)
}

func (m *EthMonitor) syncDelegator(candidatorAddr, delegatorAddr mainchain.Addr) {
	log.Infof("SyncDelegator candidate: %x, delegator: %x", candidatorAddr, delegatorAddr)

	msg := validator.NewMsgSyncDelegator(
		mainchain.Addr2Hex(candidatorAddr), mainchain.Addr2Hex(delegatorAddr), m.transactor.Key.GetAddress())
	m.transactor.AddTxMsg(msg)
}
