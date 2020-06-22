package monitor

import (
	"bytes"
	"fmt"
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

func (m *Monitor) verifyActiveChanges() {
	v, _ := validator.CLIQueryValidator(m.operator.CliCtx, staking.RouterKey, m.operator.Key.GetAddress().String())
	if v.GetStatus() != sdk.Bonded {
		log.Traceln("skip verifying changes as I am not a bonded validator")
		return
	}
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

		verified, approve := m.verifyChange(change)
		if verified {
			err = m.verifiedChanges.Set(strconv.Itoa(int(change.ID)), []byte{})
			if err != nil {
				log.Errorln("verifiedChanges Set err", err)
				continue
			}
			if approve {
				msg := sync.NewMsgApprove(change.ID, m.operator.Key.GetAddress())
				m.operator.AddTxMsg(msg)
			}
		}
	}
}

// return (verified, approve)
func (m *Monitor) verifyChange(change sync.Change) (bool, bool) {
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
		return false, false
	}
}

func (m *Monitor) verifyConfirmParamProposal(change sync.Change) (bool, bool) {
	var paramChange common.ParamChange
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &paramChange)
	logmsg := fmt.Sprintf("verify change id %d, pramameter: %s", change.ID, paramChange)

	paramValue, err := m.ethClient.DPoS.GetUIntValue(&bind.CallOpts{}, paramChange.Record.BigInt())
	if err != nil {
		log.Errorf("%s. err: %s", logmsg, err)
		return false, false
	}

	if !paramChange.NewValue.Equal(sdk.NewIntFromBigInt(paramValue)) {
		log.Errorf("%s. new value does not match mainchain value: %s", logmsg, paramValue)
		return false, false
	}

	log.Infof("%s, success", logmsg)
	return true, true
}

func (m *Monitor) verifyUpdateSidechainAddr(change sync.Change) (bool, bool) {
	var candidate validator.Candidate
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &candidate)
	logmsg := fmt.Sprintf("verify change id %d, sidechain addr for candidate: %s", change.ID, candidate)

	c, err := validator.CLIQueryCandidate(m.operator.CliCtx, validator.RouterKey, candidate.EthAddress)
	if err == nil {
		if candidate.Operator.Equals(c.Operator) {
			log.Warnf("%s. sidechain addr not changed", logmsg)
			return true, false
		}
	}

	sidechainAddr, err := m.ethClient.SGN.SidechainAddrMap(&bind.CallOpts{}, mainchain.Hex2Addr(candidate.EthAddress))
	if err != nil {
		log.Errorf("%s. query sidechain address err: %s", logmsg, err)
		return false, false
	}

	if !candidate.Operator.Equals(sdk.AccAddress(sidechainAddr)) {
		log.Errorf("%s. operator does not match mainchain value: %s", logmsg, sdk.AccAddress(sidechainAddr))
		return false, false
	}

	log.Infof("%s, success", logmsg)
	return true, true
}

func (m *Monitor) verifySyncDelegator(change sync.Change) (bool, bool) {
	var delegator validator.Delegator
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &delegator)
	logmsg := fmt.Sprintf("verify change id %d, sync delegator: %s", change.ID, delegator)

	d, err := validator.CLIQueryDelegator(m.operator.CliCtx, validator.RouterKey, delegator.CandidateAddr, delegator.DelegatorAddr)
	if err == nil {
		if delegator.DelegatedStake.Equal(d.DelegatedStake) {
			log.Warnf("%s. delegator stake not changed", logmsg)
			return true, false
		}
	}

	di, err := m.ethClient.DPoS.GetDelegatorInfo(
		&bind.CallOpts{}, mainchain.Hex2Addr(delegator.CandidateAddr), mainchain.Hex2Addr(delegator.DelegatorAddr))
	if err != nil {
		log.Errorf("%s. query delegator info err: %s", logmsg, err)
		return false, false
	}

	if delegator.DelegatedStake.BigInt().Cmp(di.DelegatedStake) != 0 {
		log.Errorf("%s. stake does not match mainchain value: %s", logmsg, di.DelegatedStake)
		return false, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifySyncValidator(change sync.Change) (bool, bool) {
	var vt staking.Validator
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &vt)

	candidateEthAddr := vt.Description.Identity
	candidate, err := validator.CLIQueryCandidate(
		m.operator.CliCtx, validator.RouterKey, candidateEthAddr)
	if err != nil {
		log.Errorf("verify change id %d, sync validator err: %s", change.ID, err)
		return false, false
	}

	logmsg := fmt.Sprintf("verify change id %d, sync validator: Operator: %s, EthAddress %x, Status %s, Token %s, Commission %s",
		change.ID, candidate.Operator.String(), mainchain.Hex2Addr(candidateEthAddr), vt.Status, vt.Tokens, vt.Commission)

	v, err := validator.CLIQueryValidator(
		m.operator.CliCtx, staking.RouterKey, candidate.Operator.String())
	if err == nil {
		if vt.Status.Equal(v.Status) && vt.Tokens.Equal(v.Tokens) &&
			vt.Commission.CommissionRates.Rate.Equal(v.Commission.CommissionRates.Rate) {
			log.Warnf("%s. validator not changed", logmsg)
			return true, false
		}
	}

	ci, err := m.ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(vt.Description.Identity))
	if err != nil {
		log.Errorf("%s. query candidate info err: %s", logmsg, err)
		return false, false
	}

	if !vt.Status.Equal(mainchain.ParseStatus(ci)) {
		log.Errorf("%s. status does not match mainchain value: %s", logmsg, mainchain.ParseStatus(ci))
		return false, false
	}

	mtk := sdk.NewIntFromBigInt(ci.StakingPool).QuoRaw(common.TokenDec)
	if !vt.Tokens.Equal(mtk) {
		log.Errorf("%s. tokens does not match mainchain value: %s", logmsg, mtk)
		return false, false
	}

	commission, err := common.NewCommission(m.ethClient, ci.CommissionRate)
	if err != nil {
		log.Errorf("%s. create new commission err: %s", logmsg, err)
		return false, false
	}

	if !vt.Commission.CommissionRates.Rate.Equal(commission.CommissionRates.Rate) {
		log.Errorf("%s. commission does not match mainchain value: %s", logmsg, commission.CommissionRates.Rate)
		return false, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifySubscribe(change sync.Change) (bool, bool) {
	var subscription subscribe.Subscription
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &subscription)
	logmsg := fmt.Sprintf("verify change id %d, subscription: %s", change.ID, subscription)

	deposit, err := m.ethClient.SGN.SubscriptionDeposits(
		&bind.CallOpts{}, mainchain.Hex2Addr(subscription.EthAddress))
	if err != nil {
		log.Errorf("%s. query subscription desposit err: %s", logmsg, err)
		return false, false
	}

	if subscription.Deposit.BigInt().Cmp(deposit) != 0 {
		log.Errorf("%s. deposit does not match mainchain value: %s", logmsg, deposit)
		return false, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifyRequest(change sync.Change) (bool, bool) {
	var request subscribe.Request
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &request)
	logmsg := fmt.Sprintf("verify change id %d, request: %s", change.ID, request)

	var signedSimplexState chain.SignedSimplexState
	err := proto.Unmarshal(request.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		log.Errorf("%s. unmarshal signedSimplexStateBytes err: %s", logmsg, err)
		return false, false
	}

	r, err := subscribe.GetRequest(m.operator.CliCtx, m.ethClient.Ledger, &signedSimplexState)
	if err != nil {
		log.Errorf("%s. get request err: %s", logmsg, err)
		return false, false
	}

	err = subscribe.VerifySignedSimplexStateSigs(request, signedSimplexState)
	if err != nil {
		log.Errorf("%s. verify sigs err: %s", logmsg, err)
		return false, false
	}

	ownerAddr, err := eth.RecoverSigner(request.SignedSimplexStateBytes, request.OwnerSig)
	if err != nil {
		log.Errorf("%s. recover signer err: %s", logmsg, err)
		return false, false
	}

	if request.SeqNum <= r.SeqNum {
		log.Errorf("%s. SeqNum not larger than mainchain value %d", logmsg, r.SeqNum)
		return false, false
	}

	if request.PeerFromIndex != r.PeerFromIndex {
		log.Errorf("%s. PeerFromIndex does not match mainchain value: %d", logmsg, r.PeerFromIndex)
		return false, false
	}

	if request.GetOwnerAddress() != mainchain.Addr2Hex(ownerAddr) {
		log.Errorf("%s. Owner sig does not match mainchain value: %x", logmsg, ownerAddr)
		return false, false
	}

	if !bytes.Equal(request.ChannelId, r.ChannelId) {
		log.Errorf("%s. ChannelId does not match mainchain value: %x", logmsg, r.ChannelId)
		return false, false
	}

	if !reflect.DeepEqual(request.PeerAddresses, r.PeerAddresses) {
		log.Errorf("%s. PeerAddresses does not match mainchain value: %s", logmsg, r.PeerAddresses)
		return false, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifyTriggerGuard(change sync.Change) (bool, bool) {
	var request subscribe.Request
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &request)
	logmsg := fmt.Sprintf("verify change id %d, trigger guard request: %s", change.ID, request)

	r, err := subscribe.CLIQueryRequest(m.operator.CliCtx, subscribe.RouterKey, request.ChannelId, request.GetOwnerAddress())
	if err != nil {
		log.Errorf("%s. query request err: %s", logmsg, err)
		return false, false
	}

	if request.TriggerTxBlkNum == r.TriggerTxBlkNum && request.DisputeTimeout == r.DisputeTimeout {
		log.Errorf("%s. TriggerTxBlkNum and DisputeTimeout not changed", logmsg)
		return true, false
	}

	triggerLog, err := subscribe.ValidateTriggerTx(m.ethClient, mainchain.Hex2Hash(request.TriggerTxHash), mainchain.Bytes2Cid(request.ChannelId))
	if err != nil {
		log.Errorf("%s. ValidateTriggerTx err: %s", logmsg, err)
		return false, false
	}

	disputeTimeout, err := m.ethClient.Ledger.GetDisputeTimeout(&bind.CallOpts{}, mainchain.Bytes2Cid(request.ChannelId))
	if err != nil {
		log.Errorf("%s. GetDisputeTimeout err: %s", logmsg, err)
		return false, false
	}

	if request.TriggerTxBlkNum != triggerLog.BlockNumber {
		log.Errorf("%s. TriggerTxBlkNum does not match mainchain value: %d", logmsg, triggerLog.BlockNumber)
		return false, false
	}
	if request.DisputeTimeout != disputeTimeout.Uint64() {
		log.Errorf("%s. DisputeTimeout does not match mainchain value: %s", logmsg, disputeTimeout)
		return false, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifyGuardProof(change sync.Change) (bool, bool) {
	var request subscribe.Request
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &request)
	logmsg := fmt.Sprintf("verify change id %d, guard proof request: %s", change.ID, request)

	if request.TriggerTxHash == "" {
		log.Errorf("%s. Request Trigger event has not been submitted", logmsg)
		return false, false
	}

	guardLog, err := subscribe.ValidateGuardTx(m.ethClient, mainchain.Hex2Hash(request.GuardTxHash), mainchain.Bytes2Cid(request.ChannelId))
	if err != nil {
		log.Errorf("%s. ValidateGuardTx err: %s", logmsg, err)
		return false, false
	}

	if guardLog.BlockNumber <= request.TriggerTxBlkNum {
		log.Errorf("%s. Invalid block number for GuardTx at at %d", logmsg, guardLog.BlockNumber)
		return false, false
	}

	err = subscribe.ValidateSnapshotSeqNum(guardLog.Data, request.PeerFromIndex, request.SeqNum)
	if err != nil {
		log.Errorf("%s. ValidateSnapshotSeqNum err: %s", logmsg, err)
		return false, false
	}

	guardSender, err := mainchain.GetTxSender(m.ethClient.Client, request.GuardTxHash)
	if err != nil {
		log.Errorf("%s. GetTxSender err: %s", logmsg, err)
		return false, false
	}

	if request.GuardTxBlkNum != guardLog.BlockNumber {
		log.Errorf("%s. GuardTxBlkNum does not match mainchain value: %d", logmsg, guardLog.BlockNumber)
		return false, false
	}
	if request.GuardSender != guardSender {
		log.Errorf("%s. GuardSender does not match mainchain value: %s", logmsg, guardSender)
		return false, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}
