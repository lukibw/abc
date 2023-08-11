package abc

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type precedence int

const (
	precedenceNone precedence = iota
	precedenceAssignment
	precedenceOr
	precedenceAnd
	precedenceEquality
	precedenceComparison
	precedenceTerm
	precedenceFactor
	precedenceUnary
	precedenceCall
	precedencePrimary
)

type parseFunction int

const (
	parseFunctionNone parseFunction = iota
	parseFunctionNumber
	parseFunctionBinary
	parseFunctionUnary
	parseFunctionGrouping
	parseFunctionLiteral
)

type parseRule struct {
	prefix     parseFunction
	infix      parseFunction
	precedence precedence
}

var parseRules = map[tokenKind]parseRule{
	tokenLeftParen:    {parseFunctionGrouping, parseFunctionNone, precedenceNone},
	tokenRightParen:   {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenLeftBrace:    {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenRightBrace:   {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenComma:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenDot:          {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenMinus:        {parseFunctionUnary, parseFunctionBinary, precedenceTerm},
	tokenPlus:         {parseFunctionNone, parseFunctionBinary, precedenceTerm},
	tokenSemicolon:    {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenSlash:        {parseFunctionNone, parseFunctionBinary, precedenceFactor},
	tokenStar:         {parseFunctionNone, parseFunctionBinary, precedenceFactor},
	tokenBang:         {parseFunctionUnary, parseFunctionNone, precedenceNone},
	tokenBangEqual:    {parseFunctionNone, parseFunctionBinary, precedenceEquality},
	tokenEqual:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenEqualEqual:   {parseFunctionNone, parseFunctionBinary, precedenceEquality},
	tokenGreater:      {parseFunctionNone, parseFunctionBinary, precedenceComparison},
	tokenGreaterEqual: {parseFunctionNone, parseFunctionBinary, precedenceComparison},
	tokenLess:         {parseFunctionNone, parseFunctionBinary, precedenceComparison},
	tokenLessEqual:    {parseFunctionNone, parseFunctionBinary, precedenceComparison},
	tokenIdentifier:   {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenString:       {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenNumber:       {parseFunctionNumber, parseFunctionNone, precedenceNone},
	tokenAnd:          {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenClass:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenElse:         {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenFalse:        {parseFunctionLiteral, parseFunctionNone, precedenceNone},
	tokenFor:          {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenFun:          {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenIf:           {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenNil:          {parseFunctionLiteral, parseFunctionNone, precedenceNone},
	tokenOr:           {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenPrint:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenReturn:       {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenSuper:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenThis:         {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenTrue:         {parseFunctionLiteral, parseFunctionNone, precedenceNone},
	tokenVar:          {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenWhile:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	tokenEof:          {parseFunctionNone, parseFunctionNone, precedenceNone},
}

type parserErrorKind int

const (
	errTooManyConstants parserErrorKind = iota
	errMissingExpr
	errMissingExprRightParen
)

var parserErrorMessages = map[parserErrorKind]string{
	errTooManyConstants:      "too many constants in one chunk",
	errMissingExpr:           "missing expression",
	errMissingExprRightParen: "missing ')' after expression",
}

func (k parserErrorKind) String() string {
	return parserErrorMessages[k]
}

type parserError struct {
	kind  parserErrorKind
	token *token
}

func (e *parserError) Error() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("[line %d] error", e.token.line))
	if e.token.kind == tokenEof {
		sb.WriteString("at end")
	} else {
		sb.WriteString(fmt.Sprintf(" at '%d'", e.token.start))
	}
	sb.WriteString(fmt.Sprintf(": %s", e.kind))
	return sb.String()
}

type parser struct {
	chunk    *Chunk
	scanner  *scanner
	previous *token
	current  *token
}

func newParser(source []rune) *parser {
	return &parser{NewChunk(), newScanner(source), nil, nil}
}

func (p *parser) newError(kind parserErrorKind) error {
	return &parserError{kind, p.previous}
}

func (p *parser) advance() error {
	p.previous = p.current
	t, err := p.scanner.token()
	if err != nil {
		return err
	}
	p.current = t
	return nil
}

func (p *parser) consume(t tokenKind, e parserErrorKind) error {
	if p.current.kind == t {
		p.advance()
		return nil
	}
	return &parserError{e, p.current}
}

func (p *parser) makeConstant(f float64) (byte, error) {
	constant := p.chunk.AddConstant(f)
	if constant > math.MaxUint8 {
		return 0, p.newError(errTooManyConstants)
	}
	return constant, nil
}

func (p *parser) emitByte(b byte) {
	p.chunk.Write(b, p.previous.line)
}

func (p *parser) emitBytes(b1, b2 byte) {
	p.emitByte(b1)
	p.emitByte(b2)
}

func (p *parser) emitReturn() {
	p.emitByte(byte(OpReturn))
}

func (p *parser) endCompiler() {
	p.emitReturn()
	fmt.Println(p.chunk)
}

func (p *parser) emitConstant(f float64) error {
	constant, err := p.makeConstant(f)
	if err != nil {
		return err
	}
	p.emitBytes(byte(OpConstant), constant)
	return nil
}

func (p *parser) binary() error {
	switch p.previous.kind {
	case tokenPlus:
		p.emitByte(byte(OpAdd))
	case tokenMinus:
		p.emitByte(byte(OpSubtract))
	case tokenStar:
		p.emitByte(byte(OpMultiply))
	case tokenSlash:
		p.emitByte(byte(OpDivide))
	case tokenBangEqual:
		p.emitBytes(byte(OpEqual), byte(OpNot))
	case tokenEqualEqual:
		p.emitByte(byte(OpEqual))
	case tokenGreater:
		p.emitByte(byte(OpGreater))
	case tokenGreaterEqual:
		p.emitBytes(byte(OpLess), byte(OpNot))
	case tokenLess:
		p.emitByte(byte(OpLess))
	case tokenLessEqual:
		p.emitBytes(byte(OpGreater), byte(OpNot))
	}
	return nil
}

func (p *parser) literal() error {
	switch p.previous.kind {
	case tokenNil:
		p.emitByte(byte(OpNil))
	case tokenFalse:
		p.emitByte(byte(OpFalse))
	case tokenTrue:
		p.emitByte(byte(OpTrue))
	}
	return nil
}

func (p *parser) number() error {
	f, err := strconv.ParseFloat(p.scanner.lexeme(p.previous), 64)
	if err != nil {
		panic("parser: cannot parse float from token lexeme")
	}
	return p.emitConstant(f)
}

func (p *parser) parseFunction(f parseFunction) error {
	switch f {
	case parseFunctionBinary:
		return p.binary()
	case parseFunctionGrouping:
		return p.grouping()
	case parseFunctionUnary:
		return p.unary()
	case parseFunctionNumber:
		return p.number()
	case parseFunctionLiteral:
		return p.literal()
	default:
		return p.newError(errMissingExpr)
	}
}

func (p *parser) parsePrecedence(min precedence) error {
	var err error
	if err = p.advance(); err != nil {
		return err
	}
	if err = p.parseFunction(parseRules[p.previous.kind].prefix); err != nil {
		return err
	}
	for min <= parseRules[p.current.kind].precedence {
		if err = p.advance(); err != nil {
			return err
		}
		if err = p.parseFunction(parseRules[p.previous.kind].infix); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) expression() error {
	return p.parsePrecedence(precedenceAssignment)
}

func (p *parser) grouping() error {
	if err := p.expression(); err != nil {
		return err
	}
	return p.consume(tokenRightParen, errMissingExprRightParen)
}

func (p *parser) unary() error {
	if err := p.parsePrecedence(precedenceUnary); err != nil {
		return err
	}
	switch p.previous.kind {
	case tokenMinus:
		p.emitByte(byte(OpNegate))
	case tokenBang:
		p.emitByte(byte(OpNot))
	}
	return nil
}
