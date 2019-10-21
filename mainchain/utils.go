package mainchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func IsBonded(candidateInfo CandidateInfo) bool {
	return candidateInfo.Status.Uint64() == Bonded
}

func ParseStatus(candidateInfo CandidateInfo) sdk.BondStatus {
	switch candidateInfo.Status.Uint64() {
	case Bonded:
		return sdk.Bonded
	case Unbonding:
		return sdk.Unbonding
	case Unbonded:
		return sdk.Unbonded
	}

	return sdk.Unbonded
}
