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

type ParamProposal struct {
	Proposer     Addr
	Deposit      *big.Int
	VoteDeadline *big.Int
	Record       *big.Int
	NewValue     *big.Int
	Status       uint8
}

// return validators, totalStakes, quorumStakes
func (ethClient *EthClient) GetValidators() (map[Addr]*ValidatorInfo, *big.Int, *big.Int, error) {
	maxValidatorNum, err := ethClient.GetMaxValidatorNum()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Get maxValidatorNum err: %w", err)
	}
	validators := make(map[Addr]*ValidatorInfo)
	totalStakes := new(big.Int)
	var i uint64
	for i = 0; i < maxValidatorNum; i++ {
		addr, err2 := ethClient.DPoS.ValidatorSet(&bind.CallOpts{}, big.NewInt(int64(i)))
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

func (ethClient *EthClient) GetParamProposal(id int64) (*ParamProposal, error) {
	res, err := ethClient.DPoS.ParamProposals(&bind.CallOpts{}, big.NewInt(id))
	if err != nil {
		return nil, err
	}
	return &ParamProposal{
		Proposer:     res.Proposer,
		Deposit:      res.Deposit,
		VoteDeadline: res.VoteDeadline,
		Record:       res.Record,
		NewValue:     res.NewValue,
		Status:       res.Status,
	}, nil
}

func (ethClient *EthClient) GetParamProposalVote(id int64, voter Addr) (uint8, error) {
	return ethClient.DPoS.GetParamProposalVote(&bind.CallOpts{}, big.NewInt(id), voter)
}

// return yesVoters, noVoters, yesStakes, totalStakes, quorumStakes
func (ethClient *EthClient) GetParamProposalVotes(id int64) ([]Addr, []Addr, *big.Int, *big.Int, *big.Int, error) {
	validators, totalStakes, quorumStakes, err := ethClient.GetValidators()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	var yesVoters, noVoters []Addr
	yesStakes := new(big.Int)
	for addr, val := range validators {
		vote, err2 := ethClient.GetParamProposalVote(id, addr)
		if err2 != nil {
			return nil, nil, nil, nil, nil, err2
		}
		if vote == VoteYes {
			yesStakes = yesStakes.Add(yesStakes, val.StakingPool)
			yesVoters = append(yesVoters, addr)
		} else if vote == VoteNo {
			noVoters = append(noVoters, addr)
		}
	}
	return yesVoters, noVoters, yesStakes, totalStakes, quorumStakes, nil
}

func (ethClient *EthClient) GetProposalDeposit() (uint64, error) {
	return ethClient.GetUIntValue(ProposalDeposit)
}

func (ethClient *EthClient) GetGovernVoteTimeout() (uint64, error) {
	return ethClient.GetUIntValue(GovernVoteTimeout)
}

func (ethClient *EthClient) GetSlashTimeout() (uint64, error) {
	return ethClient.GetUIntValue(SlashTimeout)
}

func (ethClient *EthClient) GetMinValidatorNum() (uint64, error) {
	return ethClient.GetUIntValue(MinValidatorNum)
}

func (ethClient *EthClient) GetMaxValidatorNum() (uint64, error) {
	return ethClient.GetUIntValue(MaxValidatorNum)
}

func (ethClient *EthClient) GetMinStakeInPool() (uint64, error) {
	return ethClient.GetUIntValue(MinStakeInPool)
}

func (ethClient *EthClient) GetAdvanceNoticePeriod() (uint64, error) {
	return ethClient.GetUIntValue(AdvanceNoticePeriod)
}

func (ethClient *EthClient) GetMigrationTime() (uint64, error) {
	return ethClient.GetUIntValue(MigrationTime)
}

func (ethClient *EthClient) GetUIntValue(valtype int64) (uint64, error) {
	res, err := ethClient.DPoS.GetUIntValue(&bind.CallOpts{}, big.NewInt(valtype))
	if err != nil {
		return 0, err
	}
	return res.Uint64(), err
}
