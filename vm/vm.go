package vm

import (
	"github.com/empow-blockchain/go-empow/core/contract"
	"github.com/empow-blockchain/go-empow/vm/host"
)

//go:generate mockgen -destination vm_mock.go -package vm github.com/empow-blockchain/go-empow/vm VM

// VM ...
type VM interface {
	Init() error
	Validate(contract *contract.Contract) error
	Compile(contract *contract.Contract) (string, error)
	LoadAndCall(host *host.Host, contract *contract.Contract, api string, args ...interface{}) (rtn []interface{}, cost contract.Cost, err error)
	Release()
}
