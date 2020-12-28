package types

import (
	"fmt"
	"strings"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/sgn"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/protobuf/proto"
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
	Fraction sdk.Dec `json:"fraction"`
}

func NewAccountFractionPair(account string, fraction sdk.Dec) AccountFractionPair {
	return AccountFractionPair{
		Account:  mainchain.FormatAddrHex(account),
		Fraction: fraction,
	}
}

func (amp AccountFractionPair) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Account: %s, Fraction: %v`, amp.Account, amp.Fraction))
}

type Penalty struct {
	Nonce               uint64           `json:"nonce"`
	ExpireTime          uint64           `json:"expire_time"`
	Reason              string           `json:"reason"`
	ValidatorAddr       string           `json:"validator_addr"`
	TotalPenalty        sdk.Int          `json:"totalPenalty"`
	PenalizedDelegators []AccountAmtPair `json:"penalized_delegators"`
	Beneficiaries       []AccountAmtPair `json:"beneficiaries"`
	PenaltyProtoBytes   []byte           `json:"penalty_proto_bytes"`
	Sigs                []common.Sig     `json:"sigs"`
}

func NewPenalty(nonce uint64, reason string, validatorAddr string, penalizedDelegators []AccountAmtPair,
	beneficiaryFractions []AccountFractionPair, syncerReward sdk.Int, expireTime uint64) Penalty {
	var beneficiaries []AccountAmtPair
	totalPenalty := sdk.ZeroInt()
	totalBeneficiary := sdk.ZeroInt()

	for _, penalizedDelegator := range penalizedDelegators {
		totalPenalty = totalPenalty.Add(penalizedDelegator.Amount)
	}

	if syncerReward.GTE(totalPenalty) {
		beneficiaries = append(beneficiaries, NewAccountAmtPair("1", totalPenalty))
	} else {
		beneficiaries = append(beneficiaries, NewAccountAmtPair("1", syncerReward))

		restPenalty := totalPenalty.Sub(syncerReward)

		for _, beneficiaryFraction := range beneficiaryFractions {
			amt := beneficiaryFraction.Fraction.MulInt(restPenalty).TruncateInt()
			totalBeneficiary = totalBeneficiary.Add(amt)
			beneficiaries = append(beneficiaries, NewAccountAmtPair(beneficiaryFraction.Account, amt))
		}

		restPenalty = restPenalty.Sub(totalBeneficiary)
		if restPenalty.IsPositive() {
			beneficiaries = append(beneficiaries, NewAccountAmtPair(mainchain.ZeroAddr.String(), restPenalty))
		}
	}

	return Penalty{
		Nonce:               nonce,
		ExpireTime:          expireTime,
		Reason:              reason,
		ValidatorAddr:       mainchain.FormatAddrHex(validatorAddr),
		TotalPenalty:        totalPenalty,
		Beneficiaries:       beneficiaries,
		PenalizedDelegators: penalizedDelegators,
	}
}

// implement fmt.Stringer
func (p Penalty) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Nonce: %d, Reason: %s, ValidatorAddr: %s, TotalPenalty: %s`,
		p.Nonce, p.Reason, p.ValidatorAddr, p.TotalPenalty))
}

func (p *Penalty) GenerateProtoBytes() {
	var penalizedDelegators []*sgn.AccountAmtPair
	var beneficiaries []*sgn.AccountAmtPair

	for _, penalizedDelegator := range p.PenalizedDelegators {
		penalizedDelegators = append(penalizedDelegators, &sgn.AccountAmtPair{
			Account: mainchain.Hex2Addr(penalizedDelegator.Account).Bytes(),
			Amt:     penalizedDelegator.Amount.BigInt().Bytes(),
		})
	}

	for _, beneficiary := range p.Beneficiaries {
		beneficiaries = append(beneficiaries, &sgn.AccountAmtPair{
			Account: mainchain.Hex2Addr(beneficiary.Account).Bytes(),
			Amt:     beneficiary.Amount.BigInt().Bytes(),
		})
	}

	penaltyBytes, _ := proto.Marshal(&sgn.Penalty{
		Nonce:               p.Nonce,
		ExpireTime:          p.ExpireTime,
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

	penaltyRequestBytes, _ := proto.Marshal(&sgn.PenaltyRequest{
		Penalty: p.PenaltyProtoBytes,
		Sigs:    sigs,
	})

	return penaltyRequestBytes
}
