package types

import (
	"github.com/celer-network/sgn/mainchain"
)

const (
	QuerySyncer     = "syncer"
	QueryDelegator  = "delegator"
	QueryCandidate  = "candidate"
	QueryReward     = "reward"
	QueryParameters = "parameters"
)

type QueryDelegatorParams struct {
	CandidateAddress string
	DelegatorAddress string
}

func NewQueryDelegatorParams(candidateAddress, delegatorAddress string) QueryDelegatorParams {
	return QueryDelegatorParams{
		CandidateAddress: mainchain.FormatAddrHex(candidateAddress),
		DelegatorAddress: mainchain.FormatAddrHex(delegatorAddress),
	}
}

type QueryCandidateParams struct {
	CandidateAddress string
}

func NewQueryCandidateParams(candidateAddress string) QueryCandidateParams {
	return QueryCandidateParams{
		CandidateAddress: mainchain.FormatAddrHex(candidateAddress),
	}
}

type QueryRewardParams struct {
	EthAddress string
}

func NewQueryRewardParams(ethAddress string) QueryRewardParams {
	return QueryRewardParams{
		EthAddress: mainchain.FormatAddrHex(ethAddress),
	}
}
