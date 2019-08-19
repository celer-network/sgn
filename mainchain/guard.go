// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mainchain

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// GuardABI is the input ABI used to generate the binding from.
const GuardABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"subscribe\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"subscriptionExpiration\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"securityDeposit\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// GuardBin is the compiled bytecode used for deploying new contracts.
const GuardBin = `0x608060405234801561001057600080fd5b506101be806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80630f574ba71461005157806355bf3b3314610070578063b6b55f25146100a8578063f1f5da84146100c5575b600080fd5b61006e6004803603602081101561006757600080fd5b50356100eb565b005b6100966004803603602081101561008657600080fd5b50356001600160a01b031661014d565b60408051918252519081900360200190f35b61006e600480360360208110156100be57600080fd5b503561015f565b610096600480360360208110156100db57600080fd5b50356001600160a01b0316610177565b33600081815260208190526040902054600a830490431115610129576001600160a01b03821660009081526020819052604090204382019055610148565b6001600160a01b03821660009081526020819052604090208054820190555b505050565b60006020819052908152604090205481565b33600090815260016020526040902080549091019055565b6001602052600090815260409020548156fea265627a7a72305820da6e0293c282a4d54943475accfec14cd86f5023054347be15668cbe1104432764736f6c634300050a0032`

// DeployGuard deploys a new Ethereum contract, binding an instance of Guard to it.
func DeployGuard(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Guard, error) {
	parsed, err := abi.JSON(strings.NewReader(GuardABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(GuardBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Guard{GuardCaller: GuardCaller{contract: contract}, GuardTransactor: GuardTransactor{contract: contract}, GuardFilterer: GuardFilterer{contract: contract}}, nil
}

// Guard is an auto generated Go binding around an Ethereum contract.
type Guard struct {
	GuardCaller     // Read-only binding to the contract
	GuardTransactor // Write-only binding to the contract
	GuardFilterer   // Log filterer for contract events
}

// GuardCaller is an auto generated read-only Go binding around an Ethereum contract.
type GuardCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GuardTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GuardTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GuardFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GuardFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GuardSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GuardSession struct {
	Contract     *Guard            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GuardCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GuardCallerSession struct {
	Contract *GuardCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// GuardTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GuardTransactorSession struct {
	Contract     *GuardTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GuardRaw is an auto generated low-level Go binding around an Ethereum contract.
type GuardRaw struct {
	Contract *Guard // Generic contract binding to access the raw methods on
}

// GuardCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GuardCallerRaw struct {
	Contract *GuardCaller // Generic read-only contract binding to access the raw methods on
}

// GuardTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GuardTransactorRaw struct {
	Contract *GuardTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGuard creates a new instance of Guard, bound to a specific deployed contract.
func NewGuard(address common.Address, backend bind.ContractBackend) (*Guard, error) {
	contract, err := bindGuard(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Guard{GuardCaller: GuardCaller{contract: contract}, GuardTransactor: GuardTransactor{contract: contract}, GuardFilterer: GuardFilterer{contract: contract}}, nil
}

// NewGuardCaller creates a new read-only instance of Guard, bound to a specific deployed contract.
func NewGuardCaller(address common.Address, caller bind.ContractCaller) (*GuardCaller, error) {
	contract, err := bindGuard(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GuardCaller{contract: contract}, nil
}

// NewGuardTransactor creates a new write-only instance of Guard, bound to a specific deployed contract.
func NewGuardTransactor(address common.Address, transactor bind.ContractTransactor) (*GuardTransactor, error) {
	contract, err := bindGuard(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GuardTransactor{contract: contract}, nil
}

// NewGuardFilterer creates a new log filterer instance of Guard, bound to a specific deployed contract.
func NewGuardFilterer(address common.Address, filterer bind.ContractFilterer) (*GuardFilterer, error) {
	contract, err := bindGuard(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GuardFilterer{contract: contract}, nil
}

// bindGuard binds a generic wrapper to an already deployed contract.
func bindGuard(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GuardABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Guard *GuardRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Guard.Contract.GuardCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Guard *GuardRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Guard.Contract.GuardTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Guard *GuardRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Guard.Contract.GuardTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Guard *GuardCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Guard.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Guard *GuardTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Guard.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Guard *GuardTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Guard.Contract.contract.Transact(opts, method, params...)
}

// SecurityDeposit is a free data retrieval call binding the contract method 0xf1f5da84.
//
// Solidity: function securityDeposit(address ) constant returns(uint256)
func (_Guard *GuardCaller) SecurityDeposit(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "securityDeposit", arg0)
	return *ret0, err
}

// SecurityDeposit is a free data retrieval call binding the contract method 0xf1f5da84.
//
// Solidity: function securityDeposit(address ) constant returns(uint256)
func (_Guard *GuardSession) SecurityDeposit(arg0 common.Address) (*big.Int, error) {
	return _Guard.Contract.SecurityDeposit(&_Guard.CallOpts, arg0)
}

// SecurityDeposit is a free data retrieval call binding the contract method 0xf1f5da84.
//
// Solidity: function securityDeposit(address ) constant returns(uint256)
func (_Guard *GuardCallerSession) SecurityDeposit(arg0 common.Address) (*big.Int, error) {
	return _Guard.Contract.SecurityDeposit(&_Guard.CallOpts, arg0)
}

// SubscriptionExpiration is a free data retrieval call binding the contract method 0x55bf3b33.
//
// Solidity: function subscriptionExpiration(address ) constant returns(uint256)
func (_Guard *GuardCaller) SubscriptionExpiration(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "subscriptionExpiration", arg0)
	return *ret0, err
}

// SubscriptionExpiration is a free data retrieval call binding the contract method 0x55bf3b33.
//
// Solidity: function subscriptionExpiration(address ) constant returns(uint256)
func (_Guard *GuardSession) SubscriptionExpiration(arg0 common.Address) (*big.Int, error) {
	return _Guard.Contract.SubscriptionExpiration(&_Guard.CallOpts, arg0)
}

// SubscriptionExpiration is a free data retrieval call binding the contract method 0x55bf3b33.
//
// Solidity: function subscriptionExpiration(address ) constant returns(uint256)
func (_Guard *GuardCallerSession) SubscriptionExpiration(arg0 common.Address) (*big.Int, error) {
	return _Guard.Contract.SubscriptionExpiration(&_Guard.CallOpts, arg0)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) returns()
func (_Guard *GuardTransactor) Deposit(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "deposit", _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) returns()
func (_Guard *GuardSession) Deposit(_amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.Deposit(&_Guard.TransactOpts, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) returns()
func (_Guard *GuardTransactorSession) Deposit(_amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.Deposit(&_Guard.TransactOpts, _amount)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(uint256 _amount) returns()
func (_Guard *GuardTransactor) Subscribe(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "subscribe", _amount)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(uint256 _amount) returns()
func (_Guard *GuardSession) Subscribe(_amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.Subscribe(&_Guard.TransactOpts, _amount)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(uint256 _amount) returns()
func (_Guard *GuardTransactorSession) Subscribe(_amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.Subscribe(&_Guard.TransactOpts, _amount)
}
