package compiler

import (
	"fmt"
	"strings"

	"github.com/lukibw/abc/scanner"
)

type ErrorKind int

const (
	ErrTooManyConstants ErrorKind = iota
	ErrMissingExpr
	ErrMissingExprEnd
	ErrMissingExprRightParen
)

var errorMessages = map[ErrorKind]string{
	ErrTooManyConstants:      "too many constants in one chunk",
	ErrMissingExpr:           "missing expression",
	ErrMissingExprEnd:        "missing end of expression",
	ErrMissingExprRightParen: "missing ')' after expression",
}

func (k ErrorKind) String() string {
	return errorMessages[k]
}

type Error struct {
	Kind  ErrorKind
	Token *scanner.Token
}

func (e *Error) Error() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("[line %d] compilation error", e.Token.Line))
	if e.Token.Kind == scanner.TokenEof {
		sb.WriteString(" at end")
	}
	sb.WriteString(fmt.Sprintf(": %s", e.Kind))
	return sb.String()
}
