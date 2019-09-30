package types

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/sgn"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuf "github.com/golang/protobuf/proto"
)

type Sig struct {
	Signer string `json:"signer"`
	Sig    []byte `json:"sig"`
}

func NewSig(signer string, sig []byte) Sig {
	return Sig{
		Signer: signer,
		Sig:    sig,
	}
}

// implement fmt.Stringer
func (r Sig) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Sig: %x,`, r.Sig))
}

type Reward struct {
	MiningReward     sdk.Int `json:"miningReward"`
	ServiceReward    sdk.Int `json:"serviceReward"`
	RewardProtoBytes []byte  `json:"rewardProtoBytes"` // proto msg for reward snapshot from latest intendWithdraw
	Sigs             []Sig   `json:"sigs"`
}

func NewReward() Reward {
	return Reward{}
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
func (r Reward) InitateWithdraw() {
	rewardBytes, _ := protobuf.Marshal(&sgn.Reward{
		CumulativeMiningReward:  r.MiningReward.BigInt().Bytes(),
		CumulativeServiceReward: r.ServiceReward.BigInt().Bytes(),
	})

	r.RewardProtoBytes = rewardBytes
	r.Sigs = []Sig{}
}

// Add signature to reward sigs
func (r Reward) AddSig(sig []byte) (string, error) {
	signer, err := mainchain.RecoverSigner(r.RewardProtoBytes, sig)
	if err != nil {
		return "", err
	}

	signerAddr := signer.String()
	for _, sig := range r.Sigs {
		if sig.Signer == signerAddr {
			return "", errors.New("repeated signer")
		}
	}

	r.Sigs = append(r.Sigs, NewSig(signerAddr, sig))
	return signerAddr, nil
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
