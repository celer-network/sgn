package types

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/proto/sgn"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuf "github.com/golang/protobuf/proto"
)

type Reward struct {
	MiningReward     sdk.Int      `json:"miningReward"`
	ServiceReward    sdk.Int      `json:"serviceReward"`
	RewardProtoBytes []byte       `json:"rewardProtoBytes"` // proto msg for reward snapshot from latest intendWithdraw
	Sigs             []common.Sig `json:"sigs"`
}

func NewReward() Reward {
	return Reward{
		ServiceReward: sdk.ZeroInt(),
		MiningReward:  sdk.ZeroInt(),
	}
}

// implement fmt.Stringer
func (r Reward) String() string {
	return strings.TrimSpace(fmt.Sprintf(`MiningReward: %v, ServiceReward: %v`, r.MiningReward, r.ServiceReward))
}

// Check if have new reward added
func (r Reward) HasNewReward() bool {
	if len(r.RewardProtoBytes) == 0 {
		return true
	}

	var reward sgn.Reward
	protobuf.Unmarshal(r.RewardProtoBytes, &reward)

	hasNewServiceReward := new(big.Int).SetBytes(reward.CumulativeServiceReward).Cmp(r.ServiceReward.BigInt()) != 0
	hasNewMingingReward := new(big.Int).SetBytes(reward.CumulativeMiningReward).Cmp(r.MiningReward.BigInt()) != 0
	return hasNewServiceReward || hasNewMingingReward
}

// Initiate the withdraw process
func (r *Reward) InitateWithdraw() {
	rewardBytes, _ := protobuf.Marshal(&sgn.Reward{
		CumulativeMiningReward:  r.MiningReward.BigInt().Bytes(),
		CumulativeServiceReward: r.ServiceReward.BigInt().Bytes(),
	})

	r.RewardProtoBytes = rewardBytes
	r.Sigs = []common.Sig{}
}

// Add signature to reward sigs
func (r *Reward) AddSig(sig []byte, expectedSigner string) error {
	sigs, err := common.AddSig(r.Sigs, r.RewardProtoBytes, sig, expectedSigner)
	if err != nil {
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

	rewardRequestBytes, _ := protobuf.Marshal(&sgn.RewardRequest{
		Reward: r.RewardProtoBytes,
		Sigs:   sigs,
	})

	return rewardRequestBytes
}
