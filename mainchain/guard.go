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
const GuardABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"feePerBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"subscriptionExpiration\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VALIDATOR_SET_MAX_SIZE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"celerToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"candidateProfiles\",\"outputs\":[{\"name\":\"stakes\",\"type\":\"uint256\"},{\"name\":\"sidechainAddr\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validatorSet\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"newStake\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"totalStake\",\"type\":\"uint256\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"ethAddr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"sidechainAddr\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"added\",\"type\":\"bool\"}],\"name\":\"ValidatorUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"IntendWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ConfirmWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"indemnitor\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"indemnitee\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Punish\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"consumer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"subscriptionExpiration\",\"type\":\"uint256\"}],\"name\":\"Subscription\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"stake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_sidechainAddr\",\"type\":\"bytes\"}],\"name\":\"claimValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"intendWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"confirmWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_punishRequest\",\"type\":\"bytes\"}],\"name\":\"punish\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_indemnitor\",\"type\":\"address\"},{\"name\":\"_indemnitee\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mockPunish\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"subscribe\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// GuardBin is the compiled bytecode used for deploying new contracts.
const GuardBin = `0x608060405234801561001057600080fd5b50600a600e55610a2f806100256000396000f3fe608060405234801561001057600080fd5b50600436106100ec576000357c0100000000000000000000000000000000000000000000000000000000900480637121435b116100a9578063859012a911610083578063859012a9146102a2578063c6c21e9d146102ce578063d678bded146102f2578063e64808f314610397576100ec565b80637121435b146101d057806378df5d2e146102405780637acb775714610276576100ec565b80630f574ba7146100f15780633620d149146101105780634415f47a1461018057806355bf3b331461019a57806363a20c06146101c057806368124f9a146101c8575b600080fd5b61010e6004803603602081101561010757600080fd5b50356103b4565b005b61010e6004803603602081101561012657600080fd5b81019060208101813564010000000081111561014157600080fd5b82018360208201111561015357600080fd5b8035906020019184600183028401116401000000008311171561017557600080fd5b509092509050610486565b6101886104c4565b60408051918252519081900360200190f35b610188600480360360208110156101b057600080fd5b5035600160a060020a03166104ca565b6101886104dc565b61010e6104e1565b61010e600480360360208110156101e657600080fd5b81019060208101813564010000000081111561020157600080fd5b82018360208201111561021357600080fd5b8035906020019184600183028401116401000000008311171561023557600080fd5b50909250905061057d565b61010e6004803603606081101561025657600080fd5b50600160a060020a0381358116916020810135909116906040013561061e565b61010e6004803603604081101561028c57600080fd5b5080359060200135600160a060020a031661067a565b61010e600480360360408110156102b857600080fd5b5080359060200135600160a060020a03166106f2565b6102d661077b565b60408051600160a060020a039092168252519081900360200190f35b6103186004803603602081101561030857600080fd5b5035600160a060020a031661078a565b6040518083815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561035b578181015183820152602001610343565b50505050905090810190601f1680156103885780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6102d6600480360360208110156103ad57600080fd5b5035610831565b600e54339060009083816103c457fe5b600160a060020a0384166000908152600f6020526040902054919004915043111561040b57600160a060020a0382166000908152600f60205260409020438201905561042a565b600160a060020a0382166000908152600f602052604090208054820190555b600160a060020a0382166000818152600f602090815260409182902054825193845290830186905282820152517f8eb8bcc7421f99b92ab7d727e056544fb59514f8f56251e69658c21ece8977fa9181900360600190a1505050565b73e0b6b1e22182ae2b8382bac06f5392dad89ebf0473f0d9fcb4fefdbd3e7929374b4632f8ad511bd7e360646104bd83838361061e565b5050505050565b600e5481565b600f6020526000908152604090205481565b600b81565b336000818152600d60205260409020600180820154600290920154610514928492600160a060020a03909116919061084e565b600160a060020a038082166000818152600d6020908152604091829020600181015460029091015483519485529416908301528181019290925290517f08d0283ea9a2e520a2f09611cf37ca6eb70f62e9a807e53756047dd2dc0272209181900360600190a150565b336000818152600c6020526040902061059a90600101848461095f565b5060408051600160a060020a038316815260019181018290526060602082018181529082018590527f19cf6af75cb9d1290db4d6ddf4143e0808ad392f53f9e05e38b2d2519a7d797992849287928792919060808201858580828437600083820152604051601f909101601f191690920182900397509095505050505050a1505050565b61062b838483600161084e565b60408051600160a060020a0380861682528416602082015280820183905290517f111897aba775ed4cb659e35805c453dcd8f0024cc414f560f9677bdfae30952a9181900360600190a1505050565b33610688818385600061084e565b600160a060020a038083166000818152600c60209081526040918290205482519486168552908401929092528281018690526060830191909152517f63602d0ecc7b3a0ef7ff1a116e23056662d64280355ba8031b6d0d767c4b44589181900360800190a1505050565b336000818152600d6020908152604091829020600181018054600160a060020a03871673ffffffffffffffffffffffffffffffffffffffff199091168117909155600290910186905582518481529182015280820185905290517ff0835827db57a4706bf4c0686f93f0d0bfbe30895d855f0b9585a2b3f0d1f3489181900360600190a1505050565b600054600160a060020a031681565b600c602090815260009182526040918290208054600180830180548651600293821615610100026000190190911692909204601f8101869004860283018601909652858252919492939092908301828280156108275780601f106107fc57610100808354040283529160200191610827565b820191906000526020600020905b81548152906001019060200180831161080a57829003601f168201915b5050505050905082565b600181600b811061083e57fe5b0154600160a060020a0316905081565b600081600181111561085c57fe5b141561089f57600160a060020a038085166000908152600d602090815260408083209387168352928152828220805486019055600c905220805483019055610959565b60018160018111156108ad57fe5b14156108f257600160a060020a038085166000908152600d60209081526040808320938716835292815282822080548690039055600c90522080548390039055610959565b604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f496e76616c6964206f7065726174696f6e000000000000000000000000000000604482015290519081900360640190fd5b50505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106109a05782800160ff198235161785556109cd565b828001600101855582156109cd579182015b828111156109cd5782358255916020019190600101906109b2565b506109d99291506109dd565b5090565b6109f791905b808211156109d957600081556001016109e3565b9056fea265627a7a7230582082263effc5ebe9c72d01632a263afc74d71cf2375ee70773c5435a8378c2a36164736f6c634300050a0032`

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

// CandidateProfiles is a free data retrieval call binding the contract method 0xd678bded.
//
// Solidity: function candidateProfiles(address ) constant returns(uint256 stakes, bytes sidechainAddr)
func (_Guard *GuardCaller) CandidateProfiles(opts *bind.CallOpts, arg0 common.Address) (struct {
	Stakes        *big.Int
	SidechainAddr []byte
}, error) {
	ret := new(struct {
		Stakes        *big.Int
		SidechainAddr []byte
	})
	out := ret
	err := _Guard.contract.Call(opts, out, "candidateProfiles", arg0)
	return *ret, err
}

// CandidateProfiles is a free data retrieval call binding the contract method 0xd678bded.
//
// Solidity: function candidateProfiles(address ) constant returns(uint256 stakes, bytes sidechainAddr)
func (_Guard *GuardSession) CandidateProfiles(arg0 common.Address) (struct {
	Stakes        *big.Int
	SidechainAddr []byte
}, error) {
	return _Guard.Contract.CandidateProfiles(&_Guard.CallOpts, arg0)
}

// CandidateProfiles is a free data retrieval call binding the contract method 0xd678bded.
//
// Solidity: function candidateProfiles(address ) constant returns(uint256 stakes, bytes sidechainAddr)
func (_Guard *GuardCallerSession) CandidateProfiles(arg0 common.Address) (struct {
	Stakes        *big.Int
	SidechainAddr []byte
}, error) {
	return _Guard.Contract.CandidateProfiles(&_Guard.CallOpts, arg0)
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

// ClaimValidator is a paid mutator transaction binding the contract method 0x7121435b.
//
// Solidity: function claimValidator(bytes _sidechainAddr) returns()
func (_Guard *GuardTransactor) ClaimValidator(opts *bind.TransactOpts, _sidechainAddr []byte) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "claimValidator", _sidechainAddr)
}

// ClaimValidator is a paid mutator transaction binding the contract method 0x7121435b.
//
// Solidity: function claimValidator(bytes _sidechainAddr) returns()
func (_Guard *GuardSession) ClaimValidator(_sidechainAddr []byte) (*types.Transaction, error) {
	return _Guard.Contract.ClaimValidator(&_Guard.TransactOpts, _sidechainAddr)
}

// ClaimValidator is a paid mutator transaction binding the contract method 0x7121435b.
//
// Solidity: function claimValidator(bytes _sidechainAddr) returns()
func (_Guard *GuardTransactorSession) ClaimValidator(_sidechainAddr []byte) (*types.Transaction, error) {
	return _Guard.Contract.ClaimValidator(&_Guard.TransactOpts, _sidechainAddr)
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0x68124f9a.
//
// Solidity: function confirmWithdraw() returns()
func (_Guard *GuardTransactor) ConfirmWithdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "confirmWithdraw")
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0x68124f9a.
//
// Solidity: function confirmWithdraw() returns()
func (_Guard *GuardSession) ConfirmWithdraw() (*types.Transaction, error) {
	return _Guard.Contract.ConfirmWithdraw(&_Guard.TransactOpts)
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0x68124f9a.
//
// Solidity: function confirmWithdraw() returns()
func (_Guard *GuardTransactorSession) ConfirmWithdraw() (*types.Transaction, error) {
	return _Guard.Contract.ConfirmWithdraw(&_Guard.TransactOpts)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x859012a9.
//
// Solidity: function intendWithdraw(uint256 _amount, address _candidate) returns()
func (_Guard *GuardTransactor) IntendWithdraw(opts *bind.TransactOpts, _amount *big.Int, _candidate common.Address) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "intendWithdraw", _amount, _candidate)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x859012a9.
//
// Solidity: function intendWithdraw(uint256 _amount, address _candidate) returns()
func (_Guard *GuardSession) IntendWithdraw(_amount *big.Int, _candidate common.Address) (*types.Transaction, error) {
	return _Guard.Contract.IntendWithdraw(&_Guard.TransactOpts, _amount, _candidate)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x859012a9.
//
// Solidity: function intendWithdraw(uint256 _amount, address _candidate) returns()
func (_Guard *GuardTransactorSession) IntendWithdraw(_amount *big.Int, _candidate common.Address) (*types.Transaction, error) {
	return _Guard.Contract.IntendWithdraw(&_Guard.TransactOpts, _amount, _candidate)
}

// MockPunish is a paid mutator transaction binding the contract method 0x78df5d2e.
//
// Solidity: function mockPunish(address _indemnitor, address _indemnitee, uint256 _amount) returns()
func (_Guard *GuardTransactor) MockPunish(opts *bind.TransactOpts, _indemnitor common.Address, _indemnitee common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "mockPunish", _indemnitor, _indemnitee, _amount)
}

// MockPunish is a paid mutator transaction binding the contract method 0x78df5d2e.
//
// Solidity: function mockPunish(address _indemnitor, address _indemnitee, uint256 _amount) returns()
func (_Guard *GuardSession) MockPunish(_indemnitor common.Address, _indemnitee common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.MockPunish(&_Guard.TransactOpts, _indemnitor, _indemnitee, _amount)
}

// MockPunish is a paid mutator transaction binding the contract method 0x78df5d2e.
//
// Solidity: function mockPunish(address _indemnitor, address _indemnitee, uint256 _amount) returns()
func (_Guard *GuardTransactorSession) MockPunish(_indemnitor common.Address, _indemnitee common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.MockPunish(&_Guard.TransactOpts, _indemnitor, _indemnitee, _amount)
}

// Punish is a paid mutator transaction binding the contract method 0x3620d149.
//
// Solidity: function punish(bytes _punishRequest) returns()
func (_Guard *GuardTransactor) Punish(opts *bind.TransactOpts, _punishRequest []byte) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "punish", _punishRequest)
}

// Punish is a paid mutator transaction binding the contract method 0x3620d149.
//
// Solidity: function punish(bytes _punishRequest) returns()
func (_Guard *GuardSession) Punish(_punishRequest []byte) (*types.Transaction, error) {
	return _Guard.Contract.Punish(&_Guard.TransactOpts, _punishRequest)
}

// Punish is a paid mutator transaction binding the contract method 0x3620d149.
//
// Solidity: function punish(bytes _punishRequest) returns()
func (_Guard *GuardTransactorSession) Punish(_punishRequest []byte) (*types.Transaction, error) {
	return _Guard.Contract.Punish(&_Guard.TransactOpts, _punishRequest)
}

// Stake is a paid mutator transaction binding the contract method 0x7acb7757.
//
// Solidity: function stake(uint256 _amount, address _candidate) returns()
func (_Guard *GuardTransactor) Stake(opts *bind.TransactOpts, _amount *big.Int, _candidate common.Address) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "stake", _amount, _candidate)
}

// Stake is a paid mutator transaction binding the contract method 0x7acb7757.
//
// Solidity: function stake(uint256 _amount, address _candidate) returns()
func (_Guard *GuardSession) Stake(_amount *big.Int, _candidate common.Address) (*types.Transaction, error) {
	return _Guard.Contract.Stake(&_Guard.TransactOpts, _amount, _candidate)
}

// Stake is a paid mutator transaction binding the contract method 0x7acb7757.
//
// Solidity: function stake(uint256 _amount, address _candidate) returns()
func (_Guard *GuardTransactorSession) Stake(_amount *big.Int, _candidate common.Address) (*types.Transaction, error) {
	return _Guard.Contract.Stake(&_Guard.TransactOpts, _amount, _candidate)
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
// Solidity: event ConfirmWithdraw(address delegator, address candidate, uint256 amount)
func (_Guard *GuardFilterer) FilterConfirmWithdraw(opts *bind.FilterOpts) (*GuardConfirmWithdrawIterator, error) {

	logs, sub, err := _Guard.contract.FilterLogs(opts, "ConfirmWithdraw")
	if err != nil {
		return nil, err
	}
	return &GuardConfirmWithdrawIterator{contract: _Guard.contract, event: "ConfirmWithdraw", logs: logs, sub: sub}, nil
}

// WatchConfirmWithdraw is a free log subscription operation binding the contract event 0x08d0283ea9a2e520a2f09611cf37ca6eb70f62e9a807e53756047dd2dc027220.
//
// Solidity: event ConfirmWithdraw(address delegator, address candidate, uint256 amount)
func (_Guard *GuardFilterer) WatchConfirmWithdraw(opts *bind.WatchOpts, sink chan<- *GuardConfirmWithdraw) (event.Subscription, error) {

	logs, sub, err := _Guard.contract.WatchLogs(opts, "ConfirmWithdraw")
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
	Delegator common.Address
	Candidate common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterIntendWithdraw is a free log retrieval operation binding the contract event 0xf0835827db57a4706bf4c0686f93f0d0bfbe30895d855f0b9585a2b3f0d1f348.
//
// Solidity: event IntendWithdraw(address delegator, address candidate, uint256 amount)
func (_Guard *GuardFilterer) FilterIntendWithdraw(opts *bind.FilterOpts) (*GuardIntendWithdrawIterator, error) {

	logs, sub, err := _Guard.contract.FilterLogs(opts, "IntendWithdraw")
	if err != nil {
		return nil, err
	}
	return &GuardIntendWithdrawIterator{contract: _Guard.contract, event: "IntendWithdraw", logs: logs, sub: sub}, nil
}

// WatchIntendWithdraw is a free log subscription operation binding the contract event 0xf0835827db57a4706bf4c0686f93f0d0bfbe30895d855f0b9585a2b3f0d1f348.
//
// Solidity: event IntendWithdraw(address delegator, address candidate, uint256 amount)
func (_Guard *GuardFilterer) WatchIntendWithdraw(opts *bind.WatchOpts, sink chan<- *GuardIntendWithdraw) (event.Subscription, error) {

	logs, sub, err := _Guard.contract.WatchLogs(opts, "IntendWithdraw")
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
// Solidity: event Punish(address indemnitor, address indemnitee, uint256 amount)
func (_Guard *GuardFilterer) FilterPunish(opts *bind.FilterOpts) (*GuardPunishIterator, error) {

	logs, sub, err := _Guard.contract.FilterLogs(opts, "Punish")
	if err != nil {
		return nil, err
	}
	return &GuardPunishIterator{contract: _Guard.contract, event: "Punish", logs: logs, sub: sub}, nil
}

// WatchPunish is a free log subscription operation binding the contract event 0x111897aba775ed4cb659e35805c453dcd8f0024cc414f560f9677bdfae30952a.
//
// Solidity: event Punish(address indemnitor, address indemnitee, uint256 amount)
func (_Guard *GuardFilterer) WatchPunish(opts *bind.WatchOpts, sink chan<- *GuardPunish) (event.Subscription, error) {

	logs, sub, err := _Guard.contract.WatchLogs(opts, "Punish")
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

// GuardStakeIterator is returned from FilterStake and is used to iterate over the raw logs and unpacked data for Stake events raised by the Guard contract.
type GuardStakeIterator struct {
	Event *GuardStake // Event containing the contract specifics and raw log

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
func (it *GuardStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardStake)
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
		it.Event = new(GuardStake)
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
func (it *GuardStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardStake represents a Stake event raised by the Guard contract.
type GuardStake struct {
	Delegator  common.Address
	Candidate  common.Address
	NewStake   *big.Int
	TotalStake *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterStake is a free log retrieval operation binding the contract event 0x63602d0ecc7b3a0ef7ff1a116e23056662d64280355ba8031b6d0d767c4b4458.
//
// Solidity: event Stake(address delegator, address candidate, uint256 newStake, uint256 totalStake)
func (_Guard *GuardFilterer) FilterStake(opts *bind.FilterOpts) (*GuardStakeIterator, error) {

	logs, sub, err := _Guard.contract.FilterLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return &GuardStakeIterator{contract: _Guard.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0x63602d0ecc7b3a0ef7ff1a116e23056662d64280355ba8031b6d0d767c4b4458.
//
// Solidity: event Stake(address delegator, address candidate, uint256 newStake, uint256 totalStake)
func (_Guard *GuardFilterer) WatchStake(opts *bind.WatchOpts, sink chan<- *GuardStake) (event.Subscription, error) {

	logs, sub, err := _Guard.contract.WatchLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardStake)
				if err := _Guard.contract.UnpackLog(event, "Stake", log); err != nil {
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
// Solidity: event Subscription(address consumer, uint256 amount, uint256 subscriptionExpiration)
func (_Guard *GuardFilterer) FilterSubscription(opts *bind.FilterOpts) (*GuardSubscriptionIterator, error) {

	logs, sub, err := _Guard.contract.FilterLogs(opts, "Subscription")
	if err != nil {
		return nil, err
	}
	return &GuardSubscriptionIterator{contract: _Guard.contract, event: "Subscription", logs: logs, sub: sub}, nil
}

// WatchSubscription is a free log subscription operation binding the contract event 0x8eb8bcc7421f99b92ab7d727e056544fb59514f8f56251e69658c21ece8977fa.
//
// Solidity: event Subscription(address consumer, uint256 amount, uint256 subscriptionExpiration)
func (_Guard *GuardFilterer) WatchSubscription(opts *bind.WatchOpts, sink chan<- *GuardSubscription) (event.Subscription, error) {

	logs, sub, err := _Guard.contract.WatchLogs(opts, "Subscription")
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

// GuardValidatorUpdateIterator is returned from FilterValidatorUpdate and is used to iterate over the raw logs and unpacked data for ValidatorUpdate events raised by the Guard contract.
type GuardValidatorUpdateIterator struct {
	Event *GuardValidatorUpdate // Event containing the contract specifics and raw log

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
func (it *GuardValidatorUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardValidatorUpdate)
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
		it.Event = new(GuardValidatorUpdate)
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
func (it *GuardValidatorUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardValidatorUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardValidatorUpdate represents a ValidatorUpdate event raised by the Guard contract.
type GuardValidatorUpdate struct {
	EthAddr       common.Address
	SidechainAddr []byte
	Added         bool
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterValidatorUpdate is a free log retrieval operation binding the contract event 0x19cf6af75cb9d1290db4d6ddf4143e0808ad392f53f9e05e38b2d2519a7d7979.
//
// Solidity: event ValidatorUpdate(address ethAddr, bytes sidechainAddr, bool added)
func (_Guard *GuardFilterer) FilterValidatorUpdate(opts *bind.FilterOpts) (*GuardValidatorUpdateIterator, error) {

	logs, sub, err := _Guard.contract.FilterLogs(opts, "ValidatorUpdate")
	if err != nil {
		return nil, err
	}
	return &GuardValidatorUpdateIterator{contract: _Guard.contract, event: "ValidatorUpdate", logs: logs, sub: sub}, nil
}

// WatchValidatorUpdate is a free log subscription operation binding the contract event 0x19cf6af75cb9d1290db4d6ddf4143e0808ad392f53f9e05e38b2d2519a7d7979.
//
// Solidity: event ValidatorUpdate(address ethAddr, bytes sidechainAddr, bool added)
func (_Guard *GuardFilterer) WatchValidatorUpdate(opts *bind.WatchOpts, sink chan<- *GuardValidatorUpdate) (event.Subscription, error) {

	logs, sub, err := _Guard.contract.WatchLogs(opts, "ValidatorUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardValidatorUpdate)
				if err := _Guard.contract.UnpackLog(event, "ValidatorUpdate", log); err != nil {
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
