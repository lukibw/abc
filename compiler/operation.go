package compiler

type Operation int

const (
	OperationReturn Operation = iota
	OperationConstant
	OperationNegate
	OperationPrint
	OperationPop
	OperationDefineGlobal
	OperationGetGlobal
	OperationSetGlobal
	OperationGetLocal
	OperationSetLocal
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
	OperationReturn:       "RETURN",
	OperationConstant:     "CONSTANT",
	OperationNegate:       "NEGATE",
	OperationNot:          "NOT",
	OperationAdd:          "ADD",
	OperationSubtract:     "SUBTRACT",
	OperationMultiply:     "MULTIPLY",
	OperationDivide:       "DIVIDE",
	OperationNil:          "NIL",
	OperationFalse:        "FALSE",
	OperationTrue:         "TRUE",
	OperationEqual:        "EQUAL",
	OperationGreater:      "GREATER",
	OperationLess:         "LESS",
	OperationPrint:        "PRINT",
	OperationPop:          "POP",
	OperationDefineGlobal: "DEFINE_GLOBAL",
	OperationGetGlobal:    "GET_GLOBAL",
	OperationSetGlobal:    "SET_GLOBAL",
	OperationGetLocal:     "GET_LOCAL",
	OperationSetLocal:     "SET_LOCAL",
}

func (o Operation) String() string {
	return operations[o]
}
