package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Delegator struct {
	DelegatedStake    *big.Int `json:"delegated_stake"`
	UndelegatingStake *big.Int `json:"undelegating_stake"`
}

type Candidate struct {
	StakingPool     *big.Int              `json:"staking_pool"`
	CommissionRate  *big.Int              `json:"commission_rate"`
	MinSelfStake    *big.Int              `json:"min_self_stake"`
	RateLockEndTime *big.Int              `json:"rate_lock_end_time"`
	Status          *big.Int              `json:"status"`
	UnbondTime      *big.Int              `json:"unbond_time"`
	Delegators      map[string]*Delegator `json:"delegators"`
}

type Snapshot struct {
	EndBlockNumber uint64                `json:"end_block_number"`
	Candidates     map[string]*Candidate `json:"candidates"`
}

const (
	startBlkNum   = uint64(8580124)
	queryInterval = uint64(200)

	prevSnapshotFlag = "prev-snapshot"
)

var (
	candidateMap = make(map[string]*Candidate)
)

func snapshotMainchain() error {
	snapshot := Snapshot{
		EndBlockNumber: startBlkNum,
		Candidates:     candidateMap,
	}
	prevSnapshotFile := viper.GetString(prevSnapshotFlag)
	if prevSnapshotFile != "" {
		prevSnapshotBytes, err := ioutil.ReadFile(prevSnapshotFile)
		if err != nil {
			return err
		}

		err = json.Unmarshal(prevSnapshotBytes, &snapshot)
		if err != nil {
			return err
		}
	}

	ethClient, err := common.NewEthClientFromConfig()
	if err != nil {
		return err
	}

	header, err := ethClient.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	latestBlock := header.Number.Uint64()
	dposContract := ethClient.DPoS
	err = snapshotInitializeCandidate(dposContract, snapshot.EndBlockNumber, latestBlock)
	if err != nil {
		return err
	}

	for cur := snapshot.EndBlockNumber; cur < latestBlock; {
		nextCur := cur + queryInterval
		err = snapshotDelegate(dposContract, cur, nextCur)
		if err != nil {
			return err
		}

		cur = nextCur
	}

	snapshotStaking(dposContract, header.Number)

	snapshot.EndBlockNumber = latestBlock
	res, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}

	file := fmt.Sprintf("%d-snapshot.json", latestBlock)
	err = ioutil.WriteFile(file, res, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func snapshotInitializeCandidate(dposContract *mainchain.DPoS, start, end uint64) error {
	log.Infof("Snapshot InitializeCandidate event from %d to %d", start, end)
	initializeCandidateIt, err := dposContract.FilterInitializeCandidate(&bind.FilterOpts{
		Start: start,
		End:   &end,
	}, []mainchain.Addr{})
	if err != nil {
		return err
	}

	defer initializeCandidateIt.Close()
	for initializeCandidateIt.Next() {
		initializeCandidate := initializeCandidateIt.Event
		candidate := &Candidate{
			MinSelfStake:    initializeCandidate.MinSelfStake,
			CommissionRate:  initializeCandidate.CommissionRate,
			RateLockEndTime: initializeCandidate.RateLockEndTime,
			StakingPool:     big.NewInt(0),
			Delegators:      make(map[string]*Delegator),
		}
		candidateMap[initializeCandidate.Candidate.String()] = candidate
	}

	return nil
}

func snapshotDelegate(dposContract *mainchain.DPoS, start, end uint64) error {
	log.Infof("Snapshot Delegate event from %d to %d", start, end)
	delegatedIt, err := dposContract.FilterDelegate(&bind.FilterOpts{
		Start: start,
		End:   &end,
	}, []mainchain.Addr{}, []mainchain.Addr{})
	if err != nil {
		return err
	}

	defer delegatedIt.Close()
	for delegatedIt.Next() {
		delegate := delegatedIt.Event
		candidate := candidateMap[delegate.Candidate.String()]
		_, ok := candidate.Delegators[delegate.Delegator.String()]
		if !ok {
			candidate.Delegators[delegate.Delegator.String()] = &Delegator{
				DelegatedStake:    big.NewInt(0),
				UndelegatingStake: big.NewInt(0),
			}
		}
	}

	return nil
}

func snapshotStaking(dposContract *mainchain.DPoS, blkNum *big.Int) error {
	for candidateAddr, candidate := range candidateMap {
		log.Infof("Snapshot candidate %s at %s", candidateAddr, blkNum)
		candidateInfo, err := dposContract.GetCandidateInfo(&bind.CallOpts{BlockNumber: blkNum}, mainchain.Hex2Addr(candidateAddr))
		if err != nil {
			return err
		}
		candidate.CommissionRate = candidateInfo.CommissionRate
		candidate.RateLockEndTime = candidateInfo.RateLockEndTime
		candidate.StakingPool = candidateInfo.StakingPool
		candidate.MinSelfStake = candidateInfo.MinSelfStake
		candidate.Status = candidateInfo.Status
		candidate.UnbondTime = candidateInfo.UnbondTime

		for delegatorAddr, delegator := range candidate.Delegators {
			log.Infof("Snapshot delegator %s of candidate %s at %s", delegatorAddr, candidateAddr, blkNum)
			delegatorInfo, err := dposContract.GetDelegatorInfo(&bind.CallOpts{BlockNumber: blkNum}, mainchain.Hex2Addr(candidateAddr), mainchain.Hex2Addr(delegatorAddr))
			if err != nil {
				return err
			}
			delegator.DelegatedStake = delegatorInfo.DelegatedStake
			delegator.UndelegatingStake = delegatorInfo.UndelegatingStake
		}
	}

	return nil
}

func SnapshotMainchainCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot-mainchain",
		Short: "Snapshot mainchain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return snapshotMainchain()
		},
	}
	cmd.Flags().String(prevSnapshotFlag, "", "Previous snapshot file")
	return cmd
}
