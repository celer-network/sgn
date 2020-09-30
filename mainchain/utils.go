package mainchain

import (
	"context"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func IsBonded(dposCandidateInfo DPoSCandidateInfo) bool {
	return dposCandidateInfo.Status.Uint64() == Bonded
}

func ParseStatus(dposCandidateInfo DPoSCandidateInfo) sdk.BondStatus {
	switch dposCandidateInfo.Status.Uint64() {
	case Bonded:
		return sdk.Bonded
	case Unbonding:
		return sdk.Unbonding
	case Unbonded:
		return sdk.Unbonded
	}

	return sdk.Unbonded
}

// GetEventSignature accepts the string of an event signature and return the hex
func GetEventSignature(eventSigStr string) HashType {
	return crypto.Keccak256Hash([]byte(eventSigStr))
}

// GetTxSender returns the sender address of the given transaction
func GetTxSender(ec *ethclient.Client, txHashStr string) (string, error) {
	tx, _, err := ec.TransactionByHash(context.Background(), Hex2Hash(txHashStr))
	if err != nil {
		return "", fmt.Errorf("Failed to get tx: %w", err)
	}
	msg, err := tx.AsMessage(ethtypes.NewEIP155Signer(tx.ChainId()))
	if err != nil {
		return "", fmt.Errorf("Failed to get msg: %w", err)
	}
	return Addr2Hex(msg.From()), nil
}

func GetAddressFromKeystore(ksBytes []byte) (string, error) {
	type ksStruct struct {
		Address string
	}
	var ks ksStruct
	if err := json.Unmarshal(ksBytes, &ks); err != nil {
		return "", err
	}
	return ks.Address, nil
}
