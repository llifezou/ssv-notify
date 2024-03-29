// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package liquidation

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ISSVNetworkCoreCluster is an auto generated low-level Go binding around an user-defined struct.
type ISSVNetworkCoreCluster struct {
	ValidatorCount  uint32
	NetworkFeeIndex uint64
	Index           uint64
	Active          bool
	Balance         *big.Int
}

// LiquidationMetaData contains all meta data concerning the Liquidation contract.
var LiquidationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint64[]\",\"name\":\"operatorIds\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"validatorCount\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"networkFeeIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structISSVNetworkCore.Cluster\",\"name\":\"cluster\",\"type\":\"tuple\"}],\"name\":\"liquidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint64[]\",\"name\":\"operatorIds\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"validatorCount\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"networkFeeIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structISSVNetworkCore.Cluster\",\"name\":\"cluster\",\"type\":\"tuple\"}],\"name\":\"isLiquidatable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// LiquidationABI is the input ABI used to generate the binding from.
// Deprecated: Use LiquidationMetaData.ABI instead.
var LiquidationABI = LiquidationMetaData.ABI

// Liquidation is an auto generated Go binding around an Ethereum contract.
type Liquidation struct {
	LiquidationCaller     // Read-only binding to the contract
	LiquidationTransactor // Write-only binding to the contract
	LiquidationFilterer   // Log filterer for contract events
}

// LiquidationCaller is an auto generated read-only Go binding around an Ethereum contract.
type LiquidationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LiquidationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LiquidationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LiquidationSession struct {
	Contract     *Liquidation      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LiquidationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LiquidationCallerSession struct {
	Contract *LiquidationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// LiquidationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LiquidationTransactorSession struct {
	Contract     *LiquidationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// LiquidationRaw is an auto generated low-level Go binding around an Ethereum contract.
type LiquidationRaw struct {
	Contract *Liquidation // Generic contract binding to access the raw methods on
}

// LiquidationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LiquidationCallerRaw struct {
	Contract *LiquidationCaller // Generic read-only contract binding to access the raw methods on
}

// LiquidationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LiquidationTransactorRaw struct {
	Contract *LiquidationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLiquidation creates a new instance of Liquidation, bound to a specific deployed contract.
func NewLiquidation(address common.Address, backend bind.ContractBackend) (*Liquidation, error) {
	contract, err := bindLiquidation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Liquidation{LiquidationCaller: LiquidationCaller{contract: contract}, LiquidationTransactor: LiquidationTransactor{contract: contract}, LiquidationFilterer: LiquidationFilterer{contract: contract}}, nil
}

// NewLiquidationCaller creates a new read-only instance of Liquidation, bound to a specific deployed contract.
func NewLiquidationCaller(address common.Address, caller bind.ContractCaller) (*LiquidationCaller, error) {
	contract, err := bindLiquidation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidationCaller{contract: contract}, nil
}

// NewLiquidationTransactor creates a new write-only instance of Liquidation, bound to a specific deployed contract.
func NewLiquidationTransactor(address common.Address, transactor bind.ContractTransactor) (*LiquidationTransactor, error) {
	contract, err := bindLiquidation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidationTransactor{contract: contract}, nil
}

// NewLiquidationFilterer creates a new log filterer instance of Liquidation, bound to a specific deployed contract.
func NewLiquidationFilterer(address common.Address, filterer bind.ContractFilterer) (*LiquidationFilterer, error) {
	contract, err := bindLiquidation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LiquidationFilterer{contract: contract}, nil
}

// bindLiquidation binds a generic wrapper to an already deployed contract.
func bindLiquidation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LiquidationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Liquidation *LiquidationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Liquidation.Contract.LiquidationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Liquidation *LiquidationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Liquidation.Contract.LiquidationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Liquidation *LiquidationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Liquidation.Contract.LiquidationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Liquidation *LiquidationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Liquidation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Liquidation *LiquidationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Liquidation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Liquidation *LiquidationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Liquidation.Contract.contract.Transact(opts, method, params...)
}

// IsLiquidatable is a free data retrieval call binding the contract method 0x16cff008.
//
// Solidity: function isLiquidatable(address owner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(bool)
func (_Liquidation *LiquidationCaller) IsLiquidatable(opts *bind.CallOpts, owner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (bool, error) {
	var out []interface{}
	err := _Liquidation.contract.Call(opts, &out, "isLiquidatable", owner, operatorIds, cluster)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsLiquidatable is a free data retrieval call binding the contract method 0x16cff008.
//
// Solidity: function isLiquidatable(address owner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(bool)
func (_Liquidation *LiquidationSession) IsLiquidatable(owner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (bool, error) {
	return _Liquidation.Contract.IsLiquidatable(&_Liquidation.CallOpts, owner, operatorIds, cluster)
}

// IsLiquidatable is a free data retrieval call binding the contract method 0x16cff008.
//
// Solidity: function isLiquidatable(address owner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(bool)
func (_Liquidation *LiquidationCallerSession) IsLiquidatable(owner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (bool, error) {
	return _Liquidation.Contract.IsLiquidatable(&_Liquidation.CallOpts, owner, operatorIds, cluster)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbf0f2fb2.
//
// Solidity: function liquidate(address owner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) returns()
func (_Liquidation *LiquidationTransactor) Liquidate(opts *bind.TransactOpts, owner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (*types.Transaction, error) {
	return _Liquidation.contract.Transact(opts, "liquidate", owner, operatorIds, cluster)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbf0f2fb2.
//
// Solidity: function liquidate(address owner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) returns()
func (_Liquidation *LiquidationSession) Liquidate(owner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (*types.Transaction, error) {
	return _Liquidation.Contract.Liquidate(&_Liquidation.TransactOpts, owner, operatorIds, cluster)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbf0f2fb2.
//
// Solidity: function liquidate(address owner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) returns()
func (_Liquidation *LiquidationTransactorSession) Liquidate(owner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (*types.Transaction, error) {
	return _Liquidation.Contract.Liquidate(&_Liquidation.TransactOpts, owner, operatorIds, cluster)
}
