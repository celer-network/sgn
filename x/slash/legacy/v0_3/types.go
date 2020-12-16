package v03

import (
	"github.com/celer-network/sgn/mainchain"
	v02slash "github.com/celer-network/sgn/x/slash/legacy/v0_2"
	"github.com/celer-network/sgn/x/slash/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ModuleName = "slash"

type (
	Params = types.Params

	Penalty struct {
		Nonce               uint64                    `json:"nonce"`
		Reason              string                    `json:"reason"`
		ValidatorAddr       string                    `json:"validator_addr"`
		TotalPenalty        sdk.Int                   `json:"totalPenalty"`
		PenalizedDelegators []v02slash.AccountAmtPair `json:"penalized_delegators"`
		Beneficiaries       []v02slash.AccountAmtPair `json:"beneficiaries"`
		PenaltyProtoBytes   []byte                    `json:"penalty_proto_bytes"`
		Sigs                []v02slash.Sig            `json:"sigs"`
	}

	GenesisState struct {
		Params       Params    `json:"params" yaml:"params"`
		Penalties    []Penalty `json:"penalties" yaml:"penalties"`
		PenaltyNonce uint64    `json:"penalty_nonce" yaml:"penalty_nonce"`
	}
)

const (
	DefaultPenaltyDelegatorSize = types.DefaultPenaltyDelegatorSize
)

var (
	DefaultSyncerReward = types.DefaultSyncerReward
)

func NewPenalty(nonce uint64, reason string, validatorAddr string, penalizedDelegators []v02slash.AccountAmtPair,
	beneficiaryFractions []v02slash.AccountFractionPair, syncerReward sdk.Int) Penalty {
	var beneficiaries []v02slash.AccountAmtPair
	totalPenalty := sdk.ZeroInt()
	totalBeneficiary := sdk.ZeroInt()

	for _, penalizedDelegator := range penalizedDelegators {
		totalPenalty = totalPenalty.Add(penalizedDelegator.Amount)
	}

	restPenalty := totalPenalty
	for _, beneficiaryFraction := range beneficiaryFractions {
		amt := beneficiaryFraction.Fraction.MulInt(restPenalty).TruncateInt()
		totalBeneficiary = totalBeneficiary.Add(amt)
		beneficiaries = append(beneficiaries, v02slash.NewAccountAmtPair(beneficiaryFraction.Account, amt))
	}

	restPenalty = restPenalty.Sub(totalBeneficiary)
	if restPenalty.IsPositive() {
		beneficiaries = append(beneficiaries, v02slash.NewAccountAmtPair(mainchain.ZeroAddr.String(), restPenalty))
	}

	return Penalty{
		Nonce:               nonce,
		Reason:              reason,
		ValidatorAddr:       mainchain.FormatAddrHex(validatorAddr),
		TotalPenalty:        totalPenalty,
		Beneficiaries:       beneficiaries,
		PenalizedDelegators: penalizedDelegators,
	}
}
