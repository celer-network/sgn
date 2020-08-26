package common

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type ParamChange struct {
	Record   sdk.Int `json:"record"`
	NewValue sdk.Int `json:"new_value"`
}

func NewParamChange(record, newValue sdk.Int) ParamChange {
	return ParamChange{
		Record:   record,
		NewValue: newValue,
	}
}

// implement fmt.Stringer
func (p ParamChange) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Record: %v, NewValue: %v`, p.Record, p.NewValue))
}

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
	return strings.TrimSpace(fmt.Sprintf(`Signer: %s, Sig: %x,`, r.Signer, r.Sig))
}

func AddSig(sigs []Sig, msg []byte, sig []byte, expectedSigner string) ([]Sig, error) {
	signer, err := eth.RecoverSigner(msg, sig)
	if err != nil {
		return nil, err
	}

	signerAddr := mainchain.Addr2Hex(signer)
	if signerAddr != mainchain.FormatAddrHex(expectedSigner) {
		err = fmt.Errorf("invalid signer address %s %s", signerAddr, expectedSigner)
		return nil, err
	}

	for _, s := range sigs {
		if s.Signer == signerAddr {
			err = fmt.Errorf("repeated signer %s", signerAddr)
			return nil, err
		}
	}

	return append(sigs, NewSig(signerAddr, sig)), nil
}

func NewCommission(ethClient *mainchain.EthClient, commissionRate *big.Int) (staking.Commission, error) {
	commissionBase, err := ethClient.DPoS.COMMISSIONRATEBASE(&bind.CallOpts{})
	if err != nil {
		return staking.Commission{}, err
	}

	prec := int64(len(commissionBase.String()) - 1)
	return staking.Commission{
		CommissionRates: staking.CommissionRates{
			Rate:          sdk.NewDecFromBigIntWithPrec(commissionRate, prec),
			MaxRate:       sdk.NewDec(1),
			MaxChangeRate: sdk.NewDec(1),
		},
	}, nil
}
