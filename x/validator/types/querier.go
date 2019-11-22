package types

import (
	"github.com/celer-network/sgn/mainchain"
)

const (
	QueryPuller    = "puller"
	QueryPusher    = "pusher"
	QueryDelegator = "delegator"
	QueryCandidate = "candidate"
	QueryReward    = "reward"
)

type QueryDelegatorParams struct {
	CandidateAddress string
	DelegatorAddress string
}

func NewQueryDelegatorParams(candidateAddress, delegatorAddress string) QueryDelegatorParams {
	return QueryDelegatorParams{
		CandidateAddress: mainchain.Hex2AddrHex(candidateAddress),
		DelegatorAddress: mainchain.Hex2AddrHex(delegatorAddress),
	}
}

type QueryCandidateParams struct {
	CandidateAddress string
}

func NewQueryCandidateParams(candidateAddress string) QueryCandidateParams {
	return QueryCandidateParams{
		CandidateAddress: mainchain.Hex2AddrHex(candidateAddress),
	}
}

type QueryRewardParams struct {
	EthAddress string
}

func NewQueryRewardParams(ethAddress string) QueryRewardParams {
	return QueryRewardParams{
		EthAddress: mainchain.Hex2AddrHex(ethAddress),
	}
}
