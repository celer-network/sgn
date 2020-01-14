package subscribe

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/thoas/go-funk"
)

func getRequest(ctx sdk.Context, keeper Keeper, simplexPaymentChannel entity.SimplexPaymentChannel) (Request, error) {
	request, found := keeper.GetRequest(ctx, simplexPaymentChannel.ChannelId)
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
		peerFromAddr := mainchain.Bytes2AddrHex(simplexPaymentChannel.PeerFrom)
		var peerFromIndex uint8
		if peerAddrs[0] == peerFromAddr {
			peerFromIndex = uint8(0)
		} else if peerAddrs[1] == peerFromAddr {
			peerFromIndex = uint8(1)
		} else {
			return Request{}, fmt.Errorf("Invalid peerFromAddr %s %s %s", peerFromAddr, peerAddrs[0], peerAddrs[1])
		}

		seqNum := seqNums[peerFromIndex].Uint64()
		requestGuards := getRequestGuards(ctx, keeper)
		request = NewRequest(seqNum, peerAddrs, peerFromIndex, disputeTimeout.Uint64(), requestGuards)
	}

	return request, nil
}

func getRequestGuards(ctx sdk.Context, keeper Keeper) []sdk.AccAddress {
	validators := keeper.validatorKeeper.GetValidators(ctx)
	validators = funk.Reverse(validators).([]staking.Validator)
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
			return fmt.Errorf("RecoverSigner err: %s", err)
		}

		if request.PeerAddresses[0] != mainchain.Addr2Hex(addr) {
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
