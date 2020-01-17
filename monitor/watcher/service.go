// Copyright 2018 Celer Network

package watcher

import (
	"container/heap"
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/celer-network/goutils/log"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// log watch polling as a multiplier of block number polling
	watchBlockInterval = uint64(1)
)

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

// CallbackID is the unique callback ID for deadlines and events
type CallbackID uint64

// Deadline is the metadata of a deadline
type Deadline struct {
	BlockNum *big.Int
	Callback func()
}

// DeadlineQueue is the priority queue for deadlines
type DeadlineQueue []*big.Int

func (dq DeadlineQueue) Len() int { return len(dq) }

func (dq DeadlineQueue) Less(i, j int) bool { return dq[i].Cmp(dq[j]) == -1 }

func (dq DeadlineQueue) Swap(i, j int) { dq[i], dq[j] = dq[j], dq[i] }

func (dq *DeadlineQueue) Push(x interface{}) { *dq = append(*dq, x.(*big.Int)) }

func (dq *DeadlineQueue) Pop() (popped interface{}) {
	popped = (*dq)[len(*dq)-1]
	*dq = (*dq)[:len(*dq)-1]
	return
}

func (dq *DeadlineQueue) Top() (top interface{}) {
	if len(*dq) > 0 {
		top = (*dq)[0]
	}
	return
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

type BlockNumber interface {
	GetCurrentBlockNumber() (*big.Int, error)
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

// Event is the metadata for an event
type Event struct {
	Addr       ethcommon.Address
	RawAbi     string
	Name       string
	WatchName  string
	StartBlock *big.Int
	EndBlock   *big.Int
	BlockDelay uint64
	Callback   func(CallbackID, ethtypes.Log)
	watch      *Watch
}

// Service struct stores service parameters and registered deadlines and events
type Service struct {
	watch         *WatchService           // persistent watch service
	deadlines     map[CallbackID]Deadline // deadlines
	deadlinecbs   map[string][]CallbackID // deadline callbacks
	deadlineQueue DeadlineQueue           // deadline priority queue
	events        map[CallbackID]Event    // events
	mu            sync.Mutex
	blockDelay    uint64
	// HACK HACK, the following fields should be removed
	enabled bool
	rpcAddr string
}

// NewService starts a new monitor service. Currently, if "enabled" is false,
// event monitoring will be disabled, and the IP address of the cNode given as
// "rpcAddr" will be printed.
func NewService(
	watch *WatchService, blockDelay uint64, enabled bool, rpcAddr string) *Service {
	s := &Service{
		watch:       watch,
		deadlines:   make(map[CallbackID]Deadline),
		deadlinecbs: make(map[string][]CallbackID),
		events:      make(map[CallbackID]Event),
		blockDelay:  blockDelay,
		enabled:     enabled,
		rpcAddr:     rpcAddr,
	}
	return s
}

// Init creates the event map
func (s *Service) Init() {
	heap.Init(&s.deadlineQueue)
	go s.monitorDeadlines() // start monitoring deadlines
}

// Close only set events map to empty map so all monitorEvent will exit due to isEventRemoved is true
func (s *Service) Close() {
	s.mu.Lock()
	s.events = make(map[CallbackID]Event)
	s.mu.Unlock()
}

func (s *Service) GetCurrentBlockNumber() *big.Int {
	return s.watch.GetCurrentBlockNumber()
}

// RegisterDeadline registers the deadline and returns the ID
func (s *Service) RegisterDeadline(d Deadline) CallbackID {
	// get a unique callback ID
	s.mu.Lock()
	defer s.mu.Unlock()
	var id CallbackID
	for {
		id = CallbackID(rand.Uint64())
		if _, exist := s.deadlines[id]; !exist {
			break
		}
	}

	// register deadline
	s.deadlines[id] = d
	_, ok := s.deadlinecbs[d.BlockNum.String()]
	if !ok {
		heap.Push(&s.deadlineQueue, d.BlockNum)
	}
	s.deadlinecbs[d.BlockNum.String()] = append(s.deadlinecbs[d.BlockNum.String()], id)
	return id
}

// continuously monitoring deadlines
func (s *Service) monitorDeadlines() {
	for {
		time.Sleep(2 * time.Second)
		s.mu.Lock()
		blockNumber := s.GetCurrentBlockNumber()
		for s.deadlineQueue.Len() > 0 && s.deadlineQueue.Top().(*big.Int).Cmp(blockNumber) < 1 {
			timeblock := heap.Pop(&s.deadlineQueue).(*big.Int)
			cbs, ok := s.deadlinecbs[timeblock.String()]
			if ok {
				dlCbs := make(map[CallbackID]Deadline)
				for _, id := range cbs {
					deadline, ok := s.deadlines[id]
					if ok {
						dlCbs[id] = deadline
						delete(s.deadlines, id)
					}
				}
				delete(s.deadlinecbs, timeblock.String())

				s.mu.Unlock()
				for _, deadline := range dlCbs {
					deadline.Callback()
				}
				s.mu.Lock()
			}
		}
		s.mu.Unlock()
	}
}

// Create a watch for the given event.  Use or skip using the StartBlock
// value from the event: the first time a watch is created for an event,
// the StartBlock should be used.  In follow-up re-creation of the watch
// after the previous watch was disconnected, skip using the StartBlock
// because the watch itself has persistence and knows the most up-to-date
// block to resume from instead of the original event StartBlock which is
// stale information by then.  If "reset" is enabled, the watcher ignores the
// previously stored position in the subscription which resets the stream to its
// start.
func (s *Service) createEventWatch(
	e Event, useStartBlock bool, reset bool) (*Watch, error) {
	var startBlock *big.Int
	if useStartBlock {
		startBlock = e.StartBlock
	}

	q, err := s.watch.MakeFilterQuery(e.Addr, e.RawAbi, e.Name, startBlock)
	if err != nil {
		return nil, err
	}
	return s.watch.NewWatch(e.WatchName, q, e.BlockDelay, watchBlockInterval, reset)
}

func (s *Service) Monitor(
	eventName string,
	contract Contract,
	startBlock *big.Int,
	endBlock *big.Int,
	reset bool,
	callback func(CallbackID, ethtypes.Log)) (CallbackID, error) {
	if !s.enabled {
		log.Infof("OSP (%s) not listening to on-chain logs", s.rpcAddr)
		return 0, nil
	}
	addr := contract.GetAddr()
	watchName := fmt.Sprintf("%s-%s", addr.String(), eventName)
	eventToListen := &Event{
		Addr:       addr,
		RawAbi:     contract.GetABI(),
		Name:       eventName,
		WatchName:  watchName,
		StartBlock: startBlock,
		EndBlock:   endBlock,
		BlockDelay: s.blockDelay,
		Callback:   callback,
	}

	id, err := s.MonitorEvent(*eventToListen, reset)
	if err != nil {
		log.Errorf("Cannot register event %s: %s", watchName, err)
		return 0, err
	}
	return id, nil
}

func (s *Service) MonitorEvent(e Event, reset bool) (CallbackID, error) {
	// Construct the watch now to return up-front errors to the caller.
	w, err := s.createEventWatch(e, true /* useStartBlock */, reset)
	if err != nil {
		log.Errorln("register event error:", err)
		return 0, err
	}

	// get a unique callback ID
	s.mu.Lock()
	var id CallbackID
	for {
		id = CallbackID(rand.Uint64())
		if _, exist := s.events[id]; !exist {
			break
		}
	}
	e.watch = w

	// register event
	s.events[id] = e
	s.mu.Unlock()

	go s.monitorEvent(e, id)

	return id, nil
}

func (s *Service) isEventRemoved(id CallbackID) bool {
	s.mu.Lock()
	_, ok := s.events[id]
	s.mu.Unlock()
	return !ok
}

// subscribes to events using a persistent
func (s *Service) monitorEvent(e Event, id CallbackID) {
	// WatchEvent blocks until an event is caught
	log.Debugln("monitoring event", e.Name)
	for {
		eventLog, err := e.watch.Next()
		if err != nil {
			log.Errorln("monitoring event error:", e.Name, err)
			e.watch.Close()

			var w *Watch
			for {
				if s.isEventRemoved(id) {
					e.watch.Close()
					return
				}
				w, err = s.createEventWatch(e, false /* useStartBlock */, false /* reset */)
				if err == nil {
					break
				}
				log.Errorln("recreate event watch error:", e.Name, err)
				time.Sleep(1 * time.Second)
			}

			s.mu.Lock()
			e.watch = w
			s.mu.Unlock()
			log.Debugln("event watch recreated", e.Name)
			continue
		}

		// When event log is removed due to chain re-org, just ignore it
		// TODO: emit error msg and properly roll back upon catching removed event log
		if eventLog.Removed {
			log.Warnf("Receive removed %s event log", e.Name)
			if err = e.watch.Ack(); err != nil {
				log.Errorln("monitoring event ACK error:", e.Name, err)
				e.watch.Close()
				return
			}
			continue
		}

		// Stop watching if the event was removed
		// TODO(mzhou): Also stop monitoring if timeout has passed
		if s.isEventRemoved(id) {
			e.watch.Close()
			return
		}

		e.Callback(id, eventLog)

		if err = e.watch.Ack(); err != nil {
			// This is a coding bug, just exit the loop.
			log.Errorln("monitoring event ACK error:", e.Name, err)
			e.watch.Close()
			return
		}
		if s.isEventRemoved(id) {
			e.watch.Close()
			return
		}
	}
}

// RemoveDeadline removes a deadline from the monitor
func (s *Service) RemoveDeadline(id CallbackID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Debugf("revoke deadline monitoring %d", id)
	delete(s.deadlines, id)
}

// RemoveEvent removes an event from the monitor
func (s *Service) RemoveEvent(id CallbackID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.events[id]
	if ok {
		log.Debugf("revoke event monitoring %d event %s", id, e.Name)
		e.watch.Close()
		delete(s.events, id)
	}
}
