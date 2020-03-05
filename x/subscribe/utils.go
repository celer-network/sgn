package subscribe

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func getRequest(ctx sdk.Context, keeper Keeper, simplexPaymentChannel entity.SimplexPaymentChannel) (Request, error) {
	peerFromAddr := mainchain.Bytes2AddrHex(simplexPaymentChannel.PeerFrom)
	request, found := keeper.GetRequest(ctx, simplexPaymentChannel.ChannelId, peerFromAddr)
	if !found {
		channelId := mainchain.Bytes2Cid(simplexPaymentChannel.ChannelId)

		disputeTimeout, err := keeper.ethClient.Ledger.GetDisputeTimeout(&bind.CallOpts{
			BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
		}, channelId)
		if err != nil {
			return Request{}, fmt.Errorf("GetDisputeTimeout err: %s", err)
		}

		addresses, seqNums, err := keeper.ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{
			BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
		}, channelId)
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
		request = NewRequest(simplexPaymentChannel.ChannelId, seqNum, peerAddrs, peerFromIndex, disputeTimeout.Uint64())
	}

	return request, nil
}

func getRequestGuards(ctx sdk.Context, keeper Keeper) []sdk.AccAddress {
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

	for len(requestGuards) < requestGuardCount {
		candidate := validatorCandidates[len(requestGuards)]
		candidate.RequestCount = candidate.RequestCount.AddRaw(1)
		keeper.validatorKeeper.SetCandidate(ctx, candidate)

		requestGuards = append(requestGuards, sdk.AccAddress(candidate.Operator))
	}

	return requestGuards
}

// Make sure signature match peer addresses for the channel
func verifySignedSimplexStateSigs(request Request, signedSimplexState chain.SignedSimplexState) error {
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

func getAccAddrIndex(addresses []sdk.AccAddress, targetAddress sdk.AccAddress) (index int, found bool) {
	for i, v := range addresses {
		if v.Equals(targetAddress) {
			return i, true
		}
	}
	return 0, false
}

func validateIntendSettle(txType string, ethClient *mainchain.EthClient, txHash mainchain.HashType, cid mainchain.CidType) (*ethtypes.Log, error) {
	receipt, err := ethClient.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf(txType+"TxHash is not found on mainchain. Error: %w", err)
	}

	if receipt.Status != mainchain.TxSuccess {
		return nil, fmt.Errorf(txType+"Tx failed. Error: %w", err)
	}

	log := receipt.Logs[len(receipt.Logs)-1] // IntendSettle event is the last one

	// check ledger contract
	if log.Address != ethClient.LedgerAddress {
		return nil, fmt.Errorf(txType+"Tx is not associated with ledger contract. Error: %w", err)
	}

	// check event type
	if log.Topics[0] != intendSettleEventSig {
		return nil, fmt.Errorf(txType+"Tx is not for IntendSettle event. Error: %w", err)
	}

	// check channel ID
	if log.Topics[1] != cid {
		return nil, fmt.Errorf(txType+"Tx's channel ID is wrong. Error: %w", err)
	}

	return log, nil
}

func validateIntendSettleSeqNum(logDate []byte, seqNumIndex uint8, expectedNum uint64) error {
	ledgerABI, err := abi.JSON(strings.NewReader(mainchain.CelerLedgerABI))
	if err != nil {
		return fmt.Errorf("Failed to parse CelerLedgerABI: %w", err)
	}

	var intendSettle mainchain.CelerLedgerIntendSettle
	err = ledgerABI.Unpack(&intendSettle, "IntendSettle", logDate)
	if err != nil {
		return fmt.Errorf("Failed to unpack IntendSettle event: %w", err)
	}

	if intendSettle.SeqNums[seqNumIndex].Uint64() != expectedNum {
		return fmt.Errorf("Unexpected seqNum of IntendSettle event. SeqNumIndex: %d, expected: %d, actual: %d", seqNumIndex, expectedNum, intendSettle.SeqNums[seqNumIndex].Uint64())
	}

	return nil
}
