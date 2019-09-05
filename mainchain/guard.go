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
const GuardABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"feePerBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"subscriptionExpiration\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"sidechainGoLive\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VALIDATOR_SET_MAX_SIZE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawTimeout\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"celerToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validatorSet\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minValidatorNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_celerTokenAddress\",\"type\":\"address\"},{\"name\":\"_feePerBlock\",\"type\":\"uint256\"},{\"name\":\"_withdrawTimeout\",\"type\":\"uint256\"},{\"name\":\"_minValidatorNum\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"minSelfStake\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"sidechainAddr\",\"type\":\"bytes\"}],\"name\":\"InitializeCandidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"newStake\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"totalLockedStake\",\"type\":\"uint256\"}],\"name\":\"Delegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"oldSidechainAddr\",\"type\":\"bytes\"},{\"indexed\":true,\"name\":\"newSidechainAddr\",\"type\":\"bytes\"}],\"name\":\"UpdateSidechainAddr\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"ethAddr\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"changeType\",\"type\":\"uint8\"}],\"name\":\"ValidatorChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"withdrawAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"totalLockedStake\",\"type\":\"uint256\"}],\"name\":\"IntendWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ConfirmWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"indemnitor\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"indemnitee\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Punish\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"consumer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"subscriptionExpiration\",\"type\":\"uint256\"}],\"name\":\"Subscription\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_minSelfStake\",\"type\":\"uint256\"},{\"name\":\"_sidechainAddr\",\"type\":\"bytes\"}],\"name\":\"initializeCandidate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_candidateAddr\",\"type\":\"address\"}],\"name\":\"delegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_sidechainAddr\",\"type\":\"bytes\"}],\"name\":\"updateSidechainAddr\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_candidateAddr\",\"type\":\"address\"}],\"name\":\"intendWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"confirmWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"}],\"name\":\"confirmWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"subscribe\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getValidatorNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMinStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"}],\"name\":\"getCandidateInfo\",\"outputs\":[{\"name\":\"initialized\",\"type\":\"bool\"},{\"name\":\"minSelfStake\",\"type\":\"uint256\"},{\"name\":\"sidechainAddr\",\"type\":\"bytes\"},{\"name\":\"totalLockedStake\",\"type\":\"uint256\"},{\"name\":\"isVldt\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"},{\"name\":\"_delegatorAddr\",\"type\":\"address\"}],\"name\":\"getDelegatorInfo\",\"outputs\":[{\"name\":\"lockedStake\",\"type\":\"uint256\"},{\"name\":\"intentAmounts\",\"type\":\"uint256[]\"},{\"name\":\"intentUnlockTimes\",\"type\":\"uint256[]\"},{\"name\":\"nextWithdrawIntent\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// GuardBin is the compiled bytecode used for deploying new contracts.
const GuardBin = `0x`

// DeployGuard deploys a new Ethereum contract, binding an instance of Guard to it.
func DeployGuard(auth *bind.TransactOpts, backend bind.ContractBackend, _celerTokenAddress common.Address, _feePerBlock *big.Int, _withdrawTimeout *big.Int, _minValidatorNum *big.Int) (common.Address, *types.Transaction, *Guard, error) {
	parsed, err := abi.JSON(strings.NewReader(GuardABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(GuardBin), backend, _celerTokenAddress, _feePerBlock, _withdrawTimeout, _minValidatorNum)
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

// VALIDATORSETMAXSIZE is a free data retrieval call binding the contract method 0x63a20c06.
//
// Solidity: function VALIDATOR_SET_MAX_SIZE() constant returns(uint256)
func (_Guard *GuardCaller) VALIDATORSETMAXSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "VALIDATOR_SET_MAX_SIZE")
	return *ret0, err
}

// VALIDATORSETMAXSIZE is a free data retrieval call binding the contract method 0x63a20c06.
//
// Solidity: function VALIDATOR_SET_MAX_SIZE() constant returns(uint256)
func (_Guard *GuardSession) VALIDATORSETMAXSIZE() (*big.Int, error) {
	return _Guard.Contract.VALIDATORSETMAXSIZE(&_Guard.CallOpts)
}

// VALIDATORSETMAXSIZE is a free data retrieval call binding the contract method 0x63a20c06.
//
// Solidity: function VALIDATOR_SET_MAX_SIZE() constant returns(uint256)
func (_Guard *GuardCallerSession) VALIDATORSETMAXSIZE() (*big.Int, error) {
	return _Guard.Contract.VALIDATORSETMAXSIZE(&_Guard.CallOpts)
}

// CelerToken is a free data retrieval call binding the contract method 0xc6c21e9d.
//
// Solidity: function celerToken() constant returns(address)
func (_Guard *GuardCaller) CelerToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "celerToken")
	return *ret0, err
}

// CelerToken is a free data retrieval call binding the contract method 0xc6c21e9d.
//
// Solidity: function celerToken() constant returns(address)
func (_Guard *GuardSession) CelerToken() (common.Address, error) {
	return _Guard.Contract.CelerToken(&_Guard.CallOpts)
}

// CelerToken is a free data retrieval call binding the contract method 0xc6c21e9d.
//
// Solidity: function celerToken() constant returns(address)
func (_Guard *GuardCallerSession) CelerToken() (common.Address, error) {
	return _Guard.Contract.CelerToken(&_Guard.CallOpts)
}

// FeePerBlock is a free data retrieval call binding the contract method 0x4415f47a.
//
// Solidity: function feePerBlock() constant returns(uint256)
func (_Guard *GuardCaller) FeePerBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "feePerBlock")
	return *ret0, err
}

// FeePerBlock is a free data retrieval call binding the contract method 0x4415f47a.
//
// Solidity: function feePerBlock() constant returns(uint256)
func (_Guard *GuardSession) FeePerBlock() (*big.Int, error) {
	return _Guard.Contract.FeePerBlock(&_Guard.CallOpts)
}

// FeePerBlock is a free data retrieval call binding the contract method 0x4415f47a.
//
// Solidity: function feePerBlock() constant returns(uint256)
func (_Guard *GuardCallerSession) FeePerBlock() (*big.Int, error) {
	return _Guard.Contract.FeePerBlock(&_Guard.CallOpts)
}

// GetCandidateInfo is a free data retrieval call binding the contract method 0x28bde1e1.
//
// Solidity: function getCandidateInfo(address _candidateAddr) constant returns(bool initialized, uint256 minSelfStake, bytes sidechainAddr, uint256 totalLockedStake, bool isVldt)
func (_Guard *GuardCaller) GetCandidateInfo(opts *bind.CallOpts, _candidateAddr common.Address) (struct {
	Initialized      bool
	MinSelfStake     *big.Int
	SidechainAddr    []byte
	TotalLockedStake *big.Int
	IsVldt           bool
}, error) {
	ret := new(struct {
		Initialized      bool
		MinSelfStake     *big.Int
		SidechainAddr    []byte
		TotalLockedStake *big.Int
		IsVldt           bool
	})
	out := ret
	err := _Guard.contract.Call(opts, out, "getCandidateInfo", _candidateAddr)
	return *ret, err
}

// GetCandidateInfo is a free data retrieval call binding the contract method 0x28bde1e1.
//
// Solidity: function getCandidateInfo(address _candidateAddr) constant returns(bool initialized, uint256 minSelfStake, bytes sidechainAddr, uint256 totalLockedStake, bool isVldt)
func (_Guard *GuardSession) GetCandidateInfo(_candidateAddr common.Address) (struct {
	Initialized      bool
	MinSelfStake     *big.Int
	SidechainAddr    []byte
	TotalLockedStake *big.Int
	IsVldt           bool
}, error) {
	return _Guard.Contract.GetCandidateInfo(&_Guard.CallOpts, _candidateAddr)
}

// GetCandidateInfo is a free data retrieval call binding the contract method 0x28bde1e1.
//
// Solidity: function getCandidateInfo(address _candidateAddr) constant returns(bool initialized, uint256 minSelfStake, bytes sidechainAddr, uint256 totalLockedStake, bool isVldt)
func (_Guard *GuardCallerSession) GetCandidateInfo(_candidateAddr common.Address) (struct {
	Initialized      bool
	MinSelfStake     *big.Int
	SidechainAddr    []byte
	TotalLockedStake *big.Int
	IsVldt           bool
}, error) {
	return _Guard.Contract.GetCandidateInfo(&_Guard.CallOpts, _candidateAddr)
}

// GetDelegatorInfo is a free data retrieval call binding the contract method 0xeecefef8.
//
// Solidity: function getDelegatorInfo(address _candidateAddr, address _delegatorAddr) constant returns(uint256 lockedStake, uint256[] intentAmounts, uint256[] intentUnlockTimes, uint256 nextWithdrawIntent)
func (_Guard *GuardCaller) GetDelegatorInfo(opts *bind.CallOpts, _candidateAddr common.Address, _delegatorAddr common.Address) (struct {
	LockedStake        *big.Int
	IntentAmounts      []*big.Int
	IntentUnlockTimes  []*big.Int
	NextWithdrawIntent *big.Int
}, error) {
	ret := new(struct {
		LockedStake        *big.Int
		IntentAmounts      []*big.Int
		IntentUnlockTimes  []*big.Int
		NextWithdrawIntent *big.Int
	})
	out := ret
	err := _Guard.contract.Call(opts, out, "getDelegatorInfo", _candidateAddr, _delegatorAddr)
	return *ret, err
}

// GetDelegatorInfo is a free data retrieval call binding the contract method 0xeecefef8.
//
// Solidity: function getDelegatorInfo(address _candidateAddr, address _delegatorAddr) constant returns(uint256 lockedStake, uint256[] intentAmounts, uint256[] intentUnlockTimes, uint256 nextWithdrawIntent)
func (_Guard *GuardSession) GetDelegatorInfo(_candidateAddr common.Address, _delegatorAddr common.Address) (struct {
	LockedStake        *big.Int
	IntentAmounts      []*big.Int
	IntentUnlockTimes  []*big.Int
	NextWithdrawIntent *big.Int
}, error) {
	return _Guard.Contract.GetDelegatorInfo(&_Guard.CallOpts, _candidateAddr, _delegatorAddr)
}

// GetDelegatorInfo is a free data retrieval call binding the contract method 0xeecefef8.
//
// Solidity: function getDelegatorInfo(address _candidateAddr, address _delegatorAddr) constant returns(uint256 lockedStake, uint256[] intentAmounts, uint256[] intentUnlockTimes, uint256 nextWithdrawIntent)
func (_Guard *GuardCallerSession) GetDelegatorInfo(_candidateAddr common.Address, _delegatorAddr common.Address) (struct {
	LockedStake        *big.Int
	IntentAmounts      []*big.Int
	IntentUnlockTimes  []*big.Int
	NextWithdrawIntent *big.Int
}, error) {
	return _Guard.Contract.GetDelegatorInfo(&_Guard.CallOpts, _candidateAddr, _delegatorAddr)
}

// GetMinStake is a free data retrieval call binding the contract method 0x56a3b5fa.
//
// Solidity: function getMinStake() constant returns(uint256)
func (_Guard *GuardCaller) GetMinStake(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "getMinStake")
	return *ret0, err
}

// GetMinStake is a free data retrieval call binding the contract method 0x56a3b5fa.
//
// Solidity: function getMinStake() constant returns(uint256)
func (_Guard *GuardSession) GetMinStake() (*big.Int, error) {
	return _Guard.Contract.GetMinStake(&_Guard.CallOpts)
}

// GetMinStake is a free data retrieval call binding the contract method 0x56a3b5fa.
//
// Solidity: function getMinStake() constant returns(uint256)
func (_Guard *GuardCallerSession) GetMinStake() (*big.Int, error) {
	return _Guard.Contract.GetMinStake(&_Guard.CallOpts)
}

// GetValidatorNum is a free data retrieval call binding the contract method 0x1cfe4f0b.
//
// Solidity: function getValidatorNum() constant returns(uint256)
func (_Guard *GuardCaller) GetValidatorNum(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "getValidatorNum")
	return *ret0, err
}

// GetValidatorNum is a free data retrieval call binding the contract method 0x1cfe4f0b.
//
// Solidity: function getValidatorNum() constant returns(uint256)
func (_Guard *GuardSession) GetValidatorNum() (*big.Int, error) {
	return _Guard.Contract.GetValidatorNum(&_Guard.CallOpts)
}

// GetValidatorNum is a free data retrieval call binding the contract method 0x1cfe4f0b.
//
// Solidity: function getValidatorNum() constant returns(uint256)
func (_Guard *GuardCallerSession) GetValidatorNum() (*big.Int, error) {
	return _Guard.Contract.GetValidatorNum(&_Guard.CallOpts)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address _addr) constant returns(bool)
func (_Guard *GuardCaller) IsValidator(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "isValidator", _addr)
	return *ret0, err
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address _addr) constant returns(bool)
func (_Guard *GuardSession) IsValidator(_addr common.Address) (bool, error) {
	return _Guard.Contract.IsValidator(&_Guard.CallOpts, _addr)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address _addr) constant returns(bool)
func (_Guard *GuardCallerSession) IsValidator(_addr common.Address) (bool, error) {
	return _Guard.Contract.IsValidator(&_Guard.CallOpts, _addr)
}

// MinValidatorNum is a free data retrieval call binding the contract method 0xea5976a9.
//
// Solidity: function minValidatorNum() constant returns(uint256)
func (_Guard *GuardCaller) MinValidatorNum(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "minValidatorNum")
	return *ret0, err
}

// MinValidatorNum is a free data retrieval call binding the contract method 0xea5976a9.
//
// Solidity: function minValidatorNum() constant returns(uint256)
func (_Guard *GuardSession) MinValidatorNum() (*big.Int, error) {
	return _Guard.Contract.MinValidatorNum(&_Guard.CallOpts)
}

// MinValidatorNum is a free data retrieval call binding the contract method 0xea5976a9.
//
// Solidity: function minValidatorNum() constant returns(uint256)
func (_Guard *GuardCallerSession) MinValidatorNum() (*big.Int, error) {
	return _Guard.Contract.MinValidatorNum(&_Guard.CallOpts)
}

// SidechainGoLive is a free data retrieval call binding the contract method 0x55d04500.
//
// Solidity: function sidechainGoLive() constant returns(uint256)
func (_Guard *GuardCaller) SidechainGoLive(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "sidechainGoLive")
	return *ret0, err
}

// SidechainGoLive is a free data retrieval call binding the contract method 0x55d04500.
//
// Solidity: function sidechainGoLive() constant returns(uint256)
func (_Guard *GuardSession) SidechainGoLive() (*big.Int, error) {
	return _Guard.Contract.SidechainGoLive(&_Guard.CallOpts)
}

// SidechainGoLive is a free data retrieval call binding the contract method 0x55d04500.
//
// Solidity: function sidechainGoLive() constant returns(uint256)
func (_Guard *GuardCallerSession) SidechainGoLive() (*big.Int, error) {
	return _Guard.Contract.SidechainGoLive(&_Guard.CallOpts)
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

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) constant returns(address)
func (_Guard *GuardCaller) ValidatorSet(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "validatorSet", arg0)
	return *ret0, err
}

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) constant returns(address)
func (_Guard *GuardSession) ValidatorSet(arg0 *big.Int) (common.Address, error) {
	return _Guard.Contract.ValidatorSet(&_Guard.CallOpts, arg0)
}

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) constant returns(address)
func (_Guard *GuardCallerSession) ValidatorSet(arg0 *big.Int) (common.Address, error) {
	return _Guard.Contract.ValidatorSet(&_Guard.CallOpts, arg0)
}

// WithdrawTimeout is a free data retrieval call binding the contract method 0x9c15d1a2.
//
// Solidity: function withdrawTimeout() constant returns(uint256)
func (_Guard *GuardCaller) WithdrawTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "withdrawTimeout")
	return *ret0, err
}

// WithdrawTimeout is a free data retrieval call binding the contract method 0x9c15d1a2.
//
// Solidity: function withdrawTimeout() constant returns(uint256)
func (_Guard *GuardSession) WithdrawTimeout() (*big.Int, error) {
	return _Guard.Contract.WithdrawTimeout(&_Guard.CallOpts)
}

// WithdrawTimeout is a free data retrieval call binding the contract method 0x9c15d1a2.
//
// Solidity: function withdrawTimeout() constant returns(uint256)
func (_Guard *GuardCallerSession) WithdrawTimeout() (*big.Int, error) {
	return _Guard.Contract.WithdrawTimeout(&_Guard.CallOpts)
}

// ClaimValidator is a paid mutator transaction binding the contract method 0x6e7cf85d.
//
// Solidity: function claimValidator() returns()
func (_Guard *GuardTransactor) ClaimValidator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "claimValidator")
}

// ClaimValidator is a paid mutator transaction binding the contract method 0x6e7cf85d.
//
// Solidity: function claimValidator() returns()
func (_Guard *GuardSession) ClaimValidator() (*types.Transaction, error) {
	return _Guard.Contract.ClaimValidator(&_Guard.TransactOpts)
}

// ClaimValidator is a paid mutator transaction binding the contract method 0x6e7cf85d.
//
// Solidity: function claimValidator() returns()
func (_Guard *GuardTransactorSession) ClaimValidator() (*types.Transaction, error) {
	return _Guard.Contract.ClaimValidator(&_Guard.TransactOpts)
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0xd2bfc1c7.
//
// Solidity: function confirmWithdraw(address _candidateAddr) returns()
func (_Guard *GuardTransactor) ConfirmWithdraw(opts *bind.TransactOpts, _candidateAddr common.Address) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "confirmWithdraw", _candidateAddr)
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0xd2bfc1c7.
//
// Solidity: function confirmWithdraw(address _candidateAddr) returns()
func (_Guard *GuardSession) ConfirmWithdraw(_candidateAddr common.Address) (*types.Transaction, error) {
	return _Guard.Contract.ConfirmWithdraw(&_Guard.TransactOpts, _candidateAddr)
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0xd2bfc1c7.
//
// Solidity: function confirmWithdraw(address _candidateAddr) returns()
func (_Guard *GuardTransactorSession) ConfirmWithdraw(_candidateAddr common.Address) (*types.Transaction, error) {
	return _Guard.Contract.ConfirmWithdraw(&_Guard.TransactOpts, _candidateAddr)
}

// Delegate is a paid mutator transaction binding the contract method 0x08bbb824.
//
// Solidity: function delegate(uint256 _amount, address _candidateAddr) returns()
func (_Guard *GuardTransactor) Delegate(opts *bind.TransactOpts, _amount *big.Int, _candidateAddr common.Address) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "delegate", _amount, _candidateAddr)
}

// Delegate is a paid mutator transaction binding the contract method 0x08bbb824.
//
// Solidity: function delegate(uint256 _amount, address _candidateAddr) returns()
func (_Guard *GuardSession) Delegate(_amount *big.Int, _candidateAddr common.Address) (*types.Transaction, error) {
	return _Guard.Contract.Delegate(&_Guard.TransactOpts, _amount, _candidateAddr)
}

// Delegate is a paid mutator transaction binding the contract method 0x08bbb824.
//
// Solidity: function delegate(uint256 _amount, address _candidateAddr) returns()
func (_Guard *GuardTransactorSession) Delegate(_amount *big.Int, _candidateAddr common.Address) (*types.Transaction, error) {
	return _Guard.Contract.Delegate(&_Guard.TransactOpts, _amount, _candidateAddr)
}

// InitializeCandidate is a paid mutator transaction binding the contract method 0x26c36617.
//
// Solidity: function initializeCandidate(uint256 _minSelfStake, bytes _sidechainAddr) returns()
func (_Guard *GuardTransactor) InitializeCandidate(opts *bind.TransactOpts, _minSelfStake *big.Int, _sidechainAddr []byte) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "initializeCandidate", _minSelfStake, _sidechainAddr)
}

// InitializeCandidate is a paid mutator transaction binding the contract method 0x26c36617.
//
// Solidity: function initializeCandidate(uint256 _minSelfStake, bytes _sidechainAddr) returns()
func (_Guard *GuardSession) InitializeCandidate(_minSelfStake *big.Int, _sidechainAddr []byte) (*types.Transaction, error) {
	return _Guard.Contract.InitializeCandidate(&_Guard.TransactOpts, _minSelfStake, _sidechainAddr)
}

// InitializeCandidate is a paid mutator transaction binding the contract method 0x26c36617.
//
// Solidity: function initializeCandidate(uint256 _minSelfStake, bytes _sidechainAddr) returns()
func (_Guard *GuardTransactorSession) InitializeCandidate(_minSelfStake *big.Int, _sidechainAddr []byte) (*types.Transaction, error) {
	return _Guard.Contract.InitializeCandidate(&_Guard.TransactOpts, _minSelfStake, _sidechainAddr)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x859012a9.
//
// Solidity: function intendWithdraw(uint256 _amount, address _candidateAddr) returns()
func (_Guard *GuardTransactor) IntendWithdraw(opts *bind.TransactOpts, _amount *big.Int, _candidateAddr common.Address) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "intendWithdraw", _amount, _candidateAddr)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x859012a9.
//
// Solidity: function intendWithdraw(uint256 _amount, address _candidateAddr) returns()
func (_Guard *GuardSession) IntendWithdraw(_amount *big.Int, _candidateAddr common.Address) (*types.Transaction, error) {
	return _Guard.Contract.IntendWithdraw(&_Guard.TransactOpts, _amount, _candidateAddr)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x859012a9.
//
// Solidity: function intendWithdraw(uint256 _amount, address _candidateAddr) returns()
func (_Guard *GuardTransactorSession) IntendWithdraw(_amount *big.Int, _candidateAddr common.Address) (*types.Transaction, error) {
	return _Guard.Contract.IntendWithdraw(&_Guard.TransactOpts, _amount, _candidateAddr)
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

// UpdateSidechainAddr is a paid mutator transaction binding the contract method 0xe02f39bd.
//
// Solidity: function updateSidechainAddr(bytes _sidechainAddr) returns()
func (_Guard *GuardTransactor) UpdateSidechainAddr(opts *bind.TransactOpts, _sidechainAddr []byte) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "updateSidechainAddr", _sidechainAddr)
}

// UpdateSidechainAddr is a paid mutator transaction binding the contract method 0xe02f39bd.
//
// Solidity: function updateSidechainAddr(bytes _sidechainAddr) returns()
func (_Guard *GuardSession) UpdateSidechainAddr(_sidechainAddr []byte) (*types.Transaction, error) {
	return _Guard.Contract.UpdateSidechainAddr(&_Guard.TransactOpts, _sidechainAddr)
}

// UpdateSidechainAddr is a paid mutator transaction binding the contract method 0xe02f39bd.
//
// Solidity: function updateSidechainAddr(bytes _sidechainAddr) returns()
func (_Guard *GuardTransactorSession) UpdateSidechainAddr(_sidechainAddr []byte) (*types.Transaction, error) {
	return _Guard.Contract.UpdateSidechainAddr(&_Guard.TransactOpts, _sidechainAddr)
}

// GuardConfirmWithdrawIterator is returned from FilterConfirmWithdraw and is used to iterate over the raw logs and unpacked data for ConfirmWithdraw events raised by the Guard contract.
type GuardConfirmWithdrawIterator struct {
	Event *GuardConfirmWithdraw // Event containing the contract specifics and raw log

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
func (it *GuardConfirmWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardConfirmWithdraw)
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
		it.Event = new(GuardConfirmWithdraw)
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
func (it *GuardConfirmWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardConfirmWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardConfirmWithdraw represents a ConfirmWithdraw event raised by the Guard contract.
type GuardConfirmWithdraw struct {
	Delegator common.Address
	Candidate common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterConfirmWithdraw is a free log retrieval operation binding the contract event 0x08d0283ea9a2e520a2f09611cf37ca6eb70f62e9a807e53756047dd2dc027220.
//
// Solidity: event ConfirmWithdraw(address indexed delegator, address indexed candidate, uint256 amount)
func (_Guard *GuardFilterer) FilterConfirmWithdraw(opts *bind.FilterOpts, delegator []common.Address, candidate []common.Address) (*GuardConfirmWithdrawIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _Guard.contract.FilterLogs(opts, "ConfirmWithdraw", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return &GuardConfirmWithdrawIterator{contract: _Guard.contract, event: "ConfirmWithdraw", logs: logs, sub: sub}, nil
}

// WatchConfirmWithdraw is a free log subscription operation binding the contract event 0x08d0283ea9a2e520a2f09611cf37ca6eb70f62e9a807e53756047dd2dc027220.
//
// Solidity: event ConfirmWithdraw(address indexed delegator, address indexed candidate, uint256 amount)
func (_Guard *GuardFilterer) WatchConfirmWithdraw(opts *bind.WatchOpts, sink chan<- *GuardConfirmWithdraw, delegator []common.Address, candidate []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _Guard.contract.WatchLogs(opts, "ConfirmWithdraw", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardConfirmWithdraw)
				if err := _Guard.contract.UnpackLog(event, "ConfirmWithdraw", log); err != nil {
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

// GuardDelegateIterator is returned from FilterDelegate and is used to iterate over the raw logs and unpacked data for Delegate events raised by the Guard contract.
type GuardDelegateIterator struct {
	Event *GuardDelegate // Event containing the contract specifics and raw log

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
func (it *GuardDelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardDelegate)
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
		it.Event = new(GuardDelegate)
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
func (it *GuardDelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardDelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardDelegate represents a Delegate event raised by the Guard contract.
type GuardDelegate struct {
	Delegator        common.Address
	Candidate        common.Address
	NewStake         *big.Int
	TotalLockedStake *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterDelegate is a free log retrieval operation binding the contract event 0x500599802164a08023e87ffc3eed0ba3ae60697b3083ba81d046683679d81c6b.
//
// Solidity: event Delegate(address indexed delegator, address indexed candidate, uint256 newStake, uint256 totalLockedStake)
func (_Guard *GuardFilterer) FilterDelegate(opts *bind.FilterOpts, delegator []common.Address, candidate []common.Address) (*GuardDelegateIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _Guard.contract.FilterLogs(opts, "Delegate", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return &GuardDelegateIterator{contract: _Guard.contract, event: "Delegate", logs: logs, sub: sub}, nil
}

// WatchDelegate is a free log subscription operation binding the contract event 0x500599802164a08023e87ffc3eed0ba3ae60697b3083ba81d046683679d81c6b.
//
// Solidity: event Delegate(address indexed delegator, address indexed candidate, uint256 newStake, uint256 totalLockedStake)
func (_Guard *GuardFilterer) WatchDelegate(opts *bind.WatchOpts, sink chan<- *GuardDelegate, delegator []common.Address, candidate []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _Guard.contract.WatchLogs(opts, "Delegate", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardDelegate)
				if err := _Guard.contract.UnpackLog(event, "Delegate", log); err != nil {
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

// GuardInitializeCandidateIterator is returned from FilterInitializeCandidate and is used to iterate over the raw logs and unpacked data for InitializeCandidate events raised by the Guard contract.
type GuardInitializeCandidateIterator struct {
	Event *GuardInitializeCandidate // Event containing the contract specifics and raw log

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
func (it *GuardInitializeCandidateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardInitializeCandidate)
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
		it.Event = new(GuardInitializeCandidate)
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
func (it *GuardInitializeCandidateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardInitializeCandidateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardInitializeCandidate represents a InitializeCandidate event raised by the Guard contract.
type GuardInitializeCandidate struct {
	Candidate     common.Address
	MinSelfStake  *big.Int
	SidechainAddr common.Hash
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterInitializeCandidate is a free log retrieval operation binding the contract event 0x377f6597c5132797119197fab0e953f73fd5bb109a897a11c871834af2d092a9.
//
// Solidity: event InitializeCandidate(address indexed candidate, uint256 minSelfStake, bytes indexed sidechainAddr)
func (_Guard *GuardFilterer) FilterInitializeCandidate(opts *bind.FilterOpts, candidate []common.Address, sidechainAddr [][]byte) (*GuardInitializeCandidateIterator, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	var sidechainAddrRule []interface{}
	for _, sidechainAddrItem := range sidechainAddr {
		sidechainAddrRule = append(sidechainAddrRule, sidechainAddrItem)
	}

	logs, sub, err := _Guard.contract.FilterLogs(opts, "InitializeCandidate", candidateRule, sidechainAddrRule)
	if err != nil {
		return nil, err
	}
	return &GuardInitializeCandidateIterator{contract: _Guard.contract, event: "InitializeCandidate", logs: logs, sub: sub}, nil
}

// WatchInitializeCandidate is a free log subscription operation binding the contract event 0x377f6597c5132797119197fab0e953f73fd5bb109a897a11c871834af2d092a9.
//
// Solidity: event InitializeCandidate(address indexed candidate, uint256 minSelfStake, bytes indexed sidechainAddr)
func (_Guard *GuardFilterer) WatchInitializeCandidate(opts *bind.WatchOpts, sink chan<- *GuardInitializeCandidate, candidate []common.Address, sidechainAddr [][]byte) (event.Subscription, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	var sidechainAddrRule []interface{}
	for _, sidechainAddrItem := range sidechainAddr {
		sidechainAddrRule = append(sidechainAddrRule, sidechainAddrItem)
	}

	logs, sub, err := _Guard.contract.WatchLogs(opts, "InitializeCandidate", candidateRule, sidechainAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardInitializeCandidate)
				if err := _Guard.contract.UnpackLog(event, "InitializeCandidate", log); err != nil {
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

// GuardIntendWithdrawIterator is returned from FilterIntendWithdraw and is used to iterate over the raw logs and unpacked data for IntendWithdraw events raised by the Guard contract.
type GuardIntendWithdrawIterator struct {
	Event *GuardIntendWithdraw // Event containing the contract specifics and raw log

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
func (it *GuardIntendWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardIntendWithdraw)
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
		it.Event = new(GuardIntendWithdraw)
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
func (it *GuardIntendWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardIntendWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardIntendWithdraw represents a IntendWithdraw event raised by the Guard contract.
type GuardIntendWithdraw struct {
	Delegator        common.Address
	Candidate        common.Address
	WithdrawAmount   *big.Int
	UnlockTime       *big.Int
	TotalLockedStake *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterIntendWithdraw is a free log retrieval operation binding the contract event 0x9e772df8b63d7657919bf5919c475e4033a0bd817b2468bc7ced0d962f21ded0.
//
// Solidity: event IntendWithdraw(address indexed delegator, address indexed candidate, uint256 withdrawAmount, uint256 unlockTime, uint256 totalLockedStake)
func (_Guard *GuardFilterer) FilterIntendWithdraw(opts *bind.FilterOpts, delegator []common.Address, candidate []common.Address) (*GuardIntendWithdrawIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _Guard.contract.FilterLogs(opts, "IntendWithdraw", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return &GuardIntendWithdrawIterator{contract: _Guard.contract, event: "IntendWithdraw", logs: logs, sub: sub}, nil
}

// WatchIntendWithdraw is a free log subscription operation binding the contract event 0x9e772df8b63d7657919bf5919c475e4033a0bd817b2468bc7ced0d962f21ded0.
//
// Solidity: event IntendWithdraw(address indexed delegator, address indexed candidate, uint256 withdrawAmount, uint256 unlockTime, uint256 totalLockedStake)
func (_Guard *GuardFilterer) WatchIntendWithdraw(opts *bind.WatchOpts, sink chan<- *GuardIntendWithdraw, delegator []common.Address, candidate []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _Guard.contract.WatchLogs(opts, "IntendWithdraw", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardIntendWithdraw)
				if err := _Guard.contract.UnpackLog(event, "IntendWithdraw", log); err != nil {
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

// GuardPunishIterator is returned from FilterPunish and is used to iterate over the raw logs and unpacked data for Punish events raised by the Guard contract.
type GuardPunishIterator struct {
	Event *GuardPunish // Event containing the contract specifics and raw log

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
func (it *GuardPunishIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardPunish)
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
		it.Event = new(GuardPunish)
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
func (it *GuardPunishIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardPunishIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardPunish represents a Punish event raised by the Guard contract.
type GuardPunish struct {
	Indemnitor common.Address
	Indemnitee common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPunish is a free log retrieval operation binding the contract event 0x111897aba775ed4cb659e35805c453dcd8f0024cc414f560f9677bdfae30952a.
//
// Solidity: event Punish(address indexed indemnitor, address indexed indemnitee, uint256 amount)
func (_Guard *GuardFilterer) FilterPunish(opts *bind.FilterOpts, indemnitor []common.Address, indemnitee []common.Address) (*GuardPunishIterator, error) {

	var indemnitorRule []interface{}
	for _, indemnitorItem := range indemnitor {
		indemnitorRule = append(indemnitorRule, indemnitorItem)
	}
	var indemniteeRule []interface{}
	for _, indemniteeItem := range indemnitee {
		indemniteeRule = append(indemniteeRule, indemniteeItem)
	}

	logs, sub, err := _Guard.contract.FilterLogs(opts, "Punish", indemnitorRule, indemniteeRule)
	if err != nil {
		return nil, err
	}
	return &GuardPunishIterator{contract: _Guard.contract, event: "Punish", logs: logs, sub: sub}, nil
}

// WatchPunish is a free log subscription operation binding the contract event 0x111897aba775ed4cb659e35805c453dcd8f0024cc414f560f9677bdfae30952a.
//
// Solidity: event Punish(address indexed indemnitor, address indexed indemnitee, uint256 amount)
func (_Guard *GuardFilterer) WatchPunish(opts *bind.WatchOpts, sink chan<- *GuardPunish, indemnitor []common.Address, indemnitee []common.Address) (event.Subscription, error) {

	var indemnitorRule []interface{}
	for _, indemnitorItem := range indemnitor {
		indemnitorRule = append(indemnitorRule, indemnitorItem)
	}
	var indemniteeRule []interface{}
	for _, indemniteeItem := range indemnitee {
		indemniteeRule = append(indemniteeRule, indemniteeItem)
	}

	logs, sub, err := _Guard.contract.WatchLogs(opts, "Punish", indemnitorRule, indemniteeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardPunish)
				if err := _Guard.contract.UnpackLog(event, "Punish", log); err != nil {
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

// GuardSubscriptionIterator is returned from FilterSubscription and is used to iterate over the raw logs and unpacked data for Subscription events raised by the Guard contract.
type GuardSubscriptionIterator struct {
	Event *GuardSubscription // Event containing the contract specifics and raw log

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
func (it *GuardSubscriptionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardSubscription)
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
		it.Event = new(GuardSubscription)
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
func (it *GuardSubscriptionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardSubscriptionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardSubscription represents a Subscription event raised by the Guard contract.
type GuardSubscription struct {
	Consumer               common.Address
	Amount                 *big.Int
	SubscriptionExpiration *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSubscription is a free log retrieval operation binding the contract event 0x8eb8bcc7421f99b92ab7d727e056544fb59514f8f56251e69658c21ece8977fa.
//
// Solidity: event Subscription(address indexed consumer, uint256 amount, uint256 subscriptionExpiration)
func (_Guard *GuardFilterer) FilterSubscription(opts *bind.FilterOpts, consumer []common.Address) (*GuardSubscriptionIterator, error) {

	var consumerRule []interface{}
	for _, consumerItem := range consumer {
		consumerRule = append(consumerRule, consumerItem)
	}

	logs, sub, err := _Guard.contract.FilterLogs(opts, "Subscription", consumerRule)
	if err != nil {
		return nil, err
	}
	return &GuardSubscriptionIterator{contract: _Guard.contract, event: "Subscription", logs: logs, sub: sub}, nil
}

// WatchSubscription is a free log subscription operation binding the contract event 0x8eb8bcc7421f99b92ab7d727e056544fb59514f8f56251e69658c21ece8977fa.
//
// Solidity: event Subscription(address indexed consumer, uint256 amount, uint256 subscriptionExpiration)
func (_Guard *GuardFilterer) WatchSubscription(opts *bind.WatchOpts, sink chan<- *GuardSubscription, consumer []common.Address) (event.Subscription, error) {

	var consumerRule []interface{}
	for _, consumerItem := range consumer {
		consumerRule = append(consumerRule, consumerItem)
	}

	logs, sub, err := _Guard.contract.WatchLogs(opts, "Subscription", consumerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardSubscription)
				if err := _Guard.contract.UnpackLog(event, "Subscription", log); err != nil {
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

// GuardUpdateSidechainAddrIterator is returned from FilterUpdateSidechainAddr and is used to iterate over the raw logs and unpacked data for UpdateSidechainAddr events raised by the Guard contract.
type GuardUpdateSidechainAddrIterator struct {
	Event *GuardUpdateSidechainAddr // Event containing the contract specifics and raw log

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
func (it *GuardUpdateSidechainAddrIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardUpdateSidechainAddr)
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
		it.Event = new(GuardUpdateSidechainAddr)
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
func (it *GuardUpdateSidechainAddrIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardUpdateSidechainAddrIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardUpdateSidechainAddr represents a UpdateSidechainAddr event raised by the Guard contract.
type GuardUpdateSidechainAddr struct {
	Candidate        common.Address
	OldSidechainAddr common.Hash
	NewSidechainAddr common.Hash
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpdateSidechainAddr is a free log retrieval operation binding the contract event 0x16de3299ab034ce7e21b22d55f4f9a1474bd3c4d20dbd1cc9bcd39c1ad3d5a2c.
//
// Solidity: event UpdateSidechainAddr(address indexed candidate, bytes indexed oldSidechainAddr, bytes indexed newSidechainAddr)
func (_Guard *GuardFilterer) FilterUpdateSidechainAddr(opts *bind.FilterOpts, candidate []common.Address, oldSidechainAddr [][]byte, newSidechainAddr [][]byte) (*GuardUpdateSidechainAddrIterator, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}
	var oldSidechainAddrRule []interface{}
	for _, oldSidechainAddrItem := range oldSidechainAddr {
		oldSidechainAddrRule = append(oldSidechainAddrRule, oldSidechainAddrItem)
	}
	var newSidechainAddrRule []interface{}
	for _, newSidechainAddrItem := range newSidechainAddr {
		newSidechainAddrRule = append(newSidechainAddrRule, newSidechainAddrItem)
	}

	logs, sub, err := _Guard.contract.FilterLogs(opts, "UpdateSidechainAddr", candidateRule, oldSidechainAddrRule, newSidechainAddrRule)
	if err != nil {
		return nil, err
	}
	return &GuardUpdateSidechainAddrIterator{contract: _Guard.contract, event: "UpdateSidechainAddr", logs: logs, sub: sub}, nil
}

// WatchUpdateSidechainAddr is a free log subscription operation binding the contract event 0x16de3299ab034ce7e21b22d55f4f9a1474bd3c4d20dbd1cc9bcd39c1ad3d5a2c.
//
// Solidity: event UpdateSidechainAddr(address indexed candidate, bytes indexed oldSidechainAddr, bytes indexed newSidechainAddr)
func (_Guard *GuardFilterer) WatchUpdateSidechainAddr(opts *bind.WatchOpts, sink chan<- *GuardUpdateSidechainAddr, candidate []common.Address, oldSidechainAddr [][]byte, newSidechainAddr [][]byte) (event.Subscription, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}
	var oldSidechainAddrRule []interface{}
	for _, oldSidechainAddrItem := range oldSidechainAddr {
		oldSidechainAddrRule = append(oldSidechainAddrRule, oldSidechainAddrItem)
	}
	var newSidechainAddrRule []interface{}
	for _, newSidechainAddrItem := range newSidechainAddr {
		newSidechainAddrRule = append(newSidechainAddrRule, newSidechainAddrItem)
	}

	logs, sub, err := _Guard.contract.WatchLogs(opts, "UpdateSidechainAddr", candidateRule, oldSidechainAddrRule, newSidechainAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardUpdateSidechainAddr)
				if err := _Guard.contract.UnpackLog(event, "UpdateSidechainAddr", log); err != nil {
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

// GuardValidatorChangeIterator is returned from FilterValidatorChange and is used to iterate over the raw logs and unpacked data for ValidatorChange events raised by the Guard contract.
type GuardValidatorChangeIterator struct {
	Event *GuardValidatorChange // Event containing the contract specifics and raw log

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
func (it *GuardValidatorChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardValidatorChange)
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
		it.Event = new(GuardValidatorChange)
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
func (it *GuardValidatorChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardValidatorChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardValidatorChange represents a ValidatorChange event raised by the Guard contract.
type GuardValidatorChange struct {
	EthAddr    common.Address
	ChangeType uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterValidatorChange is a free log retrieval operation binding the contract event 0x63f783ba869265648de5e70add96be9f4914e3bde064fdc19fd7e6a8ebf2f46c.
//
// Solidity: event ValidatorChange(address indexed ethAddr, uint8 indexed changeType)
func (_Guard *GuardFilterer) FilterValidatorChange(opts *bind.FilterOpts, ethAddr []common.Address, changeType []uint8) (*GuardValidatorChangeIterator, error) {

	var ethAddrRule []interface{}
	for _, ethAddrItem := range ethAddr {
		ethAddrRule = append(ethAddrRule, ethAddrItem)
	}
	var changeTypeRule []interface{}
	for _, changeTypeItem := range changeType {
		changeTypeRule = append(changeTypeRule, changeTypeItem)
	}

	logs, sub, err := _Guard.contract.FilterLogs(opts, "ValidatorChange", ethAddrRule, changeTypeRule)
	if err != nil {
		return nil, err
	}
	return &GuardValidatorChangeIterator{contract: _Guard.contract, event: "ValidatorChange", logs: logs, sub: sub}, nil
}

// WatchValidatorChange is a free log subscription operation binding the contract event 0x63f783ba869265648de5e70add96be9f4914e3bde064fdc19fd7e6a8ebf2f46c.
//
// Solidity: event ValidatorChange(address indexed ethAddr, uint8 indexed changeType)
func (_Guard *GuardFilterer) WatchValidatorChange(opts *bind.WatchOpts, sink chan<- *GuardValidatorChange, ethAddr []common.Address, changeType []uint8) (event.Subscription, error) {

	var ethAddrRule []interface{}
	for _, ethAddrItem := range ethAddr {
		ethAddrRule = append(ethAddrRule, ethAddrItem)
	}
	var changeTypeRule []interface{}
	for _, changeTypeItem := range changeType {
		changeTypeRule = append(changeTypeRule, changeTypeItem)
	}

	logs, sub, err := _Guard.contract.WatchLogs(opts, "ValidatorChange", ethAddrRule, changeTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardValidatorChange)
				if err := _Guard.contract.UnpackLog(event, "ValidatorChange", log); err != nil {
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
