package mainchain

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

func (ethClient *EthClient) SignMessage(data []byte) ([]byte, error) {
	return SignMessage(ethClient.PrivateKey, data)
}

func SignMessage(privateKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	return crypto.Sign(generatePrefixedHash(data), privateKey)
}

func RecoverSigner(data []byte, sig []byte) (Addr, error) {
	pubKey, err := crypto.SigToPub(generatePrefixedHash(data), sig)
	if err != nil {
		return ZeroAddr, err
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr, nil
}

func generatePrefixedHash(data []byte) []byte {
	return crypto.Keccak256([]byte("\x19Ethereum Signed Message:\n32"), crypto.Keccak256(data))
}
