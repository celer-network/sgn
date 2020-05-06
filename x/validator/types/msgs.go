package types

import (
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const RouterKey = ModuleName // this was defined in your key.go file

const (
	TypeMsgSetTransactors = "set_transactors"
	TypeMsgClaimValidator = "claim_validator"
	TypeMsgSyncValidator  = "sync_validator"
	TypeMsgSyncDelegator  = "sync_delegator"
	TypeMsgWithdrawReward = "withdraw_reward"
	TypeMsgSignReward     = "sign_reward"
)

type MsgSetTransactors struct {
	EthAddress  string           `json:"ethAddress"`
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

// MsgClaimValidator defines a ClaimValidator message
type MsgClaimValidator struct {
	EthAddress string         `json:"ethAddress"`
	PubKey     string         `json:"pubkey"`
	Sender     sdk.AccAddress `json:"sender"`
}

// NewMsgClaimValidator is a constructor function for MsgClaimValidator
func NewMsgClaimValidator(ethAddress string, pubkey string, sender sdk.AccAddress) MsgClaimValidator {
	return MsgClaimValidator{
		EthAddress: mainchain.FormatAddrHex(ethAddress),
		PubKey:     pubkey,
		Sender:     sender,
	}
}

// Route should return the name of the module
func (msg MsgClaimValidator) Route() string { return RouterKey }

// Type should return the action
func (msg MsgClaimValidator) Type() string { return TypeMsgClaimValidator }

// ValidateBasic runs stateless checks on the message
func (msg MsgClaimValidator) ValidateBasic() error {
	if msg.EthAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "EthAddress cannot be empty")
	}

	_, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, msg.PubKey)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgClaimValidator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgClaimValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgSyncValidator defines a SyncValidator message
type MsgSyncValidator struct {
	EthAddress string         `json:"ethAddress"`
	Sender     sdk.AccAddress `json:"sender"`
}

// NewMsgSyncValidator is a constructor function for MsgSyncValidator
func NewMsgSyncValidator(ethAddress string, sender sdk.AccAddress) MsgSyncValidator {
	return MsgSyncValidator{
		EthAddress: mainchain.FormatAddrHex(ethAddress),
		Sender:     sender,
	}
}

// Route should return the name of the module
func (msg MsgSyncValidator) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSyncValidator) Type() string { return TypeMsgSyncValidator }

// ValidateBasic runs stateless checks on the message
func (msg MsgSyncValidator) ValidateBasic() error {
	if msg.EthAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "EthAddress cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSyncValidator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSyncValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgWithdrawReward defines a SyncValidator message
type MsgWithdrawReward struct {
	EthAddress string         `json:"ethAddress"`
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
	EthAddress string         `json:"ethAddress"`
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
