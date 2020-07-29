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

// SGNABI is the input ABI used to generate the binding from.
const SGNABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"subscriptionDeposits\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"DPoSContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"redeemedServiceReward\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"celerToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"sidechainAddrMap\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"servicePool\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_celerTokenAddress\",\"type\":\"address\"},{\"name\":\"_DPoSAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"oldSidechainAddr\",\"type\":\"bytes\"},{\"indexed\":true,\"name\":\"newSidechainAddr\",\"type\":\"bytes\"}],\"name\":\"UpdateSidechainAddr\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"consumer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AddSubscriptionBalance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"cumulativeMiningReward\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"serviceReward\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"servicePool\",\"type\":\"uint256\"}],\"name\":\"RedeemReward\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_sidechainAddr\",\"type\":\"bytes\"}],\"name\":\"updateSidechainAddr\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"subscribe\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_rewardRequest\",\"type\":\"bytes\"}],\"name\":\"redeemReward\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SGNBin is the compiled bytecode used for deploying new contracts.
var SGNBin = "0x608060405234801561001057600080fd5b506040516112d33803806112d38339818101604052604081101561003357600080fd5b508051602090910151600080546001600160a01b039384166001600160a01b031991821617909155600180549390921692169190911790556112598061007a6000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063c57f666111610066578063c57f666114610183578063c6c21e9d146101a9578063e02f39bd146101b1578063e27b411014610221578063e42a06c8146102bc57610093565b80630f574ba7146100985780631e77733a146100b757806371273548146100ef57806373a6e45014610113575b600080fd5b6100b5600480360360208110156100ae57600080fd5b50356102c4565b005b6100dd600480360360208110156100cd57600080fd5b50356001600160a01b0316610439565b60408051918252519081900360200190f35b6100f761044b565b604080516001600160a01b039092168252519081900360200190f35b6100b56004803603602081101561012957600080fd5b81019060208101813564010000000081111561014457600080fd5b82018360208201111561015657600080fd5b8035906020019184600183028401116401000000008311171561017857600080fd5b50909250905061045a565b6100dd6004803603602081101561019957600080fd5b50356001600160a01b03166107c3565b6100f76107d5565b6100b5600480360360208110156101c757600080fd5b8101906020810181356401000000008111156101e257600080fd5b8201836020820111156101f457600080fd5b8035906020019184600183028401116401000000008311171561021657600080fd5b5090925090506107e4565b6102476004803603602081101561023757600080fd5b50356001600160a01b0316610a91565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610281578181015183820152602001610269565b50505050905090810190601f1680156102ae5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100dd610b2c565b600160009054906101000a90046001600160a01b03166001600160a01b031663eab2ed8c6040518163ffffffff1660e01b815260040160206040518083038186803b15801561031257600080fd5b505afa158015610326573d6000803e3d6000fd5b505050506040513d602081101561033c57600080fd5b5051610383576040805162461bcd60e51b815260206004820152601160248201527011141bd4c81a5cc81b9bdd081d985b1a59607a1b604482015290519081900360640190fd5b6003543390610398908363ffffffff610b3216565b6003556001600160a01b0381166000908152600260205260409020546103c4908363ffffffff610b3216565b6001600160a01b0380831660009081526002602052604081209290925590546103f6911682308563ffffffff610b4b16565b6040805183815290516001600160a01b038316917fac095ced75d400384d8404a37883183a56b937b8ff8195fa0c52c3ccc8bb8a49919081900360200190a25050565b60026020526000908152604090205481565b6001546001600160a01b031681565b600160009054906101000a90046001600160a01b03166001600160a01b031663eab2ed8c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156104a857600080fd5b505afa1580156104bc573d6000803e3d6000fd5b505050506040513d60208110156104d257600080fd5b5051610519576040805162461bcd60e51b815260206004820152601160248201527011141bd4c81a5cc81b9bdd081d985b1a59607a1b604482015290519081900360640190fd5b600154604051631c0efd9d60e01b8152602060048201908152602482018490526001600160a01b0390921691631c0efd9d91859185918190604401848480828437600081840152601f19601f8201169050808301925050509350505050602060405180830381600087803b15801561059057600080fd5b505af11580156105a4573d6000803e3d6000fd5b505050506040513d60208110156105ba57600080fd5b505161060d576040805162461bcd60e51b815260206004820152601c60248201527f4661696c20746f20636865636b2076616c696461746f72207369677300000000604482015290519081900360640190fd5b61061561112b565b61065483838080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610be192505050565b905061065e611145565b815161066990610d38565b80516001600160a01b03166000908152600460205260408082205490830151929350909161069c9163ffffffff610dfd16565b60408084015184516001600160a01b03166000908152600460205291909120556003549091506106d2908263ffffffff610dfd16565b6003556001548251602084015160408051630fbd844360e11b81526001600160a01b0393841660048201526024810192909252519190921691631f7b088691604480830192600092919082900301818387803b15801561073157600080fd5b505af1158015610745573d6000803e3d6000fd5b5050835160005461076893506001600160a01b031691508363ffffffff610e1216565b8151602080840151600354604080519283529282018590528183015290516001600160a01b03909216917f09251621f2e88c5e7f8df91fe1d9e9a70610e20e122945470fddd48af05104269181900360600190a25050505050565b60046020526000908152604090205481565b6000546001600160a01b031681565b600154604080516328bde1e160e01b815233600482018190529151919260009283926001600160a01b03909216916328bde1e19160248083019260e0929190829003018186803b15801561083757600080fd5b505afa15801561084b573d6000803e3d6000fd5b505050506040513d60e081101561086157600080fd5b508051606090910151909250905080156108c2576040805162461bcd60e51b815260206004820152601a60248201527f6d73672e73656e646572206973206e6f7420756e626f6e646564000000000000604482015290519081900360640190fd5b81610914576040805162461bcd60e51b815260206004820152601c60248201527f43616e646964617465206973206e6f7420696e697469616c697a656400000000604482015290519081900360640190fd5b6001600160a01b03831660009081526005602090815260409182902080548351601f60026000196101006001861615020190931692909204918201849004840281018401909452808452606093928301828280156109b35780601f10610988576101008083540402835291602001916109b3565b820191906000526020600020905b81548152906001019060200180831161099657829003601f168201915b505050506001600160a01b03861660009081526005602052604090209192506109df919050878761116f565b5085856040518083838082843760405192018290038220865190955086945091925082916020850191508083835b60208310610a2c5780518252601f199092019160209182019101610a0d565b5181516020939093036101000a6000190180199091169216919091179052604051920182900382209350506001600160a01b03881691507f16de3299ab034ce7e21b22d55f4f9a1474bd3c4d20dbd1cc9bcd39c1ad3d5a2c90600090a4505050505050565b60056020908152600091825260409182902080548351601f600260001961010060018616150201909316929092049182018490048402810184019094528084529091830182828015610b245780601f10610af957610100808354040283529160200191610b24565b820191906000526020600020905b815481529060010190602001808311610b0757829003601f168201915b505050505081565b60035481565b600082820183811015610b4457600080fd5b9392505050565b604080516323b872dd60e01b81526001600160a01b0385811660048301528481166024830152604482018490529151918616916323b872dd916064808201926020929091908290030181600087803b158015610ba657600080fd5b505af1158015610bba573d6000803e3d6000fd5b505050506040513d6020811015610bd057600080fd5b5051610bdb57600080fd5b50505050565b610be961112b565b610bf16111ed565b610bfa83610eac565b90506060610c0f82600263ffffffff610ec316565b905080600281518110610c1e57fe5b6020026020010151604051908082528060200260200182016040528015610c5957816020015b6060815260200190600190039081610c445790505b508360200181905250600081600281518110610c7157fe5b6020026020010181815250506000805b610c8a84610f53565b15610d2f57610c9884610f62565b90925090508160011415610cb657610caf84610f8f565b8552610d2a565b8160021415610d1a57610cc884610f8f565b856020015184600281518110610cda57fe5b602002602001015181518110610cec57fe5b602002602001018190525082600281518110610d0457fe5b6020908102919091010180516001019052610d2a565b610d2a848263ffffffff61101c16565b610c81565b50505050919050565b610d40611145565b610d486111ed565b610d5183610eac565b90506000805b610d6083610f53565b15610df557610d6e83610f62565b90925090508160011415610d9d57610d8d610d8884610f8f565b61107d565b6001600160a01b03168452610df0565b8160021415610dc157610db7610db284610f8f565b61108e565b6020850152610df0565b8160031415610de057610dd6610db284610f8f565b6040850152610df0565b610df0838263ffffffff61101c16565b610d57565b505050919050565b600082821115610e0c57600080fd5b50900390565b826001600160a01b031663a9059cbb83836040518363ffffffff1660e01b815260040180836001600160a01b03166001600160a01b0316815260200182815260200192505050602060405180830381600087803b158015610e7257600080fd5b505af1158015610e86573d6000803e3d6000fd5b505050506040513d6020811015610e9c57600080fd5b5051610ea757600080fd5b505050565b610eb46111ed565b60208101919091526000815290565b815160408051600184018082526020808202830101909252606092918015610ef5578160200160208202803883390190505b5091506000805b610f0586610f53565b15610f4a57610f1386610f62565b80925081935050506001848381518110610f2957fe5b602002602001018181510191508181525050610f45868261101c565b610efc565b50509092525090565b6020810151518151105b919050565b6000806000610f70846110b1565b9050600881049250806007166005811115610f8757fe5b915050915091565b60606000610f9c836110b1565b8351602085015151919250820190811115610fb657600080fd5b816040519080825280601f01601f191660200182016040528015610fe1576020820181803883390190505b50602080860151865192955091818601919083010160005b85811015611011578181015183820152602001610ff9565b505050935250919050565b600081600581111561102a57fe5b141561103f57611039826110b1565b50611079565b600281600581111561104d57fe5b141561009357600061105e836110b1565b83518101808552602085015151919250111561103957600080fd5b5050565b60006110888261110c565b92915050565b600060208251111561109f57600080fd5b50602081810151915160089103021c90565b602080820151825181019091015160009182805b600a8110156111065783811a91508060070282607f16901b8517945081608016600014156110fe57855101600101855250610f5d915050565b6001016110c5565b50600080fd5b6000815160141461111c57600080fd5b5060200151600160601b900490565b604051806040016040528060608152602001606081525090565b604051806060016040528060006001600160a01b0316815260200160008152602001600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106111b05782800160ff198235161785556111dd565b828001600101855582156111dd579182015b828111156111dd5782358255916020019190600101906111c2565b506111e9929150611207565b5090565b604051806040016040528060008152602001606081525090565b61122191905b808211156111e9576000815560010161120d565b9056fea265627a7a72305820c51a318f09cf7fea309e928e78697896da571f903214810875f0f6b083ea968c64736f6c634300050a0032"

// DeploySGN deploys a new Ethereum contract, binding an instance of SGN to it.
func DeploySGN(auth *bind.TransactOpts, backend bind.ContractBackend, _celerTokenAddress common.Address, _DPoSAddress common.Address) (common.Address, *types.Transaction, *SGN, error) {
	parsed, err := abi.JSON(strings.NewReader(SGNABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SGNBin), backend, _celerTokenAddress, _DPoSAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SGN{SGNCaller: SGNCaller{contract: contract}, SGNTransactor: SGNTransactor{contract: contract}, SGNFilterer: SGNFilterer{contract: contract}}, nil
}

// SGN is an auto generated Go binding around an Ethereum contract.
type SGN struct {
	SGNCaller     // Read-only binding to the contract
	SGNTransactor // Write-only binding to the contract
	SGNFilterer   // Log filterer for contract events
}

// SGNCaller is an auto generated read-only Go binding around an Ethereum contract.
type SGNCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SGNTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SGNTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SGNFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SGNFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SGNSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SGNSession struct {
	Contract     *SGN              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SGNCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SGNCallerSession struct {
	Contract *SGNCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SGNTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SGNTransactorSession struct {
	Contract     *SGNTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SGNRaw is an auto generated low-level Go binding around an Ethereum contract.
type SGNRaw struct {
	Contract *SGN // Generic contract binding to access the raw methods on
}

// SGNCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SGNCallerRaw struct {
	Contract *SGNCaller // Generic read-only contract binding to access the raw methods on
}

// SGNTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SGNTransactorRaw struct {
	Contract *SGNTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSGN creates a new instance of SGN, bound to a specific deployed contract.
func NewSGN(address common.Address, backend bind.ContractBackend) (*SGN, error) {
	contract, err := bindSGN(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SGN{SGNCaller: SGNCaller{contract: contract}, SGNTransactor: SGNTransactor{contract: contract}, SGNFilterer: SGNFilterer{contract: contract}}, nil
}

// NewSGNCaller creates a new read-only instance of SGN, bound to a specific deployed contract.
func NewSGNCaller(address common.Address, caller bind.ContractCaller) (*SGNCaller, error) {
	contract, err := bindSGN(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SGNCaller{contract: contract}, nil
}

// NewSGNTransactor creates a new write-only instance of SGN, bound to a specific deployed contract.
func NewSGNTransactor(address common.Address, transactor bind.ContractTransactor) (*SGNTransactor, error) {
	contract, err := bindSGN(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SGNTransactor{contract: contract}, nil
}

// NewSGNFilterer creates a new log filterer instance of SGN, bound to a specific deployed contract.
func NewSGNFilterer(address common.Address, filterer bind.ContractFilterer) (*SGNFilterer, error) {
	contract, err := bindSGN(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SGNFilterer{contract: contract}, nil
}

// bindSGN binds a generic wrapper to an already deployed contract.
func bindSGN(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SGNABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SGN *SGNRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SGN.Contract.SGNCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SGN *SGNRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SGN.Contract.SGNTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SGN *SGNRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SGN.Contract.SGNTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SGN *SGNCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SGN.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SGN *SGNTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SGN.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SGN *SGNTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SGN.Contract.contract.Transact(opts, method, params...)
}

// DPoSContract is a free data retrieval call binding the contract method 0x71273548.
//
// Solidity: function DPoSContract() constant returns(address)
func (_SGN *SGNCaller) DPoSContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SGN.contract.Call(opts, out, "DPoSContract")
	return *ret0, err
}

// DPoSContract is a free data retrieval call binding the contract method 0x71273548.
//
// Solidity: function DPoSContract() constant returns(address)
func (_SGN *SGNSession) DPoSContract() (common.Address, error) {
	return _SGN.Contract.DPoSContract(&_SGN.CallOpts)
}

// DPoSContract is a free data retrieval call binding the contract method 0x71273548.
//
// Solidity: function DPoSContract() constant returns(address)
func (_SGN *SGNCallerSession) DPoSContract() (common.Address, error) {
	return _SGN.Contract.DPoSContract(&_SGN.CallOpts)
}

// CelerToken is a free data retrieval call binding the contract method 0xc6c21e9d.
//
// Solidity: function celerToken() constant returns(address)
func (_SGN *SGNCaller) CelerToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SGN.contract.Call(opts, out, "celerToken")
	return *ret0, err
}

// CelerToken is a free data retrieval call binding the contract method 0xc6c21e9d.
//
// Solidity: function celerToken() constant returns(address)
func (_SGN *SGNSession) CelerToken() (common.Address, error) {
	return _SGN.Contract.CelerToken(&_SGN.CallOpts)
}

// CelerToken is a free data retrieval call binding the contract method 0xc6c21e9d.
//
// Solidity: function celerToken() constant returns(address)
func (_SGN *SGNCallerSession) CelerToken() (common.Address, error) {
	return _SGN.Contract.CelerToken(&_SGN.CallOpts)
}

// RedeemedServiceReward is a free data retrieval call binding the contract method 0xc57f6661.
//
// Solidity: function redeemedServiceReward(address ) constant returns(uint256)
func (_SGN *SGNCaller) RedeemedServiceReward(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SGN.contract.Call(opts, out, "redeemedServiceReward", arg0)
	return *ret0, err
}

// RedeemedServiceReward is a free data retrieval call binding the contract method 0xc57f6661.
//
// Solidity: function redeemedServiceReward(address ) constant returns(uint256)
func (_SGN *SGNSession) RedeemedServiceReward(arg0 common.Address) (*big.Int, error) {
	return _SGN.Contract.RedeemedServiceReward(&_SGN.CallOpts, arg0)
}

// RedeemedServiceReward is a free data retrieval call binding the contract method 0xc57f6661.
//
// Solidity: function redeemedServiceReward(address ) constant returns(uint256)
func (_SGN *SGNCallerSession) RedeemedServiceReward(arg0 common.Address) (*big.Int, error) {
	return _SGN.Contract.RedeemedServiceReward(&_SGN.CallOpts, arg0)
}

// ServicePool is a free data retrieval call binding the contract method 0xe42a06c8.
//
// Solidity: function servicePool() constant returns(uint256)
func (_SGN *SGNCaller) ServicePool(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SGN.contract.Call(opts, out, "servicePool")
	return *ret0, err
}

// ServicePool is a free data retrieval call binding the contract method 0xe42a06c8.
//
// Solidity: function servicePool() constant returns(uint256)
func (_SGN *SGNSession) ServicePool() (*big.Int, error) {
	return _SGN.Contract.ServicePool(&_SGN.CallOpts)
}

// ServicePool is a free data retrieval call binding the contract method 0xe42a06c8.
//
// Solidity: function servicePool() constant returns(uint256)
func (_SGN *SGNCallerSession) ServicePool() (*big.Int, error) {
	return _SGN.Contract.ServicePool(&_SGN.CallOpts)
}

// SidechainAddrMap is a free data retrieval call binding the contract method 0xe27b4110.
//
// Solidity: function sidechainAddrMap(address ) constant returns(bytes)
func (_SGN *SGNCaller) SidechainAddrMap(opts *bind.CallOpts, arg0 common.Address) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _SGN.contract.Call(opts, out, "sidechainAddrMap", arg0)
	return *ret0, err
}

// SidechainAddrMap is a free data retrieval call binding the contract method 0xe27b4110.
//
// Solidity: function sidechainAddrMap(address ) constant returns(bytes)
func (_SGN *SGNSession) SidechainAddrMap(arg0 common.Address) ([]byte, error) {
	return _SGN.Contract.SidechainAddrMap(&_SGN.CallOpts, arg0)
}

// SidechainAddrMap is a free data retrieval call binding the contract method 0xe27b4110.
//
// Solidity: function sidechainAddrMap(address ) constant returns(bytes)
func (_SGN *SGNCallerSession) SidechainAddrMap(arg0 common.Address) ([]byte, error) {
	return _SGN.Contract.SidechainAddrMap(&_SGN.CallOpts, arg0)
}

// SubscriptionDeposits is a free data retrieval call binding the contract method 0x1e77733a.
//
// Solidity: function subscriptionDeposits(address ) constant returns(uint256)
func (_SGN *SGNCaller) SubscriptionDeposits(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SGN.contract.Call(opts, out, "subscriptionDeposits", arg0)
	return *ret0, err
}

// SubscriptionDeposits is a free data retrieval call binding the contract method 0x1e77733a.
//
// Solidity: function subscriptionDeposits(address ) constant returns(uint256)
func (_SGN *SGNSession) SubscriptionDeposits(arg0 common.Address) (*big.Int, error) {
	return _SGN.Contract.SubscriptionDeposits(&_SGN.CallOpts, arg0)
}

// SubscriptionDeposits is a free data retrieval call binding the contract method 0x1e77733a.
//
// Solidity: function subscriptionDeposits(address ) constant returns(uint256)
func (_SGN *SGNCallerSession) SubscriptionDeposits(arg0 common.Address) (*big.Int, error) {
	return _SGN.Contract.SubscriptionDeposits(&_SGN.CallOpts, arg0)
}

// RedeemReward is a paid mutator transaction binding the contract method 0x73a6e450.
//
// Solidity: function redeemReward(bytes _rewardRequest) returns()
func (_SGN *SGNTransactor) RedeemReward(opts *bind.TransactOpts, _rewardRequest []byte) (*types.Transaction, error) {
	return _SGN.contract.Transact(opts, "redeemReward", _rewardRequest)
}

// RedeemReward is a paid mutator transaction binding the contract method 0x73a6e450.
//
// Solidity: function redeemReward(bytes _rewardRequest) returns()
func (_SGN *SGNSession) RedeemReward(_rewardRequest []byte) (*types.Transaction, error) {
	return _SGN.Contract.RedeemReward(&_SGN.TransactOpts, _rewardRequest)
}

// RedeemReward is a paid mutator transaction binding the contract method 0x73a6e450.
//
// Solidity: function redeemReward(bytes _rewardRequest) returns()
func (_SGN *SGNTransactorSession) RedeemReward(_rewardRequest []byte) (*types.Transaction, error) {
	return _SGN.Contract.RedeemReward(&_SGN.TransactOpts, _rewardRequest)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(uint256 _amount) returns()
func (_SGN *SGNTransactor) Subscribe(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _SGN.contract.Transact(opts, "subscribe", _amount)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(uint256 _amount) returns()
func (_SGN *SGNSession) Subscribe(_amount *big.Int) (*types.Transaction, error) {
	return _SGN.Contract.Subscribe(&_SGN.TransactOpts, _amount)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(uint256 _amount) returns()
func (_SGN *SGNTransactorSession) Subscribe(_amount *big.Int) (*types.Transaction, error) {
	return _SGN.Contract.Subscribe(&_SGN.TransactOpts, _amount)
}

// UpdateSidechainAddr is a paid mutator transaction binding the contract method 0xe02f39bd.
//
// Solidity: function updateSidechainAddr(bytes _sidechainAddr) returns()
func (_SGN *SGNTransactor) UpdateSidechainAddr(opts *bind.TransactOpts, _sidechainAddr []byte) (*types.Transaction, error) {
	return _SGN.contract.Transact(opts, "updateSidechainAddr", _sidechainAddr)
}

// UpdateSidechainAddr is a paid mutator transaction binding the contract method 0xe02f39bd.
//
// Solidity: function updateSidechainAddr(bytes _sidechainAddr) returns()
func (_SGN *SGNSession) UpdateSidechainAddr(_sidechainAddr []byte) (*types.Transaction, error) {
	return _SGN.Contract.UpdateSidechainAddr(&_SGN.TransactOpts, _sidechainAddr)
}

// UpdateSidechainAddr is a paid mutator transaction binding the contract method 0xe02f39bd.
//
// Solidity: function updateSidechainAddr(bytes _sidechainAddr) returns()
func (_SGN *SGNTransactorSession) UpdateSidechainAddr(_sidechainAddr []byte) (*types.Transaction, error) {
	return _SGN.Contract.UpdateSidechainAddr(&_SGN.TransactOpts, _sidechainAddr)
}

// SGNAddSubscriptionBalanceIterator is returned from FilterAddSubscriptionBalance and is used to iterate over the raw logs and unpacked data for AddSubscriptionBalance events raised by the SGN contract.
type SGNAddSubscriptionBalanceIterator struct {
	Event *SGNAddSubscriptionBalance // Event containing the contract specifics and raw log

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
func (it *SGNAddSubscriptionBalanceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SGNAddSubscriptionBalance)
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
		it.Event = new(SGNAddSubscriptionBalance)
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
func (it *SGNAddSubscriptionBalanceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SGNAddSubscriptionBalanceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SGNAddSubscriptionBalance represents a AddSubscriptionBalance event raised by the SGN contract.
type SGNAddSubscriptionBalance struct {
	Consumer common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddSubscriptionBalance is a free log retrieval operation binding the contract event 0xac095ced75d400384d8404a37883183a56b937b8ff8195fa0c52c3ccc8bb8a49.
//
// Solidity: event AddSubscriptionBalance(address indexed consumer, uint256 amount)
func (_SGN *SGNFilterer) FilterAddSubscriptionBalance(opts *bind.FilterOpts, consumer []common.Address) (*SGNAddSubscriptionBalanceIterator, error) {

	var consumerRule []interface{}
	for _, consumerItem := range consumer {
		consumerRule = append(consumerRule, consumerItem)
	}

	logs, sub, err := _SGN.contract.FilterLogs(opts, "AddSubscriptionBalance", consumerRule)
	if err != nil {
		return nil, err
	}
	return &SGNAddSubscriptionBalanceIterator{contract: _SGN.contract, event: "AddSubscriptionBalance", logs: logs, sub: sub}, nil
}

// WatchAddSubscriptionBalance is a free log subscription operation binding the contract event 0xac095ced75d400384d8404a37883183a56b937b8ff8195fa0c52c3ccc8bb8a49.
//
// Solidity: event AddSubscriptionBalance(address indexed consumer, uint256 amount)
func (_SGN *SGNFilterer) WatchAddSubscriptionBalance(opts *bind.WatchOpts, sink chan<- *SGNAddSubscriptionBalance, consumer []common.Address) (event.Subscription, error) {

	var consumerRule []interface{}
	for _, consumerItem := range consumer {
		consumerRule = append(consumerRule, consumerItem)
	}

	logs, sub, err := _SGN.contract.WatchLogs(opts, "AddSubscriptionBalance", consumerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SGNAddSubscriptionBalance)
				if err := _SGN.contract.UnpackLog(event, "AddSubscriptionBalance", log); err != nil {
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

// ParseAddSubscriptionBalance is a log parse operation binding the contract event 0xac095ced75d400384d8404a37883183a56b937b8ff8195fa0c52c3ccc8bb8a49.
//
// Solidity: event AddSubscriptionBalance(address indexed consumer, uint256 amount)
func (_SGN *SGNFilterer) ParseAddSubscriptionBalance(log types.Log) (*SGNAddSubscriptionBalance, error) {
	event := new(SGNAddSubscriptionBalance)
	if err := _SGN.contract.UnpackLog(event, "AddSubscriptionBalance", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SGNRedeemRewardIterator is returned from FilterRedeemReward and is used to iterate over the raw logs and unpacked data for RedeemReward events raised by the SGN contract.
type SGNRedeemRewardIterator struct {
	Event *SGNRedeemReward // Event containing the contract specifics and raw log

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
func (it *SGNRedeemRewardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SGNRedeemReward)
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
		it.Event = new(SGNRedeemReward)
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
func (it *SGNRedeemRewardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SGNRedeemRewardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SGNRedeemReward represents a RedeemReward event raised by the SGN contract.
type SGNRedeemReward struct {
	Receiver               common.Address
	CumulativeMiningReward *big.Int
	ServiceReward          *big.Int
	ServicePool            *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterRedeemReward is a free log retrieval operation binding the contract event 0x09251621f2e88c5e7f8df91fe1d9e9a70610e20e122945470fddd48af0510426.
//
// Solidity: event RedeemReward(address indexed receiver, uint256 cumulativeMiningReward, uint256 serviceReward, uint256 servicePool)
func (_SGN *SGNFilterer) FilterRedeemReward(opts *bind.FilterOpts, receiver []common.Address) (*SGNRedeemRewardIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _SGN.contract.FilterLogs(opts, "RedeemReward", receiverRule)
	if err != nil {
		return nil, err
	}
	return &SGNRedeemRewardIterator{contract: _SGN.contract, event: "RedeemReward", logs: logs, sub: sub}, nil
}

// WatchRedeemReward is a free log subscription operation binding the contract event 0x09251621f2e88c5e7f8df91fe1d9e9a70610e20e122945470fddd48af0510426.
//
// Solidity: event RedeemReward(address indexed receiver, uint256 cumulativeMiningReward, uint256 serviceReward, uint256 servicePool)
func (_SGN *SGNFilterer) WatchRedeemReward(opts *bind.WatchOpts, sink chan<- *SGNRedeemReward, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _SGN.contract.WatchLogs(opts, "RedeemReward", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SGNRedeemReward)
				if err := _SGN.contract.UnpackLog(event, "RedeemReward", log); err != nil {
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

// ParseRedeemReward is a log parse operation binding the contract event 0x09251621f2e88c5e7f8df91fe1d9e9a70610e20e122945470fddd48af0510426.
//
// Solidity: event RedeemReward(address indexed receiver, uint256 cumulativeMiningReward, uint256 serviceReward, uint256 servicePool)
func (_SGN *SGNFilterer) ParseRedeemReward(log types.Log) (*SGNRedeemReward, error) {
	event := new(SGNRedeemReward)
	if err := _SGN.contract.UnpackLog(event, "RedeemReward", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SGNUpdateSidechainAddrIterator is returned from FilterUpdateSidechainAddr and is used to iterate over the raw logs and unpacked data for UpdateSidechainAddr events raised by the SGN contract.
type SGNUpdateSidechainAddrIterator struct {
	Event *SGNUpdateSidechainAddr // Event containing the contract specifics and raw log

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
func (it *SGNUpdateSidechainAddrIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SGNUpdateSidechainAddr)
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
		it.Event = new(SGNUpdateSidechainAddr)
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
func (it *SGNUpdateSidechainAddrIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SGNUpdateSidechainAddrIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SGNUpdateSidechainAddr represents a UpdateSidechainAddr event raised by the SGN contract.
type SGNUpdateSidechainAddr struct {
	Candidate        common.Address
	OldSidechainAddr common.Hash
	NewSidechainAddr common.Hash
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpdateSidechainAddr is a free log retrieval operation binding the contract event 0x16de3299ab034ce7e21b22d55f4f9a1474bd3c4d20dbd1cc9bcd39c1ad3d5a2c.
//
// Solidity: event UpdateSidechainAddr(address indexed candidate, bytes indexed oldSidechainAddr, bytes indexed newSidechainAddr)
func (_SGN *SGNFilterer) FilterUpdateSidechainAddr(opts *bind.FilterOpts, candidate []common.Address, oldSidechainAddr [][]byte, newSidechainAddr [][]byte) (*SGNUpdateSidechainAddrIterator, error) {

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

	logs, sub, err := _SGN.contract.FilterLogs(opts, "UpdateSidechainAddr", candidateRule, oldSidechainAddrRule, newSidechainAddrRule)
	if err != nil {
		return nil, err
	}
	return &SGNUpdateSidechainAddrIterator{contract: _SGN.contract, event: "UpdateSidechainAddr", logs: logs, sub: sub}, nil
}

// WatchUpdateSidechainAddr is a free log subscription operation binding the contract event 0x16de3299ab034ce7e21b22d55f4f9a1474bd3c4d20dbd1cc9bcd39c1ad3d5a2c.
//
// Solidity: event UpdateSidechainAddr(address indexed candidate, bytes indexed oldSidechainAddr, bytes indexed newSidechainAddr)
func (_SGN *SGNFilterer) WatchUpdateSidechainAddr(opts *bind.WatchOpts, sink chan<- *SGNUpdateSidechainAddr, candidate []common.Address, oldSidechainAddr [][]byte, newSidechainAddr [][]byte) (event.Subscription, error) {

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

	logs, sub, err := _SGN.contract.WatchLogs(opts, "UpdateSidechainAddr", candidateRule, oldSidechainAddrRule, newSidechainAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SGNUpdateSidechainAddr)
				if err := _SGN.contract.UnpackLog(event, "UpdateSidechainAddr", log); err != nil {
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

// ParseUpdateSidechainAddr is a log parse operation binding the contract event 0x16de3299ab034ce7e21b22d55f4f9a1474bd3c4d20dbd1cc9bcd39c1ad3d5a2c.
//
// Solidity: event UpdateSidechainAddr(address indexed candidate, bytes indexed oldSidechainAddr, bytes indexed newSidechainAddr)
func (_SGN *SGNFilterer) ParseUpdateSidechainAddr(log types.Log) (*SGNUpdateSidechainAddr, error) {
	event := new(SGNUpdateSidechainAddr)
	if err := _SGN.contract.UnpackLog(event, "UpdateSidechainAddr", log); err != nil {
		return nil, err
	}
	return event, nil
}
