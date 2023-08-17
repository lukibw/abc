package scanner

import "fmt"

type TokenKind int

const (
	TokenLeftParen TokenKind = iota
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenDot
	TokenMinus
	TokenPlus
	TokenSemicolon
	TokenSlash
	TokenStar
	TokenBang
	TokenBangEqual
	TokenEqual
	TokenEqualEqual
	TokenGreater
	TokenGreaterEqual
	TokenLess
	TokenLessEqual
	TokenIdentifier
	TokenString
	TokenNumber
	TokenAnd
	TokenClass
	TokenElse
	TokenFalse
	TokenFor
	TokenFun
	TokenIf
	TokenNil
	TokenOr
	TokenPrint
	TokenReturn
	TokenSuper
	TokenThis
	TokenTrue
	TokenVar
	TokenWhile
	TokenEof
)

var tokenKinds = map[TokenKind]string{
	TokenLeftParen:    "(",
	TokenRightParen:   ")",
	TokenLeftBrace:    "[",
	TokenRightBrace:   "]",
	TokenComma:        ",",
	TokenDot:          ".",
	TokenMinus:        "-",
	TokenPlus:         "+",
	TokenSemicolon:    ";",
	TokenSlash:        "/",
	TokenStar:         "*",
	TokenBang:         "!",
	TokenBangEqual:    "!=",
	TokenEqual:        "=",
	TokenEqualEqual:   "==",
	TokenGreater:      ">",
	TokenGreaterEqual: ">=",
	TokenLess:         "<",
	TokenLessEqual:    "<=",
	TokenIdentifier:   "IDENTIFIER",
	TokenString:       "STRING",
	TokenNumber:       "NUMBER",
	TokenAnd:          "and",
	TokenClass:        "class",
	TokenElse:         "else",
	TokenFalse:        "false",
	TokenFor:          "for",
	TokenFun:          "fun",
	TokenIf:           "if",
	TokenNil:          "nil",
	TokenOr:           "or",
	TokenPrint:        "print",
	TokenReturn:       "return",
	TokenSuper:        "super",
	TokenThis:         "this",
	TokenTrue:         "true",
	TokenVar:          "var",
	TokenWhile:        "while",
	TokenEof:          "EOF",
}

func (k TokenKind) String() string {
	return tokenKinds[k]
}

type Token struct {
	Kind   TokenKind
	Line   int
	Lexeme string
}

func (t *Token) String() string {
	if t.Kind == TokenIdentifier || t.Kind == TokenNumber || t.Kind == TokenString {
		return fmt.Sprintf("'%s' (%s)", t.Lexeme, t.Kind)
	}
	return fmt.Sprintf("'%s'", t.Kind)
}
