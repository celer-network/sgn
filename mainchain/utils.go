package mainchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrPeersNotMatch = errors.New("channel peers not match")
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

// GetSimplexSeqNum get the mainchain simplex seqNum
func GetSimplexSeqNum(
	ledger *CelerLedger, cid CidType,
	simplexSender, simplexReceiver Addr) (seqNum uint64, err error) {
	addrs, seqNums, err := ledger.GetStateSeqNumMap(&bind.CallOpts{}, cid)
	if err != nil {
		return 0, fmt.Errorf("GetStateSeqNumMap err: %w", err)
	}
	seqIndex := 0
	var match bool
	if simplexSender == addrs[0] {
		match = (simplexReceiver == addrs[1])
	} else if simplexSender == addrs[1] {
		match = (simplexReceiver == addrs[0])
		seqIndex = 1
	}
	if !match {
		return 0, fmt.Errorf("%w: cid %x, sender %x, receiver %x",
			ErrPeersNotMatch, cid, simplexSender, simplexReceiver)
	}

	return seqNums[seqIndex].Uint64(), nil
}
