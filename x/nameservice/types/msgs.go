package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgSetName defines a SetName message
type MsgSetNumber struct {
	Sender sdk.AccAddress `json:"sender"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgSetNumber(sender sdk.AccAddress) MsgSetNumber {
	return MsgSetNumber{
		Sender: sender,
	}
}

// Route should return the name of the module
func (msg MsgSetNumber) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetNumber) Type() string { return "set_number" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetNumber) ValidateBasic() sdk.Error {
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetNumber) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetNumber) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
