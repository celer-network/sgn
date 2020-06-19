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
var SGNBin = "0x608060405234801561001057600080fd5b5060405161137b38038061137b8339818101604052604081101561003357600080fd5b50805160209091015160008054600160a060020a03938416600160a060020a031991821617909155600180549390921692169190911790556113018061007a6000396000f3fe608060405234801561001057600080fd5b50600436106100975760003560e060020a90048063c57f66611161006a578063c57f666114610187578063c6c21e9d146101ad578063e02f39bd146101b5578063e27b411014610225578063e42a06c8146102c057610097565b80630f574ba71461009c5780631e77733a146100bb57806371273548146100f357806373a6e45014610117575b600080fd5b6100b9600480360360208110156100b257600080fd5b50356102c8565b005b6100e1600480360360208110156100d157600080fd5b5035600160a060020a031661044f565b60408051918252519081900360200190f35b6100fb610461565b60408051600160a060020a039092168252519081900360200190f35b6100b96004803603602081101561012d57600080fd5b81019060208101813564010000000081111561014857600080fd5b82018360208201111561015a57600080fd5b8035906020019184600183028401116401000000008311171561017c57600080fd5b509092509050610470565b6100e16004803603602081101561019d57600080fd5b5035600160a060020a0316610820565b6100fb610832565b6100b9600480360360208110156101cb57600080fd5b8101906020810181356401000000008111156101e657600080fd5b8201836020820111156101f857600080fd5b8035906020019184600183028401116401000000008311171561021a57600080fd5b509092509050610841565b61024b6004803603602081101561023b57600080fd5b5035600160a060020a0316610b0d565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561028557818101518382015260200161026d565b50505050905090810190601f1680156102b25780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100e1610ba8565b600160009054906101000a9004600160a060020a0316600160a060020a031663eab2ed8c6040518163ffffffff1660e060020a02815260040160206040518083038186803b15801561031957600080fd5b505afa15801561032d573d6000803e3d6000fd5b505050506040513d602081101561034357600080fd5b5051610399576040805160e560020a62461bcd02815260206004820152601160248201527f44506f53206973206e6f742076616c6964000000000000000000000000000000604482015290519081900360640190fd5b60035433906103ae908363ffffffff610bae16565b600355600160a060020a0381166000908152600260205260409020546103da908363ffffffff610bae16565b600160a060020a03808316600090815260026020526040812092909255905461040c911682308563ffffffff610bc716565b604080518381529051600160a060020a038316917fac095ced75d400384d8404a37883183a56b937b8ff8195fa0c52c3ccc8bb8a49919081900360200190a25050565b60026020526000908152604090205481565b600154600160a060020a031681565b600160009054906101000a9004600160a060020a0316600160a060020a031663eab2ed8c6040518163ffffffff1660e060020a02815260040160206040518083038186803b1580156104c157600080fd5b505afa1580156104d5573d6000803e3d6000fd5b505050506040513d60208110156104eb57600080fd5b5051610541576040805160e560020a62461bcd02815260206004820152601160248201527f44506f53206973206e6f742076616c6964000000000000000000000000000000604482015290519081900360640190fd5b6001546040517f1c0efd9d00000000000000000000000000000000000000000000000000000000815260206004820190815260248201849052600160a060020a0390921691631c0efd9d91859185918190604401848480828437600081840152601f19601f8201169050808301925050509350505050602060405180830381600087803b1580156105d157600080fd5b505af11580156105e5573d6000803e3d6000fd5b505050506040513d60208110156105fb57600080fd5b5051610651576040805160e560020a62461bcd02815260206004820152601c60248201527f4661696c20746f20636865636b2076616c696461746f72207369677300000000604482015290519081900360640190fd5b6106596111d3565b61069883838080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610c7692505050565b90506106a26111ed565b81516106ad90610dcd565b8051600160a060020a0316600090815260046020526040808220549083015192935090916106e09163ffffffff610e9216565b6040808401518451600160a060020a0316600090815260046020529190912055600354909150610716908263ffffffff610e9216565b60035560015482516020840151604080517f1f7b0886000000000000000000000000000000000000000000000000000000008152600160a060020a0393841660048201526024810192909252519190921691631f7b088691604480830192600092919082900301818387803b15801561078e57600080fd5b505af11580156107a2573d6000803e3d6000fd5b505083516000546107c59350600160a060020a031691508363ffffffff610ea716565b815160208084015160035460408051928352928201859052818301529051600160a060020a03909216917f09251621f2e88c5e7f8df91fe1d9e9a70610e20e122945470fddd48af05104269181900360600190a25050505050565b60046020526000908152604090205481565b600054600160a060020a031681565b600154604080517f28bde1e10000000000000000000000000000000000000000000000000000000081523360048201819052915191926000928392600160a060020a03909216916328bde1e19160248083019260e0929190829003018186803b1580156108ad57600080fd5b505afa1580156108c1573d6000803e3d6000fd5b505050506040513d60e08110156108d757600080fd5b5080516060909101519092509050801561093b576040805160e560020a62461bcd02815260206004820152601a60248201527f6d73672e73656e646572206973206e6f7420756e626f6e646564000000000000604482015290519081900360640190fd5b81610990576040805160e560020a62461bcd02815260206004820152601c60248201527f43616e646964617465206973206e6f7420696e697469616c697a656400000000604482015290519081900360640190fd5b600160a060020a03831660009081526005602090815260409182902080548351601f6002600019610100600186161502019093169290920491820184900484028101840190945280845260609392830182828015610a2f5780601f10610a0457610100808354040283529160200191610a2f565b820191906000526020600020905b815481529060010190602001808311610a1257829003601f168201915b50505050600160a060020a0386166000908152600560205260409020919250610a5b9190508787611217565b5085856040518083838082843760405192018290038220865190955086945091925082916020850191508083835b60208310610aa85780518252601f199092019160209182019101610a89565b5181516020939093036101000a600019018019909116921691909117905260405192018290038220935050600160a060020a03881691507f16de3299ab034ce7e21b22d55f4f9a1474bd3c4d20dbd1cc9bcd39c1ad3d5a2c90600090a4505050505050565b60056020908152600091825260409182902080548351601f600260001961010060018616150201909316929092049182018490048402810184019094528084529091830182828015610ba05780601f10610b7557610100808354040283529160200191610ba0565b820191906000526020600020905b815481529060010190602001808311610b8357829003601f168201915b505050505081565b60035481565b600082820183811015610bc057600080fd5b9392505050565b604080517f23b872dd000000000000000000000000000000000000000000000000000000008152600160a060020a0385811660048301528481166024830152604482018490529151918616916323b872dd916064808201926020929091908290030181600087803b158015610c3b57600080fd5b505af1158015610c4f573d6000803e3d6000fd5b505050506040513d6020811015610c6557600080fd5b5051610c7057600080fd5b50505050565b610c7e6111d3565b610c86611295565b610c8f83610f44565b90506060610ca482600263ffffffff610f5b16565b905080600281518110610cb357fe5b6020026020010151604051908082528060200260200182016040528015610cee57816020015b6060815260200190600190039081610cd95790505b508360200181905250600081600281518110610d0657fe5b6020026020010181815250506000805b610d1f84610feb565b15610dc457610d2d84610ffa565b90925090508160011415610d4b57610d4484611027565b8552610dbf565b8160021415610daf57610d5d84611027565b856020015184600281518110610d6f57fe5b602002602001015181518110610d8157fe5b602002602001018190525082600281518110610d9957fe5b6020908102919091010180516001019052610dbf565b610dbf848263ffffffff6110b416565b610d16565b50505050919050565b610dd56111ed565b610ddd611295565b610de683610f44565b90506000805b610df583610feb565b15610e8a57610e0383610ffa565b90925090508160011415610e3257610e22610e1d84611027565b611115565b600160a060020a03168452610e85565b8160021415610e5657610e4c610e4784611027565b611126565b6020850152610e85565b8160031415610e7557610e6b610e4784611027565b6040850152610e85565b610e85838263ffffffff6110b416565b610dec565b505050919050565b600082821115610ea157600080fd5b50900390565b82600160a060020a031663a9059cbb83836040518363ffffffff1660e060020a0281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015610f0a57600080fd5b505af1158015610f1e573d6000803e3d6000fd5b505050506040513d6020811015610f3457600080fd5b5051610f3f57600080fd5b505050565b610f4c611295565b60208101919091526000815290565b815160408051600184018082526020808202830101909252606092918015610f8d578160200160208202803883390190505b5091506000805b610f9d86610feb565b15610fe257610fab86610ffa565b80925081935050506001848381518110610fc157fe5b602002602001018181510191508181525050610fdd86826110b4565b610f94565b50509092525090565b6020810151518151105b919050565b60008060006110088461114d565b905060088104925080600716600581111561101f57fe5b915050915091565b606060006110348361114d565b835160208501515191925082019081111561104e57600080fd5b816040519080825280601f01601f191660200182016040528015611079576020820181803883390190505b50602080860151865192955091818601919083010160005b858110156110a9578181015183820152602001611091565b505050935250919050565b60008160058111156110c257fe5b14156110d7576110d18261114d565b50611111565b60028160058111156110e557fe5b14156100975760006110f68361114d565b8351810180855260208501515191925011156110d157600080fd5b5050565b6000611120826111ab565b92915050565b600060208251111561113757600080fd5b506020818101519151600891030260020a900490565b602080820151825181019091015160009182805b600a8110156111a55783811a91508060070282607f169060020a0285179450816080166000141561119d57855101600101855250610ff5915050565b600101611161565b50600080fd5b600081516014146111bb57600080fd5b50602001516c01000000000000000000000000900490565b604051806040016040528060608152602001606081525090565b60405180606001604052806000600160a060020a0316815260200160008152602001600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106112585782800160ff19823516178555611285565b82800160010185558215611285579182015b8281111561128557823582559160200191906001019061126a565b506112919291506112af565b5090565b604051806040016040528060008152602001606081525090565b6112c991905b8082111561129157600081556001016112b5565b9056fea265627a7a72305820774fd440059e513aa92aec895ce9d2ec408b3742890f6e936e6dca78357d3fef64736f6c634300050a0032"

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
