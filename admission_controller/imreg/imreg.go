// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package imreg

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

// ImageRegistryMetaData contains all meta data concerning the ImageRegistry contract.
var ImageRegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"imageHash\",\"type\":\"bytes32\"}],\"name\":\"ImageRegistered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"imageHash\",\"type\":\"bytes32\"}],\"name\":\"isImageTrusted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"imageHash\",\"type\":\"bytes32\"}],\"name\":\"registerImage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"trustedImages\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b5061012f8061001c5f395ff3fe6080604052348015600e575f5ffd5b5060043610603a575f3560e01c80633f633adb14603e57806343f9405f14604f5780634ee9f2d1146082575b5f5ffd5b604d604936600460e3565b60a1565b005b606e605a36600460e3565b5f9081526020819052604090205460ff1690565b604051901515815260200160405180910390f35b606e608d36600460e3565b5f6020819052908152604090205460ff1681565b5f81815260208190526040808220805460ff191660011790555182917fbb218b24bde6ba405fba71b6e1e123ab350a965dc6ad8fc3e53b3cc00754785891a250565b5f6020828403121560f2575f5ffd5b503591905056fea2646970667358221220543d10cd6bd32dc5f3d7d132a570c3868d3e6e5483c2f681c9ab728d74b9611064736f6c634300081d0033",
}

// ImageRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ImageRegistryMetaData.ABI instead.
var ImageRegistryABI = ImageRegistryMetaData.ABI

// ImageRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ImageRegistryMetaData.Bin instead.
var ImageRegistryBin = ImageRegistryMetaData.Bin

// DeployImageRegistry deploys a new Ethereum contract, binding an instance of ImageRegistry to it.
func DeployImageRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ImageRegistry, error) {
	parsed, err := ImageRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ImageRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ImageRegistry{ImageRegistryCaller: ImageRegistryCaller{contract: contract}, ImageRegistryTransactor: ImageRegistryTransactor{contract: contract}, ImageRegistryFilterer: ImageRegistryFilterer{contract: contract}}, nil
}

// ImageRegistry is an auto generated Go binding around an Ethereum contract.
type ImageRegistry struct {
	ImageRegistryCaller     // Read-only binding to the contract
	ImageRegistryTransactor // Write-only binding to the contract
	ImageRegistryFilterer   // Log filterer for contract events
}

// ImageRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ImageRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ImageRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ImageRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ImageRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ImageRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ImageRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ImageRegistrySession struct {
	Contract     *ImageRegistry    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ImageRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ImageRegistryCallerSession struct {
	Contract *ImageRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// ImageRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ImageRegistryTransactorSession struct {
	Contract     *ImageRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ImageRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ImageRegistryRaw struct {
	Contract *ImageRegistry // Generic contract binding to access the raw methods on
}

// ImageRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ImageRegistryCallerRaw struct {
	Contract *ImageRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ImageRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ImageRegistryTransactorRaw struct {
	Contract *ImageRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewImageRegistry creates a new instance of ImageRegistry, bound to a specific deployed contract.
func NewImageRegistry(address common.Address, backend bind.ContractBackend) (*ImageRegistry, error) {
	contract, err := bindImageRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ImageRegistry{ImageRegistryCaller: ImageRegistryCaller{contract: contract}, ImageRegistryTransactor: ImageRegistryTransactor{contract: contract}, ImageRegistryFilterer: ImageRegistryFilterer{contract: contract}}, nil
}

// NewImageRegistryCaller creates a new read-only instance of ImageRegistry, bound to a specific deployed contract.
func NewImageRegistryCaller(address common.Address, caller bind.ContractCaller) (*ImageRegistryCaller, error) {
	contract, err := bindImageRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ImageRegistryCaller{contract: contract}, nil
}

// NewImageRegistryTransactor creates a new write-only instance of ImageRegistry, bound to a specific deployed contract.
func NewImageRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ImageRegistryTransactor, error) {
	contract, err := bindImageRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ImageRegistryTransactor{contract: contract}, nil
}

// NewImageRegistryFilterer creates a new log filterer instance of ImageRegistry, bound to a specific deployed contract.
func NewImageRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ImageRegistryFilterer, error) {
	contract, err := bindImageRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ImageRegistryFilterer{contract: contract}, nil
}

// bindImageRegistry binds a generic wrapper to an already deployed contract.
func bindImageRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ImageRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ImageRegistry *ImageRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ImageRegistry.Contract.ImageRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ImageRegistry *ImageRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ImageRegistry.Contract.ImageRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ImageRegistry *ImageRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ImageRegistry.Contract.ImageRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ImageRegistry *ImageRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ImageRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ImageRegistry *ImageRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ImageRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ImageRegistry *ImageRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ImageRegistry.Contract.contract.Transact(opts, method, params...)
}

// IsImageTrusted is a free data retrieval call binding the contract method 0x43f9405f.
//
// Solidity: function isImageTrusted(bytes32 imageHash) view returns(bool)
func (_ImageRegistry *ImageRegistryCaller) IsImageTrusted(opts *bind.CallOpts, imageHash [32]byte) (bool, error) {
	var out []interface{}
	err := _ImageRegistry.contract.Call(opts, &out, "isImageTrusted", imageHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsImageTrusted is a free data retrieval call binding the contract method 0x43f9405f.
//
// Solidity: function isImageTrusted(bytes32 imageHash) view returns(bool)
func (_ImageRegistry *ImageRegistrySession) IsImageTrusted(imageHash [32]byte) (bool, error) {
	return _ImageRegistry.Contract.IsImageTrusted(&_ImageRegistry.CallOpts, imageHash)
}

// IsImageTrusted is a free data retrieval call binding the contract method 0x43f9405f.
//
// Solidity: function isImageTrusted(bytes32 imageHash) view returns(bool)
func (_ImageRegistry *ImageRegistryCallerSession) IsImageTrusted(imageHash [32]byte) (bool, error) {
	return _ImageRegistry.Contract.IsImageTrusted(&_ImageRegistry.CallOpts, imageHash)
}

// TrustedImages is a free data retrieval call binding the contract method 0x4ee9f2d1.
//
// Solidity: function trustedImages(bytes32 ) view returns(bool)
func (_ImageRegistry *ImageRegistryCaller) TrustedImages(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _ImageRegistry.contract.Call(opts, &out, "trustedImages", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// TrustedImages is a free data retrieval call binding the contract method 0x4ee9f2d1.
//
// Solidity: function trustedImages(bytes32 ) view returns(bool)
func (_ImageRegistry *ImageRegistrySession) TrustedImages(arg0 [32]byte) (bool, error) {
	return _ImageRegistry.Contract.TrustedImages(&_ImageRegistry.CallOpts, arg0)
}

// TrustedImages is a free data retrieval call binding the contract method 0x4ee9f2d1.
//
// Solidity: function trustedImages(bytes32 ) view returns(bool)
func (_ImageRegistry *ImageRegistryCallerSession) TrustedImages(arg0 [32]byte) (bool, error) {
	return _ImageRegistry.Contract.TrustedImages(&_ImageRegistry.CallOpts, arg0)
}

// RegisterImage is a paid mutator transaction binding the contract method 0x3f633adb.
//
// Solidity: function registerImage(bytes32 imageHash) returns()
func (_ImageRegistry *ImageRegistryTransactor) RegisterImage(opts *bind.TransactOpts, imageHash [32]byte) (*types.Transaction, error) {
	return _ImageRegistry.contract.Transact(opts, "registerImage", imageHash)
}

// RegisterImage is a paid mutator transaction binding the contract method 0x3f633adb.
//
// Solidity: function registerImage(bytes32 imageHash) returns()
func (_ImageRegistry *ImageRegistrySession) RegisterImage(imageHash [32]byte) (*types.Transaction, error) {
	return _ImageRegistry.Contract.RegisterImage(&_ImageRegistry.TransactOpts, imageHash)
}

// RegisterImage is a paid mutator transaction binding the contract method 0x3f633adb.
//
// Solidity: function registerImage(bytes32 imageHash) returns()
func (_ImageRegistry *ImageRegistryTransactorSession) RegisterImage(imageHash [32]byte) (*types.Transaction, error) {
	return _ImageRegistry.Contract.RegisterImage(&_ImageRegistry.TransactOpts, imageHash)
}

// ImageRegistryImageRegisteredIterator is returned from FilterImageRegistered and is used to iterate over the raw logs and unpacked data for ImageRegistered events raised by the ImageRegistry contract.
type ImageRegistryImageRegisteredIterator struct {
	Event *ImageRegistryImageRegistered // Event containing the contract specifics and raw log

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
func (it *ImageRegistryImageRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ImageRegistryImageRegistered)
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
		it.Event = new(ImageRegistryImageRegistered)
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
func (it *ImageRegistryImageRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ImageRegistryImageRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ImageRegistryImageRegistered represents a ImageRegistered event raised by the ImageRegistry contract.
type ImageRegistryImageRegistered struct {
	ImageHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterImageRegistered is a free log retrieval operation binding the contract event 0xbb218b24bde6ba405fba71b6e1e123ab350a965dc6ad8fc3e53b3cc007547858.
//
// Solidity: event ImageRegistered(bytes32 indexed imageHash)
func (_ImageRegistry *ImageRegistryFilterer) FilterImageRegistered(opts *bind.FilterOpts, imageHash [][32]byte) (*ImageRegistryImageRegisteredIterator, error) {

	var imageHashRule []interface{}
	for _, imageHashItem := range imageHash {
		imageHashRule = append(imageHashRule, imageHashItem)
	}

	logs, sub, err := _ImageRegistry.contract.FilterLogs(opts, "ImageRegistered", imageHashRule)
	if err != nil {
		return nil, err
	}
	return &ImageRegistryImageRegisteredIterator{contract: _ImageRegistry.contract, event: "ImageRegistered", logs: logs, sub: sub}, nil
}

// WatchImageRegistered is a free log subscription operation binding the contract event 0xbb218b24bde6ba405fba71b6e1e123ab350a965dc6ad8fc3e53b3cc007547858.
//
// Solidity: event ImageRegistered(bytes32 indexed imageHash)
func (_ImageRegistry *ImageRegistryFilterer) WatchImageRegistered(opts *bind.WatchOpts, sink chan<- *ImageRegistryImageRegistered, imageHash [][32]byte) (event.Subscription, error) {

	var imageHashRule []interface{}
	for _, imageHashItem := range imageHash {
		imageHashRule = append(imageHashRule, imageHashItem)
	}

	logs, sub, err := _ImageRegistry.contract.WatchLogs(opts, "ImageRegistered", imageHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ImageRegistryImageRegistered)
				if err := _ImageRegistry.contract.UnpackLog(event, "ImageRegistered", log); err != nil {
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

// ParseImageRegistered is a log parse operation binding the contract event 0xbb218b24bde6ba405fba71b6e1e123ab350a965dc6ad8fc3e53b3cc007547858.
//
// Solidity: event ImageRegistered(bytes32 indexed imageHash)
func (_ImageRegistry *ImageRegistryFilterer) ParseImageRegistered(log types.Log) (*ImageRegistryImageRegistered, error) {
	event := new(ImageRegistryImageRegistered)
	if err := _ImageRegistry.contract.UnpackLog(event, "ImageRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
