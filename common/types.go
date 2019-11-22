package common

import (
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

	signerAddr := mainchain.Addr2Hex(signer)
	if signerAddr != mainchain.FormatAddrHex(expectedSigner) {
		err = fmt.Errorf("invalid signer address %s %s", signerAddr, expectedSigner)
		return
	}

	for _, s := range sigs {
		if mainchain.FormatAddrHex(s.Signer) == signerAddr {
			err = fmt.Errorf("repeated signer %s", signerAddr)
			return
		}
	}

	newSigs = append(sigs, NewSig(signerAddr, sig))
	return
}
