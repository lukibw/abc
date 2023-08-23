package compiler

import (
	"fmt"
	"strconv"

	"github.com/lukibw/abc/scanner"
)

type Compiler interface {
	Run() (*Chunk, error)
}

func New(s scanner.Scanner) Compiler {
	return &compiler{s, nil, nil, nil}
}

type compiler struct {
	scanner  scanner.Scanner
	previous *scanner.Token
	current  *scanner.Token
	chunk    *Chunk
}

func (c *compiler) makeConstant(v Value) (uint8, error) {
	i, ok := c.chunk.writeConstant(v)
	if !ok {
		return 0, &Error{ErrTooManyConstants, c.previous}
	}
	return i, nil
}

func (c *compiler) emitConstant(v Value) error {
	i, err := c.makeConstant(v)
	if err != nil {
		return err
	}
	c.chunk.writeOperation(OperationConstant)
	c.chunk.write(i)
	return nil
}

func (c *compiler) emitOperation(o Operation) {
	c.chunk.writeOperation(o)
}

func (c *compiler) emitOperations(o1, o2 Operation) {
	c.emitOperation(o1)
	c.emitOperation(o2)
}

func (c *compiler) check(k scanner.TokenKind) bool {
	return c.current.Kind == k
}

func (c *compiler) advance() error {
	c.previous = c.current
	t, err := c.scanner.Token()
	if err != nil {
		return err
	}
	c.current = t
	return nil
}

func (c *compiler) consume(t scanner.TokenKind, e ErrorKind) error {
	if c.current.Kind == t {
		return c.advance()
	}
	return &Error{e, c.current}
}

func (c *compiler) number() error {
	n, err := strconv.ParseFloat(c.previous.Lexeme, 64)
	if err != nil {
		panic("compiler: cannot parse float from token lexeme")
	}
	return c.emitConstant(NewNumber(n))
}

func (c *compiler) string() error {
	return c.emitConstant(NewString(c.previous.Lexeme[1 : len(c.previous.Lexeme)-1]))
}

func (c *compiler) binary() error {
	operator := c.previous.Kind
	if err := c.parsePrecedence(parseRules[operator].precedence + 1); err != nil {
		return err
	}
	switch operator {
	case scanner.TokenPlus:
		c.emitOperation(OperationAdd)
	case scanner.TokenMinus:
		c.emitOperation(OperationSubtract)
	case scanner.TokenStar:
		c.emitOperation(OperationMultiply)
	case scanner.TokenSlash:
		c.emitOperation(OperationDivide)
	case scanner.TokenBangEqual:
		c.emitOperations(OperationEqual, OperationNot)
	case scanner.TokenEqualEqual:
		c.emitOperation(OperationEqual)
	case scanner.TokenGreater:
		c.emitOperation(OperationGreater)
	case scanner.TokenGreaterEqual:
		c.emitOperations(OperationLess, OperationNot)
	case scanner.TokenLess:
		c.emitOperation(OperationLess)
	case scanner.TokenLessEqual:
		c.emitOperations(OperationGreater, OperationNot)
	default:
		panic(fmt.Sprintf("compiler: unexpected token kind '%s' for binary expression", operator))
	}
	return nil
}

func (c *compiler) literal() {
	value := c.previous.Kind
	switch value {
	case scanner.TokenNil:
		c.emitOperation(OperationNil)
	case scanner.TokenFalse:
		c.emitOperation(OperationFalse)
	case scanner.TokenTrue:
		c.emitOperation(OperationTrue)
	default:
		panic(fmt.Sprintf("compiler: unexpected token kind '%s' for literal expression", value))
	}
}

func (c *compiler) grouping() error {
	if err := c.expression(); err != nil {
		return err
	}
	return c.consume(scanner.TokenRightParen, ErrMissingExprRightParen)
}

func (c *compiler) unary() error {
	operator := c.previous.Kind
	if err := c.parsePrecedence(precedenceUnary); err != nil {
		return err
	}
	switch operator {
	case scanner.TokenMinus:
		c.emitOperation(OperationNegate)
	case scanner.TokenBang:
		c.emitOperation(OperationNot)
	default:
		panic(fmt.Sprintf("compiler: unexpected token kind '%s' for unary expression", operator))
	}
	return nil
}

func (c *compiler) namedVariable(t *scanner.Token, canAssign bool) error {
	i, err := c.identifierConstant(t)
	if err != nil {
		return err
	}
	if canAssign && c.check(scanner.TokenEqual) {
		if err = c.advance(); err != nil {
			return err
		}
		if err = c.expression(); err != nil {
			return err
		}
		c.emitOperation(OperationSetGlobal)
		c.chunk.write(i)
	} else {
		c.emitOperation(OperationGetGlobal)
		c.chunk.write(i)
	}
	return nil
}

func (c *compiler) variable(canAssign bool) error {
	return c.namedVariable(c.previous, canAssign)
}

func (c *compiler) parseFunction(f parseFunction, canAssign bool) error {
	switch f {
	case parseFunctionBinary:
		return c.binary()
	case parseFunctionGrouping:
		return c.grouping()
	case parseFunctionUnary:
		return c.unary()
	case parseFunctionNumber:
		return c.number()
	case parseFunctionString:
		return c.string()
	case parseFunctionLiteral:
		c.literal()
		return nil
	case parseFunctionVariable:
		return c.variable(canAssign)
	default:
		return &Error{ErrMissingExpr, c.previous}
	}
}

func (c *compiler) parsePrecedence(min precedence) error {
	var err error
	if err = c.advance(); err != nil {
		return err
	}
	canAssign := min <= precedenceAssignment
	if err = c.parseFunction(parseRules[c.previous.Kind].prefix, canAssign); err != nil {
		return err
	}
	for min <= parseRules[c.current.Kind].precedence {
		if err = c.advance(); err != nil {
			return err
		}
		if err = c.parseFunction(parseRules[c.previous.Kind].infix, canAssign); err != nil {
			return err
		}
	}
	if canAssign && c.check(scanner.TokenEqual) {
		if err = c.advance(); err != nil {
			return err
		}
		return &Error{ErrInvalidAssignTarget, c.previous}
	}
	return nil
}

func (c *compiler) expression() error {
	return c.parsePrecedence(precedenceAssignment)
}

func (c *compiler) printStatement() error {
	var err error
	if err = c.expression(); err != nil {
		return err
	}
	if err = c.consume(scanner.TokenSemicolon, ErrMissingValueSemicolon); err != nil {
		return err
	}
	c.emitOperation(OperationPrint)
	return nil
}

func (c *compiler) expressionStatement() error {
	var err error
	if err = c.expression(); err != nil {
		return err
	}
	if err = c.consume(scanner.TokenSemicolon, ErrMissingExprSemicolon); err != nil {
		return err
	}
	c.emitOperation(OperationPop)
	return nil
}

func (c *compiler) statement() error {
	if c.check(scanner.TokenPrint) {
		if err := c.advance(); err != nil {
			return err
		}
		return c.printStatement()
	}
	return c.expressionStatement()
}

func (c *compiler) identifierConstant(t *scanner.Token) (uint8, error) {
	return c.makeConstant(NewString(t.Lexeme))
}

func (c *compiler) parseVariable(k ErrorKind) (uint8, error) {
	if err := c.consume(scanner.TokenIdentifier, k); err != nil {
		return 0, err
	}
	return c.identifierConstant(c.previous)
}

func (c *compiler) defineVariable(v uint8) {
	c.emitOperation(OperationDefineGlobal)
	c.chunk.write(v)
}

func (c *compiler) varDeclaration() error {
	global, err := c.parseVariable(ErrMissingVarName)
	if err != nil {
		return err
	}
	if c.check(scanner.TokenEqual) {
		if err = c.advance(); err != nil {
			return err
		}
		if err = c.expression(); err != nil {
			return err
		}
	} else {
		c.emitOperation(OperationNil)
	}
	if err = c.consume(scanner.TokenSemicolon, ErrMissingVarSemicolon); err != nil {
		return err
	}
	c.defineVariable(global)
	return nil
}

func (c *compiler) declaration() error {
	if c.check(scanner.TokenVar) {
		if err := c.advance(); err != nil {
			return err
		}
		return c.varDeclaration()
	}
	return c.statement()
}

func (c *compiler) Run() (*Chunk, error) {
	if c.chunk == nil {
		var err error
		c.chunk = &Chunk{make([]byte, 0), make([]Value, 0)}
		if err = c.advance(); err != nil {
			return nil, err
		}
		for {
			if c.check(scanner.TokenEof) {
				if err = c.advance(); err != nil {
					return nil, err
				}
				break
			}
			if err = c.declaration(); err != nil {
				return nil, err
			}
		}
		if err = c.consume(scanner.TokenEof, ErrMissingExprEnd); err != nil {
			return nil, err
		}
		c.emitOperation(OperationReturn)
	}
	return c.chunk, nil
}
