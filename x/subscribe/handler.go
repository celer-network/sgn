package subscribe

import (
	"fmt"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "subscribe" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		logEntry := seal.NewMsgLog()
		var res *sdk.Result
		var err error
		switch msg := msg.(type) {
		case MsgRequestGuard:
			res, err = handleMsgRequestGuard(ctx, keeper, msg, logEntry)
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
func handleMsgRequestGuard(ctx sdk.Context, keeper Keeper, msg MsgRequestGuard, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()

	peerToAddr, err := eth.RecoverSigner(msg.SignedSimplexStateBytes, msg.PeerToSig)
	if err != nil {
		return nil, fmt.Errorf("Failed to recover peerTo signer: %s", err)
	}

	logEntry.EthAddress = mainchain.Addr2Hex(peerToAddr)

	signedSimplexState, simplexChannel, err := common.UnmarshalSignedSimplexStateBytes(msg.SignedSimplexStateBytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal signedSimplexStateBytes: %s", err)
	}

	cid := mainchain.Bytes2Cid(simplexChannel.ChannelId)
	logEntry.ChanInfo.ChanId = mainchain.Cid2Hex(cid)
	logEntry.ChanInfo.SeqNum = simplexChannel.SeqNum

	request, found := keeper.GetRequest(
		ctx, simplexChannel.ChannelId, mainchain.Bytes2AddrHex(simplexChannel.PeerFrom))
	if !found {
		return nil, fmt.Errorf("Failed to get request")
	}

	if mainchain.Hex2Addr(request.GetPeerToAddress()) != peerToAddr {
		return nil, fmt.Errorf("PeerTo not match stored request: %s", request.GetPeerToAddress())
	}

	err = VerifySignedSimplexStateSigs(request, signedSimplexState)
	if err != nil {
		return nil, fmt.Errorf("Failed to verify sigs: %s", err)
	}

	if simplexChannel.SeqNum < request.SeqNum {
		return nil, fmt.Errorf("Seq num smaller than stored request %d", request.SeqNum)
	}

	err = keeper.ChargeRequestFee(ctx, request.GetPeerToAddress())
	if err != nil {
		return nil, fmt.Errorf("Failed to charge request fee: %s", err)
	}

	request.SeqNum = simplexChannel.SeqNum
	request.SignedSimplexStateBytes = msg.SignedSimplexStateBytes
	keeper.SetRequest(ctx, request)

	return &sdk.Result{}, nil
}
