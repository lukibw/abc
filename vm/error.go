package vm

import "fmt"

type ErrorKind int

const (
	ErrNumberOperand ErrorKind = iota
	ErrNumberOperands
	ErrNumberOrStringOperands
)

var errorMessages = map[ErrorKind]string{
	ErrNumberOperand:          "operand must be a number",
	ErrNumberOperands:         "operands must be numbers",
	ErrNumberOrStringOperands: "operands must be two numbers or two strings",
}

func (k ErrorKind) String() string {
	return errorMessages[k]
}

type Error struct {
	Kind ErrorKind
}

func (e *Error) Error() string {
	return fmt.Sprintf("runtime error: %s", e.Kind)
}
