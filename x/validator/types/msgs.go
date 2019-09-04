package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgSyncValidator defines a SetEthAddress message
type MsgSyncValidator struct {
	EthAddress string         `json:"ethAddress"`
	PubKey     string         `json:"pubkey"`
	Sender     sdk.AccAddress `json:"sender"`
}

// NewMsgSyncValidator is a constructor function for MsgSyncValidator
func NewMsgSyncValidator(ethAddress string, pubkey string, sender sdk.AccAddress) MsgSyncValidator {
	return MsgSyncValidator{
		EthAddress: ethcommon.HexToAddress(ethAddress).String(),
		PubKey:     pubkey,
		Sender:     sender,
	}
}

// Route should return the name of the module
func (msg MsgSyncValidator) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSyncValidator) Type() string { return "set_eth_address" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSyncValidator) ValidateBasic() sdk.Error {
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
func (msg MsgSyncValidator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSyncValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
