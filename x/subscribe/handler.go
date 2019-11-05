package subscribe

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	protobuf "github.com/golang/protobuf/proto"
)

// NewHandler returns a handler for "subscribe" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSubscribe:
			return handleMsgSubscribe(ctx, keeper, msg)
		case MsgRequestGuard:
			return handleMsgRequestGuard(ctx, keeper, msg)
		case MsgGuardProof:
			return handleMsgGuardProof(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized subscribe Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to subscribe
func handleMsgSubscribe(ctx sdk.Context, keeper Keeper, msg MsgSubscribe) sdk.Result {
	deposit, err := keeper.ethClient.Guard.SubscriptionDeposits(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, ethcommon.HexToAddress(msg.EthAddress))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query subscription desposit: %s", err)).Result()
	}

	subscription, found := keeper.GetSubscription(ctx, msg.EthAddress)
	if !found {
		subscription = NewSubscription(msg.EthAddress)
	}
	subscription.Deposit = sdk.NewIntFromBigInt(deposit)
	keeper.SetSubscription(ctx, subscription)
	return sdk.Result{}
}

// Handle a message to request guard
func handleMsgRequestGuard(ctx sdk.Context, keeper Keeper, msg MsgRequestGuard) sdk.Result {
	logger := ctx.Logger()

	subscription, found := keeper.GetSubscription(ctx, msg.EthAddress)
	if !found {
		return sdk.ErrInternal("Cannot find subscription").Result()
	}

	keeper.SetSubscription(ctx, subscription)

	var signedSimplexState chain.SignedSimplexState
	err := protobuf.Unmarshal(msg.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to unmarshal signedSimplexStateBytes: %s", err)).Result()
	}

	var simplexPaymentChannel entity.SimplexPaymentChannel
	err = protobuf.Unmarshal(signedSimplexState.SimplexState, &simplexPaymentChannel)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to unmarshal simplexState: %s", err)).Result()
	}

	request, err := getRequest(ctx, keeper, simplexPaymentChannel)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to get request: %s", err)).Result()
	}
	if simplexPaymentChannel.SeqNum < request.SeqNum {
		return sdk.ErrInternal("Seq Num must be larger than previous request").Result()
	}

	err = verifySignedSimplexStateSigs(request, signedSimplexState)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to verify sigs: %s", err)).Result()
	}

	request.SeqNum = simplexPaymentChannel.SeqNum
	request.SignedSimplexStateBytes = msg.SignedSimplexStateBytes
	keeper.SetRequest(ctx, simplexPaymentChannel.ChannelId, request)
	logger.Info("Set request", "channelId", simplexPaymentChannel.ChannelId, "request seqNum", request.SeqNum)

	return sdk.Result{}
}

// Handle a message to submit guard proof
// Currently only supports that the validator sends out a tx purely for one intendSettle. (not call it via a contract or put multiple calls in one tx)
func handleMsgGuardProof(ctx sdk.Context, keeper Keeper, msg MsgGuardProof) sdk.Result {
	logger := ctx.Logger()

	request, found := keeper.GetRequest(ctx, msg.ChannelId)
	if !found {
		return sdk.ErrInternal("Cannot find request").Result()
	}

	intendSettleEventSig := mainchain.GetEventSignature("IntendSettle(bytes32,uint256[2])")

	// validate triggerTx
	triggerReceipt, err := keeper.ethClient.Client.TransactionReceipt(context.Background(), ctype.Hex2Hash(msg.TriggerTxHash))
	if err != nil {
		logger.Error("TriggerTxHash is not found on mainchain", "error", err)
		return sdk.ErrInternal("TriggerTxHash not found on mainchain").Result()
	}
	if triggerReceipt.Status != ctype.TxSuccess {
		logger.Error("TriggerTx failed")
		return sdk.ErrInternal("TriggerTx failed").Result()
	}
	log := triggerReceipt.Logs[len(triggerReceipt.Logs)-1] // IntendSettle event is the last one
	// check ledger contract
	if log.Address != keeper.ethClient.LedgerAddress {
		logger.Error("TriggerTx is not associated with ledger contract")
		return sdk.ErrInternal("TriggerTx is not associated with ledger contract").Result()
	}
	// check event type
	if log.Topics[0] != intendSettleEventSig {
		logger.Error("TriggerTx is not for IntendSettle event")
		return sdk.ErrInternal("TriggerTx is not for IntendSettle event").Result()
	}
	// check channel ID
	if log.Topics[1] != ctype.Bytes2Hash(msg.ChannelId) {
		logger.Error("TriggerTx's channel ID is wrong")
		return sdk.ErrInternal("TriggerTx's channel ID is wrong").Result()
	}
	// record time and calculate the supposed guard of the guard tx
	// TODO: (issue) need to prevent using an out-of-date triggerTx, namely an old IntendSettle event
	//     can be done by requiring the triggerTx must be after the time of submitting the request guard?
	triggerBlkNum := log.BlockNumber

	// TODO: use a function to replace duplicate code
	// validate guardTx
	guardReceipt, err := keeper.ethClient.Client.TransactionReceipt(context.Background(), ctype.Hex2Hash(msg.GuardTxHash))
	if err != nil {
		logger.Error("GuardTxHash is not found on mainchain", "error", err)
		return sdk.ErrInternal("GuardTxHash not found on mainchain").Result()
	}
	if guardReceipt.Status != ctype.TxSuccess {
		logger.Error("GuardTx failed")
		return sdk.ErrInternal("GuardTx failed").Result()
	}
	log = guardReceipt.Logs[len(guardReceipt.Logs)-1] // IntendSettle event is the last one
	// check ledger contract
	if log.Address != keeper.ethClient.LedgerAddress {
		logger.Error("GuardTx is not associated with ledger contract")
		return sdk.ErrInternal("GuardTx is not associated with ledger contract").Result()
	}
	// check event type
	if log.Topics[0] != intendSettleEventSig {
		logger.Error("GuardTx is not for IntendSettle event")
		return sdk.ErrInternal("GuardTx is not for IntendSettle event").Result()
	}
	// check channel ID
	if log.Topics[1] != ctype.Bytes2Hash(msg.ChannelId) {
		logger.Error("GuardTx's channel ID is wrong")
		return sdk.ErrInternal("GuardTx's channel ID is wrong").Result()
	}
	// check guardIntendSettle sequence number
	ledgerABI, err := abi.JSON(strings.NewReader(mainchain.CelerLedgerABI))
	if err != nil {
		logger.Error("Failed to parse CelerLedgerABI", "error", err)
		return sdk.ErrInternal("Failed to parse CelerLedgerABI").Result()
	}
	var guardIntendSettle mainchain.CelerLedgerIntendSettle
	err = ledgerABI.Unpack(&guardIntendSettle, "IntendSettle", log.Data)
	if err != nil {
		logger.Error("Failed to unpack IntendSettle event", "error", err)
		return sdk.ErrInternal("Failed to unpack IntendSettle event").Result()
	}
	if guardIntendSettle.SeqNums[request.PeerFromIndex].Uint64() != request.SeqNum {
		logger.Error("guardIntendSettle's seqNum is different from triggerIntendSettle's seqNum")
		return sdk.ErrInternal("guardIntendSettle's seqNum is different from triggerIntendSettle's seqNum").Result()
	}

	// check mainchain tx sender for reward submitter in the last stage
	guardTx, _, err := keeper.ethClient.Client.TransactionByHash(context.Background(), ctype.Hex2Hash(msg.GuardTxHash))
	guardMsg, err := guardTx.AsMessage(ethtypes.NewEIP155Signer(guardTx.ChainId()))
	if err != nil {
		logger.Error("Failed to get guardMsg")
		return sdk.ErrInternal("Failed to get guardMsg").Result()
	}
	guardEthAddrStr := ctype.Addr2HexWithPrefix(guardMsg.From())

	// check time
	guardBlkNum := log.BlockNumber
	if guardBlkNum <= triggerBlkNum {
		logger.Error("GuardTx's block number is not larger than TriggerTx's block number")
		return sdk.ErrInternal("GuardTx's block number is not larger than TriggerTx's block number").Result()
	}

	// check supposed guards
	requestGuards := request.RequestGuards
	blockNumberDiff := guardBlkNum - triggerBlkNum
	// all guards before guardIndex will be punished
	guardIndex := uint64(len(requestGuards)+1) * blockNumberDiff / request.DisputeTimeout
	if guardIndex >= uint64(len(requestGuards)) {
		guardIndex = uint64(len(requestGuards))
	}

	// TODO: (issue) what if user submits a new stateproof while old state proof being guarded

	// Set tx hashes
	request.TriggerTxHash = msg.TriggerTxHash
	request.GuardTxHash = msg.GuardTxHash
	request.GuardEthAddress = guardEthAddrStr
	keeper.SetRequest(ctx, msg.ChannelId, request)

	// Punish corresponding guards correctly.
	for i := uint64(0); i < guardIndex; i++ {
		keeper.slashKeeper.HandleGuardFailure(ctx, msg.Sender, request.RequestGuards[i])
	}

	// TODO: reward submitter in the last stage (any can submit). How?

	return sdk.Result{}
}
