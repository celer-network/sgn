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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SGNABI is the input ABI used to generate the binding from.
const SGNABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_celerTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_DPoSAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AddSubscriptionBalance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"cumulativeMiningReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"serviceReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"servicePool\",\"type\":\"uint256\"}],\"name\":\"RedeemReward\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"oldSidechainAddr\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"newSidechainAddr\",\"type\":\"bytes\"}],\"name\":\"UpdateSidechainAddr\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addPauser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"celerToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dPoSContract\",\"outputs\":[{\"internalType\":\"contractIDPoS\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isPauser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"redeemedServiceReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renouncePauser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"servicePool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"sidechainAddrMap\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"subscriptionDeposits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"drainToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_sidechainAddr\",\"type\":\"bytes\"}],\"name\":\"updateSidechainAddr\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"subscribe\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_rewardRequest\",\"type\":\"bytes\"}],\"name\":\"redeemReward\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SGNBin is the compiled bytecode used for deploying new contracts.
var SGNBin = "0x60806040523480156200001157600080fd5b50604051620019da380380620019da833981810160405260408110156200003757600080fd5b508051602090910151600080546001600160a01b03191633178082556040516001600160a01b039190911691907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a36200009d336001600160e01b03620000d816565b600280546001600160a81b0319166101006001600160a01b0394851602179055600380546001600160a01b03191691909216179055620001b9565b620000f38160016200012a60201b620015cf1790919060201c565b6040516001600160a01b038216907f6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f890600090a250565b6001600160a01b0381166200013e57600080fd5b6200015382826001600160e01b036200018316565b156200015e57600080fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b60006001600160a01b0382166200019957600080fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b61181180620001c96000396000f3fe608060405234801561001057600080fd5b506004361061012c5760003560e01c806382dc1ec4116100ad578063c6c21e9d11610071578063c6c21e9d146102f7578063e02f39bd146102ff578063e27b41101461036f578063e42a06c81461040a578063f2fde38b146104125761012c565b806382dc1ec4146102935780638456cb59146102b95780638da5cb5b146102c15780638f32d59b146102c9578063c57f6661146102d15761012c565b806346fbf68e116100f457806346fbf68e146101d15780635c975abb1461020b5780636ef8d66d14610213578063715018a61461021b57806373a6e450146102235761012c565b80630f574ba714610131578063145aa116146101505780631e77733a1461016d578063291a3eb2146101a55780633f4ba83a146101c9575b600080fd5b61014e6004803603602081101561014757600080fd5b5035610438565b005b61014e6004803603602081101561016657600080fd5b50356105c4565b6101936004803603602081101561018357600080fd5b50356001600160a01b0316610603565b60408051918252519081900360200190f35b6101ad610615565b604080516001600160a01b039092168252519081900360200190f35b61014e610624565b6101f7600480360360208110156101e757600080fd5b50356001600160a01b0316610684565b604080519115158252519081900360200190f35b6101f761069f565b61014e6106a9565b61014e6106b4565b61014e6004803603602081101561023957600080fd5b81019060208101813564010000000081111561025457600080fd5b82018360208201111561026657600080fd5b8035906020019184600183028401116401000000008311171561028857600080fd5b50909250905061070f565b61014e600480360360208110156102a957600080fd5b50356001600160a01b0316610ab1565b61014e610acc565b6101ad610b30565b6101f7610b3f565b610193600480360360208110156102e757600080fd5b50356001600160a01b0316610b50565b6101ad610b62565b61014e6004803603602081101561031557600080fd5b81019060208101813564010000000081111561033057600080fd5b82018360208201111561034257600080fd5b8035906020019184600183028401116401000000008311171561036457600080fd5b509092509050610b76565b6103956004803603602081101561038557600080fd5b50356001600160a01b0316610e23565b6040805160208082528351818301528351919283929083019185019080838360005b838110156103cf5781810151838201526020016103b7565b50505050905090810190601f1680156103fc5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610193610ebe565b61014e6004803603602081101561042857600080fd5b50356001600160a01b0316610ec4565b60025460ff161561044857600080fd5b600360009054906101000a90046001600160a01b03166001600160a01b031663eab2ed8c6040518163ffffffff1660e01b815260040160206040518083038186803b15801561049657600080fd5b505afa1580156104aa573d6000803e3d6000fd5b505050506040513d60208110156104c057600080fd5b5051610507576040805162461bcd60e51b815260206004820152601160248201527011141bd4c81a5cc81b9bdd081d985b1a59607a1b604482015290519081900360640190fd5b600554339061051c908363ffffffff610ede16565b6005556001600160a01b038116600090815260046020526040902054610548908363ffffffff610ede16565b6001600160a01b03808316600090815260046020526040902091909155600254610581916101009091041682308563ffffffff610ef716565b6040805183815290516001600160a01b038316917fac095ced75d400384d8404a37883183a56b937b8ff8195fa0c52c3ccc8bb8a49919081900360200190a25050565b60025460ff166105d357600080fd5b6105db610b3f565b6105e457600080fd5b6002546106009061010090046001600160a01b03163383610f8d565b50565b60046020526000908152604090205481565b6003546001600160a01b031681565b61062d33610684565b61063657600080fd5b60025460ff1661064557600080fd5b6002805460ff191690556040805133815290517f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa9181900360200190a1565b600061069760018363ffffffff61102716565b90505b919050565b60025460ff165b90565b6106b23361105c565b565b6106bc610b3f565b6106c557600080fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b60025460ff161561071f57600080fd5b600360009054906101000a90046001600160a01b03166001600160a01b031663eab2ed8c6040518163ffffffff1660e01b815260040160206040518083038186803b15801561076d57600080fd5b505afa158015610781573d6000803e3d6000fd5b505050506040513d602081101561079757600080fd5b50516107de576040805162461bcd60e51b815260206004820152601160248201527011141bd4c81a5cc81b9bdd081d985b1a59607a1b604482015290519081900360640190fd5b600354604051631c0efd9d60e01b8152602060048201908152602482018490526001600160a01b0390921691631c0efd9d91859185918190604401848480828437600081840152601f19601f8201169050808301925050509350505050602060405180830381600087803b15801561085557600080fd5b505af1158015610869573d6000803e3d6000fd5b505050506040513d602081101561087f57600080fd5b50516108bc5760405162461bcd60e51b815260040180806020018281038252602281526020018061178c6022913960400191505060405180910390fd5b6108c4611695565b61090383838080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506110a492505050565b905061090d6116af565b8151610918906111fb565b80516001600160a01b03166000908152600660205260408082205490830151929350909161094b9163ffffffff6112c016565b905080600554101561098e5760405162461bcd60e51b815260040180806020018281038252602f8152602001806117ae602f913960400191505060405180910390fd5b60408083015183516001600160a01b03166000908152600660205291909120556005546109c1908263ffffffff6112c016565b6005556003548251602084015160408051630fbd844360e11b81526001600160a01b0393841660048201526024810192909252519190921691631f7b088691604480830192600092919082900301818387803b158015610a2057600080fd5b505af1158015610a34573d6000803e3d6000fd5b50508351600254610a56935061010090046001600160a01b0316915083610f8d565b8151602080840151600554604080519283529282018590528183015290516001600160a01b03909216917f09251621f2e88c5e7f8df91fe1d9e9a70610e20e122945470fddd48af05104269181900360600190a25050505050565b610aba33610684565b610ac357600080fd5b610600816112d5565b610ad533610684565b610ade57600080fd5b60025460ff1615610aee57600080fd5b6002805460ff191660011790556040805133815290517f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2589181900360200190a1565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b60066020526000908152604090205481565b60025461010090046001600160a01b031681565b600354604080516328bde1e160e01b815233600482018190529151919260009283926001600160a01b03909216916328bde1e19160248083019260e0929190829003018186803b158015610bc957600080fd5b505afa158015610bdd573d6000803e3d6000fd5b505050506040513d60e0811015610bf357600080fd5b50805160609091015190925090508015610c54576040805162461bcd60e51b815260206004820152601a60248201527f6d73672e73656e646572206973206e6f7420756e626f6e646564000000000000604482015290519081900360640190fd5b81610ca6576040805162461bcd60e51b815260206004820152601c60248201527f43616e646964617465206973206e6f7420696e697469616c697a656400000000604482015290519081900360640190fd5b6001600160a01b03831660009081526007602090815260409182902080548351601f6002600019610100600186161502019093169290920491820184900484028101840190945280845260609392830182828015610d455780601f10610d1a57610100808354040283529160200191610d45565b820191906000526020600020905b815481529060010190602001808311610d2857829003601f168201915b505050506001600160a01b0386166000908152600760205260409020919250610d7191905087876116d9565b5085856040518083838082843760405192018290038220865190955086945091925082916020850191508083835b60208310610dbe5780518252601f199092019160209182019101610d9f565b5181516020939093036101000a6000190180199091169216919091179052604051920182900382209350506001600160a01b03881691507f16de3299ab034ce7e21b22d55f4f9a1474bd3c4d20dbd1cc9bcd39c1ad3d5a2c90600090a4505050505050565b60076020908152600091825260409182902080548351601f600260001961010060018616150201909316929092049182018490048402810184019094528084529091830182828015610eb65780601f10610e8b57610100808354040283529160200191610eb6565b820191906000526020600020905b815481529060010190602001808311610e9957829003601f168201915b505050505081565b60055481565b610ecc610b3f565b610ed557600080fd5b6106008161131d565b600082820183811015610ef057600080fd5b9392505050565b604080516323b872dd60e01b81526001600160a01b0385811660048301528481166024830152604482018490529151918616916323b872dd916064808201926020929091908290030181600087803b158015610f5257600080fd5b505af1158015610f66573d6000803e3d6000fd5b505050506040513d6020811015610f7c57600080fd5b5051610f8757600080fd5b50505050565b826001600160a01b031663a9059cbb83836040518363ffffffff1660e01b815260040180836001600160a01b03166001600160a01b0316815260200182815260200192505050602060405180830381600087803b158015610fed57600080fd5b505af1158015611001573d6000803e3d6000fd5b505050506040513d602081101561101757600080fd5b505161102257600080fd5b505050565b60006001600160a01b03821661103c57600080fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b61106d60018263ffffffff61138b16565b6040516001600160a01b038216907fcd265ebaf09df2871cc7bd4133404a235ba12eff2041bb89d9c714a2621c7c7e90600090a250565b6110ac611695565b6110b4611757565b6110bd836113d3565b905060606110d282600263ffffffff6113ea16565b9050806002815181106110e157fe5b602002602001015160405190808252806020026020018201604052801561111c57816020015b60608152602001906001900390816111075790505b50836020018190525060008160028151811061113457fe5b6020026020010181815250506000805b61114d8461147a565b156111f25761115b84611486565b9092509050816001141561117957611172846114b3565b85526111ed565b81600214156111dd5761118b846114b3565b85602001518460028151811061119d57fe5b6020026020010151815181106111af57fe5b6020026020010181905250826002815181106111c757fe5b60209081029190910101805160010190526111ed565b6111ed848263ffffffff61154016565b611144565b50505050919050565b6112036116af565b61120b611757565b611214836113d3565b90506000805b6112238361147a565b156112b85761123183611486565b909250905081600114156112605761125061124b846114b3565b6115a1565b6001600160a01b031684526112b3565b81600214156112845761127a611275846114b3565b6115ac565b60208501526112b3565b81600314156112a357611299611275846114b3565b60408501526112b3565b6112b3838263ffffffff61154016565b61121a565b505050919050565b6000828211156112cf57600080fd5b50900390565b6112e660018263ffffffff6115cf16565b6040516001600160a01b038216907f6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f890600090a250565b6001600160a01b03811661133057600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b03811661139e57600080fd5b6113a88282611027565b6113b157600080fd5b6001600160a01b0316600090815260209190915260409020805460ff19169055565b6113db611757565b60208101919091526000815290565b81516040805160018401808252602080820283010190925260609291801561141c578160200160208202803883390190505b5091506000805b61142c8661147a565b156114715761143a86611486565b8092508193505050600184838151811061145057fe5b60200260200101818151019150818152505061146c8682611540565b611423565b50509092525090565b60208101515190511090565b60008060006114948461161b565b90506008810492508060071660058111156114ab57fe5b915050915091565b606060006114c08361161b565b83516020850151519192508201908111156114da57600080fd5b816040519080825280601f01601f191660200182016040528015611505576020820181803883390190505b50602080860151865192955091818601919083010160005b8581101561153557818101518382015260200161151d565b505050935250919050565b600081600581111561154e57fe5b14156115635761155d8261161b565b5061159d565b600281600581111561157157fe5b141561012c5760006115828361161b565b83518101808552602085015151919250111561155d57600080fd5b5050565b600061069782611676565b60006020825111156115bd57600080fd5b50602081810151915160089103021c90565b6001600160a01b0381166115e257600080fd5b6115ec8282611027565b156115f657600080fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b602080820151825181019091015160009182805b600a8110156116705783811a91508060070282607f16901b8517945081608016600014156116685785510160010185525061069a915050565b60010161162f565b50600080fd5b6000815160141461168657600080fd5b5060200151600160601b900490565b604051806040016040528060608152602001606081525090565b604051806060016040528060006001600160a01b0316815260200160008152602001600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061171a5782800160ff19823516178555611747565b82800160010185558215611747579182015b8281111561174757823582559160200191906001019061172c565b50611753929150611771565b5090565b604051806040016040528060008152602001606081525090565b6106a691905b80821115611753576000815560010161177756fe56616c696461746f72207369677320766572696669636174696f6e206661696c65645365727669636520706f6f6c20697320736d616c6c6572207468616e206e6577207365727669636520726577617264a265627a7a723158200ef48f0fa8c2c39f174999979f900f8eed58c0a95e79562343d0db5c5f94678a64736f6c63430005110032"

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

// CelerToken is a free data retrieval call binding the contract method 0xc6c21e9d.
//
// Solidity: function celerToken() view returns(address)
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
// Solidity: function celerToken() view returns(address)
func (_SGN *SGNSession) CelerToken() (common.Address, error) {
	return _SGN.Contract.CelerToken(&_SGN.CallOpts)
}

// CelerToken is a free data retrieval call binding the contract method 0xc6c21e9d.
//
// Solidity: function celerToken() view returns(address)
func (_SGN *SGNCallerSession) CelerToken() (common.Address, error) {
	return _SGN.Contract.CelerToken(&_SGN.CallOpts)
}

// DPoSContract is a free data retrieval call binding the contract method 0x291a3eb2.
//
// Solidity: function dPoSContract() view returns(address)
func (_SGN *SGNCaller) DPoSContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SGN.contract.Call(opts, out, "dPoSContract")
	return *ret0, err
}

// DPoSContract is a free data retrieval call binding the contract method 0x291a3eb2.
//
// Solidity: function dPoSContract() view returns(address)
func (_SGN *SGNSession) DPoSContract() (common.Address, error) {
	return _SGN.Contract.DPoSContract(&_SGN.CallOpts)
}

// DPoSContract is a free data retrieval call binding the contract method 0x291a3eb2.
//
// Solidity: function dPoSContract() view returns(address)
func (_SGN *SGNCallerSession) DPoSContract() (common.Address, error) {
	return _SGN.Contract.DPoSContract(&_SGN.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_SGN *SGNCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SGN.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_SGN *SGNSession) IsOwner() (bool, error) {
	return _SGN.Contract.IsOwner(&_SGN.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_SGN *SGNCallerSession) IsOwner() (bool, error) {
	return _SGN.Contract.IsOwner(&_SGN.CallOpts)
}

// IsPauser is a free data retrieval call binding the contract method 0x46fbf68e.
//
// Solidity: function isPauser(address account) view returns(bool)
func (_SGN *SGNCaller) IsPauser(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SGN.contract.Call(opts, out, "isPauser", account)
	return *ret0, err
}

// IsPauser is a free data retrieval call binding the contract method 0x46fbf68e.
//
// Solidity: function isPauser(address account) view returns(bool)
func (_SGN *SGNSession) IsPauser(account common.Address) (bool, error) {
	return _SGN.Contract.IsPauser(&_SGN.CallOpts, account)
}

// IsPauser is a free data retrieval call binding the contract method 0x46fbf68e.
//
// Solidity: function isPauser(address account) view returns(bool)
func (_SGN *SGNCallerSession) IsPauser(account common.Address) (bool, error) {
	return _SGN.Contract.IsPauser(&_SGN.CallOpts, account)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SGN *SGNCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SGN.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SGN *SGNSession) Owner() (common.Address, error) {
	return _SGN.Contract.Owner(&_SGN.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SGN *SGNCallerSession) Owner() (common.Address, error) {
	return _SGN.Contract.Owner(&_SGN.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SGN *SGNCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SGN.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SGN *SGNSession) Paused() (bool, error) {
	return _SGN.Contract.Paused(&_SGN.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SGN *SGNCallerSession) Paused() (bool, error) {
	return _SGN.Contract.Paused(&_SGN.CallOpts)
}

// RedeemedServiceReward is a free data retrieval call binding the contract method 0xc57f6661.
//
// Solidity: function redeemedServiceReward(address ) view returns(uint256)
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
// Solidity: function redeemedServiceReward(address ) view returns(uint256)
func (_SGN *SGNSession) RedeemedServiceReward(arg0 common.Address) (*big.Int, error) {
	return _SGN.Contract.RedeemedServiceReward(&_SGN.CallOpts, arg0)
}

// RedeemedServiceReward is a free data retrieval call binding the contract method 0xc57f6661.
//
// Solidity: function redeemedServiceReward(address ) view returns(uint256)
func (_SGN *SGNCallerSession) RedeemedServiceReward(arg0 common.Address) (*big.Int, error) {
	return _SGN.Contract.RedeemedServiceReward(&_SGN.CallOpts, arg0)
}

// ServicePool is a free data retrieval call binding the contract method 0xe42a06c8.
//
// Solidity: function servicePool() view returns(uint256)
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
// Solidity: function servicePool() view returns(uint256)
func (_SGN *SGNSession) ServicePool() (*big.Int, error) {
	return _SGN.Contract.ServicePool(&_SGN.CallOpts)
}

// ServicePool is a free data retrieval call binding the contract method 0xe42a06c8.
//
// Solidity: function servicePool() view returns(uint256)
func (_SGN *SGNCallerSession) ServicePool() (*big.Int, error) {
	return _SGN.Contract.ServicePool(&_SGN.CallOpts)
}

// SidechainAddrMap is a free data retrieval call binding the contract method 0xe27b4110.
//
// Solidity: function sidechainAddrMap(address ) view returns(bytes)
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
// Solidity: function sidechainAddrMap(address ) view returns(bytes)
func (_SGN *SGNSession) SidechainAddrMap(arg0 common.Address) ([]byte, error) {
	return _SGN.Contract.SidechainAddrMap(&_SGN.CallOpts, arg0)
}

// SidechainAddrMap is a free data retrieval call binding the contract method 0xe27b4110.
//
// Solidity: function sidechainAddrMap(address ) view returns(bytes)
func (_SGN *SGNCallerSession) SidechainAddrMap(arg0 common.Address) ([]byte, error) {
	return _SGN.Contract.SidechainAddrMap(&_SGN.CallOpts, arg0)
}

// SubscriptionDeposits is a free data retrieval call binding the contract method 0x1e77733a.
//
// Solidity: function subscriptionDeposits(address ) view returns(uint256)
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
// Solidity: function subscriptionDeposits(address ) view returns(uint256)
func (_SGN *SGNSession) SubscriptionDeposits(arg0 common.Address) (*big.Int, error) {
	return _SGN.Contract.SubscriptionDeposits(&_SGN.CallOpts, arg0)
}

// SubscriptionDeposits is a free data retrieval call binding the contract method 0x1e77733a.
//
// Solidity: function subscriptionDeposits(address ) view returns(uint256)
func (_SGN *SGNCallerSession) SubscriptionDeposits(arg0 common.Address) (*big.Int, error) {
	return _SGN.Contract.SubscriptionDeposits(&_SGN.CallOpts, arg0)
}

// AddPauser is a paid mutator transaction binding the contract method 0x82dc1ec4.
//
// Solidity: function addPauser(address account) returns()
func (_SGN *SGNTransactor) AddPauser(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _SGN.contract.Transact(opts, "addPauser", account)
}

// AddPauser is a paid mutator transaction binding the contract method 0x82dc1ec4.
//
// Solidity: function addPauser(address account) returns()
func (_SGN *SGNSession) AddPauser(account common.Address) (*types.Transaction, error) {
	return _SGN.Contract.AddPauser(&_SGN.TransactOpts, account)
}

// AddPauser is a paid mutator transaction binding the contract method 0x82dc1ec4.
//
// Solidity: function addPauser(address account) returns()
func (_SGN *SGNTransactorSession) AddPauser(account common.Address) (*types.Transaction, error) {
	return _SGN.Contract.AddPauser(&_SGN.TransactOpts, account)
}

// DrainToken is a paid mutator transaction binding the contract method 0x145aa116.
//
// Solidity: function drainToken(uint256 _amount) returns()
func (_SGN *SGNTransactor) DrainToken(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _SGN.contract.Transact(opts, "drainToken", _amount)
}

// DrainToken is a paid mutator transaction binding the contract method 0x145aa116.
//
// Solidity: function drainToken(uint256 _amount) returns()
func (_SGN *SGNSession) DrainToken(_amount *big.Int) (*types.Transaction, error) {
	return _SGN.Contract.DrainToken(&_SGN.TransactOpts, _amount)
}

// DrainToken is a paid mutator transaction binding the contract method 0x145aa116.
//
// Solidity: function drainToken(uint256 _amount) returns()
func (_SGN *SGNTransactorSession) DrainToken(_amount *big.Int) (*types.Transaction, error) {
	return _SGN.Contract.DrainToken(&_SGN.TransactOpts, _amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SGN *SGNTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SGN.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SGN *SGNSession) Pause() (*types.Transaction, error) {
	return _SGN.Contract.Pause(&_SGN.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SGN *SGNTransactorSession) Pause() (*types.Transaction, error) {
	return _SGN.Contract.Pause(&_SGN.TransactOpts)
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

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SGN *SGNTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SGN.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SGN *SGNSession) RenounceOwnership() (*types.Transaction, error) {
	return _SGN.Contract.RenounceOwnership(&_SGN.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SGN *SGNTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SGN.Contract.RenounceOwnership(&_SGN.TransactOpts)
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x6ef8d66d.
//
// Solidity: function renouncePauser() returns()
func (_SGN *SGNTransactor) RenouncePauser(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SGN.contract.Transact(opts, "renouncePauser")
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x6ef8d66d.
//
// Solidity: function renouncePauser() returns()
func (_SGN *SGNSession) RenouncePauser() (*types.Transaction, error) {
	return _SGN.Contract.RenouncePauser(&_SGN.TransactOpts)
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x6ef8d66d.
//
// Solidity: function renouncePauser() returns()
func (_SGN *SGNTransactorSession) RenouncePauser() (*types.Transaction, error) {
	return _SGN.Contract.RenouncePauser(&_SGN.TransactOpts)
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

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SGN *SGNTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SGN.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SGN *SGNSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SGN.Contract.TransferOwnership(&_SGN.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SGN *SGNTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SGN.Contract.TransferOwnership(&_SGN.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SGN *SGNTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SGN.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SGN *SGNSession) Unpause() (*types.Transaction, error) {
	return _SGN.Contract.Unpause(&_SGN.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SGN *SGNTransactorSession) Unpause() (*types.Transaction, error) {
	return _SGN.Contract.Unpause(&_SGN.TransactOpts)
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

// SGNOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SGN contract.
type SGNOwnershipTransferredIterator struct {
	Event *SGNOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SGNOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SGNOwnershipTransferred)
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
		it.Event = new(SGNOwnershipTransferred)
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
func (it *SGNOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SGNOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SGNOwnershipTransferred represents a OwnershipTransferred event raised by the SGN contract.
type SGNOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SGN *SGNFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SGNOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SGN.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SGNOwnershipTransferredIterator{contract: _SGN.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SGN *SGNFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SGNOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SGN.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SGNOwnershipTransferred)
				if err := _SGN.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SGN *SGNFilterer) ParseOwnershipTransferred(log types.Log) (*SGNOwnershipTransferred, error) {
	event := new(SGNOwnershipTransferred)
	if err := _SGN.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SGNPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the SGN contract.
type SGNPausedIterator struct {
	Event *SGNPaused // Event containing the contract specifics and raw log

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
func (it *SGNPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SGNPaused)
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
		it.Event = new(SGNPaused)
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
func (it *SGNPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SGNPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SGNPaused represents a Paused event raised by the SGN contract.
type SGNPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SGN *SGNFilterer) FilterPaused(opts *bind.FilterOpts) (*SGNPausedIterator, error) {

	logs, sub, err := _SGN.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &SGNPausedIterator{contract: _SGN.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SGN *SGNFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *SGNPaused) (event.Subscription, error) {

	logs, sub, err := _SGN.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SGNPaused)
				if err := _SGN.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SGN *SGNFilterer) ParsePaused(log types.Log) (*SGNPaused, error) {
	event := new(SGNPaused)
	if err := _SGN.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SGNPauserAddedIterator is returned from FilterPauserAdded and is used to iterate over the raw logs and unpacked data for PauserAdded events raised by the SGN contract.
type SGNPauserAddedIterator struct {
	Event *SGNPauserAdded // Event containing the contract specifics and raw log

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
func (it *SGNPauserAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SGNPauserAdded)
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
		it.Event = new(SGNPauserAdded)
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
func (it *SGNPauserAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SGNPauserAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SGNPauserAdded represents a PauserAdded event raised by the SGN contract.
type SGNPauserAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPauserAdded is a free log retrieval operation binding the contract event 0x6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f8.
//
// Solidity: event PauserAdded(address indexed account)
func (_SGN *SGNFilterer) FilterPauserAdded(opts *bind.FilterOpts, account []common.Address) (*SGNPauserAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SGN.contract.FilterLogs(opts, "PauserAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &SGNPauserAddedIterator{contract: _SGN.contract, event: "PauserAdded", logs: logs, sub: sub}, nil
}

// WatchPauserAdded is a free log subscription operation binding the contract event 0x6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f8.
//
// Solidity: event PauserAdded(address indexed account)
func (_SGN *SGNFilterer) WatchPauserAdded(opts *bind.WatchOpts, sink chan<- *SGNPauserAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SGN.contract.WatchLogs(opts, "PauserAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SGNPauserAdded)
				if err := _SGN.contract.UnpackLog(event, "PauserAdded", log); err != nil {
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

// ParsePauserAdded is a log parse operation binding the contract event 0x6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f8.
//
// Solidity: event PauserAdded(address indexed account)
func (_SGN *SGNFilterer) ParsePauserAdded(log types.Log) (*SGNPauserAdded, error) {
	event := new(SGNPauserAdded)
	if err := _SGN.contract.UnpackLog(event, "PauserAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SGNPauserRemovedIterator is returned from FilterPauserRemoved and is used to iterate over the raw logs and unpacked data for PauserRemoved events raised by the SGN contract.
type SGNPauserRemovedIterator struct {
	Event *SGNPauserRemoved // Event containing the contract specifics and raw log

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
func (it *SGNPauserRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SGNPauserRemoved)
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
		it.Event = new(SGNPauserRemoved)
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
func (it *SGNPauserRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SGNPauserRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SGNPauserRemoved represents a PauserRemoved event raised by the SGN contract.
type SGNPauserRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPauserRemoved is a free log retrieval operation binding the contract event 0xcd265ebaf09df2871cc7bd4133404a235ba12eff2041bb89d9c714a2621c7c7e.
//
// Solidity: event PauserRemoved(address indexed account)
func (_SGN *SGNFilterer) FilterPauserRemoved(opts *bind.FilterOpts, account []common.Address) (*SGNPauserRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SGN.contract.FilterLogs(opts, "PauserRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &SGNPauserRemovedIterator{contract: _SGN.contract, event: "PauserRemoved", logs: logs, sub: sub}, nil
}

// WatchPauserRemoved is a free log subscription operation binding the contract event 0xcd265ebaf09df2871cc7bd4133404a235ba12eff2041bb89d9c714a2621c7c7e.
//
// Solidity: event PauserRemoved(address indexed account)
func (_SGN *SGNFilterer) WatchPauserRemoved(opts *bind.WatchOpts, sink chan<- *SGNPauserRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SGN.contract.WatchLogs(opts, "PauserRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SGNPauserRemoved)
				if err := _SGN.contract.UnpackLog(event, "PauserRemoved", log); err != nil {
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

// ParsePauserRemoved is a log parse operation binding the contract event 0xcd265ebaf09df2871cc7bd4133404a235ba12eff2041bb89d9c714a2621c7c7e.
//
// Solidity: event PauserRemoved(address indexed account)
func (_SGN *SGNFilterer) ParsePauserRemoved(log types.Log) (*SGNPauserRemoved, error) {
	event := new(SGNPauserRemoved)
	if err := _SGN.contract.UnpackLog(event, "PauserRemoved", log); err != nil {
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

// SGNUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the SGN contract.
type SGNUnpausedIterator struct {
	Event *SGNUnpaused // Event containing the contract specifics and raw log

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
func (it *SGNUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SGNUnpaused)
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
		it.Event = new(SGNUnpaused)
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
func (it *SGNUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SGNUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SGNUnpaused represents a Unpaused event raised by the SGN contract.
type SGNUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SGN *SGNFilterer) FilterUnpaused(opts *bind.FilterOpts) (*SGNUnpausedIterator, error) {

	logs, sub, err := _SGN.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &SGNUnpausedIterator{contract: _SGN.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SGN *SGNFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *SGNUnpaused) (event.Subscription, error) {

	logs, sub, err := _SGN.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SGNUnpaused)
				if err := _SGN.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SGN *SGNFilterer) ParseUnpaused(log types.Log) (*SGNUnpaused, error) {
	event := new(SGNUnpaused)
	if err := _SGN.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
