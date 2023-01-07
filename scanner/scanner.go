package scanner

import (
	"glox/errors"
	"glox/token"
	"strconv"
)

type Scanner struct {
    source string
    tokens []token.Token
    start int
    current int
    line int
}

func NewScanner(source string) Scanner {
    return Scanner{source: source, tokens: []token.Token{}, start: 0, current: 0, line: 1}
}

func (s *Scanner) ScanTokens() []token.Token {
    for !s.isAtEnd() {
        s.start = s.current
        s.scanToken()
    }

    s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))

    return s.tokens
}

func (s Scanner) isAtEnd() bool {
    return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
    c := s.advance()

    switch c {
    case '(':
        s.addToken(token.LEFT_PAREN)
    case ')':
        s.addToken(token.RIGHT_PAREN)
    case '{':
        s.addToken(token.LEFT_BRACE)
    case '}':
        s.addToken(token.RIGHT_BRACE)
    case ',':
        s.addToken(token.COMMA)
    case '.':
        s.addToken(token.DOT)
    case '-':
        s.addToken(token.MINUS)
    case '+':
        s.addToken(token.PLUS)
    case ';':
        s.addToken(token.SEMICOLON)
    case '*':
        s.addToken(token.STAR)
    case '!':
        t := token.BANG
        if s.match('=') {
            t = token.BANG_EQUAL
        }
        s.addToken(t);
    case '=':
        t := token.EQUAL
        if s.match('=') {
            t = token.EQUAL_EQUAL
        }
        s.addToken(t);
    case '<':
        t := token.LESS
        if s.match('=') {
            t = token.LESS_EQUAL
        }
        s.addToken(t);
    case '>':
        t := token.GREATER
        if s.match('=') {
            t = token.GREATER_EQUAL
        }
        s.addToken(t);

    case '/':
        if s.match('/') {
            for s.peek() != '\n' && !s.isAtEnd() {
                s.advance()
            }
        } else {
            s.addToken(token.SLASH)
        }

    case ' ', '\r', '\t':
        break

    case '\n':
        s.line++

    case '"':
        s.string()

    default:
        if isDigit(c) {
            s.number()
        } else if isAlpha(c) {
            s.identifier()
        }else {
            errors.ErrorAt(s.line, "Unexpected character.")
        }
    }
}

func (s *Scanner) advance() byte {
    c := s.source[s.current]
    s.current++

    return c
}

func (s *Scanner) addToken(tokenType token.TokenType) {
    s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal interface{}) {
    text := s.source[s.start:s.current]
    s.tokens = append(s.tokens, token.NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) match(expected byte) bool {
    if s.isAtEnd() {
        return false
    }

    if s.source[s.current] != expected {
        return false
    }

    s.current++

    return true
}

func (s Scanner) peek() byte {
    if s.isAtEnd() {
        return 0
    }

    return s.source[s.current]
}

func (s Scanner) peekNext() byte {
    if s.current + 1 >= len(s.source) {
        return 0
    }

    return s.source[s.current + 1]
}

func (s *Scanner) string() {
    for s.peek() != '"' && !s.isAtEnd() {
        if s.peek() == '\n' {
            s.line++
        }

        s.advance()
    }

    if s.isAtEnd() {
        errors.ErrorAt(s.line, "Unterminated string.")
        return
    }

    s.advance()
    
    value := s.source[s.start + 1: s.current - 1]

    s.addTokenWithLiteral(token.STRING, value)
}

func (s *Scanner) number() {
    for isDigit(s.peek()) {
        s.advance()
    }

    if s.peek() == '.' && isDigit(s.peekNext()) {
        s.advance()

        for isDigit(s.peek()) {
            s.advance()
        }
    }

    if f, err := strconv.ParseFloat(s.source[s.start:s.current], 64); err == nil {
        s.addTokenWithLiteral(token.NUMBER, f)
    }
}

func (s *Scanner) identifier() {
    for isAlphaNumeric(s.peek()) {
        s.advance()
    }
    
    text := s.source[s.start:s.current]
    tokenType, isKeyword := keywords[text]

    if !isKeyword {
        tokenType = token.IDENTIFIER
    }

    s.addToken(tokenType)
}
 
func isDigit(c byte) bool {
    return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func isAlphaNumeric(c byte) bool {
    return isDigit(c) || isAlpha(c)
}
