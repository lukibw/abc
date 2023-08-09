package abc

import (
	"fmt"
	"strings"
)

type OpCode byte

const (
	OpReturn OpCode = iota
	OpConstant
)

func (c OpCode) String() string {
	switch c {
	case OpReturn:
		return "OP_RETURN"
	case OpConstant:
		return "OP_CONSTANT"
	default:
		return ""
	}
}

func (c OpCode) Offset() int {
	switch c {
	case OpReturn:
		return 1
	case OpConstant:
		return 2
	default:
		return 0
	}
}

type Chunk struct {
	code      []byte
	lines     []int
	constants []float64
}

func NewChunk() *Chunk {
	return &Chunk{make([]byte, 0), make([]int, 0), make([]float64, 0)}
}

func (c *Chunk) Write(value byte, line int) {
	c.code = append(c.code, value)
	c.lines = append(c.lines, line)
}

func (c *Chunk) AddConstant(value float64) byte {
	c.constants = append(c.constants, value)
	return byte(len(c.constants) - 1)
}

func (c *Chunk) String() string {
	i := 0
	sb := strings.Builder{}
	for i < len(c.code) {
		sb.WriteString(fmt.Sprintf("%04d ", i))
		if i > 0 && c.lines[i] == c.lines[i-1] {
			sb.WriteString("   | ")
		} else {
			sb.WriteString(fmt.Sprintf("%4d ", c.lines[i]))
		}
		instruction := OpCode(c.code[i])
		switch instruction {
		case OpReturn:
			sb.WriteString(fmt.Sprintf("%s\n", instruction))
		case OpConstant:
			constant := c.code[i+1]
			sb.WriteString(fmt.Sprintf("%-16s %4d '%g'\n", instruction, constant, c.constants[constant]))
		}
		i += instruction.Offset()
	}
	return sb.String()
}
