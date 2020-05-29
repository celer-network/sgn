package common

import (
	"fmt"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ParamChange struct {
	Record   sdk.Int `json:"record"`
	NewValue sdk.Int `json:"newValue"`
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

func AddSig(sigs []Sig, msg []byte, sig []byte, expectedSigner string) (newSigs []Sig, err error) {
	signer, err := mainchain.RecoverSigner(msg, sig)
	if err != nil {
		return
	}

	signerAddr := mainchain.Addr2Hex(signer)
	if signerAddr != mainchain.FormatAddrHex(expectedSigner) {
		err = fmt.Errorf("invalid signer address %s %s", signerAddr, expectedSigner)
		return
	}

	for _, s := range sigs {
		if s.Signer == signerAddr {
			err = fmt.Errorf("repeated signer %s", signerAddr)
			return
		}
	}

	newSigs = append(sigs, NewSig(signerAddr, sig))
	return
}
