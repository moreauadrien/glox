package token

import "fmt"

type Token struct {
    tokenType TokenType
    Lexeme string
    literal interface{}
    line int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token{
    return Token{tokenType: tokenType, Lexeme: lexeme, literal: literal, line: line}
}

func (t Token) String() string {
    return fmt.Sprintf("%v %v %v", t.tokenType, t.Lexeme, t.literal)
}
