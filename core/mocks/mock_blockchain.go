// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/empow-blockchain/go-empow/core/block (interfaces: Chain)

// Package core_mock is a generated GoMock package.
package core_mock

import (
	gomock "github.com/golang/mock/gomock"
	block "github.com/empow-blockchain/go-empow/core/block"
	tx "github.com/empow-blockchain/go-empow/core/tx"
	reflect "reflect"
)

// MockChain is a mock of Chain interface
type MockChain struct {
	ctrl     *gomock.Controller
	recorder *MockChainMockRecorder
}

// MockChainMockRecorder is the mock recorder for MockChain
type MockChainMockRecorder struct {
	mock *MockChain
}

// NewMockChain creates a new mock instance
func NewMockChain(ctrl *gomock.Controller) *MockChain {
	mock := &MockChain{ctrl: ctrl}
	mock.recorder = &MockChainMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockChain) EXPECT() *MockChainMockRecorder {
	return m.recorder
}

// AllDelaytx mocks base method
func (m *MockChain) AllDelaytx() ([]*tx.Tx, error) {
	ret := m.ctrl.Call(m, "AllDelaytx")
	ret0, _ := ret[0].([]*tx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllDelaytx indicates an expected call of AllDelaytx
func (mr *MockChainMockRecorder) AllDelaytx() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllDelaytx", reflect.TypeOf((*MockChain)(nil).AllDelaytx))
}

// CheckLength mocks base method
func (m *MockChain) CheckLength() {
	m.ctrl.Call(m, "CheckLength")
}

// CheckLength indicates an expected call of CheckLength
func (mr *MockChainMockRecorder) CheckLength() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLength", reflect.TypeOf((*MockChain)(nil).CheckLength))
}

// Close mocks base method
func (m *MockChain) Close() {
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockChainMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockChain)(nil).Close))
}

// Draw mocks base method
func (m *MockChain) Draw(arg0, arg1 int64) string {
	ret := m.ctrl.Call(m, "Draw", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// Draw indicates an expected call of Draw
func (mr *MockChainMockRecorder) Draw(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Draw", reflect.TypeOf((*MockChain)(nil).Draw), arg0, arg1)
}

// GetBlockByHash mocks base method
func (m *MockChain) GetBlockByHash(arg0 []byte) (*block.Block, error) {
	ret := m.ctrl.Call(m, "GetBlockByHash", arg0)
	ret0, _ := ret[0].(*block.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockByHash indicates an expected call of GetBlockByHash
func (mr *MockChainMockRecorder) GetBlockByHash(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByHash", reflect.TypeOf((*MockChain)(nil).GetBlockByHash), arg0)
}

// GetBlockByNumber mocks base method
func (m *MockChain) GetBlockByNumber(arg0 int64) (*block.Block, error) {
	ret := m.ctrl.Call(m, "GetBlockByNumber", arg0)
	ret0, _ := ret[0].(*block.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockByNumber indicates an expected call of GetBlockByNumber
func (mr *MockChainMockRecorder) GetBlockByNumber(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByNumber", reflect.TypeOf((*MockChain)(nil).GetBlockByNumber), arg0)
}

// GetBlockNumberByTxHash mocks base method
func (m *MockChain) GetBlockNumberByTxHash(arg0 []byte) (int64, error) {
	ret := m.ctrl.Call(m, "GetBlockNumberByTxHash", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockNumberByTxHash indicates an expected call of GetBlockNumberByTxHash
func (mr *MockChainMockRecorder) GetBlockNumberByTxHash(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockNumberByTxHash", reflect.TypeOf((*MockChain)(nil).GetBlockNumberByTxHash), arg0)
}

// GetHashByNumber mocks base method
func (m *MockChain) GetHashByNumber(arg0 int64) ([]byte, error) {
	ret := m.ctrl.Call(m, "GetHashByNumber", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHashByNumber indicates an expected call of GetHashByNumber
func (mr *MockChainMockRecorder) GetHashByNumber(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHashByNumber", reflect.TypeOf((*MockChain)(nil).GetHashByNumber), arg0)
}

// GetReceipt mocks base method
func (m *MockChain) GetReceipt(arg0 []byte) (*tx.TxReceipt, error) {
	ret := m.ctrl.Call(m, "GetReceipt", arg0)
	ret0, _ := ret[0].(*tx.TxReceipt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReceipt indicates an expected call of GetReceipt
func (mr *MockChainMockRecorder) GetReceipt(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReceipt", reflect.TypeOf((*MockChain)(nil).GetReceipt), arg0)
}

// GetReceiptByTxHash mocks base method
func (m *MockChain) GetReceiptByTxHash(arg0 []byte) (*tx.TxReceipt, error) {
	ret := m.ctrl.Call(m, "GetReceiptByTxHash", arg0)
	ret0, _ := ret[0].(*tx.TxReceipt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReceiptByTxHash indicates an expected call of GetReceiptByTxHash
func (mr *MockChainMockRecorder) GetReceiptByTxHash(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReceiptByTxHash", reflect.TypeOf((*MockChain)(nil).GetReceiptByTxHash), arg0)
}

// GetTx mocks base method
func (m *MockChain) GetTx(arg0 []byte) (*tx.Tx, error) {
	ret := m.ctrl.Call(m, "GetTx", arg0)
	ret0, _ := ret[0].(*tx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTx indicates an expected call of GetTx
func (mr *MockChainMockRecorder) GetTx(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTx", reflect.TypeOf((*MockChain)(nil).GetTx), arg0)
}

// HasReceipt mocks base method
func (m *MockChain) HasReceipt(arg0 []byte) (bool, error) {
	ret := m.ctrl.Call(m, "HasReceipt", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasReceipt indicates an expected call of HasReceipt
func (mr *MockChainMockRecorder) HasReceipt(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasReceipt", reflect.TypeOf((*MockChain)(nil).HasReceipt), arg0)
}

// HasTx mocks base method
func (m *MockChain) HasTx(arg0 []byte) (bool, error) {
	ret := m.ctrl.Call(m, "HasTx", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasTx indicates an expected call of HasTx
func (mr *MockChainMockRecorder) HasTx(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasTx", reflect.TypeOf((*MockChain)(nil).HasTx), arg0)
}

// Length mocks base method
func (m *MockChain) Length() int64 {
	ret := m.ctrl.Call(m, "Length")
	ret0, _ := ret[0].(int64)
	return ret0
}

// Length indicates an expected call of Length
func (mr *MockChainMockRecorder) Length() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Length", reflect.TypeOf((*MockChain)(nil).Length))
}

// Push mocks base method
func (m *MockChain) Push(arg0 *block.Block) error {
	ret := m.ctrl.Call(m, "Push", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Push indicates an expected call of Push
func (mr *MockChainMockRecorder) Push(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockChain)(nil).Push), arg0)
}

// SetLength mocks base method
func (m *MockChain) SetLength(arg0 int64) {
	m.ctrl.Call(m, "SetLength", arg0)
}

// SetLength indicates an expected call of SetLength
func (mr *MockChainMockRecorder) SetLength(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLength", reflect.TypeOf((*MockChain)(nil).SetLength), arg0)
}

// Size mocks base method
func (m *MockChain) Size() (int64, error) {
	ret := m.ctrl.Call(m, "Size")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Size indicates an expected call of Size
func (mr *MockChainMockRecorder) Size() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockChain)(nil).Size))
}

// Top mocks base method
func (m *MockChain) Top() (*block.Block, error) {
	ret := m.ctrl.Call(m, "Top")
	ret0, _ := ret[0].(*block.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Top indicates an expected call of Top
func (mr *MockChainMockRecorder) Top() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Top", reflect.TypeOf((*MockChain)(nil).Top))
}

// TxTotal mocks base method
func (m *MockChain) TxTotal() int64 {
	ret := m.ctrl.Call(m, "TxTotal")
	ret0, _ := ret[0].(int64)
	return ret0
}

// TxTotal indicates an expected call of TxTotal
func (mr *MockChainMockRecorder) TxTotal() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxTotal", reflect.TypeOf((*MockChain)(nil).TxTotal))
}
