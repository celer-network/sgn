package types

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
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
		CandidateAddress: ethcommon.HexToAddress(candidateAddress).String(),
		DelegatorAddress: ethcommon.HexToAddress(delegatorAddress).String(),
	}
}

type QueryCandidateParams struct {
	CandidateAddress string
}

func NewQueryCandidateParams(candidateAddress string) QueryCandidateParams {
	return QueryCandidateParams{
		CandidateAddress: ethcommon.HexToAddress(candidateAddress).String(),
	}
}

type QueryRewardParams struct {
	EthAddress string
}

func NewQueryRewardParams(ethAddress string) QueryRewardParams {
	return QueryRewardParams{
		EthAddress: ethcommon.HexToAddress(ethAddress).String(),
	}
}
