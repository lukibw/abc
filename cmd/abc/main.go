package main

import (
	"fmt"
	"log"

	"github.com/lukibw/abc"
)

func main() {
	chunk := abc.NewChunk()
	constant := chunk.AddConstant(1.2)
	chunk.Write(byte(abc.OpConstant), 123)
	chunk.Write(constant, 123)
	constant = chunk.AddConstant(3.4)
	chunk.Write(byte(abc.OpConstant), 123)
	chunk.Write(constant, 123)
	chunk.Write(byte(abc.OpAdd), 123)
	constant = chunk.AddConstant(5.6)
	chunk.Write(byte(abc.OpConstant), 123)
	chunk.Write(constant, 123)
	chunk.Write(byte(abc.OpDivide), 123)
	chunk.Write(byte(abc.OpNegate), 123)
	chunk.Write(byte(abc.OpReturn), 123)
	fmt.Println("== test chunk ==")
	fmt.Print(chunk.String())
	vm := abc.NewVM(chunk)
	if err := vm.Interpret(); err != nil {
		log.Fatal(err)
	}
}
