package mainchain

import (
	"crypto/ecdsa"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (ethClient *EthClient) SignMessage(data []byte) ([]byte, error) {
	return crypto.Sign(generatePrefixedHash(data), ethClient.PrivateKey)
}

func SignMessage(privateKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	sig, err := crypto.Sign(generatePrefixedHash(data), privateKey)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func RecoverSigner(data []byte, sig []byte) (ethcommon.Address, error) {
	pubKey, err := crypto.SigToPub(generatePrefixedHash(data), sig)
	if err != nil {
		return ethcommon.Address{}, err
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr, nil
}

func generatePrefixedHash(data []byte) []byte {
	return crypto.Keccak256([]byte("\x19Ethereum Signed Message:\n32"), crypto.Keccak256(data))
}
