package abc

import (
	"errors"
	"fmt"
)

func isFalsey(v value) bool {
	return v.isNil() || (v.isBoolean() && !v.asBoolean())
}

type VM struct {
	chunk *Chunk
	stack *stack
}

func NewVM(chunk *Chunk) *VM {
	return &VM{chunk, newStack()}
}

func (vm *VM) binary(f func(x, y float64) float64) error {
	if !vm.stack.peek(0).isNumber() || !vm.stack.peek(1).isNumber() {
		return errors.New("operands must be numbers")
	}
	b := vm.stack.pop().asNumber()
	a := vm.stack.pop().asNumber()
	vm.stack.push(newNumber(f(a, b)))
	return nil
}

func (vm *VM) comparison(f func(x, y float64) bool) error {
	if !vm.stack.peek(0).isNumber() || !vm.stack.peek(1).isNumber() {
		return errors.New("operands must be numbers")
	}
	b := vm.stack.pop().asNumber()
	a := vm.stack.pop().asNumber()
	vm.stack.push(newBoolean(f(a, b)))
	return nil
}

func (vm *VM) Interpret() error {
	i := 0
	for {
		instruction := vm.chunk.code[i]
		i++
		switch instruction {
		case byte(OpReturn):
			fmt.Println(vm.stack.pop())
			return nil
		case byte(OpConstant):
			constant := vm.chunk.constants[vm.chunk.code[i]]
			i++
			vm.stack.push(newNumber(constant))
		case byte(OpNegate):
			if !vm.stack.peek(0).isNumber() {
				return errors.New("operand must be a number")
			}
			vm.stack.push(newNumber(-vm.stack.pop().asNumber()))
		case byte(OpAdd):
			if err := vm.binary(func(x, y float64) float64 { return x + y }); err != nil {
				return err
			}
		case byte(OpSubtract):
			if err := vm.binary(func(x, y float64) float64 { return x - y }); err != nil {
				return err
			}
		case byte(OpMultiply):
			if err := vm.binary(func(x, y float64) float64 { return x * y }); err != nil {
				return err
			}
		case byte(OpDivide):
			if err := vm.binary(func(x, y float64) float64 { return x / y }); err != nil {
				return err
			}
		case byte(OpNil):
			vm.stack.push(newNil())
		case byte(OpFalse):
			vm.stack.push(newBoolean(false))
		case byte(OpTrue):
			vm.stack.push(newBoolean(true))
		case byte(OpNot):
			vm.stack.push(newBoolean(isFalsey(vm.stack.pop())))
		case byte(OpEqual):
			b := vm.stack.pop()
			a := vm.stack.pop()
			vm.stack.push(newBoolean(a == b))
		case byte(OpGreater):
			if err := vm.comparison(func(x, y float64) bool { return x > y }); err != nil {
				return err
			}
		case byte(OpLess):
			if err := vm.comparison(func(x, y float64) bool { return x < y }); err != nil {
				return err
			}
		}
	}
}
