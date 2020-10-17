package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/sgn"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/protobuf/proto"
)

type Reward struct {
	Receiver         string       `json:"receiver"`
	MiningReward     sdk.Int      `json:"mining_reward"`
	ServiceReward    sdk.Int      `json:"service_reward"`
	RewardProtoBytes []byte       `json:"reward_proto_bytes"` // proto msg for reward snapshot from latest intendWithdraw
	LastWithdrawTime time.Time    `json:"last_withdraw_time"` // last time the user triggers withdraw
	Sigs             []common.Sig `json:"sigs"`
}

func NewReward(receiver string) Reward {
	return Reward{
		Receiver:      mainchain.FormatAddrHex(receiver),
		ServiceReward: sdk.ZeroInt(),
		MiningReward:  sdk.ZeroInt(),
	}
}

// implement fmt.Stringer
func (r Reward) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Receiver: %s, MiningReward: %v, ServiceReward: %v`, r.Receiver, r.MiningReward, r.ServiceReward))
}

// Initiate the withdraw process
func (r *Reward) InitateWithdraw(now time.Time) {
	rewardBytes, _ := proto.Marshal(
		&sgn.Reward{
			Receiver:                mainchain.Hex2Bytes(r.Receiver),
			CumulativeMiningReward:  r.MiningReward.BigInt().Bytes(),
			CumulativeServiceReward: r.ServiceReward.BigInt().Bytes(),
		})

	r.RewardProtoBytes = rewardBytes
	r.Sigs = []common.Sig{}
	r.LastWithdrawTime = now
}

// Add signature to reward sigs
func (r *Reward) AddSig(sig []byte, expectedSigner string) error {
	sigs, err := common.AddSig(r.Sigs, r.RewardProtoBytes, sig, expectedSigner)
	if err != nil {
		log.Error(err)
		return err
	}

	r.Sigs = sigs
	return nil
}

// Generate rewardRequest msg
func (r Reward) GetRewardRequest() []byte {
	var sigs [][]byte
	for _, sig := range r.Sigs {
		sigs = append(sigs, sig.Sig)
	}

	rewardRequestBytes, _ := proto.Marshal(&sgn.RewardRequest{
		Reward: r.RewardProtoBytes,
		Sigs:   sigs,
	})

	return rewardRequestBytes
}

type RewardEpoch struct {
	StartHeight   int64   `json:"start_height"`
	MiningReward  sdk.Int `json:"mining_reward"`
	ServiceReward sdk.Int `json:"service_reward"`
}

func NewRewardEpoch(height int64) RewardEpoch {
	return RewardEpoch{
		StartHeight:   height,
		MiningReward:  sdk.ZeroInt(),
		ServiceReward: sdk.ZeroInt(),
	}
}

type PendingReward struct {
	CandidateAddr string  `json:"candidate_addr"`
	MiningReward  sdk.Int `json:"mining_reward"`
	ServiceReward sdk.Int `json:"service_reward"`
}

func NewPendingReward(ethAddress string) PendingReward {
	return PendingReward{
		CandidateAddr: ethAddress,
		MiningReward:  sdk.ZeroInt(),
		ServiceReward: sdk.ZeroInt(),
	}
}

func (pr PendingReward) IsZero() bool {
	return pr.MiningReward.IsZero() && pr.ServiceReward.IsZero()
}

type RewardStats struct {
	TotalMiningReward  sdk.Int `json:"total_mining_reward"`
	TotalServiceReward sdk.Int `json:"total_service_reward"`
	NumReceiver        uint    `json:"num_receiver"`
	NumWithdrawer      uint    `json:"num_withdrawer"`
	MaxReward          sdk.Int `json:"max_reward"`
	MaxRewardReceiver  string  `json:"max_reward_receiver"`
}

func NewRewardStats() RewardStats {
	return RewardStats{
		TotalMiningReward:  sdk.ZeroInt(),
		TotalServiceReward: sdk.ZeroInt(),
		MaxReward:          sdk.ZeroInt(),
	}
}
