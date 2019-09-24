package subscribe

import (
	"errors"
	"math/big"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func getRequest(ctx sdk.Context, keeper Keeper, simplexPaymentChannel proto.SimplexPaymentChannel) (Request, error) {
	request, found := keeper.GetRequest(ctx, simplexPaymentChannel.ChannelId)
	if !found {
		channelId := [32]byte{}
		copy(channelId[:], simplexPaymentChannel.ChannelId)
		addresses, seqNums, err := keeper.ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{
			BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
		}, channelId)
		if err != nil {
			return Request{}, err
		}

		peerAddresses := []string{addresses[0].String(), addresses[1].String()}
		peerFromAddress := ethcommon.BytesToAddress(simplexPaymentChannel.PeerFrom).String()
		var peerFromIndex uint
		if peerAddresses[0] == peerFromAddress {
			peerFromIndex = 0
		} else if peerAddresses[1] == peerFromAddress {
			peerFromIndex = 1
		} else {
			return Request{}, errors.New("peerFrom is not valid address")
		}

		request = NewRequest(seqNums[peerFromIndex].Uint64(), peerAddresses, peerFromIndex)
	}

	return request, nil
}

func verifySignedSimplexStateSigs(request Request, signedSimplexState proto.SignedSimplexState) error {
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
