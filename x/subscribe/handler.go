package subscribe

import (
	"fmt"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang/protobuf/proto"
)

// NewHandler returns a handler for "subscribe" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		logEntry := seal.NewMsgLog()
		var res *sdk.Result
		var err error
		switch msg := msg.(type) {
		case MsgGuardRequest:
			res, err = handleMsgGuardRequest(ctx, keeper, msg, logEntry)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
		if err != nil {
			logEntry.Error = append(logEntry.Error, err.Error())
		}

		seal.CommitMsgLog(logEntry)
		return res, err
	}
}

// Handle a message to request guard
func handleMsgGuardRequest(ctx sdk.Context, keeper Keeper, msg MsgGuardRequest, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	err := keeper.ChargeRequestFee(ctx, msg.EthAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to charge request fee: %s", err)
	}

	var signedSimplexState chain.SignedSimplexState
	err = proto.Unmarshal(msg.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal signedSimplexStateBytes: %s", err)
	}

	var simplexChannel entity.SimplexPaymentChannel
	err = proto.Unmarshal(signedSimplexState.SimplexState, &simplexChannel)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal simplexState: %s", err)
	}

	cid := mainchain.Bytes2Cid(simplexChannel.ChannelId)
	logEntry.ChanInfo.ChanId = mainchain.Cid2Hex(cid)
	logEntry.ChanInfo.SeqNum = simplexChannel.SeqNum

	request, found := keeper.GetRequest(
		ctx, simplexChannel.ChannelId, mainchain.Bytes2AddrHex(simplexChannel.PeerFrom))
	if !found {
		return nil, fmt.Errorf("Failed to get request: %s", err)
	}

	err = VerifySignedSimplexStateSigs(request, signedSimplexState)
	if err != nil {
		return nil, fmt.Errorf("Failed to verify sigs: %s", err)
	}

	if simplexChannel.SeqNum < request.SeqNum {
		return nil, fmt.Errorf("Seq num smaller than previous request %d", request.SeqNum)
	}

	ownerAddr, err := eth.RecoverSigner(request.SignedSimplexStateBytes, request.OwnerSig)
	if err != nil {
		return nil, fmt.Errorf("Failed to recover signer: %s", err)
	}

	if request.GetOwnerAddress() != mainchain.Addr2Hex(ownerAddr) {
		return nil, fmt.Errorf("Owner sig not match: %x", ownerAddr)
	}

	request.SeqNum = simplexChannel.SeqNum
	request.SignedSimplexStateBytes = msg.SignedSimplexStateBytes
	keeper.SetRequest(ctx, request)

	return &sdk.Result{}, nil
}
