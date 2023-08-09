package abc

import "fmt"

type VM struct {
	chunk *Chunk
	stack *stack
}

func NewVM(chunk *Chunk) *VM {
	return &VM{chunk, newStack()}
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
			vm.stack.push(constant)
		case byte(OpNegate):
			vm.stack.push(-vm.stack.pop())
		case byte(OpAdd):
			b := vm.stack.pop()
			a := vm.stack.pop()
			vm.stack.push(a + b)
		case byte(OpSubtract):
			b := vm.stack.pop()
			a := vm.stack.pop()
			vm.stack.push(a - b)
		case byte(OpMultiply):
			b := vm.stack.pop()
			a := vm.stack.pop()
			vm.stack.push(a * b)
		case byte(OpDivide):
			b := vm.stack.pop()
			a := vm.stack.pop()
			vm.stack.push(a / b)
		}
	}
}
