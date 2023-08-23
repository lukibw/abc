package compiler

import "math"

type Chunk struct {
	Code      []byte
	Constants []Value
}

func (c *Chunk) write(b byte) {
	c.Code = append(c.Code, b)
}

func (c *Chunk) writeOperation(o Operation) {
	c.write(byte(o))
}

func (c *Chunk) writeConstant(v Value) (uint8, bool) {
	if len(c.Constants) >= math.MaxUint8 {
		return 0, false
	}
	c.Constants = append(c.Constants, v)
	return uint8(len(c.Constants) - 1), true
}
