package subscribe

import (
	"errors"
	"math/big"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func getRequest(ctx sdk.Context, keeper Keeper, simplexPaymentChannel entity.SimplexPaymentChannel) (Request, error) {
	logger := ctx.Logger()

	request, found := keeper.GetRequest(ctx, simplexPaymentChannel.ChannelId)
	if !found {
		channelId := [32]byte{}
		copy(channelId[:], simplexPaymentChannel.ChannelId)

		disputeTimeout, err := keeper.ethClient.Ledger.GetDisputeTimeout(&bind.CallOpts{
			BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
		}, channelId)
		if err != nil {
			return Request{}, err
		}

		addresses, seqNums, err := keeper.ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{
			BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
		}, channelId)
		if err != nil {
			return Request{}, err
		}

		peerAddresses := []string{addresses[0].String(), addresses[1].String()}
		peerFromAddress := ethcommon.BytesToAddress(simplexPaymentChannel.PeerFrom).String()
		var peerFromIndex uint8
		if peerAddresses[0] == peerFromAddress {
			peerFromIndex = uint8(0)
		} else if peerAddresses[1] == peerFromAddress {
			peerFromIndex = uint8(1)
		} else {
			logger.Error("peerFrom is neither peerAddresses[0] nor peerAddresses[1]", "peerFrom", peerFromAddress, "peerAddresses[0]", peerAddresses[0], "peerAddresses[1]", peerAddresses[1], "channelId", channelId)
			return Request{}, errors.New("peerFrom is not valid address")
		}

		seqNum := seqNums[peerFromIndex].Uint64()
		requestGuards := getRequestGuards(ctx, keeper)
		request = NewRequest(seqNum, peerAddresses, peerFromIndex, disputeTimeout.Uint64(), requestGuards)
	}

	return request, nil
}

func getRequestGuards(ctx sdk.Context, keeper Keeper) []sdk.AccAddress {
	validators := keeper.validatorKeeper.GetValidators(ctx)
	requestGuardId := keeper.GetRequestGuardId(ctx)
	requestGuardCount := keeper.RequestGuardCount(ctx)
	requestGuards := []sdk.AccAddress{}

	for uint64(len(requestGuards)) < requestGuardCount {
		requestGuards = append(requestGuards, sdk.AccAddress(validators[requestGuardId].OperatorAddress))
		requestGuardId = (requestGuardId + 1) % uint8(len(validators))
	}

	keeper.SetRequestGuardId(ctx, requestGuardId)
	return requestGuards
}

// Make sure signature match peer addresses for the channel
func verifySignedSimplexStateSigs(request Request, signedSimplexState chain.SignedSimplexState) error {
	if len(signedSimplexState.Sigs) != 2 {
		return errors.New("incorrect sigs count")
	}

	for i := 0; i < 2; i++ {
		addr, err := mainchain.RecoverSigner(signedSimplexState.SimplexState, signedSimplexState.Sigs[0])
		if err != nil {
			return err
		}

		if request.PeerAddresses[0] != addr.String() {
			return errors.New("invalid sig")
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
