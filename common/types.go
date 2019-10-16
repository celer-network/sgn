package common

import (
	"errors"
	"fmt"
	"strings"

	"github.com/celer-network/sgn/mainchain"
)

// CProfile struct is based on github.com/goCeler/common (commit ID: d7335ae321b67150d92de18f6589f1d1fd8b0910)
// CProfile contains configurations for CelerClient/OSP
type CProfile struct {
	ETHInstance        string `json:"ethInstance"`
	SvrETHAddr         string `json:"svrEthAddr"`
	WalletAddr         string `json:"walletAddr"`
	LedgerAddr         string `json:"ledgerAddr"`
	VirtResolverAddr   string `json:"virtResolverAddr"`
	EthPoolAddr        string `json:"ethPoolAddr"`
	PayResolverAddr    string `json:"payResolverAddr"`
	PayRegistryAddr    string `json:"payRegistryAddr"`
	RouterRegistryAddr string `json:"routerRegistryAddr"`
	SvrRPC             string `json:"svrRpc"`
	SelfRPC            string `json:"selfRpc,omitempty"`
	StoreDir           string `json:"storeDir,omitempty"`
	StoreSql           string `json:"storeSql,omitempty"`
	WebPort            string `json:"webPort,omitempty"`
	WsOrigin           string `json:"wsOrigin,omitempty"`
	ChainId            int64  `json:"chainId"`
	BlockDelayNum      uint64 `json:"blockDelayNum"`
	IsOSP              bool   `json:"isOsp,omitempty"`
	ListenOnChain      bool   `json:"listenOnChain,omitempty"`
	PollingInterval    uint64 `json:"pollingInterval"`
	DisputeTimeout     uint64 `json:"disputeTimeout"`
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
