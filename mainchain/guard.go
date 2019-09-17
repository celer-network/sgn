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
const GuardABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"blameTimeout\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minTotalStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VALIDATOR_SET_MAX_SIZE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"miningPool\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"sidechainGoLiveTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"celerToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"subscriptionFees\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validatorSet\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minValidatorNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_celerTokenAddress\",\"type\":\"address\"},{\"name\":\"_blameTimeout\",\"type\":\"uint256\"},{\"name\":\"_minValidatorNum\",\"type\":\"uint256\"},{\"name\":\"_minTotalStake\",\"type\":\"uint256\"},{\"name\":\"_sidechainGoLiveTime\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"minSelfStake\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"sidechainAddr\",\"type\":\"bytes\"}],\"name\":\"InitializeCandidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"newStake\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"totalStake\",\"type\":\"uint256\"}],\"name\":\"Delegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"oldSidechainAddr\",\"type\":\"bytes\"},{\"indexed\":true,\"name\":\"newSidechainAddr\",\"type\":\"bytes\"}],\"name\":\"UpdateSidechainAddr\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"ethAddr\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"changeType\",\"type\":\"uint8\"}],\"name\":\"ValidatorChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawFromUnbondedCandidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"withdrawAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"totalStake\",\"type\":\"uint256\"}],\"name\":\"IntendWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ConfirmWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"indemnitor\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"indemnitee\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Punish\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"consumer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AddSubscriptionBalance\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_minSelfStake\",\"type\":\"uint256\"},{\"name\":\"_sidechainAddr\",\"type\":\"bytes\"}],\"name\":\"initializeCandidate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"delegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_sidechainAddr\",\"type\":\"bytes\"}],\"name\":\"updateSidechainAddr\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFromUnbondedCandidate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"intendWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"}],\"name\":\"confirmWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"subscribe\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getValidatorNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMinStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"}],\"name\":\"getCandidateInfo\",\"outputs\":[{\"name\":\"initialized\",\"type\":\"bool\"},{\"name\":\"minSelfStake\",\"type\":\"uint256\"},{\"name\":\"sidechainAddr\",\"type\":\"bytes\"},{\"name\":\"totalStake\",\"type\":\"uint256\"},{\"name\":\"isVldt\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"},{\"name\":\"_delegatorAddr\",\"type\":\"address\"}],\"name\":\"getDelegatorInfo\",\"outputs\":[{\"name\":\"stake\",\"type\":\"uint256\"},{\"name\":\"intentAmounts\",\"type\":\"uint256[]\"},{\"name\":\"intentUnlockTimes\",\"type\":\"uint256[]\"},{\"name\":\"nextWithdrawIntent\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// GuardBin is the compiled bytecode used for deploying new contracts.
const GuardBin = `0x608060405234801561001057600080fd5b50604051611c19380380611c19833981810160405260a081101561003357600080fd5b50805160208201516040830151606084015160809094015160008054600160a060020a03909516600160a060020a031990951694909417909355600191909155600255600491909155600355611b8b8061008e6000396000f3fe608060405234801561001057600080fd5b506004361061015f576000357c010000000000000000000000000000000000000000000000000000000090048063785f8ffd116100d5578063d56ba03b11610099578063d56ba03b146103dc578063e02f39bd146103e4578063e64808f314610454578063ea5976a914610471578063eecefef814610479578063facd743b1461054e5761015f565b8063785f8ffd146103325780639ff296ca1461035e578063bb9053d014610366578063c6c21e9d14610392578063d2bfc1c7146103b65761015f565b806328bde1e11161012757806328bde1e11461024857806356a3b5fa1461030a5780635d281ec71461031257806363a20c061461031a5780636e7cf85d14610322578063733975971461032a5761015f565b8063026e402b146101645780630f574ba7146101925780631cfe4f0b146101af5780632218d255146101c957806326c36617146101d1575b600080fd5b6101906004803603604081101561017a57600080fd5b50600160a060020a038135169060200135610588565b005b610190600480360360208110156101a857600080fd5b50356106d0565b6101b7610804565b60408051918252519081900360200190f35b6101b7610847565b610190600480360360408110156101e757600080fd5b8135919081019060408101602082013564010000000081111561020957600080fd5b82018360208201111561021b57600080fd5b8035906020019184600183028401116401000000008311171561023d57600080fd5b50909250905061084d565b61026e6004803603602081101561025e57600080fd5b5035600160a060020a031661092f565b60405180861515151581526020018581526020018060200184815260200183151515158152602001828103825285818151815260200191508051906020019080838360005b838110156102cb5781810151838201526020016102b3565b50505050905090810190601f1680156102f85780820380516001836020036101000a031916815260200191505b50965050505050505060405180910390f35b6101b7610a0a565b6101b7610ab0565b6101b7610ab6565b610190610abb565b6101b7610d6a565b6101906004803603604081101561034857600080fd5b50600160a060020a038135169060200135610d70565b6101b7610f79565b6101906004803603604081101561037c57600080fd5b50600160a060020a038135169060200135610f7f565b61039a611090565b60408051600160a060020a039092168252519081900360200190f35b610190600480360360208110156103cc57600080fd5b5035600160a060020a031661109f565b6101b7611215565b610190600480360360208110156103fa57600080fd5b81019060208101813564010000000081111561041557600080fd5b82018360208201111561042757600080fd5b8035906020019184600183028401116401000000008311171561044957600080fd5b50909250905061121b565b61039a6004803603602081101561046a57600080fd5b503561143e565b6101b761145b565b6104a76004803603604081101561048f57600080fd5b50600160a060020a0381358116916020013516611461565b604051808581526020018060200180602001848152602001838103835286818151815260200191508051906020019060200280838360005b838110156104f75781810151838201526020016104df565b50505050905001838103825285818151815260200191508051906020019060200280838360005b8381101561053657818101518382015260200161051e565b50505050905001965050505050505060405180910390f35b6105746004803603602081101561056457600080fd5b5035600160a060020a0316611592565b604080519115158252519081900360200190f35b81600160a060020a0381166105d5576040805160e560020a62461bcd0281526020600482015260096024820152600080516020611b0c833981519152604482015290519081900360640190fd5b600160a060020a0383166000908152601260205260409020805460ff16610646576040805160e560020a62461bcd02815260206004820152601c60248201527f43616e646964617465206973206e6f7420696e697469616c697a656400000000604482015290519081900360640190fd5b3361065482828660006115c9565b60005461067290600160a060020a031682308763ffffffff6116b916565b84600160a060020a031681600160a060020a03167f500599802164a08023e87ffc3eed0ba3ae60697b3083ba81d046683679d81c6b868560030154604051808381526020018281526020019250505060405180910390a35050505050565b60035442101561072a576040805160e560020a62461bcd02815260206004820152601560248201527f53696465636861696e206973206e6f74206c6976650000000000000000000000604482015290519081900360640190fd5b600254610735610804565b101561078b576040805160e560020a62461bcd02815260206004820152601260248201527f546f6f206665772076616c696461746f72730000000000000000000000000000604482015290519081900360640190fd5b60055433906107a0908363ffffffff61176216565b6005556000546107c190600160a060020a031682308563ffffffff6116b916565b604080518381529051600160a060020a038316917fac095ced75d400384d8404a37883183a56b937b8ff8195fa0c52c3ccc8bb8a49919081900360200190a25050565b600080805b600b811015610840576000600782600b811061082157fe5b0154600160a060020a031614610838576001909101905b600101610809565b5090505b90565b60015481565b336000908152601260205260409020805460ff16156108b6576040805160e560020a62461bcd02815260206004820152601860248201527f43616e64696461746520697320696e697469616c697a65640000000000000000604482015290519081900360640190fd5b805460ff19166001908117825581018490556108d6600282018484611a59565b50828260405180838380828437604080519390910183900383208a845290519095503394507f377f6597c5132797119197fab0e953f73fd5bb109a897a11c871834af2d092a9935091829003602001919050a350505050565b600160a060020a038116600090815260126020908152604080832080546001808301546002808501805487516101009582161595909502600019011691909104601f810188900488028401880190965285835260ff9093169690956060959194859490939290918301828280156109e75780601f106109bc576101008083540402835291602001916109e7565b820191906000526020600020905b8154815290600101906020018083116109ca57829003601f168201915b50505050509350806003015492506109fe87611592565b91505091939590929450565b600754600160a060020a031660009081526012602052604081206003015460015b600b8110156108405781610a3e57610840565b8160126000600784600b8110610a5057fe5b0154600160a060020a031681526020810191909152604001600020600301541015610aa85760126000600783600b8110610a8657fe5b0154600160a060020a0316815260208101919091526040016000206003015491505b600101610a2b565b60045481565b600b81565b336000818152601260205260409020805460ff16610b23576040805160e560020a62461bcd02815260206004820152601c60248201527f43616e646964617465206973206e6f7420696e697469616c697a656400000000604482015290519081900360640190fd5b6000600582015460ff166002811115610b3857fe5b14610b4257600080fd5b60045481600301541015610ba0576040805160e560020a62461bcd02815260206004820152601660248201527f4e6f7420656e6f75676820746f74616c207374616b6500000000000000000000604482015290519081900360640190fd5b6001810154600160a060020a03831660009081526004830160205260409020541015610c16576040805160e560020a62461bcd02815260206004820152601560248201527f4e6f7420656e6f7567682073656c66207374616b650000000000000000000000604482015290519081900360640190fd5b600754600160a060020a031660009081526012602052604081206003015460015b600b811015610d2d5784600160a060020a0316600782600b8110610c5757fe5b0154600160a060020a03161415610cb8576040805160e560020a62461bcd02815260206004820152601860248201527f416c726561647920696e2076616c696461746f72207365740000000000000000604482015290519081900360640190fd5b8160126000600784600b8110610cca57fe5b0154600160a060020a031681526020810191909152604001600020600301541015610d255780925060126000600783600b8110610d0357fe5b0154600160a060020a0316815260208101919091526040016000206003015491505b600101610c37565b506000600783600b8110610d3d57fe5b0154600160a060020a031690508015610d5957610d598361177b565b610d63858461184c565b5050505050565b60065481565b81600160a060020a038116610dbd576040805160e560020a62461bcd0281526020600482015260096024820152600080516020611b0c833981519152604482015290519081900360640190fd5b600160a060020a038316600090815260126020908152604080832033808552600482019093529220909190610df582848760016115c9565b610dfd611ad7565b8581526001600584015460ff166002811115610e1557fe5b1415610e8e57600084600160a060020a031688600160a060020a0316148015610e42575060018401548354105b6004546003860154919250118180610e575750805b15610e6d57610e6d610e688a6118f9565b61177b565b600154610e8190429063ffffffff61176216565b602084015250610ef29050565b6002600584015460ff166002811115610ea357fe5b1415610eb85760068301546020820152610ef2565b60405160e560020a62461bcd02815260040180806020018281038252602b815260200180611b2c602b913960400191505060405180910390fd5b600180830180548083018255600091825260209182902084516002909202019081558184015192018290556003850154604080518a815292830193909352818301529051600160a060020a03808a1692908716917f9e772df8b63d7657919bf5919c475e4033a0bd817b2468bc7ced0d962f21ded09181900360600190a350505050505050565b60035481565b81600160a060020a038116610fcc576040805160e560020a62461bcd0281526020600482015260096024820152600080516020611b0c833981519152604482015290519081900360640190fd5b600160a060020a038316600090815260126020526040812090600582015460ff166002811115610ff857fe5b1461100257600080fd5b336000818152600483016020526040902061102083838760016115c9565b60005461103d90600160a060020a0316838763ffffffff61198e16565b85600160a060020a031682600160a060020a03167f585e40624b400c05be4193af453d2fd2e69facd17163bda6afd44546f3dbbaa8876040518082815260200191505060405180910390a3505050505050565b600054600160a060020a031681565b80600160a060020a0381166110ec576040805160e560020a62461bcd0281526020600482015260096024820152600080516020611b0c833981519152604482015290519081900360640190fd5b600160a060020a038216600090815260126020908152604080832033808552600490910190925282206001810154600282015492939192909142915b838110156111a35784600101818154811061113f57fe5b90600052602060002090600202016001015483111561118f5761118885600101828154811061116a57fe5b6000918252602090912060029091020154839063ffffffff61176216565b915061119b565b600285018190556111a3565b600101611128565b506000546111c190600160a060020a0316868363ffffffff61198e16565b86600160a060020a031685600160a060020a03167f08d0283ea9a2e520a2f09611cf37ca6eb70f62e9a807e53756047dd2dc027220836040518082815260200191505060405180910390a350505050505050565b60055481565b3361122581611592565b1561127a576040805160e560020a62461bcd02815260206004820152601760248201527f6d73672e73656e6465722069732076616c696461746f72000000000000000000604482015290519081900360640190fd5b600160a060020a0381166000908152601260205260409020805460ff166112eb576040805160e560020a62461bcd02815260206004820152601c60248201527f43616e646964617465206973206e6f7420696e697469616c697a656400000000604482015290519081900360640190fd5b600281810180546040805160206001841615610100026000190190931694909404601f810183900483028501830190915280845260609392918301828280156113755780601f1061134a57610100808354040283529160200191611375565b820191906000526020600020905b81548152906001019060200180831161135857829003601f168201915b5093945061138d935050506002840190508686611a59565b5084846040518083838082843760405192018290038220865190955086945091925082916020850191508083835b602083106113da5780518252601f1990920191602091820191016113bb565b5181516020939093036101000a600019018019909116921691909117905260405192018290038220935050600160a060020a03871691507f16de3299ab034ce7e21b22d55f4f9a1474bd3c4d20dbd1cc9bcd39c1ad3d5a2c90600090a45050505050565b600781600b811061144b57fe5b0154600160a060020a0316905081565b60025481565b600160a060020a038083166000908152601260209081526040808320938516835260049093018152828220600181015484518181528184028101909301909452919260609283928592908180156114c2578160200160208202803883390190505b509450806040519080825280602002602001820160405280156114ef578160200160208202803883390190505b50935060005b600183015481101561157a5782600101818154811061151057fe5b90600052602060002090600202016000015486828151811061152e57fe5b60200260200101818152505082600101818154811061154957fe5b90600052602060002090600202016001015485828151811061156757fe5b60209081029190910101526001016114f5565b50508054600290910154909793965091945090925050565b60006001600160a060020a03831660009081526012602052604090206005015460ff1660028111156115c057fe5b1490505b919050565b60008160018111156115d757fe5b141561164257600160a060020a0383166000908152600485016020526040902054611608908363ffffffff61176216565b600160a060020a03841660009081526004860160205260409020556003840154611638908363ffffffff61176216565b60038501556116b3565b600181600181111561165057fe5b14156116b157600160a060020a0383166000908152600485016020526040902054611681908363ffffffff611a4416565b600160a060020a03841660009081526004860160205260409020556003840154611638908363ffffffff611a4416565bfe5b50505050565b604080517f23b872dd000000000000000000000000000000000000000000000000000000008152600160a060020a0385811660048301528481166024830152604482018490529151918616916323b872dd916064808201926020929091908290030181600087803b15801561172d57600080fd5b505af1158015611741573d6000803e3d6000fd5b505050506040513d602081101561175757600080fd5b50516116b357600080fd5b60008282018381101561177457600080fd5b9392505050565b6000600782600b811061178a57fe5b0154600160a060020a03169050806117a25750611849565b600782600b81106117af57fe5b01805473ffffffffffffffffffffffffffffffffffffffff19169055600160a060020a0381166000908152601260205260409020600501805460ff191660021790556001546117ff904290611762565b600160a060020a0382166000818152601260205260408082206006019390935591516001927f63f783ba869265648de5e70add96be9f4914e3bde064fdc19fd7e6a8ebf2f46c91a3505b50565b6000600782600b811061185b57fe5b0154600160a060020a03161461187057600080fd5b81600782600b811061187e57fe5b01805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03928316179055821660008181526012602052604080822060058101805460ff19166001179055600601829055519091907f63f783ba869265648de5e70add96be9f4914e3bde064fdc19fd7e6a8ebf2f46c908390a35050565b6000805b600b81101561193d5782600160a060020a0316600782600b811061191d57fe5b0154600160a060020a031614156119355790506115c4565b6001016118fd565b506040805160e560020a62461bcd02815260206004820152601360248201527f6e6f207375636820612076616c696461746f7200000000000000000000000000604482015290519081900360640190fd5b82600160a060020a031663a9059cbb83836040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018083600160a060020a0316600160a060020a0316815260200182815260200192505050602060405180830381600087803b158015611a0a57600080fd5b505af1158015611a1e573d6000803e3d6000fd5b505050506040513d6020811015611a3457600080fd5b5051611a3f57600080fd5b505050565b600082821115611a5357600080fd5b50900390565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611a9a5782800160ff19823516178555611ac7565b82800160010185558215611ac7579182015b82811115611ac7578235825591602001919060010190611aac565b50611ad3929150611af1565b5090565b604051806040016040528060008152602001600081525090565b61084491905b80821115611ad35760008155600101611af756fe302061646472657373000000000000000000000000000000000000000000000043616e64696461746520737461747573206973206e6f7420426f6e646564206f7220556e626f6e64696e67a265627a7a72305820ba798da195c09fa1d1158dff1d323f263edf762c3f62826d521a630a0c905d9464736f6c634300050a0032`

// DeployGuard deploys a new Ethereum contract, binding an instance of Guard to it.
func DeployGuard(auth *bind.TransactOpts, backend bind.ContractBackend, _celerTokenAddress common.Address, _blameTimeout *big.Int, _minValidatorNum *big.Int, _minTotalStake *big.Int, _sidechainGoLiveTime *big.Int) (common.Address, *types.Transaction, *Guard, error) {
	parsed, err := abi.JSON(strings.NewReader(GuardABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(GuardBin), backend, _celerTokenAddress, _blameTimeout, _minValidatorNum, _minTotalStake, _sidechainGoLiveTime)
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

// BlameTimeout is a free data retrieval call binding the contract method 0x2218d255.
//
// Solidity: function blameTimeout() constant returns(uint256)
func (_Guard *GuardCaller) BlameTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "blameTimeout")
	return *ret0, err
}

// BlameTimeout is a free data retrieval call binding the contract method 0x2218d255.
//
// Solidity: function blameTimeout() constant returns(uint256)
func (_Guard *GuardSession) BlameTimeout() (*big.Int, error) {
	return _Guard.Contract.BlameTimeout(&_Guard.CallOpts)
}

// BlameTimeout is a free data retrieval call binding the contract method 0x2218d255.
//
// Solidity: function blameTimeout() constant returns(uint256)
func (_Guard *GuardCallerSession) BlameTimeout() (*big.Int, error) {
	return _Guard.Contract.BlameTimeout(&_Guard.CallOpts)
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

// GetCandidateInfo is a free data retrieval call binding the contract method 0x28bde1e1.
//
// Solidity: function getCandidateInfo(address _candidateAddr) constant returns(bool initialized, uint256 minSelfStake, bytes sidechainAddr, uint256 totalStake, bool isVldt)
func (_Guard *GuardCaller) GetCandidateInfo(opts *bind.CallOpts, _candidateAddr common.Address) (struct {
	Initialized   bool
	MinSelfStake  *big.Int
	SidechainAddr []byte
	TotalStake    *big.Int
	IsVldt        bool
}, error) {
	ret := new(struct {
		Initialized   bool
		MinSelfStake  *big.Int
		SidechainAddr []byte
		TotalStake    *big.Int
		IsVldt        bool
	})
	out := ret
	err := _Guard.contract.Call(opts, out, "getCandidateInfo", _candidateAddr)
	return *ret, err
}

// GetCandidateInfo is a free data retrieval call binding the contract method 0x28bde1e1.
//
// Solidity: function getCandidateInfo(address _candidateAddr) constant returns(bool initialized, uint256 minSelfStake, bytes sidechainAddr, uint256 totalStake, bool isVldt)
func (_Guard *GuardSession) GetCandidateInfo(_candidateAddr common.Address) (struct {
	Initialized   bool
	MinSelfStake  *big.Int
	SidechainAddr []byte
	TotalStake    *big.Int
	IsVldt        bool
}, error) {
	return _Guard.Contract.GetCandidateInfo(&_Guard.CallOpts, _candidateAddr)
}

// GetCandidateInfo is a free data retrieval call binding the contract method 0x28bde1e1.
//
// Solidity: function getCandidateInfo(address _candidateAddr) constant returns(bool initialized, uint256 minSelfStake, bytes sidechainAddr, uint256 totalStake, bool isVldt)
func (_Guard *GuardCallerSession) GetCandidateInfo(_candidateAddr common.Address) (struct {
	Initialized   bool
	MinSelfStake  *big.Int
	SidechainAddr []byte
	TotalStake    *big.Int
	IsVldt        bool
}, error) {
	return _Guard.Contract.GetCandidateInfo(&_Guard.CallOpts, _candidateAddr)
}

// GetDelegatorInfo is a free data retrieval call binding the contract method 0xeecefef8.
//
// Solidity: function getDelegatorInfo(address _candidateAddr, address _delegatorAddr) constant returns(uint256 stake, uint256[] intentAmounts, uint256[] intentUnlockTimes, uint256 nextWithdrawIntent)
func (_Guard *GuardCaller) GetDelegatorInfo(opts *bind.CallOpts, _candidateAddr common.Address, _delegatorAddr common.Address) (struct {
	Stake              *big.Int
	IntentAmounts      []*big.Int
	IntentUnlockTimes  []*big.Int
	NextWithdrawIntent *big.Int
}, error) {
	ret := new(struct {
		Stake              *big.Int
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
// Solidity: function getDelegatorInfo(address _candidateAddr, address _delegatorAddr) constant returns(uint256 stake, uint256[] intentAmounts, uint256[] intentUnlockTimes, uint256 nextWithdrawIntent)
func (_Guard *GuardSession) GetDelegatorInfo(_candidateAddr common.Address, _delegatorAddr common.Address) (struct {
	Stake              *big.Int
	IntentAmounts      []*big.Int
	IntentUnlockTimes  []*big.Int
	NextWithdrawIntent *big.Int
}, error) {
	return _Guard.Contract.GetDelegatorInfo(&_Guard.CallOpts, _candidateAddr, _delegatorAddr)
}

// GetDelegatorInfo is a free data retrieval call binding the contract method 0xeecefef8.
//
// Solidity: function getDelegatorInfo(address _candidateAddr, address _delegatorAddr) constant returns(uint256 stake, uint256[] intentAmounts, uint256[] intentUnlockTimes, uint256 nextWithdrawIntent)
func (_Guard *GuardCallerSession) GetDelegatorInfo(_candidateAddr common.Address, _delegatorAddr common.Address) (struct {
	Stake              *big.Int
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

// MinTotalStake is a free data retrieval call binding the contract method 0x5d281ec7.
//
// Solidity: function minTotalStake() constant returns(uint256)
func (_Guard *GuardCaller) MinTotalStake(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "minTotalStake")
	return *ret0, err
}

// MinTotalStake is a free data retrieval call binding the contract method 0x5d281ec7.
//
// Solidity: function minTotalStake() constant returns(uint256)
func (_Guard *GuardSession) MinTotalStake() (*big.Int, error) {
	return _Guard.Contract.MinTotalStake(&_Guard.CallOpts)
}

// MinTotalStake is a free data retrieval call binding the contract method 0x5d281ec7.
//
// Solidity: function minTotalStake() constant returns(uint256)
func (_Guard *GuardCallerSession) MinTotalStake() (*big.Int, error) {
	return _Guard.Contract.MinTotalStake(&_Guard.CallOpts)
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

// MiningPool is a free data retrieval call binding the contract method 0x73397597.
//
// Solidity: function miningPool() constant returns(uint256)
func (_Guard *GuardCaller) MiningPool(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "miningPool")
	return *ret0, err
}

// MiningPool is a free data retrieval call binding the contract method 0x73397597.
//
// Solidity: function miningPool() constant returns(uint256)
func (_Guard *GuardSession) MiningPool() (*big.Int, error) {
	return _Guard.Contract.MiningPool(&_Guard.CallOpts)
}

// MiningPool is a free data retrieval call binding the contract method 0x73397597.
//
// Solidity: function miningPool() constant returns(uint256)
func (_Guard *GuardCallerSession) MiningPool() (*big.Int, error) {
	return _Guard.Contract.MiningPool(&_Guard.CallOpts)
}

// SidechainGoLiveTime is a free data retrieval call binding the contract method 0x9ff296ca.
//
// Solidity: function sidechainGoLiveTime() constant returns(uint256)
func (_Guard *GuardCaller) SidechainGoLiveTime(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "sidechainGoLiveTime")
	return *ret0, err
}

// SidechainGoLiveTime is a free data retrieval call binding the contract method 0x9ff296ca.
//
// Solidity: function sidechainGoLiveTime() constant returns(uint256)
func (_Guard *GuardSession) SidechainGoLiveTime() (*big.Int, error) {
	return _Guard.Contract.SidechainGoLiveTime(&_Guard.CallOpts)
}

// SidechainGoLiveTime is a free data retrieval call binding the contract method 0x9ff296ca.
//
// Solidity: function sidechainGoLiveTime() constant returns(uint256)
func (_Guard *GuardCallerSession) SidechainGoLiveTime() (*big.Int, error) {
	return _Guard.Contract.SidechainGoLiveTime(&_Guard.CallOpts)
}

// SubscriptionFees is a free data retrieval call binding the contract method 0xd56ba03b.
//
// Solidity: function subscriptionFees() constant returns(uint256)
func (_Guard *GuardCaller) SubscriptionFees(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Guard.contract.Call(opts, out, "subscriptionFees")
	return *ret0, err
}

// SubscriptionFees is a free data retrieval call binding the contract method 0xd56ba03b.
//
// Solidity: function subscriptionFees() constant returns(uint256)
func (_Guard *GuardSession) SubscriptionFees() (*big.Int, error) {
	return _Guard.Contract.SubscriptionFees(&_Guard.CallOpts)
}

// SubscriptionFees is a free data retrieval call binding the contract method 0xd56ba03b.
//
// Solidity: function subscriptionFees() constant returns(uint256)
func (_Guard *GuardCallerSession) SubscriptionFees() (*big.Int, error) {
	return _Guard.Contract.SubscriptionFees(&_Guard.CallOpts)
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

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address _candidateAddr, uint256 _amount) returns()
func (_Guard *GuardTransactor) Delegate(opts *bind.TransactOpts, _candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "delegate", _candidateAddr, _amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address _candidateAddr, uint256 _amount) returns()
func (_Guard *GuardSession) Delegate(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.Delegate(&_Guard.TransactOpts, _candidateAddr, _amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address _candidateAddr, uint256 _amount) returns()
func (_Guard *GuardTransactorSession) Delegate(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.Delegate(&_Guard.TransactOpts, _candidateAddr, _amount)
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

// IntendWithdraw is a paid mutator transaction binding the contract method 0x785f8ffd.
//
// Solidity: function intendWithdraw(address _candidateAddr, uint256 _amount) returns()
func (_Guard *GuardTransactor) IntendWithdraw(opts *bind.TransactOpts, _candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "intendWithdraw", _candidateAddr, _amount)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x785f8ffd.
//
// Solidity: function intendWithdraw(address _candidateAddr, uint256 _amount) returns()
func (_Guard *GuardSession) IntendWithdraw(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.IntendWithdraw(&_Guard.TransactOpts, _candidateAddr, _amount)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x785f8ffd.
//
// Solidity: function intendWithdraw(address _candidateAddr, uint256 _amount) returns()
func (_Guard *GuardTransactorSession) IntendWithdraw(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.IntendWithdraw(&_Guard.TransactOpts, _candidateAddr, _amount)
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

// WithdrawFromUnbondedCandidate is a paid mutator transaction binding the contract method 0xbb9053d0.
//
// Solidity: function withdrawFromUnbondedCandidate(address _candidateAddr, uint256 _amount) returns()
func (_Guard *GuardTransactor) WithdrawFromUnbondedCandidate(opts *bind.TransactOpts, _candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.contract.Transact(opts, "withdrawFromUnbondedCandidate", _candidateAddr, _amount)
}

// WithdrawFromUnbondedCandidate is a paid mutator transaction binding the contract method 0xbb9053d0.
//
// Solidity: function withdrawFromUnbondedCandidate(address _candidateAddr, uint256 _amount) returns()
func (_Guard *GuardSession) WithdrawFromUnbondedCandidate(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.WithdrawFromUnbondedCandidate(&_Guard.TransactOpts, _candidateAddr, _amount)
}

// WithdrawFromUnbondedCandidate is a paid mutator transaction binding the contract method 0xbb9053d0.
//
// Solidity: function withdrawFromUnbondedCandidate(address _candidateAddr, uint256 _amount) returns()
func (_Guard *GuardTransactorSession) WithdrawFromUnbondedCandidate(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Guard.Contract.WithdrawFromUnbondedCandidate(&_Guard.TransactOpts, _candidateAddr, _amount)
}

// GuardAddSubscriptionBalanceIterator is returned from FilterAddSubscriptionBalance and is used to iterate over the raw logs and unpacked data for AddSubscriptionBalance events raised by the Guard contract.
type GuardAddSubscriptionBalanceIterator struct {
	Event *GuardAddSubscriptionBalance // Event containing the contract specifics and raw log

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
func (it *GuardAddSubscriptionBalanceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardAddSubscriptionBalance)
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
		it.Event = new(GuardAddSubscriptionBalance)
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
func (it *GuardAddSubscriptionBalanceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardAddSubscriptionBalanceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardAddSubscriptionBalance represents a AddSubscriptionBalance event raised by the Guard contract.
type GuardAddSubscriptionBalance struct {
	Consumer common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddSubscriptionBalance is a free log retrieval operation binding the contract event 0xac095ced75d400384d8404a37883183a56b937b8ff8195fa0c52c3ccc8bb8a49.
//
// Solidity: event AddSubscriptionBalance(address indexed consumer, uint256 amount)
func (_Guard *GuardFilterer) FilterAddSubscriptionBalance(opts *bind.FilterOpts, consumer []common.Address) (*GuardAddSubscriptionBalanceIterator, error) {

	var consumerRule []interface{}
	for _, consumerItem := range consumer {
		consumerRule = append(consumerRule, consumerItem)
	}

	logs, sub, err := _Guard.contract.FilterLogs(opts, "AddSubscriptionBalance", consumerRule)
	if err != nil {
		return nil, err
	}
	return &GuardAddSubscriptionBalanceIterator{contract: _Guard.contract, event: "AddSubscriptionBalance", logs: logs, sub: sub}, nil
}

// WatchAddSubscriptionBalance is a free log subscription operation binding the contract event 0xac095ced75d400384d8404a37883183a56b937b8ff8195fa0c52c3ccc8bb8a49.
//
// Solidity: event AddSubscriptionBalance(address indexed consumer, uint256 amount)
func (_Guard *GuardFilterer) WatchAddSubscriptionBalance(opts *bind.WatchOpts, sink chan<- *GuardAddSubscriptionBalance, consumer []common.Address) (event.Subscription, error) {

	var consumerRule []interface{}
	for _, consumerItem := range consumer {
		consumerRule = append(consumerRule, consumerItem)
	}

	logs, sub, err := _Guard.contract.WatchLogs(opts, "AddSubscriptionBalance", consumerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardAddSubscriptionBalance)
				if err := _Guard.contract.UnpackLog(event, "AddSubscriptionBalance", log); err != nil {
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
	Delegator  common.Address
	Candidate  common.Address
	NewStake   *big.Int
	TotalStake *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDelegate is a free log retrieval operation binding the contract event 0x500599802164a08023e87ffc3eed0ba3ae60697b3083ba81d046683679d81c6b.
//
// Solidity: event Delegate(address indexed delegator, address indexed candidate, uint256 newStake, uint256 totalStake)
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
// Solidity: event Delegate(address indexed delegator, address indexed candidate, uint256 newStake, uint256 totalStake)
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
	Delegator      common.Address
	Candidate      common.Address
	WithdrawAmount *big.Int
	UnlockTime     *big.Int
	TotalStake     *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterIntendWithdraw is a free log retrieval operation binding the contract event 0x9e772df8b63d7657919bf5919c475e4033a0bd817b2468bc7ced0d962f21ded0.
//
// Solidity: event IntendWithdraw(address indexed delegator, address indexed candidate, uint256 withdrawAmount, uint256 unlockTime, uint256 totalStake)
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
// Solidity: event IntendWithdraw(address indexed delegator, address indexed candidate, uint256 withdrawAmount, uint256 unlockTime, uint256 totalStake)
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

// GuardWithdrawFromUnbondedCandidateIterator is returned from FilterWithdrawFromUnbondedCandidate and is used to iterate over the raw logs and unpacked data for WithdrawFromUnbondedCandidate events raised by the Guard contract.
type GuardWithdrawFromUnbondedCandidateIterator struct {
	Event *GuardWithdrawFromUnbondedCandidate // Event containing the contract specifics and raw log

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
func (it *GuardWithdrawFromUnbondedCandidateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GuardWithdrawFromUnbondedCandidate)
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
		it.Event = new(GuardWithdrawFromUnbondedCandidate)
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
func (it *GuardWithdrawFromUnbondedCandidateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GuardWithdrawFromUnbondedCandidateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GuardWithdrawFromUnbondedCandidate represents a WithdrawFromUnbondedCandidate event raised by the Guard contract.
type GuardWithdrawFromUnbondedCandidate struct {
	Delegator common.Address
	Candidate common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawFromUnbondedCandidate is a free log retrieval operation binding the contract event 0x585e40624b400c05be4193af453d2fd2e69facd17163bda6afd44546f3dbbaa8.
//
// Solidity: event WithdrawFromUnbondedCandidate(address indexed delegator, address indexed candidate, uint256 amount)
func (_Guard *GuardFilterer) FilterWithdrawFromUnbondedCandidate(opts *bind.FilterOpts, delegator []common.Address, candidate []common.Address) (*GuardWithdrawFromUnbondedCandidateIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _Guard.contract.FilterLogs(opts, "WithdrawFromUnbondedCandidate", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return &GuardWithdrawFromUnbondedCandidateIterator{contract: _Guard.contract, event: "WithdrawFromUnbondedCandidate", logs: logs, sub: sub}, nil
}

// WatchWithdrawFromUnbondedCandidate is a free log subscription operation binding the contract event 0x585e40624b400c05be4193af453d2fd2e69facd17163bda6afd44546f3dbbaa8.
//
// Solidity: event WithdrawFromUnbondedCandidate(address indexed delegator, address indexed candidate, uint256 amount)
func (_Guard *GuardFilterer) WatchWithdrawFromUnbondedCandidate(opts *bind.WatchOpts, sink chan<- *GuardWithdrawFromUnbondedCandidate, delegator []common.Address, candidate []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _Guard.contract.WatchLogs(opts, "WithdrawFromUnbondedCandidate", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GuardWithdrawFromUnbondedCandidate)
				if err := _Guard.contract.UnpackLog(event, "WithdrawFromUnbondedCandidate", log); err != nil {
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
