// Code generated - DO NOT EDIT.
// This file is a generated binding with CelerLedgerBin cleared

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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// CelerLedgerABI is the input ABI used to generate the binding from.
const CelerLedgerABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_ethPool\",\"type\":\"address\"},{\"name\":\"_payRegistry\",\"type\":\"address\"},{\"name\":\"_celerWallet\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"tokenType\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"peerAddrs\",\"type\":\"address[2]\"},{\"indexed\":false,\"name\":\"initialDeposits\",\"type\":\"uint256[2]\"}],\"name\":\"OpenChannel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"peerAddrs\",\"type\":\"address[2]\"},{\"indexed\":false,\"name\":\"deposits\",\"type\":\"uint256[2]\"},{\"indexed\":false,\"name\":\"withdrawals\",\"type\":\"uint256[2]\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"seqNums\",\"type\":\"uint256[2]\"}],\"name\":\"SnapshotStates\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"seqNums\",\"type\":\"uint256[2]\"}],\"name\":\"IntendSettle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"payId\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"peerFrom\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ClearOnePay\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"settleBalance\",\"type\":\"uint256[2]\"}],\"name\":\"ConfirmSettle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"ConfirmSettleFail\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"IntendWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"withdrawnAmount\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"recipientChannelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"deposits\",\"type\":\"uint256[2]\"},{\"indexed\":false,\"name\":\"withdrawals\",\"type\":\"uint256[2]\"}],\"name\":\"ConfirmWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"VetoWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"withdrawnAmount\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"recipientChannelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"deposits\",\"type\":\"uint256[2]\"},{\"indexed\":false,\"name\":\"withdrawals\",\"type\":\"uint256[2]\"},{\"indexed\":false,\"name\":\"seqNum\",\"type\":\"uint256\"}],\"name\":\"CooperativeWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"settleBalance\",\"type\":\"uint256[2]\"}],\"name\":\"CooperativeSettle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"newLedgerAddr\",\"type\":\"address\"}],\"name\":\"MigrateChannelTo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"oldLedgerAddr\",\"type\":\"address\"}],\"name\":\"MigrateChannelFrom\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenAddrs\",\"type\":\"address[]\"},{\"name\":\"_limits\",\"type\":\"uint256[]\"}],\"name\":\"setBalanceLimits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"disableBalanceLimits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"enableBalanceLimits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_openRequest\",\"type\":\"bytes\"}],\"name\":\"openChannel\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"},{\"name\":\"_receiver\",\"type\":\"address\"},{\"name\":\"_transferFromAmount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_channelIds\",\"type\":\"bytes32[]\"},{\"name\":\"_receivers\",\"type\":\"address[]\"},{\"name\":\"_transferFromAmounts\",\"type\":\"uint256[]\"}],\"name\":\"depositInBatch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_signedSimplexStateArray\",\"type\":\"bytes\"}],\"name\":\"snapshotStates\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"},{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_recipientChannelId\",\"type\":\"bytes32\"}],\"name\":\"intendWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"confirmWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"vetoWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_cooperativeWithdrawRequest\",\"type\":\"bytes\"}],\"name\":\"cooperativeWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_signedSimplexStateArray\",\"type\":\"bytes\"}],\"name\":\"intendSettle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"},{\"name\":\"_peerFrom\",\"type\":\"address\"},{\"name\":\"_payIdList\",\"type\":\"bytes\"}],\"name\":\"clearPays\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"confirmSettle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_settleRequest\",\"type\":\"bytes\"}],\"name\":\"cooperativeSettle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_migrationRequest\",\"type\":\"bytes\"}],\"name\":\"migrateChannelTo\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_fromLedgerAddr\",\"type\":\"address\"},{\"name\":\"_migrationRequest\",\"type\":\"bytes\"}],\"name\":\"migrateChannelFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getSettleFinalizedTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getTokenContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getTokenType\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getChannelStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getCooperativeWithdrawSeqNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getTotalBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getBalanceMap\",\"outputs\":[{\"name\":\"\",\"type\":\"address[2]\"},{\"name\":\"\",\"type\":\"uint256[2]\"},{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getChannelMigrationArgs\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getPeersMigrationInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"address[2]\"},{\"name\":\"\",\"type\":\"uint256[2]\"},{\"name\":\"\",\"type\":\"uint256[2]\"},{\"name\":\"\",\"type\":\"uint256[2]\"},{\"name\":\"\",\"type\":\"uint256[2]\"},{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getDisputeTimeout\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getMigratedTo\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getStateSeqNumMap\",\"outputs\":[{\"name\":\"\",\"type\":\"address[2]\"},{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getTransferOutMap\",\"outputs\":[{\"name\":\"\",\"type\":\"address[2]\"},{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getNextPayIdListHashMap\",\"outputs\":[{\"name\":\"\",\"type\":\"address[2]\"},{\"name\":\"\",\"type\":\"bytes32[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getLastPayResolveDeadlineMap\",\"outputs\":[{\"name\":\"\",\"type\":\"address[2]\"},{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getPendingPayOutMap\",\"outputs\":[{\"name\":\"\",\"type\":\"address[2]\"},{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelId\",\"type\":\"bytes32\"}],\"name\":\"getWithdrawIntent\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_channelStatus\",\"type\":\"uint256\"}],\"name\":\"getChannelStatusNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getEthPool\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPayRegistry\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getCelerWallet\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenAddr\",\"type\":\"address\"}],\"name\":\"getBalanceLimit\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getBalanceLimitsEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// CelerLedgerBin is the compiled bytecode used for deploying new contracts.
var CelerLedgerBin = ""

// DeployCelerLedger deploys a new Ethereum contract, binding an instance of CelerLedger to it.
func DeployCelerLedger(auth *bind.TransactOpts, backend bind.ContractBackend, _ethPool common.Address, _payRegistry common.Address, _celerWallet common.Address) (common.Address, *types.Transaction, *CelerLedger, error) {
	parsed, err := abi.JSON(strings.NewReader(CelerLedgerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CelerLedgerBin), backend, _ethPool, _payRegistry, _celerWallet)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CelerLedger{CelerLedgerCaller: CelerLedgerCaller{contract: contract}, CelerLedgerTransactor: CelerLedgerTransactor{contract: contract}, CelerLedgerFilterer: CelerLedgerFilterer{contract: contract}}, nil
}

// CelerLedger is an auto generated Go binding around an Ethereum contract.
type CelerLedger struct {
	CelerLedgerCaller     // Read-only binding to the contract
	CelerLedgerTransactor // Write-only binding to the contract
	CelerLedgerFilterer   // Log filterer for contract events
}

// CelerLedgerCaller is an auto generated read-only Go binding around an Ethereum contract.
type CelerLedgerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CelerLedgerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CelerLedgerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CelerLedgerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CelerLedgerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CelerLedgerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CelerLedgerSession struct {
	Contract     *CelerLedger      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CelerLedgerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CelerLedgerCallerSession struct {
	Contract *CelerLedgerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// CelerLedgerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CelerLedgerTransactorSession struct {
	Contract     *CelerLedgerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// CelerLedgerRaw is an auto generated low-level Go binding around an Ethereum contract.
type CelerLedgerRaw struct {
	Contract *CelerLedger // Generic contract binding to access the raw methods on
}

// CelerLedgerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CelerLedgerCallerRaw struct {
	Contract *CelerLedgerCaller // Generic read-only contract binding to access the raw methods on
}

// CelerLedgerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CelerLedgerTransactorRaw struct {
	Contract *CelerLedgerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCelerLedger creates a new instance of CelerLedger, bound to a specific deployed contract.
func NewCelerLedger(address common.Address, backend bind.ContractBackend) (*CelerLedger, error) {
	contract, err := bindCelerLedger(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CelerLedger{CelerLedgerCaller: CelerLedgerCaller{contract: contract}, CelerLedgerTransactor: CelerLedgerTransactor{contract: contract}, CelerLedgerFilterer: CelerLedgerFilterer{contract: contract}}, nil
}

// NewCelerLedgerCaller creates a new read-only instance of CelerLedger, bound to a specific deployed contract.
func NewCelerLedgerCaller(address common.Address, caller bind.ContractCaller) (*CelerLedgerCaller, error) {
	contract, err := bindCelerLedger(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerCaller{contract: contract}, nil
}

// NewCelerLedgerTransactor creates a new write-only instance of CelerLedger, bound to a specific deployed contract.
func NewCelerLedgerTransactor(address common.Address, transactor bind.ContractTransactor) (*CelerLedgerTransactor, error) {
	contract, err := bindCelerLedger(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerTransactor{contract: contract}, nil
}

// NewCelerLedgerFilterer creates a new log filterer instance of CelerLedger, bound to a specific deployed contract.
func NewCelerLedgerFilterer(address common.Address, filterer bind.ContractFilterer) (*CelerLedgerFilterer, error) {
	contract, err := bindCelerLedger(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerFilterer{contract: contract}, nil
}

// bindCelerLedger binds a generic wrapper to an already deployed contract.
func bindCelerLedger(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CelerLedgerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CelerLedger *CelerLedgerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CelerLedger.Contract.CelerLedgerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CelerLedger *CelerLedgerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CelerLedger.Contract.CelerLedgerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CelerLedger *CelerLedgerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CelerLedger.Contract.CelerLedgerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CelerLedger *CelerLedgerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CelerLedger.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CelerLedger *CelerLedgerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CelerLedger.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CelerLedger *CelerLedgerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CelerLedger.Contract.contract.Transact(opts, method, params...)
}

// GetBalanceLimit is a free data retrieval call binding the contract method 0xec7c637d.
//
// Solidity: function getBalanceLimit(address _tokenAddr) constant returns(uint256)
func (_CelerLedger *CelerLedgerCaller) GetBalanceLimit(opts *bind.CallOpts, _tokenAddr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getBalanceLimit", _tokenAddr)
	return *ret0, err
}

// GetBalanceLimit is a free data retrieval call binding the contract method 0xec7c637d.
//
// Solidity: function getBalanceLimit(address _tokenAddr) constant returns(uint256)
func (_CelerLedger *CelerLedgerSession) GetBalanceLimit(_tokenAddr common.Address) (*big.Int, error) {
	return _CelerLedger.Contract.GetBalanceLimit(&_CelerLedger.CallOpts, _tokenAddr)
}

// GetBalanceLimit is a free data retrieval call binding the contract method 0xec7c637d.
//
// Solidity: function getBalanceLimit(address _tokenAddr) constant returns(uint256)
func (_CelerLedger *CelerLedgerCallerSession) GetBalanceLimit(_tokenAddr common.Address) (*big.Int, error) {
	return _CelerLedger.Contract.GetBalanceLimit(&_CelerLedger.CallOpts, _tokenAddr)
}

// GetBalanceLimitsEnabled is a free data retrieval call binding the contract method 0x2b559ecc.
//
// Solidity: function getBalanceLimitsEnabled() constant returns(bool)
func (_CelerLedger *CelerLedgerCaller) GetBalanceLimitsEnabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getBalanceLimitsEnabled")
	return *ret0, err
}

// GetBalanceLimitsEnabled is a free data retrieval call binding the contract method 0x2b559ecc.
//
// Solidity: function getBalanceLimitsEnabled() constant returns(bool)
func (_CelerLedger *CelerLedgerSession) GetBalanceLimitsEnabled() (bool, error) {
	return _CelerLedger.Contract.GetBalanceLimitsEnabled(&_CelerLedger.CallOpts)
}

// GetBalanceLimitsEnabled is a free data retrieval call binding the contract method 0x2b559ecc.
//
// Solidity: function getBalanceLimitsEnabled() constant returns(bool)
func (_CelerLedger *CelerLedgerCallerSession) GetBalanceLimitsEnabled() (bool, error) {
	return _CelerLedger.Contract.GetBalanceLimitsEnabled(&_CelerLedger.CallOpts)
}

// GetBalanceMap is a free data retrieval call binding the contract method 0xbd480cb7.
//
// Solidity: function getBalanceMap(bytes32 _channelId) constant returns(address[2], uint256[2], uint256[2])
func (_CelerLedger *CelerLedgerCaller) GetBalanceMap(opts *bind.CallOpts, _channelId [32]byte) ([2]common.Address, [2]*big.Int, [2]*big.Int, error) {
	var (
		ret0 = new([2]common.Address)
		ret1 = new([2]*big.Int)
		ret2 = new([2]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _CelerLedger.contract.Call(opts, out, "getBalanceMap", _channelId)
	return *ret0, *ret1, *ret2, err
}

// GetBalanceMap is a free data retrieval call binding the contract method 0xbd480cb7.
//
// Solidity: function getBalanceMap(bytes32 _channelId) constant returns(address[2], uint256[2], uint256[2])
func (_CelerLedger *CelerLedgerSession) GetBalanceMap(_channelId [32]byte) ([2]common.Address, [2]*big.Int, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetBalanceMap(&_CelerLedger.CallOpts, _channelId)
}

// GetBalanceMap is a free data retrieval call binding the contract method 0xbd480cb7.
//
// Solidity: function getBalanceMap(bytes32 _channelId) constant returns(address[2], uint256[2], uint256[2])
func (_CelerLedger *CelerLedgerCallerSession) GetBalanceMap(_channelId [32]byte) ([2]common.Address, [2]*big.Int, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetBalanceMap(&_CelerLedger.CallOpts, _channelId)
}

// GetCelerWallet is a free data retrieval call binding the contract method 0xa099a39f.
//
// Solidity: function getCelerWallet() constant returns(address)
func (_CelerLedger *CelerLedgerCaller) GetCelerWallet(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getCelerWallet")
	return *ret0, err
}

// GetCelerWallet is a free data retrieval call binding the contract method 0xa099a39f.
//
// Solidity: function getCelerWallet() constant returns(address)
func (_CelerLedger *CelerLedgerSession) GetCelerWallet() (common.Address, error) {
	return _CelerLedger.Contract.GetCelerWallet(&_CelerLedger.CallOpts)
}

// GetCelerWallet is a free data retrieval call binding the contract method 0xa099a39f.
//
// Solidity: function getCelerWallet() constant returns(address)
func (_CelerLedger *CelerLedgerCallerSession) GetCelerWallet() (common.Address, error) {
	return _CelerLedger.Contract.GetCelerWallet(&_CelerLedger.CallOpts)
}

// GetChannelMigrationArgs is a free data retrieval call binding the contract method 0x2f0ac304.
//
// Solidity: function getChannelMigrationArgs(bytes32 _channelId) constant returns(uint256, uint256, address, uint256)
func (_CelerLedger *CelerLedgerCaller) GetChannelMigrationArgs(opts *bind.CallOpts, _channelId [32]byte) (*big.Int, *big.Int, common.Address, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new(common.Address)
		ret3 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _CelerLedger.contract.Call(opts, out, "getChannelMigrationArgs", _channelId)
	return *ret0, *ret1, *ret2, *ret3, err
}

// GetChannelMigrationArgs is a free data retrieval call binding the contract method 0x2f0ac304.
//
// Solidity: function getChannelMigrationArgs(bytes32 _channelId) constant returns(uint256, uint256, address, uint256)
func (_CelerLedger *CelerLedgerSession) GetChannelMigrationArgs(_channelId [32]byte) (*big.Int, *big.Int, common.Address, *big.Int, error) {
	return _CelerLedger.Contract.GetChannelMigrationArgs(&_CelerLedger.CallOpts, _channelId)
}

// GetChannelMigrationArgs is a free data retrieval call binding the contract method 0x2f0ac304.
//
// Solidity: function getChannelMigrationArgs(bytes32 _channelId) constant returns(uint256, uint256, address, uint256)
func (_CelerLedger *CelerLedgerCallerSession) GetChannelMigrationArgs(_channelId [32]byte) (*big.Int, *big.Int, common.Address, *big.Int, error) {
	return _CelerLedger.Contract.GetChannelMigrationArgs(&_CelerLedger.CallOpts, _channelId)
}

// GetChannelStatus is a free data retrieval call binding the contract method 0xcc0b94b7.
//
// Solidity: function getChannelStatus(bytes32 _channelId) constant returns(uint8)
func (_CelerLedger *CelerLedgerCaller) GetChannelStatus(opts *bind.CallOpts, _channelId [32]byte) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getChannelStatus", _channelId)
	return *ret0, err
}

// GetChannelStatus is a free data retrieval call binding the contract method 0xcc0b94b7.
//
// Solidity: function getChannelStatus(bytes32 _channelId) constant returns(uint8)
func (_CelerLedger *CelerLedgerSession) GetChannelStatus(_channelId [32]byte) (uint8, error) {
	return _CelerLedger.Contract.GetChannelStatus(&_CelerLedger.CallOpts, _channelId)
}

// GetChannelStatus is a free data retrieval call binding the contract method 0xcc0b94b7.
//
// Solidity: function getChannelStatus(bytes32 _channelId) constant returns(uint8)
func (_CelerLedger *CelerLedgerCallerSession) GetChannelStatus(_channelId [32]byte) (uint8, error) {
	return _CelerLedger.Contract.GetChannelStatus(&_CelerLedger.CallOpts, _channelId)
}

// GetChannelStatusNum is a free data retrieval call binding the contract method 0x307d6f96.
//
// Solidity: function getChannelStatusNum(uint256 _channelStatus) constant returns(uint256)
func (_CelerLedger *CelerLedgerCaller) GetChannelStatusNum(opts *bind.CallOpts, _channelStatus *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getChannelStatusNum", _channelStatus)
	return *ret0, err
}

// GetChannelStatusNum is a free data retrieval call binding the contract method 0x307d6f96.
//
// Solidity: function getChannelStatusNum(uint256 _channelStatus) constant returns(uint256)
func (_CelerLedger *CelerLedgerSession) GetChannelStatusNum(_channelStatus *big.Int) (*big.Int, error) {
	return _CelerLedger.Contract.GetChannelStatusNum(&_CelerLedger.CallOpts, _channelStatus)
}

// GetChannelStatusNum is a free data retrieval call binding the contract method 0x307d6f96.
//
// Solidity: function getChannelStatusNum(uint256 _channelStatus) constant returns(uint256)
func (_CelerLedger *CelerLedgerCallerSession) GetChannelStatusNum(_channelStatus *big.Int) (*big.Int, error) {
	return _CelerLedger.Contract.GetChannelStatusNum(&_CelerLedger.CallOpts, _channelStatus)
}

// GetCooperativeWithdrawSeqNum is a free data retrieval call binding the contract method 0xf0c73d70.
//
// Solidity: function getCooperativeWithdrawSeqNum(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerCaller) GetCooperativeWithdrawSeqNum(opts *bind.CallOpts, _channelId [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getCooperativeWithdrawSeqNum", _channelId)
	return *ret0, err
}

// GetCooperativeWithdrawSeqNum is a free data retrieval call binding the contract method 0xf0c73d70.
//
// Solidity: function getCooperativeWithdrawSeqNum(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerSession) GetCooperativeWithdrawSeqNum(_channelId [32]byte) (*big.Int, error) {
	return _CelerLedger.Contract.GetCooperativeWithdrawSeqNum(&_CelerLedger.CallOpts, _channelId)
}

// GetCooperativeWithdrawSeqNum is a free data retrieval call binding the contract method 0xf0c73d70.
//
// Solidity: function getCooperativeWithdrawSeqNum(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerCallerSession) GetCooperativeWithdrawSeqNum(_channelId [32]byte) (*big.Int, error) {
	return _CelerLedger.Contract.GetCooperativeWithdrawSeqNum(&_CelerLedger.CallOpts, _channelId)
}

// GetDisputeTimeout is a free data retrieval call binding the contract method 0xe6322df7.
//
// Solidity: function getDisputeTimeout(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerCaller) GetDisputeTimeout(opts *bind.CallOpts, _channelId [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getDisputeTimeout", _channelId)
	return *ret0, err
}

// GetDisputeTimeout is a free data retrieval call binding the contract method 0xe6322df7.
//
// Solidity: function getDisputeTimeout(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerSession) GetDisputeTimeout(_channelId [32]byte) (*big.Int, error) {
	return _CelerLedger.Contract.GetDisputeTimeout(&_CelerLedger.CallOpts, _channelId)
}

// GetDisputeTimeout is a free data retrieval call binding the contract method 0xe6322df7.
//
// Solidity: function getDisputeTimeout(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerCallerSession) GetDisputeTimeout(_channelId [32]byte) (*big.Int, error) {
	return _CelerLedger.Contract.GetDisputeTimeout(&_CelerLedger.CallOpts, _channelId)
}

// GetEthPool is a free data retrieval call binding the contract method 0xd75f960e.
//
// Solidity: function getEthPool() constant returns(address)
func (_CelerLedger *CelerLedgerCaller) GetEthPool(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getEthPool")
	return *ret0, err
}

// GetEthPool is a free data retrieval call binding the contract method 0xd75f960e.
//
// Solidity: function getEthPool() constant returns(address)
func (_CelerLedger *CelerLedgerSession) GetEthPool() (common.Address, error) {
	return _CelerLedger.Contract.GetEthPool(&_CelerLedger.CallOpts)
}

// GetEthPool is a free data retrieval call binding the contract method 0xd75f960e.
//
// Solidity: function getEthPool() constant returns(address)
func (_CelerLedger *CelerLedgerCallerSession) GetEthPool() (common.Address, error) {
	return _CelerLedger.Contract.GetEthPool(&_CelerLedger.CallOpts)
}

// GetLastPayResolveDeadlineMap is a free data retrieval call binding the contract method 0x9f1fad83.
//
// Solidity: function getLastPayResolveDeadlineMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerCaller) GetLastPayResolveDeadlineMap(opts *bind.CallOpts, _channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	var (
		ret0 = new([2]common.Address)
		ret1 = new([2]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _CelerLedger.contract.Call(opts, out, "getLastPayResolveDeadlineMap", _channelId)
	return *ret0, *ret1, err
}

// GetLastPayResolveDeadlineMap is a free data retrieval call binding the contract method 0x9f1fad83.
//
// Solidity: function getLastPayResolveDeadlineMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerSession) GetLastPayResolveDeadlineMap(_channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetLastPayResolveDeadlineMap(&_CelerLedger.CallOpts, _channelId)
}

// GetLastPayResolveDeadlineMap is a free data retrieval call binding the contract method 0x9f1fad83.
//
// Solidity: function getLastPayResolveDeadlineMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerCallerSession) GetLastPayResolveDeadlineMap(_channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetLastPayResolveDeadlineMap(&_CelerLedger.CallOpts, _channelId)
}

// GetMigratedTo is a free data retrieval call binding the contract method 0xc38a325d.
//
// Solidity: function getMigratedTo(bytes32 _channelId) constant returns(address)
func (_CelerLedger *CelerLedgerCaller) GetMigratedTo(opts *bind.CallOpts, _channelId [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getMigratedTo", _channelId)
	return *ret0, err
}

// GetMigratedTo is a free data retrieval call binding the contract method 0xc38a325d.
//
// Solidity: function getMigratedTo(bytes32 _channelId) constant returns(address)
func (_CelerLedger *CelerLedgerSession) GetMigratedTo(_channelId [32]byte) (common.Address, error) {
	return _CelerLedger.Contract.GetMigratedTo(&_CelerLedger.CallOpts, _channelId)
}

// GetMigratedTo is a free data retrieval call binding the contract method 0xc38a325d.
//
// Solidity: function getMigratedTo(bytes32 _channelId) constant returns(address)
func (_CelerLedger *CelerLedgerCallerSession) GetMigratedTo(_channelId [32]byte) (common.Address, error) {
	return _CelerLedger.Contract.GetMigratedTo(&_CelerLedger.CallOpts, _channelId)
}

// GetNextPayIdListHashMap is a free data retrieval call binding the contract method 0xcd3a1be6.
//
// Solidity: function getNextPayIdListHashMap(bytes32 _channelId) constant returns(address[2], bytes32[2])
func (_CelerLedger *CelerLedgerCaller) GetNextPayIdListHashMap(opts *bind.CallOpts, _channelId [32]byte) ([2]common.Address, [2][32]byte, error) {
	var (
		ret0 = new([2]common.Address)
		ret1 = new([2][32]byte)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _CelerLedger.contract.Call(opts, out, "getNextPayIdListHashMap", _channelId)
	return *ret0, *ret1, err
}

// GetNextPayIdListHashMap is a free data retrieval call binding the contract method 0xcd3a1be6.
//
// Solidity: function getNextPayIdListHashMap(bytes32 _channelId) constant returns(address[2], bytes32[2])
func (_CelerLedger *CelerLedgerSession) GetNextPayIdListHashMap(_channelId [32]byte) ([2]common.Address, [2][32]byte, error) {
	return _CelerLedger.Contract.GetNextPayIdListHashMap(&_CelerLedger.CallOpts, _channelId)
}

// GetNextPayIdListHashMap is a free data retrieval call binding the contract method 0xcd3a1be6.
//
// Solidity: function getNextPayIdListHashMap(bytes32 _channelId) constant returns(address[2], bytes32[2])
func (_CelerLedger *CelerLedgerCallerSession) GetNextPayIdListHashMap(_channelId [32]byte) ([2]common.Address, [2][32]byte, error) {
	return _CelerLedger.Contract.GetNextPayIdListHashMap(&_CelerLedger.CallOpts, _channelId)
}

// GetPayRegistry is a free data retrieval call binding the contract method 0x312ea2c6.
//
// Solidity: function getPayRegistry() constant returns(address)
func (_CelerLedger *CelerLedgerCaller) GetPayRegistry(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getPayRegistry")
	return *ret0, err
}

// GetPayRegistry is a free data retrieval call binding the contract method 0x312ea2c6.
//
// Solidity: function getPayRegistry() constant returns(address)
func (_CelerLedger *CelerLedgerSession) GetPayRegistry() (common.Address, error) {
	return _CelerLedger.Contract.GetPayRegistry(&_CelerLedger.CallOpts)
}

// GetPayRegistry is a free data retrieval call binding the contract method 0x312ea2c6.
//
// Solidity: function getPayRegistry() constant returns(address)
func (_CelerLedger *CelerLedgerCallerSession) GetPayRegistry() (common.Address, error) {
	return _CelerLedger.Contract.GetPayRegistry(&_CelerLedger.CallOpts)
}

// GetPeersMigrationInfo is a free data retrieval call binding the contract method 0x88f41465.
//
// Solidity: function getPeersMigrationInfo(bytes32 _channelId) constant returns(address[2], uint256[2], uint256[2], uint256[2], uint256[2], uint256[2])
func (_CelerLedger *CelerLedgerCaller) GetPeersMigrationInfo(opts *bind.CallOpts, _channelId [32]byte) ([2]common.Address, [2]*big.Int, [2]*big.Int, [2]*big.Int, [2]*big.Int, [2]*big.Int, error) {
	var (
		ret0 = new([2]common.Address)
		ret1 = new([2]*big.Int)
		ret2 = new([2]*big.Int)
		ret3 = new([2]*big.Int)
		ret4 = new([2]*big.Int)
		ret5 = new([2]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
		ret5,
	}
	err := _CelerLedger.contract.Call(opts, out, "getPeersMigrationInfo", _channelId)
	return *ret0, *ret1, *ret2, *ret3, *ret4, *ret5, err
}

// GetPeersMigrationInfo is a free data retrieval call binding the contract method 0x88f41465.
//
// Solidity: function getPeersMigrationInfo(bytes32 _channelId) constant returns(address[2], uint256[2], uint256[2], uint256[2], uint256[2], uint256[2])
func (_CelerLedger *CelerLedgerSession) GetPeersMigrationInfo(_channelId [32]byte) ([2]common.Address, [2]*big.Int, [2]*big.Int, [2]*big.Int, [2]*big.Int, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetPeersMigrationInfo(&_CelerLedger.CallOpts, _channelId)
}

// GetPeersMigrationInfo is a free data retrieval call binding the contract method 0x88f41465.
//
// Solidity: function getPeersMigrationInfo(bytes32 _channelId) constant returns(address[2], uint256[2], uint256[2], uint256[2], uint256[2], uint256[2])
func (_CelerLedger *CelerLedgerCallerSession) GetPeersMigrationInfo(_channelId [32]byte) ([2]common.Address, [2]*big.Int, [2]*big.Int, [2]*big.Int, [2]*big.Int, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetPeersMigrationInfo(&_CelerLedger.CallOpts, _channelId)
}

// GetPendingPayOutMap is a free data retrieval call binding the contract method 0x0165cef8.
//
// Solidity: function getPendingPayOutMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerCaller) GetPendingPayOutMap(opts *bind.CallOpts, _channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	var (
		ret0 = new([2]common.Address)
		ret1 = new([2]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _CelerLedger.contract.Call(opts, out, "getPendingPayOutMap", _channelId)
	return *ret0, *ret1, err
}

// GetPendingPayOutMap is a free data retrieval call binding the contract method 0x0165cef8.
//
// Solidity: function getPendingPayOutMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerSession) GetPendingPayOutMap(_channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetPendingPayOutMap(&_CelerLedger.CallOpts, _channelId)
}

// GetPendingPayOutMap is a free data retrieval call binding the contract method 0x0165cef8.
//
// Solidity: function getPendingPayOutMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerCallerSession) GetPendingPayOutMap(_channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetPendingPayOutMap(&_CelerLedger.CallOpts, _channelId)
}

// GetSettleFinalizedTime is a free data retrieval call binding the contract method 0x09b65d86.
//
// Solidity: function getSettleFinalizedTime(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerCaller) GetSettleFinalizedTime(opts *bind.CallOpts, _channelId [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getSettleFinalizedTime", _channelId)
	return *ret0, err
}

// GetSettleFinalizedTime is a free data retrieval call binding the contract method 0x09b65d86.
//
// Solidity: function getSettleFinalizedTime(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerSession) GetSettleFinalizedTime(_channelId [32]byte) (*big.Int, error) {
	return _CelerLedger.Contract.GetSettleFinalizedTime(&_CelerLedger.CallOpts, _channelId)
}

// GetSettleFinalizedTime is a free data retrieval call binding the contract method 0x09b65d86.
//
// Solidity: function getSettleFinalizedTime(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerCallerSession) GetSettleFinalizedTime(_channelId [32]byte) (*big.Int, error) {
	return _CelerLedger.Contract.GetSettleFinalizedTime(&_CelerLedger.CallOpts, _channelId)
}

// GetStateSeqNumMap is a free data retrieval call binding the contract method 0x666a6d65.
//
// Solidity: function getStateSeqNumMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerCaller) GetStateSeqNumMap(opts *bind.CallOpts, _channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	var (
		ret0 = new([2]common.Address)
		ret1 = new([2]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _CelerLedger.contract.Call(opts, out, "getStateSeqNumMap", _channelId)
	return *ret0, *ret1, err
}

// GetStateSeqNumMap is a free data retrieval call binding the contract method 0x666a6d65.
//
// Solidity: function getStateSeqNumMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerSession) GetStateSeqNumMap(_channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetStateSeqNumMap(&_CelerLedger.CallOpts, _channelId)
}

// GetStateSeqNumMap is a free data retrieval call binding the contract method 0x666a6d65.
//
// Solidity: function getStateSeqNumMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerCallerSession) GetStateSeqNumMap(_channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetStateSeqNumMap(&_CelerLedger.CallOpts, _channelId)
}

// GetTokenContract is a free data retrieval call binding the contract method 0x2e2a5a02.
//
// Solidity: function getTokenContract(bytes32 _channelId) constant returns(address)
func (_CelerLedger *CelerLedgerCaller) GetTokenContract(opts *bind.CallOpts, _channelId [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getTokenContract", _channelId)
	return *ret0, err
}

// GetTokenContract is a free data retrieval call binding the contract method 0x2e2a5a02.
//
// Solidity: function getTokenContract(bytes32 _channelId) constant returns(address)
func (_CelerLedger *CelerLedgerSession) GetTokenContract(_channelId [32]byte) (common.Address, error) {
	return _CelerLedger.Contract.GetTokenContract(&_CelerLedger.CallOpts, _channelId)
}

// GetTokenContract is a free data retrieval call binding the contract method 0x2e2a5a02.
//
// Solidity: function getTokenContract(bytes32 _channelId) constant returns(address)
func (_CelerLedger *CelerLedgerCallerSession) GetTokenContract(_channelId [32]byte) (common.Address, error) {
	return _CelerLedger.Contract.GetTokenContract(&_CelerLedger.CallOpts, _channelId)
}

// GetTokenType is a free data retrieval call binding the contract method 0x979a9b5e.
//
// Solidity: function getTokenType(bytes32 _channelId) constant returns(uint8)
func (_CelerLedger *CelerLedgerCaller) GetTokenType(opts *bind.CallOpts, _channelId [32]byte) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getTokenType", _channelId)
	return *ret0, err
}

// GetTokenType is a free data retrieval call binding the contract method 0x979a9b5e.
//
// Solidity: function getTokenType(bytes32 _channelId) constant returns(uint8)
func (_CelerLedger *CelerLedgerSession) GetTokenType(_channelId [32]byte) (uint8, error) {
	return _CelerLedger.Contract.GetTokenType(&_CelerLedger.CallOpts, _channelId)
}

// GetTokenType is a free data retrieval call binding the contract method 0x979a9b5e.
//
// Solidity: function getTokenType(bytes32 _channelId) constant returns(uint8)
func (_CelerLedger *CelerLedgerCallerSession) GetTokenType(_channelId [32]byte) (uint8, error) {
	return _CelerLedger.Contract.GetTokenType(&_CelerLedger.CallOpts, _channelId)
}

// GetTotalBalance is a free data retrieval call binding the contract method 0x69d5dd67.
//
// Solidity: function getTotalBalance(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerCaller) GetTotalBalance(opts *bind.CallOpts, _channelId [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "getTotalBalance", _channelId)
	return *ret0, err
}

// GetTotalBalance is a free data retrieval call binding the contract method 0x69d5dd67.
//
// Solidity: function getTotalBalance(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerSession) GetTotalBalance(_channelId [32]byte) (*big.Int, error) {
	return _CelerLedger.Contract.GetTotalBalance(&_CelerLedger.CallOpts, _channelId)
}

// GetTotalBalance is a free data retrieval call binding the contract method 0x69d5dd67.
//
// Solidity: function getTotalBalance(bytes32 _channelId) constant returns(uint256)
func (_CelerLedger *CelerLedgerCallerSession) GetTotalBalance(_channelId [32]byte) (*big.Int, error) {
	return _CelerLedger.Contract.GetTotalBalance(&_CelerLedger.CallOpts, _channelId)
}

// GetTransferOutMap is a free data retrieval call binding the contract method 0xd927bfc4.
//
// Solidity: function getTransferOutMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerCaller) GetTransferOutMap(opts *bind.CallOpts, _channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	var (
		ret0 = new([2]common.Address)
		ret1 = new([2]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _CelerLedger.contract.Call(opts, out, "getTransferOutMap", _channelId)
	return *ret0, *ret1, err
}

// GetTransferOutMap is a free data retrieval call binding the contract method 0xd927bfc4.
//
// Solidity: function getTransferOutMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerSession) GetTransferOutMap(_channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetTransferOutMap(&_CelerLedger.CallOpts, _channelId)
}

// GetTransferOutMap is a free data retrieval call binding the contract method 0xd927bfc4.
//
// Solidity: function getTransferOutMap(bytes32 _channelId) constant returns(address[2], uint256[2])
func (_CelerLedger *CelerLedgerCallerSession) GetTransferOutMap(_channelId [32]byte) ([2]common.Address, [2]*big.Int, error) {
	return _CelerLedger.Contract.GetTransferOutMap(&_CelerLedger.CallOpts, _channelId)
}

// GetWithdrawIntent is a free data retrieval call binding the contract method 0x76bff117.
//
// Solidity: function getWithdrawIntent(bytes32 _channelId) constant returns(address, uint256, uint256, bytes32)
func (_CelerLedger *CelerLedgerCaller) GetWithdrawIntent(opts *bind.CallOpts, _channelId [32]byte) (common.Address, *big.Int, *big.Int, [32]byte, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
		ret3 = new([32]byte)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _CelerLedger.contract.Call(opts, out, "getWithdrawIntent", _channelId)
	return *ret0, *ret1, *ret2, *ret3, err
}

// GetWithdrawIntent is a free data retrieval call binding the contract method 0x76bff117.
//
// Solidity: function getWithdrawIntent(bytes32 _channelId) constant returns(address, uint256, uint256, bytes32)
func (_CelerLedger *CelerLedgerSession) GetWithdrawIntent(_channelId [32]byte) (common.Address, *big.Int, *big.Int, [32]byte, error) {
	return _CelerLedger.Contract.GetWithdrawIntent(&_CelerLedger.CallOpts, _channelId)
}

// GetWithdrawIntent is a free data retrieval call binding the contract method 0x76bff117.
//
// Solidity: function getWithdrawIntent(bytes32 _channelId) constant returns(address, uint256, uint256, bytes32)
func (_CelerLedger *CelerLedgerCallerSession) GetWithdrawIntent(_channelId [32]byte) (common.Address, *big.Int, *big.Int, [32]byte, error) {
	return _CelerLedger.Contract.GetWithdrawIntent(&_CelerLedger.CallOpts, _channelId)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_CelerLedger *CelerLedgerCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_CelerLedger *CelerLedgerSession) IsOwner() (bool, error) {
	return _CelerLedger.Contract.IsOwner(&_CelerLedger.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_CelerLedger *CelerLedgerCallerSession) IsOwner() (bool, error) {
	return _CelerLedger.Contract.IsOwner(&_CelerLedger.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_CelerLedger *CelerLedgerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CelerLedger.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_CelerLedger *CelerLedgerSession) Owner() (common.Address, error) {
	return _CelerLedger.Contract.Owner(&_CelerLedger.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_CelerLedger *CelerLedgerCallerSession) Owner() (common.Address, error) {
	return _CelerLedger.Contract.Owner(&_CelerLedger.CallOpts)
}

// ClearPays is a paid mutator transaction binding the contract method 0xfd0a1a61.
//
// Solidity: function clearPays(bytes32 _channelId, address _peerFrom, bytes _payIdList) returns()
func (_CelerLedger *CelerLedgerTransactor) ClearPays(opts *bind.TransactOpts, _channelId [32]byte, _peerFrom common.Address, _payIdList []byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "clearPays", _channelId, _peerFrom, _payIdList)
}

// ClearPays is a paid mutator transaction binding the contract method 0xfd0a1a61.
//
// Solidity: function clearPays(bytes32 _channelId, address _peerFrom, bytes _payIdList) returns()
func (_CelerLedger *CelerLedgerSession) ClearPays(_channelId [32]byte, _peerFrom common.Address, _payIdList []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.ClearPays(&_CelerLedger.TransactOpts, _channelId, _peerFrom, _payIdList)
}

// ClearPays is a paid mutator transaction binding the contract method 0xfd0a1a61.
//
// Solidity: function clearPays(bytes32 _channelId, address _peerFrom, bytes _payIdList) returns()
func (_CelerLedger *CelerLedgerTransactorSession) ClearPays(_channelId [32]byte, _peerFrom common.Address, _payIdList []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.ClearPays(&_CelerLedger.TransactOpts, _channelId, _peerFrom, _payIdList)
}

// ConfirmSettle is a paid mutator transaction binding the contract method 0xc7ff8625.
//
// Solidity: function confirmSettle(bytes32 _channelId) returns()
func (_CelerLedger *CelerLedgerTransactor) ConfirmSettle(opts *bind.TransactOpts, _channelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "confirmSettle", _channelId)
}

// ConfirmSettle is a paid mutator transaction binding the contract method 0xc7ff8625.
//
// Solidity: function confirmSettle(bytes32 _channelId) returns()
func (_CelerLedger *CelerLedgerSession) ConfirmSettle(_channelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.ConfirmSettle(&_CelerLedger.TransactOpts, _channelId)
}

// ConfirmSettle is a paid mutator transaction binding the contract method 0xc7ff8625.
//
// Solidity: function confirmSettle(bytes32 _channelId) returns()
func (_CelerLedger *CelerLedgerTransactorSession) ConfirmSettle(_channelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.ConfirmSettle(&_CelerLedger.TransactOpts, _channelId)
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0x7e9a7a3e.
//
// Solidity: function confirmWithdraw(bytes32 _channelId) returns()
func (_CelerLedger *CelerLedgerTransactor) ConfirmWithdraw(opts *bind.TransactOpts, _channelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "confirmWithdraw", _channelId)
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0x7e9a7a3e.
//
// Solidity: function confirmWithdraw(bytes32 _channelId) returns()
func (_CelerLedger *CelerLedgerSession) ConfirmWithdraw(_channelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.ConfirmWithdraw(&_CelerLedger.TransactOpts, _channelId)
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0x7e9a7a3e.
//
// Solidity: function confirmWithdraw(bytes32 _channelId) returns()
func (_CelerLedger *CelerLedgerTransactorSession) ConfirmWithdraw(_channelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.ConfirmWithdraw(&_CelerLedger.TransactOpts, _channelId)
}

// CooperativeSettle is a paid mutator transaction binding the contract method 0x09683c03.
//
// Solidity: function cooperativeSettle(bytes _settleRequest) returns()
func (_CelerLedger *CelerLedgerTransactor) CooperativeSettle(opts *bind.TransactOpts, _settleRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "cooperativeSettle", _settleRequest)
}

// CooperativeSettle is a paid mutator transaction binding the contract method 0x09683c03.
//
// Solidity: function cooperativeSettle(bytes _settleRequest) returns()
func (_CelerLedger *CelerLedgerSession) CooperativeSettle(_settleRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.CooperativeSettle(&_CelerLedger.TransactOpts, _settleRequest)
}

// CooperativeSettle is a paid mutator transaction binding the contract method 0x09683c03.
//
// Solidity: function cooperativeSettle(bytes _settleRequest) returns()
func (_CelerLedger *CelerLedgerTransactorSession) CooperativeSettle(_settleRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.CooperativeSettle(&_CelerLedger.TransactOpts, _settleRequest)
}

// CooperativeWithdraw is a paid mutator transaction binding the contract method 0xd757abd2.
//
// Solidity: function cooperativeWithdraw(bytes _cooperativeWithdrawRequest) returns()
func (_CelerLedger *CelerLedgerTransactor) CooperativeWithdraw(opts *bind.TransactOpts, _cooperativeWithdrawRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "cooperativeWithdraw", _cooperativeWithdrawRequest)
}

// CooperativeWithdraw is a paid mutator transaction binding the contract method 0xd757abd2.
//
// Solidity: function cooperativeWithdraw(bytes _cooperativeWithdrawRequest) returns()
func (_CelerLedger *CelerLedgerSession) CooperativeWithdraw(_cooperativeWithdrawRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.CooperativeWithdraw(&_CelerLedger.TransactOpts, _cooperativeWithdrawRequest)
}

// CooperativeWithdraw is a paid mutator transaction binding the contract method 0xd757abd2.
//
// Solidity: function cooperativeWithdraw(bytes _cooperativeWithdrawRequest) returns()
func (_CelerLedger *CelerLedgerTransactorSession) CooperativeWithdraw(_cooperativeWithdrawRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.CooperativeWithdraw(&_CelerLedger.TransactOpts, _cooperativeWithdrawRequest)
}

// Deposit is a paid mutator transaction binding the contract method 0xd954863c.
//
// Solidity: function deposit(bytes32 _channelId, address _receiver, uint256 _transferFromAmount) returns()
func (_CelerLedger *CelerLedgerTransactor) Deposit(opts *bind.TransactOpts, _channelId [32]byte, _receiver common.Address, _transferFromAmount *big.Int) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "deposit", _channelId, _receiver, _transferFromAmount)
}

// Deposit is a paid mutator transaction binding the contract method 0xd954863c.
//
// Solidity: function deposit(bytes32 _channelId, address _receiver, uint256 _transferFromAmount) returns()
func (_CelerLedger *CelerLedgerSession) Deposit(_channelId [32]byte, _receiver common.Address, _transferFromAmount *big.Int) (*types.Transaction, error) {
	return _CelerLedger.Contract.Deposit(&_CelerLedger.TransactOpts, _channelId, _receiver, _transferFromAmount)
}

// Deposit is a paid mutator transaction binding the contract method 0xd954863c.
//
// Solidity: function deposit(bytes32 _channelId, address _receiver, uint256 _transferFromAmount) returns()
func (_CelerLedger *CelerLedgerTransactorSession) Deposit(_channelId [32]byte, _receiver common.Address, _transferFromAmount *big.Int) (*types.Transaction, error) {
	return _CelerLedger.Contract.Deposit(&_CelerLedger.TransactOpts, _channelId, _receiver, _transferFromAmount)
}

// DepositInBatch is a paid mutator transaction binding the contract method 0xe5780db2.
//
// Solidity: function depositInBatch(bytes32[] _channelIds, address[] _receivers, uint256[] _transferFromAmounts) returns()
func (_CelerLedger *CelerLedgerTransactor) DepositInBatch(opts *bind.TransactOpts, _channelIds [][32]byte, _receivers []common.Address, _transferFromAmounts []*big.Int) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "depositInBatch", _channelIds, _receivers, _transferFromAmounts)
}

// DepositInBatch is a paid mutator transaction binding the contract method 0xe5780db2.
//
// Solidity: function depositInBatch(bytes32[] _channelIds, address[] _receivers, uint256[] _transferFromAmounts) returns()
func (_CelerLedger *CelerLedgerSession) DepositInBatch(_channelIds [][32]byte, _receivers []common.Address, _transferFromAmounts []*big.Int) (*types.Transaction, error) {
	return _CelerLedger.Contract.DepositInBatch(&_CelerLedger.TransactOpts, _channelIds, _receivers, _transferFromAmounts)
}

// DepositInBatch is a paid mutator transaction binding the contract method 0xe5780db2.
//
// Solidity: function depositInBatch(bytes32[] _channelIds, address[] _receivers, uint256[] _transferFromAmounts) returns()
func (_CelerLedger *CelerLedgerTransactorSession) DepositInBatch(_channelIds [][32]byte, _receivers []common.Address, _transferFromAmounts []*big.Int) (*types.Transaction, error) {
	return _CelerLedger.Contract.DepositInBatch(&_CelerLedger.TransactOpts, _channelIds, _receivers, _transferFromAmounts)
}

// DisableBalanceLimits is a paid mutator transaction binding the contract method 0xe063913c.
//
// Solidity: function disableBalanceLimits() returns()
func (_CelerLedger *CelerLedgerTransactor) DisableBalanceLimits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "disableBalanceLimits")
}

// DisableBalanceLimits is a paid mutator transaction binding the contract method 0xe063913c.
//
// Solidity: function disableBalanceLimits() returns()
func (_CelerLedger *CelerLedgerSession) DisableBalanceLimits() (*types.Transaction, error) {
	return _CelerLedger.Contract.DisableBalanceLimits(&_CelerLedger.TransactOpts)
}

// DisableBalanceLimits is a paid mutator transaction binding the contract method 0xe063913c.
//
// Solidity: function disableBalanceLimits() returns()
func (_CelerLedger *CelerLedgerTransactorSession) DisableBalanceLimits() (*types.Transaction, error) {
	return _CelerLedger.Contract.DisableBalanceLimits(&_CelerLedger.TransactOpts)
}

// EnableBalanceLimits is a paid mutator transaction binding the contract method 0x83c8f8b8.
//
// Solidity: function enableBalanceLimits() returns()
func (_CelerLedger *CelerLedgerTransactor) EnableBalanceLimits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "enableBalanceLimits")
}

// EnableBalanceLimits is a paid mutator transaction binding the contract method 0x83c8f8b8.
//
// Solidity: function enableBalanceLimits() returns()
func (_CelerLedger *CelerLedgerSession) EnableBalanceLimits() (*types.Transaction, error) {
	return _CelerLedger.Contract.EnableBalanceLimits(&_CelerLedger.TransactOpts)
}

// EnableBalanceLimits is a paid mutator transaction binding the contract method 0x83c8f8b8.
//
// Solidity: function enableBalanceLimits() returns()
func (_CelerLedger *CelerLedgerTransactorSession) EnableBalanceLimits() (*types.Transaction, error) {
	return _CelerLedger.Contract.EnableBalanceLimits(&_CelerLedger.TransactOpts)
}

// IntendSettle is a paid mutator transaction binding the contract method 0x130d33fe.
//
// Solidity: function intendSettle(bytes _signedSimplexStateArray) returns()
func (_CelerLedger *CelerLedgerTransactor) IntendSettle(opts *bind.TransactOpts, _signedSimplexStateArray []byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "intendSettle", _signedSimplexStateArray)
}

// IntendSettle is a paid mutator transaction binding the contract method 0x130d33fe.
//
// Solidity: function intendSettle(bytes _signedSimplexStateArray) returns()
func (_CelerLedger *CelerLedgerSession) IntendSettle(_signedSimplexStateArray []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.IntendSettle(&_CelerLedger.TransactOpts, _signedSimplexStateArray)
}

// IntendSettle is a paid mutator transaction binding the contract method 0x130d33fe.
//
// Solidity: function intendSettle(bytes _signedSimplexStateArray) returns()
func (_CelerLedger *CelerLedgerTransactorSession) IntendSettle(_signedSimplexStateArray []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.IntendSettle(&_CelerLedger.TransactOpts, _signedSimplexStateArray)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x8942ecb2.
//
// Solidity: function intendWithdraw(bytes32 _channelId, uint256 _amount, bytes32 _recipientChannelId) returns()
func (_CelerLedger *CelerLedgerTransactor) IntendWithdraw(opts *bind.TransactOpts, _channelId [32]byte, _amount *big.Int, _recipientChannelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "intendWithdraw", _channelId, _amount, _recipientChannelId)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x8942ecb2.
//
// Solidity: function intendWithdraw(bytes32 _channelId, uint256 _amount, bytes32 _recipientChannelId) returns()
func (_CelerLedger *CelerLedgerSession) IntendWithdraw(_channelId [32]byte, _amount *big.Int, _recipientChannelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.IntendWithdraw(&_CelerLedger.TransactOpts, _channelId, _amount, _recipientChannelId)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x8942ecb2.
//
// Solidity: function intendWithdraw(bytes32 _channelId, uint256 _amount, bytes32 _recipientChannelId) returns()
func (_CelerLedger *CelerLedgerTransactorSession) IntendWithdraw(_channelId [32]byte, _amount *big.Int, _recipientChannelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.IntendWithdraw(&_CelerLedger.TransactOpts, _channelId, _amount, _recipientChannelId)
}

// MigrateChannelFrom is a paid mutator transaction binding the contract method 0x2e3c517a.
//
// Solidity: function migrateChannelFrom(address _fromLedgerAddr, bytes _migrationRequest) returns()
func (_CelerLedger *CelerLedgerTransactor) MigrateChannelFrom(opts *bind.TransactOpts, _fromLedgerAddr common.Address, _migrationRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "migrateChannelFrom", _fromLedgerAddr, _migrationRequest)
}

// MigrateChannelFrom is a paid mutator transaction binding the contract method 0x2e3c517a.
//
// Solidity: function migrateChannelFrom(address _fromLedgerAddr, bytes _migrationRequest) returns()
func (_CelerLedger *CelerLedgerSession) MigrateChannelFrom(_fromLedgerAddr common.Address, _migrationRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.MigrateChannelFrom(&_CelerLedger.TransactOpts, _fromLedgerAddr, _migrationRequest)
}

// MigrateChannelFrom is a paid mutator transaction binding the contract method 0x2e3c517a.
//
// Solidity: function migrateChannelFrom(address _fromLedgerAddr, bytes _migrationRequest) returns()
func (_CelerLedger *CelerLedgerTransactorSession) MigrateChannelFrom(_fromLedgerAddr common.Address, _migrationRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.MigrateChannelFrom(&_CelerLedger.TransactOpts, _fromLedgerAddr, _migrationRequest)
}

// MigrateChannelTo is a paid mutator transaction binding the contract method 0xe0a515b7.
//
// Solidity: function migrateChannelTo(bytes _migrationRequest) returns(bytes32)
func (_CelerLedger *CelerLedgerTransactor) MigrateChannelTo(opts *bind.TransactOpts, _migrationRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "migrateChannelTo", _migrationRequest)
}

// MigrateChannelTo is a paid mutator transaction binding the contract method 0xe0a515b7.
//
// Solidity: function migrateChannelTo(bytes _migrationRequest) returns(bytes32)
func (_CelerLedger *CelerLedgerSession) MigrateChannelTo(_migrationRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.MigrateChannelTo(&_CelerLedger.TransactOpts, _migrationRequest)
}

// MigrateChannelTo is a paid mutator transaction binding the contract method 0xe0a515b7.
//
// Solidity: function migrateChannelTo(bytes _migrationRequest) returns(bytes32)
func (_CelerLedger *CelerLedgerTransactorSession) MigrateChannelTo(_migrationRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.MigrateChannelTo(&_CelerLedger.TransactOpts, _migrationRequest)
}

// OpenChannel is a paid mutator transaction binding the contract method 0x93b7b3ce.
//
// Solidity: function openChannel(bytes _openRequest) returns()
func (_CelerLedger *CelerLedgerTransactor) OpenChannel(opts *bind.TransactOpts, _openRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "openChannel", _openRequest)
}

// OpenChannel is a paid mutator transaction binding the contract method 0x93b7b3ce.
//
// Solidity: function openChannel(bytes _openRequest) returns()
func (_CelerLedger *CelerLedgerSession) OpenChannel(_openRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.OpenChannel(&_CelerLedger.TransactOpts, _openRequest)
}

// OpenChannel is a paid mutator transaction binding the contract method 0x93b7b3ce.
//
// Solidity: function openChannel(bytes _openRequest) returns()
func (_CelerLedger *CelerLedgerTransactorSession) OpenChannel(_openRequest []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.OpenChannel(&_CelerLedger.TransactOpts, _openRequest)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CelerLedger *CelerLedgerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CelerLedger *CelerLedgerSession) RenounceOwnership() (*types.Transaction, error) {
	return _CelerLedger.Contract.RenounceOwnership(&_CelerLedger.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CelerLedger *CelerLedgerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CelerLedger.Contract.RenounceOwnership(&_CelerLedger.TransactOpts)
}

// SetBalanceLimits is a paid mutator transaction binding the contract method 0xa8580cab.
//
// Solidity: function setBalanceLimits(address[] _tokenAddrs, uint256[] _limits) returns()
func (_CelerLedger *CelerLedgerTransactor) SetBalanceLimits(opts *bind.TransactOpts, _tokenAddrs []common.Address, _limits []*big.Int) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "setBalanceLimits", _tokenAddrs, _limits)
}

// SetBalanceLimits is a paid mutator transaction binding the contract method 0xa8580cab.
//
// Solidity: function setBalanceLimits(address[] _tokenAddrs, uint256[] _limits) returns()
func (_CelerLedger *CelerLedgerSession) SetBalanceLimits(_tokenAddrs []common.Address, _limits []*big.Int) (*types.Transaction, error) {
	return _CelerLedger.Contract.SetBalanceLimits(&_CelerLedger.TransactOpts, _tokenAddrs, _limits)
}

// SetBalanceLimits is a paid mutator transaction binding the contract method 0xa8580cab.
//
// Solidity: function setBalanceLimits(address[] _tokenAddrs, uint256[] _limits) returns()
func (_CelerLedger *CelerLedgerTransactorSession) SetBalanceLimits(_tokenAddrs []common.Address, _limits []*big.Int) (*types.Transaction, error) {
	return _CelerLedger.Contract.SetBalanceLimits(&_CelerLedger.TransactOpts, _tokenAddrs, _limits)
}

// SnapshotStates is a paid mutator transaction binding the contract method 0x4102b9a8.
//
// Solidity: function snapshotStates(bytes _signedSimplexStateArray) returns()
func (_CelerLedger *CelerLedgerTransactor) SnapshotStates(opts *bind.TransactOpts, _signedSimplexStateArray []byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "snapshotStates", _signedSimplexStateArray)
}

// SnapshotStates is a paid mutator transaction binding the contract method 0x4102b9a8.
//
// Solidity: function snapshotStates(bytes _signedSimplexStateArray) returns()
func (_CelerLedger *CelerLedgerSession) SnapshotStates(_signedSimplexStateArray []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.SnapshotStates(&_CelerLedger.TransactOpts, _signedSimplexStateArray)
}

// SnapshotStates is a paid mutator transaction binding the contract method 0x4102b9a8.
//
// Solidity: function snapshotStates(bytes _signedSimplexStateArray) returns()
func (_CelerLedger *CelerLedgerTransactorSession) SnapshotStates(_signedSimplexStateArray []byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.SnapshotStates(&_CelerLedger.TransactOpts, _signedSimplexStateArray)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CelerLedger *CelerLedgerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CelerLedger *CelerLedgerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CelerLedger.Contract.TransferOwnership(&_CelerLedger.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CelerLedger *CelerLedgerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CelerLedger.Contract.TransferOwnership(&_CelerLedger.TransactOpts, newOwner)
}

// VetoWithdraw is a paid mutator transaction binding the contract method 0x255aab59.
//
// Solidity: function vetoWithdraw(bytes32 _channelId) returns()
func (_CelerLedger *CelerLedgerTransactor) VetoWithdraw(opts *bind.TransactOpts, _channelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.contract.Transact(opts, "vetoWithdraw", _channelId)
}

// VetoWithdraw is a paid mutator transaction binding the contract method 0x255aab59.
//
// Solidity: function vetoWithdraw(bytes32 _channelId) returns()
func (_CelerLedger *CelerLedgerSession) VetoWithdraw(_channelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.VetoWithdraw(&_CelerLedger.TransactOpts, _channelId)
}

// VetoWithdraw is a paid mutator transaction binding the contract method 0x255aab59.
//
// Solidity: function vetoWithdraw(bytes32 _channelId) returns()
func (_CelerLedger *CelerLedgerTransactorSession) VetoWithdraw(_channelId [32]byte) (*types.Transaction, error) {
	return _CelerLedger.Contract.VetoWithdraw(&_CelerLedger.TransactOpts, _channelId)
}

// CelerLedgerClearOnePayIterator is returned from FilterClearOnePay and is used to iterate over the raw logs and unpacked data for ClearOnePay events raised by the CelerLedger contract.
type CelerLedgerClearOnePayIterator struct {
	Event *CelerLedgerClearOnePay // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerClearOnePayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerClearOnePay)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerClearOnePay)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerClearOnePayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerClearOnePayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerClearOnePay represents a ClearOnePay event raised by the CelerLedger contract.
type CelerLedgerClearOnePay struct {
	ChannelId [32]byte
	PayId     [32]byte
	PeerFrom  common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClearOnePay is a free log retrieval operation binding the contract event 0x33252d4bc5cee2ad248475e8c39239a79dc64b2691c9ca1a63ff9af0c75b8776.
//
// Solidity: event ClearOnePay(bytes32 indexed channelId, bytes32 indexed payId, address indexed peerFrom, uint256 amount)
func (_CelerLedger *CelerLedgerFilterer) FilterClearOnePay(opts *bind.FilterOpts, channelId [][32]byte, payId [][32]byte, peerFrom []common.Address) (*CelerLedgerClearOnePayIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var payIdRule []interface{}
	for _, payIdItem := range payId {
		payIdRule = append(payIdRule, payIdItem)
	}
	var peerFromRule []interface{}
	for _, peerFromItem := range peerFrom {
		peerFromRule = append(peerFromRule, peerFromItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "ClearOnePay", channelIdRule, payIdRule, peerFromRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerClearOnePayIterator{contract: _CelerLedger.contract, event: "ClearOnePay", logs: logs, sub: sub}, nil
}

// WatchClearOnePay is a free log subscription operation binding the contract event 0x33252d4bc5cee2ad248475e8c39239a79dc64b2691c9ca1a63ff9af0c75b8776.
//
// Solidity: event ClearOnePay(bytes32 indexed channelId, bytes32 indexed payId, address indexed peerFrom, uint256 amount)
func (_CelerLedger *CelerLedgerFilterer) WatchClearOnePay(opts *bind.WatchOpts, sink chan<- *CelerLedgerClearOnePay, channelId [][32]byte, payId [][32]byte, peerFrom []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var payIdRule []interface{}
	for _, payIdItem := range payId {
		payIdRule = append(payIdRule, payIdItem)
	}
	var peerFromRule []interface{}
	for _, peerFromItem := range peerFrom {
		peerFromRule = append(peerFromRule, peerFromItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "ClearOnePay", channelIdRule, payIdRule, peerFromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerClearOnePay)
				if err := _CelerLedger.contract.UnpackLog(event, "ClearOnePay", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClearOnePay is a log parse operation binding the contract event 0x33252d4bc5cee2ad248475e8c39239a79dc64b2691c9ca1a63ff9af0c75b8776.
//
// Solidity: event ClearOnePay(bytes32 indexed channelId, bytes32 indexed payId, address indexed peerFrom, uint256 amount)
func (_CelerLedger *CelerLedgerFilterer) ParseClearOnePay(log types.Log) (*CelerLedgerClearOnePay, error) {
	event := new(CelerLedgerClearOnePay)
	if err := _CelerLedger.contract.UnpackLog(event, "ClearOnePay", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerConfirmSettleIterator is returned from FilterConfirmSettle and is used to iterate over the raw logs and unpacked data for ConfirmSettle events raised by the CelerLedger contract.
type CelerLedgerConfirmSettleIterator struct {
	Event *CelerLedgerConfirmSettle // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerConfirmSettleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerConfirmSettle)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerConfirmSettle)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerConfirmSettleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerConfirmSettleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerConfirmSettle represents a ConfirmSettle event raised by the CelerLedger contract.
type CelerLedgerConfirmSettle struct {
	ChannelId     [32]byte
	SettleBalance [2]*big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterConfirmSettle is a free log retrieval operation binding the contract event 0x728ddd8c5acda5947c34db8d759c66ae70884f526ff9b93637d351b012ef3206.
//
// Solidity: event ConfirmSettle(bytes32 indexed channelId, uint256[2] settleBalance)
func (_CelerLedger *CelerLedgerFilterer) FilterConfirmSettle(opts *bind.FilterOpts, channelId [][32]byte) (*CelerLedgerConfirmSettleIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "ConfirmSettle", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerConfirmSettleIterator{contract: _CelerLedger.contract, event: "ConfirmSettle", logs: logs, sub: sub}, nil
}

// WatchConfirmSettle is a free log subscription operation binding the contract event 0x728ddd8c5acda5947c34db8d759c66ae70884f526ff9b93637d351b012ef3206.
//
// Solidity: event ConfirmSettle(bytes32 indexed channelId, uint256[2] settleBalance)
func (_CelerLedger *CelerLedgerFilterer) WatchConfirmSettle(opts *bind.WatchOpts, sink chan<- *CelerLedgerConfirmSettle, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "ConfirmSettle", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerConfirmSettle)
				if err := _CelerLedger.contract.UnpackLog(event, "ConfirmSettle", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfirmSettle is a log parse operation binding the contract event 0x728ddd8c5acda5947c34db8d759c66ae70884f526ff9b93637d351b012ef3206.
//
// Solidity: event ConfirmSettle(bytes32 indexed channelId, uint256[2] settleBalance)
func (_CelerLedger *CelerLedgerFilterer) ParseConfirmSettle(log types.Log) (*CelerLedgerConfirmSettle, error) {
	event := new(CelerLedgerConfirmSettle)
	if err := _CelerLedger.contract.UnpackLog(event, "ConfirmSettle", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerConfirmSettleFailIterator is returned from FilterConfirmSettleFail and is used to iterate over the raw logs and unpacked data for ConfirmSettleFail events raised by the CelerLedger contract.
type CelerLedgerConfirmSettleFailIterator struct {
	Event *CelerLedgerConfirmSettleFail // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerConfirmSettleFailIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerConfirmSettleFail)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerConfirmSettleFail)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerConfirmSettleFailIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerConfirmSettleFailIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerConfirmSettleFail represents a ConfirmSettleFail event raised by the CelerLedger contract.
type CelerLedgerConfirmSettleFail struct {
	ChannelId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterConfirmSettleFail is a free log retrieval operation binding the contract event 0xa6549eb18490d42e7ec93f42115d1ee11b706d04077be9597034dd73ec8bcb36.
//
// Solidity: event ConfirmSettleFail(bytes32 indexed channelId)
func (_CelerLedger *CelerLedgerFilterer) FilterConfirmSettleFail(opts *bind.FilterOpts, channelId [][32]byte) (*CelerLedgerConfirmSettleFailIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "ConfirmSettleFail", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerConfirmSettleFailIterator{contract: _CelerLedger.contract, event: "ConfirmSettleFail", logs: logs, sub: sub}, nil
}

// WatchConfirmSettleFail is a free log subscription operation binding the contract event 0xa6549eb18490d42e7ec93f42115d1ee11b706d04077be9597034dd73ec8bcb36.
//
// Solidity: event ConfirmSettleFail(bytes32 indexed channelId)
func (_CelerLedger *CelerLedgerFilterer) WatchConfirmSettleFail(opts *bind.WatchOpts, sink chan<- *CelerLedgerConfirmSettleFail, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "ConfirmSettleFail", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerConfirmSettleFail)
				if err := _CelerLedger.contract.UnpackLog(event, "ConfirmSettleFail", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfirmSettleFail is a log parse operation binding the contract event 0xa6549eb18490d42e7ec93f42115d1ee11b706d04077be9597034dd73ec8bcb36.
//
// Solidity: event ConfirmSettleFail(bytes32 indexed channelId)
func (_CelerLedger *CelerLedgerFilterer) ParseConfirmSettleFail(log types.Log) (*CelerLedgerConfirmSettleFail, error) {
	event := new(CelerLedgerConfirmSettleFail)
	if err := _CelerLedger.contract.UnpackLog(event, "ConfirmSettleFail", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerConfirmWithdrawIterator is returned from FilterConfirmWithdraw and is used to iterate over the raw logs and unpacked data for ConfirmWithdraw events raised by the CelerLedger contract.
type CelerLedgerConfirmWithdrawIterator struct {
	Event *CelerLedgerConfirmWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerConfirmWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerConfirmWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerConfirmWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerConfirmWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerConfirmWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerConfirmWithdraw represents a ConfirmWithdraw event raised by the CelerLedger contract.
type CelerLedgerConfirmWithdraw struct {
	ChannelId          [32]byte
	WithdrawnAmount    *big.Int
	Receiver           common.Address
	RecipientChannelId [32]byte
	Deposits           [2]*big.Int
	Withdrawals        [2]*big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterConfirmWithdraw is a free log retrieval operation binding the contract event 0xe8110b4ee08638c48f6a4d5f726927df4e541893efa9d2c2c47a6b889041826e.
//
// Solidity: event ConfirmWithdraw(bytes32 indexed channelId, uint256 withdrawnAmount, address indexed receiver, bytes32 indexed recipientChannelId, uint256[2] deposits, uint256[2] withdrawals)
func (_CelerLedger *CelerLedgerFilterer) FilterConfirmWithdraw(opts *bind.FilterOpts, channelId [][32]byte, receiver []common.Address, recipientChannelId [][32]byte) (*CelerLedgerConfirmWithdrawIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var recipientChannelIdRule []interface{}
	for _, recipientChannelIdItem := range recipientChannelId {
		recipientChannelIdRule = append(recipientChannelIdRule, recipientChannelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "ConfirmWithdraw", channelIdRule, receiverRule, recipientChannelIdRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerConfirmWithdrawIterator{contract: _CelerLedger.contract, event: "ConfirmWithdraw", logs: logs, sub: sub}, nil
}

// WatchConfirmWithdraw is a free log subscription operation binding the contract event 0xe8110b4ee08638c48f6a4d5f726927df4e541893efa9d2c2c47a6b889041826e.
//
// Solidity: event ConfirmWithdraw(bytes32 indexed channelId, uint256 withdrawnAmount, address indexed receiver, bytes32 indexed recipientChannelId, uint256[2] deposits, uint256[2] withdrawals)
func (_CelerLedger *CelerLedgerFilterer) WatchConfirmWithdraw(opts *bind.WatchOpts, sink chan<- *CelerLedgerConfirmWithdraw, channelId [][32]byte, receiver []common.Address, recipientChannelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var recipientChannelIdRule []interface{}
	for _, recipientChannelIdItem := range recipientChannelId {
		recipientChannelIdRule = append(recipientChannelIdRule, recipientChannelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "ConfirmWithdraw", channelIdRule, receiverRule, recipientChannelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerConfirmWithdraw)
				if err := _CelerLedger.contract.UnpackLog(event, "ConfirmWithdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfirmWithdraw is a log parse operation binding the contract event 0xe8110b4ee08638c48f6a4d5f726927df4e541893efa9d2c2c47a6b889041826e.
//
// Solidity: event ConfirmWithdraw(bytes32 indexed channelId, uint256 withdrawnAmount, address indexed receiver, bytes32 indexed recipientChannelId, uint256[2] deposits, uint256[2] withdrawals)
func (_CelerLedger *CelerLedgerFilterer) ParseConfirmWithdraw(log types.Log) (*CelerLedgerConfirmWithdraw, error) {
	event := new(CelerLedgerConfirmWithdraw)
	if err := _CelerLedger.contract.UnpackLog(event, "ConfirmWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerCooperativeSettleIterator is returned from FilterCooperativeSettle and is used to iterate over the raw logs and unpacked data for CooperativeSettle events raised by the CelerLedger contract.
type CelerLedgerCooperativeSettleIterator struct {
	Event *CelerLedgerCooperativeSettle // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerCooperativeSettleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerCooperativeSettle)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerCooperativeSettle)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerCooperativeSettleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerCooperativeSettleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerCooperativeSettle represents a CooperativeSettle event raised by the CelerLedger contract.
type CelerLedgerCooperativeSettle struct {
	ChannelId     [32]byte
	SettleBalance [2]*big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterCooperativeSettle is a free log retrieval operation binding the contract event 0x6c666557dc97fd52cd2d9d6dd6d109e501ffdb831abeecf13aafeeaf762ee1fd.
//
// Solidity: event CooperativeSettle(bytes32 indexed channelId, uint256[2] settleBalance)
func (_CelerLedger *CelerLedgerFilterer) FilterCooperativeSettle(opts *bind.FilterOpts, channelId [][32]byte) (*CelerLedgerCooperativeSettleIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "CooperativeSettle", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerCooperativeSettleIterator{contract: _CelerLedger.contract, event: "CooperativeSettle", logs: logs, sub: sub}, nil
}

// WatchCooperativeSettle is a free log subscription operation binding the contract event 0x6c666557dc97fd52cd2d9d6dd6d109e501ffdb831abeecf13aafeeaf762ee1fd.
//
// Solidity: event CooperativeSettle(bytes32 indexed channelId, uint256[2] settleBalance)
func (_CelerLedger *CelerLedgerFilterer) WatchCooperativeSettle(opts *bind.WatchOpts, sink chan<- *CelerLedgerCooperativeSettle, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "CooperativeSettle", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerCooperativeSettle)
				if err := _CelerLedger.contract.UnpackLog(event, "CooperativeSettle", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCooperativeSettle is a log parse operation binding the contract event 0x6c666557dc97fd52cd2d9d6dd6d109e501ffdb831abeecf13aafeeaf762ee1fd.
//
// Solidity: event CooperativeSettle(bytes32 indexed channelId, uint256[2] settleBalance)
func (_CelerLedger *CelerLedgerFilterer) ParseCooperativeSettle(log types.Log) (*CelerLedgerCooperativeSettle, error) {
	event := new(CelerLedgerCooperativeSettle)
	if err := _CelerLedger.contract.UnpackLog(event, "CooperativeSettle", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerCooperativeWithdrawIterator is returned from FilterCooperativeWithdraw and is used to iterate over the raw logs and unpacked data for CooperativeWithdraw events raised by the CelerLedger contract.
type CelerLedgerCooperativeWithdrawIterator struct {
	Event *CelerLedgerCooperativeWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerCooperativeWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerCooperativeWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerCooperativeWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerCooperativeWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerCooperativeWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerCooperativeWithdraw represents a CooperativeWithdraw event raised by the CelerLedger contract.
type CelerLedgerCooperativeWithdraw struct {
	ChannelId          [32]byte
	WithdrawnAmount    *big.Int
	Receiver           common.Address
	RecipientChannelId [32]byte
	Deposits           [2]*big.Int
	Withdrawals        [2]*big.Int
	SeqNum             *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterCooperativeWithdraw is a free log retrieval operation binding the contract event 0x1b87d077d9b706e42883b454b67730633fd6b4b29f9a9cf5f57c278c54f51c8f.
//
// Solidity: event CooperativeWithdraw(bytes32 indexed channelId, uint256 withdrawnAmount, address indexed receiver, bytes32 indexed recipientChannelId, uint256[2] deposits, uint256[2] withdrawals, uint256 seqNum)
func (_CelerLedger *CelerLedgerFilterer) FilterCooperativeWithdraw(opts *bind.FilterOpts, channelId [][32]byte, receiver []common.Address, recipientChannelId [][32]byte) (*CelerLedgerCooperativeWithdrawIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var recipientChannelIdRule []interface{}
	for _, recipientChannelIdItem := range recipientChannelId {
		recipientChannelIdRule = append(recipientChannelIdRule, recipientChannelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "CooperativeWithdraw", channelIdRule, receiverRule, recipientChannelIdRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerCooperativeWithdrawIterator{contract: _CelerLedger.contract, event: "CooperativeWithdraw", logs: logs, sub: sub}, nil
}

// WatchCooperativeWithdraw is a free log subscription operation binding the contract event 0x1b87d077d9b706e42883b454b67730633fd6b4b29f9a9cf5f57c278c54f51c8f.
//
// Solidity: event CooperativeWithdraw(bytes32 indexed channelId, uint256 withdrawnAmount, address indexed receiver, bytes32 indexed recipientChannelId, uint256[2] deposits, uint256[2] withdrawals, uint256 seqNum)
func (_CelerLedger *CelerLedgerFilterer) WatchCooperativeWithdraw(opts *bind.WatchOpts, sink chan<- *CelerLedgerCooperativeWithdraw, channelId [][32]byte, receiver []common.Address, recipientChannelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var recipientChannelIdRule []interface{}
	for _, recipientChannelIdItem := range recipientChannelId {
		recipientChannelIdRule = append(recipientChannelIdRule, recipientChannelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "CooperativeWithdraw", channelIdRule, receiverRule, recipientChannelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerCooperativeWithdraw)
				if err := _CelerLedger.contract.UnpackLog(event, "CooperativeWithdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCooperativeWithdraw is a log parse operation binding the contract event 0x1b87d077d9b706e42883b454b67730633fd6b4b29f9a9cf5f57c278c54f51c8f.
//
// Solidity: event CooperativeWithdraw(bytes32 indexed channelId, uint256 withdrawnAmount, address indexed receiver, bytes32 indexed recipientChannelId, uint256[2] deposits, uint256[2] withdrawals, uint256 seqNum)
func (_CelerLedger *CelerLedgerFilterer) ParseCooperativeWithdraw(log types.Log) (*CelerLedgerCooperativeWithdraw, error) {
	event := new(CelerLedgerCooperativeWithdraw)
	if err := _CelerLedger.contract.UnpackLog(event, "CooperativeWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the CelerLedger contract.
type CelerLedgerDepositIterator struct {
	Event *CelerLedgerDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerDeposit represents a Deposit event raised by the CelerLedger contract.
type CelerLedgerDeposit struct {
	ChannelId   [32]byte
	PeerAddrs   [2]common.Address
	Deposits    [2]*big.Int
	Withdrawals [2]*big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xb63f5dc096f516663ffb5ef2b611f0e2acca8617a868c2a3653cba5e3ed0e92c.
//
// Solidity: event Deposit(bytes32 indexed channelId, address[2] peerAddrs, uint256[2] deposits, uint256[2] withdrawals)
func (_CelerLedger *CelerLedgerFilterer) FilterDeposit(opts *bind.FilterOpts, channelId [][32]byte) (*CelerLedgerDepositIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "Deposit", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerDepositIterator{contract: _CelerLedger.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xb63f5dc096f516663ffb5ef2b611f0e2acca8617a868c2a3653cba5e3ed0e92c.
//
// Solidity: event Deposit(bytes32 indexed channelId, address[2] peerAddrs, uint256[2] deposits, uint256[2] withdrawals)
func (_CelerLedger *CelerLedgerFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *CelerLedgerDeposit, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "Deposit", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerDeposit)
				if err := _CelerLedger.contract.UnpackLog(event, "Deposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposit is a log parse operation binding the contract event 0xb63f5dc096f516663ffb5ef2b611f0e2acca8617a868c2a3653cba5e3ed0e92c.
//
// Solidity: event Deposit(bytes32 indexed channelId, address[2] peerAddrs, uint256[2] deposits, uint256[2] withdrawals)
func (_CelerLedger *CelerLedgerFilterer) ParseDeposit(log types.Log) (*CelerLedgerDeposit, error) {
	event := new(CelerLedgerDeposit)
	if err := _CelerLedger.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerIntendSettleIterator is returned from FilterIntendSettle and is used to iterate over the raw logs and unpacked data for IntendSettle events raised by the CelerLedger contract.
type CelerLedgerIntendSettleIterator struct {
	Event *CelerLedgerIntendSettle // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerIntendSettleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerIntendSettle)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerIntendSettle)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerIntendSettleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerIntendSettleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerIntendSettle represents a IntendSettle event raised by the CelerLedger contract.
type CelerLedgerIntendSettle struct {
	ChannelId [32]byte
	SeqNums   [2]*big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterIntendSettle is a free log retrieval operation binding the contract event 0x296143e7e25aa055fbb871702776a67da540876e2be721d5c38ba23c97c90d64.
//
// Solidity: event IntendSettle(bytes32 indexed channelId, uint256[2] seqNums)
func (_CelerLedger *CelerLedgerFilterer) FilterIntendSettle(opts *bind.FilterOpts, channelId [][32]byte) (*CelerLedgerIntendSettleIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "IntendSettle", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerIntendSettleIterator{contract: _CelerLedger.contract, event: "IntendSettle", logs: logs, sub: sub}, nil
}

// WatchIntendSettle is a free log subscription operation binding the contract event 0x296143e7e25aa055fbb871702776a67da540876e2be721d5c38ba23c97c90d64.
//
// Solidity: event IntendSettle(bytes32 indexed channelId, uint256[2] seqNums)
func (_CelerLedger *CelerLedgerFilterer) WatchIntendSettle(opts *bind.WatchOpts, sink chan<- *CelerLedgerIntendSettle, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "IntendSettle", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerIntendSettle)
				if err := _CelerLedger.contract.UnpackLog(event, "IntendSettle", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIntendSettle is a log parse operation binding the contract event 0x296143e7e25aa055fbb871702776a67da540876e2be721d5c38ba23c97c90d64.
//
// Solidity: event IntendSettle(bytes32 indexed channelId, uint256[2] seqNums)
func (_CelerLedger *CelerLedgerFilterer) ParseIntendSettle(log types.Log) (*CelerLedgerIntendSettle, error) {
	event := new(CelerLedgerIntendSettle)
	if err := _CelerLedger.contract.UnpackLog(event, "IntendSettle", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerIntendWithdrawIterator is returned from FilterIntendWithdraw and is used to iterate over the raw logs and unpacked data for IntendWithdraw events raised by the CelerLedger contract.
type CelerLedgerIntendWithdrawIterator struct {
	Event *CelerLedgerIntendWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerIntendWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerIntendWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerIntendWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerIntendWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerIntendWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerIntendWithdraw represents a IntendWithdraw event raised by the CelerLedger contract.
type CelerLedgerIntendWithdraw struct {
	ChannelId [32]byte
	Receiver  common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterIntendWithdraw is a free log retrieval operation binding the contract event 0x97883669625c4ff7f5432b4ca33fe75fb5fee985deb196a967e5758f846170fe.
//
// Solidity: event IntendWithdraw(bytes32 indexed channelId, address indexed receiver, uint256 amount)
func (_CelerLedger *CelerLedgerFilterer) FilterIntendWithdraw(opts *bind.FilterOpts, channelId [][32]byte, receiver []common.Address) (*CelerLedgerIntendWithdrawIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "IntendWithdraw", channelIdRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerIntendWithdrawIterator{contract: _CelerLedger.contract, event: "IntendWithdraw", logs: logs, sub: sub}, nil
}

// WatchIntendWithdraw is a free log subscription operation binding the contract event 0x97883669625c4ff7f5432b4ca33fe75fb5fee985deb196a967e5758f846170fe.
//
// Solidity: event IntendWithdraw(bytes32 indexed channelId, address indexed receiver, uint256 amount)
func (_CelerLedger *CelerLedgerFilterer) WatchIntendWithdraw(opts *bind.WatchOpts, sink chan<- *CelerLedgerIntendWithdraw, channelId [][32]byte, receiver []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "IntendWithdraw", channelIdRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerIntendWithdraw)
				if err := _CelerLedger.contract.UnpackLog(event, "IntendWithdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIntendWithdraw is a log parse operation binding the contract event 0x97883669625c4ff7f5432b4ca33fe75fb5fee985deb196a967e5758f846170fe.
//
// Solidity: event IntendWithdraw(bytes32 indexed channelId, address indexed receiver, uint256 amount)
func (_CelerLedger *CelerLedgerFilterer) ParseIntendWithdraw(log types.Log) (*CelerLedgerIntendWithdraw, error) {
	event := new(CelerLedgerIntendWithdraw)
	if err := _CelerLedger.contract.UnpackLog(event, "IntendWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerMigrateChannelFromIterator is returned from FilterMigrateChannelFrom and is used to iterate over the raw logs and unpacked data for MigrateChannelFrom events raised by the CelerLedger contract.
type CelerLedgerMigrateChannelFromIterator struct {
	Event *CelerLedgerMigrateChannelFrom // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerMigrateChannelFromIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerMigrateChannelFrom)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerMigrateChannelFrom)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerMigrateChannelFromIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerMigrateChannelFromIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerMigrateChannelFrom represents a MigrateChannelFrom event raised by the CelerLedger contract.
type CelerLedgerMigrateChannelFrom struct {
	ChannelId     [32]byte
	OldLedgerAddr common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMigrateChannelFrom is a free log retrieval operation binding the contract event 0x141a72a1d915a7c4205104b6e564cc991aa827c5f2c672a5d6a1da8bef99d6eb.
//
// Solidity: event MigrateChannelFrom(bytes32 indexed channelId, address indexed oldLedgerAddr)
func (_CelerLedger *CelerLedgerFilterer) FilterMigrateChannelFrom(opts *bind.FilterOpts, channelId [][32]byte, oldLedgerAddr []common.Address) (*CelerLedgerMigrateChannelFromIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var oldLedgerAddrRule []interface{}
	for _, oldLedgerAddrItem := range oldLedgerAddr {
		oldLedgerAddrRule = append(oldLedgerAddrRule, oldLedgerAddrItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "MigrateChannelFrom", channelIdRule, oldLedgerAddrRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerMigrateChannelFromIterator{contract: _CelerLedger.contract, event: "MigrateChannelFrom", logs: logs, sub: sub}, nil
}

// WatchMigrateChannelFrom is a free log subscription operation binding the contract event 0x141a72a1d915a7c4205104b6e564cc991aa827c5f2c672a5d6a1da8bef99d6eb.
//
// Solidity: event MigrateChannelFrom(bytes32 indexed channelId, address indexed oldLedgerAddr)
func (_CelerLedger *CelerLedgerFilterer) WatchMigrateChannelFrom(opts *bind.WatchOpts, sink chan<- *CelerLedgerMigrateChannelFrom, channelId [][32]byte, oldLedgerAddr []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var oldLedgerAddrRule []interface{}
	for _, oldLedgerAddrItem := range oldLedgerAddr {
		oldLedgerAddrRule = append(oldLedgerAddrRule, oldLedgerAddrItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "MigrateChannelFrom", channelIdRule, oldLedgerAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerMigrateChannelFrom)
				if err := _CelerLedger.contract.UnpackLog(event, "MigrateChannelFrom", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMigrateChannelFrom is a log parse operation binding the contract event 0x141a72a1d915a7c4205104b6e564cc991aa827c5f2c672a5d6a1da8bef99d6eb.
//
// Solidity: event MigrateChannelFrom(bytes32 indexed channelId, address indexed oldLedgerAddr)
func (_CelerLedger *CelerLedgerFilterer) ParseMigrateChannelFrom(log types.Log) (*CelerLedgerMigrateChannelFrom, error) {
	event := new(CelerLedgerMigrateChannelFrom)
	if err := _CelerLedger.contract.UnpackLog(event, "MigrateChannelFrom", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerMigrateChannelToIterator is returned from FilterMigrateChannelTo and is used to iterate over the raw logs and unpacked data for MigrateChannelTo events raised by the CelerLedger contract.
type CelerLedgerMigrateChannelToIterator struct {
	Event *CelerLedgerMigrateChannelTo // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerMigrateChannelToIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerMigrateChannelTo)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerMigrateChannelTo)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerMigrateChannelToIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerMigrateChannelToIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerMigrateChannelTo represents a MigrateChannelTo event raised by the CelerLedger contract.
type CelerLedgerMigrateChannelTo struct {
	ChannelId     [32]byte
	NewLedgerAddr common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMigrateChannelTo is a free log retrieval operation binding the contract event 0xdefb8a94bbfc44ef5297b035407a7dd1314f369e39c3301f5b90f8810fb9fe4f.
//
// Solidity: event MigrateChannelTo(bytes32 indexed channelId, address indexed newLedgerAddr)
func (_CelerLedger *CelerLedgerFilterer) FilterMigrateChannelTo(opts *bind.FilterOpts, channelId [][32]byte, newLedgerAddr []common.Address) (*CelerLedgerMigrateChannelToIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var newLedgerAddrRule []interface{}
	for _, newLedgerAddrItem := range newLedgerAddr {
		newLedgerAddrRule = append(newLedgerAddrRule, newLedgerAddrItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "MigrateChannelTo", channelIdRule, newLedgerAddrRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerMigrateChannelToIterator{contract: _CelerLedger.contract, event: "MigrateChannelTo", logs: logs, sub: sub}, nil
}

// WatchMigrateChannelTo is a free log subscription operation binding the contract event 0xdefb8a94bbfc44ef5297b035407a7dd1314f369e39c3301f5b90f8810fb9fe4f.
//
// Solidity: event MigrateChannelTo(bytes32 indexed channelId, address indexed newLedgerAddr)
func (_CelerLedger *CelerLedgerFilterer) WatchMigrateChannelTo(opts *bind.WatchOpts, sink chan<- *CelerLedgerMigrateChannelTo, channelId [][32]byte, newLedgerAddr []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var newLedgerAddrRule []interface{}
	for _, newLedgerAddrItem := range newLedgerAddr {
		newLedgerAddrRule = append(newLedgerAddrRule, newLedgerAddrItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "MigrateChannelTo", channelIdRule, newLedgerAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerMigrateChannelTo)
				if err := _CelerLedger.contract.UnpackLog(event, "MigrateChannelTo", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMigrateChannelTo is a log parse operation binding the contract event 0xdefb8a94bbfc44ef5297b035407a7dd1314f369e39c3301f5b90f8810fb9fe4f.
//
// Solidity: event MigrateChannelTo(bytes32 indexed channelId, address indexed newLedgerAddr)
func (_CelerLedger *CelerLedgerFilterer) ParseMigrateChannelTo(log types.Log) (*CelerLedgerMigrateChannelTo, error) {
	event := new(CelerLedgerMigrateChannelTo)
	if err := _CelerLedger.contract.UnpackLog(event, "MigrateChannelTo", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerOpenChannelIterator is returned from FilterOpenChannel and is used to iterate over the raw logs and unpacked data for OpenChannel events raised by the CelerLedger contract.
type CelerLedgerOpenChannelIterator struct {
	Event *CelerLedgerOpenChannel // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerOpenChannelIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerOpenChannel)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerOpenChannel)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerOpenChannelIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerOpenChannelIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerOpenChannel represents a OpenChannel event raised by the CelerLedger contract.
type CelerLedgerOpenChannel struct {
	ChannelId       [32]byte
	TokenType       *big.Int
	TokenAddress    common.Address
	PeerAddrs       [2]common.Address
	InitialDeposits [2]*big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterOpenChannel is a free log retrieval operation binding the contract event 0x9d9f66221370175606b4085f28a419b201c9b6dafd9e0c4520e5bf69ea3e166d.
//
// Solidity: event OpenChannel(bytes32 indexed channelId, uint256 tokenType, address indexed tokenAddress, address[2] peerAddrs, uint256[2] initialDeposits)
func (_CelerLedger *CelerLedgerFilterer) FilterOpenChannel(opts *bind.FilterOpts, channelId [][32]byte, tokenAddress []common.Address) (*CelerLedgerOpenChannelIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	var tokenAddressRule []interface{}
	for _, tokenAddressItem := range tokenAddress {
		tokenAddressRule = append(tokenAddressRule, tokenAddressItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "OpenChannel", channelIdRule, tokenAddressRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerOpenChannelIterator{contract: _CelerLedger.contract, event: "OpenChannel", logs: logs, sub: sub}, nil
}

// WatchOpenChannel is a free log subscription operation binding the contract event 0x9d9f66221370175606b4085f28a419b201c9b6dafd9e0c4520e5bf69ea3e166d.
//
// Solidity: event OpenChannel(bytes32 indexed channelId, uint256 tokenType, address indexed tokenAddress, address[2] peerAddrs, uint256[2] initialDeposits)
func (_CelerLedger *CelerLedgerFilterer) WatchOpenChannel(opts *bind.WatchOpts, sink chan<- *CelerLedgerOpenChannel, channelId [][32]byte, tokenAddress []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	var tokenAddressRule []interface{}
	for _, tokenAddressItem := range tokenAddress {
		tokenAddressRule = append(tokenAddressRule, tokenAddressItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "OpenChannel", channelIdRule, tokenAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerOpenChannel)
				if err := _CelerLedger.contract.UnpackLog(event, "OpenChannel", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOpenChannel is a log parse operation binding the contract event 0x9d9f66221370175606b4085f28a419b201c9b6dafd9e0c4520e5bf69ea3e166d.
//
// Solidity: event OpenChannel(bytes32 indexed channelId, uint256 tokenType, address indexed tokenAddress, address[2] peerAddrs, uint256[2] initialDeposits)
func (_CelerLedger *CelerLedgerFilterer) ParseOpenChannel(log types.Log) (*CelerLedgerOpenChannel, error) {
	event := new(CelerLedgerOpenChannel)
	if err := _CelerLedger.contract.UnpackLog(event, "OpenChannel", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CelerLedger contract.
type CelerLedgerOwnershipTransferredIterator struct {
	Event *CelerLedgerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerOwnershipTransferred represents a OwnershipTransferred event raised by the CelerLedger contract.
type CelerLedgerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CelerLedger *CelerLedgerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CelerLedgerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerOwnershipTransferredIterator{contract: _CelerLedger.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CelerLedger *CelerLedgerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CelerLedgerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerOwnershipTransferred)
				if err := _CelerLedger.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CelerLedger *CelerLedgerFilterer) ParseOwnershipTransferred(log types.Log) (*CelerLedgerOwnershipTransferred, error) {
	event := new(CelerLedgerOwnershipTransferred)
	if err := _CelerLedger.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerSnapshotStatesIterator is returned from FilterSnapshotStates and is used to iterate over the raw logs and unpacked data for SnapshotStates events raised by the CelerLedger contract.
type CelerLedgerSnapshotStatesIterator struct {
	Event *CelerLedgerSnapshotStates // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerSnapshotStatesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerSnapshotStates)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerSnapshotStates)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerSnapshotStatesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerSnapshotStatesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerSnapshotStates represents a SnapshotStates event raised by the CelerLedger contract.
type CelerLedgerSnapshotStates struct {
	ChannelId [32]byte
	SeqNums   [2]*big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSnapshotStates is a free log retrieval operation binding the contract event 0xd0793cc4198bf052a6d91a9a1273c4af39f02a91b0e19029477511c278c5b271.
//
// Solidity: event SnapshotStates(bytes32 indexed channelId, uint256[2] seqNums)
func (_CelerLedger *CelerLedgerFilterer) FilterSnapshotStates(opts *bind.FilterOpts, channelId [][32]byte) (*CelerLedgerSnapshotStatesIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "SnapshotStates", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerSnapshotStatesIterator{contract: _CelerLedger.contract, event: "SnapshotStates", logs: logs, sub: sub}, nil
}

// WatchSnapshotStates is a free log subscription operation binding the contract event 0xd0793cc4198bf052a6d91a9a1273c4af39f02a91b0e19029477511c278c5b271.
//
// Solidity: event SnapshotStates(bytes32 indexed channelId, uint256[2] seqNums)
func (_CelerLedger *CelerLedgerFilterer) WatchSnapshotStates(opts *bind.WatchOpts, sink chan<- *CelerLedgerSnapshotStates, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "SnapshotStates", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerSnapshotStates)
				if err := _CelerLedger.contract.UnpackLog(event, "SnapshotStates", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSnapshotStates is a log parse operation binding the contract event 0xd0793cc4198bf052a6d91a9a1273c4af39f02a91b0e19029477511c278c5b271.
//
// Solidity: event SnapshotStates(bytes32 indexed channelId, uint256[2] seqNums)
func (_CelerLedger *CelerLedgerFilterer) ParseSnapshotStates(log types.Log) (*CelerLedgerSnapshotStates, error) {
	event := new(CelerLedgerSnapshotStates)
	if err := _CelerLedger.contract.UnpackLog(event, "SnapshotStates", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CelerLedgerVetoWithdrawIterator is returned from FilterVetoWithdraw and is used to iterate over the raw logs and unpacked data for VetoWithdraw events raised by the CelerLedger contract.
type CelerLedgerVetoWithdrawIterator struct {
	Event *CelerLedgerVetoWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CelerLedgerVetoWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerLedgerVetoWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CelerLedgerVetoWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CelerLedgerVetoWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerLedgerVetoWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerLedgerVetoWithdraw represents a VetoWithdraw event raised by the CelerLedger contract.
type CelerLedgerVetoWithdraw struct {
	ChannelId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVetoWithdraw is a free log retrieval operation binding the contract event 0x9a8a5493b616f074b3f754b5fd66049c8e7980f01547289e5e31808485c6002c.
//
// Solidity: event VetoWithdraw(bytes32 indexed channelId)
func (_CelerLedger *CelerLedgerFilterer) FilterVetoWithdraw(opts *bind.FilterOpts, channelId [][32]byte) (*CelerLedgerVetoWithdrawIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.FilterLogs(opts, "VetoWithdraw", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &CelerLedgerVetoWithdrawIterator{contract: _CelerLedger.contract, event: "VetoWithdraw", logs: logs, sub: sub}, nil
}

// WatchVetoWithdraw is a free log subscription operation binding the contract event 0x9a8a5493b616f074b3f754b5fd66049c8e7980f01547289e5e31808485c6002c.
//
// Solidity: event VetoWithdraw(bytes32 indexed channelId)
func (_CelerLedger *CelerLedgerFilterer) WatchVetoWithdraw(opts *bind.WatchOpts, sink chan<- *CelerLedgerVetoWithdraw, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _CelerLedger.contract.WatchLogs(opts, "VetoWithdraw", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerLedgerVetoWithdraw)
				if err := _CelerLedger.contract.UnpackLog(event, "VetoWithdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVetoWithdraw is a log parse operation binding the contract event 0x9a8a5493b616f074b3f754b5fd66049c8e7980f01547289e5e31808485c6002c.
//
// Solidity: event VetoWithdraw(bytes32 indexed channelId)
func (_CelerLedger *CelerLedgerFilterer) ParseVetoWithdraw(log types.Log) (*CelerLedgerVetoWithdraw, error) {
	event := new(CelerLedgerVetoWithdraw)
	if err := _CelerLedger.contract.UnpackLog(event, "VetoWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}
