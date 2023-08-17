package compiler

import "github.com/lukibw/abc/scanner"

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

var parseRules = map[scanner.TokenKind]parseRule{
	scanner.TokenLeftParen:    {parseFunctionGrouping, parseFunctionNone, precedenceNone},
	scanner.TokenRightParen:   {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenLeftBrace:    {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenRightBrace:   {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenComma:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenDot:          {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenMinus:        {parseFunctionUnary, parseFunctionBinary, precedenceTerm},
	scanner.TokenPlus:         {parseFunctionNone, parseFunctionBinary, precedenceTerm},
	scanner.TokenSemicolon:    {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenSlash:        {parseFunctionNone, parseFunctionBinary, precedenceFactor},
	scanner.TokenStar:         {parseFunctionNone, parseFunctionBinary, precedenceFactor},
	scanner.TokenBang:         {parseFunctionUnary, parseFunctionNone, precedenceNone},
	scanner.TokenBangEqual:    {parseFunctionNone, parseFunctionBinary, precedenceEquality},
	scanner.TokenEqual:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenEqualEqual:   {parseFunctionNone, parseFunctionBinary, precedenceEquality},
	scanner.TokenGreater:      {parseFunctionNone, parseFunctionBinary, precedenceComparison},
	scanner.TokenGreaterEqual: {parseFunctionNone, parseFunctionBinary, precedenceComparison},
	scanner.TokenLess:         {parseFunctionNone, parseFunctionBinary, precedenceComparison},
	scanner.TokenLessEqual:    {parseFunctionNone, parseFunctionBinary, precedenceComparison},
	scanner.TokenIdentifier:   {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenString:       {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenNumber:       {parseFunctionNumber, parseFunctionNone, precedenceNone},
	scanner.TokenAnd:          {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenClass:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenElse:         {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenFalse:        {parseFunctionLiteral, parseFunctionNone, precedenceNone},
	scanner.TokenFor:          {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenFun:          {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenIf:           {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenNil:          {parseFunctionLiteral, parseFunctionNone, precedenceNone},
	scanner.TokenOr:           {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenPrint:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenReturn:       {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenSuper:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenThis:         {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenTrue:         {parseFunctionLiteral, parseFunctionNone, precedenceNone},
	scanner.TokenVar:          {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenWhile:        {parseFunctionNone, parseFunctionNone, precedenceNone},
	scanner.TokenEof:          {parseFunctionNone, parseFunctionNone, precedenceNone},
}
