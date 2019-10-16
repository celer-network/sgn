// Copyright 2018 Celer Network

package utils

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/celer-network/sgn/goceler-copy/common"
	"github.com/celer-network/sgn/goceler-copy/ctype"
	"github.com/celer-network/sgn/proto/entity"
	"github.com/celer-network/sgn/testing/log"
	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/protobuf/jsonpb"
	proto "github.com/golang/protobuf/proto"
	"github.com/shopspring/decimal"
)

// Dec2HexStr decimal string to hex
func Dec2HexStr(dec string) string {
	i := new(big.Int)
	i.SetString(dec, 10)
	return i.Text(16)
}

// Hex2DecStr hex string to decimal
func Hex2DecStr(hex string) string {
	i := new(big.Int)
	i.SetString(hex, 16)
	return i.Text(10)
}

func BytesToBigInt(in []byte) *big.Int {
	ret := big.NewInt(0)
	ret.SetBytes(in)
	return ret
}

// convert decimal wei string to big.Int
func Wei2BigInt(wei string) *big.Int {
	i := big.NewInt(0)
	_, ok := i.SetString(wei, 10)
	if !ok {
		return nil
	}
	return i
}

// float in 10e18 wei to wei
func Float2Wei(f float64) *big.Int {
	if f < 0 {
		return nil
	}
	wei := decimal.NewFromFloat(f).Mul(decimal.NewFromFloat(10).Pow(decimal.NewFromFloat(18)))
	weiInt := new(big.Int)
	weiInt.SetString(wei.String(), 10)
	return weiInt
}

// left padding
func Pad(origin []byte, n int) []byte {
	m := len(origin)
	padded := make([]byte, n)
	pn := n - m
	for i := m - 1; i >= 0; i-- {
		padded[pn+i] = origin[i]
	}
	return padded
}

func TryLock(m *sync.Mutex) bool {
	const mutexLocked = 1 << iota
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(m)), 0, mutexLocked)
}

func ValidateAndFormatAddress(address string) (string, error) {
	if !ethcommon.IsHexAddress(address) {
		return "", errors.New("Invalid address")
	}
	return ctype.Bytes2Hex(ctype.Hex2Bytes(address)), nil
}

// GetTokenAddrStr returns string for tokenInfo
func GetTokenAddrStr(tokenInfo *entity.TokenInfo) string {
	switch tktype := tokenInfo.TokenType; tktype {
	case entity.TokenType_ETH:
		return common.EthContractAddr
	case entity.TokenType_ERC20:
		return ctype.Bytes2Hex(tokenInfo.TokenAddress)
	}
	return ""
}

// GetTokenInfoFromAddress returns TokenInfo from tkaddr
// only support ERC20 for now
func GetTokenInfoFromAddress(tkaddr ctype.Addr) *entity.TokenInfo {
	tkInfo := new(entity.TokenInfo)
	if tkaddr == ctype.ZeroAddr {
		tkInfo.TokenType = entity.TokenType_ETH
	} else {
		tkInfo.TokenType = entity.TokenType_ERC20
		tkInfo.TokenAddress = tkaddr.Bytes()
	}
	return tkInfo
}

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

// Serialize a protobuf to json string
func PbToJSONString(pb proto.Message) string {
	m := jsonpb.Marshaler{}
	ret, err := m.MarshalToString(pb)
	if err != nil {
		return ""
	}
	return ret
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
