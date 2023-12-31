// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package multicall

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
	_ = abi.ConvertType
)

// MulticallMetaData contains all meta data concerning the Multicall contract.
var MulticallMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"to\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"value\",\"type\":\"uint256[]\"}],\"name\":\"multiSigcall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b50610b0b8061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610028575f3560e01c8062d2475f1461002c575b5f80fd5b61004660048036038101906100419190610670565b610048565b005b5f5b82518110156100f2575f848483815181106100685761006761070c565b5b60200260200101518484815181106100835761008261070c565b5b602002602001015160405160200161009d93929190610757565b60405160208183030381529060405290506100e0866323b872dd60e01b836040516020016100cc929190610843565b6040516020818303038152906040526100f9565b50806100eb90610897565b905061004a565b5050505050565b5f610123828473ffffffffffffffffffffffffffffffffffffffff1661018e90919063ffffffff16565b90505f8151141580156101475750808060200190518101906101459190610913565b155b1561018957826040517f5274afe7000000000000000000000000000000000000000000000000000000008152600401610180919061093e565b60405180910390fd5b505050565b60606101d183835f6040518060400160405280601e81526020017f416464726573733a206c6f772d6c6576656c2063616c6c206661696c656400008152506101d9565b905092915050565b60608247101561021e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610215906109d7565b60405180910390fd5b5f808673ffffffffffffffffffffffffffffffffffffffff16858760405161024691906109f5565b5f6040518083038185875af1925050503d805f8114610280576040519150601f19603f3d011682016040523d82523d5f602084013e610285565b606091505b5091509150610296878383876102a2565b92505050949350505050565b60608315610303575f8351036102fb576102bb85610316565b6102fa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102f190610a55565b60405180910390fd5b5b82905061030e565b61030d8383610338565b5b949350505050565b5f808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b5f8251111561034a5781518083602001fd5b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161037e9190610ab5565b60405180910390fd5b5f604051905090565b5f80fd5b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6103c182610398565b9050919050565b5f6103d2826103b7565b9050919050565b6103e2816103c8565b81146103ec575f80fd5b50565b5f813590506103fd816103d9565b92915050565b61040c816103b7565b8114610416575f80fd5b50565b5f8135905061042781610403565b92915050565b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b61047782610431565b810181811067ffffffffffffffff8211171561049657610495610441565b5b80604052505050565b5f6104a8610387565b90506104b4828261046e565b919050565b5f67ffffffffffffffff8211156104d3576104d2610441565b5b602082029050602081019050919050565b5f80fd5b5f6104fa6104f5846104b9565b61049f565b9050808382526020820190506020840283018581111561051d5761051c6104e4565b5b835b8181101561054657806105328882610419565b84526020840193505060208101905061051f565b5050509392505050565b5f82601f8301126105645761056361042d565b5b81356105748482602086016104e8565b91505092915050565b5f67ffffffffffffffff82111561059757610596610441565b5b602082029050602081019050919050565b5f819050919050565b6105ba816105a8565b81146105c4575f80fd5b50565b5f813590506105d5816105b1565b92915050565b5f6105ed6105e88461057d565b61049f565b905080838252602082019050602084028301858111156106105761060f6104e4565b5b835b81811015610639578061062588826105c7565b845260208401935050602081019050610612565b5050509392505050565b5f82601f8301126106575761065661042d565b5b81356106678482602086016105db565b91505092915050565b5f805f806080858703121561068857610687610390565b5b5f610695878288016103ef565b94505060206106a687828801610419565b935050604085013567ffffffffffffffff8111156106c7576106c6610394565b5b6106d387828801610550565b925050606085013567ffffffffffffffff8111156106f4576106f3610394565b5b61070087828801610643565b91505092959194509250565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b610742816103b7565b82525050565b610751816105a8565b82525050565b5f60608201905061076a5f830186610739565b6107776020830185610739565b6107846040830184610748565b949350505050565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b5f819050919050565b6107d16107cc8261078c565b6107b7565b82525050565b5f81519050919050565b5f81905092915050565b5f5b838110156108085780820151818401526020810190506107ed565b5f8484015250505050565b5f61081d826107d7565b61082781856107e1565b93506108378185602086016107eb565b80840191505092915050565b5f61084e82856107c0565b60048201915061085e8284610813565b91508190509392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6108a1826105a8565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036108d3576108d261086a565b5b600182019050919050565b5f8115159050919050565b6108f2816108de565b81146108fc575f80fd5b50565b5f8151905061090d816108e9565b92915050565b5f6020828403121561092857610927610390565b5b5f610935848285016108ff565b91505092915050565b5f6020820190506109515f830184610739565b92915050565b5f82825260208201905092915050565b7f416464726573733a20696e73756666696369656e742062616c616e636520666f5f8201527f722063616c6c0000000000000000000000000000000000000000000000000000602082015250565b5f6109c1602683610957565b91506109cc82610967565b604082019050919050565b5f6020820190508181035f8301526109ee816109b5565b9050919050565b5f610a008284610813565b915081905092915050565b7f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000005f82015250565b5f610a3f601d83610957565b9150610a4a82610a0b565b602082019050919050565b5f6020820190508181035f830152610a6c81610a33565b9050919050565b5f81519050919050565b5f610a8782610a73565b610a918185610957565b9350610aa18185602086016107eb565b610aaa81610431565b840191505092915050565b5f6020820190508181035f830152610acd8184610a7d565b90509291505056fea2646970667358221220472ecb6d9e91c980bc84cc284bd7b693ffd2e236286b5468fc7f10f210567a5b64736f6c63430008140033",
}

// MulticallABI is the input ABI used to generate the binding from.
// Deprecated: Use MulticallMetaData.ABI instead.
var MulticallABI = MulticallMetaData.ABI

// MulticallBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MulticallMetaData.Bin instead.
var MulticallBin = MulticallMetaData.Bin

// DeployMulticall deploys a new Ethereum contract, binding an instance of Multicall to it.
func DeployMulticall(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Multicall, error) {
	parsed, err := MulticallMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MulticallBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Multicall{MulticallCaller: MulticallCaller{contract: contract}, MulticallTransactor: MulticallTransactor{contract: contract}, MulticallFilterer: MulticallFilterer{contract: contract}}, nil
}

// Multicall is an auto generated Go binding around an Ethereum contract.
type Multicall struct {
	MulticallCaller     // Read-only binding to the contract
	MulticallTransactor // Write-only binding to the contract
	MulticallFilterer   // Log filterer for contract events
}

// MulticallCaller is an auto generated read-only Go binding around an Ethereum contract.
type MulticallCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MulticallTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MulticallTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MulticallFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MulticallFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MulticallSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MulticallSession struct {
	Contract     *Multicall        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MulticallCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MulticallCallerSession struct {
	Contract *MulticallCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// MulticallTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MulticallTransactorSession struct {
	Contract     *MulticallTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// MulticallRaw is an auto generated low-level Go binding around an Ethereum contract.
type MulticallRaw struct {
	Contract *Multicall // Generic contract binding to access the raw methods on
}

// MulticallCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MulticallCallerRaw struct {
	Contract *MulticallCaller // Generic read-only contract binding to access the raw methods on
}

// MulticallTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MulticallTransactorRaw struct {
	Contract *MulticallTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMulticall creates a new instance of Multicall, bound to a specific deployed contract.
func NewMulticall(address common.Address, backend bind.ContractBackend) (*Multicall, error) {
	contract, err := bindMulticall(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Multicall{MulticallCaller: MulticallCaller{contract: contract}, MulticallTransactor: MulticallTransactor{contract: contract}, MulticallFilterer: MulticallFilterer{contract: contract}}, nil
}

// NewMulticallCaller creates a new read-only instance of Multicall, bound to a specific deployed contract.
func NewMulticallCaller(address common.Address, caller bind.ContractCaller) (*MulticallCaller, error) {
	contract, err := bindMulticall(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MulticallCaller{contract: contract}, nil
}

// NewMulticallTransactor creates a new write-only instance of Multicall, bound to a specific deployed contract.
func NewMulticallTransactor(address common.Address, transactor bind.ContractTransactor) (*MulticallTransactor, error) {
	contract, err := bindMulticall(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MulticallTransactor{contract: contract}, nil
}

// NewMulticallFilterer creates a new log filterer instance of Multicall, bound to a specific deployed contract.
func NewMulticallFilterer(address common.Address, filterer bind.ContractFilterer) (*MulticallFilterer, error) {
	contract, err := bindMulticall(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MulticallFilterer{contract: contract}, nil
}

// bindMulticall binds a generic wrapper to an already deployed contract.
func bindMulticall(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MulticallMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Multicall *MulticallRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multicall.Contract.MulticallCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Multicall *MulticallRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multicall.Contract.MulticallTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Multicall *MulticallRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multicall.Contract.MulticallTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Multicall *MulticallCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multicall.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Multicall *MulticallTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multicall.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Multicall *MulticallTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multicall.Contract.contract.Transact(opts, method, params...)
}

// MultiSigcall is a paid mutator transaction binding the contract method 0x00d2475f.
//
// Solidity: function multiSigcall(address token, address from, address[] to, uint256[] value) returns()
func (_Multicall *MulticallTransactor) MultiSigcall(opts *bind.TransactOpts, token common.Address, from common.Address, to []common.Address, value []*big.Int) (*types.Transaction, error) {
	return _Multicall.contract.Transact(opts, "multiSigcall", token, from, to, value)
}

// MultiSigcall is a paid mutator transaction binding the contract method 0x00d2475f.
//
// Solidity: function multiSigcall(address token, address from, address[] to, uint256[] value) returns()
func (_Multicall *MulticallSession) MultiSigcall(token common.Address, from common.Address, to []common.Address, value []*big.Int) (*types.Transaction, error) {
	return _Multicall.Contract.MultiSigcall(&_Multicall.TransactOpts, token, from, to, value)
}

// MultiSigcall is a paid mutator transaction binding the contract method 0x00d2475f.
//
// Solidity: function multiSigcall(address token, address from, address[] to, uint256[] value) returns()
func (_Multicall *MulticallTransactorSession) MultiSigcall(token common.Address, from common.Address, to []common.Address, value []*big.Int) (*types.Transaction, error) {
	return _Multicall.Contract.MultiSigcall(&_Multicall.TransactOpts, token, from, to, value)
}
