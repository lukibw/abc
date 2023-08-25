package compiler

import (
	"fmt"
	"math"
	"strconv"

	"github.com/lukibw/abc/scanner"
)

type Compiler interface {
	Run() (*Chunk, error)
}

func New(s scanner.Scanner) Compiler {
	return &compiler{s, nil, nil, make([]local, 0), 0, nil}
}

type local struct {
	name  *scanner.Token
	depth int
}

type compiler struct {
	scanner    scanner.Scanner
	previous   *scanner.Token
	current    *scanner.Token
	locals     []local
	scopeDepth int
	chunk      *Chunk
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

func (c *compiler) emitLoop(start int) error {
	c.emitOperation(OperationLoop)
	offset := len(c.chunk.Code) - start + 2
	if offset > math.MaxUint16 {
		return &Error{ErrTooBigLoop, c.previous}
	}
	c.chunk.write(byte((offset >> 8) & 0xff))
	c.chunk.write(byte(offset & 0xff))
	return nil
}

func (c *compiler) emitJump(o Operation) int {
	c.emitOperation(o)
	c.chunk.write(0xff)
	c.chunk.write(0xff)
	return len(c.chunk.Code) - 2
}

func (c *compiler) patchJump(offset int) error {
	jump := len(c.chunk.Code) - offset - 2
	if jump > math.MaxUint16 {
		return &Error{ErrTooBigJump, c.previous}
	}
	c.chunk.Code[offset] = byte((jump >> 8) & 0xff)
	c.chunk.Code[offset+1] = byte((jump & 0xff))
	return nil
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

func (c *compiler) beginScope() {
	c.scopeDepth++
}

func (c *compiler) endScope() {
	c.scopeDepth--
	for len(c.locals) > 0 && c.locals[len(c.locals)-1].depth > c.scopeDepth {
		c.emitOperation(OperationPop)
		c.locals = c.locals[:len(c.locals)-1]
	}
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

func (c *compiler) resolveLocal(t *scanner.Token) (int, error) {
	for i := len(c.locals) - 1; i >= 0; i-- {
		if t.Lexeme == c.locals[i].name.Lexeme {
			if c.locals[i].depth == -1 {
				return 0, &Error{ErrVarOwnInitializer, c.previous}
			}
			return i, nil
		}
	}
	return -1, nil
}

func (c *compiler) namedVariable(t *scanner.Token, canAssign bool) error {
	var getOp, setOp Operation
	i, err := c.resolveLocal(t)
	if err != nil {
		return err
	}
	if i != -1 {
		getOp = OperationGetLocal
		setOp = OperationSetLocal
	} else {
		x, err := c.identifierConstant(t)
		if err != nil {
			return err
		}
		i = int(x)
		getOp = OperationGetGlobal
		setOp = OperationSetGlobal
	}
	if canAssign && c.check(scanner.TokenEqual) {
		if err = c.advance(); err != nil {
			return err
		}
		if err = c.expression(); err != nil {
			return err
		}
		c.emitOperation(setOp)
		c.chunk.write(uint8(i))
	} else {
		c.emitOperation(getOp)
		c.chunk.write(uint8(i))
	}
	return nil
}

func (c *compiler) variable(canAssign bool) error {
	return c.namedVariable(c.previous, canAssign)
}

func (c *compiler) and() error {
	endJump := c.emitJump(OperationJumpIfFalse)
	c.emitOperation(OperationPop)
	if err := c.parsePrecedence(precedenceAnd); err != nil {
		return err
	}
	return c.patchJump(endJump)
}

func (c *compiler) or() error {
	elseJump := c.emitJump(OperationJumpIfFalse)
	endJump := c.emitJump(OperationJump)
	var err error
	if err = c.patchJump(elseJump); err != nil {
		return err
	}
	c.emitOperation(OperationPop)
	if err = c.parsePrecedence(precedenceOr); err != nil {
		return err
	}
	return c.patchJump(endJump)
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
	case parseFunctionAnd:
		return c.and()
	case parseFunctionOr:
		return c.or()
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

func (c *compiler) ifStatement() error {
	var err error
	if err = c.consume(scanner.TokenLeftParen, ErrIfLeftParen); err != nil {
		return err
	}
	if err = c.expression(); err != nil {
		return err
	}
	if err = c.consume(scanner.TokenRightParen, ErrIfRightParen); err != nil {
		return err
	}
	thenJump := c.emitJump(OperationJumpIfFalse)
	c.emitOperation(OperationPop)
	if err = c.statement(); err != nil {
		return err
	}
	elseJump := c.emitJump(OperationJump)
	if err = c.patchJump(thenJump); err != nil {
		return err
	}
	c.emitOperation(OperationPop)
	if c.check(scanner.TokenElse) {
		if err = c.advance(); err != nil {
			return err
		}
		if err = c.statement(); err != nil {
			return err
		}
	}
	return c.patchJump(elseJump)
}

func (c *compiler) whileStatement() error {
	loopStart := len(c.chunk.Code)
	var err error
	if err = c.consume(scanner.TokenLeftParen, ErrWhileLeftParen); err != nil {
		return err
	}
	if err = c.expression(); err != nil {
		return err
	}
	if err = c.consume(scanner.TokenRightParen, ErrWhileRightParen); err != nil {
		return err
	}
	exitJump := c.emitJump(OperationJumpIfFalse)
	c.emitOperation(OperationPop)
	if err = c.statement(); err != nil {
		return err
	}
	if err = c.emitLoop(loopStart); err != nil {
		return err
	}
	if err = c.patchJump(exitJump); err != nil {
		return err
	}
	c.emitOperation(OperationPop)
	return nil
}

func (c *compiler) forStatement() error {
	c.beginScope()
	var err error
	if err = c.consume(scanner.TokenLeftParen, ErrForLeftParen); err != nil {
		return err
	}
	if c.check(scanner.TokenSemicolon) {
		if err = c.advance(); err != nil {
			return err
		}
	} else if c.check(scanner.TokenVar) {
		if err = c.advance(); err != nil {
			return err
		}
		if err = c.varDeclaration(); err != nil {
			return err
		}
	} else {
		if err = c.expressionStatement(); err != nil {
			return err
		}
	}
	loopStart := len(c.chunk.Code)
	exitJump := -1
	if c.check(scanner.TokenSemicolon) {
		if err = c.advance(); err != nil {
			return err
		}
	} else {
		if err = c.expression(); err != nil {
			return err
		}
		if err = c.consume(scanner.TokenSemicolon, ErrForConditionSemicolon); err != nil {
			return err
		}
		exitJump = c.emitJump(OperationJumpIfFalse)
		c.emitOperation(OperationPop)
	}
	if c.check(scanner.TokenRightParen) {
		if err = c.advance(); err != nil {
			return err
		}
	} else {
		bodyJump := c.emitJump(OperationJump)
		incrementStart := len(c.chunk.Code)
		if err = c.expression(); err != nil {
			return err
		}
		c.emitOperation(OperationPop)
		if err = c.consume(scanner.TokenRightParen, ErrForRightParen); err != nil {
			return err
		}
		if err = c.emitLoop(loopStart); err != nil {
			return err
		}
		loopStart = incrementStart
		if err = c.patchJump(bodyJump); err != nil {
			return err
		}
	}
	if err = c.statement(); err != nil {
		return err
	}
	if err = c.emitLoop(loopStart); err != nil {
		return err
	}
	if exitJump != -1 {
		if err = c.patchJump(exitJump); err != nil {
			return err
		}
		c.emitOperation(OperationPop)
	}
	c.endScope()
	return nil
}

func (c *compiler) block() error {
	var err error
	for !c.check(scanner.TokenRightBrace) && !c.check(scanner.TokenEof) {
		if err = c.declaration(); err != nil {
			return err
		}
	}
	return c.consume(scanner.TokenRightBrace, ErrMissingBlockRightBrace)
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
	switch {
	case c.check(scanner.TokenPrint):
		if err := c.advance(); err != nil {
			return err
		}
		return c.printStatement()
	case c.check(scanner.TokenIf):
		if err := c.advance(); err != nil {
			return err
		}
		return c.ifStatement()
	case c.check(scanner.TokenWhile):
		if err := c.advance(); err != nil {
			return err
		}
		return c.whileStatement()
	case c.check(scanner.TokenFor):
		if err := c.advance(); err != nil {
			return err
		}
		return c.forStatement()
	case c.check(scanner.TokenLeftBrace):
		if err := c.advance(); err != nil {
			return err
		}
		c.beginScope()
		if err := c.block(); err != nil {
			return err
		}
		c.endScope()
		return nil
	default:
		return c.expressionStatement()
	}
}

func (c *compiler) identifierConstant(t *scanner.Token) (uint8, error) {
	return c.makeConstant(NewString(t.Lexeme))
}

func (c *compiler) addLocal(name *scanner.Token) error {
	if len(c.locals) > math.MaxUint8 {
		return &Error{ErrTooManyLocals, c.previous}
	}
	c.locals = append(c.locals, local{name, -1})
	return nil
}

func (c *compiler) declareVariable() error {
	if c.scopeDepth == 0 {
		return nil
	}
	for i := len(c.locals) - 1; i >= 0; i-- {
		if c.locals[i].depth != -1 && c.locals[i].depth < c.scopeDepth {
			break
		}
		if c.previous.Lexeme == c.locals[i].name.Lexeme {
			return &Error{ErrVarAlreadyDefined, c.previous}
		}
	}
	return c.addLocal(c.previous)
}

func (c *compiler) parseVariable(k ErrorKind) (uint8, error) {
	var err error
	if err = c.consume(scanner.TokenIdentifier, k); err != nil {
		return 0, err
	}
	if err = c.declareVariable(); err != nil {
		return 0, err
	}
	if c.scopeDepth > 0 {
		return 0, nil
	}
	return c.identifierConstant(c.previous)
}

func (c *compiler) markInitialized() {
	c.locals[len(c.locals)-1].depth = c.scopeDepth
}

func (c *compiler) defineVariable(v uint8) {
	if c.scopeDepth > 0 {
		c.markInitialized()
		return
	}
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
