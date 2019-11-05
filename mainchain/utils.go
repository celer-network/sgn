package mainchain

import (
	"github.com/celer-network/sgn/ctype"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
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

// GetEventSignature accepts the string of an event signature and return the hex
func GetEventSignature(eventSigStr string) ctype.HashType {
	return crypto.Keccak256Hash([]byte(eventSigStr))
}
