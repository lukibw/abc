package compiler

type Operation int

const (
	OperationReturn Operation = iota
	OperationConstant
	OperationNegate
	OperationNot
	OperationAdd
	OperationSubtract
	OperationMultiply
	OperationDivide
	OperationNil
	OperationFalse
	OperationTrue
	OperationEqual
	OperationGreater
	OperationLess
)

var operations = map[Operation]string{
	OperationReturn:   "RETURN",
	OperationConstant: "CONSTANT",
	OperationNegate:   "NEGATE",
	OperationNot:      "NOT",
	OperationAdd:      "ADD",
	OperationSubtract: "SUBTRACT",
	OperationMultiply: "MULTIPLY",
	OperationDivide:   "DIVIDE",
	OperationNil:      "NIL",
	OperationFalse:    "FALSE",
	OperationTrue:     "TRUE",
	OperationEqual:    "EQUAL",
	OperationGreater:  "GREATER",
	OperationLess:     "LESS",
}

func (o Operation) String() string {
	return operations[o]
}
