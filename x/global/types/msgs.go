package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const RouterKey = ModuleName // this was defined in your key.go file

const (
	TypeMsgSyncBlock = "sync_block"
)

// MsgSyncBlock defines a SyncBlock message
type MsgSyncBlock struct {
	BlockNumber uint64         `json:"blockNumber"`
	Sender      sdk.AccAddress `json:"sender"`
}

// NewMsgSyncBlock is a constructor function for MsgSyncBlock
func NewMsgSyncBlock(blockNumber uint64, sender sdk.AccAddress) MsgSyncBlock {
	return MsgSyncBlock{
		BlockNumber: blockNumber,
		Sender:      sender,
	}
}

// Route should return the name of the module
func (msg MsgSyncBlock) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSyncBlock) Type() string { return TypeMsgSyncBlock }

// ValidateBasic runs stateless checks on the message
func (msg MsgSyncBlock) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSyncBlock) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSyncBlock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
