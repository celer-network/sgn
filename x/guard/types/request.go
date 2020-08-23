package types

import (
	"fmt"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ChanStatus uint8

const (
	// state of simplex channel guard request
	ChanStatus_Idle        ChanStatus = 0
	ChanStatus_Withdrawing ChanStatus = 1
	ChanStatus_Settling    ChanStatus = 2
	ChanStatus_Settled     ChanStatus = 3
)

func (status ChanStatus) String() string {
	switch status {
	case ChanStatus_Idle:
		return "Idle"
	case ChanStatus_Withdrawing:
		return "Withdraw"
	case ChanStatus_Settling:
		return "Settling"
	case ChanStatus_Settled:
		return "Settled"
	default:
		return "Invalid"
	}
}

type Request struct {
	ChannelId               []byte           `json:"channel_id"`
	SeqNum                  uint64           `json:"seq_num"`
	SimplexSender           string           `json:"simplex_sender"`
	SimplexReceiver         string           `json:"simplex_receiver"`
	SignedSimplexStateBytes []byte           `json:"signed_simplex_state_bytes"`
	DisputeTimeout          uint64           `json:"dispute_timeout"`
	Status                  ChanStatus       `json:"status"`
	AssignedGuards          []sdk.AccAddress `json:"assigned_guards"`
	TriggerTxHash           string           `json:"trigger_tx_hash"`
	TriggerTxBlkNum         uint64           `json:"trigger_tx_blk_num"`
	GuardTxHash             string           `json:"guard_tx_hash"`
	GuardTxBlkNum           uint64           `json:"guard_tx_blk_num"`
	GuardSender             string           `json:"guard_sender"`
}

func NewRequest(
	channelId []byte, seqNum uint64,
	senderAddr, receiverAddr mainchain.Addr,
	signedSimplexStateBytes []byte,
	disputeTimeout uint64) Request {
	return Request{
		ChannelId:               channelId,
		SeqNum:                  seqNum,
		SimplexSender:           mainchain.Addr2Hex(senderAddr),
		SimplexReceiver:         mainchain.Addr2Hex(receiverAddr),
		SignedSimplexStateBytes: signedSimplexStateBytes,
		DisputeTimeout:          disputeTimeout,
	}
}

func (r Request) String() string {
	out := fmt.Sprintf(`ChannelId: %x, SeqNum: %d, SimplexSender: %s, SimplexReceiver: %s, DisputeTimeout: %d, Status: %s`,
		r.ChannelId, r.SeqNum, r.SimplexSender, r.SimplexReceiver, r.DisputeTimeout, r.Status)
	if r.TriggerTxHash != "" {
		out += fmt.Sprintf(`, TriggerTxHash: %s, TriggerTxBlkNum: %d`,
			r.TriggerTxHash, r.TriggerTxBlkNum)
	}
	if r.GuardTxHash != "" {
		out += fmt.Sprintf(`, GuardTxHash: %s, GuardTxBlkNum: %d, GuardSender: %s`,
			r.GuardTxHash, r.GuardTxBlkNum, r.GuardSender)
	}
	return strings.TrimSpace(out)
}

type InitRequest struct {
	SignedSimplexStateBytes []byte `json:"signed_simplex_state_bytes"`
	SimplexReceiverSig      []byte `json:"simplex_receiver_sig"`
	DisputeTimeout          uint64 `json:"dispute_timeout"`
}

func NewInitRequest(signedSimplexStateBytes, simplexReceiverSig []byte, disputeTimeout uint64) *InitRequest {
	return &InitRequest{
		SignedSimplexStateBytes: signedSimplexStateBytes,
		SimplexReceiverSig:      simplexReceiverSig,
		DisputeTimeout:          disputeTimeout,
	}
}

type GuardTrigger struct {
	ChannelId       []byte     `json:"channel_id"`
	SimplexReceiver string     `json:"simplex_receiver"`
	TriggerTxHash   string     `json:"trigger_tx_hash"`
	TriggerTxBlkNum uint64     `json:"trigger_tx_blk_num"`
	TriggerSeqNum   uint64     `json:"trigger_seq_num"`
	Status          ChanStatus `json:"status"`
}

func (gt GuardTrigger) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ChannelId: %x, SimplexReceiver: %s, TriggerTxHash: %s, TriggerTxBlkNum: %d, TriggerSeqNum: %d, Status: %s`,
		gt.ChannelId, gt.SimplexReceiver, gt.TriggerTxHash, gt.TriggerTxBlkNum, gt.TriggerSeqNum, gt.Status))
}

func NewGuardTrigger(
	channelId mainchain.CidType,
	simplexReceiver mainchain.Addr,
	triggerTxHash mainchain.HashType,
	triggerTxBlkNum uint64,
	triggerSeqNum uint64,
	status ChanStatus) *GuardTrigger {
	return &GuardTrigger{
		ChannelId:       channelId.Bytes(),
		SimplexReceiver: mainchain.Addr2Hex(simplexReceiver),
		TriggerTxHash:   triggerTxHash.Hex(),
		TriggerTxBlkNum: triggerTxBlkNum,
		TriggerSeqNum:   triggerSeqNum,
		Status:          status,
	}
}

type GuardProof struct {
	ChannelId       []byte     `json:"channel_id"`
	SimplexReceiver string     `json:"simplex_receiver"`
	GuardTxHash     string     `json:"guard_tx_hash"`
	GuardTxBlkNum   uint64     `json:"guard_tx_blk_num"`
	GuardSender     string     `json:"guard_sender"`
	Status          ChanStatus `json:"status"`
}

func (gp GuardProof) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ChannelId: %x, SimplexReceiver: %s, GuardTxHash: %s, GuardTxBlkNum: %d, GuardSender: %s, Status: %s`,
		gp.ChannelId, gp.SimplexReceiver, gp.GuardTxHash, gp.GuardTxBlkNum, gp.GuardSender, gp.Status))
}

func NewGuardProof(
	channelId mainchain.CidType,
	simplexReceiver mainchain.Addr,
	guardTxHash mainchain.HashType,
	guardTxBlkNum uint64,
	guardSender mainchain.Addr,
	status ChanStatus) *GuardProof {
	return &GuardProof{
		ChannelId:       channelId.Bytes(),
		SimplexReceiver: mainchain.Addr2Hex(simplexReceiver),
		GuardTxHash:     guardTxHash.Hex(),
		GuardTxBlkNum:   guardTxBlkNum,
		GuardSender:     mainchain.Addr2Hex(guardSender),
		Status:          status,
	}
}
