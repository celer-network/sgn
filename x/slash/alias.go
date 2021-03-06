package slash

import (
	"github.com/celer-network/sgn/x/slash/client/cli"
	"github.com/celer-network/sgn/x/slash/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

const (
	ModuleName                 = types.ModuleName
	RouterKey                  = types.RouterKey
	StoreKey                   = types.StoreKey
	QueryPenalty               = types.QueryPenalty
	QueryPenalties             = types.QueryPenalties
	QueryParameters            = types.QueryParameters
	AttributeKeyNonce          = types.AttributeKeyNonce
	AttributeValueGuardFailure = types.AttributeValueGuardFailure
	AttributeValueDepositBurn  = types.AttributeValueDepositBurn
)

var (
	ModuleCdc              = types.ModuleCdc
	RegisterCodec          = types.RegisterCodec
	NewAccountAmtPair      = types.NewAccountAmtPair
	NewAccountFractionPair = types.NewAccountFractionPair
	NewPenalty             = types.NewPenalty
	NewQueryPenaltyParams  = types.NewQueryPenaltyParams
	NewMsgSignPenalty      = types.NewMsgSignPenalty
	ActionPenalty          = types.ActionPenalty
	EventTypeSlash         = slashing.EventTypeSlash
	PenaltyNonceKey        = types.PenaltyNonceKey
	GetPenaltyKey          = types.GetPenaltyKey
	PenaltyKeyPrefix       = types.PenaltyKeyPrefix
	CLIQueryPenalty        = cli.QueryPenalty
	CLIQueryPenaltyRequest = cli.QueryPenaltyRequest
	DefaultParams          = types.DefaultParams
)

type (
	AccountAmtPair      = types.AccountAmtPair
	AccountFractionPair = types.AccountFractionPair
	Penalty             = types.Penalty
	Params              = types.Params
	QueryPenaltyParams  = types.QueryPenaltyParams
	MsgSignPenalty      = types.MsgSignPenalty
)
