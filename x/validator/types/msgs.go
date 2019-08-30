package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgClaimValidator defines a SetEthAddress message
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
func (msg MsgClaimValidator) Type() string { return "set_eth_address" }

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
