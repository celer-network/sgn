package types

import (
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgSubscribe defines a Subscribe message
type MsgSubscribe struct {
	EthAddress string         `json:"ethAddress"`
	Sender     sdk.AccAddress `json:"sender"`
}

// NewMsgSubscribe is a constructor function for MsgSubscribe
func NewMsgSubscribe(ethAddress string, sender sdk.AccAddress) MsgSubscribe {
	return MsgSubscribe{
		EthAddress: mainchain.FormatAddrHex(ethAddress),
		Sender:     sender,
	}
}

// Route should return the name of the module
func (msg MsgSubscribe) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSubscribe) Type() string { return "subscribe" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSubscribe) ValidateBasic() sdk.Error {
	if msg.EthAddress == "" {
		return sdk.ErrUnknownRequest("Eth adress cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSubscribe) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSubscribe) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

type MsgRequestGuard struct {
	EthAddress              string         `json:"ethAddress"`
	SignedSimplexStateBytes []byte         `json:"signedSimplexStateBytes"`
	Sender                  sdk.AccAddress `json:"sender"`
}

// NewMsgRequestGuard is a constructor function for MsgRequestGuard
func NewMsgRequestGuard(ethAddress string, signedSimplexStateBytes []byte, sender sdk.AccAddress) MsgRequestGuard {
	return MsgRequestGuard{
		EthAddress:              mainchain.FormatAddrHex(ethAddress),
		SignedSimplexStateBytes: signedSimplexStateBytes,
		Sender:                  sender,
	}
}

// Route should return the name of the module
func (msg MsgRequestGuard) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestGuard) Type() string { return "request_guard" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRequestGuard) ValidateBasic() sdk.Error {
	if len(msg.SignedSimplexStateBytes) == 0 {
		return sdk.ErrUnknownRequest("SignedSimplexStateBytes cannot be empty")
	}

	if msg.EthAddress == "" {
		return sdk.ErrUnknownRequest("EthAddress cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestGuard) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRequestGuard) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgGuardProof defines a Subscribe message
type MsgGuardProof struct {
	ChannelId     []byte         `json:"channelId"`
	PeerFrom      string         `json:"peerFrom"`
	TriggerTxHash string         `json:"triggerTxHash"` // intendSettle tx with lower sequence number
	GuardTxHash   string         `json:"guardTxHash"`   // intendSettle tx to guard user's state proof
	Sender        sdk.AccAddress `json:"sender"`
}

// NewMsgGuardProof is a constructor function for MsgGuardProof
func NewMsgGuardProof(channelId []byte, peerFrom string, triggerTxHash, guardTxHash string, sender sdk.AccAddress) MsgGuardProof {
	return MsgGuardProof{
		ChannelId:     channelId,
		PeerFrom:      mainchain.FormatAddrHex(peerFrom),
		TriggerTxHash: triggerTxHash,
		GuardTxHash:   guardTxHash,
		Sender:        sender,
	}
}

// Route should return the name of the module
func (msg MsgGuardProof) Route() string { return RouterKey }

// Type should return the action
func (msg MsgGuardProof) Type() string { return "guard_proof" }

// ValidateBasic runs stateless checks on the message
func (msg MsgGuardProof) ValidateBasic() sdk.Error {
	if len(msg.ChannelId) == 0 {
		return sdk.ErrUnknownRequest("channelId cannot be empty")
	}

	if msg.TriggerTxHash == "" {
		return sdk.ErrUnknownRequest("obsolete tx hash cannot be empty")
	}

	if msg.GuardTxHash == "" {
		return sdk.ErrUnknownRequest("guard tx hash cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgGuardProof) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgGuardProof) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
