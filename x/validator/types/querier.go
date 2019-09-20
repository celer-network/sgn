package types

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	QueryPuller    = "puller"
	QueryPusher    = "pusher"
	QueryDelegator = "delegator"
	QueryCandidate = "candidate"
)

type QueryDelegatorParams struct {
	CandidateAddress string
	DelegatorAddress string
}

func NewQueryDelegatorParams(candidateAddress, delegatorAddress string) QueryDelegatorParams {
	return QueryDelegatorParams{
		CandidateAddress: ethcommon.HexToAddress(candidateAddress).String(),
		DelegatorAddress: ethcommon.HexToAddress(delegatorAddress).String(),
	}
}

type QueryCandidateParams struct {
	CandidateAddress string
	Seq              int64
}

func NewQueryCandidateParams(candidateAddress string, seq int64) QueryCandidateParams {
	return QueryCandidateParams{
		CandidateAddress: ethcommon.HexToAddress(candidateAddress).String(),
		Seq:              seq,
	}
}
