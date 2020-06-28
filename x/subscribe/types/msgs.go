package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const RouterKey = ModuleName // this was defined in your key.go file

type MsgGuardRequest struct {
	SignedSimplexStateBytes []byte         `json:"signedSimplexStateBytes"`
	PeerToSig               []byte         `json:"peerToSig"`
	Sender                  sdk.AccAddress `json:"sender"`
}

// NewMsgGuardRequest is a constructor function for MsgGuardRequest
func NewMsgGuardRequest(signedSimplexStateBytes, peerToSig []byte, sender sdk.AccAddress) MsgGuardRequest {
	return MsgGuardRequest{
		SignedSimplexStateBytes: signedSimplexStateBytes,
		PeerToSig:               peerToSig,
		Sender:                  sender,
	}
}

// Route should return the name of the module
func (msg MsgGuardRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgGuardRequest) Type() string { return "request_guard" }

// ValidateBasic runs stateless checks on the message
func (msg MsgGuardRequest) ValidateBasic() error {
	if len(msg.SignedSimplexStateBytes) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "SignedSimplexStateBytes cannot be empty")
	}

	if len(msg.PeerToSig) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "PeerToSig cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgGuardRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgGuardRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
