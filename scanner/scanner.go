package scanner

type Scanner interface {
	Token() (*Token, error)
}

func New(source []byte) Scanner {
	return &scanner{source, 0, 0, 1}
}

type scanner struct {
	source  []byte
	start   int
	current int
	line    int
}

func (s *scanner) newToken(k TokenKind) *Token {
	return &Token{k, s.line, string(s.source[s.start:s.current])}
}

func (s *scanner) newError(k ErrorKind) error {
	return &Error{k, s.line}
}

func (s *scanner) isAtEnd() bool {
	return s.current == len(s.source)
}

func (s *scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *scanner) peekNext() byte {
	if s.current > len(s.source)-2 {
		return 0
	}
	return s.source[s.current+1]
}

func (s *scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *scanner) match(b byte) bool {
	if s.isAtEnd() || s.peek() != b {
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

func (s *scanner) string() (*Token, error) {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return nil, s.newError(ErrUnterminatedString)
	}
	s.advance()
	return s.newToken(TokenString), nil
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (s *scanner) number() (*Token, error) {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	return s.newToken(TokenNumber), nil
}

var keywords = map[string]TokenKind{
	"and":    TokenAnd,
	"class":  TokenClass,
	"else":   TokenElse,
	"if":     TokenIf,
	"nil":    TokenNil,
	"or":     TokenOr,
	"print":  TokenPrint,
	"return": TokenReturn,
	"super":  TokenSuper,
	"var":    TokenVar,
	"while":  TokenWhile,
	"false":  TokenFalse,
	"for":    TokenFor,
	"fun":    TokenFun,
	"this":   TokenThis,
	"true":   TokenTrue,
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func (s *scanner) identifier() (*Token, error) {
	for isAlpha(s.peek()) || isDigit(s.peek()) {
		s.advance()
	}
	if keyword, ok := keywords[string(s.source[s.start:s.current])]; ok {
		return s.newToken(keyword), nil
	}
	return s.newToken(TokenIdentifier), nil
}

func (s *scanner) Token() (*Token, error) {
	s.skipWhitespace()
	s.start = s.current
	if s.isAtEnd() {
		return s.newToken(TokenEof), nil
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
		return s.newToken(TokenLeftParen), nil
	case ')':
		return s.newToken(TokenRightParen), nil
	case '{':
		return s.newToken(TokenLeftBrace), nil
	case '}':
		return s.newToken(TokenRightBrace), nil
	case ';':
		return s.newToken(TokenSemicolon), nil
	case ',':
		return s.newToken(TokenComma), nil
	case '.':
		return s.newToken(TokenDot), nil
	case '-':
		return s.newToken(TokenMinus), nil
	case '+':
		return s.newToken(TokenPlus), nil
	case '/':
		return s.newToken(TokenSlash), nil
	case '*':
		return s.newToken(TokenStar), nil
	case '!':
		if s.match('=') {
			return s.newToken(TokenBangEqual), nil
		} else {
			return s.newToken(TokenBang), nil
		}
	case '=':
		if s.match('=') {
			return s.newToken(TokenEqualEqual), nil
		} else {
			return s.newToken(TokenEqual), nil
		}
	case '<':
		if s.match('=') {
			return s.newToken(TokenLessEqual), nil
		} else {
			return s.newToken(TokenLess), nil
		}
	case '>':
		if s.match('=') {
			return s.newToken(TokenGreaterEqual), nil
		} else {
			return s.newToken(TokenGreater), nil
		}
	case '"':
		return s.string()
	default:
		return nil, s.newError(ErrUnexpectedCharacter)
	}
}
