package monitor

import (
	"math/big"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/validator"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	maxSlashRetry = 5
)

func (m *Monitor) processPenaltyQueue() {
	if !m.isSyncer() {
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

	if !m.validatePenaltySigs(penaltyEvent.Nonce) {
		log.Infof("Penalty %d does not have enough sigs", penaltyEvent.Nonce)
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
					log.Infof("Slash transaction %x succeeded", receipt.TxHash)
				} else {
					log.Errorf("Slash transaction %x failed", receipt.TxHash)
				}
			},
			OnError: func(tx *ethtypes.Transaction, err error) {
				log.Errorf("Slash transaction %x err: %s", tx.Hash(), err)
			},
		},
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return m.EthClient.DPoS.Slash(opts, penaltyRequest)
		},
	)
	if err != nil {
		if penaltyEvent.RetryCount < maxSlashRetry {
			penaltyEvent.RetryCount = penaltyEvent.RetryCount + 1
			err = m.dbSet(GetPenaltyKey(penaltyEvent.Nonce), penaltyEvent.MustMarshal())
			if err != nil {
				log.Errorln("db Set err", err)
			}
			return
		}
		log.Errorln("Slash err", err)
		return
	}
	log.Infoln("Slash tx submitted", tx.Hash().Hex())
}

func (m *Monitor) validatePenaltySigs(nonce uint64) bool {
	penalty, err := slash.CLIQueryPenalty(m.Transactor.CliCtx, slash.StoreKey, nonce)
	if err != nil {
		log.Errorln("QueryPenalty err", err)
		return false
	}

	signedValidators := mapset.NewSet()
	for _, sig := range penalty.Sigs {
		signedValidators.Add(sig.Signer)
	}

	validators, err := validator.CLIQueryBondedValidators(m.Transactor.CliCtx, validator.StoreKey)
	if err != nil {
		log.Errorln("QueryBondedValidators err", err)
		return false
	}

	totalStake := sdk.ZeroInt()
	votingStake := sdk.ZeroInt()
	for _, v := range validators {
		totalStake = totalStake.Add(v.BondedTokens())

		if signedValidators.Contains(v.Description.Identity) {
			votingStake = votingStake.Add(v.BondedTokens())
		}
	}

	return votingStake.GTE(totalStake.MulRaw(2).QuoRaw(3))
}
