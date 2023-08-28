package vm

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/lukibw/abc/compiler"
)

type VM interface {
	Run() error
}

func New(c compiler.Compiler, logger *log.Logger) (VM, error) {
	chunk, err := c.Run()
	if err != nil {
		return nil, err
	}
	return &vm{false, 0, logger, chunk, make([]compiler.Value, 0), sync.Mutex{}, make(map[string]compiler.Value)}, nil
}

type vm struct {
	isEnd   bool
	i       int
	logger  *log.Logger
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

func (vm *vm) readOperation() compiler.Operation {
	return compiler.Operation(vm.chunk.Code[vm.i])
}

func (vm *vm) readSlot() byte {
	return vm.chunk.Code[vm.i+1]
}

func (vm *vm) readConstant() compiler.Value {
	return vm.chunk.Constants[vm.readSlot()]
}

func (vm *vm) readJump() int {
	return int(uint16(vm.chunk.Code[vm.i+1])<<8 | uint16(vm.chunk.Code[vm.i+2]))
}

func (vm *vm) debug() {
	var sb strings.Builder
	o := vm.readOperation()
	sb.WriteString(fmt.Sprintf("%04d | %-16s |", vm.i, o))
	switch o {
	case compiler.OperationJump, compiler.OperationJumpIfFalse, compiler.OperationLoop:
		sb.WriteString(fmt.Sprintf(" %d", vm.readJump()))
	case compiler.OperationGetLocal, compiler.OperationSetLocal:
		sb.WriteString(fmt.Sprintf(" %d", vm.readSlot()))
	case compiler.OperationConstant, compiler.OperationDefineGlobal, compiler.OperationGetGlobal, compiler.OperationSetGlobal:
		sb.WriteString(fmt.Sprintf(" %s", vm.readConstant()))
	}
	sb.WriteRune('\n')
	vm.logger.Print(sb.String())
}

func (vm *vm) execute() error {
	o := vm.readOperation()
	switch o {
	case compiler.OperationJump:
		vm.i += vm.readJump()
	case compiler.OperationJumpIfFalse:
		if vm.peek(0).IsFalsey() {
			vm.i += vm.readJump()
		}
	case compiler.OperationLoop:
		vm.i -= vm.readJump()
	case compiler.OperationGetLocal:
		vm.push(vm.stack[vm.readSlot()])
	case compiler.OperationSetLocal:
		vm.stack[vm.readSlot()] = vm.peek(0)
	case compiler.OperationSetGlobal:
		constant := vm.readConstant()
		_, ok := vm.globals[constant.AsString()]
		if !ok {
			return &Error{ErrUndefinedVar}
		}
		vm.globals[constant.AsString()] = vm.peek(0)
	case compiler.OperationGetGlobal:
		value, ok := vm.globals[vm.readConstant().AsString()]
		if !ok {
			return &Error{ErrUndefinedVar}
		}
		vm.push(value)
	case compiler.OperationDefineGlobal:
		vm.globals[vm.readConstant().AsString()] = vm.pop()
	case compiler.OperationConstant:
		vm.push(vm.readConstant())
	case compiler.OperationPrint:
		fmt.Println(vm.pop())
	case compiler.OperationPop:
		vm.pop()
	case compiler.OperationReturn:
		vm.isEnd = true
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
		if err := vm.binary(func(x, y float64) float64 { return x - y }); err != nil {
			return err
		}
	case compiler.OperationMultiply:
		if err := vm.binary(func(x, y float64) float64 { return x * y }); err != nil {
			return err
		}
	case compiler.OperationDivide:
		if err := vm.binary(func(x, y float64) float64 { return x / y }); err != nil {
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
		if err := vm.comparison(func(x, y float64) bool { return x > y }); err != nil {
			return err
		}
	case compiler.OperationLess:
		if err := vm.comparison(func(x, y float64) bool { return x < y }); err != nil {
			return err
		}
	}
	switch o {
	case compiler.OperationJump, compiler.OperationJumpIfFalse, compiler.OperationLoop:
		vm.i += 3
	case compiler.OperationGetLocal, compiler.OperationSetLocal, compiler.OperationConstant, compiler.OperationDefineGlobal, compiler.OperationGetGlobal, compiler.OperationSetGlobal:
		vm.i += 2
	default:
		vm.i++
	}
	return nil
}

func (vm *vm) Run() error {
	var err error
	for !vm.isEnd {
		vm.debug()
		if err = vm.execute(); err != nil {
			return err
		}
	}
	return nil
}
