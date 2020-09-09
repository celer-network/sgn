package types

import (
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

const RouterKey = ModuleName // this was defined in your key.go file

const (
	TypeMsgSetTransactors           = "set_transactors"
	TypeMsgEditCandidateDescription = "edit_candidate_description"
	TypeMsgWithdrawReward           = "withdraw_reward"
	TypeMsgSignReward               = "sign_reward"
)

type MsgSetTransactors struct {
	EthAddress  string           `json:"eth_address"`
	Transactors []sdk.AccAddress `json:"transactors"`
	Sender      sdk.AccAddress   `json:"sender"`
}

// NewMsgSetTransactors is a constructor function for MsgSetTransactors
func NewMsgSetTransactors(ethAddress string, transactors []sdk.AccAddress, sender sdk.AccAddress) MsgSetTransactors {
	return MsgSetTransactors{
		EthAddress:  mainchain.FormatAddrHex(ethAddress),
		Transactors: transactors,
		Sender:      sender,
	}
}

// Route should return the name of the module
func (msg MsgSetTransactors) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetTransactors) Type() string { return TypeMsgSetTransactors }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetTransactors) ValidateBasic() error {
	if msg.EthAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "EthAddress cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	for _, transactor := range msg.Transactors {
		if transactor.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, transactor.String())
		}
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetTransactors) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetTransactors) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

type MsgEditCandidateDescription struct {
	EthAddress  string              `json:"eth_address"`
	Description staking.Description `json:"description"`
	Sender      sdk.AccAddress      `json:"sender"`
}

// NewMsgEditCandidateDescription is a constructor function for MsgEditCandidateDescription
func NewMsgEditCandidateDescription(ethAddress string, description staking.Description, sender sdk.AccAddress) MsgEditCandidateDescription {
	return MsgEditCandidateDescription{
		EthAddress:  mainchain.FormatAddrHex(ethAddress),
		Description: description,
		Sender:      sender,
	}
}

// Route should return the name of the module
func (msg MsgEditCandidateDescription) Route() string { return RouterKey }

// Type should return the action
func (msg MsgEditCandidateDescription) Type() string { return TypeMsgEditCandidateDescription }

// ValidateBasic runs stateless checks on the message
func (msg MsgEditCandidateDescription) ValidateBasic() error {
	if msg.EthAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "EthAddress cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	if msg.Description == (staking.Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgEditCandidateDescription) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgEditCandidateDescription) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgWithdrawReward defines a SyncValidator message
type MsgWithdrawReward struct {
	EthAddress string         `json:"eth_address"`
	Sender     sdk.AccAddress `json:"sender"`
}

func NewMsgWithdrawReward(ethAddress string, sender sdk.AccAddress) MsgWithdrawReward {
	return MsgWithdrawReward{
		EthAddress: mainchain.FormatAddrHex(ethAddress),
		Sender:     sender,
	}
}

// Route should return the name of the module
func (msg MsgWithdrawReward) Route() string { return RouterKey }

// Type should return the action
func (msg MsgWithdrawReward) Type() string { return TypeMsgWithdrawReward }

// ValidateBasic runs stateless checks on the message
func (msg MsgWithdrawReward) ValidateBasic() error {
	if msg.EthAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "EthAddress cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgWithdrawReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgWithdrawReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgSignReward defines a SyncValidator message
type MsgSignReward struct {
	EthAddress string         `json:"eth_address"`
	Sig        []byte         `json:"sig"`
	Sender     sdk.AccAddress `json:"sender"`
}

func NewMsgSignReward(ethAddress string, sig []byte, sender sdk.AccAddress) MsgSignReward {
	return MsgSignReward{
		EthAddress: mainchain.FormatAddrHex(ethAddress),
		Sig:        sig,
		Sender:     sender,
	}
}

// Route should return the name of the module
func (msg MsgSignReward) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSignReward) Type() string { return TypeMsgSignReward }

// ValidateBasic runs stateless checks on the message
func (msg MsgSignReward) ValidateBasic() error {
	if msg.EthAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "EthAddress cannot be empty")
	}

	if len(msg.Sig) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Sig cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSignReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSignReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
