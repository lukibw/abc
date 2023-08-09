package abc

import (
	"fmt"
)

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

type tokenKind int

const (
	tokenLeftParen tokenKind = iota
	tokenRightParen
	tokenLeftBrace
	tokenRightBrace
	tokenComma
	tokenDot
	tokenMinus
	tokenPlus
	tokenSemicolon
	tokenSlash
	tokenStar
	tokenBang
	tokenBangEqual
	tokenEqual
	tokenEqualEqual
	tokenGreater
	tokenGreaterEqual
	tokenLess
	tokenLessEqual
	tokenIdentifier
	tokenString
	tokenNumber
	tokenAnd
	tokenClass
	tokenElse
	tokenFalse
	tokenFor
	tokenFun
	tokenIf
	tokenNil
	tokenOr
	tokenPrint
	tokenReturn
	tokenSuper
	tokenThis
	tokenTrue
	tokenVar
	tokenWhile
	tokenEof
)

var keywords = map[string]tokenKind{
	"and":    tokenAnd,
	"class":  tokenClass,
	"else":   tokenElse,
	"if":     tokenIf,
	"nil":    tokenNil,
	"or":     tokenOr,
	"print":  tokenPrint,
	"return": tokenReturn,
	"super":  tokenSuper,
	"var":    tokenVar,
	"while":  tokenWhile,
}

type token struct {
	kind   tokenKind
	start  int
	length int
	line   int
}

type scannerErrorKind int

const (
	errUnexpectedCharacter scannerErrorKind = iota
	errUnterminatedString
)

var scannerErrorMessages = map[scannerErrorKind]string{
	errUnexpectedCharacter: "unexpected character",
	errUnterminatedString:  "unterminated string",
}

func (k scannerErrorKind) String() string {
	return scannerErrorMessages[k]
}

type scannerError struct {
	kind scannerErrorKind
	line int
}

func (e *scannerError) Error() string {
	return fmt.Sprintf("[line %d] %s", e.line, e.kind)
}

type scanner struct {
	source  []rune
	start   int
	current int
	line    int
}

func newScanner(source []rune) *scanner {
	return &scanner{source, 0, 0, 1}
}

func (s *scanner) newToken(kind tokenKind) token {
	return token{kind, s.start, s.current - s.start, s.line}
}

func (s *scanner) newError(kind scannerErrorKind) error {
	return &scannerError{kind, s.line}
}

func (s *scanner) isAtEnd() bool {
	return s.current == len(s.source)
}

func (s *scanner) peek() rune {
	return s.source[s.current]
}

func (s *scanner) peekNext() rune {
	if s.current > len(s.source)-2 {
		return 0
	}
	return s.source[s.current+1]
}

func (s *scanner) advance() rune {
	s.current++
	return s.source[s.current-1]
}

func (s *scanner) match(r rune) bool {
	if s.isAtEnd() || s.peek() != r {
		return false
	}
	s.current++
	return true
}

func (s *scanner) skipWhitespace() {
	for {
		switch s.peek() {
		case ' ', '\t', '\r':
			s.advance()
		case '\n':
			s.advance()
			s.line++
		case '/':
			if s.peekNext() == '/' {
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (s *scanner) string() (token, error) {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() != '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return token{}, s.newError(errUnterminatedString)
	}
	s.advance()
	return s.newToken(tokenString), nil
}

func (s *scanner) number() (token, error) {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	return s.newToken(tokenNumber), nil
}

func (s *scanner) identifier() (token, error) {
	for isAlpha(s.peek()) || isDigit(s.peek()) {
		s.advance()
	}
	if keyword, ok := keywords[string(s.source[s.start:s.current])]; ok {
		return s.newToken(keyword), nil
	}
	return s.newToken(tokenIdentifier), nil
}

func (s *scanner) token() (token, error) {
	s.skipWhitespace()
	s.start = s.current
	if s.isAtEnd() {
		return s.newToken(tokenEof), nil
	}
	r := s.advance()
	if isDigit(r) {
		return s.number()
	}
	if isAlpha(r) {
		return s.identifier()
	}
	switch r {
	case '(':
		return s.newToken(tokenLeftParen), nil
	case ')':
		return s.newToken(tokenRightParen), nil
	case '{':
		return s.newToken(tokenLeftBrace), nil
	case '}':
		return s.newToken(tokenRightBrace), nil
	case ';':
		return s.newToken(tokenSemicolon), nil
	case ',':
		return s.newToken(tokenComma), nil
	case '.':
		return s.newToken(tokenDot), nil
	case '-':
		return s.newToken(tokenMinus), nil
	case '+':
		return s.newToken(tokenPlus), nil
	case '/':
		return s.newToken(tokenSlash), nil
	case '*':
		return s.newToken(tokenStar), nil
	case '!':
		if s.match('=') {
			return s.newToken(tokenBangEqual), nil
		} else {
			return s.newToken(tokenBang), nil
		}
	case '=':
		if s.match('=') {
			return s.newToken(tokenEqualEqual), nil
		} else {
			return s.newToken(tokenEqual), nil
		}
	case '<':
		if s.match('=') {
			return s.newToken(tokenLessEqual), nil
		} else {
			return s.newToken(tokenLess), nil
		}
	case '>':
		if s.match('=') {
			return s.newToken(tokenGreaterEqual), nil
		} else {
			return s.newToken(tokenGreater), nil
		}
	case '"':
		return s.string()
	default:
		return token{}, s.newError(errUnexpectedCharacter)
	}
}
