package mainchain

import (
	"github.com/celer-network/sgn-contract/bindings/go/sgncontracts"
	"github.com/ethereum/go-ethereum/ethclient"
)

type LedgerContract struct {
	*CelerLedger
	Address Addr
}

func NewLedgerContract(address Addr, client *ethclient.Client) (*LedgerContract, error) {
	ledger, err := NewCelerLedger(address, client)
	if err != nil {
		return nil, err
	}
	return &LedgerContract{
		CelerLedger: ledger,
		Address:     address,
	}, nil
}

func (c *LedgerContract) GetAddr() Addr {
	return c.Address
}

func (c *LedgerContract) GetABI() string {
	return CelerLedgerABI
}

type DposContract struct {
	*sgncontracts.DPoS
	Address Addr
}

func NewDposContract(address Addr, client *ethclient.Client) (*DposContract, error) {
	dpos, err := sgncontracts.NewDPoS(address, client)
	if err != nil {
		return nil, err
	}
	return &DposContract{
		DPoS:    dpos,
		Address: address,
	}, nil
}

func (c *DposContract) GetAddr() Addr {
	return c.Address
}

func (c *DposContract) GetABI() string {
	return sgncontracts.DPoSABI
}

type SgnContract struct {
	*sgncontracts.SGN
	Address Addr
}

func NewSgnContract(address Addr, client *ethclient.Client) (*SgnContract, error) {
	sgn, err := sgncontracts.NewSGN(address, client)
	if err != nil {
		return nil, err
	}
	return &SgnContract{
		SGN:     sgn,
		Address: address,
	}, nil
}

func (c *SgnContract) GetAddr() Addr {
	return c.Address
}

func (c *SgnContract) GetABI() string {
	return sgncontracts.SGNABI
}
