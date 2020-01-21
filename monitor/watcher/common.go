package watcher

import (
	"context"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// LogEventID tracks the position of a watch event in the event log.
type LogEventID struct {
	BlockNumber uint64 // Number of the block containing the event
	Index       int64  // Index of the event within the block
}

type Contract interface {
	GetAddr() ethcommon.Address
	GetABI() string
	GetETHClient() *ethclient.Client
	SendTransaction(*bind.TransactOpts, string, ...interface{}) (*ethtypes.Transaction, error)
	CallFunc(interface{}, string, ...interface{}) error
	WatchEvent(string, *bind.WatchOpts, <-chan bool) (ethtypes.Log, error)
	FilterEvent(string, *bind.FilterOpts, interface{}) (*EventIterator, error)
	ParseEvent(string, ethtypes.Log, interface{}) error
}

// BoundContract is a binding object for Ethereum smart contract
// It contains *bind.BoundContract (in go-ethereum) as an embedding
type BoundContract struct {
	*bind.BoundContract
	addr ethcommon.Address
	abi  string
	conn *ethclient.Client
}

// NewBoundContract creates a new contract binding
func NewBoundContract(
	conn *ethclient.Client,
	addr ethcommon.Address,
	rawABI string) (*BoundContract, error) {
	parsedABI, err := abi.JSON(strings.NewReader((rawABI)))
	return &BoundContract{
		bind.NewBoundContract(addr, parsedABI, conn, conn, conn),
		addr,
		rawABI,
		conn,
	}, err
}

// GetAddr returns contract addr
func (c *BoundContract) GetAddr() ethcommon.Address {
	return c.addr
}

// GetABI returns contract abi
func (c *BoundContract) GetABI() string {
	return c.abi
}

// GetETHClient return ethereum client
func (c *BoundContract) GetETHClient() *ethclient.Client {
	return c.conn
}

// SendTransaction sends transactions to smart contract via bound contract
func (c *BoundContract) SendTransaction(
	auth *bind.TransactOpts,
	method string,
	params ...interface{}) (*ethtypes.Transaction, error) {
	return c.Transact(auth, method, params...)
}

// CallFunc invokes a view-only contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (c *BoundContract) CallFunc(
	result interface{},
	method string,
	params ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	return c.Call(&bind.CallOpts{Context: ctx}, result, method, params...)
}

// WatchEvent subscribes to future events
// This function blocks until an event is catched or done signal is received
func (c *BoundContract) WatchEvent(
	name string,
	opts *bind.WatchOpts,
	done <-chan bool) (ethtypes.Log, error) {

	logs, sub, err := c.WatchLogs(opts, name)
	if err != nil {
		return ethtypes.Log{}, err
	}
	defer sub.Unsubscribe()
	select {
	case log := <-logs:
		return log, nil
	case <-done:
	}
	return ethtypes.Log{}, nil
}

// FilterEvent gets historical events
// This function returns an iterator over historical events
func (c *BoundContract) FilterEvent(
	name string,
	opts *bind.FilterOpts,
	event interface{}) (*EventIterator, error) {
	logs, sub, err := c.FilterLogs(opts, name)
	if err != nil {
		return nil, err
	}
	return &EventIterator{
		Contract: c,
		Event:    event,
		Name:     name,
		Logs:     logs,
		Sub:      sub,
	}, nil
}

// ParseEvent parses the catched event according to the event template
func (c *BoundContract) ParseEvent(
	name string,
	log ethtypes.Log,
	event interface{}) error {
	err := c.UnpackLog(event, name, log)
	return err
}

// EventIterator is returned from FilterEvent and is used to iterate over the raw logs and unpacked data
type EventIterator struct {
	Event    interface{}           // Event containing the contract specifics and raw log
	Contract *BoundContract        // Generic contract to use for unpacking event data
	Name     string                // Event name to use for unpacking event data
	Logs     chan ethtypes.Log     // Log channel receiving the found contract events
	Sub      ethereum.Subscription // Subscription for errors, completion and termination
	Done     bool                  // Whether the subscription completed delivering logs
	Fail     error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EventIterator) Next() (ethtypes.Log, bool) {
	// If the iterator failed, stop iterating
	if it.Fail != nil {
		return ethtypes.Log{}, false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.Done {
		select {
		case log := <-it.Logs:
			return log, true

		default:
			return ethtypes.Log{}, false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.Logs:
		return log, true

	case err := <-it.Sub.Err():
		it.Done = true
		it.Fail = err
		return it.Next()
	}
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventIterator) Close() error {
	it.Sub.Unsubscribe()
	return nil
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EventIterator) Error() error {
	return it.Fail
}
