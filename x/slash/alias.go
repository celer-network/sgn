package slash

import (
	"github.com/celer-network/sgn/x/slash/client/cli"
	"github.com/celer-network/sgn/x/slash/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

const (
	ModuleName                 = types.ModuleName
	RouterKey                  = types.RouterKey
	StoreKey                   = types.StoreKey
	QueryPenalty               = types.QueryPenalty
	QueryParameters            = types.QueryParameters
	AttributeKeyNonce          = types.AttributeKeyNonce
	AttributeValueGuardFailure = types.AttributeValueGuardFailure
)

var (
	ModuleCdc              = types.ModuleCdc
	RegisterCodec          = types.RegisterCodec
	NewAccountAmtPair      = types.NewAccountAmtPair
	NewAccountFractionPair  = types.NewAccountFractionPair
	NewPenalty             = types.NewPenalty
	NewQueryPenaltyParams  = types.NewQueryPenaltyParams
	NewMsgSignPenalty      = types.NewMsgSignPenalty
	EventTypeSlash         = slashingTypes.EventTypeSlash
	PenaltyNonceKey        = types.PenaltyNonceKey
	ActionPenalty          = types.ActionPenalty
	GetPenaltyKey          = types.GetPenaltyKey
	CLIQueryPenalty        = cli.QueryPenalty
	CLIQueryPenaltyRequest = cli.QueryPenaltyRequest
	DefaultParams          = types.DefaultParams
)

type (
	AccountAmtPair     = types.AccountAmtPair
	AccountFractionPair = types.AccountFractionPair
	Penalty            = types.Penalty
	Params             = types.Params
	QueryPenaltyParams = types.QueryPenaltyParams
	MsgSignPenalty     = types.MsgSignPenalty
)
