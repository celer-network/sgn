package mainchain

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type ValidatorInfo struct {
	StakingPool     *big.Int
	MinSelfStake    *big.Int
	CommissionRate  *big.Int
	RateLockEndTime *big.Int
}

// return validators, totalStakes, quorumStakes
func (ethClient *EthClient) GetValidators() (map[Addr]*ValidatorInfo, *big.Int, *big.Int, error) {
	maxValidatorNum, err := ethClient.DPoS.GetUIntValue(&bind.CallOpts{}, big.NewInt(MaxValidatorNum))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Get maxValidatorNum err: %w", err)
	}
	validators := make(map[Addr]*ValidatorInfo)
	totalStakes := new(big.Int)
	var i int64
	for i = 0; i < maxValidatorNum.Int64(); i++ {
		addr, err2 := ethClient.DPoS.ValidatorSet(&bind.CallOpts{}, big.NewInt(i))
		if err2 != nil {
			return nil, nil, nil, fmt.Errorf("Get validator %d err: %w", i, err2)
		}
		if addr != ZeroAddr {
			info, err3 := ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, addr)
			if err3 != nil {
				return nil, nil, nil, fmt.Errorf("Get validator %x info err: %w", addr, err3)
			}
			validators[addr] = &ValidatorInfo{
				StakingPool:     info.StakingPool,
				MinSelfStake:    info.MinSelfStake,
				CommissionRate:  info.CommissionRate,
				RateLockEndTime: info.RateLockEndTime,
			}
			totalStakes = totalStakes.Add(totalStakes, info.StakingPool)
		}
	}
	quorumStakes := new(big.Int).Mul(totalStakes, big.NewInt(2))
	quorumStakes = quorumStakes.Quo(quorumStakes, big.NewInt(3))
	quorumStakes = quorumStakes.Add(quorumStakes, big.NewInt(1))

	return validators, totalStakes, quorumStakes, nil
}

// return checkedStakes, totalStakes, quorumStakes
func (ethClient *EthClient) CheckVotingPower(addrs []Addr) (*big.Int, *big.Int, *big.Int, error) {
	validators, totalStakes, quorumStakes, err := ethClient.GetValidators()
	if err != nil {
		return nil, nil, nil, err
	}

	checked := make(map[Addr]bool)
	checkedStakes := new(big.Int)
	for _, addr := range addrs {
		if checked[addr] {
			continue
		}
		validator, exist := validators[addr]
		if exist {
			checkedStakes = checkedStakes.Add(checkedStakes, validator.StakingPool)
		}
		checked[addr] = true
	}

	return checkedStakes, totalStakes, quorumStakes, nil
}
