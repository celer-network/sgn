package guard

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	intendSettleEventSig   = mainchain.GetEventSignature("IntendSettle(bytes32,uint256[2])")
	intendWithdrawEventSig = mainchain.GetEventSignature("IntendWithdraw(bytes32,address,uint256)")
	snapshotStatesEventSig = mainchain.GetEventSignature("SnapshotStates(bytes32,uint256[2])")
)

func VerifySimplexStateSigs(signedSimplexState *chain.SignedSimplexState, simplexSender, simplexReceiver mainchain.Addr) error {
	if len(signedSimplexState.Sigs) != 2 {
		return fmt.Errorf("incorrect sigs count %d", len(signedSimplexState.Sigs))
	}

	var addrs []mainchain.Addr
	if bytes.Compare(simplexSender.Bytes(), simplexReceiver.Bytes()) < 0 {
		addrs = []mainchain.Addr{simplexSender, simplexReceiver}
	} else {
		addrs = []mainchain.Addr{simplexReceiver, simplexSender}
	}

	for i := 0; i < 2; i++ {
		addr, err := eth.RecoverSigner(signedSimplexState.SimplexState, signedSimplexState.Sigs[i])
		if err != nil {
			return fmt.Errorf("RecoverSigner err: %s", err)
		}

		if addrs[i] != addr {
			return fmt.Errorf("invalid eth signer %d %s %x", i, addrs[i], addr)
		}
	}

	return nil
}

func ValidateTriggerTx(ethClient *mainchain.EthClient, txHash mainchain.HashType, cid mainchain.CidType) (*ethtypes.Log, error) {
	receipt, err := ethClient.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("Trigger TxHash is not found on mainchain. Error: %w", err)
	}

	if receipt.Status != mainchain.TxSuccess {
		return nil, fmt.Errorf("Trigger Tx failed. Error: %w", err)
	}

	log := receipt.Logs[len(receipt.Logs)-1] // IntendSettle/IntendWithdraw event is the last one

	// check ledger contract
	if log.Address != ethClient.LedgerAddress {
		return nil, fmt.Errorf("Trigger Tx is not associated with ledger contract. Error: %w", err)
	}

	// check event type
	if log.Topics[0] != intendSettleEventSig && log.Topics[0] != intendWithdrawEventSig {
		return nil, fmt.Errorf("Trigger Tx is not for IntendSettle/IntendWithdraw event. Error: %w", err)
	}

	// check channel ID
	if log.Topics[1] != cid {
		return nil, fmt.Errorf("Trigger Tx's channel ID is wrong. Error: %w", err)
	}

	return log, nil
}

func ValidateGuardTx(ethClient *mainchain.EthClient, txHash mainchain.HashType, cid mainchain.CidType) (*ethtypes.Log, error) {
	receipt, err := ethClient.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("Guard TxHash is not found on mainchain. Error: %w", err)
	}

	if receipt.Status != mainchain.TxSuccess {
		return nil, fmt.Errorf("Guard Tx failed. Error: %w", err)
	}

	log := receipt.Logs[len(receipt.Logs)-1] // SnapshotStates event is the last one

	// check ledger contract
	if log.Address != ethClient.LedgerAddress {
		return nil, fmt.Errorf("Guard Tx is not associated with ledger contract. Error: %w", err)
	}

	// check event type
	if log.Topics[0] != intendSettleEventSig && log.Topics[0] != snapshotStatesEventSig {
		return nil, fmt.Errorf("Guard Tx is not for IntendSettle/SnapshotStates event. Error: %w", err)
	}

	// check channel ID
	if log.Topics[1] != cid {
		return nil, fmt.Errorf("Guard Tx's channel ID is wrong. Error: %w", err)
	}

	return log, nil
}

func GetRequestGuards(ctx sdk.Context, keeper Keeper) []sdk.AccAddress {
	validatorCandidates := keeper.validatorKeeper.GetValidatorCandidates(ctx)
	sort.Slice(validatorCandidates, func(i, j int) bool {
		validatorCandidate0 := validatorCandidates[i]
		validatorCandidate1 := validatorCandidates[j]
		reqStakeRatio0 := validatorCandidate0.RequestCount.ToDec().QuoInt(validatorCandidate0.StakingPool)
		reqStakeRatio1 := validatorCandidate1.RequestCount.ToDec().QuoInt(validatorCandidate1.StakingPool)

		if !reqStakeRatio0.Equal(reqStakeRatio1) {
			return reqStakeRatio0.LT(reqStakeRatio1)
		}

		return validatorCandidate0.StakingPool.LT(validatorCandidate1.StakingPool)
	})

	requestGuardCount := int(keeper.RequestGuardCount(ctx))
	requestGuards := []sdk.AccAddress{}

	for len(requestGuards) < requestGuardCount && len(requestGuards) < len(validatorCandidates) {
		candidate := validatorCandidates[len(requestGuards)]
		candidate.RequestCount = candidate.RequestCount.AddRaw(1)
		keeper.validatorKeeper.SetCandidate(ctx, candidate)

		requestGuards = append(requestGuards, candidate.Operator)
	}

	return requestGuards
}

func ValidateGuardProofSeqNum(logDate []byte, seqNumIndex uint8, expectedNum uint64) error {
	ledgerABI, err := abi.JSON(strings.NewReader(mainchain.CelerLedgerABI))
	if err != nil {
		return fmt.Errorf("Failed to parse CelerLedgerABI: %w", err)
	}

	// TODO: currently relying on the fact the SnapshotStates and IntendSettle events have the same struct
	// should use different event parser.
	var snapshotStates mainchain.CelerLedgerSnapshotStates
	err = ledgerABI.Unpack(&snapshotStates, "SnapshotStates", logDate)
	if err != nil {
		return fmt.Errorf("Failed to unpack SnapshotStates event: %w", err)
	}

	if snapshotStates.SeqNums[seqNumIndex].Uint64() != expectedNum {
		return fmt.Errorf("Unexpected seqNum of SnapshotStates event. SeqNumIndex: %d, expected: %d, actual: %d", seqNumIndex, expectedNum, snapshotStates.SeqNums[seqNumIndex].Uint64())
	}

	return nil
}

func PrintSimplexChannel(simplex *entity.SimplexPaymentChannel) string {
	return fmt.Sprintf("cid: %x, from: %x, seq: %d", simplex.GetChannelId(), simplex.GetPeerFrom(), simplex.GetSeqNum())
}