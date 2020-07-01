package monitor

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

		done, approve := m.verifyChange(change)
		if done {
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

// return (done, approve)
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
	case sync.InitGuardRequest:
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
	var subscription guard.Subscription
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
	var request guard.InitRequest
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &request)
	logmsg := fmt.Sprintf("verify change id %d, init request", change.ID)

	signedSimplexState, simplexChannel, err := common.UnmarshalSignedSimplexStateBytes(request.SignedSimplexStateBytes)
	if err != nil {
		log.Errorf("%s. unmarshal signedSimplexStateBytes err: %s", logmsg, err)
		return true, false
	}

	cid := mainchain.Bytes2Cid(simplexChannel.ChannelId)
	simplexReceiver, err := eth.RecoverSigner(request.SignedSimplexStateBytes, request.SimplexReceiverSig)
	if err != nil {
		log.Errorf("%s. recover signer err: %s", logmsg, err)
		return true, false
	}

	logmsg = logmsg + ". " + guard.PrintSimplexChannel(simplexChannel) + " " + fmt.Sprintf("to: %x", simplexReceiver)

	_, err = m.getRequest(simplexChannel.ChannelId, mainchain.Addr2Hex(simplexReceiver))
	if err == nil {
		log.Errorf("%s. request already initiated", logmsg)
		return true, false
	} else if !strings.Contains(err.Error(), common.ErrRecordNotFound.Error()) {
		log.Errorf("%s. getRequest err: %s", logmsg, err)
		return false, false
	}

	// verify signature
	simplexSender := mainchain.Bytes2Addr(simplexChannel.PeerFrom)
	err = guard.VerifySimplexStateSigs(signedSimplexState, simplexSender, simplexReceiver)
	if err != nil {
		log.Errorf("%s. verify sigs err: %s", logmsg, err)
		return true, false
	}

	// verify addr
	addrs, seqNums, err := m.ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{}, cid)
	if err != nil {
		log.Errorf("%s. GetStateSeqNumMap err: %s", logmsg, err)
		return false, false
	}
	seqIndex := 0
	var match bool
	if simplexSender == addrs[0] {
		match = (simplexReceiver == addrs[1])
	} else if simplexSender == addrs[1] {
		match = (simplexReceiver == addrs[0])
		seqIndex = 1
	}
	if !match {
		log.Errorf("%s. simplex addrs not match mainchain value %x %x", logmsg, addrs[0], addrs[1])
		return true, false
	}

	// verify seq
	if simplexChannel.SeqNum <= seqNums[seqIndex].Uint64() {
		log.Errorf("%s. SeqNum not larger than mainchain value %s", logmsg, seqNums[seqIndex])
		return true, false
	}

	// verify dispute timeout
	disputeTimeout, err := m.ethClient.Ledger.GetDisputeTimeout(&bind.CallOpts{}, cid)
	if err != nil {
		log.Errorf("%s. get dispute timeout err: %s", logmsg, err)
		return false, false
	}
	if disputeTimeout.Uint64() != request.DisputeTimeout {
		log.Errorf("%s. ispute timeout not match mainchain value %s", logmsg, disputeTimeout)
		return true, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifyTriggerGuard(change sync.Change) (bool, bool) {
	var trigger guard.GuardTrigger
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &trigger)
	logmsg := fmt.Sprintf("verify change id %d, trigger guard request: %s", change.ID, trigger)

	r, err := guard.CLIQueryRequest(m.operator.CliCtx, guard.RouterKey, trigger.ChannelId, trigger.SimplexReceiver)
	if err != nil {
		log.Errorf("%s. query request err: %s", logmsg, err)
		return false, false
	}

	if trigger.TriggerTxBlkNum == r.TriggerTxBlkNum {
		log.Errorf("%s. TriggerTxBlkNum not changed", logmsg)
		return true, false
	}

	triggerLog, err := guard.ValidateTriggerTx(
		m.ethClient, mainchain.Hex2Hash(trigger.TriggerTxHash), mainchain.Bytes2Cid(trigger.ChannelId))
	if err != nil {
		log.Errorf("%s. ValidateTriggerTx err: %s", logmsg, err)
		return false, false
	}

	if trigger.TriggerTxBlkNum != triggerLog.BlockNumber {
		log.Errorf("%s. TriggerTxBlkNum does not match mainchain value: %d", logmsg, triggerLog.BlockNumber)
		return false, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifyGuardProof(change sync.Change) (bool, bool) {
	var proof guard.GuardProof
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &proof)
	logmsg := fmt.Sprintf("verify change id %d, guard proof request: %s", change.ID, proof)

	r, err := guard.CLIQueryRequest(m.operator.CliCtx, guard.RouterKey, proof.ChannelId, proof.SimplexReceiver)
	if err != nil {
		log.Errorf("%s. query request err: %s", logmsg, err)
		return false, false
	}

	if r.TriggerTxHash == "" {
		log.Errorf("%s. Request Trigger event has not been submitted", logmsg)
		return true, false
	}

	guardLog, err := guard.ValidateGuardTx(
		m.ethClient, mainchain.Hex2Hash(proof.GuardTxHash), mainchain.Bytes2Cid(proof.ChannelId))
	if err != nil {
		log.Errorf("%s. ValidateGuardTx err: %s", logmsg, err)
		return false, false
	}

	if guardLog.BlockNumber <= r.TriggerTxBlkNum {
		log.Errorf("%s. Invalid block number for GuardTx at at %d", logmsg, guardLog.BlockNumber)
		return false, false
	}
	// TOTO: more check on request

	seqIndex := 0
	if bytes.Compare(mainchain.Hex2Addr(r.SimplexSender).Bytes(), mainchain.Hex2Addr(r.SimplexReceiver).Bytes()) > 0 {
		seqIndex = 1
	}
	err = guard.ValidateGuardProofSeqNum(guardLog.Data, uint8(seqIndex), r.SeqNum)
	if err != nil {
		log.Errorf("%s. ValidateGuardProofSeqNum err: %s", logmsg, err)
		return false, false
	}

	guardSender, err := mainchain.GetTxSender(m.ethClient.Client, proof.GuardTxHash)
	if err != nil {
		log.Errorf("%s. GetTxSender err: %s", logmsg, err)
		return false, false
	}

	if proof.GuardTxBlkNum != guardLog.BlockNumber {
		log.Errorf("%s. GuardTxBlkNum does not match mainchain value: %d", logmsg, guardLog.BlockNumber)
		return false, false
	}
	if proof.GuardSender != guardSender {
		log.Errorf("%s. GuardSender does not match mainchain value: %s", logmsg, guardSender)
		return false, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}
