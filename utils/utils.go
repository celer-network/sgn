// This file is based on github.com/goCeler/utils (commit ID: d7335ae321b67150d92de18f6589f1d1fd8b0910)

// Copyright 2018 Celer Network

package utils

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/celer-network/sgn/testing/log"
	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TODO(mzhou): Remove this once cEnv is cleaned up
func WaitMined(ctx context.Context, ec *ethclient.Client,
	tx *ethtypes.Transaction, blockDelay uint64) (*ethtypes.Receipt, error) {
	return WaitMinedWithTxHash(ctx, ec, tx.Hash().Hex(), blockDelay)
}

// WaitMined waits for tx to be mined on the blockchain
// It returns tx receipt when the tx has been mined and enough block confirmations have passed
func WaitMinedWithTxHash(ctx context.Context, ec *ethclient.Client,
	txHash string, blockDelay uint64) (*ethtypes.Receipt, error) {
	// an error possibly returned when a transaction is pending
	const missingFieldErr = "missing required field 'transactionHash' for Log"

	if ec == nil {
		return nil, errors.New("nil ethclient")
	}
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	// wait tx to be mined
	txHashBytes := ethcommon.HexToHash(txHash)
	for {
		receipt, rerr := ec.TransactionReceipt(ctx, txHashBytes)
		if rerr == nil {
			log.Debugf("Transaction mined. Waiting for %d block confirmations", blockDelay)
			if blockDelay == 0 {
				return receipt, rerr
			}
			break
		} else if rerr == ethereum.NotFound || rerr.Error() == missingFieldErr {
			// Wait for the next round
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-queryTicker.C:
			}
		} else {
			return receipt, rerr
		}
	}
	// wait for enough block confirmations
	ddl := big.NewInt(0)
	latestBlockHeader, err := ec.HeaderByNumber(ctx, nil)
	if err == nil {
		ddl.Add(new(big.Int).SetUint64(blockDelay), latestBlockHeader.Number)
	}
	for {
		latestBlockHeader, err := ec.HeaderByNumber(ctx, nil)
		if err == nil && ddl.Cmp(latestBlockHeader.Number) < 0 {
			receipt, rerr := ec.TransactionReceipt(ctx, txHashBytes)
			if rerr == nil {
				log.Debugln("tx confirmed!")
				return receipt, rerr
			} else if rerr == ethereum.NotFound || rerr.Error() == missingFieldErr {
				return nil, errors.New("tx is dropped due to chain re-org")
			} else {
				return receipt, rerr
			}
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
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
