package types

import (
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgGuardProof defines a Subscribe message
type MsgGuardProof struct {
	ChannelId []byte         `json:"channelId"`
	PeerFrom  string         `json:"peerFrom"`
	TxHash    string         `json:"txHash"` // intendSettle tx to guard user's state proof
	Sender    sdk.AccAddress `json:"sender"`
}

// NewMsgGuardProof is a constructor function for MsgGuardProof
func NewMsgGuardProof(channelId []byte, peerFrom string, txHash string, sender sdk.AccAddress) MsgGuardProof {
	return MsgGuardProof{
		ChannelId: channelId,
		PeerFrom:  mainchain.FormatAddrHex(peerFrom),
		TxHash:    txHash,
		Sender:    sender,
	}
}

// Route should return the name of the module
func (msg MsgGuardProof) Route() string { return RouterKey }

// Type should return the action
func (msg MsgGuardProof) Type() string { return "guard_proof" }

// ValidateBasic runs stateless checks on the message
func (msg MsgGuardProof) ValidateBasic() error {
	if len(msg.ChannelId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "channelId cannot be empty")
	}

	if msg.PeerFrom == "" {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "peerFrom cannot be empty")
	}

	if msg.TxHash == "" {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "guard tx hash cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
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
