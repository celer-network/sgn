package v02

import (
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ModuleName = "slash"

type (
	Params struct {
		SignedBlocksWindow        int64   `json:"signed_blocks_window" yaml:"signed_blocks_window"`
		MinSignedPerWindow        sdk.Dec `json:"min_signed_per_window" yaml:"min_signed_per_window"`
		SlashFractionDoubleSign   sdk.Dec `json:"slash_fraction_double_sign" yaml:"slash_fraction_double_sign"`
		SlashFractionDowntime     sdk.Dec `json:"slash_fraction_downtime" yaml:"slash_fraction_downtime"`
		SlashFractionGuardFailure sdk.Dec `json:"slash_fraction_guard_failure" yaml:"slash_fraction_guard_failure"`
		FallbackGuardReward       sdk.Dec `json:"fallback_guard_reward" yaml:"fallback_guard_reward"`
	}

	AccountAmtPair struct {
		Account string  `json:"account"`
		Amount  sdk.Int `json:"amount"`
	}

	AccountFractionPair struct {
		Account  string  `json:"account"`
		Fraction sdk.Dec `json:"percent"`
	}

	Sig struct {
		Signer string `json:"signer"`
		Sig    []byte `json:"sig"`
	}

	Penalty struct {
		Nonce               uint64                `json:"nonce"`
		ValidatorAddr       string                `json:"validator_addr"`
		Reason              string                `json:"reason"`
		PenalizedDelegators []AccountAmtPair      `json:"penalized_delegators"`
		Beneficiaries       []AccountFractionPair `json:"beneficiaries"`
		PenaltyProtoBytes   []byte                `json:"penalty_proto_bytes"`
		Sigs                []Sig                 `json:"sigs"`
	}

	GenesisState struct {
		Params       Params    `json:"params" yaml:"params"`
		Penalties    []Penalty `json:"penalties" yaml:"penalties"`
		PenaltyNonce uint64    `json:"penalty_nonce" yaml:"penalty_nonce"`
	}
)

func NewAccountAmtPair(account string, amount sdk.Int) AccountAmtPair {
	return AccountAmtPair{
		Account: mainchain.FormatAddrHex(account),
		Amount:  amount,
	}
}
