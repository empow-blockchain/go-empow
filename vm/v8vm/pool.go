package v8

import (
	"github.com/empow-blockchain/go-empow/core/contract"
	"github.com/empow-blockchain/go-empow/vm/host"
)

type vmPoolType int

const (
	// CompileVMPool maintains a pool of compile vm instance
	CompileVMPool vmPoolType = iota
	// RunVMPool maintains a pool of run vm instance
	RunVMPool
)

// VMPool manage all V8VM instance.
type VMPool struct {
	compilePoolSize int
	runPoolSize     int
	compilePoolBuff chan *VM
	runPoolBuff     chan *VM
	jsPath          string
}

// NewVMPool create new VMPool instance.
func NewVMPool(compilePoolSize, runPoolSize int) *VMPool {
	return &VMPool{
		compilePoolSize: compilePoolSize,
		runPoolSize:     runPoolSize,
		compilePoolBuff: make(chan *VM, compilePoolSize),
		runPoolBuff:     make(chan *VM, runPoolSize),
	}
}

func (vmp *VMPool) getCompileVM() *VM {
	vm := <-vmp.compilePoolBuff
	vm.refCount++
	return vm
}

func (vmp *VMPool) getRunVM() *VM {
	vm := <-vmp.runPoolBuff
	vm.refCount++
	return vm
}

// Init init VMPool.
func (vmp *VMPool) Init() error {
	// Fill vmPoolBuffer
	for i := 0; i < vmp.compilePoolSize; i++ {
		var e = NewVMWithChannel(CompileVMPool, vmp.jsPath, vmp.compilePoolBuff)
		vmp.compilePoolBuff <- e
	}
	for i := 0; i < vmp.runPoolSize; i++ {
		var e = NewVMWithChannel(RunVMPool, vmp.jsPath, vmp.runPoolBuff)
		vmp.runPoolBuff <- e
	}
	return nil
}

// SetJSPath set standard Javascript library path.
func (vmp *VMPool) SetJSPath(path string) {
	vmp.jsPath = path
}

// Validate js code and abi.
func (vmp *VMPool) Validate(contract *contract.Contract) error {
	vm := vmp.getCompileVM()
	defer func() {
		go vm.recycle(CompileVMPool)
	}()

	return vm.validate(contract)
}

// Compile compile js code to binary.
func (vmp *VMPool) Compile(contract *contract.Contract) (string, error) {
	vm := vmp.getCompileVM()
	defer func() {
		go vm.recycle(CompileVMPool)
	}()

	return vm.compile(contract)
}

// LoadAndCall load compiled Javascript code and run code with specified api and args
func (vmp *VMPool) LoadAndCall(host *host.Host, contract *contract.Contract, api string, args ...interface{}) (rtn []interface{}, cost contract.Cost, err error) {
	vm := vmp.getRunVM()
	defer func() {
		go vm.recycle(RunVMPool)
	}()

	vm.setHost(host)
	preparedCode, _ := vm.setContract(contract, api, args)

	return vm.execute(preparedCode)
}

// Release release all V8VM instance in VMPool
func (vmp *VMPool) Release() {
	close(vmp.compilePoolBuff)
	for e := range vmp.compilePoolBuff {
		e.release()
	}

	close(vmp.runPoolBuff)
	for e := range vmp.runPoolBuff {
		e.release()
	}
}
