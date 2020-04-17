package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Governance message types and routes
const (
	TypeMsgVote         = "vote"
	TypeMsgSubmitChange = "submit_change"
)

// MsgSubmitChange defines a message to create a sync change
type MsgSubmitChange struct {
	Type   string         `json:"type" yaml:"type"`
	Data   []byte         `json:"data" yaml:"data"`     //  Initial deposit paid by sender. Must be strictly positive
	Sender sdk.AccAddress `json:"sender" yaml:"sender"` //  Address of the sender
}

// NewMsgSubmitChange creates a new MsgSubmitChange instance
func NewMsgSubmitChange(changeType string, data []byte, sender sdk.AccAddress) MsgSubmitChange {
	return MsgSubmitChange{changeType, data, sender}
}

// Route implements Msg
func (msg MsgSubmitChange) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgSubmitChange) Type() string { return TypeMsgSubmitChange }

// ValidateBasic implements Msg
func (msg MsgSubmitChange) ValidateBasic() error {
	if msg.Type == "" {
		return sdkerrors.Wrap(ErrInvalidChangeType, "missing type")
	}
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}
	if len(msg.Data) > 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidChangeData, "data length must be larger than 0")
	}

	return msg.Content.ValidateBasic()
}

// String implements the Stringer interface
func (msg MsgSubmitChange) String() string {
	return fmt.Sprintf(`Submit Change Message:
  Type:         %s
  Sender: %s
`, msg.Type, msg.Sender)
}

// GetSignBytes implements Msg
func (msg MsgSubmitChange) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgSubmitChange) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgVote defines a message to cast a vote
type MsgVote struct {
	ChangeID uint64         `json:"change_id" yaml:"change_id"` // ID of the change
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`       //  address of the sender
}

// NewMsgVote creates a message to cast a vote on an active change
func NewMsgVote(sender sdk.AccAddress, changeID uint64) MsgVote {
	return MsgVote{changeID, sender}
}

// Route implements Msg
func (msg MsgVote) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgVote) Type() string { return TypeMsgVote }

// ValidateBasic implements Msg
func (msg MsgVote) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgVote) String() string {
	return fmt.Sprintf(`Vote Message:
  Change ID: %d
`, msg.ChangeID)
}

// GetSignBytes implements Msg
func (msg MsgVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
