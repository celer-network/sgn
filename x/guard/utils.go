package guard

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
			return fmt.Errorf("invalid eth signer %d %x %x", i, addrs[i], addr)
		}
	}

	return nil
}

func AssignGuards(ctx sdk.Context, keeper Keeper) []sdk.AccAddress {
	validatorCandidates := keeper.validatorKeeper.GetValidatorCandidates(ctx)
	sort.Slice(validatorCandidates, func(i, j int) bool {
		validatorCandidate0 := validatorCandidates[i]
		validatorCandidate1 := validatorCandidates[j]
		stakeReqRatio0 := validatorCandidate0.StakingPool.Quo(validatorCandidate0.RequestCount.AddRaw(1))
		stakeReqRatio1 := validatorCandidate1.StakingPool.Quo(validatorCandidate1.RequestCount.AddRaw(1))

		if !stakeReqRatio0.Equal(stakeReqRatio1) {
			return stakeReqRatio0.GT(stakeReqRatio1)
		}
		return validatorCandidate0.StakingPool.GT(validatorCandidate1.StakingPool)
	})

	requestGuardCount := int(keeper.RequestGuardCount(ctx))
	assignedGuards := []sdk.AccAddress{}
	for len(assignedGuards) < requestGuardCount && len(assignedGuards) < len(validatorCandidates) {
		candidate := validatorCandidates[len(assignedGuards)]
		candidate.RequestCount = candidate.RequestCount.AddRaw(1)
		keeper.validatorKeeper.SetCandidate(ctx, candidate)

		assignedGuards = append(assignedGuards, candidate.ValAccount)
	}

	return assignedGuards
}

func PrintSimplexChannel(simplex *entity.SimplexPaymentChannel) string {
	return fmt.Sprintf("cid: %x, from: %x, seq: %d", simplex.GetChannelId(), simplex.GetPeerFrom(), simplex.GetSeqNum())
}
