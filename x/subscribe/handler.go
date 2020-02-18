package subscribe

import (
	"fmt"
	"math/big"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/protobuf/proto"
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
		case MsgSubscribe:
			res, err = handleMsgSubscribe(ctx, keeper, msg, logEntry)
		case MsgRequestGuard:
			res, err = handleMsgRequestGuard(ctx, keeper, msg, logEntry)
		case MsgIntendSettle:
			res, err = handleMsgIntendSettle(ctx, keeper, msg, logEntry)
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

// Handle a message to subscribe
func handleMsgSubscribe(ctx sdk.Context, keeper Keeper, msg MsgSubscribe, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	res := &sdk.Result{}
	deposit, err := keeper.ethClient.Guard.SubscriptionDeposits(
		&bind.CallOpts{BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx))},
		mainchain.Hex2Addr(msg.EthAddress))
	if err != nil {
		return res, fmt.Errorf("Failed to query subscription desposit: %s", err)
	}
	logEntry.Deposit = deposit.String()

	subscription, found := keeper.GetSubscription(ctx, msg.EthAddress)
	if !found {
		subscription = NewSubscription(msg.EthAddress)
	}
	subscription.Deposit = sdk.NewIntFromBigInt(deposit)
	keeper.SetSubscription(ctx, subscription)
	return res, nil
}

// Handle a message to request guard
func handleMsgRequestGuard(ctx sdk.Context, keeper Keeper, msg MsgRequestGuard, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	res := &sdk.Result{}
	err := keeper.ChargeRequestFee(ctx, msg.EthAddress)
	if err != nil {
		return res, fmt.Errorf("Failed to charge request fee: %s", err)
	}

	var signedSimplexState chain.SignedSimplexState
	err = proto.Unmarshal(msg.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		return res, fmt.Errorf("Failed to unmarshal signedSimplexStateBytes: %s", err)
	}

	var simplexPaymentChannel entity.SimplexPaymentChannel
	err = proto.Unmarshal(signedSimplexState.SimplexState, &simplexPaymentChannel)
	if err != nil {
		return res, fmt.Errorf("Failed to unmarshal simplexState: %s", err)
	}

	// reject guard request if the channel is not Operable
	// TODO: is this sufficient to handle the racing condition of one guard request and one IntendSettle event
	cid := mainchain.Bytes2Cid(simplexPaymentChannel.ChannelId)
	logEntry.ChanInfo.ChanId = mainchain.Cid2Hex(cid)
	logEntry.ChanInfo.SeqNum = simplexPaymentChannel.SeqNum

	status, err := keeper.ethClient.Ledger.GetChannelStatus(
		&bind.CallOpts{BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx))}, cid)
	if err != nil {
		return res, fmt.Errorf("Failed to query channel status: %s", err)
	}
	if status != mainchain.OperableChannel {
		return res, fmt.Errorf("Channel status is not Operable")
	}

	request, err := getRequest(ctx, keeper, simplexPaymentChannel)
	if err != nil {
		return res, fmt.Errorf("Failed to get request: %s", err)
	}
	if simplexPaymentChannel.SeqNum < request.SeqNum {
		return res, fmt.Errorf("Seq num smaller than previous request %d", request.SeqNum)
	}

	err = verifySignedSimplexStateSigs(request, signedSimplexState)
	if err != nil {
		return res, fmt.Errorf("Failed to verify sigs: %s", err)
	}

	request.SeqNum = simplexPaymentChannel.SeqNum
	request.SignedSimplexStateBytes = msg.SignedSimplexStateBytes
	keeper.SetRequest(ctx, request)

	return res, nil
}

func handleMsgIntendSettle(ctx sdk.Context, keeper Keeper, msg MsgIntendSettle, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.ChanInfo.ChanId = mainchain.Bytes2Hex(msg.ChannelId)
	logEntry.ChanInfo.PeerFrom = msg.PeerFrom
	logEntry.ChanInfo.TriggerTxHash = msg.TxHash

	res := &sdk.Result{}
	request, found := keeper.GetRequest(ctx, msg.ChannelId, msg.PeerFrom)
	if !found {
		return res, fmt.Errorf("Cannot find request for channel ID")
	}

	_, err := validateIntendSettle(
		"Trigger", keeper.ethClient, mainchain.Hex2Hash(msg.TxHash), mainchain.Bytes2Cid(msg.ChannelId))
	if err != nil {
		return res, err
	}

	request.TriggerTxHash = msg.TxHash
	request.RequestGuards = getRequestGuards(ctx, keeper)
	keeper.SetRequest(ctx, request)

	return res, nil
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

	triggerLog, err := validateIntendSettle(
		"Trigger", keeper.ethClient, mainchain.Hex2Hash(request.TriggerTxHash), mainchain.Bytes2Cid(msg.ChannelId))
	if err != nil {
		return res, err
	}

	guardLog, err := validateIntendSettle(
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
