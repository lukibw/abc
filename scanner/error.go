package scanner

import "fmt"

type ErrorKind int

const (
	ErrUnexpectedCharacter ErrorKind = iota
	ErrUnterminatedString
)

var errorMessages = map[ErrorKind]string{
	ErrUnexpectedCharacter: "unexpected character",
	ErrUnterminatedString:  "unterminated string",
}

func (k ErrorKind) String() string {
	return errorMessages[k]
}

type Error struct {
	Kind ErrorKind
	Line int
}

func (e *Error) Error() string {
	return fmt.Sprintf("[line %d] compilation error: %s", e.Line, e.Kind)
}
