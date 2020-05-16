package monitor

import (
	"bytes"
	"reflect"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/protobuf/proto"
)

func (m *EthMonitor) verifyChange(change sync.Change) bool {
	switch change.Type {
	case sync.SyncBlock:
		return m.verifySyncBlock(change)
	case sync.Subscribe:
		return m.verifySubscribe(change)
	case sync.Request:
		return m.verifyRequest(change)
	case sync.IntendSettle:
		return m.verifyIntendSettle(change)
	case sync.GuardProof:
		return m.verifyGuardProof(change)
	case sync.UpdateSidechainAddr:
		return m.verifyUpdateSidechainAddr(change)
	case sync.SyncDelegator:
		return m.verifySyncDelegator(change)
	case sync.SyncValidator:
		return m.verifySyncValidator(change)
	default:
		return false
	}
}

func (m *EthMonitor) verifySyncBlock(change sync.Change) bool {
	var block global.Block
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &block)

	syncedBlock, err := m.getLatestBlock()
	if err != nil {
		return false
	}

	log.Infoln("Verify SyncBlock", block, syncedBlock)
	return block.Number <= m.blkNum.Uint64() && block.Number > syncedBlock.Number
}

func (m *EthMonitor) verifySubscribe(change sync.Change) bool {
	var subscription subscribe.Subscription
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &subscription)
	log.Infoln("Verify subscription", subscription)

	deposit, err := m.ethClient.SGN.SubscriptionDeposits(
		&bind.CallOpts{},
		mainchain.Hex2Addr(subscription.EthAddress))
	if err != nil {
		log.Errorf("Failed to query subscription desposit: %s", err)
		return false
	}

	return subscription.Deposit.BigInt().Cmp(deposit) == 0
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

	r, err := subscribe.GetRequest(m.operator.CliCtx, m.ethClient, &signedSimplexState)
	if err != nil {
		log.Errorln("Failed to get request through SignedSimplexStateBytes:", err)
		return false
	}

	err = subscribe.VerifySignedSimplexStateSigs(request, signedSimplexState)
	if err != nil {
		log.Infoln("Failed to verify sigs:", err)
		return false
	}

	return request.SeqNum > r.SeqNum && request.PeerFromIndex == r.PeerFromIndex &&
		bytes.Compare(request.ChannelId, r.ChannelId) && reflect.DeepEqual(request.PeerAddresses, r.PeerAddresses)
}

func (m *EthMonitor) verifyIntendSettle(change sync.Change) bool {
	var request subscribe.Request
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &request)
	log.Infoln("Verify IntendSettle", request)

	triggerLog, err := subscribe.ValidateIntendSettle(
		"Trigger", m.ethClient, mainchain.Hex2Hash(request.TriggerTxHash), mainchain.Bytes2Cid(request.ChannelId))
	if err != nil {
		log.Errorln(err)
		return false
	}

	disputeTimeout, err := m.ethClient.Ledger.GetDisputeTimeout(&bind.CallOpts{}, mainchain.Bytes2Cid(request.ChannelId))
	if err != nil {
		log.Errorln("GetDisputeTimeout err:", err)
		return false
	}

	return request.TriggerTxBlkNum == triggerLog.BlockNumber && request.DisputeTimeout == disputeTimeout.Uint64()
}

func (m *EthMonitor) verifyGuardProof(change sync.Change) bool {
	var request subscribe.Request
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &request)
	log.Infoln("Verify GuardProof", request)

	if request.TriggerTxHash == "" {
		log.Errorln("IntendSettle Trigger event has not been submitted")
		return false
	}

	guardLog, err := subscribe.ValidateIntendSettle(
		"Guard", m.ethClient, mainchain.Hex2Hash(request.GuardTxHash), mainchain.Bytes2Cid(request.ChannelId))
	if err != nil {
		log.Errorln(err)
		return false
	}

	if guardLog.BlockNumber <= request.TriggerTxBlkNum {
		log.Errorf("Invalid block number for GuardTx %d TriggerTx %d", guardLog.BlockNumber, request.TriggerTxBlkNum)
		return false
	}

	err = subscribe.ValidateIntendSettleSeqNum(guardLog.Data, request.PeerFromIndex, request.SeqNum)
	if err != nil {
		log.Errorln(err)
		return false
	}

	guardSender, err := mainchain.GetTxSender(m.ethClient.Client, request.GuardTxHash)
	if err != nil {
		log.Errorln(err)
		return false
	}

	return request.GuardTxBlkNum == guardLog.BlockNumber && request.GuardSender == guardSender
}

func (m *EthMonitor) verifyUpdateSidechainAddr(change sync.Change) bool {
	var candidate validator.Candidate
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &candidate)
	log.Infoln("Verify candidate", candidate)

	sidechainAddr, err := m.ethClient.SGN.SidechainAddrMap(&bind.CallOpts{}, mainchain.Hex2Addr(candidate.EthAddress))
	if err != nil {
		log.Errorln("Query sidechain ddress error:", err)
		return false
	}

	return candidate.Operator.Equals(sdk.AccAddress(sidechainAddr))
}

func (m *EthMonitor) verifySyncDelegator(change sync.Change) bool {
	var delegator validator.Delegator
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &delegator)
	log.Infoln("Verify delegator", delegator)

	di, err := m.ethClient.DPoS.GetDelegatorInfo(&bind.CallOpts{},
		mainchain.Hex2Addr(delegator.CandidateAddr), mainchain.Hex2Addr(delegator.DelegatorAddr))
	if err != nil {
		log.Errorf("Failed to query delegator info: %s", err)
		return false
	}

	return delegator.DelegatedStake.BigInt().Cmp(di.DelegatedStake) == 0
}

func (m *EthMonitor) verifySyncValidator(change sync.Change) bool {
	var validator staking.Validator
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(change.Data, &validator)
	log.Infoln("Verify validator", validator)

	ci, err := m.ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(validator.Description.Identity))
	if err != nil {
		log.Errorf("Failed to query candidate info: %s", err)
		return false
	}

	return validator.Status.Equal(mainchain.ParseStatus(ci)) &&
		validator.Tokens.Equal(sdk.NewIntFromBigInt(ci.StakingPool).QuoRaw(common.TokenDec))
}
