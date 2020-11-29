package ethtx

import (
	"math/big"

	esTxManager "github.com/celer-network/eth-services/txmanager"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"
)

type ContractInfo struct {
	Address gethCommon.Address
	ABI     abi.ABI
}

func (s *TxSender) SnapshotStates(
	signedSimplexStateArrayBytes []byte,
	handler esTxManager.JobHandler,
) error {
	data, err := s.ledgerInfo.ABI.Pack("snapshotStates", signedSimplexStateArrayBytes)
	if err != nil {
		return err
	}
	return s.sendTx(s.from, s.ledgerInfo, big.NewInt(0), data, handler)
}

func (s *TxSender) IntendSettle(
	signedSimplexStateArrayBytes []byte,
	handler esTxManager.JobHandler,
) error {
	data, err := s.ledgerInfo.ABI.Pack("intendSettle", signedSimplexStateArrayBytes)
	if err != nil {
		return err
	}
	return s.sendTx(s.from, s.ledgerInfo, big.NewInt(0), data, handler)
}
