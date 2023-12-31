package compiler

import (
	"fmt"
	"strings"

	"github.com/lukibw/abc/scanner"
)

type ErrorKind int

const (
	ErrTooManyConstants ErrorKind = iota
	ErrTooManyLocals
	ErrTooBigJump
	ErrTooBigLoop
	ErrVarAlreadyDefined
	ErrVarOwnInitializer
	ErrIfLeftParen
	ErrIfRightParen
	ErrWhileLeftParen
	ErrWhileRightParen
	ErrForLeftParen
	ErrForRightParen
	ErrForConditionSemicolon
	ErrInvalidAssignTarget
	ErrMissingVarName
	ErrMissingVarSemicolon
	ErrMissingValueSemicolon
	ErrMissingExpr
	ErrMissingExprEnd
	ErrMissingExprRightParen
	ErrMissingExprSemicolon
	ErrMissingBlockRightBrace
)

var errorMessages = map[ErrorKind]string{
	ErrTooManyConstants:       "too many constants in one chunk",
	ErrTooManyLocals:          "too many local variables in function",
	ErrTooBigJump:             "too much code to jump over",
	ErrTooBigLoop:             "loop body too large",
	ErrVarAlreadyDefined:      "already a variable with this name in this scope",
	ErrVarOwnInitializer:      "cannot read local variable in its own intializer",
	ErrIfLeftParen:            "missing '(' after 'if'",
	ErrIfRightParen:           "missing ')' after condition",
	ErrWhileLeftParen:         "missing '(' after 'while'",
	ErrWhileRightParen:        "missing ')' after condition",
	ErrForLeftParen:           "missing '(' after 'for'",
	ErrForRightParen:          "missing ')' after for clauses",
	ErrForConditionSemicolon:  "missing ';' after loop condition",
	ErrInvalidAssignTarget:    "invalid assignment target",
	ErrMissingVarName:         "missing variable name",
	ErrMissingVarSemicolon:    "missing ';' after variable declaration",
	ErrMissingValueSemicolon:  "missing ';' after value",
	ErrMissingExpr:            "missing expression",
	ErrMissingExprEnd:         "missing end of expression",
	ErrMissingExprRightParen:  "missing ')' after expression",
	ErrMissingExprSemicolon:   "missing ';' after expression",
	ErrMissingBlockRightBrace: "missing '}' after block",
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
