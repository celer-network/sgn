package types

import (
	"fmt"
	"strings"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/sgn"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		Account: mainchain.FormatAddrHex(account),
		Amount:  amount,
	}
}

func (amp AccountAmtPair) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Account: %s, Amount: %v`, amp.Account, amp.Amount))
}

type AccountFractionPair struct {
	Account  string  `json:"account"`
	Fraction sdk.Dec `json:"percent"`
}

func NewAccountFractionPair(account string, fraction sdk.Dec) AccountFractionPair {
	return AccountFractionPair{
		Account:  mainchain.FormatAddrHex(account),
		Fraction: fraction,
	}
}

type Penalty struct {
	Nonce               uint64                `json:"nonce"`
	ValidatorAddr       string                `json:"validatorAddr"`
	Reason              string                `json:"reason"`
	PenalizedDelegators []AccountAmtPair      `json:"penalizedDelegators"`
	Beneficiaries       []AccountFractionPair `json:"beneficiaries"`
	PenaltyProtoBytes   []byte                `json:"penaltyProtoBytes"`
	Sigs                []common.Sig          `json:"sigs"`
}

func NewPenalty(nonce uint64, reason string, validatorAddr string) Penalty {
	return Penalty{
		Nonce:         nonce,
		Reason:        reason,
		ValidatorAddr: mainchain.FormatAddrHex(validatorAddr),
	}
}

// implement fmt.Stringer
func (p Penalty) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Nonce: %d, ValidatorAddr: %s, Reason: %s`, p.Nonce, p.ValidatorAddr, p.Reason))
}

func (p Penalty) GenerateProtoBytes() {
	var penalizedDelegators []*sgn.AccountAmtPair
	var beneficiaries []*sgn.AccountAmtPair
	totalPenalty := sdk.ZeroInt()
	totalBeneficiary := sdk.ZeroInt()

	for _, penalizedDelegator := range p.PenalizedDelegators {
		totalPenalty = totalPenalty.Add(penalizedDelegator.Amount)
		penalizedDelegators = append(penalizedDelegators, &sgn.AccountAmtPair{
			Account: mainchain.Hex2Addr(penalizedDelegator.Account).Bytes(),
			Amt:     penalizedDelegator.Amount.BigInt().Bytes(),
		})
	}

	for _, beneficiary := range p.Beneficiaries {
		amt := beneficiary.Fraction.MulInt(totalPenalty).TruncateInt()
		totalBeneficiary = totalBeneficiary.Add(amt)
		beneficiaries = append(beneficiaries, &sgn.AccountAmtPair{
			Account: mainchain.Hex2Addr(beneficiary.Account).Bytes(),
			Amt:     amt.BigInt().Bytes(),
		})
	}

	restPenalty := totalPenalty.Sub(totalBeneficiary)
	if restPenalty.IsPositive() {
		beneficiaries = append(beneficiaries, &sgn.AccountAmtPair{
			Account: mainchain.Hex2Addr(emptyAddr).Bytes(),
			Amt:     restPenalty.BigInt().Bytes(),
		})
	}

	penaltyBytes, _ := protobuf.Marshal(&sgn.Penalty{
		Nonce:               p.Nonce,
		ExpireTime:          expireTime,
		ValidatorAddress:    mainchain.Hex2Addr(p.ValidatorAddr).Bytes(),
		PenalizedDelegators: penalizedDelegators,
		Beneficiaries:       beneficiaries,
	})

	p.PenaltyProtoBytes = penaltyBytes
}

// Add signature to penalty sigs
func (p *Penalty) AddSig(sig []byte, expectedSigner string) error {
	sigs, err := common.AddSig(p.Sigs, p.PenaltyProtoBytes, sig, expectedSigner)
	if err != nil {
		log.Error(err)
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
