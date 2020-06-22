package monitor

import (
	"math/big"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/x/slash"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	maxPunishRetry = 5
)

func (m *Monitor) processPenaltyQueue() {
	if !m.isPusher() {
		return
	}

	iterator, err := m.db.Iterator(PenaltyKeyPrefix, storetypes.PrefixEndBytes(PenaltyKeyPrefix))
	if err != nil {
		log.Errorln("Create db iterator err", err)
		return
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewPenaltyEventFromBytes(iterator.Value())
		err = m.db.Delete(iterator.Key())
		if err != nil {
			log.Errorln("db Delete err", err)
			continue
		}
		m.submitPenalty(event)
	}
}

func (m *Monitor) submitPenalty(penaltyEvent PenaltyEvent) {
	log.Infoln("Process Penalty", penaltyEvent.Nonce)

	used, err := m.ethClient.DPoS.UsedPenaltyNonce(
		&bind.CallOpts{BlockNumber: sdk.NewIntFromUint64(m.secureBlkNum).BigInt()},
		big.NewInt(int64(penaltyEvent.Nonce)))
	if err != nil {
		log.Errorln("Get usedPenaltyNonce err", err)
		return
	}

	if used {
		log.Infof("Penalty %d has been used", penaltyEvent.Nonce)
		return
	}

	penaltyRequest, err := slash.CLIQueryPenaltyRequest(m.operator.CliCtx, slash.StoreKey, penaltyEvent.Nonce)
	if err != nil {
		log.Errorln("QueryPenaltyRequest err", err)
		return
	}

	tx, err := m.ethClient.Transactor.Transact(
		&eth.TransactionStateHandler{
			OnMined: func(receipt *ethtypes.Receipt) {
				if receipt.Status == ethtypes.ReceiptStatusSuccessful {
					log.Infof("Punish transaction %x succeeded", receipt.TxHash)
				} else {
					log.Errorf("Punish transaction %x failed", receipt.TxHash)
				}
			},
			OnError: func(tx *ethtypes.Transaction, err error) {
				log.Errorf("Punish transaction %x err: %s", tx.Hash(), err)
			},
		},
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return m.ethClient.DPoS.Punish(opts, penaltyRequest)
		},
	)
	if err != nil {
		if penaltyEvent.RetryCount < maxPunishRetry {
			penaltyEvent.RetryCount = penaltyEvent.RetryCount + 1
			err = m.db.Set(GetPenaltyKey(penaltyEvent.Nonce), penaltyEvent.MustMarshal())
			if err != nil {
				log.Errorln("db Set err", err)
			}
			return
		}
		log.Errorln("Punish err", err)
		return
	}
	log.Infoln("Punish tx submitted", tx.Hash().Hex())
}
