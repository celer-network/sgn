package types

import (
	"fmt"
	"strings"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/proto/sgn"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	protobuf "github.com/golang/protobuf/proto"
)

const (
	emptyAddr  = "0x0000000000000000000000000000000000000000"
	expireTime = ^uint64(0)
)

type AccountAmtPair struct {
	Account string  `json:"account"`
	Amount  sdk.Int `json:"amount"`
}

func NewAccountAmtPair(account string, amount sdk.Int) AccountAmtPair {
	return AccountAmtPair{
		Account: account,
		Amount:  amount,
	}
}

type AccountPercentPair struct {
	Account string  `json:"account"`
	Percent sdk.Dec `json:"percent"`
}

func NewAccountPercentPair(account string, percent sdk.Dec) AccountPercentPair {
	return AccountPercentPair{
		Account: account,
		Percent: percent,
	}
}

type Penalty struct {
	Nonce               uint64               `json:"nonce"`
	ValidatorAddr       string               `json:"validatorAddr"`
	Reason              string               `json:"reason"`
	PenalizedDelegators []AccountAmtPair     `json:"penalizedDelegators"`
	Beneficiaries       []AccountPercentPair `json:"beneficiaries"`
	PenaltyProtoBytes   []byte               `json:"penaltyProtoBytes"`
	Sigs                []common.Sig         `json:"sigs"`
}

func NewPenalty(nonce uint64, reason string, validatorAddr string) Penalty {
	return Penalty{
		Nonce:         nonce,
		Reason:        reason,
		ValidatorAddr: validatorAddr,
	}
}

// implement fmt.Stringer
func (p Penalty) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Nonce: %d, Reason: %s`, p.Nonce, p.Reason))
}

func (p Penalty) GenerateProtoBytes() {
	var penalizedDelegators []*sgn.AccountAmtPair
	var beneficiaries []*sgn.AccountAmtPair
	totalPenalty := sdk.ZeroInt()
	totalBeneficiary := sdk.ZeroInt()

	for _, penalizedDelegator := range p.PenalizedDelegators {
		totalPenalty = totalPenalty.Add(penalizedDelegator.Amount)
		penalizedDelegators = append(penalizedDelegators, &sgn.AccountAmtPair{
			Account: ethcommon.HexToAddress(penalizedDelegator.Account).Bytes(),
			Amt:     penalizedDelegator.Amount.BigInt().Bytes(),
		})
	}

	for _, beneficiary := range p.Beneficiaries {
		amt := beneficiary.Percent.MulInt(totalPenalty).TruncateInt()
		totalBeneficiary = totalBeneficiary.Add(amt)
		beneficiaries = append(beneficiaries, &sgn.AccountAmtPair{
			Account: ethcommon.HexToAddress(beneficiary.Account).Bytes(),
			Amt:     amt.BigInt().Bytes(),
		})
	}

	restPenalty := totalPenalty.Sub(totalBeneficiary)
	if restPenalty.IsPositive() {
		beneficiaries = append(beneficiaries, &sgn.AccountAmtPair{
			Account: ethcommon.HexToAddress(emptyAddr).Bytes(),
			Amt:     restPenalty.BigInt().Bytes(),
		})
	}

	penaltyBytes, _ := protobuf.Marshal(&sgn.Penalty{
		Nonce:               p.Nonce,
		ExpireTime:          expireTime,
		ValidatorAddress:    ethcommon.HexToAddress(p.ValidatorAddr).Bytes(),
		PenalizedDelegators: penalizedDelegators,
		Beneficiaries:       beneficiaries,
	})

	p.PenaltyProtoBytes = penaltyBytes
}

// Add signature to penalty sigs
func (p *Penalty) AddSig(sig []byte, expectedSigner string) error {
	sigs, err := common.AddSig(p.Sigs, p.PenaltyProtoBytes, sig, expectedSigner)
	if err != nil {
		return err
	}

	p.Sigs = sigs
	return nil
}

// Generate penaltyRequest msg
func (p Penalty) GetPenaltyRequest() []byte {
	var sigs [][]byte
	for _, sig := range p.Sigs {
		sigs = append(sigs, sig.Sig)
	}

	penaltyRequestBytes, _ := protobuf.Marshal(&sgn.PenaltyRequest{
		Penalty: p.PenaltyProtoBytes,
		Sigs:    sigs,
	})

	return penaltyRequestBytes
}
