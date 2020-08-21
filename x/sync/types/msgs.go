package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Sync message types and routes
const (
	TypeMsgApprove      = "approve_change"
	TypeMsgSubmitChange = "submit_change"

	ConfirmParamProposal = "confirm_param_proposal"
	UpdateSidechainAddr  = "update_sidechain_addr"
	SyncDelegator        = "sync_delegator"
	SyncValidator        = "sync_validator"
	Subscribe            = "subscribe"
	InitGuardRequest     = "init_guard_request"
	GuardTrigger         = "guard_trigger"
	GuardProof           = "guard_proof"
)

// MsgSubmitChange defines a message to create a sync change
type MsgSubmitChange struct {
	ChangeType string         `json:"changeType" yaml:"changeType"`
	Data       []byte         `json:"data" yaml:"data"`         // The serialized information to sync
	BlockNum   uint64         `json:"blockNum" yaml:"blockNum"` // Claimed mainchain block number when the infomation was queried
	Sender     sdk.AccAddress `json:"sender" yaml:"sender"`     // Address of the sender
}

// NewMsgSubmitChange creates a new MsgSubmitChange instance
func NewMsgSubmitChange(changeType string, data []byte, blockNum uint64, sender sdk.AccAddress) MsgSubmitChange {
	return MsgSubmitChange{changeType, data, blockNum, sender}
}

// Route implements Msg
func (msg MsgSubmitChange) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgSubmitChange) Type() string { return TypeMsgSubmitChange }

// ValidateBasic implements Msg
func (msg MsgSubmitChange) ValidateBasic() error {
	if msg.ChangeType == "" {
		return sdkerrors.Wrap(ErrInvalidChangeType, "missing type")
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	if len(msg.Data) == 0 {
		return sdkerrors.Wrap(ErrInvalidChangeData, "data length must be larger than 0")
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgSubmitChange) String() string {
	return fmt.Sprintf(`Submit Change Message: ChangeType:%s, BlockNum:%d, Sender:%s`,
		msg.ChangeType, msg.BlockNum, msg.Sender)
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

// MsgApprove defines a message to cast a vote
type MsgApprove struct {
	ChangeID uint64         `json:"change_id" yaml:"change_id"` // ID of the change
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`       // Address of the sender
}

// NewMsgApprove creates a message to cast a vote on an active change
func NewMsgApprove(changeID uint64, sender sdk.AccAddress) MsgApprove {
	return MsgApprove{changeID, sender}
}

// Route implements Msg
func (msg MsgApprove) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgApprove) Type() string { return TypeMsgApprove }

// ValidateBasic implements Msg
func (msg MsgApprove) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgApprove) String() string {
	return fmt.Sprintf("Approve Message: Change ID: %d", msg.ChangeID)
}

// GetSignBytes implements Msg
func (msg MsgApprove) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgApprove) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
