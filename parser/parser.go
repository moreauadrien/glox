package parser

import (
	"fmt"
	"glox/errors"
	"glox/token"
	"glox/tree"
)

type Parser struct {
    tokens []token.Token
    current int
}

func NewParser(tokens []token.Token) Parser {
    return Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() []tree.Stmt {
    statements := []tree.Stmt{}

    for !p.isAtEnd() {
        statements = append(statements, p.declaration())
    }

    return statements
}

func (p *Parser) declaration() tree.Stmt {
    var stmt tree.Stmt
    var err error

    if p.match(token.VAR) {
        stmt, err = p.varDeclaration()
    } else {
        stmt, err = p.statements()
    }

    if err != nil {
        p.synchronize()
        return nil
    }

    return stmt
}


func (p *Parser) statements() (tree.Stmt, error) {
    if p.match(token.PRINT) {
        return p.printStatement()
    }

    if p.match(token.LEFT_BRACE) {
        block, err := p.block()

        if err != nil {
            return nil, err
        }

        return tree.NewBlock(block), nil
    }

    return p.expressionStatement()
}

func (p *Parser) printStatement() (tree.Stmt, error) {
    val, err := p.expression()
    if err != nil {
        return nil, err
    }

    if _, err = p.consume(token.SEMICOLON, "Expected ';' after value."); err != nil {
        return nil, err
    }

    return tree.NewPrint(val), nil
}

func (p *Parser) block() ([]tree.Stmt, error) {
    statements := []tree.Stmt{}

    for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
        statements = append(statements, p.declaration())
    }

    if _, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block."); err != nil {
        return nil, err
    }

    return statements, nil
}

func (p *Parser) varDeclaration() (tree.Stmt, error) {
    name, err := p.consume(token.IDENTIFIER, "Expect variable name.")

    if err != nil {
        return nil, err
    }

    var initializer tree.Expr
    if p.match(token.EQUAL) {
        initializer, err = p.expression()

        if err != nil {
            return nil, err
        }
    }

    _, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")

    if err != nil {
        return nil, err
    }

    return tree.NewVar(name, initializer), nil
}

func (p *Parser) expressionStatement() (tree.Stmt, error){
    expr, err := p.expression()

    if err != nil {
        return nil, err
    }

    if _, err = p.consume(token.SEMICOLON, "Expected ';' after value."); err != nil {
        return nil, err
    }

    return tree.NewExpression(expr), nil
}

func (p *Parser) expression() (tree.Expr, error) {
    return p.assignement()
}

func (p *Parser) assignement() (tree.Expr, error) {
    expr, err := p.equality()

    if err != nil {
        return nil, err
    }

    if p.match(token.EQUAL) {
        equals := p.previous()
        value, err := p.assignement()

        if err != nil {
            return nil, err
        }

        if variable, ok := expr.(*tree.Variable); ok {
            name := variable.Name
            return tree.NewAssign(name, value), nil
        }

        errors.Error(equals, "Invalid assignement target.")
    }

    return expr, nil
}

func (p *Parser) equality() (tree.Expr, error) {
    expr, err := p.comparison()

    if err != nil {
        return nil, err
    }

    for p.match(token.EQUAL_EQUAL, token.BANG_EQUAL) {
        op := p.previous()
        right, err := p.comparison()

        if err != nil {
            return nil, err
        }

        expr = tree.NewBinary(expr, op, right)
    }

    return expr, nil
}

func (p *Parser) comparison() (tree.Expr, error) {
    expr, err := p.term()
    
    if err != nil {
        return nil, err
    }

    for p.match(token.LESS, token.LESS_EQUAL, token.GREATER, token.GREATER_EQUAL) {
        op := p.previous()
        right, err := p.term()

        if err != nil {
            return nil, err
        }

        expr = tree.NewBinary(expr, op, right)
    }

    return expr, nil
}

func (p *Parser) term() (tree.Expr, error) {
    expr, err := p.factor()

    if err != nil {
        return nil, err
    }

    for p.match(token.PLUS, token.MINUS) {
        op := p.previous()
        right, err := p.factor()

        if err != nil {
            return nil, err
        }

        expr = tree.NewBinary(expr, op, right)
    }

    return expr, nil
}

func (p *Parser) factor() (tree.Expr, error) {
    expr, err := p.unary()

    if err != nil {
        return nil, err
    }

    for p.match(token.SLASH, token.STAR) {
        op := p.previous()
        right, err := p.unary()

        if err != nil {
            return nil, err
        }

        expr = tree.NewBinary(expr, op, right)
    }

    return expr, nil
}

func (p *Parser) unary() (tree.Expr, error) {
    if p.match(token.BANG, token.MINUS) {
        expr, err := p.unary()
        if err != nil {
            return nil, err
        }

        return tree.NewUnary(p.advance(), expr), nil
    }

    return p.primary()
}

func (p *Parser) primary() (tree.Expr, error){
    if p.match(token.TRUE) {
        return tree.NewLiteral(true), nil
    }

    if p.match(token.FALSE) {
        return tree.NewLiteral(false), nil
    }

    if p.match(token.STRING, token.NUMBER) {
        return tree.NewLiteral(p.previous().Literal()), nil
    }

    if p.match(token.IDENTIFIER) {
        return tree.NewVariable(p.previous()), nil
    }

    if p.match(token.LEFT_PAREN) {
        expr, err := p.expression()

        if err != nil {
            return nil, err
       }

        p.consume(token.RIGHT_PAREN, "Expected ')' after expression.")

        return tree.NewGrouping(expr), nil
    }

    return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) consume(tokenType token.TokenType, message string) (token.Token, error) {
    if p.check(tokenType) {
        return p.advance(), nil
    }

    err := p.error(p.peek(), message)

    return token.Token{}, err
}

func (p *Parser) synchronize() {
    p.advance()

    for !p.isAtEnd() {
        if p.previous().Type() == token.SEMICOLON {
            return
        }

        switch p.previous().Type() {
        case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
            return
        }

        p.advance()
    }
}

func (p *Parser) error(t token.Token, message string) error {
    errors.Error(t, message)
    return fmt.Errorf("")
}

func (p *Parser) match(types ...token.TokenType) bool {
    for _, t := range types {
        if p.check(t) {
            p.advance()
            return true
        }
    }

    return false
}

func (p Parser) check(t token.TokenType) bool {
    if p.isAtEnd() {
        return false
    }

    return p.tokens[p.current].Type() == t
}

func (p *Parser) advance() token.Token {
    if !p.isAtEnd() {
        p.current++
    }

    return p.previous()
}

func (p Parser) isAtEnd() bool {
    return p.peek().Type() == token.EOF
}

func (p Parser) peek() token.Token {
    return p.tokens[p.current]
}

func (p Parser) previous() token.Token {
    return p.tokens[p.current - 1]
}
