package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

type MsgSignPenalty struct {
	Nonce  uint64         `json:"nonce"`
	Sig    []byte         `json:"sig"`
	Sender sdk.AccAddress `json:"sender"`
}

func NewMsgSignPenalty(nonce uint64, sig []byte, sender sdk.AccAddress) MsgSignPenalty {
	return MsgSignPenalty{
		Nonce:  nonce,
		Sig:    sig,
		Sender: sender,
	}
}

// Route should return the name of the module
func (msg MsgSignPenalty) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSignPenalty) Type() string { return "sign_penalty" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSignPenalty) ValidateBasic() sdk.Error {
	if len(msg.Sig) == 0 {
		return sdk.ErrUnknownRequest("Sig cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSignPenalty) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSignPenalty) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
