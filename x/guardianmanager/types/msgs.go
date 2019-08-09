package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgSetName defines a SetName message
type MsgDeposit struct {
	EthAddress string         `json:"ethAddress"`
	Sender     sdk.AccAddress `json:"sender"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgDeposit(ethAddress string, sender sdk.AccAddress) MsgDeposit {
	return MsgDeposit{
		EthAddress: ethAddress,
		Sender:     sender,
	}
}

// Route should return the name of the module
func (msg MsgDeposit) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDeposit) Type() string { return "deposit" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDeposit) ValidateBasic() sdk.Error {
	if msg.EthAddress == "" {
		return sdk.ErrUnknownRequest("Eth adress cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
