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

type Penalty struct {
	Nonce               uint64           `json:"nonce"`
	ValidatorAddr       string           `json:"validatorAddr"`
	PenalizedDelegators []AccountAmtPair `json:"penalizedDelegators"`
	PenaltyProtoBytes   []byte           `json:"penaltyProtoBytes"`
	Sigs                []common.Sig     `json:"sigs"`
}

func NewPenalty(nonce uint64, validatorAddr string) Penalty {
	return Penalty{
		Nonce:         nonce,
		ValidatorAddr: validatorAddr,
	}
}

// implement fmt.Stringer
func (p Penalty) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Nonce: %d`, p.Nonce))
}

func (p Penalty) GenerateProtoBytes() {
	var penalizedDelegators []*sgn.AccountAmtPair
	for _, penalizedDelegator := range p.PenalizedDelegators {
		pd := &sgn.AccountAmtPair{
			Account: ethcommon.HexToAddress(penalizedDelegator.Account).Bytes(),
			Amt:     penalizedDelegator.Amount.BigInt().Bytes(),
		}
		penalizedDelegators = append(penalizedDelegators, pd)
	}

	penaltyBytes, _ := protobuf.Marshal(&sgn.Penalty{
		Nonce:               p.Nonce,
		ValidatorAddress:    ethcommon.HexToAddress(p.ValidatorAddr).Bytes(),
		PenalizedDelegators: penalizedDelegators,
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
