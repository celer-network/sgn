package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const RouterKey = ModuleName // this was defined in your key.go file

type MsgRequestGuard struct {
	SignedSimplexStateBytes []byte         `json:"signedSimplexStateBytes"`
	ReceiverSig               []byte         `json:"receiverSig"`
	Sender                  sdk.AccAddress `json:"sender"`
}

// NewMsgRequestGuard is a constructor function for MsgRequestGuard
func NewMsgRequestGuard(signedSimplexStateBytes, receiverSig []byte, sender sdk.AccAddress) MsgRequestGuard {
	return MsgRequestGuard{
		SignedSimplexStateBytes: signedSimplexStateBytes,
		ReceiverSig:               receiverSig,
		Sender:                  sender,
	}
}

// Route should return the name of the module
func (msg MsgRequestGuard) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestGuard) Type() string { return "request_guard" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRequestGuard) ValidateBasic() error {
	if len(msg.SignedSimplexStateBytes) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "SignedSimplexStateBytes cannot be empty")
	}

	if len(msg.ReceiverSig) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "ReceiverSig cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
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
