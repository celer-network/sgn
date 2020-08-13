package monitor

import (
	"math/big"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/x/slash"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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
	var keys, vals [][]byte
	m.lock.RLock()
	iterator, err := m.db.Iterator(PenaltyKeyPrefix, storetypes.PrefixEndBytes(PenaltyKeyPrefix))
	if err != nil {
		log.Errorln("Create db iterator err", err)
		return
	}
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, iterator.Key())
		vals = append(vals, iterator.Value())
	}
	iterator.Close()
	m.lock.RUnlock()

	for i, key := range keys {
		event := NewPenaltyEventFromBytes(vals[i])
		err = m.dbDelete(key)
		if err != nil {
			log.Errorln("db Delete err", err)
			continue
		}
		m.submitPenalty(event)
	}
}

func (m *Monitor) submitPenalty(penaltyEvent PenaltyEvent) {
	log.Infoln("Process Penalty", penaltyEvent.Nonce)

	used, err := m.EthClient.DPoS.UsedPenaltyNonce(&bind.CallOpts{}, big.NewInt(int64(penaltyEvent.Nonce)))
	if err != nil {
		log.Errorln("Get usedPenaltyNonce err", err)
		return
	}

	if used {
		log.Infof("Penalty %d has been used", penaltyEvent.Nonce)
		return
	}

	penaltyRequest, err := slash.CLIQueryPenaltyRequest(m.Transactor.CliCtx, slash.StoreKey, penaltyEvent.Nonce)
	if err != nil {
		log.Errorln("QueryPenaltyRequest err", err)
		return
	}

	tx, err := m.EthClient.Transactor.Transact(
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
			return m.EthClient.DPoS.Punish(opts, penaltyRequest)
		},
	)
	if err != nil {
		if penaltyEvent.RetryCount < maxPunishRetry {
			penaltyEvent.RetryCount = penaltyEvent.RetryCount + 1
			err = m.dbSet(GetPenaltyKey(penaltyEvent.Nonce), penaltyEvent.MustMarshal())
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
