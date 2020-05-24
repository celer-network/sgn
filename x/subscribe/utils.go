package subscribe

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	coscontext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/protobuf/proto"
)

var (
	intendSettleEventSig   = mainchain.GetEventSignature("IntendSettle(bytes32,uint256[2])")
	intendWithdrawEventSig = mainchain.GetEventSignature("IntendWithdraw(bytes32,address,uint256)")
	snapshotStatesEventSig = mainchain.GetEventSignature("SnapshotStates(bytes32,uint256[2])")
)

func GetRequest(cliCtx coscontext.CLIContext, ethClient *mainchain.EthClient, signedSimplexState *chain.SignedSimplexState) (Request, error) {
	var simplexPaymentChannel entity.SimplexPaymentChannel
	err := proto.Unmarshal(signedSimplexState.SimplexState, &simplexPaymentChannel)
	if err != nil {
		return Request{}, fmt.Errorf("Failed to unmarshal simplexState: %s", err)
	}

	peerFromAddr := mainchain.Bytes2AddrHex(simplexPaymentChannel.PeerFrom)
	request, err := CLIQueryRequest(cliCtx, RouterKey, simplexPaymentChannel.ChannelId, peerFromAddr)
	if err != nil {
		channelId := mainchain.Bytes2Cid(simplexPaymentChannel.ChannelId)
		addresses, seqNums, err := ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{}, channelId)
		if err != nil {
			return Request{}, fmt.Errorf("GetStateSeqNumMap err: %s", err)
		}

		peerAddrs := []string{mainchain.Addr2Hex(addresses[0]), mainchain.Addr2Hex(addresses[1])}
		var peerFromIndex uint8
		if peerAddrs[0] == peerFromAddr {
			peerFromIndex = uint8(0)
		} else if peerAddrs[1] == peerFromAddr {
			peerFromIndex = uint8(1)
		} else {
			return Request{}, fmt.Errorf("Invalid peerFromAddr %s %s %s", peerFromAddr, peerAddrs[0], peerAddrs[1])
		}

		seqNum := seqNums[peerFromIndex].Uint64()
		request = NewRequest(simplexPaymentChannel.ChannelId, seqNum, peerAddrs, peerFromIndex)
	}

	return request, nil
}

// Make sure signature match peer addresses for the channel
func VerifySignedSimplexStateSigs(request Request, signedSimplexState chain.SignedSimplexState) error {
	if len(signedSimplexState.Sigs) != 2 {
		return fmt.Errorf("incorrect sigs count %d", len(signedSimplexState.Sigs))
	}

	for i := 0; i < 2; i++ {
		addr, err := mainchain.RecoverSigner(signedSimplexState.SimplexState, signedSimplexState.Sigs[i])
		if err != nil {
			return fmt.Errorf("RecoverSigner err: %s", err)
		}

		if request.PeerAddresses[i] != mainchain.Addr2Hex(addr) {
			return fmt.Errorf("invalid eth signer %d %s %s", i, request.PeerAddresses[i], mainchain.Addr2Hex(addr))
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

		requestGuards = append(requestGuards, sdk.AccAddress(candidate.Operator))
	}

	return requestGuards
}

func ValidateSnapshotSeqNum(logDate []byte, seqNumIndex uint8, expectedNum uint64) error {
	ledgerABI, err := abi.JSON(strings.NewReader(mainchain.CelerLedgerABI))
	if err != nil {
		return fmt.Errorf("Failed to parse CelerLedgerABI: %w", err)
	}

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
