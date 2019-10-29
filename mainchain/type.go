package mainchain

import "math/big"

type CandidateInfo struct {
	Initialized   bool
	MinSelfStake  *big.Int
	SidechainAddr []byte
	StakingPool   *big.Int
	Status        *big.Int
	UnbondTime    *big.Int
}
