package monitor

import (
	"bytes"
	"context"
	"errors"
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

var (
	intendSettleEventSig   = mainchain.GetEventSignature("IntendSettle(bytes32,uint256[2])")
	intendWithdrawEventSig = mainchain.GetEventSignature("IntendWithdraw(bytes32,address,uint256)")
	snapshotStatesEventSig = mainchain.GetEventSignature("SnapshotStates(bytes32,uint256[2])")
)

func (m *Monitor) verifyActiveChanges() {
	v, _ := validator.CLIQueryValidator(m.Transactor.CliCtx, staking.RouterKey, m.Transactor.Key.GetAddress().String())
	if v.GetStatus() != sdk.Bonded {
		log.Traceln("skip verifying changes as I am not a bonded validator")
		return
	}
	activeChanges, err := sync.CLIQueryActiveChanges(m.Transactor.CliCtx, sync.RouterKey)
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
				msg := sync.NewMsgApprove(change.ID, m.Transactor.Key.GetAddress())
				m.Transactor.AddTxMsg(msg)
			}
		}
	}
}

func (m *Monitor) verifyChange(change sync.Change) (done, approve bool) {
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
		return m.verifyInitGuardRequest(change)
	case sync.GuardTrigger:
		return m.verifyGuardTrigger(change)
	case sync.GuardProof:
		return m.verifyGuardProof(change)
	default:
		return false, false
	}
}

func (m *Monitor) verifyConfirmParamProposal(change sync.Change) (done, approve bool) {
	var paramChange common.ParamChange
	m.Transactor.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &paramChange)
	logmsg := fmt.Sprintf("verify change id %d, pramameter: %s", change.ID, paramChange)

	paramValue, err := m.EthClient.DPoS.GetUIntValue(&bind.CallOpts{}, paramChange.Record.BigInt())
	if err != nil {
		log.Errorf("%s. err: %s", logmsg, err)
		return false, false
	}

	if !paramChange.NewValue.Equal(sdk.NewIntFromBigInt(paramValue)) {
		if m.cmpBlkNum(change.BlockNum) == 1 {
			log.Infof("%s. new value not match mainchain value: %s", logmsg, paramValue)
			return true, false
		}
		log.Infof("%s. mainchain block not passed, value: %s", logmsg, paramValue)
		return false, false
	}

	log.Infof("%s, success", logmsg)
	return true, true
}

func (m *Monitor) verifyUpdateSidechainAddr(change sync.Change) (done, approve bool) {
	var candidate validator.Candidate
	m.Transactor.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &candidate)
	logmsg := fmt.Sprintf("verify change id %d, sidechain addr for candidate: %s", change.ID, candidate)

	c, err := validator.CLIQueryCandidate(m.Transactor.CliCtx, validator.RouterKey, candidate.EthAddress)
	if err == nil {
		if candidate.ValAccount.Equals(c.ValAccount) {
			log.Infof("%s. sidechain addr already updated", logmsg)
			return true, false
		}
	}

	sidechainAddr, err := m.EthClient.SGN.SidechainAddrMap(&bind.CallOpts{}, mainchain.Hex2Addr(candidate.EthAddress))
	if err != nil {
		log.Errorf("%s. query sidechain address err: %s", logmsg, err)
		return false, false
	}

	acctAddr := sdk.AccAddress(sidechainAddr)
	err = sdk.VerifyAddressFormat(acctAddr)
	if err != nil {
		log.Errorf("%s. verify address format err: %s", logmsg, err)
		return true, false
	}

	if !candidate.ValAccount.Equals(acctAddr) {
		if m.cmpBlkNum(change.BlockNum) == 1 {
			log.Infof("%s. validator account not match mainchain value: %s", logmsg, acctAddr)
			return true, false
		}
		log.Infof("%s. mainchain block not passed, validator account: %s", logmsg, acctAddr)
		return false, false
	}

	log.Infof("%s, success", logmsg)
	return true, true
}

func (m *Monitor) verifySyncDelegator(change sync.Change) (done, approve bool) {
	var delegator validator.Delegator
	m.Transactor.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &delegator)
	logmsg := fmt.Sprintf("verify change id %d, sync delegator: %s", change.ID, delegator)

	d, err := validator.CLIQueryDelegator(m.Transactor.CliCtx, validator.RouterKey, delegator.CandidateAddr, delegator.DelegatorAddr)
	if err == nil {
		if delegator.DelegatedStake.Equal(d.DelegatedStake) {
			log.Infof("%s. delegator stake already updated", logmsg)
			return true, false
		}
	} else if strings.Contains(err.Error(), common.ErrRecordNotFound.Error()) {
		if delegator.DelegatedStake.IsZero() {
			log.Debugf("%s. delegator with zero stake already removed", logmsg)
			return true, false
		}
	}

	di, err := m.EthClient.DPoS.GetDelegatorInfo(
		&bind.CallOpts{}, mainchain.Hex2Addr(delegator.CandidateAddr), mainchain.Hex2Addr(delegator.DelegatorAddr))
	if err != nil {
		log.Errorf("%s. query delegator info err: %s", logmsg, err)
		return false, false
	}

	if delegator.DelegatedStake.BigInt().Cmp(di.DelegatedStake) != 0 {
		if m.cmpBlkNum(change.BlockNum) == 1 {
			log.Infof("%s. stake not match mainchain value: %s", logmsg, di.DelegatedStake)
			return true, false
		}
		log.Infof("%s. mainchain block not passed, stake: %s", logmsg, di.DelegatedStake)
		return false, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifySyncValidator(change sync.Change) (done, approve bool) {
	var newVal staking.Validator
	m.Transactor.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &newVal)

	candidateEthAddr := newVal.Description.Identity
	candidate, err := validator.CLIQueryCandidate(
		m.Transactor.CliCtx, validator.RouterKey, candidateEthAddr)
	if err != nil {
		log.Errorf("verify change id %d, sync validator err: %s", change.ID, err)
		return false, false
	}

	logmsg := fmt.Sprintf("verify change id %d, sync validator: Account %s, EthAddress %x, Status %s, Token %s, Commission %s",
		change.ID, candidate.ValAccount, mainchain.Hex2Addr(candidateEthAddr), newVal.Status, newVal.Tokens, newVal.Commission.Rate)

	storedVal, err := validator.CLIQueryValidator(
		m.Transactor.CliCtx, staking.RouterKey, candidate.ValAccount.String())
	if err == nil {
		if newVal.Status.Equal(storedVal.Status) && newVal.Tokens.Equal(candidate.StakingPool) &&
			newVal.Commission.Rate.Equal(storedVal.Commission.Rate) {
			log.Infof("%s. validator already updated", logmsg)
			return true, false
		}
	}

	ci, err := m.EthClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(newVal.Description.Identity))
	if err != nil {
		log.Errorf("%s. query candidate info err: %s", logmsg, err)
		return false, false
	}

	if !newVal.Status.Equal(mainchain.ParseStatus(ci)) {
		if m.cmpBlkNum(change.BlockNum) == 1 {
			log.Infof("%s. status not match mainchain value: %s", logmsg, mainchain.ParseStatus(ci))
			return true, false
		}
		log.Infof("%s. mainchain block not passed, status: %s", logmsg, mainchain.ParseStatus(ci))
		return false, false
	}

	if !newVal.Tokens.Equal(sdk.NewIntFromBigInt(ci.StakingPool)) {
		if m.cmpBlkNum(change.BlockNum) == 1 {
			log.Infof("%s. staking pool not match mainchain value: %s", logmsg, ci.StakingPool)
			return true, false
		}
		log.Infof("%s. mainchain block not passed, token: %s", logmsg, ci.StakingPool)
		return false, false
	}

	commission, err := common.NewCommission(m.EthClient, ci.CommissionRate)
	if err != nil {
		log.Errorf("%s. create new commission err: %s", logmsg, err)
		return false, false
	}

	if !newVal.Commission.Rate.Equal(commission.Rate) {
		if m.cmpBlkNum(change.BlockNum) == 1 {
			log.Infof("%s. commission not match mainchain value: %s", logmsg, commission.Rate)
			return true, false
		}
		log.Infof("%s. mainchain block not passed, commission: %s", logmsg, commission.Rate)
		return false, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifySubscribe(change sync.Change) (done, approve bool) {
	var subscription guard.Subscription
	m.Transactor.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &subscription)
	logmsg := fmt.Sprintf("verify change id %d, subscription: %s", change.ID, subscription)

	deposit, err := m.EthClient.SGN.SubscriptionDeposits(
		&bind.CallOpts{}, mainchain.Hex2Addr(subscription.EthAddress))
	if err != nil {
		log.Errorf("%s. query subscription desposit err: %s", logmsg, err)
		return false, false
	}

	if subscription.Deposit.BigInt().Cmp(deposit) != 0 {
		if m.cmpBlkNum(change.BlockNum) == 1 {
			log.Infof("%s. deposit not match mainchain value: %s", logmsg, deposit)
			return true, false
		}
		log.Infof("%s. mainchain block not passed, deposit: %s", logmsg, deposit)
		return false, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifyInitGuardRequest(change sync.Change) (done, approve bool) {
	var request guard.InitRequest
	m.Transactor.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &request)
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

	logmsg += fmt.Sprintf(". %s, to %x, disputeTimeout: %d",
		guard.PrintSimplexChannel(simplexChannel), simplexReceiver, request.DisputeTimeout)

	_, err = m.getGuardRequest(simplexChannel.ChannelId, mainchain.Addr2Hex(simplexReceiver))
	if err == nil {
		log.Infof("%s. request already initiated", logmsg)
		return true, false
	} else if !strings.Contains(err.Error(), common.ErrRecordNotFound.Error()) {
		log.Errorf("%s. getGuardRequest err: %s", logmsg, err)
		return false, false
	}

	// verify signature
	simplexSender := mainchain.Bytes2Addr(simplexChannel.PeerFrom)
	err = guard.VerifySimplexStateSigs(signedSimplexState, simplexSender, simplexReceiver)
	if err != nil {
		log.Errorf("%s. verify sigs err: %s", logmsg, err)
		return true, false
	}

	// verify channel state
	chanState, err := m.EthClient.GetLedger().GetChannelStatus(&bind.CallOpts{}, cid)
	if err != nil {
		log.Errorf("%s. GetChannelStatus err: %s", logmsg, err)
		return false, false
	}
	if chanState != mainchain.ChannelStatus_OPERABLE {
		log.Errorf("%s. channel not in operable state: %d", logmsg, chanState)
		return true, false
	}

	// verify addr and seq
	seqNum, err := mainchain.GetSimplexSeqNum(m.EthClient.GetLedger(), cid, simplexSender, simplexReceiver)
	if err != nil {
		log.Errorf("%s. GetSimplexSeqNum err: %s", logmsg, err)
		return errors.Is(err, mainchain.ErrPeersNotMatch), false
	}
	if simplexChannel.SeqNum <= seqNum {
		log.Errorf("%s. SeqNum not larger than mainchain value %d", logmsg, seqNum)
		return true, false
	}

	// verify dispute timeout
	disputeTimeout, err := m.EthClient.GetLedger().GetDisputeTimeout(&bind.CallOpts{}, cid)
	if err != nil {
		log.Errorf("%s. get dispute timeout err: %s", logmsg, err)
		return false, false
	}
	if disputeTimeout.Uint64() != request.DisputeTimeout {
		log.Errorf("%s. dispute timeout not match mainchain value %s", logmsg, disputeTimeout)
		return true, false
	}
	guardParams, err := guard.CLIQueryParams(m.Transactor.CliCtx, guard.RouterKey)
	if err != nil {
		log.Errorf("%s. query guard params err: %s", logmsg, err)
		return false, false
	}
	if request.DisputeTimeout < guardParams.MinDisputeTimeout {
		log.Errorf("%s. dispute timeout smaller than min value %d", logmsg, guardParams.MinDisputeTimeout)
		return true, false
	}

	// verify not in active unilateral withdraw state
	wrecv, _, wblk, _, err := m.EthClient.GetLedger().GetWithdrawIntent(&bind.CallOpts{}, cid)
	if wrecv != mainchain.ZeroAddr {
		if m.getCurrentBlockNumber().Uint64() <= wblk.Uint64()+request.DisputeTimeout {
			log.Errorf("%s. channel has pending unilateral withdrawal request", logmsg)
			return true, false
		}
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifyGuardTrigger(change sync.Change) (done, approve bool) {
	var trigger guard.GuardTrigger
	m.Transactor.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &trigger)
	logmsg := fmt.Sprintf("verify change id %d, trigger guard request: %s", change.ID, trigger)

	// verify request exists
	r, err := guard.CLIQueryRequest(m.Transactor.CliCtx, guard.RouterKey, trigger.ChannelId, trigger.SimplexReceiver)
	if err != nil {
		log.Errorf("%s. query request err: %s", logmsg, err)
		return false, false
	}

	// verify trigger block number
	if trigger.TriggerTxBlkNum <= r.TriggerTxBlkNum {
		log.Infof("%s. TriggerTxBlkNum not greater than stored value %d", logmsg, r.TriggerTxBlkNum)
		return true, false
	}

	// verify seqNum
	cid := mainchain.Bytes2Cid(trigger.ChannelId)
	simplexSender := mainchain.Hex2Addr(r.SimplexSender)
	simplexReceiver := mainchain.Hex2Addr(r.SimplexReceiver)
	seqNum, err := mainchain.GetSimplexSeqNum(m.EthClient.GetLedger(), cid, simplexSender, simplexReceiver)
	if err != nil {
		log.Errorf("%s. GetSimplexSeqNum err: %s", logmsg, err)
		return errors.Is(err, mainchain.ErrPeersNotMatch), false
	}
	if r.SeqNum <= seqNum {
		if m.cmpBlkNum(change.BlockNum) == 1 {
			log.Errorf("%s. Stored SeqNum %d not larger than mainchain value %d", logmsg, r.SeqNum, seqNum)
			return true, false
		}
		log.Infof("%s. mainchain block not passed, stored SeqNum %d, mainchain value %d", logmsg, r.SeqNum, seqNum)
		return false, false
	}

	// verify triggerSeqNum
	if trigger.TriggerSeqNum > seqNum {
		log.Errorf("%s. TriggerSeqNum greater than mainchain value %d", logmsg, seqNum)
		return true, false
	}

	// verify onchain trasaction receipt and status
	receipt, err := m.EthClient.Client.TransactionReceipt(context.Background(), mainchain.Hex2Hash(trigger.TriggerTxHash))
	if err != nil {
		log.Errorf("%s. Get trigger tx receipt err: %s", logmsg, err)
		return false, false
	}
	if receipt.Status != mainchain.TxSuccess {
		log.Errorf("%s. Trigger tx failed: %d", logmsg, receipt.Status)
		return true, false
	}
	triggerLog := receipt.Logs[len(receipt.Logs)-1] // IntendSettle/IntendWithdraw event is the last one

	// verify transaction contract address
	if triggerLog.Address != m.EthClient.GetLedger().Address {
		log.Errorf("%s. Trigger tx contract address not match: %x", logmsg, triggerLog.Address)
		return true, false
	}

	// verify transaction event type and current request status
	if triggerLog.Topics[0] == intendSettleEventSig {
		if trigger.Status != guard.ChanStatus_Settling {
			log.Errorf("%s. Trigger guard state should be settling", logmsg)
			return true, false
		}
		if r.Status == guard.ChanStatus_Settling || r.Status == guard.ChanStatus_Settled {
			log.Errorf("%s. Invalid ChanStatus current state: %d", logmsg, r.Status)
			return true, false
		}
	} else if triggerLog.Topics[0] == intendWithdrawEventSig {
		if trigger.Status != guard.ChanStatus_Withdrawing {
			log.Errorf("%s. Trigger guard state should be withdraw", logmsg)
			return true, false
		}
		if r.Status != guard.ChanStatus_Idle {
			log.Errorf("%s. Invalid ChanStatus current state: %d", logmsg, r.Status)
			return true, false
		}
		event, err2 := m.EthClient.GetLedger().ParseIntendWithdraw(*triggerLog)
		if err2 != nil {
			log.Errorf("%s. ParseIntendWithdraw event err %s", logmsg, err2)
			return true, false
		}
		if event.Receiver == mainchain.Hex2Addr(r.SimplexReceiver) {
			log.Errorf("%s intendWithdraw receiver is the simplex receiver", logmsg)
			return true, false
		}
	} else {
		log.Errorf("%s. Trigger Tx is not for IntendSettle/IntendWithdraw event", logmsg)
		return true, false
	}

	// verify transaction channel ID
	if triggerLog.Topics[1] != mainchain.Bytes2Cid(trigger.ChannelId) {
		log.Errorf("%s. Trigger Tx channel ID not match, %x", logmsg, triggerLog.Topics[1])
		return true, false
	}

	// verify transaction block number
	if trigger.TriggerTxBlkNum != triggerLog.BlockNumber {
		log.Errorf("%s. TriggerTxBlkNum not match mainchain value: %d", logmsg, triggerLog.BlockNumber)
		return true, false
	}

	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) verifyGuardProof(change sync.Change) (done, approve bool) {
	var proof guard.GuardProof
	m.Transactor.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &proof)
	logmsg := fmt.Sprintf("verify change id %d, guard proof request: %s", change.ID, proof)

	// verify request exists
	r, err := guard.CLIQueryRequest(m.Transactor.CliCtx, guard.RouterKey, proof.ChannelId, proof.SimplexReceiver)
	if err != nil {
		log.Errorf("%s. query request err: %s", logmsg, err)
		return false, false
	}
	if r.TriggerTxHash == "" {
		log.Errorf("%s. Request Trigger event has not been submitted", logmsg)
		return true, false
	}

	// verify onchain trasaction receipt and status
	receipt, err := m.EthClient.Client.TransactionReceipt(context.Background(), mainchain.Hex2Hash(proof.GuardTxHash))
	if err != nil {
		log.Errorf("%s. Get trigger transaction receipt err: %s", logmsg, err)
		return false, false
	}
	if receipt.Status != mainchain.TxSuccess {
		log.Errorf("%s. Trigger transaction failed: %d", logmsg, receipt.Status)
		return true, false
	}
	guardLog := receipt.Logs[len(receipt.Logs)-1] // IntendSettle/SnapshotStates event is the last one

	// verify transaction contract address
	if guardLog.Address != m.EthClient.GetLedger().Address {
		log.Errorf("%s. Guard tx contract address not match: %x", logmsg, guardLog.Address)
		return true, false
	}

	// verify transaction event type
	if guardLog.Topics[0] != intendSettleEventSig && guardLog.Topics[0] != snapshotStatesEventSig {
		log.Errorf("%s. Guard Tx is not for IntendSettle/SnapshotStates event", logmsg)
		return true, false
	}

	// verify transaction channel ID
	if guardLog.Topics[1] != mainchain.Bytes2Cid(proof.ChannelId) {
		log.Errorf("%s. Guard Tx channel ID not match, %x", logmsg, guardLog.Topics[1])
		return true, false
	}

	// verify transaction block number
	if guardLog.BlockNumber <= r.TriggerTxBlkNum {
		log.Errorf("%s. Invalid block number for GuardTx at at %d", logmsg, guardLog.BlockNumber)
		return true, false
	}
	if guardLog.BlockNumber != proof.GuardTxBlkNum {
		log.Errorf("%s. GuardTxBlkNum not match mainchain value: %d", logmsg, guardLog.BlockNumber)
		return true, false
	}

	// verify channel sequence number and current request status
	seqIndex := 0
	if bytes.Compare(mainchain.Hex2Addr(r.SimplexSender).Bytes(), mainchain.Hex2Addr(r.SimplexReceiver).Bytes()) > 0 {
		seqIndex = 1
	}
	var seqNum uint64
	if r.Status == guard.ChanStatus_Settling {
		if proof.Status != guard.ChanStatus_Settled {
			log.Errorf("%s. Proof guard state should be settled", logmsg)
			return true, false
		}
		event, err2 := m.EthClient.GetLedger().ParseIntendSettle(*guardLog)
		if err2 != nil {
			log.Errorf("%s. ParseIntendSettle event err %s", logmsg, err2)
			return true, false
		}
		seqNum = event.SeqNums[seqIndex].Uint64()
	} else if r.Status == guard.ChanStatus_Withdrawing {
		if proof.Status != guard.ChanStatus_Idle {
			log.Errorf("%s. Proof guard state should be idle", logmsg)
			return true, false
		}
		event, err2 := m.EthClient.GetLedger().ParseSnapshotStates(*guardLog)
		if err2 != nil {
			log.Errorf("%s. ParseSnapshotStates event err %s", logmsg, err2)
			return true, false
		}
		seqNum = event.SeqNums[seqIndex].Uint64()
	} else {
		log.Errorf("%s. Current guard state is not settling or withdraw, %d", logmsg, r.Status)
		return true, false
	}
	if seqNum != r.SeqNum {
		log.Errorf("SeqNum not match, expected: %d, actual: %d", r.SeqNum, seqNum)
		return true, false
	}

	// verify guard sender
	guardSender, err := mainchain.GetTxSender(m.EthClient.Client, proof.GuardTxHash)
	if err != nil {
		log.Errorf("%s. GetTxSender err: %s", logmsg, err)
		return false, false
	}
	if proof.GuardSender != guardSender {
		log.Errorf("%s. GuardSender not match mainchain value: %s", logmsg, guardSender)
		return true, false
	}

	// TOTO: more check on request
	log.Infof("%s. success", logmsg)
	return true, true
}

func (m *Monitor) cmpBlkNum(blkNum uint64) int8 {
	currentBlkNum := m.getCurrentBlockNumber().Uint64()
	if currentBlkNum > blkNum {
		return 1
	} else if currentBlkNum < blkNum {
		return -1
	}
	return 0
}
