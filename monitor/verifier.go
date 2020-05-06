package monitor

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (m *EthMonitor) verifyChange(change sync.Change) bool {
	switch change.Type {
	case sync.SyncBlock:
		return m.verifySyncBlock(change.Data)
	case sync.Subscribe:
		return m.verifySubscribe(change.Data)
	case sync.UpdateSidechainAddr:
		return m.verifyUpdateSidechainAddr(change.Data)
	case sync.SyncDelegator:
		return m.verifySyncDelegator(change.Data)
	case sync.SyncValidator:
		return m.verifySyncValidator(change.Data)
	default:
		return false
	}
}

func (m *EthMonitor) verifySyncBlock(data []byte) bool {
	var block global.Block
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(data, &block)

	syncedBlock, err := m.getLatestBlock()
	if err != nil {
		return false
	}

	log.Infoln("Verify SyncBlock", block, syncedBlock)
	return block.Number <= m.blkNum.Uint64() && block.Number > syncedBlock.Number
}

func (m *EthMonitor) verifySubscribe(data []byte) bool {
	var subscription subscribe.Subscription
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(data, &subscription)
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

func (m *EthMonitor) verifyUpdateSidechainAddr(data []byte) bool {
	var candidate validator.Candidate
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(data, &candidate)
	log.Infoln("Verify candidate", candidate)

	sidechainAddr, err := m.ethClient.SGN.SidechainAddrMap(&bind.CallOpts{}, mainchain.Hex2Addr(candidate.EthAddress))
	if err != nil {
		log.Errorln("Query sidechain ddress error:", err)
		return false
	}

	return candidate.Operator.Equals(sdk.AccAddress(sidechainAddr))
}

func (m *EthMonitor) verifySyncDelegator(data []byte) bool {
	var delegator validator.Delegator
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(data, &delegator)
	log.Infoln("Verify delegator", delegator)

	di, err := m.ethClient.DPoS.GetDelegatorInfo(&bind.CallOpts{},
		mainchain.Hex2Addr(delegator.CandidateAddr), mainchain.Hex2Addr(delegator.DelegatorAddr))
	if err != nil {
		log.Errorf("Failed to query delegator info: %s", err)
		return false
	}

	return delegator.DelegatedStake.BigInt().Cmp(di.DelegatedStake) == 0
}

func (m *EthMonitor) verifySyncValidator(data []byte) bool {
	var validator staking.Validator
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(data, &validator)
	log.Infoln("Verify validator", validator)

	ci, err := m.ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, validator.Description.Identity)
	if err != nil {
		log.Errorf("Failed to query candidate info: %s", err)
		return
	}

	di, err := m.ethClient.DPoS.GetDelegatorInfo(&bind.CallOpts{},
		mainchain.Hex2Addr(delegator.CandidateAddr), mainchain.Hex2Addr(delegator.DelegatorAddr))
	if err != nil {
		log.Errorf("Failed to query delegator info: %s", err)
		return false
	}

	return delegator.DelegatedStake.BigInt().Cmp(di.DelegatedStake) == 0
}
