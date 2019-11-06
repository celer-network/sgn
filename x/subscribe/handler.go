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
	"github.com/celer-network/sgn/x/subscribe/types"
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

	// reject guard request if the channel is not Operable
	cid := ctype.Bytes2Cid(simplexPaymentChannel.ChannelId)
	status, err := keeper.ethClient.Ledger.GetChannelStatus(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, cid)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query channel status: %s. Channel ID: %s", err, cid.Hex())).Result()
	}
	if status != uint8(1) {
		return sdk.ErrInternal(fmt.Sprintf("Channel status is not Operable. Channel ID: %s", cid.Hex())).Result()
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

	// validate triggerTx
	triggerLog, err := validateIntendSettle("Trigger", keeper.ethClient, ctype.Hex2Hash(msg.TriggerTxHash), ctype.Bytes2Cid(msg.ChannelId))
	if err != nil {
		logger.Error(err.Error())
		return sdk.ErrInternal(err.Error()).Result()
	}

	// validate guardTx
	guardLog, err := validateIntendSettle("Guard", keeper.ethClient, ctype.Hex2Hash(msg.GuardTxHash), ctype.Bytes2Cid(msg.ChannelId))
	if err != nil {
		logger.Error(err.Error())
		return sdk.ErrInternal(err.Error()).Result()
	}

	// check block numbers
	if guardLog.BlockNumber <= triggerLog.BlockNumber {
		logger.Error("GuardTx's block number is not larger than TriggerTx's block number")
		return sdk.ErrInternal("GuardTx's block number is not larger than TriggerTx's block number").Result()
	}

	// check guardIntendSettle sequence number
	ledgerABI, err := abi.JSON(strings.NewReader(mainchain.CelerLedgerABI))
	if err != nil {
		logger.Error("Failed to parse CelerLedgerABI", "error", err)
		return sdk.ErrInternal("Failed to parse CelerLedgerABI").Result()
	}
	var guardIntendSettle mainchain.CelerLedgerIntendSettle
	err = ledgerABI.Unpack(&guardIntendSettle, "IntendSettle", guardLog.Data)
	if err != nil {
		logger.Error("Failed to unpack IntendSettle event", "error", err)
		return sdk.ErrInternal("Failed to unpack IntendSettle event").Result()
	}
	if guardIntendSettle.SeqNums[request.PeerFromIndex].Uint64() != request.SeqNum {
		logger.Error("guardIntendSettle's seqNum is different from triggerIntendSettle's seqNum")
		return sdk.ErrInternal("guardIntendSettle's seqNum is different from triggerIntendSettle's seqNum").Result()
	}

	// get mainchain tx sender for rewarding guard in the last stage
	guardTx, _, err := keeper.ethClient.Client.TransactionByHash(context.Background(), ctype.Hex2Hash(msg.GuardTxHash))
	guardMsg, err := guardTx.AsMessage(ethtypes.NewEIP155Signer(guardTx.ChainId()))
	if err != nil {
		logger.Error("Failed to get guardMsg")
		return sdk.ErrInternal("Failed to get guardMsg").Result()
	}
	guardEthAddrStr := ctype.Addr2HexWithPrefix(guardMsg.From())

	// set tx hashes
	request.TriggerTxHash = msg.TriggerTxHash
	request.GuardTxHash = msg.GuardTxHash
	request.GuardEthAddress = guardEthAddrStr
	keeper.SetRequest(ctx, msg.ChannelId, request)

	// get supposed guards
	requestGuards := request.RequestGuards
	blockNumberDiff := guardLog.BlockNumber - triggerLog.BlockNumber
	// all guards before guardIndex will be punished
	guardIndex := (len(requestGuards) + 1) * int(blockNumberDiff) / int(request.DisputeTimeout)

	// punish corresponding guards and reward corresponding validator
	var rewardValidator sdk.AccAddress
	if guardIndex < len(requestGuards) {
		rewardValidator = request.RequestGuards[guardIndex]
	} else {
		rewardCandidate, found := keeper.validatorKeeper.GetCandidate(ctx, guardEthAddrStr)
		if found {
			_, found = getAccAddrIndex(request.RequestGuards, rewardCandidate.Operator)
			if !found {
				rewardValidator = rewardCandidate.Operator
			}
		}

		guardIndex = len(requestGuards)
	}
	for i := 0; i < guardIndex; i++ {
		keeper.slashKeeper.HandleGuardFailure(ctx, rewardValidator, request.RequestGuards[i])
	}

	return sdk.Result{}
}

func validateIntendSettle(txType string, ethClient *mainchain.EthClient, txHash ctype.HashType, cid ctype.CidType) (*ethtypes.Log, error) {
	receipt, err := ethClient.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return &ethtypes.Log{}, fmt.Errorf(txType+"TxHash is not found on mainchain. Error: %w", err)
	}
	if receipt.Status != ctype.TxSuccess {
		return &ethtypes.Log{}, fmt.Errorf(txType+"Tx failed. Error: %w", err)
	}
	log := receipt.Logs[len(receipt.Logs)-1] // IntendSettle event is the last one

	// check ledger contract
	if log.Address != ethClient.LedgerAddress {
		return &ethtypes.Log{}, fmt.Errorf(txType+"Tx is not associated with ledger contract. Error: %w", err)
	}
	// check event type
	if log.Topics[0] != types.IntendSettleEventSig {
		return &ethtypes.Log{}, fmt.Errorf(txType+"Tx is not for IntendSettle event. Error: %w", err)
	}
	// check channel ID
	if log.Topics[1] != cid {
		return &ethtypes.Log{}, fmt.Errorf(txType+"Tx's channel ID is wrong. Error: %w", err)
	}

	return log, nil
}
