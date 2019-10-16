// This file is based on github.com/goCeler/ctype (commit ID: d7335ae321b67150d92de18f6589f1d1fd8b0910)

// Copyright 2018-2019 Celer Network

// util package to handle various types and hex string, bytes etc
package ctype

/*
Terms in this package:
Hex: hex string. Hex2xxx accepts with or without 0x prefix. xxxToHex always without 0x
Bytes: []byte, mostly for interacting with protobuf
Addr: go-ethereum/common.Address [20]byte
*/

import (
	"encoding/hex"
	"math/big"

	ec "github.com/ethereum/go-ethereum/common"
)

var (
	// ZeroAddr is all 0s
	ZeroAddr Addr
	// ZeroAddrHex is string of 20 0s
	ZeroAddrHex = Addr2Hex(ZeroAddr)
	// ZeroBigInt is big.NewInt(0)
	ZeroBigInt = big.NewInt(0)
	// ZeroCid is all 0s
	ZeroCid CidType
)

// PayIDType is the ID type for pays
type PayIDType = ec.Hash

// CidType is the type for payment channel ID
// Note we need to change all cid.Hex() to Cid2Hex() because Hash.Hex() has 0x prefix
type CidType = ec.Hash

// Addr is alias to geth common.Address
type Addr = ec.Address

// ========== Hex/Bytes ==========

// Hex2Bytes supports hex string with or without 0x prefix
// Calls hex.DecodeString directly and ignore err
// similar to ec.FromHex but better
func Hex2Bytes(s string) (b []byte) {
	if len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		s = s[2:]
	}
	// hex.DecodeString expects an even-length string
	if len(s)%2 == 1 {
		s = "0" + s
	}
	b, _ = hex.DecodeString(s)
	return b
}

// Bytes2Hex returns hex string without 0x prefix
func Bytes2Hex(b []byte) string {
	return hex.EncodeToString(b)
}

// ========== Address ==========

// Hex2Addr accepts hex string with or without 0x prefix and return Addr
func Hex2Addr(s string) Addr {
	return ec.BytesToAddress(Hex2Bytes(s))
}

// Addr2Hex returns hex without 0x
func Addr2Hex(a Addr) string {
	return Bytes2Hex(a[:])
}

// Bytes2Addr returns Address from b
// Addr.Bytes() does the reverse
func Bytes2Addr(b []byte) Addr {
	return ec.BytesToAddress(b)
}

// ========== CidType ==========

// Bytes2Cid converts bytes to CidType
func Bytes2Cid(b []byte) CidType {
	return ec.BytesToHash(b)
}

// Cid2Hex returns hex without 0x prefix
func Cid2Hex(p CidType) string {
	return Bytes2Hex(p[:])
}

// Hex2Cid accepts hex string with or without 0x prefix and return CidType
func Hex2Cid(s string) CidType {
	return ec.BytesToHash(Hex2Bytes(s))
}
