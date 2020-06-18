package monitor

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/celer-network/goutils/eth"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/protobuf/proto"
)

func (m *EthMonitor) verifyActiveChanges() {
	activeChanges, err := sync.CLIQueryActiveChanges(m.operator.CliCtx, sync.RouterKey)
	if err != nil {
		log.Errorln("Query active changes error:", err)
		return
	}

	for _, change := range activeChanges {
		_, err = m.verifiedChanges.Get(strconv.Itoa(int(change.ID)))
		if err == nil {
			continue
		}

		if m.verifyChange(change) {
			err = m.verifiedChanges.Set(strconv.Itoa(int(change.ID)), []byte{})
			if err != nil {
				log.Errorln("verifiedChanges Set err", err)
				continue
			}
			msg := sync.NewMsgApprove(change.ID, m.operator.Key.GetAddress())
			m.operator.AddTxMsg(msg)
		}
	}
}

func (m *EthMonitor) verifyChange(change sync.Change) bool {
	switch change.Type {
	case sync.ConfirmParamProposal:
		return m.verifyConfirmParamProposal(change)
	case sync.UpdateSidechainAddr:
		return m.verifyUpdateSidechainAddr(change)
	case sync.SyncDelegator:
		return m.verifySyncDelegator(change)
	case sync.SyncValidator:
		return m.verifySyncValidator(change)
	case sync.Subscribe:
		return m.verifySubscribe(change)
	case sync.Request:
		return m.verifyRequest(change)
	case sync.TriggerGuard:
		return m.verifyTriggerGuard(change)
	case sync.GuardProof:
		return m.verifyGuardProof(change)
	default:
		return false
	}
}

func (m *EthMonitor) verifyConfirmParamProposal(change sync.Change) bool {
	var paramChange common.ParamChange
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &paramChange)
	log.Infoln("Verify paramChange", paramChange)

	paramValue, err := m.ethClient.DPoS.GetUIntValue(&bind.CallOpts{}, paramChange.Record.BigInt())
	if err != nil {
		log.Errorln("Query param value error:", err)
		return false
	}

	if !paramChange.NewValue.Equal(sdk.NewIntFromBigInt(paramValue)) {
		log.Errorln("Param newValue does not match maichain value")
		return false
	}

	return true
}

func (m *EthMonitor) verifyUpdateSidechainAddr(change sync.Change) bool {
	var candidate validator.Candidate
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &candidate)
	log.Infoln("Verify candidate", candidate)

	c, err := validator.CLIQueryCandidate(m.operator.CliCtx, validator.RouterKey, candidate.EthAddress)
	if err == nil {
		if candidate.Operator.Equals(c.Operator) {
			log.Errorln("Invalid change for the same Operator value")
			return false
		}
	}

	sidechainAddr, err := m.ethClient.SGN.SidechainAddrMap(&bind.CallOpts{}, mainchain.Hex2Addr(candidate.EthAddress))
	if err != nil {
		log.Errorln("Query sidechain address error:", err)
		return false
	}

	if !candidate.Operator.Equals(sdk.AccAddress(sidechainAddr)) {
		log.Errorln("Operator does not match maichain value")
		return false
	}

	return true
}

func (m *EthMonitor) verifySyncDelegator(change sync.Change) bool {
	var delegator validator.Delegator
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &delegator)
	log.Infoln("Verify sync delegator", delegator)

	d, err := validator.CLIQueryDelegator(m.operator.CliCtx, validator.RouterKey, delegator.CandidateAddr, delegator.DelegatorAddr)
	if err == nil {
		if delegator.DelegatedStake.Equal(d.DelegatedStake) {
			log.Errorln("Invalid change for the same DelegatedStake value")
			return false
		}
	}

	di, err := m.ethClient.DPoS.GetDelegatorInfo(&bind.CallOpts{},
		mainchain.Hex2Addr(delegator.CandidateAddr), mainchain.Hex2Addr(delegator.DelegatorAddr))
	if err != nil {
		log.Errorln("Failed to query delegator info:", err)
		return false
	}

	if delegator.DelegatedStake.BigInt().Cmp(di.DelegatedStake) != 0 {
		log.Errorln("DelegatedStake does not match maichain value")
		return false
	}

	return true
}

func (m *EthMonitor) verifySyncValidator(change sync.Change) bool {
	var vt staking.Validator
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &vt)
	log.Infoln("Verify sync validator", vt)

	candidateEthAddr := vt.Description.Identity
	candidate, err := validator.CLIQueryCandidate(
		m.operator.CliCtx,
		validator.RouterKey,
		candidateEthAddr,
	)
	if err != nil {
		log.Errorln("Failed to query candidate:", err)
		return false
	}

	v, err := validator.CLIQueryValidator(
		m.operator.CliCtx,
		staking.RouterKey,
		candidate.Operator.String(),
	)
	if err == nil {
		if vt.Status.Equal(v.Status) && vt.Tokens.Equal(v.Tokens) && vt.Commission.Equal(v.Commission) {
			log.Errorln("Invalid change for the same Status/Tokens/Commission value")
			return false
		}
	}

	ci, err := m.ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(vt.Description.Identity))
	if err != nil {
		log.Errorln("Failed to query candidate info:", err)
		return false
	}

	if !vt.Status.Equal(mainchain.ParseStatus(ci)) {
		log.Errorln("Status does not match maichain value")
		return false
	}

	if !vt.Tokens.Equal(sdk.NewIntFromBigInt(ci.StakingPool).QuoRaw(common.TokenDec)) {
		log.Errorln("Tokens does not match maichain value")
		return false
	}

	commission, err := common.NewCommission(m.ethClient, ci.CommissionRate)
	if err != nil {
		log.Errorln("Failed to create new commission:", err)
		return false
	}

	if !vt.Commission.CommissionRates.Rate.Equal(commission.CommissionRates.Rate) {
		log.Errorln("Commission does not match maichain value")
		return false
	}

	return true
}

func (m *EthMonitor) verifySubscribe(change sync.Change) bool {
	var subscription subscribe.Subscription
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &subscription)
	log.Infoln("Verify subscription", subscription)

	deposit, err := m.ethClient.SGN.SubscriptionDeposits(
		&bind.CallOpts{},
		mainchain.Hex2Addr(subscription.EthAddress))
	if err != nil {
		log.Errorln("Failed to query subscription desposit:", err)
		return false
	}

	if subscription.Deposit.BigInt().Cmp(deposit) != 0 {
		log.Errorln("Deposit does not match maichain value")
		return false
	}

	return true
}

func (m *EthMonitor) verifyRequest(change sync.Change) bool {
	var request subscribe.Request
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &request)
	log.Infoln("Verify request", request)

	var signedSimplexState chain.SignedSimplexState
	err := proto.Unmarshal(request.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		log.Errorln("Failed to unmarshal signedSimplexStateBytes:", err)
		return false
	}

	r, err := subscribe.GetRequest(m.operator.CliCtx, m.ethClient.Ledger, &signedSimplexState)
	if err != nil {
		log.Errorln("Failed to get request through SignedSimplexStateBytes:", err)
		return false
	}

	err = subscribe.VerifySignedSimplexStateSigs(request, signedSimplexState)
	if err != nil {
		log.Errorln("Failed to verify sigs:", err)
		return false
	}

	ownerAddr, err := eth.RecoverSigner(request.SignedSimplexStateBytes, request.OwnerSig)
	if err != nil {
		log.Errorln("Failed to recover signer:", err)
		return false
	}

	if request.SeqNum <= r.SeqNum {
		log.Errorln("SeqNum is smaller than expected")
		return false
	}

	if request.PeerFromIndex != r.PeerFromIndex {
		log.Errorln("PeerFromIndex does not match maichain value")
		return false
	}

	if request.GetOwnerAddress() != mainchain.Addr2Hex(ownerAddr) {
		log.Errorln("Owner sig does not match maichain value")
		return false
	}

	if !bytes.Equal(request.ChannelId, r.ChannelId) {
		log.Errorln("ChannelId does not match maichain value")
		return false
	}

	if !reflect.DeepEqual(request.PeerAddresses, r.PeerAddresses) {
		log.Errorln("PeerAddresses does not match maichain value")
		return false
	}

	return true
}

func (m *EthMonitor) verifyTriggerGuard(change sync.Change) bool {
	var request subscribe.Request
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &request)
	log.Infoln("Verify TriggerGuard", request)

	r, err := subscribe.CLIQueryRequest(m.operator.CliCtx, subscribe.RouterKey, request.ChannelId, request.GetOwnerAddress())
	if err != nil {
		log.Errorln("Query request error:", err)
		return false
	}

	if request.TriggerTxBlkNum == r.TriggerTxBlkNum && request.DisputeTimeout == r.DisputeTimeout {
		log.Errorln("Invalid change for the same TriggerTxBlkNum/DisputeTimeout value")
		return false
	}

	triggerLog, err := subscribe.ValidateTriggerTx(m.ethClient, mainchain.Hex2Hash(request.TriggerTxHash), mainchain.Bytes2Cid(request.ChannelId))
	if err != nil {
		log.Errorln(err)
		return false
	}

	disputeTimeout, err := m.ethClient.Ledger.GetDisputeTimeout(&bind.CallOpts{}, mainchain.Bytes2Cid(request.ChannelId))
	if err != nil {
		log.Errorln("GetDisputeTimeout err:", err)
		return false
	}

	if request.TriggerTxBlkNum != triggerLog.BlockNumber {
		log.Errorln("TriggerTxBlkNum does not match maichain value")
		return false
	}
	if request.DisputeTimeout != disputeTimeout.Uint64() {
		log.Errorln("DisputeTimeout does not match maichain value")
		return false
	}

	return true
}

func (m *EthMonitor) verifyGuardProof(change sync.Change) bool {
	var request subscribe.Request
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &request)
	log.Infoln("Verify GuardProof", request)

	if request.TriggerTxHash == "" {
		log.Errorln("Request Trigger event has not been submitted")
		return false
	}

	guardLog, err := subscribe.ValidateGuardTx(m.ethClient, mainchain.Hex2Hash(request.GuardTxHash), mainchain.Bytes2Cid(request.ChannelId))
	if err != nil {
		log.Errorln(err)
		return false
	}

	if guardLog.BlockNumber <= request.TriggerTxBlkNum {
		log.Errorf("Invalid block number for GuardTx %d TriggerTx %d", guardLog.BlockNumber, request.TriggerTxBlkNum)
		return false
	}

	err = subscribe.ValidateSnapshotSeqNum(guardLog.Data, request.PeerFromIndex, request.SeqNum)
	if err != nil {
		log.Errorln(err)
		return false
	}

	guardSender, err := mainchain.GetTxSender(m.ethClient.Client, request.GuardTxHash)
	if err != nil {
		log.Errorln(err)
		return false
	}

	if request.GuardTxBlkNum != guardLog.BlockNumber {
		log.Errorln("GuardTxBlkNum does not match maichain value")
		return false
	}
	if request.GuardSender != guardSender {
		log.Errorln("GuardSender does not match maichain value")
		return false
	}

	return true
}
