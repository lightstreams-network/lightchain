// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package distribution

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

// DistributionABI is the input ABI used to generate the binding from.
const DistributionABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_beneficiary\",\"type\":\"address\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_oldBeneficiary\",\"type\":\"address\"},{\"name\":\"_newBeneficiary\",\"type\":\"address\"}],\"name\":\"changeDepositBeneficiary\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"}]"

// DistributionBin is the compiled bytecode used for deploying new contracts.
const DistributionBin = `0x608060405234801561001057600080fd5b50600080546001600160a01b031916331790556102f0806100326000396000f3fe60806040526004361061003f5760003560e01c80633ccfd60b146100415780638da5cb5b14610049578063f340fa011461007a578063f5897c95146100a0575b005b61003f6100db565b34801561005557600080fd5b5061005e610139565b604080516001600160a01b039092168252519081900360200190f35b61003f6004803603602081101561009057600080fd5b50356001600160a01b0316610148565b3480156100ac57600080fd5b5061003f600480360360408110156100c357600080fd5b506001600160a01b03813581169160200135166101f3565b336000908152600160205260409020546100f457600080fd5b33600081815260016020526040808220805490839055905190929183156108fc02918491818181858888f19350505050158015610135573d6000803e3d6000fd5b5050565b6000546001600160a01b031681565b6000546001600160a01b0316331461019457604051600160e51b62461bcd0281526004018080602001828103825260228152602001806102a36022913960400191505060405180910390fd5b6001600160a01b0381166101a757600080fd5b6001600160a01b038116600090815260016020526040902054156101ca57600080fd5b600034116101d757600080fd5b6001600160a01b03166000908152600160205260409020349055565b6000546001600160a01b0316331461023f57604051600160e51b62461bcd0281526004018080602001828103825260228152602001806102a36022913960400191505060405180910390fd5b6001600160a01b03811661025257600080fd5b6001600160a01b0381166000908152600160205260409020541561027557600080fd5b6001600160a01b03918216600090815260016020526040808220805490839055929093168152919091205556fe4f6e6c79206f776e65722063616e2063616c6c20746869732066756e6374696f6e2ea165627a7a723058202df16dc5c51c157597f87259d5861698338d0a3675c57241f524f30c57a261a50029`

// DeployDistribution deploys a new Ethereum contract, binding an instance of Distribution to it.
func DeployDistribution(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Distribution, error) {
	parsed, err := abi.JSON(strings.NewReader(DistributionABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DistributionBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Distribution{DistributionCaller: DistributionCaller{contract: contract}, DistributionTransactor: DistributionTransactor{contract: contract}, DistributionFilterer: DistributionFilterer{contract: contract}}, nil
}

// Distribution is an auto generated Go binding around an Ethereum contract.
type Distribution struct {
	DistributionCaller     // Read-only binding to the contract
	DistributionTransactor // Write-only binding to the contract
	DistributionFilterer   // Log filterer for contract events
}

// DistributionCaller is an auto generated read-only Go binding around an Ethereum contract.
type DistributionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DistributionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DistributionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DistributionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DistributionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DistributionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DistributionSession struct {
	Contract     *Distribution     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DistributionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DistributionCallerSession struct {
	Contract *DistributionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// DistributionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DistributionTransactorSession struct {
	Contract     *DistributionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// DistributionRaw is an auto generated low-level Go binding around an Ethereum contract.
type DistributionRaw struct {
	Contract *Distribution // Generic contract binding to access the raw methods on
}

// DistributionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DistributionCallerRaw struct {
	Contract *DistributionCaller // Generic read-only contract binding to access the raw methods on
}

// DistributionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DistributionTransactorRaw struct {
	Contract *DistributionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDistribution creates a new instance of Distribution, bound to a specific deployed contract.
func NewDistribution(address common.Address, backend bind.ContractBackend) (*Distribution, error) {
	contract, err := bindDistribution(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Distribution{DistributionCaller: DistributionCaller{contract: contract}, DistributionTransactor: DistributionTransactor{contract: contract}, DistributionFilterer: DistributionFilterer{contract: contract}}, nil
}

// NewDistributionCaller creates a new read-only instance of Distribution, bound to a specific deployed contract.
func NewDistributionCaller(address common.Address, caller bind.ContractCaller) (*DistributionCaller, error) {
	contract, err := bindDistribution(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DistributionCaller{contract: contract}, nil
}

// NewDistributionTransactor creates a new write-only instance of Distribution, bound to a specific deployed contract.
func NewDistributionTransactor(address common.Address, transactor bind.ContractTransactor) (*DistributionTransactor, error) {
	contract, err := bindDistribution(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DistributionTransactor{contract: contract}, nil
}

// NewDistributionFilterer creates a new log filterer instance of Distribution, bound to a specific deployed contract.
func NewDistributionFilterer(address common.Address, filterer bind.ContractFilterer) (*DistributionFilterer, error) {
	contract, err := bindDistribution(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DistributionFilterer{contract: contract}, nil
}

// bindDistribution binds a generic wrapper to an already deployed contract.
func bindDistribution(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DistributionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Distribution *DistributionRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Distribution.Contract.DistributionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Distribution *DistributionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Distribution.Contract.DistributionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Distribution *DistributionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Distribution.Contract.DistributionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Distribution *DistributionCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Distribution.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Distribution *DistributionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Distribution.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Distribution *DistributionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Distribution.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Distribution *DistributionCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Distribution.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Distribution *DistributionSession) Owner() (common.Address, error) {
	return _Distribution.Contract.Owner(&_Distribution.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Distribution *DistributionCallerSession) Owner() (common.Address, error) {
	return _Distribution.Contract.Owner(&_Distribution.CallOpts)
}

// ChangeDepositBeneficiary is a paid mutator transaction binding the contract method 0xf5897c95.
//
// Solidity: function changeDepositBeneficiary(address _oldBeneficiary, address _newBeneficiary) returns()
func (_Distribution *DistributionTransactor) ChangeDepositBeneficiary(opts *bind.TransactOpts, _oldBeneficiary common.Address, _newBeneficiary common.Address) (*types.Transaction, error) {
	return _Distribution.contract.Transact(opts, "changeDepositBeneficiary", _oldBeneficiary, _newBeneficiary)
}

// ChangeDepositBeneficiary is a paid mutator transaction binding the contract method 0xf5897c95.
//
// Solidity: function changeDepositBeneficiary(address _oldBeneficiary, address _newBeneficiary) returns()
func (_Distribution *DistributionSession) ChangeDepositBeneficiary(_oldBeneficiary common.Address, _newBeneficiary common.Address) (*types.Transaction, error) {
	return _Distribution.Contract.ChangeDepositBeneficiary(&_Distribution.TransactOpts, _oldBeneficiary, _newBeneficiary)
}

// ChangeDepositBeneficiary is a paid mutator transaction binding the contract method 0xf5897c95.
//
// Solidity: function changeDepositBeneficiary(address _oldBeneficiary, address _newBeneficiary) returns()
func (_Distribution *DistributionTransactorSession) ChangeDepositBeneficiary(_oldBeneficiary common.Address, _newBeneficiary common.Address) (*types.Transaction, error) {
	return _Distribution.Contract.ChangeDepositBeneficiary(&_Distribution.TransactOpts, _oldBeneficiary, _newBeneficiary)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _beneficiary) returns()
func (_Distribution *DistributionTransactor) Deposit(opts *bind.TransactOpts, _beneficiary common.Address) (*types.Transaction, error) {
	return _Distribution.contract.Transact(opts, "deposit", _beneficiary)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _beneficiary) returns()
func (_Distribution *DistributionSession) Deposit(_beneficiary common.Address) (*types.Transaction, error) {
	return _Distribution.Contract.Deposit(&_Distribution.TransactOpts, _beneficiary)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _beneficiary) returns()
func (_Distribution *DistributionTransactorSession) Deposit(_beneficiary common.Address) (*types.Transaction, error) {
	return _Distribution.Contract.Deposit(&_Distribution.TransactOpts, _beneficiary)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Distribution *DistributionTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Distribution.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Distribution *DistributionSession) Withdraw() (*types.Transaction, error) {
	return _Distribution.Contract.Withdraw(&_Distribution.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Distribution *DistributionTransactorSession) Withdraw() (*types.Transaction, error) {
	return _Distribution.Contract.Withdraw(&_Distribution.TransactOpts)
}
