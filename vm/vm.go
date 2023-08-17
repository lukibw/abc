package vm

import (
	"fmt"
	"sync"

	"github.com/lukibw/abc/compiler"
)

type VM interface {
	Run() error
}

func New(c compiler.Compiler) (VM, error) {
	chunk, err := c.Run()
	if err != nil {
		return nil, err
	}
	return &vm{chunk, make([]value, 0), sync.Mutex{}}, nil
}

type vm struct {
	chunk *compiler.Chunk
	stack []value
	mutex sync.Mutex
}

func (vm *vm) push(v value) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()
	vm.stack = append(vm.stack, v)
}

func (vm *vm) pop() value {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()
	item := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return item
}

func (vm *vm) peek(distance int) value {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()
	return vm.stack[len(vm.stack)-1-distance]
}

func (vm *vm) binary(f func(x, y float64) float64) error {
	if !vm.peek(0).isNumber() || !vm.peek(1).isNumber() {
		return &Error{ErrNumberOperands}
	}
	b := vm.pop().asNumber()
	a := vm.pop().asNumber()
	vm.push(newNumber(f(a, b)))
	return nil
}

func (vm *vm) comparison(f func(x, y float64) bool) error {
	if !vm.peek(0).isNumber() || !vm.peek(1).isNumber() {
		return &Error{ErrNumberOperands}
	}
	b := vm.pop().asNumber()
	a := vm.pop().asNumber()
	vm.push(newBoolean(f(a, b)))
	return nil
}

func (vm *vm) Run() error {
	var err error
	i := 0
	for {
		o := compiler.Operation(vm.chunk.Code[i])
		fmt.Printf("%04d %s", i, o)
		i++
		switch o {
		case compiler.OperationReturn:
			fmt.Println()
			fmt.Println(vm.pop())
			return nil
		case compiler.OperationConstant:
			j := vm.chunk.Code[i]
			constant := vm.chunk.Constants[j]
			fmt.Printf(" %d %g", j, constant)
			i++
			vm.push(newNumber(constant))
		case compiler.OperationNegate:
			if !vm.peek(0).isNumber() {
				return &Error{ErrNumberOperand}
			}
			vm.push(newNumber(-vm.pop().asNumber()))
		case compiler.OperationAdd:
			if err = vm.binary(func(x, y float64) float64 { return x + y }); err != nil {
				return err
			}
		case compiler.OperationSubtract:
			if err = vm.binary(func(x, y float64) float64 { return x - y }); err != nil {
				return err
			}
		case compiler.OperationMultiply:
			if err = vm.binary(func(x, y float64) float64 { return x * y }); err != nil {
				return err
			}
		case compiler.OperationDivide:
			if err = vm.binary(func(x, y float64) float64 { return x / y }); err != nil {
				return err
			}
		case compiler.OperationNil:
			vm.push(newNil())
		case compiler.OperationFalse:
			vm.push(newBoolean(false))
		case compiler.OperationTrue:
			vm.push(newBoolean(true))
		case compiler.OperationNot:
			vm.push(newBoolean(vm.pop().isFalsey()))
		case compiler.OperationEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(newBoolean(a == b))
		case compiler.OperationGreater:
			if err = vm.comparison(func(x, y float64) bool { return x > y }); err != nil {
				return err
			}
		case compiler.OperationLess:
			if err = vm.comparison(func(x, y float64) bool { return x < y }); err != nil {
				return err
			}
		}
		fmt.Println()
	}
}
