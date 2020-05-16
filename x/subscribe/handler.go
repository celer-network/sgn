package subscribe

import (
	"fmt"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	intendSettleEventSig = mainchain.GetEventSignature("IntendSettle(bytes32,uint256[2])")
)

// NewHandler returns a handler for "subscribe" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		logEntry := seal.NewMsgLog()
		var res *sdk.Result
		var err error
		switch msg := msg.(type) {
		case MsgGuardProof:
			res, err = handleMsgGuardProof(ctx, keeper, msg, logEntry)
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

// Handle a message to submit guard proof
// Currently only supports that the validator sends out a tx purely for one intendSettle.
// (not call it via a contract or put multiple calls in one tx)
func handleMsgGuardProof(ctx sdk.Context, keeper Keeper, msg MsgGuardProof, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.ChanInfo.ChanId = mainchain.Bytes2Hex(msg.ChannelId)
	logEntry.ChanInfo.PeerFrom = msg.PeerFrom
	logEntry.ChanInfo.GuardTxHash = msg.TxHash

	res := &sdk.Result{}
	request, found := keeper.GetRequest(ctx, msg.ChannelId, msg.PeerFrom)
	if !found {
		return res, fmt.Errorf("Cannot find request for channel ID")
	}

	if request.TriggerTxHash == "" {
		return res, fmt.Errorf("IntendSettle Trigger event has not been submitted")
	}

	triggerLog, err := ValidateIntendSettle(
		"Trigger", keeper.ethClient, mainchain.Hex2Hash(request.TriggerTxHash), mainchain.Bytes2Cid(msg.ChannelId))
	if err != nil {
		return res, err
	}

	guardLog, err := ValidateIntendSettle(
		"Guard", keeper.ethClient, mainchain.Hex2Hash(msg.TxHash), mainchain.Bytes2Cid(msg.ChannelId))
	if err != nil {
		return res, err
	}

	// check block numbers
	if guardLog.BlockNumber <= triggerLog.BlockNumber {
		return res, fmt.Errorf("Invalid block number for GuardTx %d TriggerTx %d", guardLog.BlockNumber, triggerLog.BlockNumber)
	}

	err = validateIntendSettleSeqNum(guardLog.Data, request.PeerFromIndex, request.SeqNum)
	if err != nil {
		return res, err
	}

	requestGuards := request.RequestGuards
	blockNumberDiff := guardLog.BlockNumber - triggerLog.BlockNumber
	guardIndex := (len(requestGuards) + 1) * int(blockNumberDiff) / int(request.DisputeTimeout)
	logEntry.ChanInfo.GuardIndex = uint32(guardIndex)
	logEntry.ChanInfo.TriggerBlknum = triggerLog.BlockNumber
	logEntry.ChanInfo.GuardBlknum = guardLog.BlockNumber

	var rewardValidator sdk.AccAddress
	if guardIndex < len(requestGuards) {
		rewardValidator = request.RequestGuards[guardIndex]
	} else {
		// get mainchain tx sender in the last stage for rewarding
		guardEthAddrStr, err := mainchain.GetTxSender(keeper.ethClient.Client, msg.TxHash)
		if err != nil {
			return res, fmt.Errorf("GetTxSender err: %s", err)
		}
		rewardCandidate, found := keeper.validatorKeeper.GetCandidate(ctx, guardEthAddrStr)
		if found {
			_, found = getAccAddrIndex(request.RequestGuards, rewardCandidate.Operator)
			if !found {
				rewardValidator = rewardCandidate.Operator
			}
		}

		guardIndex = len(requestGuards)
	}

	request.GuardTxHash = msg.TxHash
	keeper.SetRequest(ctx, request)

	// punish corresponding guards and reward corresponding validator
	for i := 0; i < guardIndex; i++ {
		keeper.slashKeeper.HandleGuardFailure(ctx, rewardValidator, request.RequestGuards[i])
	}

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}
