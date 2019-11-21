package common

import (
	"errors"
	"fmt"
	"strings"

	"github.com/celer-network/sgn/mainchain"
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

func AddSig(sigs []Sig, msg []byte, sig []byte, expectedSigner string) (newSigs []Sig, err error) {
	signer, err := mainchain.RecoverSigner(msg, sig)
	if err != nil {
		return
	}

	signerAddr := signer.String()
	if signerAddr != expectedSigner {
		err = errors.New("invalid signer address")
		return
	}

	for _, s := range sigs {
		if s.Signer == signerAddr {
			err = errors.New("repeated signer")
			return
		}
	}

	newSigs = append(sigs, NewSig(signerAddr, sig))
	return
}
