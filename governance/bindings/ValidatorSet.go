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
const ValidatorSetABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_vPubKeyAddresses\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes20\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"vPubKeyAddress\",\"type\":\"bytes20\"},{\"name\":\"vAddress\",\"type\":\"address\"}],\"name\":\"addValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_nextVersion\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setFreezeStatus\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_freeze\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"vPubKeyAddress\",\"type\":\"bytes20\"},{\"name\":\"vAddress\",\"type\":\"address\"}],\"name\":\"removeValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"address\"}],\"name\":\"_setNextVersionAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"validatorPubKey\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes20\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorSetSize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"vPubKeyAddress\",\"type\":\"bytes20\"}],\"name\":\"validatorAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"vPubKeyAddress\",\"type\":\"bytes20\"},{\"indexed\":false,\"name\":\"vAddress\",\"type\":\"address\"}],\"name\":\"ValidatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"vPubKeyAddress\",\"type\":\"bytes20\"},{\"indexed\":false,\"name\":\"vAddress\",\"type\":\"address\"}],\"name\":\"ValidatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"nAddress\",\"type\":\"address\"}],\"name\":\"SetNextVersion\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Freeze\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// ValidatorSetBin is the compiled bytecode used for deploying new contracts.
const ValidatorSetBin = `60806040526000600360006101000a81548160ff021916908315150217905550336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3611314806100ea6000396000f3fe6080604052600436106100d0576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806310860542146100d55780632112c5ed14610142578063369b4b96146101ac57806341edb7751461020357806365ed50e914610240578063707f6f941461026f578063715018a6146102d95780638da5cb5b146102f05780638f32d59b14610347578063b268496614610376578063f047a316146103c7578063f2fde38b14610434578063f800db5014610485578063ffd5ffbe146104b0575b600080fd5b3480156100e157600080fd5b5061010e600480360360208110156100f857600080fd5b810190808035906020019092919050505061053a565b60405180826bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200191505060405180910390f35b34801561014e57600080fd5b506101aa6004803603604081101561016557600080fd5b8101908080356bffffffffffffffffffffffff19169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610571565b005b3480156101b857600080fd5b506101c1610834565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561020f57600080fd5b5061023e6004803603602081101561022657600080fd5b8101908080351515906020019092919050505061085a565b005b34801561024c57600080fd5b5061025561091f565b604051808215151515815260200191505060405180910390f35b34801561027b57600080fd5b506102d76004803603604081101561029257600080fd5b8101908080356bffffffffffffffffffffffff19169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610932565b005b3480156102e557600080fd5b506102ee610d04565b005b3480156102fc57600080fd5b50610305610e3f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561035357600080fd5b5061035c610e68565b604051808215151515815260200191505060405180910390f35b34801561038257600080fd5b506103c56004803603602081101561039957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610ebf565b005b3480156103d357600080fd5b50610400600480360360208110156103ea57600080fd5b8101908080359060200190929190505050610fe2565b60405180826bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200191505060405180910390f35b34801561044057600080fd5b506104836004803603602081101561045757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061101e565b005b34801561049157600080fd5b5061049a6110a6565b6040518082815260200191505060405180910390f35b3480156104bc57600080fd5b506104f8600480360360208110156104d357600080fd5b8101908080356bffffffffffffffffffffffff191690602001909291905050506110b3565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60028181548110151561054957fe5b906000526020600020016000915054906101000a90046c010000000000000000000000000281565b610579610e68565b15156105ed576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b60001515600360009054906101000a900460ff16151514151561060f57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415151561064b57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff1660016000846bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415156106d757600080fd5b8060016000846bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060028290806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff02191690836c0100000000000000000000000090040217905550507fd46e7b5750f59856e49d3251a33178e173a5bcb73801c4c7a80f870d961f7c02828260405180836bffffffffffffffffffffffff19166bffffffffffffffffffffffff191681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a15050565b600360019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b610862610e68565b15156108d6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b80600360006101000a81548160ff0219169083151502179055507f615acbaede366d76a8b8cb2a9ada6a71495f0786513d71aa97aaf0c3910b78de60405160405180910390a150565b600360009054906101000a900460ff1681565b61093a610e68565b15156109ae576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b60001515600360009054906101000a900460ff1615151415156109d057600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515610a0c57600080fd5b8073ffffffffffffffffffffffffffffffffffffffff1660016000846bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141515610a9757600080fd5b600060016000846bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600080905060008090505b600160028054905003811015610c1b57836bffffffffffffffffffffffff1916600282815481101515610b4257fe5b9060005260206000200160009054906101000a90046c01000000000000000000000000026bffffffffffffffffffffffff19161415610b8057600191505b8115610c0e57600260018201815481101515610b9857fe5b9060005260206000200160009054906101000a90046c0100000000000000000000000002600282815481101515610bcb57fe5b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff02191690836c01000000000000000000000000900402179055505b8080600101915050610b13565b506002600160028054905003815481101515610c3357fe5b9060005260206000200160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556002805480919060019003610c759190611297565b507f173d885d15691a4b17a9fe2da8efd3339c47991bbf4a76af4fc345c433e59f38838360405180836bffffffffffffffffffffffff19166bffffffffffffffffffffffff191681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a1505050565b610d0c610e68565b1515610d80576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a360008060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614905090565b610ec7610e68565b1515610f3b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b80600360016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507fb694c35dc1cdbab3df571f07de838161b83be5527eacfcb8a4b589b70943874581604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a150565b6000600282815481101515610ff357fe5b9060005260206000200160009054906101000a90046c01000000000000000000000000029050919050565b611026610e68565b151561109a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b6110a38161110e565b50565b6000600280549050905090565b600060016000836bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141515156111d9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001807f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206181526020017f646472657373000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b8154818355818111156112be578183600052602060002091820191016112bd91906112c3565b5b505050565b6112e591905b808211156112e15760008160009055506001016112c9565b5090565b9056fea165627a7a72305820dc4dabfbd294a0c3624b2d592da44dacc06d9082be610dabda7b61e4fcbcfd060029`

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

// VPubKeyAddresses is a free data retrieval call binding the contract method 0x10860542.
//
// Solidity: function _vPubKeyAddresses(uint256 ) constant returns(bytes20)
func (_ValidatorSet *ValidatorSetCaller) VPubKeyAddresses(opts *bind.CallOpts, arg0 *big.Int) ([20]byte, error) {
	var (
		ret0 = new([20]byte)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "_vPubKeyAddresses", arg0)
	return *ret0, err
}

// VPubKeyAddresses is a free data retrieval call binding the contract method 0x10860542.
//
// Solidity: function _vPubKeyAddresses(uint256 ) constant returns(bytes20)
func (_ValidatorSet *ValidatorSetSession) VPubKeyAddresses(arg0 *big.Int) ([20]byte, error) {
	return _ValidatorSet.Contract.VPubKeyAddresses(&_ValidatorSet.CallOpts, arg0)
}

// VPubKeyAddresses is a free data retrieval call binding the contract method 0x10860542.
//
// Solidity: function _vPubKeyAddresses(uint256 ) constant returns(bytes20)
func (_ValidatorSet *ValidatorSetCallerSession) VPubKeyAddresses(arg0 *big.Int) ([20]byte, error) {
	return _ValidatorSet.Contract.VPubKeyAddresses(&_ValidatorSet.CallOpts, arg0)
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

// ValidatorAddress is a free data retrieval call binding the contract method 0xffd5ffbe.
//
// Solidity: function validatorAddress(bytes20 vPubKeyAddress) constant returns(address)
func (_ValidatorSet *ValidatorSetCaller) ValidatorAddress(opts *bind.CallOpts, vPubKeyAddress [20]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "validatorAddress", vPubKeyAddress)
	return *ret0, err
}

// ValidatorAddress is a free data retrieval call binding the contract method 0xffd5ffbe.
//
// Solidity: function validatorAddress(bytes20 vPubKeyAddress) constant returns(address)
func (_ValidatorSet *ValidatorSetSession) ValidatorAddress(vPubKeyAddress [20]byte) (common.Address, error) {
	return _ValidatorSet.Contract.ValidatorAddress(&_ValidatorSet.CallOpts, vPubKeyAddress)
}

// ValidatorAddress is a free data retrieval call binding the contract method 0xffd5ffbe.
//
// Solidity: function validatorAddress(bytes20 vPubKeyAddress) constant returns(address)
func (_ValidatorSet *ValidatorSetCallerSession) ValidatorAddress(vPubKeyAddress [20]byte) (common.Address, error) {
	return _ValidatorSet.Contract.ValidatorAddress(&_ValidatorSet.CallOpts, vPubKeyAddress)
}

// ValidatorPubKey is a free data retrieval call binding the contract method 0xf047a316.
//
// Solidity: function validatorPubKey(uint256 index) constant returns(bytes20)
func (_ValidatorSet *ValidatorSetCaller) ValidatorPubKey(opts *bind.CallOpts, index *big.Int) ([20]byte, error) {
	var (
		ret0 = new([20]byte)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "validatorPubKey", index)
	return *ret0, err
}

// ValidatorPubKey is a free data retrieval call binding the contract method 0xf047a316.
//
// Solidity: function validatorPubKey(uint256 index) constant returns(bytes20)
func (_ValidatorSet *ValidatorSetSession) ValidatorPubKey(index *big.Int) ([20]byte, error) {
	return _ValidatorSet.Contract.ValidatorPubKey(&_ValidatorSet.CallOpts, index)
}

// ValidatorPubKey is a free data retrieval call binding the contract method 0xf047a316.
//
// Solidity: function validatorPubKey(uint256 index) constant returns(bytes20)
func (_ValidatorSet *ValidatorSetCallerSession) ValidatorPubKey(index *big.Int) ([20]byte, error) {
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
// Solidity: function _setNextVersionAddress(address value) returns()
func (_ValidatorSet *ValidatorSetTransactor) SetNextVersionAddress(opts *bind.TransactOpts, value common.Address) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "_setNextVersionAddress", value)
}

// SetNextVersionAddress is a paid mutator transaction binding the contract method 0xb2684966.
//
// Solidity: function _setNextVersionAddress(address value) returns()
func (_ValidatorSet *ValidatorSetSession) SetNextVersionAddress(value common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.SetNextVersionAddress(&_ValidatorSet.TransactOpts, value)
}

// SetNextVersionAddress is a paid mutator transaction binding the contract method 0xb2684966.
//
// Solidity: function _setNextVersionAddress(address value) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) SetNextVersionAddress(value common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.SetNextVersionAddress(&_ValidatorSet.TransactOpts, value)
}

// AddValidator is a paid mutator transaction binding the contract method 0x2112c5ed.
//
// Solidity: function addValidator(bytes20 vPubKeyAddress, address vAddress) returns()
func (_ValidatorSet *ValidatorSetTransactor) AddValidator(opts *bind.TransactOpts, vPubKeyAddress [20]byte, vAddress common.Address) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "addValidator", vPubKeyAddress, vAddress)
}

// AddValidator is a paid mutator transaction binding the contract method 0x2112c5ed.
//
// Solidity: function addValidator(bytes20 vPubKeyAddress, address vAddress) returns()
func (_ValidatorSet *ValidatorSetSession) AddValidator(vPubKeyAddress [20]byte, vAddress common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.AddValidator(&_ValidatorSet.TransactOpts, vPubKeyAddress, vAddress)
}

// AddValidator is a paid mutator transaction binding the contract method 0x2112c5ed.
//
// Solidity: function addValidator(bytes20 vPubKeyAddress, address vAddress) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) AddValidator(vPubKeyAddress [20]byte, vAddress common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.AddValidator(&_ValidatorSet.TransactOpts, vPubKeyAddress, vAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x707f6f94.
//
// Solidity: function removeValidator(bytes20 vPubKeyAddress, address vAddress) returns()
func (_ValidatorSet *ValidatorSetTransactor) RemoveValidator(opts *bind.TransactOpts, vPubKeyAddress [20]byte, vAddress common.Address) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "removeValidator", vPubKeyAddress, vAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x707f6f94.
//
// Solidity: function removeValidator(bytes20 vPubKeyAddress, address vAddress) returns()
func (_ValidatorSet *ValidatorSetSession) RemoveValidator(vPubKeyAddress [20]byte, vAddress common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.RemoveValidator(&_ValidatorSet.TransactOpts, vPubKeyAddress, vAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x707f6f94.
//
// Solidity: function removeValidator(bytes20 vPubKeyAddress, address vAddress) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) RemoveValidator(vPubKeyAddress [20]byte, vAddress common.Address) (*types.Transaction, error) {
	return _ValidatorSet.Contract.RemoveValidator(&_ValidatorSet.TransactOpts, vPubKeyAddress, vAddress)
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
// Solidity: function setFreezeStatus(bool value) returns()
func (_ValidatorSet *ValidatorSetTransactor) SetFreezeStatus(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "setFreezeStatus", value)
}

// SetFreezeStatus is a paid mutator transaction binding the contract method 0x41edb775.
//
// Solidity: function setFreezeStatus(bool value) returns()
func (_ValidatorSet *ValidatorSetSession) SetFreezeStatus(value bool) (*types.Transaction, error) {
	return _ValidatorSet.Contract.SetFreezeStatus(&_ValidatorSet.TransactOpts, value)
}

// SetFreezeStatus is a paid mutator transaction binding the contract method 0x41edb775.
//
// Solidity: function setFreezeStatus(bool value) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) SetFreezeStatus(value bool) (*types.Transaction, error) {
	return _ValidatorSet.Contract.SetFreezeStatus(&_ValidatorSet.TransactOpts, value)
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
	NAddress common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetNextVersion is a free log retrieval operation binding the contract event 0xb694c35dc1cdbab3df571f07de838161b83be5527eacfcb8a4b589b709438745.
//
// Solidity: event SetNextVersion(address nAddress)
func (_ValidatorSet *ValidatorSetFilterer) FilterSetNextVersion(opts *bind.FilterOpts) (*ValidatorSetSetNextVersionIterator, error) {

	logs, sub, err := _ValidatorSet.contract.FilterLogs(opts, "SetNextVersion")
	if err != nil {
		return nil, err
	}
	return &ValidatorSetSetNextVersionIterator{contract: _ValidatorSet.contract, event: "SetNextVersion", logs: logs, sub: sub}, nil
}

// WatchSetNextVersion is a free log subscription operation binding the contract event 0xb694c35dc1cdbab3df571f07de838161b83be5527eacfcb8a4b589b709438745.
//
// Solidity: event SetNextVersion(address nAddress)
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
	VPubKeyAddress [20]byte
	VAddress       common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterValidatorAdded is a free log retrieval operation binding the contract event 0xd46e7b5750f59856e49d3251a33178e173a5bcb73801c4c7a80f870d961f7c02.
//
// Solidity: event ValidatorAdded(bytes20 vPubKeyAddress, address vAddress)
func (_ValidatorSet *ValidatorSetFilterer) FilterValidatorAdded(opts *bind.FilterOpts) (*ValidatorSetValidatorAddedIterator, error) {

	logs, sub, err := _ValidatorSet.contract.FilterLogs(opts, "ValidatorAdded")
	if err != nil {
		return nil, err
	}
	return &ValidatorSetValidatorAddedIterator{contract: _ValidatorSet.contract, event: "ValidatorAdded", logs: logs, sub: sub}, nil
}

// WatchValidatorAdded is a free log subscription operation binding the contract event 0xd46e7b5750f59856e49d3251a33178e173a5bcb73801c4c7a80f870d961f7c02.
//
// Solidity: event ValidatorAdded(bytes20 vPubKeyAddress, address vAddress)
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
	VPubKeyAddress [20]byte
	VAddress       common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterValidatorRemoved is a free log retrieval operation binding the contract event 0x173d885d15691a4b17a9fe2da8efd3339c47991bbf4a76af4fc345c433e59f38.
//
// Solidity: event ValidatorRemoved(bytes20 vPubKeyAddress, address vAddress)
func (_ValidatorSet *ValidatorSetFilterer) FilterValidatorRemoved(opts *bind.FilterOpts) (*ValidatorSetValidatorRemovedIterator, error) {

	logs, sub, err := _ValidatorSet.contract.FilterLogs(opts, "ValidatorRemoved")
	if err != nil {
		return nil, err
	}
	return &ValidatorSetValidatorRemovedIterator{contract: _ValidatorSet.contract, event: "ValidatorRemoved", logs: logs, sub: sub}, nil
}

// WatchValidatorRemoved is a free log subscription operation binding the contract event 0x173d885d15691a4b17a9fe2da8efd3339c47991bbf4a76af4fc345c433e59f38.
//
// Solidity: event ValidatorRemoved(bytes20 vPubKeyAddress, address vAddress)
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
