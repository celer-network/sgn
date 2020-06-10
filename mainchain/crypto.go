package mainchain

import (
	"github.com/celer-network/goutils/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func RecoverSigner(data []byte, sig []byte) (common.Address, error) {
	if len(sig) == 65 { // we could return zeroAddr if len not 65
		if sig[64] == 27 || sig[64] == 28 {
			// SigToPub only expect v to be 0 or 1,
			// see https://github.com/ethereum/go-ethereum/blob/v1.8.23/internal/ethapi/api.go#L468.
			// we've been ok as our own code only has v 0 or 1, but using external signer may cause issue
			// we also fix v in celersdk.PublishSignedResult to be extra safe
			sig[64] -= 27
		}
	}
	pubKey, err := crypto.SigToPub(eth.GeneratePrefixedHash(data), sig)
	if err != nil {
		return common.Address{}, err
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr, nil
}
