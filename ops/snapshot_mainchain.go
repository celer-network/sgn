package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

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
	err = snapshotInitializeCandidate(dposContract)
	if err != nil {
		return err
	}

	err = snapshotUpdateCommissionRate(dposContract)
	if err != nil {
		return err
	}

	err = snapshotUpdateMinSelfStake(dposContract)
	if err != nil {
		return err
	}

	for cur := snapshot.EndBlockNumber; cur < latestBlock; {
		nextCur := cur + queryInterval
		err = snapshotUpdateDelegatedStake(dposContract, cur, nextCur)
		if err != nil {
			return err
		}

		cur = nextCur
	}

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

func snapshotInitializeCandidate(dposContract *mainchain.DPoS) error {
	initializeCandidateIt, err := dposContract.FilterInitializeCandidate(&bind.FilterOpts{}, []mainchain.Addr{})
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

func snapshotUpdateCommissionRate(dposContract *mainchain.DPoS) error {
	updateCommissionRateIt, err := dposContract.FilterUpdateCommissionRate(&bind.FilterOpts{}, []mainchain.Addr{})
	if err != nil {
		return err
	}

	defer updateCommissionRateIt.Close()
	for updateCommissionRateIt.Next() {
		updateCommissionRate := updateCommissionRateIt.Event
		candidate := candidateMap[updateCommissionRate.Candidate.String()]
		candidate.CommissionRate = updateCommissionRate.NewRate
		candidate.RateLockEndTime = updateCommissionRate.NewLockEndTime
	}

	return nil
}

func snapshotUpdateMinSelfStake(dposContract *mainchain.DPoS) error {
	minSelfStakeIt, err := dposContract.FilterUpdateMinSelfStake(&bind.FilterOpts{}, []mainchain.Addr{})
	if err != nil {
		return err
	}

	defer minSelfStakeIt.Close()
	for minSelfStakeIt.Next() {
		minSelfStake := minSelfStakeIt.Event
		candidate := candidateMap[minSelfStake.Candidate.String()]
		candidate.MinSelfStake = minSelfStake.MinSelfStake
	}

	return nil
}

func snapshotUpdateDelegatedStake(dposContract *mainchain.DPoS, start, end uint64) error {
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
		candidate := candidateMap[updateDelegatedStake.Candidate.String()]
		candidate.StakingPool = updateDelegatedStake.CandidatePool

		delegator := getDelegator(candidate, updateDelegatedStake.Delegator.String())
		delegator.DelegatedStake = updateDelegatedStake.DelegatorStake
	}

	return nil
}

func getDelegator(candidate *Candidate, delegatorAddr string) *Delegator {
	delegator, ok := candidate.Delegators[delegatorAddr]
	if !ok {
		delegator = &Delegator{
			DelegatedStake:    big.NewInt(0),
			UndelegatingStake: big.NewInt(0),
		}

		candidate.Delegators[delegatorAddr] = delegator
	}

	return delegator
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
