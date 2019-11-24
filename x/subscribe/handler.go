package subscribe

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	protobuf "github.com/golang/protobuf/proto"
)

var (
	intendSettleEventSig = mainchain.GetEventSignature("IntendSettle(bytes32,uint256[2])")
)

// NewHandler returns a handler for "subscribe" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		logEntry := seal.NewMsgLog()
		var res sdk.Result
		var err error
		switch msg := msg.(type) {
		case MsgSubscribe:
			res, err = handleMsgSubscribe(ctx, keeper, msg, logEntry)
		case MsgRequestGuard:
			res, err = handleMsgRequestGuard(ctx, keeper, msg, logEntry)
		case MsgGuardProof:
			res, err = handleMsgGuardProof(ctx, keeper, msg, logEntry)
		default:
			errMsg := fmt.Sprintf("Unrecognized subscribe Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
		if err != nil {
			logEntry.Error = append(logEntry.Error, err.Error())
			seal.CommitMsgLog(logEntry)
			return sdk.ErrInternal(err.Error()).Result()
		}
		seal.CommitMsgLog(logEntry)
		return res
	}
}

// Handle a message to subscribe
func handleMsgSubscribe(ctx sdk.Context, keeper Keeper, msg MsgSubscribe, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	res := sdk.Result{}
	deposit, err := keeper.ethClient.Guard.SubscriptionDeposits(
		&bind.CallOpts{BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx))},
		mainchain.Hex2Addr(msg.EthAddress))
	if err != nil {
		return res, fmt.Errorf("Failed to query subscription desposit: %s", err)
	}

	subscription, found := keeper.GetSubscription(ctx, msg.EthAddress)
	if !found {
		subscription = NewSubscription(msg.EthAddress)
	}
	subscription.Deposit = sdk.NewIntFromBigInt(deposit)
	keeper.SetSubscription(ctx, subscription)
	return res, nil
}

// Handle a message to request guard
func handleMsgRequestGuard(ctx sdk.Context, keeper Keeper, msg MsgRequestGuard, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	res := sdk.Result{}
	err := keeper.ChargeRequestFee(ctx, msg.EthAddress)
	if err != nil {
		return res, fmt.Errorf("Failed to charge request fee: %s", err)
	}

	var signedSimplexState chain.SignedSimplexState
	err = protobuf.Unmarshal(msg.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		return res, fmt.Errorf("Failed to unmarshal signedSimplexStateBytes: %s", err)
	}

	var simplexPaymentChannel entity.SimplexPaymentChannel
	err = protobuf.Unmarshal(signedSimplexState.SimplexState, &simplexPaymentChannel)
	if err != nil {
		return res, fmt.Errorf("Failed to unmarshal simplexState: %s", err)
	}

	// reject guard request if the channel is not Operable
	// TODO: is this sufficient to handle the racing condition of one guard request and one IntendSettle event
	cid := mainchain.Bytes2Cid(simplexPaymentChannel.ChannelId)
	logEntry.ChanId = mainchain.Cid2Hex(cid)
	logEntry.ChanSeqNum = simplexPaymentChannel.SeqNum

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
	keeper.SetRequest(ctx, simplexPaymentChannel.ChannelId, request)

	return res, nil
}

// Handle a message to submit guard proof
// Currently only supports that the validator sends out a tx purely for one intendSettle.
// (not call it via a contract or put multiple calls in one tx)
func handleMsgGuardProof(ctx sdk.Context, keeper Keeper, msg MsgGuardProof, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.TriggerTxHash = msg.TriggerTxHash
	logEntry.GuardTxHash = msg.GuardTxHash
	logEntry.ChanId = mainchain.Bytes2Hex(msg.ChannelId)

	res := sdk.Result{}
	request, found := keeper.GetRequest(ctx, msg.ChannelId)
	if !found {
		return res, fmt.Errorf("Cannot find request for channel ID")
	}

	triggerLog, err := validateIntendSettle(
		"Trigger", keeper.ethClient, mainchain.Hex2Hash(msg.TriggerTxHash), mainchain.Bytes2Cid(msg.ChannelId))
	if err != nil {
		return res, err
	}

	guardLog, err := validateIntendSettle(
		"Guard", keeper.ethClient, mainchain.Hex2Hash(msg.GuardTxHash), mainchain.Bytes2Cid(msg.ChannelId))
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

	var rewardValidator sdk.AccAddress
	if guardIndex < len(requestGuards) {
		rewardValidator = request.RequestGuards[guardIndex]
	} else {
		// get mainchain tx sender in the last stage for rewarding
		guardEthAddrStr, err := mainchain.GetTxSender(keeper.ethClient.Client, msg.GuardTxHash)
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

	// set tx hashes
	request.TriggerTxHash = msg.TriggerTxHash
	request.GuardTxHash = msg.GuardTxHash
	keeper.SetRequest(ctx, msg.ChannelId, request)

	// punish corresponding guards and reward corresponding validator
	for i := 0; i < guardIndex; i++ {
		keeper.slashKeeper.HandleGuardFailure(ctx, rewardValidator, request.RequestGuards[i])
	}

	return res, nil
}

func validateIntendSettle(txType string, ethClient *mainchain.EthClient, txHash mainchain.HashType, cid mainchain.CidType) (*ethtypes.Log, error) {
	receipt, err := ethClient.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf(txType+"TxHash is not found on mainchain. Error: %w", err)
	}
	if receipt.Status != mainchain.TxSuccess {
		return nil, fmt.Errorf(txType+"Tx failed. Error: %w", err)
	}
	log := receipt.Logs[len(receipt.Logs)-1] // IntendSettle event is the last one

	// check ledger contract
	if log.Address != ethClient.LedgerAddress {
		return nil, fmt.Errorf(txType+"Tx is not associated with ledger contract. Error: %w", err)
	}
	// check event type
	if log.Topics[0] != intendSettleEventSig {
		return nil, fmt.Errorf(txType+"Tx is not for IntendSettle event. Error: %w", err)
	}
	// check channel ID
	if log.Topics[1] != cid {
		return nil, fmt.Errorf(txType+"Tx's channel ID is wrong. Error: %w", err)
	}

	return log, nil
}

func validateIntendSettleSeqNum(logDate []byte, seqNumIndex uint8, expectedNum uint64) error {
	ledgerABI, err := abi.JSON(strings.NewReader(mainchain.CelerLedgerABI))
	if err != nil {
		return fmt.Errorf("Failed to parse CelerLedgerABI: %w", err)
	}
	var intendSettle mainchain.CelerLedgerIntendSettle
	err = ledgerABI.Unpack(&intendSettle, "IntendSettle", logDate)
	if err != nil {
		return fmt.Errorf("Failed to unpack IntendSettle event: %w", err)
	}
	if intendSettle.SeqNums[seqNumIndex].Uint64() != expectedNum {
		return fmt.Errorf("Unexpected seqNum of IntendSettle event. SeqNumIndex: %d, expected: %d, actual: %d", seqNumIndex, expectedNum, intendSettle.SeqNums[seqNumIndex].Uint64())
	}

	return nil
}
