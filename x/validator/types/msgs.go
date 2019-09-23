package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const RouterKey = ModuleName // this was defined in your key.go file

const (
	TypeMsgInitializeCandidate = "initialize_candidate"
	TypeMsgClaimValidator      = "claim_validator"
	TypeMsgSyncValidator       = "sync_validator"
	TypeMsgSyncDelegator       = "sync_delegator"
	TypeMsgWithdrawReward      = "withdraw_reward"
)

// MsgInitializeCandidate defines a InitializeCandidate message
type MsgInitializeCandidate struct {
	EthAddress string         `json:"ethAddress"`
	Sender     sdk.AccAddress `json:"sender"`
}

// NewMsgInitializeCandidate is a constructor function for MsgInitializeCandidate
func NewMsgInitializeCandidate(ethAddress string, sender sdk.AccAddress) MsgInitializeCandidate {
	return MsgInitializeCandidate{
		EthAddress: ethcommon.HexToAddress(ethAddress).String(),
		Sender:     sender,
	}
}

// Route should return the name of the module
func (msg MsgInitializeCandidate) Route() string { return RouterKey }

// Type should return the action
func (msg MsgInitializeCandidate) Type() string { return TypeMsgInitializeCandidate }

// ValidateBasic runs stateless checks on the message
func (msg MsgInitializeCandidate) ValidateBasic() sdk.Error {
	if msg.EthAddress == "" {
		return sdk.ErrUnknownRequest("EthAddress cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgInitializeCandidate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgInitializeCandidate) GetSigners() []sdk.AccAddress {
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
		EthAddress: ethcommon.HexToAddress(ethAddress).String(),
		PubKey:     pubkey,
		Sender:     sender,
	}
}

// Route should return the name of the module
func (msg MsgClaimValidator) Route() string { return RouterKey }

// Type should return the action
func (msg MsgClaimValidator) Type() string { return TypeMsgClaimValidator }

// ValidateBasic runs stateless checks on the message
func (msg MsgClaimValidator) ValidateBasic() sdk.Error {
	if msg.EthAddress == "" {
		return sdk.ErrUnknownRequest("EthAddress cannot be empty")
	}

	_, err := sdk.GetConsPubKeyBech32(msg.PubKey)
	if err != nil {
		return sdk.ErrInvalidPubKey(err.Error())
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
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
		EthAddress: ethcommon.HexToAddress(ethAddress).String(),
		Sender:     sender,
	}
}

// Route should return the name of the module
func (msg MsgSyncValidator) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSyncValidator) Type() string { return TypeMsgSyncValidator }

// ValidateBasic runs stateless checks on the message
func (msg MsgSyncValidator) ValidateBasic() sdk.Error {
	if msg.EthAddress == "" {
		return sdk.ErrUnknownRequest("EthAddress cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
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

// MsgSyncDelegator defines a SyncDelegator message
type MsgSyncDelegator struct {
	CandidateAddress string         `json:"candidateAddress"`
	DelegatorAddress string         `json:"delegatorAddress"`
	Sender           sdk.AccAddress `json:"sender"`
}

// NewMsgSyncDelegator is a constructor function for MsgSyncDelegator
func NewMsgSyncDelegator(candidateAddress, delegatorAddress string, sender sdk.AccAddress) MsgSyncDelegator {
	return MsgSyncDelegator{
		CandidateAddress: ethcommon.HexToAddress(candidateAddress).String(),
		DelegatorAddress: ethcommon.HexToAddress(delegatorAddress).String(),
		Sender:           sender,
	}
}

// Route should return the name of the module
func (msg MsgSyncDelegator) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSyncDelegator) Type() string { return TypeMsgSyncDelegator }

// ValidateBasic runs stateless checks on the message
func (msg MsgSyncDelegator) ValidateBasic() sdk.Error {
	if msg.CandidateAddress == "" {
		return sdk.ErrUnknownRequest("CandidateAddress cannot be empty")
	}

	if msg.DelegatorAddress == "" {
		return sdk.ErrUnknownRequest("DelegatorAddress cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSyncDelegator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSyncDelegator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgWithdrawReward defines a SyncValidator message
type MsgWithdrawReward struct {
	EthAddress string         `json:"ethAddress"`
	Sender     sdk.AccAddress `json:"sender"`
}

func NewMsgWithdrawReward(ethAddress string, sender sdk.AccAddress) MsgWithdrawReward {
	return MsgWithdrawReward{
		EthAddress: ethcommon.HexToAddress(ethAddress).String(),
		Sender:     sender,
	}
}

// Route should return the name of the module
func (msg MsgWithdrawReward) Route() string { return RouterKey }

// Type should return the action
func (msg MsgWithdrawReward) Type() string { return TypeMsgWithdrawReward }

// ValidateBasic runs stateless checks on the message
func (msg MsgWithdrawReward) ValidateBasic() sdk.Error {
	if msg.EthAddress == "" {
		return sdk.ErrUnknownRequest("EthAddress cannot be empty")
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
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
