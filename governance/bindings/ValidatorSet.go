// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// ValidatorSetABI is the input ABI used to generate the binding from.
const ValidatorSetABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_validatorPubKeys\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pubKey\",\"type\":\"string\"},{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"removeValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_nextVersion\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pubKey\",\"type\":\"string\"},{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"addValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"bool\"}],\"name\":\"setFreezeStatus\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_freeze\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"address\"}],\"name\":\"_setNextVersionAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_pubKey\",\"type\":\"string\"}],\"name\":\"validatorAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"validatorPubKey\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorSetSize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_key\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_pubKey\",\"type\":\"string\"}],\"name\":\"ValidatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_key\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_pubKey\",\"type\":\"string\"}],\"name\":\"ValidatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Freeze\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"SetNextVersion\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// ValidatorSetBin is the compiled bytecode used for deploying new contracts.
const ValidatorSetBin = `60806040526000600360006101000a81548160ff021916908315150217905550336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3611893806100ea6000396000f3fe6080604052600436106100d0576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806308d98bb3146100d55780630e62d9eb14610189578063369b4b96146102715780633e8bb9a0146102c857806341edb775146103b057806365ed50e9146103ed578063715018a61461041c5780638da5cb5b146104335780638f32d59b1461048a578063b2684966146104b9578063d4b0d70a1461050a578063f047a31614610612578063f2fde38b146106c6578063f800db5014610717575b600080fd5b3480156100e157600080fd5b5061010e600480360360208110156100f857600080fd5b8101908080359060200190929190505050610742565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561014e578082015181840152602081019050610133565b50505050905090810190601f16801561017b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561019557600080fd5b5061026f600480360360408110156101ac57600080fd5b81019080803590602001906401000000008111156101c957600080fd5b8201836020820111156101db57600080fd5b803590602001918460018302840111640100000000831117156101fd57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506107fd565b005b34801561027d57600080fd5b50610286610c20565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102d457600080fd5b506103ae600480360360408110156102eb57600080fd5b810190808035906020019064010000000081111561030857600080fd5b82018360208201111561031a57600080fd5b8035906020019184600183028401116401000000008311171561033c57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610c46565b005b3480156103bc57600080fd5b506103eb600480360360208110156103d357600080fd5b81019080803515159060200190929190505050610f19565b005b3480156103f957600080fd5b50610402610fde565b604051808215151515815260200191505060405180910390f35b34801561042857600080fd5b50610431610ff1565b005b34801561043f57600080fd5b5061044861112c565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561049657600080fd5b5061049f611155565b604051808215151515815260200191505060405180910390f35b3480156104c557600080fd5b50610508600480360360208110156104dc57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506111ac565b005b34801561051657600080fd5b506105d06004803603602081101561052d57600080fd5b810190808035906020019064010000000081111561054a57600080fd5b82018360208201111561055c57600080fd5b8035906020019184600183028401116401000000008311171561057e57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506112cf565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561061e57600080fd5b5061064b6004803603602081101561063557600080fd5b8101908080359060200190929190505050611319565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561068b578082015181840152602081019050610670565b50505050905090810190601f1680156106b85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156106d257600080fd5b50610715600480360360208110156106e957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506113d4565b005b34801561072357600080fd5b5061072c61145c565b6040518082815260200191505060405180910390f35b60028181548110151561075157fe5b906000526020600020016000915090508054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107f55780601f106107ca576101008083540402835291602001916107f5565b820191906000526020600020905b8154815290600101906020018083116107d857829003601f168201915b505050505081565b610805611155565b1515610879576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b60001515600360009054906101000a900460ff16151514151561089b57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141515156108d757600080fd5b602882511415156108e757600080fd5b60006108f283611469565b90508173ffffffffffffffffffffffffffffffffffffffff166001600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614151561096157600080fd5b60006001600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600080905060008090505b600160028054905003811015610afe5782610a8a6002838154811015156109e257fe5b906000526020600020018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610a805780601f10610a5557610100808354040283529160200191610a80565b820191906000526020600020905b815481529060010190602001808311610a6357829003601f168201915b5050505050611469565b1415610a9557600191505b8115610af157600260018201815481101515610aad57fe5b90600052602060002001600282815481101515610ac657fe5b906000526020600020019080546001816001161561010002031660029004610aef92919061169b565b505b80806001019150506109bf565b506002600160028054905003815481101515610b1657fe5b906000526020600020016000610b2c9190611722565b6002805480919060019003610b41919061176a565b507ff10b1758a748f201874627a3eb8553553531350314c7da7d9a13f2ae457e165d828486604051808481526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610bde578082015181840152602081019050610bc3565b50505050905090810190601f168015610c0b5780820380516001836020036101000a031916815260200191505b5094505050505060405180910390a150505050565b600360019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b610c4e611155565b1515610cc2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b60001515600360009054906101000a900460ff161515141515610ce457600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515610d2057600080fd5b60288251141515610d3057600080fd5b6000610d3b83611469565b9050600073ffffffffffffffffffffffffffffffffffffffff166001600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141515610dab57600080fd5b816001600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506002839080600181540180825580915050906001820390600052602060002001600090919290919091509080519060200190610e3a929190611796565b50507f940f4642cd572b4a8d9cf8545df0a213cbf8bee450f9ea7c781ede70650c4ab6818385604051808481526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610ed8578082015181840152602081019050610ebd565b50505050905090810190601f168015610f055780820380516001836020036101000a031916815260200191505b5094505050505060405180910390a1505050565b610f21611155565b1515610f95576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b80600360006101000a81548160ff0219169083151502179055507f615acbaede366d76a8b8cb2a9ada6a71495f0786513d71aa97aaf0c3910b78de60405160405180910390a150565b600360009054906101000a900460ff1681565b610ff9611155565b151561106d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a360008060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614905090565b6111b4611155565b1515611228576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b80600360016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507fb694c35dc1cdbab3df571f07de838161b83be5527eacfcb8a4b589b70943874581604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a150565b6000806112db83611469565b90506001600082815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16915050919050565b606060028281548110151561132a57fe5b906000526020600020018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156113c85780601f1061139d576101008083540402835291602001916113c8565b820191906000526020600020905b8154815290600101906020018083116113ab57829003601f168201915b50505050509050919050565b6113dc611155565b1515611450576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b61145981611512565b50565b6000600280549050905090565b60006002826040518082805190602001908083835b6020831015156114a3578051825260208201915060208101905060208303925061147e565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa1580156114e5573d6000803e3d6000fd5b5050506040513d60208110156114fa57600080fd5b81019080805190602001909291905050509050919050565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141515156115dd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001807f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206181526020017f646472657373000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106116d45780548555611711565b8280016001018555821561171157600052602060002091601f016020900482015b828111156117105782548255916001019190600101906116f5565b5b50905061171e9190611816565b5090565b50805460018160011615610100020316600290046000825580601f106117485750611767565b601f0160209004906000526020600020908101906117669190611816565b5b50565b81548183558181111561179157818360005260206000209182019101611790919061183b565b5b505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106117d757805160ff1916838001178555611805565b82800160010185558215611805579182015b828111156118045782518255916020019190600101906117e9565b5b5090506118129190611816565b5090565b61183891905b8082111561183457600081600090555060010161181c565b5090565b90565b61186491905b8082111561186057600081816118579190611722565b50600101611841565b5090565b9056fea165627a7a723058208ce14dd1f1521597299ab99b71d6a22a3a90b4e81921889046068f03590767550029`

// DeployValidatorSet deploys a new Ethereum contract, binding an instance of ValidatorSet to it.
func DeployValidatorSet(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ValidatorSet, error) {
	parsed, err := abi.JSON(strings.NewReader(ValidatorSetABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ValidatorSetBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ValidatorSet{ValidatorSetCaller: ValidatorSetCaller{contract: contract}, ValidatorSetTransactor: ValidatorSetTransactor{contract: contract}, ValidatorSetFilterer: ValidatorSetFilterer{contract: contract}}, nil
}

// ValidatorSet is an auto generated Go binding around an Ethereum contract.
type ValidatorSet struct {
	ValidatorSetCaller     // Read-only binding to the contract
	ValidatorSetTransactor // Write-only binding to the contract
	ValidatorSetFilterer   // Log filterer for contract events
}

// ValidatorSetCaller is an auto generated read-only Go binding around an Ethereum contract.
type ValidatorSetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorSetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ValidatorSetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorSetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValidatorSetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorSetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValidatorSetSession struct {
	Contract     *ValidatorSet     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ValidatorSetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValidatorSetCallerSession struct {
	Contract *ValidatorSetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ValidatorSetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValidatorSetTransactorSession struct {
	Contract     *ValidatorSetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ValidatorSetRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValidatorSetRaw struct {
	Contract *ValidatorSet // Generic contract binding to access the raw methods on
}

// ValidatorSetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValidatorSetCallerRaw struct {
	Contract *ValidatorSetCaller // Generic read-only contract binding to access the raw methods on
}

// ValidatorSetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValidatorSetTransactorRaw struct {
	Contract *ValidatorSetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValidatorSet creates a new instance of ValidatorSet, bound to a specific deployed contract.
func NewValidatorSet(address common.Address, backend bind.ContractBackend) (*ValidatorSet, error) {
	contract, err := bindValidatorSet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ValidatorSet{ValidatorSetCaller: ValidatorSetCaller{contract: contract}, ValidatorSetTransactor: ValidatorSetTransactor{contract: contract}, ValidatorSetFilterer: ValidatorSetFilterer{contract: contract}}, nil
}

// NewValidatorSetCaller creates a new read-only instance of ValidatorSet, bound to a specific deployed contract.
func NewValidatorSetCaller(address common.Address, caller bind.ContractCaller) (*ValidatorSetCaller, error) {
	contract, err := bindValidatorSet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorSetCaller{contract: contract}, nil
}

// NewValidatorSetTransactor creates a new write-only instance of ValidatorSet, bound to a specific deployed contract.
func NewValidatorSetTransactor(address common.Address, transactor bind.ContractTransactor) (*ValidatorSetTransactor, error) {
	contract, err := bindValidatorSet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorSetTransactor{contract: contract}, nil
}

// NewValidatorSetFilterer creates a new log filterer instance of ValidatorSet, bound to a specific deployed contract.
func NewValidatorSetFilterer(address common.Address, filterer bind.ContractFilterer) (*ValidatorSetFilterer, error) {
	contract, err := bindValidatorSet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValidatorSetFilterer{contract: contract}, nil
}

// bindValidatorSet binds a generic wrapper to an already deployed contract.
func bindValidatorSet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ValidatorSetABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidatorSet *ValidatorSetRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ValidatorSet.Contract.ValidatorSetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidatorSet *ValidatorSetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorSet.Contract.ValidatorSetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidatorSet *ValidatorSetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ValidatorSet.Contract.ValidatorSetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidatorSet *ValidatorSetCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ValidatorSet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidatorSet *ValidatorSetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorSet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidatorSet *ValidatorSetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ValidatorSet.Contract.contract.Transact(opts, method, params...)
}

// Freeze is a free data retrieval call binding the contract method 0x65ed50e9.
//
// Solidity: function _freeze() constant returns(bool)
func (_ValidatorSet *ValidatorSetCaller) Freeze(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "_freeze")
	return *ret0, err
}

// Freeze is a free data retrieval call binding the contract method 0x65ed50e9.
//
// Solidity: function _freeze() constant returns(bool)
func (_ValidatorSet *ValidatorSetSession) Freeze() (bool, error) {
	return _ValidatorSet.Contract.Freeze(&_ValidatorSet.CallOpts)
}

// Freeze is a free data retrieval call binding the contract method 0x65ed50e9.
//
// Solidity: function _freeze() constant returns(bool)
func (_ValidatorSet *ValidatorSetCallerSession) Freeze() (bool, error) {
	return _ValidatorSet.Contract.Freeze(&_ValidatorSet.CallOpts)
}

// NextVersion is a free data retrieval call binding the contract method 0x369b4b96.
//
// Solidity: function _nextVersion() constant returns(address)
func (_ValidatorSet *ValidatorSetCaller) NextVersion(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "_nextVersion")
	return *ret0, err
}

// NextVersion is a free data retrieval call binding the contract method 0x369b4b96.
//
// Solidity: function _nextVersion() constant returns(address)
func (_ValidatorSet *ValidatorSetSession) NextVersion() (common.Address, error) {
	return _ValidatorSet.Contract.NextVersion(&_ValidatorSet.CallOpts)
}

// NextVersion is a free data retrieval call binding the contract method 0x369b4b96.
//
// Solidity: function _nextVersion() constant returns(address)
func (_ValidatorSet *ValidatorSetCallerSession) NextVersion() (common.Address, error) {
	return _ValidatorSet.Contract.NextVersion(&_ValidatorSet.CallOpts)
}

// ValidatorPubKeys is a free data retrieval call binding the contract method 0x08d98bb3.
//
// Solidity: function _validatorPubKeys(uint256 ) constant returns(string)
func (_ValidatorSet *ValidatorSetCaller) ValidatorPubKeys(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "_validatorPubKeys", arg0)
	return *ret0, err
}

// ValidatorPubKeys is a free data retrieval call binding the contract method 0x08d98bb3.
//
// Solidity: function _validatorPubKeys(uint256 ) constant returns(string)
func (_ValidatorSet *ValidatorSetSession) ValidatorPubKeys(arg0 *big.Int) (string, error) {
	return _ValidatorSet.Contract.ValidatorPubKeys(&_ValidatorSet.CallOpts, arg0)
}

// ValidatorPubKeys is a free data retrieval call binding the contract method 0x08d98bb3.
//
// Solidity: function _validatorPubKeys(uint256 ) constant returns(string)
func (_ValidatorSet *ValidatorSetCallerSession) ValidatorPubKeys(arg0 *big.Int) (string, error) {
	return _ValidatorSet.Contract.ValidatorPubKeys(&_ValidatorSet.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_ValidatorSet *ValidatorSetCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_ValidatorSet *ValidatorSetSession) IsOwner() (bool, error) {
	return _ValidatorSet.Contract.IsOwner(&_ValidatorSet.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_ValidatorSet *ValidatorSetCallerSession) IsOwner() (bool, error) {
	return _ValidatorSet.Contract.IsOwner(&_ValidatorSet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ValidatorSet *ValidatorSetCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ValidatorSet *ValidatorSetSession) Owner() (common.Address, error) {
	return _ValidatorSet.Contract.Owner(&_ValidatorSet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ValidatorSet *ValidatorSetCallerSession) Owner() (common.Address, error) {
	return _ValidatorSet.Contract.Owner(&_ValidatorSet.CallOpts)
}

// ValidatorAddress is a free data retrieval call binding the contract method 0xd4b0d70a.
//
// Solidity: function validatorAddress(string _pubKey) constant returns(address)
func (_ValidatorSet *ValidatorSetCaller) ValidatorAddress(opts *bind.CallOpts, _pubKey string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "validatorAddress", _pubKey)
	return *ret0, err
}

// ValidatorAddress is a free data retrieval call binding the contract method 0xd4b0d70a.
//
// Solidity: function validatorAddress(string _pubKey) constant returns(address)
func (_ValidatorSet *ValidatorSetSession) ValidatorAddress(_pubKey string) (common.Address, error) {
	return _ValidatorSet.Contract.ValidatorAddress(&_ValidatorSet.CallOpts, _pubKey)
}

// ValidatorAddress is a free data retrieval call binding the contract method 0xd4b0d70a.
//
// Solidity: function validatorAddress(string _pubKey) constant returns(address)
func (_ValidatorSet *ValidatorSetCallerSession) ValidatorAddress(_pubKey string) (common.Address, error) {
	return _ValidatorSet.Contract.ValidatorAddress(&_ValidatorSet.CallOpts, _pubKey)
}

// ValidatorPubKey is a free data retrieval call binding the contract method 0xf047a316.
//
// Solidity: function validatorPubKey(uint256 index) constant returns(string)
func (_ValidatorSet *ValidatorSetCaller) ValidatorPubKey(opts *bind.CallOpts, index *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "validatorPubKey", index)
	return *ret0, err
}

// ValidatorPubKey is a free data retrieval call binding the contract method 0xf047a316.
//
// Solidity: function validatorPubKey(uint256 index) constant returns(string)
func (_ValidatorSet *ValidatorSetSession) ValidatorPubKey(index *big.Int) (string, error) {
	return _ValidatorSet.Contract.ValidatorPubKey(&_ValidatorSet.CallOpts, index)
}

// ValidatorPubKey is a free data retrieval call binding the contract method 0xf047a316.
//
// Solidity: function validatorPubKey(uint256 index) constant returns(string)
func (_ValidatorSet *ValidatorSetCallerSession) ValidatorPubKey(index *big.Int) (string, error) {
	return _ValidatorSet.Contract.ValidatorPubKey(&_ValidatorSet.CallOpts, index)
}

// ValidatorSetSize is a free data retrieval call binding the contract method 0xf800db50.
//
// Solidity: function validatorSetSize() constant returns(uint256)
func (_ValidatorSet *ValidatorSetCaller) ValidatorSetSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "validatorSetSize")
	return *ret0, err
}

// ValidatorSetSize is a free data retrieval call binding the contract method 0xf800db50.
//
// Solidity: function validatorSetSize() constant returns(uint256)
func (_ValidatorSet *ValidatorSetSession) ValidatorSetSize() (*big.Int, error) {
	return _ValidatorSet.Contract.ValidatorSetSize(&_ValidatorSet.CallOpts)
}

// ValidatorSetSize is a free data retrieval call binding the contract method 0xf800db50.
//
// Solidity: function validatorSetSize() constant returns(uint256)
func (_ValidatorSet *ValidatorSetCallerSession) ValidatorSetSize() (*big.Int, error) {
	return _ValidatorSet.Contract.ValidatorSetSize(&_ValidatorSet.CallOpts)
}

// SetNextVersionAddress is a paid mutator transaction binding the contract method 0xb2684966.
//
// Solidity: function _setNextVersionAddress(address _value) returns()
func (_ValidatorSet *ValidatorSetTransactor) SetNextVersionAddress(opts *bind.TransactOpts, _value common.Address) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "_setNextVersionAddress", _value)
}

// SetNextVersionAddress is a paid mutator transaction binding the contract method 0xb2684966.
//
// Solidity: function _setNextVersionAddress(address _value) returns()
func (_ValidatorSet *ValidatorSetSession) SetNextVersionAddress(_value common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.SetNextVersionAddress(&_ValidatorSet.TransactOpts, _value)
}

// SetNextVersionAddress is a paid mutator transaction binding the contract method 0xb2684966.
//
// Solidity: function _setNextVersionAddress(address _value) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) SetNextVersionAddress(_value common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.SetNextVersionAddress(&_ValidatorSet.TransactOpts, _value)
}

// AddValidator is a paid mutator transaction binding the contract method 0x3e8bb9a0.
//
// Solidity: function addValidator(string _pubKey, address _address) returns()
func (_ValidatorSet *ValidatorSetTransactor) AddValidator(opts *bind.TransactOpts, _pubKey string, _address common.Address) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "addValidator", _pubKey, _address)
}

// AddValidator is a paid mutator transaction binding the contract method 0x3e8bb9a0.
//
// Solidity: function addValidator(string _pubKey, address _address) returns()
func (_ValidatorSet *ValidatorSetSession) AddValidator(_pubKey string, _address common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.AddValidator(&_ValidatorSet.TransactOpts, _pubKey, _address)
}

// AddValidator is a paid mutator transaction binding the contract method 0x3e8bb9a0.
//
// Solidity: function addValidator(string _pubKey, address _address) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) AddValidator(_pubKey string, _address common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.AddValidator(&_ValidatorSet.TransactOpts, _pubKey, _address)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x0e62d9eb.
//
// Solidity: function removeValidator(string _pubKey, address _address) returns()
func (_ValidatorSet *ValidatorSetTransactor) RemoveValidator(opts *bind.TransactOpts, _pubKey string, _address common.Address) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "removeValidator", _pubKey, _address)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x0e62d9eb.
//
// Solidity: function removeValidator(string _pubKey, address _address) returns()
func (_ValidatorSet *ValidatorSetSession) RemoveValidator(_pubKey string, _address common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.RemoveValidator(&_ValidatorSet.TransactOpts, _pubKey, _address)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x0e62d9eb.
//
// Solidity: function removeValidator(string _pubKey, address _address) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) RemoveValidator(_pubKey string, _address common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.RemoveValidator(&_ValidatorSet.TransactOpts, _pubKey, _address)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ValidatorSet *ValidatorSetTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ValidatorSet *ValidatorSetSession) RenounceOwnership() (*types.Transaction, error) {
	return _ValidatorSet.Contract.RenounceOwnership(&_ValidatorSet.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ValidatorSet *ValidatorSetTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ValidatorSet.Contract.RenounceOwnership(&_ValidatorSet.TransactOpts)
}

// SetFreezeStatus is a paid mutator transaction binding the contract method 0x41edb775.
//
// Solidity: function setFreezeStatus(bool _value) returns()
func (_ValidatorSet *ValidatorSetTransactor) SetFreezeStatus(opts *bind.TransactOpts, _value bool) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "setFreezeStatus", _value)
}

// SetFreezeStatus is a paid mutator transaction binding the contract method 0x41edb775.
//
// Solidity: function setFreezeStatus(bool _value) returns()
func (_ValidatorSet *ValidatorSetSession) SetFreezeStatus(_value bool) (*types.Transaction, error) {
	return _ValidatorSet.Contract.SetFreezeStatus(&_ValidatorSet.TransactOpts, _value)
}

// SetFreezeStatus is a paid mutator transaction binding the contract method 0x41edb775.
//
// Solidity: function setFreezeStatus(bool _value) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) SetFreezeStatus(_value bool) (*types.Transaction, error) {
	return _ValidatorSet.Contract.SetFreezeStatus(&_ValidatorSet.TransactOpts, _value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ValidatorSet *ValidatorSetTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ValidatorSet *ValidatorSetSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.TransferOwnership(&_ValidatorSet.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.TransferOwnership(&_ValidatorSet.TransactOpts, newOwner)
}

// ValidatorSetFreezeIterator is returned from FilterFreeze and is used to iterate over the raw logs and unpacked data for Freeze events raised by the ValidatorSet contract.
type ValidatorSetFreezeIterator struct {
	Event *ValidatorSetFreeze // Event containing the contract specifics and raw log

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
func (it *ValidatorSetFreezeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorSetFreeze)
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
		it.Event = new(ValidatorSetFreeze)
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
func (it *ValidatorSetFreezeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorSetFreezeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorSetFreeze represents a Freeze event raised by the ValidatorSet contract.
type ValidatorSetFreeze struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterFreeze is a free log retrieval operation binding the contract event 0x615acbaede366d76a8b8cb2a9ada6a71495f0786513d71aa97aaf0c3910b78de.
//
// Solidity: event Freeze()
func (_ValidatorSet *ValidatorSetFilterer) FilterFreeze(opts *bind.FilterOpts) (*ValidatorSetFreezeIterator, error) {

	logs, sub, err := _ValidatorSet.contract.FilterLogs(opts, "Freeze")
	if err != nil {
		return nil, err
	}
	return &ValidatorSetFreezeIterator{contract: _ValidatorSet.contract, event: "Freeze", logs: logs, sub: sub}, nil
}

// WatchFreeze is a free log subscription operation binding the contract event 0x615acbaede366d76a8b8cb2a9ada6a71495f0786513d71aa97aaf0c3910b78de.
//
// Solidity: event Freeze()
func (_ValidatorSet *ValidatorSetFilterer) WatchFreeze(opts *bind.WatchOpts, sink chan<- *ValidatorSetFreeze) (event.Subscription, error) {

	logs, sub, err := _ValidatorSet.contract.WatchLogs(opts, "Freeze")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorSetFreeze)
				if err := _ValidatorSet.contract.UnpackLog(event, "Freeze", log); err != nil {
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

// ValidatorSetOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ValidatorSet contract.
type ValidatorSetOwnershipTransferredIterator struct {
	Event *ValidatorSetOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ValidatorSetOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorSetOwnershipTransferred)
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
		it.Event = new(ValidatorSetOwnershipTransferred)
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
func (it *ValidatorSetOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorSetOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorSetOwnershipTransferred represents a OwnershipTransferred event raised by the ValidatorSet contract.
type ValidatorSetOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ValidatorSet *ValidatorSetFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ValidatorSetOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ValidatorSet.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorSetOwnershipTransferredIterator{contract: _ValidatorSet.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ValidatorSet *ValidatorSetFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ValidatorSetOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ValidatorSet.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorSetOwnershipTransferred)
				if err := _ValidatorSet.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ValidatorSetSetNextVersionIterator is returned from FilterSetNextVersion and is used to iterate over the raw logs and unpacked data for SetNextVersion events raised by the ValidatorSet contract.
type ValidatorSetSetNextVersionIterator struct {
	Event *ValidatorSetSetNextVersion // Event containing the contract specifics and raw log

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
func (it *ValidatorSetSetNextVersionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorSetSetNextVersion)
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
		it.Event = new(ValidatorSetSetNextVersion)
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
func (it *ValidatorSetSetNextVersionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorSetSetNextVersionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorSetSetNextVersion represents a SetNextVersion event raised by the ValidatorSet contract.
type ValidatorSetSetNextVersion struct {
	Address common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSetNextVersion is a free log retrieval operation binding the contract event 0xb694c35dc1cdbab3df571f07de838161b83be5527eacfcb8a4b589b709438745.
//
// Solidity: event SetNextVersion(address _address)
func (_ValidatorSet *ValidatorSetFilterer) FilterSetNextVersion(opts *bind.FilterOpts) (*ValidatorSetSetNextVersionIterator, error) {

	logs, sub, err := _ValidatorSet.contract.FilterLogs(opts, "SetNextVersion")
	if err != nil {
		return nil, err
	}
	return &ValidatorSetSetNextVersionIterator{contract: _ValidatorSet.contract, event: "SetNextVersion", logs: logs, sub: sub}, nil
}

// WatchSetNextVersion is a free log subscription operation binding the contract event 0xb694c35dc1cdbab3df571f07de838161b83be5527eacfcb8a4b589b709438745.
//
// Solidity: event SetNextVersion(address _address)
func (_ValidatorSet *ValidatorSetFilterer) WatchSetNextVersion(opts *bind.WatchOpts, sink chan<- *ValidatorSetSetNextVersion) (event.Subscription, error) {

	logs, sub, err := _ValidatorSet.contract.WatchLogs(opts, "SetNextVersion")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorSetSetNextVersion)
				if err := _ValidatorSet.contract.UnpackLog(event, "SetNextVersion", log); err != nil {
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

// ValidatorSetValidatorAddedIterator is returned from FilterValidatorAdded and is used to iterate over the raw logs and unpacked data for ValidatorAdded events raised by the ValidatorSet contract.
type ValidatorSetValidatorAddedIterator struct {
	Event *ValidatorSetValidatorAdded // Event containing the contract specifics and raw log

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
func (it *ValidatorSetValidatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorSetValidatorAdded)
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
		it.Event = new(ValidatorSetValidatorAdded)
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
func (it *ValidatorSetValidatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorSetValidatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorSetValidatorAdded represents a ValidatorAdded event raised by the ValidatorSet contract.
type ValidatorSetValidatorAdded struct {
	Key     [32]byte
	Address common.Address
	PubKey  string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterValidatorAdded is a free log retrieval operation binding the contract event 0x940f4642cd572b4a8d9cf8545df0a213cbf8bee450f9ea7c781ede70650c4ab6.
//
// Solidity: event ValidatorAdded(bytes32 _key, address _address, string _pubKey)
func (_ValidatorSet *ValidatorSetFilterer) FilterValidatorAdded(opts *bind.FilterOpts) (*ValidatorSetValidatorAddedIterator, error) {

	logs, sub, err := _ValidatorSet.contract.FilterLogs(opts, "ValidatorAdded")
	if err != nil {
		return nil, err
	}
	return &ValidatorSetValidatorAddedIterator{contract: _ValidatorSet.contract, event: "ValidatorAdded", logs: logs, sub: sub}, nil
}

// WatchValidatorAdded is a free log subscription operation binding the contract event 0x940f4642cd572b4a8d9cf8545df0a213cbf8bee450f9ea7c781ede70650c4ab6.
//
// Solidity: event ValidatorAdded(bytes32 _key, address _address, string _pubKey)
func (_ValidatorSet *ValidatorSetFilterer) WatchValidatorAdded(opts *bind.WatchOpts, sink chan<- *ValidatorSetValidatorAdded) (event.Subscription, error) {

	logs, sub, err := _ValidatorSet.contract.WatchLogs(opts, "ValidatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorSetValidatorAdded)
				if err := _ValidatorSet.contract.UnpackLog(event, "ValidatorAdded", log); err != nil {
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

// ValidatorSetValidatorRemovedIterator is returned from FilterValidatorRemoved and is used to iterate over the raw logs and unpacked data for ValidatorRemoved events raised by the ValidatorSet contract.
type ValidatorSetValidatorRemovedIterator struct {
	Event *ValidatorSetValidatorRemoved // Event containing the contract specifics and raw log

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
func (it *ValidatorSetValidatorRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorSetValidatorRemoved)
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
		it.Event = new(ValidatorSetValidatorRemoved)
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
func (it *ValidatorSetValidatorRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorSetValidatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorSetValidatorRemoved represents a ValidatorRemoved event raised by the ValidatorSet contract.
type ValidatorSetValidatorRemoved struct {
	Key     [32]byte
	Address common.Address
	PubKey  string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterValidatorRemoved is a free log retrieval operation binding the contract event 0xf10b1758a748f201874627a3eb8553553531350314c7da7d9a13f2ae457e165d.
//
// Solidity: event ValidatorRemoved(bytes32 _key, address _address, string _pubKey)
func (_ValidatorSet *ValidatorSetFilterer) FilterValidatorRemoved(opts *bind.FilterOpts) (*ValidatorSetValidatorRemovedIterator, error) {

	logs, sub, err := _ValidatorSet.contract.FilterLogs(opts, "ValidatorRemoved")
	if err != nil {
		return nil, err
	}
	return &ValidatorSetValidatorRemovedIterator{contract: _ValidatorSet.contract, event: "ValidatorRemoved", logs: logs, sub: sub}, nil
}

// WatchValidatorRemoved is a free log subscription operation binding the contract event 0xf10b1758a748f201874627a3eb8553553531350314c7da7d9a13f2ae457e165d.
//
// Solidity: event ValidatorRemoved(bytes32 _key, address _address, string _pubKey)
func (_ValidatorSet *ValidatorSetFilterer) WatchValidatorRemoved(opts *bind.WatchOpts, sink chan<- *ValidatorSetValidatorRemoved) (event.Subscription, error) {

	logs, sub, err := _ValidatorSet.contract.WatchLogs(opts, "ValidatorRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorSetValidatorRemoved)
				if err := _ValidatorSet.contract.UnpackLog(event, "ValidatorRemoved", log); err != nil {
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
