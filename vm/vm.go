package vm

import (
	"fmt"
	"strings"
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
	return &vm{chunk, make([]compiler.Value, 0), sync.Mutex{}, make(map[string]compiler.Value)}, nil
}

type vm struct {
	chunk   *compiler.Chunk
	stack   []compiler.Value
	mutex   sync.Mutex
	globals map[string]compiler.Value
}

func (vm *vm) push(v compiler.Value) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()
	vm.stack = append(vm.stack, v)
}

func (vm *vm) pop() compiler.Value {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()
	item := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return item
}

func (vm *vm) peek(distance int) compiler.Value {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()
	return vm.stack[len(vm.stack)-1-distance]
}

func (vm *vm) binary(f func(x, y float64) float64) error {
	if !vm.peek(0).IsNumber() || !vm.peek(1).IsNumber() {
		return &Error{ErrNumberOperands}
	}
	b := vm.pop().AsNumber()
	a := vm.pop().AsNumber()
	vm.push(compiler.NewNumber(f(a, b)))
	return nil
}

func (vm *vm) comparison(f func(x, y float64) bool) error {
	if !vm.peek(0).IsNumber() || !vm.peek(1).IsNumber() {
		return &Error{ErrNumberOperands}
	}
	b := vm.pop().AsNumber()
	a := vm.pop().AsNumber()
	vm.push(compiler.NewBoolean(f(a, b)))
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
		case compiler.OperationSetGlobal:
			j := vm.chunk.Code[i]
			constant := vm.chunk.Constants[j]
			fmt.Printf(" %d %s", j, constant)
			i++
			_, ok := vm.globals[constant.AsString()]
			if !ok {
				return &Error{ErrUndefinedVar}
			}
			vm.globals[constant.AsString()] = vm.peek(0)
		case compiler.OperationGetGlobal:
			j := vm.chunk.Code[i]
			constant := vm.chunk.Constants[j]
			fmt.Printf(" %d %s", j, constant)
			i++
			value, ok := vm.globals[constant.AsString()]
			if !ok {
				return &Error{ErrUndefinedVar}
			}
			vm.push(value)
		case compiler.OperationDefineGlobal:
			j := vm.chunk.Code[i]
			constant := vm.chunk.Constants[j]
			fmt.Printf(" %d %s", j, constant)
			i++
			vm.globals[constant.AsString()] = vm.pop()
		case compiler.OperationPrint:
			fmt.Println(vm.pop())
		case compiler.OperationPop:
			vm.pop()
		case compiler.OperationReturn:
			fmt.Println()
			return nil
		case compiler.OperationConstant:
			j := vm.chunk.Code[i]
			constant := vm.chunk.Constants[j]
			fmt.Printf(" %d %s", j, constant)
			i++
			vm.push(constant)
		case compiler.OperationNegate:
			if !vm.peek(0).IsNumber() {
				return &Error{ErrNumberOperand}
			}
			vm.push(compiler.NewNumber(-vm.pop().AsNumber()))
		case compiler.OperationAdd:
			b := vm.peek(0)
			a := vm.peek(1)
			areStrings := a.IsString() && b.IsString()
			areNumbers := a.IsNumber() && b.IsNumber()
			if !areStrings && !areNumbers {
				return &Error{ErrNumberOrStringOperands}
			}
			b = vm.pop()
			a = vm.pop()
			if areStrings {
				var sb strings.Builder
				sb.WriteString(a.AsString())
				sb.WriteString(b.AsString())
				vm.push(compiler.NewString(sb.String()))
			} else {
				vm.push(compiler.NewNumber(a.AsNumber() + b.AsNumber()))
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
			vm.push(compiler.NewNil())
		case compiler.OperationFalse:
			vm.push(compiler.NewBoolean(false))
		case compiler.OperationTrue:
			vm.push(compiler.NewBoolean(true))
		case compiler.OperationNot:
			vm.push(compiler.NewBoolean(vm.pop().IsFalsey()))
		case compiler.OperationEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(compiler.NewBoolean(a == b))
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
