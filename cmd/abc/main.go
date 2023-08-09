package main

import (
	"fmt"

	"github.com/lukibw/abc"
)

func main() {
	c := abc.NewChunk()
	constant := c.AddConstant(1.2)
	c.Write(byte(abc.OpConstant), 123)
	c.Write(constant, 123)
	c.Write(byte(abc.OpReturn), 123)
	fmt.Println("== test chunk ==")
	fmt.Print(c.String())
}
