package abc

import (
	"fmt"
	"strings"
)

type OpCode byte

const (
	OpReturn OpCode = iota
	OpConstant
	OpNegate
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpNil
	OpFalse
	OpTrue
	OpNot
	OpEqual
	OpGreater
	OpLess
)

func (c OpCode) String() string {
	switch c {
	case OpReturn:
		return "return"
	case OpConstant:
		return "constant"
	case OpNegate:
		return "negate"
	case OpAdd:
		return "add"
	case OpSubtract:
		return "subtract"
	case OpMultiply:
		return "multiply"
	case OpDivide:
		return "divide"
	case OpNil:
		return "nil"
	case OpFalse:
		return "false"
	case OpTrue:
		return "true"
	case OpNot:
		return "not"
	case OpEqual:
		return "equal"
	case OpGreater:
		return "greater"
	case OpLess:
		return "less"
	default:
		return ""
	}
}

func (c OpCode) Offset() int {
	switch c {
	case OpConstant:
		return 2
	default:
		return 1
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

func (c *Chunk) AddConstant(f float64) byte {
	c.constants = append(c.constants, f)
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
		case OpConstant:
			constant := c.code[i+1]
			sb.WriteString(fmt.Sprintf("%-16s %4d '%g'\n", instruction, constant, c.constants[constant]))
		default:
			sb.WriteString(fmt.Sprintf("%s\n", instruction))
		}
		i += instruction.Offset()
	}
	return sb.String()
}
