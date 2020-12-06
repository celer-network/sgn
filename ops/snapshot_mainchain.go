package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/celer-network/goutils/log"
	 "github.com/celer-network/sgn-contract/bindings/go/sgncontracts"
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
	prevSnapshotFlag  = "prev-snapshot"
	startBlkNumFlag   = "start-blk-number"
	queryIntervalFlag = "query-interval"
)

var (
	candidateMap              = make(map[string]*Candidate)
	candidateDelegatorChanges = make(map[string](map[string]bool))
)

func snapshotMainchain() error {
	snapshot := Snapshot{
		EndBlockNumber: viper.GetUint64(startBlkNumFlag),
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

	queryInterval := viper.GetUint64(queryIntervalFlag)
	for cur := snapshot.EndBlockNumber; cur < latestBlock; {
		nextCur := cur + queryInterval
		err = snapshotUpdateDelegatedStake(dposContract, cur, nextCur)
		if err != nil {
			return err
		}

		cur = nextCur
	}

	syncCandidateDelegators(dposContract, header.Number)

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

func snapshotInitializeCandidate(dposContract *sgncontracts.DPoS, start, end uint64) error {
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

func snapshotUpdateDelegatedStake(dposContract *sgncontracts.DPoS, start, end uint64) error {
	log.Infof("Snapshot UpdateDelegatedStake event from %d to %d", start, end)
	updateDelegatedStakeIt, err := dposContract.FilterUpdateDelegatedStake(&bind.FilterOpts{
		Start: start,
		End:   &end,
	}, []mainchain.Addr{}, []mainchain.Addr{})
	if err != nil {
		return err
	}

	defer updateDelegatedStakeIt.Close()
	for updateDelegatedStakeIt.Next() {
		updateDelegatedStake := updateDelegatedStakeIt.Event
		candidateAddr := updateDelegatedStake.Candidate.String()
		delegatorAddr := updateDelegatedStake.Delegator.String()
		candidate := candidateMap[candidateAddr]
		_, ok := candidate.Delegators[delegatorAddr]
		if !ok {
			candidate.Delegators[delegatorAddr] = &Delegator{
				DelegatedStake:    big.NewInt(0),
				UndelegatingStake: big.NewInt(0),
			}
		}

		_, ok = candidateDelegatorChanges[candidateAddr]
		if !ok {
			candidateDelegatorChanges[candidateAddr] = make(map[string]bool)
		}

		candidateDelegatorChanges[candidateAddr][delegatorAddr] = true
	}

	return nil
}

func syncCandidateDelegators(dposContract *sgncontracts.DPoS, blkNum *big.Int) error {
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
	}

	for candidateAddr, delegatorChanges := range candidateDelegatorChanges {
		candidate := candidateMap[candidateAddr]

		for delegatorAddr := range delegatorChanges {
			log.Infof("Snapshot delegator %s of candidate %s at %s", delegatorAddr, candidateAddr, blkNum)
			delegatorInfo, err := dposContract.GetDelegatorInfo(&bind.CallOpts{BlockNumber: blkNum}, mainchain.Hex2Addr(candidateAddr), mainchain.Hex2Addr(delegatorAddr))
			if err != nil {
				return err
			}

			delegator := candidate.Delegators[delegatorAddr]
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
	cmd.Flags().Uint64(startBlkNumFlag, 8580124, "Block number start to query")
	cmd.Flags().Uint64(queryIntervalFlag, 10000, "Number of blocks between each query")
	return cmd
}
